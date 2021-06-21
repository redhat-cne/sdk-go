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

package hwevent

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/redhat-cne/sdk-go/pkg/channel"
	"github.com/redhat-cne/sdk-go/pkg/pubsub"

	"github.com/redhat-cne/sdk-go/pkg/hwevent"
)

//PublishCloudEventToLog .. publish event data to a log
func PublishCloudEventToLog(e cloudevents.Event) {
	log.Infof("Publishing event to log %#v", e)
}

//CloudNativeEvent gets Cloud Native Event object
func CloudNativeEvent() hwevent.Event {
	return hwevent.Event{Type: "HwEvent"}
}

//CloudNativeData gets Cloud Native Event object
func CloudNativeData() hwevent.Data {
	return hwevent.Data{}
}

//SendEventToLog ...
func SendEventToLog(e hwevent.Event) {
	log.Infof("Publishing event to log %#v", e)
}

//SendNewEventToDataChannel send created publisher information for QDR to process
func SendNewEventToDataChannel(inChan chan<- *channel.DataChan, address string, e *cloudevents.Event) {
	// go ahead and create QDR to this address
	inChan <- &channel.DataChan{
		Address: address,
		Data:    e,
		Status:  channel.NEW,
		Type:    channel.EVENT,
	}
}

//SendStatusToDataChannel send created publisher information for QDR to process
func SendStatusToDataChannel(inChan chan<- *channel.DataChan, status channel.Status, address string) {
	// go ahead and create QDR to this address
	inChan <- &channel.DataChan{
		Address: address,
		Type:    channel.EVENT,
		Status:  status,
	}
}

// SendCloudEventsToDataChannel sends data event in cloudevents format to data channel
func SendCloudEventsToDataChannel(inChan chan<- *channel.DataChan, status channel.Status, address string, e cloudevents.Event) {
	inChan <- &channel.DataChan{
		Address: address,
		Data:    &e,
		Status:  status,
		Type:    channel.EVENT,
	}
}

//CreateCloudEvents create new cloud event from cloud native events and pubsub
func CreateCloudEvents(e hwevent.Event, ps pubsub.PubSub) (*cloudevents.Event, error) {
	ce := cloudevents.NewEvent(cloudevents.VersionV03)
	ce.SetTime(e.GetTime())
	ce.SetType(e.Type)
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	ce.SetSource(ps.Resource) // bus address
	ce.SetSpecVersion(cloudevents.VersionV03)
	ce.SetID(uuid.New().String())
	if err := ce.SetData(cloudevents.ApplicationJSON, e.GetData()); err != nil {
		return nil, err
	}
	return &ce, nil
}

// GetCloudNativeEvents  get event data from cloud events object if its valid else return error
func GetCloudNativeEvents(ce cloudevents.Event) (e hwevent.Event, err error) {
	if ce.Data() == nil {
		return e, fmt.Errorf("event data is empty")
	}
	data := hwevent.Data{}
	log.Infof("DZK 0620-1")
	log.Infof("DZK 0620-1 %v", ce.Data())
	if err = json.Unmarshal(ce.Data(), &data); err != nil {
		return
	}
	//data = ce.Data()
	log.Infof("DZK 0620-2")
	e.SetDataContentType(hwevent.ApplicationJSON)
	e.SetTime(ce.Time())
	e.SetType(ce.Type())
	e.SetData(data)
	return
}
