package lagoon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Fact struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Source      string `json:"source"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

const SourceName string = "insights:image:solr"
const FactMaxValueLength int = 300

var project string
var environment string

func MustHaveEnvVars() {
	project = os.Getenv("LAGOON_PROJECT")
	environment = os.Getenv("LAGOON_ENVIRONMENT")
	if project == "" || environment == "" {
		log.Fatal("project & environment name required; please ensure both " +
			"LAGOON_PROJECT & LAGOON_ENVIRONMENT are set")
	}
}

const DefaultLagoonInsightsTokenLocation = "/var/run/secrets/lagoon/dynamic/insights-token/INSIGHTS_TOKEN"

// Each environment stores token on disk.
func GetBearerTokenFromDisk(tokenLocation string) (string, error) {
	// First, we check that the token exists on disk.
	_, err := os.Stat(tokenLocation)
	if err != nil {
		return "", fmt.Errorf("Unable to load insights token from disk")
	}

	b, err := os.ReadFile(tokenLocation)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(b), "\n"), nil
}

func ProcessFacts(info map[string]string, insightsEndpoint string) error {
	facts := []Fact{}

	for n, v := range info {
		if len(v) == 0 {
			continue
		}

		facts = append(facts, Fact{
			Name:        n,
			Value:       v,
			Source:      SourceName,
			Description: "Solr metric for " + n,
			Category:    "Solr",
		})
	}

	// Attempt to look up bearer token.
	bearerToken, err := GetBearerTokenFromDisk(DefaultLagoonInsightsTokenLocation)
	if err == nil { // we have a token, and so we can proceed via the internal service call.
		err = FactsToInsightsRemote(facts, insightsEndpoint, bearerToken)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	fmt.Println("Successfully pushed facts to Lagoon Remote.")
	return nil
}

// Send POST request to Insights Remote.
func FactsToInsightsRemote(facts []Fact, serviceEndpoint string, bearerToken string) error {
	bodyString, err := json.Marshal(facts)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodPost, serviceEndpoint+"/facts", bytes.NewBuffer(bodyString))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Errorf("there was an error sending the facts to '%s' : %s", serviceEndpoint, string(bytes))
	}
	return nil
}
