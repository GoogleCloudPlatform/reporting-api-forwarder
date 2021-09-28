# Reporting API forwarder

This is an unoffical Google project to demonstrate the [Reporting API](https://www.w3.org/TR/reporting/)'s endpoint server.
This server works as a web server that receives all data from Reporting API and forwards them to monitoring tools such as [Google Cloud Monitoring](https://cloud.google.com/monitoring).
Also, for local testing purpose, you can try [Prometheus](https://prometheus.io/) with some additional configurations.

## How to try this sample

### Prerequisites

* [Docker Engine](https://docs.docker.com/engine/install/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### With Cloud Monitoring
#### Prepare report endpoint server

This repository contains two subdirectories:

* `forwarder`: the root of source code that contains Go web server
* `collector`: the config file for OpenTelemetry Collector
* `prometheus`: (optional) the config file for Prometheus

To confirm the details of each subdirectories, please find and read `README` files for detailed instructions.

To run this sample, you need to set up Google Cloud default credentials. Run the following commands for this demo.

```console
gcloud auth application-default login
chmod a+r ~/.config/gcloud/application_default_credentials.json
```

After setting up [Docker Compose](https://docs.docker.com/compose/) (detailed in [README](./collector/README.md)), launch the Reporting API endpoint system with the following command:

```console
$ docker-compose up
```

Then you will see the logs like the followings:

```console
 docker-compose up
WARNING: Found orphan containers (reporting-api-forwarder_prometheus_1, reporting-api-forwarder_reporter_1) for this project. If you removed or renamed this service in your compose file, you can run this command with the --remove-orphans flag to clean it up.
Starting reporting-api-forwarder_collector_1 ... done
Starting reporting-api-forwarder_forwarder_1 ... done
Attaching to reporting-api-forwarder_forwarder_1, reporting-api-forwarder_collector_1
forwarder_1  | {"severity":"INFO","timestamp":"2021-09-28T07:36:23.080139909Z","message":"Starting Reporting API forwarder: version 0.1.0.dev"}
forwarder_1  | ⇨ https server started on [::]:30443
collector_1  | 2021-09-28T07:36:23.473Z info    service/collector.go:303        Starting otelcontribcol...      {"Version": "v0.34.0", "NumCPU": 12}
collector_1  | 2021-09-28T07:36:23.474Z info    service/collector.go:242        Loading configuration...
...(omit)...
collector_1  | 2021-09-28T07:36:23.480Z info    service/collector.go:218        Everything is ready. Begin running and processing data.
```

As you can confirm in `docker-compose.yaml`, this demo uses the following ports:

|**Port in localhost**|**Service using the port**|**Port in services**|**Purpose**|
|30443|forwarder|30443|HTTPS server endpoint|
|4317|collector|4317|OTLP gRPC endpoint|
|4318|collector|4318|OTLP HTTP endpoint|
|9990|collector|9990|Prometheus exporter endpoint|
|8888|collector|8888|Colletor's metrics exporting endpoint|
|9090|prometheus|9090|Prometheus UI's endpoint|

**Note:** The instruction to change file permission of the credntial is just for this demonstration purpose, so after running this demo, please overwrite the permission as it was.

```
chmod go-r ~/.config/gcloud/application_default_credentials.json
```

#### Remix glitch samples

After creating [Glitch](https://glitch.com) account, remix the following two projects.

1. https://glitch.com/edit/#!/reporting-api-demo
2. https://glitch.com/edit/#!/intervention-generator

`reporting-api-demo` is the report generator, and `intervention-generator` is the source to generate intervention error that is called via the iframe DOM in the `reporting-api-demo`.
For the name of remix projects, the recommendation is to add your own suffix to both like `reporting-api-demo-<suffix>`, such as `reporting-api-demo-yoshifumi` and `intervention-generator-yoshifumi`

Once you remixed the two project, you need to modify the value of endpoint URLs.

In `reporting-api-demo-<suffix>`, open `server.js`, find the following lines and fix the value.

```javascript
const REPORTING_ENDPOINT_BASE = "https://<the URL where you run the endpoint server>"
...(omit)...
const INTERVENTION_GENERATOR_URL = "https://intervention-generator-<suffix>.glitch.me/";
```

In `intervention-generator-<suffix>`, open `server.js`, find the following lines and fix the value.

```javascript
const REPORTING_ENDPOINT = "https://<the URL where you run the endpoint server>/default";
```

**NOTE:** This demo requires `OPTION` HTTP method to work, so some handy solutions such as ngrok can't be an option. [This article](https://web.dev/how-to-use-local-https/) explains how to setup local environment for HTTPS connection from glitch to the forwarder.

#### Confirm data reception

If the endpoint system is working correctly, you will see the following outputs on your console:


```
collector_1  | 2021-09-05T05:08:02.622Z DEBUG   loggingexporter/logging_exporter.go:66  ResourceMetrics #0
collector_1  | Resource labels:
collector_1  |      -> service.name: STRING(unknown_service:server)
collector_1  |      -> telemetry.sdk.language: STRING(go)
collector_1  |      -> telemetry.sdk.name: STRING(opentelemetry)
collector_1  |      -> telemetry.sdk.version: STRING(1.0.0-RC3)
collector_1  | InstrumentationLibraryMetrics #0
collector_1  | InstrumentationLibrary reporting-api-forwarder 0.1.0.dev
collector_1  | Metric #0
collector_1  | Descriptor:
collector_1  |      -> Name: report.count
collector_1  |      -> Description: number of reports
collector_1  |      -> Unit: call
collector_1  |      -> DataType: Sum
collector_1  |      -> IsMonotonic: true
collector_1  |      -> AggregationTemporality: AGGREGATION_TEMPORALITY_CUMULATIVE
collector_1  | NumberDataPoints #0
collector_1  | Data point attributes:
collector_1  |      -> column-number: INT(24)
collector_1  |      -> dispotision: STRING(enforce)
collector_1  |      -> generated: INT(1630818474100)
collector_1  |      -> line-number: INT(14)
collector_1  |      -> message: STRING(Permissions policy violation: microphone is not allowed in this document.)
collector_1  |      -> policy-id: STRING(microphone)
collector_1  |      -> source-file: STRING(https://reporting-api-demo-otel.glitch.me/script.js)
collector_1  |      -> type: STRING(permissions-policy-violation)
collector_1  |      -> url: STRING(https://reporting-api-demo-otel.glitch.me/page)
collector_1  |      -> useragent: STRING(Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4628.3 Safari/537.36)
collector_1  | StartTimestamp: 2021-09-05 04:49:42.432981665 +0000 UTC
collector_1  | Timestamp: 2021-09-05 05:08:02.440571589 +0000 UTC
collector_1  | Value: 1
collector_1  | NumberDataPoints #1
collector_1  | Data point attributes:
collector_1  |      -> column-number: INT(12)
...(omit)...
```

This is because the default configuration of OpenTelemetry Collector is set up with [loggingexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/loggingexporter) as well as [googlecloudexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/googlecloudexporter).

This means, your Reporting API endpoint system is working properly and it should be sending the corresponding report data to Google Cloud Monitoring.


#### Confirm Google Cloud Monitoring

Open your Google Cloud Console and navigate yourself to Metrics Explorer in Google Cloud Monitoring. Then, try the following values in the forms:

* Resource Type: `global`
* Metric: filter `reporting-api/count` -> select `custom.googleapis.com/opencensus/reporting-api/count`
* Group by: `type`
* Aligner: `delta`

![Metrics Explorer](./static/image/metrics-explorer-1.png "the form in Metrics Explorer")

![Metrics Explorer](./static/image/metrics-explorer-2.png "the chart in Metrics Explorer")


### With Prometheus

This sample also provides the configuration to play with [Prometheus](https://prometheus.io/).

First, you need to configure `docker-compose.yaml`. You can find the section for `prometheus` commmented out in the file. Just uncomment the section and you will see the followings:

```yaml
...(omit)...
## Uncomment `prometheus` section if you want to try locally
  prometheus:
    image: prom/prometheus:v2.30.0
    volumes:
    - ./prometheus/config.yaml:/etc/prometheus/config.yaml
    entrypoint: ['prometheus', '--config.file', '/etc/prometheus/config.yaml']
    ports:
    - "9090:9090"
```

Second, you need to configure OpenTelemetry Collector to expose metrics in the Prometheus format via HTTP. You also just need to uncomment the section for Prometheus in `collector/config.yaml`. Please refer to [README.md](./collector/README.md) for the details.

After applying the changes, you can run `docker-compose` and you will see the logs from Docker Compose.

```console
$ docker-compose up
```

If the configuration is working correctly and you are getting the reports from the browser, you can access [the query page (http://localhost:9090)](http://localhost:9090/graph?g0.expr=sum%20by%20(type)%20(reporting_api_demo_reporting_api_count)&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=5m) and you will see the line chart like below:

![Prometheus UI](./static/image/prometheus-1.png "the chart on Prometheus UI")
