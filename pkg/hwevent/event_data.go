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
	"strings"
)

// The structs efined here are based on Event v1.4, which is part of
// Redfish Schema Bundle 2019.2.
// Reference: https://redfish.dmtf.org/schemas/v1/Event.v1_4_1.json

// EventRecord is defined in Redfish Event_v1_4_1_EventRecord
// Required fields: EventType, MessageId, MemberId
type EventRecord struct {
	// The Actions property shall contain the available actions
	// for this resource.
	Actions []byte `json:"Actions,omitempty"`
	// *deprecated* This property has been Deprecated in favor of Context
	// found at the root level of the object.
	Context string `json:"Context,omitempty"`
	// This value is the identifier used to correlate events that
	// came from the same cause.
	EventGroupID int `json:"EventGroupId,omitempty"`
	// The value of this property shall indicate a unique identifier
	// for the event, the format of which is implementation dependent.
	EventID string `json:"EventId,omitempty"`
	// This is time the event occurred.
	EventTimestamp string `json:"EventTimestamp,omitempty"`
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
	Message string `json:"Message,omitempty"`
	// This array of message arguments are substituted for the arguments
	// in the message when looked up in the message registry.
	MessageArgs []string `json:"MessageArgs,omitempty"`
	// This property shall be a key into message registry as
	// described in the Redfish specification.
	MessageID string `json:"MessageId"`
	// This is the manufacturer/provider specific extension
	Oem []byte `json:"Oem,omitempty"`
	// This indicates the resource that originated the condition that
	// caused the event to be generated.
	OriginOfCondition string `json:"OriginOfCondition,omitempty"`
	//  This is the severity of the event.
	Severity string `json:"Severity,omitempty"`
}

// String returns a pretty-printed representation of the EventRecord.
func (e EventRecord) String() string {
	b := strings.Builder{}
	if e.Actions != nil {
		b.WriteString("Actions: " + string(e.Actions.Oem) + "\n")
	}
	if e.Context != "" {
		b.WriteString("Context: " + e.Context + "\n")
	}
	b.WriteString("EventGroupId: " + string(e.EventGroupID) + "\n")
	if e.EventID != "" {
		b.WriteString("EventId: " + e.EventID + "\n")
	}
	if e.EventTimestamp != "" {
		b.WriteString("EventTimestamp: " + e.EventTimestamp + "\n")
	}
	b.WriteString("EventType: " + e.EventType + "\n")
	b.WriteString("MemberId: " + e.MemberID + "\n")
	if e.Message != "" {
		b.WriteString("Message: " + e.Message + "\n")
	}
	if e.MessageArgs != nil {
		b.WriteString("MessageArgs: ")
		for _, arg := range e.MessageArgs {
			b.WriteString(arg + ", ")
		}
		b.WriteString("\n")
	}
	b.WriteString("MessageId: " + e.MessageID + "\n")
	if e.Oem != nil {
		b.WriteString("Oem: " + string(e.Oem) + "\n")
	}
	return b.String()
}

// RedfishEvent Event_v1_4_1_Event
// The Event schema describes the JSON payload received by an Event
// Destination (which has subscribed to event notification) when events occurs.  This
// resource contains data about event(s), including descriptions, severity and
// MessageId reference to a Message Registry that can be accessed for further
// information.
// Required fields: @odata.type, Events, Id, Name
type RedfishEvent struct {
	OdataContext string `json:"@odata.context,omitempty"`
	OdataType    string `json:"@odata.type"`
	// The available actions for this resource.
	Actions []byte `json:"Actions,omitempty"`
	// A context can be supplied at subscription time.  This property
	// is the context value supplied by the subscriber.
	Context     string        `json:"Context,omitempty"`
	Description string        `json:"Description,omitempty"`
	Events      []EventRecord `json:"Events"`
	ID          string        `json:"Id"`
	Name        string        `json:"Name"`
	// This is the manufacturer/provider specific extension
	Oem []byte `json:"Oem,omitempty"`
}

// String returns a pretty-printed representation of the RedfishEvent.
func (e RedfishEvent) String() string {
	b := strings.Builder{}
	if e.OdataContext != "" {
		b.WriteString("@odata.context: " + e.OdataContext + "\n")
	}
	b.WriteString("@odata.type: " + e.OdataType + "\n")
	if e.Actions != nil {
		b.WriteString("Actions: " + string(e.Actions.Oem) + "\n")
	}
	if e.Context != "" {
		b.WriteString("Context: " + e.Context + "\n")
	}
	b.WriteString("Id: " + e.ID + "\n")
	b.WriteString("Name: " + e.Name + "\n")
	if e.Oem != nil {
		b.WriteString("Oem: " + string(e.Oem) + "\n")
	}
	for i, e := range e.Events {
		b.WriteString("Events[" + string(i) + "]:\n")
		b.WriteString(e.String())
	}
	return b.String()
}

// Data ... cloud native events data
// Data Json payload is as follows,
//{
//	"version": "v1.4.0",
//
//}
type Data struct {
	Version string       `json:"version" example:"v1"`
	Data    RedfishEvent `json:"data"`
}

// String returns a pretty-printed representation of the Data.
func (d Data) String() string {
	b := strings.Builder{}
	b.WriteString("version: " + d.Version + "\n")
	b.WriteString("data:\n")
	b.WriteString(d.Data.String())
	return b.String()
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
func (d *Data) SetData(b RedfishEvent) {
	d.Data = b
}
