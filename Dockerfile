FROM --platform=$BUILDPLATFORM golang:1.21 AS builder

ARG VERSION
ARG COMMIT
ARG SOLR_HOST
ARG SOLR_PORT
ARG SOLR_CORE
ARG LAGOON_INSIGHTS_ENDPOINT

ADD . $GOPATH/src/github.com/salsadigitalauorg/lagoon-solr-metrics/

WORKDIR $GOPATH/src/github.com/salsadigitalauorg/lagoon-solr-metrics

ENV CGO_ENABLED 0

ARG TARGETOS TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} && \
    go mod tidy && \
    go generate ./... && \
    go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o build/lagoon-solr-metrics

FROM scratch

COPY --from=builder /go/src/github.com/salsadigitalauorg/lagoon-solr-metrics/build/lagoon-solr-metrics /usr/local/bin/lagoon-solr-metrics

ENV SOLR_HOST=${SOLR_HOST}
ENV SOLR_PORT=${SOLR_PORT}
ENV SOLR_CORE=${SOLR_CORE}
ENV LAGOON_INSIGHTS_ENDPOINT=${LAGOON_INSIGHTS_ENDPOINT}

ENTRYPOINT ["/usr/local/bin/lagoon-solr-metrics"]