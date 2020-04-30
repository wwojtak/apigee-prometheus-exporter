package exporter

import (
	"apigee-prometheus-exporter/pkg/metrics"
	"apigee-prometheus-exporter/pkg/token"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	apigeeToken       = new(token.ApigeeToken)
	httpRequestsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "apigee",
			Name:      "http_requests_total",
			Help:      "Total number of processed requests",
		},
		[]string{"env", "proxy", "region", "status_code"},
	)
	httpLatencyPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "apigee",
			Name:      "http_total_latency_milliseconds",
			Help:      "95th percentile latency for requests to backend targets",
		},
		[]string{"env", "proxy", "region"},
	)
)

func itofloat(value interface{}) float64 {
	result, ok := value.(float64)
	if ok {
		return result
	}
	return -1
}

func SetupMetricsServer() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpLatencyPercentile)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			if time.Now().Unix() >= apigeeToken.Expires {
				var err error
				apigeeToken, err = token.ApigeeClient(apigeeToken)
				if err != nil {
					log.Fatal("Token request failed, see logs for details")
				}
			}

			// TODO: Rework logic so it's reusable
			tra := metrics.GetTrafficMetrics(apigeeToken)
			for s := range tra.Results[0].Series {
				t := tra.Results[0].Series[s]
				v := itofloat(t.Values[0][1])
				// Apigee sends metrics with proxy name "(not set)" for any invalid requests
				httpRequestsTotal.WithLabelValues(
					t.Tags.Env,
					t.Tags.Proxy,
					t.Tags.Region,
					t.Tags.StatusCode,
				).Set(v)
			}

			lat := metrics.GetLatencyMetrics(apigeeToken)
			for s := range lat.Results[0].Series {
				t := lat.Results[0].Series[s]
				// TODO: Figure out how to handle "null" from Apigee
				v := itofloat(t.Values[0][1])
				if v > 0 {
					httpLatencyPercentile.WithLabelValues(
						t.Tags.Env,
						t.Tags.Proxy,
						t.Tags.Region,
					).Set(v)
				}
			}

			time.Sleep(time.Minute * 5)
			// Delete all the metrics before next GET
			httpRequestsTotal.Reset()
			httpLatencyPercentile.Reset()
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
