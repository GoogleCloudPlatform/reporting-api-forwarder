module github.com/GoogleCloudPlatform/reporting-api-forwarder

go 1.16

require (
	github.com/labstack/echo/v4 v4.9.0
	github.com/rs/zerolog v1.25.0
	go.opentelemetry.io/otel v1.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.43.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.27.0
	go.opentelemetry.io/otel/metric v1.27.0
	go.opentelemetry.io/otel/sdk/metric v1.27.0
)
