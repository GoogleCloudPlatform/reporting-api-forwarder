# Reporting API forwarder

This direcotry contains a web server that receives reports via Reporting API and forward them to [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) in [OTLP](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md).

Because both Reporting API and OTLP are vendor agnostic specification and tools, you are able to utilize data emit from this forwarder with arbitrary backends that has the integration with OpenTelemetry Collector.

In this sample, the default configuration is with [Google Cloud Monitoring](https://cloud.google.com/monitoring).

## certification files

The forwarder expects to run with TLS with a custom domain (e.g. `endpoint.example.com`), and you need to put certification file and key file in `cert` directory.
