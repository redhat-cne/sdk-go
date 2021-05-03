package channel

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/redhat-cne/sdk-go/pkg/event"
)

// DataChan ...
type DataChan struct {
	Address string
	Data    *cloudevents.Event
	Status  Status
	//Type defines type of data (Notification,Metric,Status)
	Type Type
	OnReceiveFn  func(e cloudevents.Event)
	// OnReceiveOverrideFn Optional for event, but override for status pings.This is an override function on receiving msg by amqp listener,
	// if not set then the data is sent to out channel and processed by side car  default method
	OnReceiveOverrideFn func(e cloudevents.Event) error
	// ProcessEventFn  Optional, this allows to customize message handler thar was received at the out channel
	ProcessEventFn func(e event.Event) error
}
