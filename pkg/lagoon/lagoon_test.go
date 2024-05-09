package lagoon_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/salsadigitalauorg/lagoon-solr-metrics/pkg/internal"
	"github.com/salsadigitalauorg/lagoon-solr-metrics/pkg/lagoon"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMustHaveEnvVars(t *testing.T) {
	assert := assert.New(t)

	origOutput := logrus.StandardLogger().Out
	defer func() {
		logrus.StandardLogger().ExitFunc = nil
		logrus.SetOutput(origOutput)
	}()

	var buf bytes.Buffer
	logrus.SetOutput(&buf)

	var fatal bool
	logrus.StandardLogger().ExitFunc = func(retCode int) {
		fatal = (retCode == 1)
	}

	t.Run("varsSet", func(t *testing.T) {
		os.Setenv("LAGOON_PROJECT", "foo")
		os.Setenv("LAGOON_ENVIRONMENT", "bar")
		defer func() {
			os.Unsetenv("LAGOON_PROJECT")
			os.Unsetenv("LAGOON_ENVIRONMENT")
		}()

		lagoon.MustHaveEnvVars()
		assert.False(fatal)
	})

	t.Run("varsUnset", func(t *testing.T) {
		lagoon.MustHaveEnvVars()
		assert.True(fatal)
		assert.Contains(buf.String(), "project & environment name required; please ensure both LAGOON_PROJECT & LAGOON_ENVIRONMENT are set")
	})
}

func Test_GetBearerTokenFromDisk(t *testing.T) {
	type args struct {
		tokenLocation string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successfully loads token from disk",
			args: args{
				tokenLocation: "./testdata/insightsbearertoken",
			},
			want:    "bearertokentext",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lagoon.GetBearerTokenFromDisk(tt.args.tokenLocation)
			if !tt.wantErr(t, err, fmt.Sprintf("getBearerTokenFromDisk(%v)", tt.args.tokenLocation)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getBearerTokenFromDisk(%v)", tt.args.tokenLocation)
		})
	}
}

func Test_FactsToInsightsRemote(t *testing.T) {
	type args struct {
		facts       []lagoon.Fact
		bearerToken string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       assert.ErrorAssertionFunc
		testBodyEqual bool
	}{
		{
			name: "Successful post",
			args: args{
				facts: []lagoon.Fact{
					{
						Name:        "version",
						Value:       "9.5.0",
						Source:      "insights:image:solr",
						Description: "Solr metric for version",
						Category:    "Solr",
					},
				},
				bearerToken: "bearertoken",
			},
			wantErr:       assert.NoError,
			testBodyEqual: true,
		},
	}
	for _, tt := range tests {

		mockServerData := internal.MockInsightsRemoteTestState{}

		serv := internal.MockRemoteInsightsServer(&mockServerData)

		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, lagoon.FactsToInsightsRemote(tt.args.facts, serv.URL, tt.args.bearerToken), fmt.Sprintf("factsToInsightsRemote(%v, %v)", tt.args.facts, tt.args.bearerToken))
			if tt.testBodyEqual == true { //let's check the data we sent through seems correct
				bodyString := mockServerData.LastCallBody
				var facts []lagoon.Fact
				err := json.Unmarshal([]byte(bodyString), &facts)
				assert.NoError(t, err)
				assert.Equalf(t, tt.args.facts, facts, "Unmarshalled Body not Equal")
			}
		})
	}
}
