package metrics

import (
	"apigee-prometheus-exporter/pkg/token"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Result struct {
	Results []Series `json:"results"`
}

type Series struct {
	Series []Metrics `json:"series"`
}

type Metrics struct {
	Name    string          `json:"name"`
	Tags    Tags            `json:"tags"`
	Columns []string        `json:"columns"`
	Values  [][]interface{} `json:"values"`
}

type Tags struct {
	Env        string `json:"env"`
	Org        string `json:"org"`
	Percentile string `json:"percentile"`
	Proxy      string `json:"proxy"`
	Region     string `json:"region"`
	StatusCode string `json:"statusCode"`
}

func GetTrafficMetrics(apigeeToken *token.ApigeeToken) Result {
	metricsURL := "https://apimonitoring.enterprise.apigee.com/metrics/traffic"
	var bearer = "Bearer " + apigeeToken.AccessToken
	req, err := http.NewRequest("GET", metricsURL, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	query := req.URL.Query()
	query.Set("org", "YOUR-ORG-HERE")
	query.Set("interval", "5m")
	query.Set("groupBy", "env,proxy,statusCode")
	query.Set("from", "-5m")
	query.Set("select", "count")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Error(err)
	}
	return result
}

func GetLatencyMetrics(apigeeToken *token.ApigeeToken) Result {
	metricsURL := "https://apimonitoring.enterprise.apigee.com/metrics/latency"
	var bearer = "Bearer " + apigeeToken.AccessToken
	req, err := http.NewRequest("GET", metricsURL, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	query := req.URL.Query()
	query.Set("org", "YOUR-ORG-HERE")
	query.Set("percentile", "95")
	query.Set("interval", "5m")
	query.Set("windowsize", "1m")
	query.Set("select", "totalLatency")
	query.Set("groupBy", "env,region,proxy")
	query.Set("from", "-5m")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Error(err)
	}
	return result
}
