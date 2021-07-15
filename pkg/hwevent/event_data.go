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
	"fmt"
)

// The structs are defined based on Event v1.4, which is part of
// Redfish Schema Bundle 2019.2.
// Reference: https://redfish.dmtf.org/schemas/v1/Event.v1_4_0.yaml

// EventRecord is defined in Redfish Event_v1_4_0_EventRecord
type EventRecord struct {
	// The Actions property shall contain the available actions
	// for this resource.
	Actions []string `json:"Actions"`
	// *deprecated* This property has been Deprecated in favor of Context
	// found at the root level of the object.
	Context string `json:"Context"`
	// This value is the identifier used to correlate events that
	// came from the same cause.
	EventGroupId int `json:"EventGroupId"`
	// The value of this property shall indicate a unique identifier
	// for the event, the format of which is implementation dependent.
	EventID string `json:"EventId"`
	// This is time the event occurred.
	EventTimestamp string `json:"EventTimestamp"`
	// *deprecated* This property has been deprecated.  Starting Redfish
	// Spec 1.6 (Event 1.3), subscriptions are based on RegistryId and ResourceType
	// and not EventType.
	// This indicates the type of event sent, according to the definitions
	// in the EventService.
	EventType string `json:"EventType"`
	// This is the identifier for the member within the collection.
	MemberID string `json:"MemberId"`
	// This property shall contain an optional human readable
	// message.
	Message string `json:"Message"`
	// This array of message arguments are substituted for the arguments
	// in the message when looked up in the message registry.
	MessageArgs []string `json:"MessageArgs"`
	// This property shall be a key into message registry as
	// described in the Redfish specification.
	MessageID string `json:"MessageId"`
	//  This is the severity of the event.
	Severity string `json:"Severity"`
}

// RedfishEvent Event_v1_4_0_Event
// The Event schema describes the JSON payload received by an Event
// Destination (which has subscribed to event notification) when events occurs.  This
// resource contains data about event(s), including descriptions, severity and
// MessageId reference to a Message Registry that can be accessed for further
// information.
type RedfishEvent struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	// This property shall contain a client supplied context
	// for the Event Destination to which this event is being sent.
	Context string `json:"Context"`
	Events  []struct {
	} `json:"Events"`
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// Data ... cloud native events data
// Data Json payload is as follows,
//{
//	"version": "v1.0",
//
//}
type Data struct {
	Version string `json:"version" example:"v1"`
	Data    []byte `json:"data"`
}

// SetVersion  ...
func (d *Data) SetVersion(s string) error {
	d.Version = s
	if s == "" {
		err := fmt.Errorf("version cannot be empty")
		return err
	}
	return nil
}

// GetVersion ...
func (d *Data) GetVersion() string {
	return d.Version
}

// SetData ...
func (d *Data) SetData(b []byte) {
	d.Data = b
}
