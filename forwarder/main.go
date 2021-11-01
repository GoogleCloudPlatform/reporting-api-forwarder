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
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	// Using echo for easy CORS support
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
	instrumentationVersion = "0.1.0"
	instrumentationName    = "reporting-api-forwarder"
)

var (
	collectorAddr string
	certFile      string
	keyFile       string
	enableTLS     bool = true

	// meter is the global meter to send Reporting API relevant metrics
	meter metric.Meter

	// reportCounter is the counter of reports
	reportCounter metric.Int64Counter
)

func init() {
	meter := metric.Must(global.Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(instrumentationVersion),
	))
	reportCounter = meter.NewInt64Counter(
		"reporting-api/count",
		metric.WithDescription("number of reports"),
		metric.WithUnit("call"),
	)
	collectorAddr = os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "127.0.0.1:4317"
	}
	certFile = os.Getenv("CERT_FILE")
	if certFile == "" {
		certFile = "cert/cert.pem"
	}
	keyFile = os.Getenv("KEY_FILE")
	if keyFile == "" {
		keyFile = "cert/key.pem"
	}
	if os.Getenv("ENABLE_TLS") == "0" {
		enableTLS = false
	}
}

// installPipeline sets up the initial pipeline for exporting metrics via OpenTelemetry.
func installPipeline(ctx context.Context) func() {
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(collectorAddr),
	)
	exporter, err := otlpmetric.New(ctx, client)
	if err != nil {
		logger.Fatal().Msgf("failed to create stdout metric exporter: %v", err)
	}

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithExactDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
	)
	global.SetMeterProvider(pusher)
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
	logger.Info().Msgf("Collector endpoint: %s", collectorAddr)

	ctx := context.Background()
	shutdown := installPipeline(ctx)
	defer shutdown()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	// In order to hand Reporting API, the reporting endpoint needs to handle CORS
	// TODO(yoshifumi): consider adding the config to set allow origin externally
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/", rootHandler)
	e.POST("/main", mainHandler)
	e.POST("/default", defaultHandler)
	e.GET("/healthz", healthzHandler)

	if enableTLS {
		if err := e.StartTLS(":30443", certFile, keyFile); err != nil {
			logger.Fatal().Msgf("failure occured during HTTP server launch process: %v."+
				"Check docker cache if you added files already and running this app in docker.", err)
		}
	} else {
		if err := e.Start(":8080"); err != nil {
			logger.Fatal().Msgf("failure occured on launching HTTP server: %v", err)
		}
	}
}

func rootHandler(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("%v: The reporting endpoint is /default", time.Now()))
}

func mainHandler(c echo.Context) error {
	return handleReportRequest(c)
}

func defaultHandler(c echo.Context) error {
	return handleReportRequest(c)
}

func healthzHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
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
	if err := r.Body.Close(); err != nil {
		return err
	}
	logger.Info().Msgf("accepted %v bytes report", len(data))

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
	kvs := []attribute.KeyValue{
		attribute.String("type", r.Typ),
		attribute.String("url", r.URL),
		attribute.String("useragent", r.UserAgent),
		attribute.Int64("generated", now-int64(r.Age)),
	}
	body := r.Body
	switch r.Typ {
	case "csp-violation":
		kvs = append(kvs,
			attribute.String("blocked-url", body.BlockedURL),
			attribute.String("dispotision", body.Disposition),
			attribute.String("document-url", body.DocumentURL),
			attribute.String("effective-directive", body.EffectiveDirective),
			attribute.String("original-policy", body.OriginalPolicy),
			attribute.String("referrer", body.Referrer),
			attribute.String("sample", body.Sample),
			attribute.Int("status-code", body.StatusCode),
		)
	case "deprecation":
		kvs = append(kvs,
			attribute.Int("column-number", body.ColumnNumber),
			attribute.Int("line-number", body.LineNumber),
			attribute.String("id", body.ID),
			attribute.String("message", body.Message),
			attribute.String("source-file", body.SourceFile),
		)
	case "permissions-policy-violation", "document-policy-violation":
		kvs = append(kvs,
			attribute.Int("column-number", body.ColumnNumber),
			attribute.Int("line-number", body.LineNumber),
			attribute.String("id", body.ID),
			attribute.String("message", body.Message),
			attribute.String("policy-id", body.PolicyID),
			attribute.String("source-file", body.SourceFile),
		)
	case "coep":
		kvs = append(kvs,
			attribute.String("blocked-url", body.BlockedURL),
			attribute.String("dispotision", body.Disposition),
			attribute.String("destination", body.Destination),
		)
	case "intervention":
		kvs = append(kvs,
			attribute.Int("column-number", body.ColumnNumber),
			attribute.Int("line-number", body.LineNumber),
			attribute.String("id", body.ID),
			attribute.String("message", body.Message),
			attribute.String("source-file", body.SourceFile),
		)
	default:
	}
	return kvs
}
