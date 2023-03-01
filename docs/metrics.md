---
Title: Metrics
---

SDK-GO populates [Prometheus][prometheus]  collectors for metrics reporting. The metrics can be used for real-time monitoring and debugging.
sdk-go metrics collector does not persist its metrics; if a member restarts, the metrics will be reset.

The simplest way to see the available metrics is to cURL the metrics endpoint `/metrics`. The format is described [here](http://prometheus.io/docs/instrumenting/exposition_formats/).

Follow the [Prometheus getting started doc](http://prometheus.io/docs/introduction/getting_started/) to spin up a Prometheus server to collect metrics.

The naming of metrics follows the suggested [Prometheus best practices](http://prometheus.io/docs/practices/naming/).

A metric name has an `cne`  prefix as its namespace, and a subsystem prefix .

###Registering collector in your application
The collector needs to be registered in the consuming application by calling `RegisterMetrics()`  method from `sdk-go/pkg/localmetrics package`


## cne namespace metrics

The metrics under the `cne` prefix are for monitoring .  If there is any change of these metrics, it will be included in release notes.


### Metrics

These metrics describe the status of the cloud native events, publisher and subscriptions .

All these metrics are prefixed with `cne_`

| Name                        | Description                                              | Type    |
|-----------------------------|----------------------------------------------------------|---------|
| cne_transport_events_received | Metric to get number of events received  by the transport.   | Gauge |
| cne_transport_events_published       | Metric to get number of events published by the transport.  | Gauge   |
| cne_transport_connection_reset   | Metric to get number of connection resets.  | Gauge   |
| cne_transport_sender             | Metric to get number of sender created.  | Gauge   |
| cne_transport_receiver           | Metric to get number of receiver created.  | Gauge   |
| cne_transport_status_check_published | Metric to get number of status check published by the transport | Gauge |

`cne_transport_events_received` -  The number of events received by the transport protocol, and their status by address.

Example
``` 
# HELP cne_transport_events_received Metric to get number of events received  by the transport
# TYPE cne_transport_events_received gauge
cne_transport_events_received{address="/news-service/finance",status="success"} 8
cne_transport_events_received{address="/news-service/sports",status="success"} 8
```

`cne_transport_events_published` -  This metrics indicates number of events that were published via transport , grouped by address and status.

Example
```
# HELP cne_transport_events_published Metric to get number of events published by the transport
# TYPE cne_transport_events_published gauge
cne_transport_events_published{address="/news-service/finance",status="connection reset"} 1
cne_transport_events_published{address="/news-service/finance",status="success"} 8
cne_transport_events_published{address="/news-service/sports",status="connection reset"} 1
cne_transport_events_published{address="/news-service/sports",status="success"} 8
```

`cne_transport_connection_reset` -  This metrics indicates number of types transport connection was reset

Example
```
# HELP cne_transport_connection_reset Metric to get number of connection resets
# TYPE cne_transport_connection_reset gauge
cne_transport_connection_reset 1
```

`cne_transport_sender` -  This metrics indicates number of transport sender objects were created , grouped by address and status.

Example
```
# HELP cne_transport_sender Metric to get number of sender active
# TYPE cne_transport_sender gauge
cne_transport_sender{address="/news-service/finance",status="active"} 1
cne_transport_sender{address="/news-service/sports",status="active"} 1
```

`cne_transport_receiver` -  This metrics indicates number of transport receiver objects were created, grouped by address and status.

Example
```
# HELP cne_transport_receiver Metric to get number of receiver active
# TYPE cne_transport_receiver gauge
cne_transport_receiver{address="/news-service/finance",status="active"} 1
cne_transport_receiver{address="/news-service/sports",status="active"} 1
```

`cne_transport_status_check_published` -  This metrics indicates status check that were published via transport , grouped by address and status.

Example
```
# HELP cne_transport_status_check_published  Metric to get number of status check published by the transport
# TYPE cne_transport_status_check_published gauge
cne_transport_status_check_published{address="/news-service/finance/status",status="failed"} 1
cne_transport_status_check_published{address="/news-service/sports/status",status="connection reset"} 1
cne_transport_status_check_published{address="/news-service/sports/status",status="success"} 1
```



