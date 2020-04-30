# Introduction 
Project is transforming data from Apigee metrics endpoint into prometheus-friendly format. Currently supported metrics
are 95th percentile latency and total traffic counts divided into separate environment, region, status code and proxy.

# Getting Started
Installation: provided manifest that pulls image from CR and enables scraping with prometheus k8s auto-discovery.
Requires `apigee-monitoring-credentials` secret with username and password fields. This is only used to obtain OAuth
token, after that refresh token is used.

# Build and Test
Use provided dockerfile for building the image. When using kuberenetes deployment, specify image registry in manifest.
Update hardcoded organization value to match your Apigee org in pkg/metrics/metrics.go

# Contribute

Adding new metrics should be straightforward, see `metrics.go` for example of obtaining metrics from Apigee
and `exporter.go` on configuring [Prometheus metric](https://prometheus.io/docs/concepts/metric_types/)
