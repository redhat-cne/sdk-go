// Copyright 2021 The Cloud Native Events Authors
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

package redfish_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/redhat-cne/sdk-go/v2/pkg/event"
	"github.com/redhat-cne/sdk-go/v2/pkg/event/redfish"
	"github.com/redhat-cne/sdk-go/v2/pkg/types"
	"github.com/stretchr/testify/require"
)

var (
	JSON_EVENT_TMP0100 = map[string]interface{}{
		"@odata.context": "/redfish/v1/$metadata#Event.Event",
		"@odata.id":      "/redfish/v1/EventService/Events/5e004f5a-e3d1-11eb-ae9c-3448edf18a38",
		"@odata.type":    "#Event.v1_3_0.Event",
		"Context":        "any string is valid",
		"Events": []interface{}{
			map[string]interface{}{
				"Context":                 "any string is valid",
				"EventId":                 "2162",
				"EventTimestamp":          "2021-07-13T15:07:59+0300",
				"EventType":               "Alert",
				"MemberId":                "615703",
				"Message":                 "The system board Inlet temperature is less than the lower warning threshold.",
				"MessageArgs":             []string{"Inlet"},
				"MessageArgs@odata.count": 1,
				"MessageId":               "TMP0100",
				// Do not use []byte here since json.Marshal from standard library will encode []byte with base64
				"OriginOfCondition": map[string]interface{}{
					"@odata.id": "/redfish/v1/Systems/System.Embedded.1",
				},
				"Severity": "Warning",
			},
		},
		"Id":   "5e004f5a-e3d1-11eb-ae9c-3448edf18a38",
		"Name": "Event Array",
		"Oem": map[string]interface{}{
			"Dell": map[string]interface{}{
				"@odata.type":    "#DellEvent.v1_0_0.DellEvent",
				"ServerHostname": "",
			},
		},
	}
)

func TestMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	schemaURL := "http://example.com/schema"
	_type := string(redfish.Alert)
	_source := "/cluster/node/nodename/redfish/event"
	version := "v1"
	data := event.Data{}
	value := event.DataValue{
		Resource:  resource,
		DataType:  event.NOTIFICATION,
		ValueType: event.REDFISH_EVENT,
		Value:     REDFISH_EVENT_TMP0100,
	}
	data.SetVersion(version) //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

	testCases := map[string]struct {
		event   event.Event
		want    map[string]interface{}
		wantErr *string
	}{
		"struct Data v1": {
			event: func() event.Event {
				e := event.Event{Type: string(redfish.Alert)}
				e.SetDataContentType(event.ApplicationJSON)
				_ = e.SetDataSchema(schemaURL)
				e.Time = &now
				e.SetType(_type)
				e.SetSource(_source)
				e.SetData(data)
				return e
			}(),
			want: map[string]interface{}{
				"dataContentType": "application/json",
				"data": map[string]interface{}{
					"values": []interface{}{
						map[string]interface{}{
							"ResourceAddress": resource,
							"data_type":       event.NOTIFICATION,
							// NOTE: Marshal results in compact JSON format without whitespaces
							"value":      JSON_EVENT_TMP0100,
							"value_type": event.REDFISH_EVENT,
						},
					},
					"version": "v1",
				},
				"id":         "",
				"time":       now.Format(time.RFC3339Nano),
				"type":       _type,
				"source":     _source,
				"dataSchema": schemaURL,
			},
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			e := tc.event
			gotBytes, err := json.Marshal(e)
			if tc.wantErr != nil {
				require.Error(t, err, *tc.wantErr)
				return
			}
			assertJSONEquals(t, tc.want, gotBytes)
		})
	}
}

func assertJSONEquals(t *testing.T, want map[string]interface{}, got []byte) {
	gotToCompare := event.Event{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal `want` to make sure the types are correct
	// NOTE: json.Marshal from the standard `encoding/json` library is used here
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	wantToCompare := event.Event{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
