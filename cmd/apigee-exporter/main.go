package main

import (
	"apigee-prometheus-exporter/pkg/exporter"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	exporter.SetupMetricsServer()

}
