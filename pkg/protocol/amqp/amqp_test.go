// Copyright 2020 The Cloud Native Events Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package amqp_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	cetypes "github.com/cloudevents/sdk-go/v2/types"
	amqp1 "github.com/redhat-cne/sdk-go/pkg/protocol/amqp"
	"github.com/redhat-cne/sdk-go/pkg/types"

	cneevent "github.com/redhat-cne/sdk-go/pkg/event"
	"github.com/redhat-cne/sdk-go/pkg/event/ptp"

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
	_type           = string(ptp.PtpStateChange)
	resourceAddress = "/test1/test1"
)

// CloudNativeEvents generates cloud events for testing
func CloudNativeEvents() cneevent.Event {
	data := cneevent.Data{}
	value := cneevent.DataValue{
		Resource:  resourceAddress,
		DataType:  cneevent.NOTIFICATION,
		ValueType: cneevent.ENUMERATION,
		Value:     ptp.ACQUIRING_SYNC,
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
		Value:     ptp.ACQUIRING_SYNC,
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

func TestDeleteSender(t *testing.T) {
	time.Sleep(250 * time.Millisecond)
	addr := "test/sender/delete"
	s := "amqp://localhost:5672"

	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	assert.Equal(t, len(server.Listeners), 0)
	wg := sync.WaitGroup{}
	server.QDRRouter(&wg)

	// create a listener
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.PUBLISHER,
	}
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, len(server.Senders))

	// send data
	in <- &channel.DataChan{
		Address: addr,
		Status:  channel.DELETE,
		Type:    channel.PUBLISHER,
	}
	time.Sleep(250 * time.Millisecond)

	// read data
	assert.Equal(t, 0, len(server.Senders))
	close(closeCh)
}

func TestDeleteListener(t *testing.T) {
	addr := "test/listener/delete"
	s := "amqp://localhost:5672"

	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	assert.Equal(t, len(server.Listeners), 0)
	wg := sync.WaitGroup{}
	server.QDRRouter(&wg)

	// create a listener
	in <- &channel.DataChan{
		Address:             addr,
		Type:                channel.SUBSCRIBER,
		Status:              channel.NEW,
		ProcessEventFn:      func(e interface{}) error { return nil },
		OnReceiveOverrideFn: func(e cloudevents.Event, dataChan *channel.DataChan) error { return nil },
	}
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, len(server.Listeners))
	// send data
	in <- &channel.DataChan{
		Address: addr,
		Status:  channel.DELETE,
		Type:    channel.SUBSCRIBER,
	}
	// read data
	time.Sleep(250 * time.Millisecond)
	// read data
	assert.Equal(t, 0, len(server.Listeners))
	close(closeCh)
}

func TestSendSuccessStatus(t *testing.T) {
	time.Sleep(250 * time.Millisecond)
	addr := "test/test/success"
	s := "amqp://localhost:5672"

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	server.QDRRouter(&wg)

	// send status, this will create status listener
	// always you need to define how you handle status  when it is received
	// do not override for events
	in <- &channel.DataChan{
		Address:             fmt.Sprintf("%s/%s", addr, "status"),
		Status:              channel.NEW,
		Type:                channel.SUBSCRIBER,
		ProcessEventFn:      func(e interface{}) error { return nil },
		OnReceiveOverrideFn: func(e cloudevents.Event, dataChan *channel.DataChan) error { return nil },
	}

	// create a sender
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Type:    channel.PUBLISHER,
	}

	// ping for status, this will  send the  status check ping to the address
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.STATUS,
	}
	// read data
	d := <-out
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)
	close(closeCh)
	//waitTimeout(&wg, timeout)
}

func TestSendFailureStatus(t *testing.T) {
	addr := "test/test/failure"
	s := "amqp://localhost:5672"

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	server.QDRRouter(&wg)

	// send status, this will create status listener
	// always you need to define how you handle status  when it is received
	// do not override for events
	in <- &channel.DataChan{
		Address:        fmt.Sprintf("%s/%s", addr, "status"),
		Status:         channel.NEW,
		Type:           channel.SUBSCRIBER,
		ProcessEventFn: func(e interface{}) error { return fmt.Errorf("EVENT PROCESS ERROR") },
		OnReceiveOverrideFn: func(e cloudevents.Event, dataChan *channel.DataChan) error {
			return fmt.Errorf("STATUS RECEEIVE ERROR")
		},
	}

	// create a sender
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Type:    channel.PUBLISHER,
	}

	// ping for status, this will  send the  status check ping to the address
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.STATUS,
	}
	// read data
	d := <-out
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.FAILED, d.Status)
	close(closeCh)
}

func TestSendEvent(t *testing.T) {
	addr := "test/test/event"
	s := "amqp://localhost:5672"

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := amqp1.InitServer(s, in, out, closeCh)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	server.QDRRouter(&wg)

	// create a sender
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.PUBLISHER,
	}

	// create a listener
	in <- &channel.DataChan{
		Address: addr,
		Type:    channel.SUBSCRIBER,
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
	assert.Equal(t, channel.EVENT, d.Type)
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
		Type:                channel.SUBSCRIBER,
		ProcessEventFn:      func(e interface{}) error { return nil },
		OnReceiveOverrideFn: func(e cloudevents.Event, dataChan *channel.DataChan) error { return nil },
	}

	// create a sender
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Type:    channel.PUBLISHER,
	}

	// ping for status, this will  send the  status check ping to the address
	in <- &channel.DataChan{
		Address: fmt.Sprintf("%s/%s", addr, "status"),
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}
	// read data
	d = <-out
	assert.Equal(t, channel.EVENT, d.Type)
	close(closeCh)
}
