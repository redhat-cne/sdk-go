package channel

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// DataEvent ...
type DataEvent struct {
	Address     string
	Data        *cloudevents.Event
	EventStatus EventStatus
	EndPointURI string
	StatusCh    *ListenerChannel
	//EventType defines type of data (Notification,Metric,Status)
	EventType EventType
}

//EventStatus specifies status of the event
type EventStatus int

const (
	// SUCCEED if the event is posted successfully
	SUCCEED EventStatus = 1
	//FAILED if the event  failed to post
	FAILED EventStatus = 2
	//NEW if the event is new for the consumer
	NEW EventStatus = 0
)

//EventStatus specifies status of the event
type EventType int

const (
	// CONSUMER  the data is consumer type
	CONSUMER EventType = 1
	//STATUS  the data is status type
	STATUS EventType = 2
	//PRODUCER  the data is producer type
	PRODUCER EventType = 0
	//EVENT  the data is event type
	EVENT EventType = 3
)
