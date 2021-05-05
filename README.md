# sdk-go
GO SDK for Cloud Native Events API

[![go-doc](https://godoc.org/github.com/redhat-cne/sdk-go?status.svg)](https://godoc.org/github.com/redhat-cne/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/redhat-cne/sdk-go)](https://goreportcard.com/report/github.com/redhat-cne/sdk-go)
[![LICENSE](https://img.shields.io/github/license/redhat-cne/sdk-go.svg)](https://github.com/redhat-cne/sdk-go/blob/main/LICENSE)


### This SDK is used to develop Cloud Event Proxy
https://github.com/redhat-cne/cloud-event-proxy


## Collecting metrics with Prometheus
Cloud native events SDK-go comes with following metrics collectors .
1. Number of events received  by the transport
2. Number of events published by the transport.
3. Number of connection resets.
4. Number of sender created
5. Number of receiver created
      
[SDK-GO Metrics details ](docs/metrics.md)



