module github.com/GoogleCloudPlatform/reporting-api-forwarder

go 1.16

require (
	github.com/labstack/echo/v4 v4.6.1
	github.com/rs/zerolog v1.25.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.24.0
	go.opentelemetry.io/otel/metric v0.24.0
	go.opentelemetry.io/otel/sdk/metric v0.24.0
)
