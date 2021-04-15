package amqp_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	cetypes "github.com/cloudevents/sdk-go/v2/types"
	amqp1 "github.com/redhat-cne/sdk-go/pkg/protocol/amqp"
	"github.com/redhat-cne/sdk-go/pkg/types"

	cneevent "github.com/redhat-cne/sdk-go/pkg/event"

	"github.com/redhat-cne/sdk-go/pkg/channel"

	"sync"

	"github.com/stretchr/testify/assert"

	"testing"
)

func strptr(s string) *string { return &s }

var (
	ceSource        = cetypes.URIRef{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/source"}}
	ceTimestamp     = cetypes.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	cneTimestamp    = types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	ceSchema        = cetypes.URI{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/schema"}}
	_type           = "ptp_status_type"
	resourceAddress = "/test1/test1"
)

// CloudNativeEvents generates cloud events for testing
func CloudNativeEvents() cneevent.Event {
	data := cneevent.Data{}
	value := cneevent.DataValue{
		Resource:  resourceAddress,
		DataType:  cneevent.NOTIFICATION,
		ValueType: cneevent.ENUMERATION,
		Value:     cneevent.ACQUIRING_SYNC,
	}
	data.SetVersion("1.0")   //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

	cne := cneevent.Event{
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
	data := cneevent.Data{}
	value := cneevent.DataValue{
		Resource:  resourceAddress,
		DataType:  cneevent.NOTIFICATION,
		ValueType: cneevent.ENUMERATION,
		Value:     cneevent.ACQUIRING_SYNC,
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

func TestSendSuccessStatus(t *testing.T) {
	addr := "test/test2"
	s := "amqp://localhost:5672"

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan bool)
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	go server.QDRRouter(&wg)

	// send status, this will create status listener
	// always you need to define how you handle status  when it is received
	// do not override for events
	in <- &channel.DataChan{
		Address:             fmt.Sprintf("%s/%s", addr, "status"),
		Status:              channel.NEW,
		Type:                channel.LISTENER,
		ProcessEventFn:      func(e cneevent.Event) error { return nil },
		OnReceiveOverrideFn: func(e cloudevents.Event) error { return nil },
	}

	// create a sender
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Type:    channel.SENDER,
	}

	// ping for status, this will  send the  status check ping to the address
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}

	log.Printf("Reading out channel")
	// read data
	d := <-out
	log.Printf("Processing out channel")
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.SUCCEED, d.Status)
	log.Printf("sending close")
	closeCh <- true
	wg.Wait()

}

func TestSendFailureStatus(t *testing.T) {
	addr := "test/test2"
	s := "amqp://localhost:5672"

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan bool)
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	go server.QDRRouter(&wg)

	// send status, this will create status listener
	// always you need to define how you handle status  when it is received
	// do not override for events
	in <- &channel.DataChan{
		Address:             fmt.Sprintf("%s/%s", addr, "status"),
		Status:              channel.NEW,
		Type:                channel.LISTENER,
		ProcessEventFn:      func(e cneevent.Event) error { return fmt.Errorf("EVENT PROCESS ERROR") },
		OnReceiveOverrideFn: func(e cloudevents.Event) error { return fmt.Errorf("STATUS RECEEIVE ERROR") },
	}

	// create a sender
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Type:    channel.SENDER,
	}

	// ping for status, this will  send the  status check ping to the address
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}

	log.Printf("Reading out channel")
	// read data
	d := <-out
	log.Printf("Processing out channel")
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.FAILED, d.Status)
	log.Printf("sending close")
	closeCh <- true
	wg.Wait()

}

func TestSendEvent(t *testing.T) {
	addr := "test/test2"
	s := "amqp://localhost:5672"

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan bool)
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	go server.QDRRouter(&wg)

	// create a sender
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.SENDER,
	}

	// create a listener
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.LISTENER,
	}

	// send data
	in <- &channel.DataChan{
		Address: addr,
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}
	// read data
	d := <-out
	log.Printf("Processing out channel %v", d)
	dd := cneevent.Data{}
	err = json.Unmarshal(e.Data(), &dd)
	assert.Nil(t, err)
	assert.Equal(t, dd.Version, "1.0")

	// send status, this will create status listener
	// always you need to define how you handle status  when it is received
	// do not override for events
	in <- &channel.DataChan{
		Address:             fmt.Sprintf("%s/%s", addr, "status"),
		Status:              channel.NEW,
		Type:                channel.LISTENER,
		ProcessEventFn:      func(e cneevent.Event) error { return nil },
		OnReceiveOverrideFn: func(e cloudevents.Event) error { return nil },
	}

	// create a sender
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Type:    channel.SENDER,
	}

	// ping for status, this will  send the  status check ping to the address
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}

	log.Printf("Reading out channel")
	// read data
	d = <-out
	log.Printf("Processing out channel")
	assert.Equal(t, channel.EVENT, d.Type)

	log.Printf("sending close")
	closeCh <- true
	wg.Wait()

}

func TestDeleteListener(t *testing.T) {
	addr := "test/test2"
	s := "amqp://localhost:5672"

	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan bool)
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	assert.Equal(t, len(server.Listeners), 0)
	wg := sync.WaitGroup{}
	go server.QDRRouter(&wg)

	// create a listener
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.LISTENER,
	}
	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, len(server.Listeners))

	// send data
	in <- &channel.DataChan{
		Address: addr,
		Status:  channel.DELETE,
		Type:    channel.LISTENER,
	}
	time.Sleep(2 * time.Second)
	// read data
	assert.Equal(t, 0, len(server.Listeners))

}

func TestDeleteSender(t *testing.T) {
	addr := "test/test2"
	s := "amqp://localhost:5672"

	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan bool)
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	assert.Equal(t, len(server.Listeners), 0)
	wg := sync.WaitGroup{}
	go server.QDRRouter(&wg)

	// create a listener
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.SENDER,
	}
	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, len(server.Senders))

	// send data
	in <- &channel.DataChan{
		Address: addr,
		Status:  channel.DELETE,
		Type:    channel.SENDER,
	}
	time.Sleep(2 * time.Second)
	// read data
	assert.Equal(t, 0, len(server.Senders))

}
