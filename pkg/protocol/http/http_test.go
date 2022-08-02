package http_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	cetypes "github.com/cloudevents/sdk-go/v2/types"
	"github.com/google/uuid"
	"github.com/redhat-cne/sdk-go/pkg/channel"
	cneevent "github.com/redhat-cne/sdk-go/pkg/event"
	"github.com/redhat-cne/sdk-go/pkg/event/ptp"
	ceHttp "github.com/redhat-cne/sdk-go/pkg/protocol/http"
	"github.com/redhat-cne/sdk-go/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/pkg/store"
	subscriber "github.com/redhat-cne/sdk-go/pkg/subscriber"
	"github.com/redhat-cne/sdk-go/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func strptr(s string) *string { return &s }

var (
	storePath         = "."
	serverClientID, _ = uuid.Parse("d73555b4-b01e-4802-89f1-f058f215a1f8")
	clientClientID, _ = uuid.Parse("5202d2c4-f652-4974-b24d-1934f0d819e3")
	subscriptionOneID = "123e4567-e89b-12d3-a456-426614174001"
	subscriptionTwoID = "123e4567-e89b-12d3-a456-426614174002"
	serverAddress     = "http://localhost:8086"
	clientAddress     = "http://localhost:8087"
	hostPort          = 8086
	clientPort        = 8087

	subscriptionOne = &pubsub.PubSub{
		ID:       subscriptionOneID,
		Resource: "test/test/1",
	}

	subscriptionTwo = &pubsub.PubSub{
		ID:       subscriptionTwoID,
		Resource: "test/test/2",
	}

	subscriberWithOneEventCheck = subscriber.Subscriber{
		ClientID: clientClientID.String(),
		SubStore: &store.PubSubStore{
			RWMutex: sync.RWMutex{},
			Store:   map[string]*pubsub.PubSub{subscriptionOneID: subscriptionOne},
		},
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8087"}},
		Status:      1,
	}
)
var (
	ceSource     = cetypes.URIRef{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/source"}}
	ceTimestamp  = cetypes.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	cneTimestamp = types.Timestamp{Time: time.Date(2020, 03, 21, 12, 34, 56, 780000000, time.UTC)}
	ceSchema     = cetypes.URI{URL: url.URL{Scheme: "http", Host: "example.com", Path: "/schema"}}
	_type        = string(ptp.PtpStateChange)
)

// CloudNativeEvents generates cloud events for testing
func CloudNativeEvents() cneevent.Event {
	data := cneevent.Data{}
	value := cneevent.DataValue{
		Resource:  _type,
		DataType:  cneevent.NOTIFICATION,
		ValueType: cneevent.ENUMERATION,
		Value:     ptp.ACQUIRING_SYNC,
	}
	data.SetVersion("1.0")   //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

	cne := cneevent.Event{
		ID:              "123",
		Type:            _type,
		Source:          subscriptionOne.Resource,
		DataContentType: strptr(event.ApplicationJSON),
		Time:            &cneTimestamp,
		DataSchema:      nil,
		Data:            &data,
	}

	return cne
}

// CloudNativeEvents generates cloud events for testing
func subscriberCloudEvents(dataType string) *cloudevents.Event {

	e := cloudevents.Event{
		Context: cloudevents.EventContextV1{
			Type:       dataType,
			Source:     ceSource,
			ID:         "full-event",
			Time:       &ceTimestamp,
			DataSchema: &ceSchema,
			Subject:    strptr("topic"),
		}.AsV1(),
	}
	e.SetData(cloudevents.ApplicationJSON, subscriberWithOneEventCheck)
	return &e
}

func statusCloudEvents() *cloudevents.Event {

	e := cloudevents.Event{
		Context: cloudevents.EventContextV1{
			Type:       channel.STATUS.String(),
			Source:     ceSource,
			ID:         "full-event",
			Time:       &ceTimestamp,
			DataSchema: &ceSchema,
			Subject:    strptr("topic"),
		}.AsV1(),
	}
	e.SetData(cloudevents.ApplicationJSON, subscriptionTwo)
	return &e
}

//CloudEvents return cloud events objects
func CloudEvents() cloudevents.Event {
	data := cneevent.Data{}
	value := cneevent.DataValue{
		Resource:  subscriptionOne.Resource,
		DataType:  cneevent.NOTIFICATION,
		ValueType: cneevent.ENUMERATION,
		Value:     ptp.ACQUIRING_SYNC,
	}
	data.SetVersion("1.0")   //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

	e := cloudevents.Event{
		Context: cloudevents.EventContextV1{
			Type:       _type,
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

// client  registers with server and ask for status , also receive any event that was generated
func createClient(clientS *ceHttp.Server, closeCh chan struct{}, withStatus bool, clientOutChannel chan *channel.DataChan) {
	in := make(chan *channel.DataChan, 10)

	clientS, err := ceHttp.InitServer(clientAddress, clientPort, storePath, in, clientOutChannel, closeCh, nil, nil)
	if err != nil {
		log.Infof("error creating client ")
	}
	clientS.ClientID = clientClientID
	clientS.RegisterPublishers(serverAddress)
	wg := sync.WaitGroup{}
	time.Sleep(250 * time.Millisecond)
	// Start the server and channel processor
	clientS.Start(&wg)
	clientS.HTTPProcessor(&wg)
	time.Sleep(250 * time.Millisecond)
	// create a subscription
	in <- &channel.DataChan{
		ID:      subscriptionOneID,
		Address: subscriptionOne.Resource,
		Type:    channel.SUBSCRIBER,
	}
	time.Sleep(250 * time.Millisecond)
	if withStatus {
		// ping for status, this will  send the  status check ping to the address
		in <- &channel.DataChan{
			Address: subscriptionOne.Resource,
			Status:  channel.NEW,
			Type:    channel.STATUS,
		}
	}

	<-closeCh

	/*
		s.SetType(channel.STATUS.String())
		log.Infof("sending STATUS request")
		if e2 := clientS.Post(serverAddress, s); e2 != nil {
			log.Infof("failed to ping for status %s", e2.Error())
		} else {
			log.Infof("ping was successfull ")
		}
	*/
}
func TestSubscribeCreated(t *testing.T) {
	in := make(chan *channel.DataChan, 10)
	out := make(chan *channel.DataChan, 10)
	closeCh := make(chan struct{})
	eventChannel := make(chan *channel.DataChan, 10)
	server, err := ceHttp.InitServer(serverAddress, hostPort, storePath, in, out, closeCh, nil, nil)
	if err != nil {
		t.Skipf("http failed(%#v): %v", server, err)
	}
	server.ClientID = serverClientID

	wg := sync.WaitGroup{}
	// Start the server and channel proceesor
	server.Start(&wg)
	server.HTTPProcessor(&wg)
	var clientS *ceHttp.Server
	go createClient(clientS, closeCh, false, eventChannel)
	time.Sleep(500 * time.Millisecond)
	<-out
	assert.Equal(t, 1, len(server.Sender))
	d := <-eventChannel
	assert.Equal(t, channel.SUBSCRIBER, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)

	close(closeCh)

}

func TestSendEvent(t *testing.T) {
	//closeClient := make(chan struct{})
	//createClient(clientAddress, closeClient)
	time.Sleep(2 * time.Second)
	e := CloudEvents()
	in := make(chan *channel.DataChan, 10)
	out := make(chan *channel.DataChan, 10)
	clientOutChannel := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := ceHttp.InitServer(serverAddress, hostPort, storePath, in, out, closeCh, nil, nil)
	if err != nil {
		t.Skipf("http failed(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	// Start the server and channel proceesor
	server.Start(&wg)
	server.HTTPProcessor(&wg)
	time.Sleep(500 * time.Millisecond)
	var clientS *ceHttp.Server
	go createClient(clientS, closeCh, false, clientOutChannel)
	//  read what server has in outChannel
	<-out
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, len(server.Sender))
	// read what client put in out channel
	d := <-clientOutChannel
	assert.Equal(t, channel.SUBSCRIBER, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)

	// send event
	in <- &channel.DataChan{
		Address: subscriptionOne.Resource,
		Data:    &e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}
	// read event
	log.Info("waiting for event channel from the client when it received the event")
	d = <-clientOutChannel // a client needs to break out or else it will be holding it forever
	assert.Equal(t, channel.EVENT, d.Type)
	dd := cneevent.Data{}
	err = json.Unmarshal(e.Data(), &dd)
	assert.Nil(t, err)
	assert.Equal(t, dd.Version, "1.0")

	log.Info("waiting for event response")
	d = <-out
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)

	time.Sleep(250 * time.Millisecond)
	close(closeCh)
	//close(closeClient)
}

func TestSendSuccessStatus(t *testing.T) {
	//time.Sleep(250 * time.Millisecond)

	//closeClient := make(chan struct{})
	//createClient(clientAddress, closeClient)

	e := CloudEvents()
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	clientOutChannel := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := ceHttp.InitServer(serverAddress, hostPort, storePath, in, out, closeCh, func(e cloudevents.Event, dataChan *channel.DataChan) error {
		dataChan.Address = clientAddress
		e.SetType(channel.EVENT.String())
		if err := ceHttp.Post(fmt.Sprintf("%s/event", clientAddress), e); err != nil {
			log.Errorf("error %s sending event %v at  %s", err, e, clientAddress)
			return err
		}
		return nil
	}, nil)
	if err != nil {
		t.Skipf("http failed(%#v): %v", server, err)
	}
	wg := sync.WaitGroup{}
	// Start the server and channel proceesor
	server.Start(&wg)
	server.HTTPProcessor(&wg)

	// create a sender
	var clientS *ceHttp.Server
	go createClient(clientS, closeCh, true, clientOutChannel)
	<-out
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, len(server.Sender))
	// read what client put in out channel
	d := <-clientOutChannel
	assert.Equal(t, channel.SUBSCRIBER, d.Type)
	assert.Equal(t, channel.SUCCESS, d.Status)
	// read Status
	d = <-clientOutChannel
	assert.Equal(t, channel.EVENT, d.Type)
	assert.Equal(t, channel.NEW, d.Status)
	<-out
	// read event
	log.Info("waiting for event channel from the client when it received the event")
	d = <-clientOutChannel // a client needs to break out or else it will be holding it forever
	assert.Equal(t, channel.STATUS, d.Type)
	dd := cneevent.Data{}
	err = json.Unmarshal(e.Data(), &dd)
	assert.Nil(t, err)
	assert.Equal(t, dd.Version, "1.0")

	close(closeCh)
	//waitTimeout(&wg, timeout)
}

/*
func TestDeleteSender(t *testing.T) {
	in := make(chan *channel.DataChan)
	out := make(chan *channel.DataChan)
	closeCh := make(chan struct{})
	server, err := ceHttp.InitServer(serverAddress, hostPort, storePath, in, out, closeCh, nil)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", server, err)
	}

	wg := sync.WaitGroup{}
	server.HTTPProcessor(&wg)

	// create a sender
	in <- &channel.DataChan{
		Address: clientAddress,
		Type:    channel.PUBLISHER,
		Data:    SubscriberCloudEvents(),
	}

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, len(server.Sender))

	// send data
	in <- &channel.DataChan{
		Address: clientAddress,
		Status:  channel.DELETE,
		Type:    channel.PUBLISHER,
		Data:    SubscriberCloudEvents(),
	}
	time.Sleep(250 * time.Millisecond)

	// read data
	assert.Equal(t, 0, len(server.Sender))
	close(closeCh)
}

*/

func TestTeardown(t *testing.T) {
	_ = os.Remove(fmt.Sprintf("./%s.json", clientClientID.String()))
}
