# Reporting API forwarder

This direcotry contains a web server that receives reports via Reporting API and forwards them to [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) in [OTLP](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md).

Because both Reporting API and OTLP are vendor-agnostic specifications and tools, you are able to utilize data emit from this forwarder with arbitrary backends that has the integration with OpenTelemetry Collector.

In this sample, the default configuration is with [Google Cloud Monitoring](https://cloud.google.com/monitoring).

## Server requirements

Reporting API expects some technical conditions:

* The endpoint receives reports via HTTP `POST` method
* `Content-Type` HTTP header is `application/reports+json`
* HTTP body is JSON (the root object is an array)
* Cookie enabled
* CORS with preflight

## TLS support

The forwarder supports TLS because CORS often requires the server side in TLS to avoid mixed content.

If you are to run the forwarder with TLS and to serve directly to the public with your custom domain (e.g. `endpoint.example.com`) or in localhost, you need to put certification file and key file in `cert` directory.

If you just run it behind a HTTPS proxy and don't need the forwarder to be HTTPS server, run it with `ENABLE_TLS=0` environment variable.
