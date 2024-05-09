FROM --platform=$BUILDPLATFORM golang:1.21 AS builder

ARG VERSION
ARG COMMIT

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
