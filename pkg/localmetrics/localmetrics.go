package localmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricStatus metrics status
type MetricStatus string

const (
	// SUCCESS ...
	SUCCESS MetricStatus = "success"
	// FAILED ...
	FAILED MetricStatus = "failed"
	// CONNECTION_RESET ...
	CONNECTION_RESET MetricStatus = "connection reset"
)

var (

	//amqpEventReceivedCount ...  Total no of events received by the transport
	amqpEventReceivedCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cne_events_amqp_received",
			Help: "Metric to get number of events received  by the transport",
		}, []string{"address", "status"})
	//amqpEventPublishedCount ...  Total no of events published by the transport
	amqpEventPublishedCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cne_events_amqp_published",
			Help: "Metric to get number of events published by the transport",
		}, []string{"address", "status"})

	//amqpConnectionResetCount ...  Total no of connection resets
	amqpConnectionResetCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cne_amqp_connections_resets",
			Help: "Metric to get number of connection resets",
		}, []string{})

	//amqpSenderCount ...  Total no of events published by the transport
	amqpSenderCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cne_amqp_sender",
			Help: "Metric to get number of sender created",
		}, []string{"address", "status"})

	//amqpReceiverCount ...  Total no of events published by the transport
	amqpReceiverCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cne_amqp_receiver",
			Help: "Metric to get number of receiver created",
		}, []string{"address", "status"})
)

// RegisterMetrics ...
func RegisterMetrics() {
	prometheus.MustRegister(amqpEventReceivedCount)
	prometheus.MustRegister(amqpEventPublishedCount)
	prometheus.MustRegister(amqpConnectionResetCount)
	prometheus.MustRegister(amqpSenderCount)
	prometheus.MustRegister(amqpReceiverCount)
}

// UpdateTransportConnectionResetCount ...
func UpdateTransportConnectionResetCount(val int) {
	amqpConnectionResetCount.With(prometheus.Labels{}).Add(float64(val))
}

// UpdateEventReceivedCount ...
func UpdateEventReceivedCount(address string, status MetricStatus, val int) {
	amqpEventReceivedCount.With(
		prometheus.Labels{"address": address, "status": string(status)}).Add(float64(val))
}

// UpdateEventCreatedCount ...
func UpdateEventCreatedCount(address string, status MetricStatus, val int) {
	amqpEventPublishedCount.With(
		prometheus.Labels{"address": address, "status": string(status)}).Add(float64(val))
}

// UpdateSenderCreatedCount ...
func UpdateSenderCreatedCount(address string, status MetricStatus, val int) {
	amqpSenderCount.With(
		prometheus.Labels{"address": address, "status": string(status)}).Add(float64(val))
}

// UpdateReceiverCreatedCount ...
func UpdateReceiverCreatedCount(address string, status MetricStatus, val int) {
	amqpReceiverCount.With(
		prometheus.Labels{"address": address, "status": string(status)}).Add(float64(val))
}
