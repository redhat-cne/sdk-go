package amqp_test

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	ce_types "github.com/cloudevents/sdk-go/v2/types"
	"github.com/redhat-cne/sdk-go/types"
	"log"
	"net/url"
	"time"

	cne_event "github.com/redhat-cne/sdk-go/event"

	"github.com/redhat-cne/sdk-go/channel"

	"github.com/stretchr/testify/assert"
	"sync"

	amqp1 "github.com/redhat-cne/sdk-go/protocol/amqp"

	"testing"
)

func strptr(s string) *string { return &s }

var (
	ceSource        = ce_types.URIRef{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/source"}}
	ceTimestamp     = ce_types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	cneTimestamp    = types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	ceSchema        = ce_types.URI{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/schema"}}
	_type           = "ptp_status_type"
	resourceAddress = "/test1/test1"
)

// CloudNativeEvents generates cloud events for testing
func CloudNativeEvents() cne_event.Event {
	data := cne_event.Data{}
	value := cne_event.DataValue{
		Resource:  resourceAddress,
		DataType:  cne_event.NOTIFICATION,
		ValueType: cne_event.ENUMERATION,
		Value:     cne_event.GNSS_ACQUIRING_SYNC,
	}
	data.SetVersion("1.0")   //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

	cne := cne_event.Event{
		ID:              "123",
		Type:            _type,
		DataContentType: strptr(event.ApplicationJSON),
		Time:            &cneTimestamp,
		Data:            &data,
	}

	return cne
}

//CloudEvents return cloud events objects
func CloudEvents() cloudevents.Event {
	data := cne_event.Data{}
	value := cne_event.DataValue{
		Resource:  resourceAddress,
		DataType:  cne_event.NOTIFICATION,
		ValueType: cne_event.ENUMERATION,
		Value:     cne_event.GNSS_ACQUIRING_SYNC,
	}
	data.SetVersion("1.0")   //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

	e := cloudevents.Event{
		Context: cloudevents.EventContextV1{
			Type:       "com.example.FullEvent",
			Source:     ceSource,
			ID:         "full-event",
			Time:       &ceTimestamp,
			DataSchema: &ceSchema,
			Subject:    strptr("topic"),
		}.AsV1(),
	}
	cne := CloudNativeEvents()

	_ = e.SetData(cloudevents.ApplicationJSON, cne.Data)

	return e
}

func TestSendEvent(t *testing.T) {

	addr := "test/test2"
	s := "amqp://localhost:5672"

	event := CloudEvents()
	in := make(chan channel.DataEvent)
	out := make(chan channel.DataEvent)
	close := make(chan bool)
	server, err := amqp1.InitServer(s, in, out, close)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	go server.QDRRouter(&wg)

	// create a sender
	in <- channel.DataEvent{
		Address:   addr,
		EventType: channel.SENDER,
	}

	// create a listener
	in <- channel.DataEvent{
		Address:   addr,
		EventType: channel.LISTENER,
	}

	// send data
	in <- channel.DataEvent{
		Address:     addr,
		Data:        &event,
		EventStatus: channel.NEW,
		EndPointURI: "http://localhost",
		EventType:   channel.EVENT,
	}

	// read data
	d := <-out
	log.Printf("Processing out channel")
	assert.Equal(t, d.EventType, channel.EVENT)
	assert.Equal(t, d.Address, addr)
	dd := cne_event.Data{}
	err = json.Unmarshal(event.Data(), &dd)
	assert.Nil(t, err)
	assert.Equal(t, dd.Version, "1.0")
	close <- true

}
