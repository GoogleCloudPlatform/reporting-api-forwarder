# Reporting API forwarder

This is an unoffical Google project to demonstrate the [Reporting API](https://www.w3.org/TR/reporting/)'s endpoint server.
This server works as a web server that receives all data from Reporting API and forwards them to monitoring tools such as [Google Cloud Monitoring](https://cloud.google.com/monitoring).
Also, for local testing purpose, you can try [Prometheus](https://prometheus.io/) with some additional configurations.


## Components

This repository contains the following subdirectories:

* `forwarder`: the main component that receives reports and forwards them to OpenTelemetry Collector
* `collector`: the config file for OpenTelemetry Collector
* `prometheus`: (optional) the config file for Prometheus
* `grafana`: (optional) the config file for Grafana

To confirm the details of each subdirectories, please find and read `README` files for detailed instructions.

The interaction of each components is described as below:

![Diagram](./static/image/overall-diagram.png "diagram")
## How to try this sample

### Prerequisites

* [Docker Engine](https://docs.docker.com/engine/install/)
  * If you run the forwarder in Docker or to run `docker compose`

* [Docker Compose](https://docs.docker.com/compose/install/)
  * For the cases "With Cloud Monitorig" and "with Prometheus"

* Google Cloud project and [`gcloud`](https://cloud.google.com/sdk/docs/install) command
  * For the case "With Cloud Monitoring" and "Running on Google Kubernetes Engine"

* [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl) and [`skaffold`](https://skaffold.dev/docs/install/)
  * For the case "Running on Google Kubernetes Engine"

`kubectl` and `skaffold` are also available via `gcloud components install` command. See [the document](https://cloud.google.com/sdk/docs/components) on how to install those command via `gcloud`.

### With Cloud Monitoring

#### Prepare Google Cloud project

In order to run this sample, you need to set up a Google Cloud project.

First, you need to create default credentials. Run the following commands for this demo.

```console
gcloud auth application-default login
chmod a+r ~/.config/gcloud/application_default_credentials.json
```
> **Note:** The instruction to change file permission of the credntial is just for this demonstration purpose, so after running this demo, please overwrite the permission as it was.
>
> ```console
> chmod go-r ~/.config/gcloud/application_default_credentials.json
> ```

Then you need a project where you send the reports to.

```console
gcloud projects create reporting-api-demo
```

Once the command is successfully done, enable "Cloud Monitoring API" in the project.

```console
gcloud services eanble monitoring
```

#### Prepare forwarder

After setting up Google Cloud project and credentials, you need to run the fowarder (reporting API endpoint).
Specify TLS certification files with the environment variables `CERT_FILE` and `KEY_FILE`. The default value is set to `cert/cert.pem` and `cert/key.pem`.

Then you need to add TLS certification files to `forwarder/cert` directory. Add cert file and key file there with the name `cert.pem` and `key.pem`.

After setting up [Docker Compose](https://docs.docker.com/compose/) (detailed in [README](./collector/README.md)), launch the Reporting API endpoint system with the following command:

```console
$ PROJECT_ID=$(gcloud config get-value project) docker-compose up
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

|Port in localhost|Service using the port|Port in services|Purpose|
|-----------------|----------------------|----------------|-------|
|30443|forwarder|30443|HTTPS server endpoint|
|4317|collector|4317|OTLP gRPC endpoint|
|4318|collector|4318|OTLP HTTP endpoint|
|9990|collector|9990|Prometheus exporter endpoint|
|8888|collector|8888|Colletor's metrics exporting endpoint|
|9090|prometheus|9090|Prometheus UI's endpoint|

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

> **NOTE:** This demo requires `OPTION` HTTP method to work, so some handy solutions such as ngrok can't be an option. [This article](https://web.dev/how-to-use-local-https/) explains how to setup local environment for HTTPS connection from glitch to the forwarder.

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


### With Prometheus and Grafana

This sample also provides the configuration to play with [Prometheus](https://prometheus.io/) and [Grafana](https://grafana.com/)

You can run `docker-compose` with the configuration file for Prometheus and you will see the logs from Docker Compose.

```console
$ docker-compose -f docker-compose-prometheus.yaml up
```

This docker-compose uses the following ports:

|Port in localhost|Service using the port|Port in services|Purpose|
|-----------------|----------------------|----------------|-------|
|30443|forwarder|30443|HTTPS server endpoint|
|4317|collector|4317|OTLP gRPC endpoint|
|4318|collector|4318|OTLP HTTP endpoint|
|9990|collector|9990|Prometheus exporter endpoint|
|8888|collector|8888|Colletor's metrics exporting endpoint|
|9090|prometheus|9090|Prometheus UI's endpoint|


If the configuration is working correctly and you are getting the reports from the browser, you can access [the query page (http://localhost:9090)](http://localhost:9090/graph?g0.expr=sum%20by%20(type)%20(reporting_api_demo_reporting_api_count)&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=5m) and you will see the line chart like below:

![Prometheus UI](./static/image/prometheus-1.png "the chart on Prometheus UI")

Also you can access to [the pre-defined dashboard of Grafana](http://localhost:3000/d/uhbXu_H7z/reporting-api?orgId=1).


### Running on Google Kubernetes Engine

This repository also contains the way to run the demo on [Google Kubernetes Engine (GKE)](https://cloud.google.com/kubernetes-engine) to share the demo in public with the team. Note that this is just a demo, so **do not use this demo as-is for production.**
#### Prerequisites

To deploy this demo onto a GKE Cluster, you need to have the following commands in your environment.

* `gcloud`
* [Docker Engine](https://docs.docker.com/engine/install/)
* [`kubectl`](https://kubernetes.io/docs/reference/kubectl/overview/)
* [`skaffold`](https://skaffold.dev/)

Make sure to install these commands in advance. Also, the demo needs to have **a custom domain** (eg. `reporting.example.com`). You will need to add a DNS record for it later on.

#### Set up Google Cloud

In order to deploy and run the demo on GKE, you need to have a couple of set up in addition:

* Prepare a custom domain for the demo
* Enabling APIs for GCE, GKE, Cloud Monitoring, Cloud Logging, Container Registry
* Create a GKE cluster
* Reserve static IP address
* Initialize `kubectl` and `skaffold` command

##### Prepare a custom domain

You can use any registerer to obtain the custom domain for the demo if you don't have any. If you already have one, you may want to add a subdomain for the demo. The external IP address will be reserved later.
##### Enabling Google Cloud APIs

```console
gcloud services enable compute.googleapis.com
gcloud services enable container.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable logging.googleapis.com
gcloud services enable monitoring.googleapis.com
```

##### Create a GKE cluster

You can create a GKE cluster via [Cloud Console](https://console.cloud.google.com/kubernetes) or `gcloud` command. Please refer to [the official document](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-zonal-cluster). Here is the sample command to create a cluster for this demo.

```console
$ gcloud container clusters create reporting-api-demo-cluster \
  --release-channel=stable \
  --zone=us-central1-b \
  --machine-type=e2-standard-2 \
  --num-nodes=2

Creating cluster reporting-api-demo-cluster in us-central1-b...done.
Created [https://container.googleapis.com/v1/projects/<project_id>/zones/us-central1-b/clusters/reporting-api-demo-cluster].
To inspect the contents of your cluster, go to: https://console.cloud.google.com/kubernetes/workload_/gcloud/us-central1-b/reporting-api-demo-cluster?project=<project_id>
kubeconfig entry generated for reporting-api-demo-cluster.

NAME                        LOCATION       MASTER_VERSION    MASTER_IP       MACHINE_TYPE   NODE_VERSION      NUM_NODES  STATUS
reporting-api-demo-cluster  us-central1-b  1.19.13-gke.1200  34.133.154.139  e2-standard-2  1.19.13-gke.1200  2          RUNNING
```

Note that `gcloud container clusters create` command automatically generates the necessary kubeconfig entry for `kubectl` to work with the cluster.
##### Reserve static IP address and set DNS record

In this demo, you will use "VPC Network" setting to obetain external static IP address. See the details in [this document.](https://cloud.google.com/compute/docs/ip-addresses/reserve-static-external-ip-address)

You can reserve a static external IP with the following command:

```console
$ gcloud compute addresses create reporting-api \
  --global \
  --ip-version IPV4

Created [https://www.googleapis.com/compute/v1/projects/<project_id>/global/addresses/reporting-api].
```

After seeing this successful message, you can confirm the reserved IP address.

```console
$ gcloud compute addresses describe reportin-api \
  --global \
  --format="value(address)"

34.117.46.72
```

Note this IP address for DNS record setting. For example, in the case of [Google Domains](https://domains.google/), you can configure the DNS record as follows:

![DNS record configurations](./static/image/dns-record.png "Google Domains DNS record edit")

For testing purpose, it's better to set TTL shorter such as 60s.

##### Initialize `kubectl` command

So that `kubectl` (the Kubernetes administration tool) can communicate with the GKE cluster, you need to pass the credential to it. [Google Cloud's document](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl) explains how to configure this in detail. Though already kubeconfig entry is set as mentioned in "Create a GKE cluster" section, just run the command just in case.

```console
$ gcloud container clusters get-credentials reporting-api-demo-cluster
Fetching cluster endpoint and auth data.
kubeconfig entry generated for reporting-api-demo-cluster.
```

Now you can confirm if the Kubernetes context is properly set.

```console
$ kubectl config current-context
gke_<project_id>_us-central1-b_reporting-api-demo-cluster
```

##### Initialize `skaffold` command

Finally, you can configure `skaffold`. [`skaffold`](https://skaffold.dev/) is a handy tool that manages the development iterations of container images and Kuberenetes deployment.

For this demo, you need to configure the default container registry where `skaffold` push container images to and the GKE cluster pull those out of.

As described in [this document](https://skaffold.dev/docs/environment/image-registries/), set the default repositry.

```console
$ skaffold config set default-repo gcr.io/$(gcloud config get-value project)/reporting-api
```

By doing this, `skaffold` command parses the image path in the Kubernetes manifests, and tries to search the image from the repository by completing the default repository path. For example, if your manifest file is configured as follows:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: forwarder
spec:
  selector:
    matchLabels:
      app: forwarder
  template:
    metadata:
      labels:
        app: forwarder
    spec:
      terminationGracePeriodSeconds: 5
      containers:
      - name: forwarder
        image: forwarder
        ports:
        - containerPort: 8080
        ...(omit)...
```

then, the `spec.template.spec.containers[0].iamge` is recognized as `gcr.io/<project_id>/reporting-api/forwarder` instead of `forwarder`.

#### Edit `certificate.yaml`

Before running `skaffold`, you need to specify your custom domain name to `${PROJECT_ROOT}/gke/manifests/certificate.yaml` so that [Google-managed SSL certificates](https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs) will be available for your domain.

Open the `certificate.yaml` file and edit the line of `spec.domains[0]` to have your custom domain.

```yaml
apiVersion: networking.gke.io/v1
kind: ManagedCertificate
metadata:
  name: reporting-api-cert
spec:
  domains:
  - REPLACE_WITH_YOUR_COSTOM_DOMAIN # eg. reporting.example.com
```

#### Run `skaffold`

Now you are ready to run `skaffold` to deploy the Kubernetes cluster to GKE.

```console
$ cd ${PROJECT_ROOT}/gke
$ skaffold run --tail

Generating tags...
 - forwarder -> gcr.io/<project_id>/reporting-api/forwarder:7ba4376-dirty
 - collector -> gcr.io/<project_id>/reporting-api/collector:7ba4376-dirty
Checking cache...
 - forwarder: Not found. Building
 - collector: Not found. Building
Starting build...
Building [forwarder]...
Sending build context to Docker daemon  17.32MB
Step 1/11 : FROM golang:1.17-bullseye as builder
 ---> 9f8b89ee4475
Step 2/11 : WORKDIR /src

...(omit)...

 - deployment.apps/collector created
 - service/collector created
Waiting for deployments to stabilize...
 - deployment/collector: creating container collector
    - pod/collector-5cbfb8dfd4-bbjgc: creating container collector
 - deployment/forwarder: creating container forwarder
    - pod/forwarder-b866cdf54-gh6np: creating container forwarder
 - deployment/forwarder is ready. [1/2 deployment(s) still pending]
 - deployment/collector is ready.
Deployments stabilized in 1 minute 1.717 second
You can also run [skaffold run --tail] to get the logs
```

One thing to note is that it takes 15-20 minutes for the certificate manager to take effect, so you may not able to access to your custom domain (eg. `https://reporting.example.com/`) via HTTPS.
You can confirm the status of the certificate provisioning by `kubectl describe managedcertificate` command.

As you see in `${PROJECT_ROOT}/gke/manifests/certificate.yaml`, your certificate is managed with the name `reporting-api-cert`.

```console
$ kubectl describe managedcertificate reporting-api-cert
Name:         reporting-api-cert
Namespace:    default
Labels:       app.kubernetes.io/managed-by=skaffold
              skaffold.dev/run-id=85a43ea1-345c-43d3-9676-9ddf398ccf81
Annotations:  <none>
API Version:  networking.gke.io/v1

...(omit)...

Spec:
  Domains:
    reporting.example.com
Status:
  Certificate Name:    mcrt-2476e85b-385e-4a93-b41c-77c668cd47ce
  Certificate Status:  Active
  Domain Status:
    Domain:     reporting.example.com
    Status:     Active
  Expire Time:  2022-01-12T15:50:28.000-08:00
Events:
  Type    Reason  Age   From                            Message
  ----    ------  ----  ----                            -------
  Normal  Create  17m   managed-certificate-controller  Create SslCertificate mcrt-2476e85b-385e-4a93-b41c-77c668cd47ce
```

Confirm the value of "Status > Doomain Status > Status" got "Active". (If it is any other values such as "Provisioning", the certificate is not ready yet.)
Then, you should be able to access `https://reporting.example.com/healthz` and get the result `"OK"` with the status code 200.

```console
$ curl -vvvv "https://reporting.example.com/healthz"
*   Trying 35.244.202.19:443...
* Connected to reporting.example.com (35.244.212.19) port 443 (#0)
* ALPN, offering h2
... (omit) ...
* Connection state changed (MAX_CONCURRENT_STREAMS == 100)!
< HTTP/2 200
< content-type: text/plain; charset=UTF-8
< vary: Origin
< date: Fri, 15 Oct 2021 11:53:08 GMT
< content-length: 2
< via: 1.1 google
< alt-svc: clear
<
* Connection #0 to host reporting.example.com left intact
OK%
```

### Standalone

If you already uses OpenTelemetry Collector and have the monitoring backend, you can run `forwarder` alone.

```console
cd forwarder
go build -o forwarder
COLLECTOR_ADDR=example.com:4317 ./forwarder
```

Make sure to put your certificate files (`cert.pem` and `key.pem`) in the `cert` file relative to the executor.
Otherwise, specify the path to both files with the environment variables `CERT_FILE` and `KEY_FILE`.


## Troubleshooting

### `PermissionDenied` error for Cloud Monitoring

If you see `code = PermissionDenied` error on sending metrics to Google Cloud Monitoring, like:

```console
collector_1   | 2021-09-29T04:03:11.441Z        info    exporterhelper/queued_retry.go:231      Exporting failed. Will retry the request after interval.        {"kind": "exporter", "name": "googlecloud", "error": "rpc error: code = PermissionDenied desc = Permission monitoring.metricDescriptors.create denied (or the resource may not exist).", "interval": "16.812982588s"}
```

You can check the following configurations:

* If the application default credentials is the one for the account you want to use
* If the Cloud Monitoring API is enabled in the project

### `forwarder` container image doesn't refresh

If you replace the certification files and observe no changes with `forwarder` container image, `docker-compose` should be running the cache. Try the following command and see how it goes.

```console
docker-compose build --no-cache
```

### Permission error of config files

If you face the following permission errors of reading config files, just add read access of those files to others: `chmod a+r path/to/config.file`.

```console
collector_1   | 2021-09-29T06:04:41.745Z        info    service/collector.go:242        Loading configuration...
collector_1   | Error: cannot load configuration's parser: error loading config file "/etc/otel/config.yaml": unable to read the file /etc/otel/config.yaml: open /etc/otel/config.yaml: permission denied
```