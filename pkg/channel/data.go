package channel

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// DataChan ...
type DataChan struct {
	Address  string
	Data     *cloudevents.Event
	Status   Status
	StatusCh *ListenerChannel
	//Type defines type of data (Notification,Metric,Status)
	Type Type
}
