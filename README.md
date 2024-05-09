# Lagoon Solr Metrics

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/salsadigitalauorg/lagoon-solr-metrics)
[![Go Report Card](https://goreportcard.com/badge/github.com/salsadigitalauorg/lagoon-solr-metrics)](https://goreportcard.com/report/github.com/salsadigitalauorg/lagoon-solr-metrics)
[![Coverage Status](https://coveralls.io/repos/github/salsadigitalauorg/lagoon-solr-metrics/badge.svg?branch=master)](https://coveralls.io/github/salsadigitalauorg/lagoon-solr-metrics?branch=master)
[![Release](https://img.shields.io/github/v/release/salsadigitalauorg/lagoon-solr-metrics)](https://github.com/salsadigitalauorg/lagoon-solr-metrics/releases/latest)

## Installation


```sh
curl -L -o lagoon-solr-metrics https://github.com/salsadigitalauorg/lagoon-solr-metrics/releases/latest/download/lagoon-solr-metrics-$(uname -s)-$(uname -m)
chmod +x lagoon-solr-metrics
mv lagoon-solr-metrics /usr/local/bin/lagoon-solr-metrics
```

### Docker

Run directly from a docker image:
```sh
docker run --rm ghcr.io/salsadigitalauorg/lagoon-solr-metrics:latest lagoon-solr-metrics
```

Or add to your docker image:
```Dockerfile
COPY --from=ghcr.io/salsadigitalauorg/lagoon-solr-metrics:latest /usr/local/bin/lagoon-solr-metrics /usr/local/bin/lagoon-solr-metrics
```

## Usage
The application needs a few environmental variables set so that it can know where to collect the metrics from and where to push it.

```
SOLR_HOST=solr
SOLR_PORT=8983
SOLR_CORE=drupal
LAGOON_INSIGHTS_ENDPOINT="http://lagoon-remote-insights-remote.lagoon.svc"
```

The above are the defaults that the application will use, but you can override them as needed in your runtime environment.


## Local development

### Build
```sh
git clone git@github.com:salsadigitalauorg/lagoon-solr-metrics.git && cd lagoon-solr-metrics
go generate ./...
go build -ldflags="-s -w" -o build/lagoon-solr-metrics .
go run . -h
```

### Run tests
```sh
go generate ./...
go test -v ./... -coverprofile=build/coverage.out
```

View coverage results:
```sh
go tool cover -html=build/coverage.out
```
