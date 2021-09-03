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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

const (
	instrumentationVersion = "0.1.0.dev"
	instrumentationName    = "github.com/GoogleCloudPlatform/reporting-api-forwarder"
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
	reportCounter, err = meter.NewInt64Counter("request.count")
	if err != nil {
		panic(err)
	}
}

func installPipeline(ctx context.Context) func() {
	exporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		logger.Fatal().Msgf("failed to create stdoutmetric exporter: %v", err)
	}

	pusher := controller.New(
		processor.New(
			simple.NewWithInexpensiveDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
	)
	if err = pusher.Start(ctx); err != nil {
		logger.Fatal().Msgf("failed to start push controller: %v", err)
	}
	global.SetMeterProvider(pusher.MeterProvider())

	return func() {
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
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/", rootHandler)
	e.POST("/main", mainHandler)
	e.POST("/default", defaultHandler)
	if err := e.Start(":3000"); err != nil {
		logger.Fatal().Msgf("failure occured during HTTP server launch process: %v", err)
	}
}

func rootHandler(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()

	reportCounter.Add(ctx, 1)
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
	r := c.Request()
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

	var buf []map[string]interface{}
	err = json.Unmarshal(data, &buf)
	if err != nil {
		logger.Error().Msgf("error on parsing JSON: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info().Msgf("%v", buf)
	return c.String(http.StatusOK, "OK")
}
