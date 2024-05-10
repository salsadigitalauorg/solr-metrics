package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/salsadigitalauorg/lagoon-solr-metrics/pkg/lagoon"

	log "github.com/sirupsen/logrus"
)

var (
	solrHost         string
	solrPort         string
	solrCore         string
	insightsEndpoint string
)

// Init http client with default timeout.
var httpClient = &http.Client{Timeout: 10 * time.Second}

type Solr struct {
	Metrics map[string]any `json:"metrics"`
}

// Struct that only holds selected metrics.
type Metrics struct {
	Name      string `json:"CORE.coreName,omitempty"`
	StartTime string `json:"CORE.startTime,omitempty"`
	NumDocs   int    `json:"SEARCHER.searcher.numDocs,omitempty"`
	IndexSize string `json:"INDEX.size,omitempty"`
	Version   string `json:"CONTAINER.version.specification,omitempty"`
}

// A copy of Metrics struct but with flattened json keys.
type Info struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	NumDocs   int    `json:"numDocs"`
	IndexSize string `json:"indexSize"`
	Version   string `json:"version"`
}

func main() {
	// Extract environment variables.
	parseEnvVars()

	// Request data from Solr.
	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%s/solr/admin/metrics", solrHost, solrPort))
	if err != nil {
		log.Fatal(err)
	}

	// Ensure response body is always closed.
	defer resp.Body.Close()

	// Extract selected data from Solr.
	metrics := parseSolrData(resp)
	// Convert metrics to a flattened structure.
	info := Info(metrics)
	if err := transcode(metrics, &info); err != nil {
		log.Fatal(err)
	}

	// Push data to Remote Insights.
	var dataMap map[string]string
	if err := transcode(info, &dataMap); err != nil {
		log.Fatal(err)
	}

	if err := lagoon.ProcessFacts(dataMap, insightsEndpoint); err != nil {
		log.Fatal(err)
	}
}

func parseSolrData(resp *http.Response) Metrics {
	// Parse raw Solr API json.
	solr := Solr{}
	json.NewDecoder(resp.Body).Decode(&solr)

	// Initialise the metrics container.
	metrics := Metrics{}

	// Solr Core metrics.
	if err := transcode(solr.Metrics["solr.core."+solrCore], &metrics); err != nil {
		log.Fatal(err)
	}
	// Solr Node metrics.
	if err := transcode(solr.Metrics["solr.node"], &metrics); err != nil {
		log.Fatal(err)
	}

	return metrics
}

// Parse data to a struct.
func transcode(in any, out interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, &out)
	return nil
}

// Read and apply supported environment variables.
func parseEnvVars() {
	solrHostEnv := os.Getenv("SOLR_HOST")
	solrHost = "solr" // default.
	if solrHostEnv != "" {
		solrHost = solrHostEnv
	}

	solrPortEnv := os.Getenv("SOLR_PORT")
	solrPort = "8983" // default.
	if solrPortEnv != "" {
		solrPort = solrPortEnv
	}

	solrCoreEnv := os.Getenv("SOLR_CORE")
	solrCore = "drupal" // default.
	if solrCoreEnv != "" {
		solrCore = solrCoreEnv
	}

	insightsEndpointEnv := os.Getenv("LAGOON_INSIGHTS_ENDPOINT")
	insightsEndpoint = "http://lagoon-remote-insights-remote.lagoon.svc" // default.
	if insightsEndpointEnv != "" {
		insightsEndpoint = insightsEndpointEnv
	}
}
