module github.com/GoogleCloudPlatform/reporting-api-forwarder

go 1.16

require (
	github.com/labstack/echo/v4 v4.5.0
	github.com/rs/zerolog v1.23.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.22.0
	go.opentelemetry.io/otel/metric v0.22.0
	go.opentelemetry.io/otel/sdk/metric v0.22.0
)
