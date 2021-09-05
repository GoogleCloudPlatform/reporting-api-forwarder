# OpenTelemetry Collector for Reporting API Forwarder

This demo uses [OpenTelemetry contrib collector](https://github.com/open-telemetry/opentelemetry-collector-contrib) that bundles [Google Cloud Exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/googlecloudexporter).

You can modify the sample `config.yaml` to meet your backend. In such a case, you will also need to find or build OpenTelemetry Collector that bundles the exporter to your backend monitoring products. For details, please refer to [this blog post](https://medium.com/opentelemetry/building-your-own-opentelemetry-collector-distribution-42337e994b63).
## Setup

The [README of Google Cloud Exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/googlecloudexporter) gives the detailed set up process. Please follow step 1, 2 and 3 in the README in advance of the following steps.

This section explains the specific setup process for this demo.

### Change the permission of the credential files

After obtaining the application default credentials in step 3 of the `googlecloudexporter` setting, you need to change the file permission so that Docker Compose can look up the file.

**NOTE: This operation is only for the demonstration purpose. Do not apply this to your production environment. Also, please remove the permission after trying this demo.**

```
chmod u+r ~/.config/gcloud/application_default_credentials.json
```

### Install Docker Compose

_You can ignore this section if you already have Docker Compose on your environment._

This demo requires [Docker Compose](https://docs.docker.com/compose/). If you don't have Docker Compose on your environment, please follow [the official installation guide](https://docs.docker.com/compose/install/).


### Run Docker Compose

After the setup, try running Docker Compose. You will see something like the followings:

```
$ cd ${PROJECT_ROOT}
$ docker-compose up
Starting reporting-api-forwarder_forwarder_1 ... done
Starting reporting-api-forwarder_collector_1 ... done
Attaching to reporting-api-forwarder_forwarder_1, reporting-api-forwarder_collector_1
forwarder_1  | {"severity":"INFO","timestamp":"2021-09-05T04:48:20.64896508Z","message":"Starting Reporting API forwarder: version 0.1.0.dev"}
forwarder_1  | ⇨ https server started on [::]:30443
collector_1  | 2021-09-05T04:48:20.942Z info    service/collector.go:303        Starting otelcontribcol...  {"Version": "v0.34.0", "NumCPU": 12}
collector_1  | 2021-09-05T04:48:20.943Z info    service/collector.go:242        Loading configuration...
...(omit)...
collector_1  | 2021-09-05T04:48:20.950Z info    service/telemetry.go:95 Serving Prometheus metrics      {"address": ":8888", "level": 0, "service.instance.id": "5e407fc5-b3dd-4255-941d-e76f0dbe9dd4"}
collector_1  | 2021-09-05T04:48:20.950Z info    service/collector.go:218        Everything is ready. Begin running and processing data.
```

### Integration check

See the rest in [README](../README.md) in the project root.
