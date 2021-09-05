// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

const (
	instrumentationVersion = "0.1.0.dev"
	instrumentationName    = "reporting-api-forwarder"
	collectorHost          = "collector" // refer to docker-compose.yaml
	collectorPort          = "4317"
)

var (
	meter metric.Meter

	reportCounter metric.Int64Counter
)

func init() {
	meter = global.GetMeterProvider().Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(instrumentationVersion),
	)
	var err error
	reportCounter, err = meter.NewInt64Counter(
		"report.count",
		metric.WithDescription("number of reports"),
		metric.WithUnit("call"),
	)
	if err != nil {
		panic(err)
	}
}

func installPipeline(ctx context.Context) func() {
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(collectorHost+":"+collectorPort),
	)
	exporter, err := otlpmetric.New(ctx, client)
	if err != nil {
		logger.Fatal().Msgf("failed to create stdoutmetric exporter: %v", err)
	}

	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
	)
	global.SetMeterProvider(pusher.MeterProvider())
	if err = pusher.Start(ctx); err != nil {
		logger.Fatal().Msgf("failed to start push controller: %v", err)
	}

	return func() {
		if err := exporter.Shutdown(ctx); err != nil {
			logger.Fatal().Msgf("failed to stop OTLP metric client: %v", err)
		}
		if err := pusher.Stop(ctx); err != nil {
			logger.Fatal().Msgf("failed to stop push controller: %v", err)
		}
	}
}

func main() {
	logger.Info().Msgf("Starting Reporting API forwarder: version %s", instrumentationVersion)

	ctx := context.Background()
	shutdown := installPipeline(ctx)
	defer shutdown()

	e := echo.New()
	e.HideBanner = true
	// In order to hand Reporting API, the reporting endpoint needs to handle CORS
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/", rootHandler)
	e.POST("/main", mainHandler)
	e.POST("/default", defaultHandler)
	if err := e.StartTLS(":30443", "cert/cert.pem", "cert/key.pem"); err != nil {
		logger.Fatal().Msgf("failure occured during HTTP server launch process: %v", err)
	}
}

func rootHandler(c echo.Context) error {
	r := c.Request()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error().Msgf("failed to read request body: %v", err)
		return err
	}
	defer r.Body.Close()
	logger.Info().RawJSON("report", data)
	return c.String(http.StatusOK, string(data))
}

func mainHandler(c echo.Context) error {
	return handleReportRequest(c)
}

func defaultHandler(c echo.Context) error {
	return handleReportRequest(c)
}

func handleReportRequest(c echo.Context) error {
	now := time.Now()
	r := c.Request()
	ctx := r.Context()

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/reports+json" {
		logger.Error().Msgf("Content-Type header is not application/reports+json: %v", r.Header)
		return c.String(http.StatusBadRequest, "Content-Type not supported. The Content-Type must be application/reports+json.")
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error().Msgf("error on reading data: %v", err)
	}
	defer r.Body.Close()

	var buf []report
	err = json.Unmarshal(data, &buf)
	if err != nil {
		logger.Error().Msgf("error on parsing JSON: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	for _, r := range buf {
		kvs := extractKeyValues(r, now)
		reportCounter.Add(ctx,
			1,
			kvs...,
		)
	}
	return c.String(http.StatusOK, "OK")
}

// extractKeyValues returns a slice of KeyValues based on the metadata
// embedded in body in the report.
func extractKeyValues(r report, t time.Time) []attribute.KeyValue {
	// NOTE: from Go1.17, use time.Time#UnixMilli
	now := t.UnixNano() / int64(time.Millisecond)
	kvs := make([]attribute.KeyValue, 4)
	kvs[0] = attribute.String("type", r.Typ)
	kvs[1] = attribute.String("url", r.URL)
	kvs[2] = attribute.String("useragent", r.UserAgent)
	kvs[3] = attribute.Int64("generated", now-int64(r.Age))
	body := r.Body
	switch r.Typ {
	case "csp-violation":
		kvs = append(kvs,
			attribute.String("blocked-url", body["blockedURL"].(string)),
			attribute.String("dispotision", body["disposition"].(string)),
			attribute.String("document-url", body["documentURL"].(string)),
			attribute.String("effective-directive", body["effectiveDirective"].(string)),
			attribute.String("original-policy", body["originalPolicy"].(string)),
			attribute.String("referrer", body["referrer"].(string)),
			attribute.String("sample", body["sample"].(string)),
			attribute.Int("status-code", int(body["statusCode"].(float64))),
		)
	case "deprecation":
		kvs = append(kvs,
			attribute.Int("column-number", int(body["columnNumber"].(float64))),
			attribute.Int("line-number", int(body["lineNumber"].(float64))),
			attribute.String("id", body["id"].(string)),
			attribute.String("message", body["message"].(string)),
			attribute.String("source-file", body["sourceFile"].(string)),
		)
	case "permissions-policy-violation":
		kvs = append(kvs,
			attribute.Int("column-number", int(body["columnNumber"].(float64))),
			attribute.Int("line-number", int(body["lineNumber"].(float64))),
			attribute.String("dispotision", body["disposition"].(string)),
			attribute.String("message", body["message"].(string)),
			attribute.String("policy-id", body["policyId"].(string)),
			attribute.String("source-file", body["sourceFile"].(string)),
		)
	case "document-policy-violation":
		kvs = append(kvs,
			attribute.Int("column-number", int(body["columnNumber"].(float64))),
			attribute.Int("line-number", int(body["lineNumber"].(float64))),
			attribute.String("dispotision", body["disposition"].(string)),
			attribute.String("message", body["message"].(string)),
			attribute.String("policy-id", body["policyId"].(string)),
			attribute.String("source-file", body["sourceFile"].(string)),
		)
	case "coep":
		kvs = append(kvs,
			attribute.String("blocked-url", body["blockedURL"].(string)),
			attribute.String("dispotision", body["disposition"].(string)),
			attribute.String("destination", body["destination"].(string)),
		)
	case "intervention":
		kvs = append(kvs,
			attribute.Int("column-number", int(body["columnNumber"].(float64))),
			attribute.Int("line-number", int(body["lineNumber"].(float64))),
			attribute.String("id", body["id"].(string)),
			attribute.String("message", body["message"].(string)),
			attribute.String("source-file", body["sourceFile"].(string)),
		)
	default:
	}
	return kvs
}
