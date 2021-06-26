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

| Name                                                  | Description                                              | Type    |
|-------------------------------------------------------|----------------------------------------------------------|---------|
| cne_events_amqp_received          | Metric to get number of events received  by the transport.   | Gauge |
| cne_events_amqp_published     | Metric to get number of events published by the transport.  | Gauge   |
| cne_amqp_connection_reset     | Metric to get number of connection resets.  | Gauge   |
| cne_amqp_sender     | Metric to get number of sender created.  | Gauge   |
| cne_amqp_receiver     | Metric to get number of receiver created.  | Gauge   |
| cne_status_check_amqp_published | Metric to get number of status check published by the transport | Gauge |

`cne_events_amqp_received` -  The number of events received by the amqp protocol, and their status by address.

Example
```json 
# HELP cne_events_amqp_received Metric to get number of events received  by the transport
# TYPE cne_events_amqp_received gauge
cne_events_amqp_received{address="/news-service/finance",status="success"} 8
cne_events_amqp_received{address="/news-service/sports",status="success"} 8
```

`cne_events_amqp_published` -  This metrics indicates number of events that were published via amqp , grouped by address and status.

Example
```json
# HELP cne_events_amqp_published Metric to get number of events published by the transport
# TYPE cne_events_amqp_published gauge
cne_events_amqp_published{address="/news-service/finance",status="connection reset"} 1
cne_events_amqp_published{address="/news-service/finance",status="success"} 8
cne_events_amqp_published{address="/news-service/sports",status="connection reset"} 1
cne_events_amqp_published{address="/news-service/sports",status="success"} 8
```

`cne_amqp_connection_reset` -  This metrics indicates number of types amqp connection was reset

Example
```json
# HELP cne_amqp_connection_reset Metric to get number of connection resets
# TYPE cne_amqp_connection_reset gauge
cne_amqp_connection_reset 1
```

`cne_amqp_sender` -  This metrics indicates number of amqp sender objects were created , grouped by address and status.

Example
```json
# HELP cne_amqp_sender Metric to get number of sender active
# TYPE cne_amqp_sender gauge
cne_amqp_sender{address="/news-service/finance",status="active"} 1
cne_amqp_sender{address="/news-service/sports",status="active"} 1
```

`cne_amqp_receiver` -  This metrics indicates number of amqp receiver objects were created, grouped by address and status.

Example
```json
# HELP cne_amqp_receiver Metric to get number of receiver active
# TYPE cne_amqp_receiver gauge
cne_amqp_receiver{address="/news-service/finance",status="active"} 1
cne_amqp_receiver{address="/news-service/sports",status="active"} 1
```

`cne_status_check_amqp_published` -  This metrics indicates status check that were published via amqp , grouped by address and status.

Example
```json
# HELP cne_status_check_amqp_published  Metric to get number of status check published by the transport
# TYPE cne_status_check_amqp_published gauge
cne_status_check_amqp_published{address="/news-service/finance/status",status="failed"} 1
cne_status_check_amqp_published{address="/news-service/sports/status",status="connection reset"} 1
cne_status_check_amqp_published{address="/news-service/sports/status",status="success"} 1
```



