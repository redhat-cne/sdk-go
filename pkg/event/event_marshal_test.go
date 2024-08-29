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

package event_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/redhat-cne/sdk-go/v2/pkg/event"
	"github.com/redhat-cne/sdk-go/v2/pkg/event/ptp"
	"github.com/redhat-cne/sdk-go/v2/pkg/types"
	v1 "github.com/redhat-cne/sdk-go/v2/v1/event"
)

func TestMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	schemaURL := "http://example.com/schema"
	resource := "/cluster/node/ptp"
	_type := string(ptp.PtpStateChange)
	_source := "/cluster/node/example.com/ptp/clock_realtime"
	version := "v1"
	data := event.Data{}
	value := []event.DataValue{{
		Resource:  resource,
		DataType:  event.NOTIFICATION,
		ValueType: event.ENUMERATION,
		Value:     ptp.FREERUN,
	}, {
		Resource:  resource,
		DataType:  event.METRIC,
		ValueType: event.DECIMAL,
		Value:     10.7,
	}}
	data.SetVersion(version) //nolint:errcheck
	data.Values = value      //nolint:errcheck
	fmt.Print(data)

	testCases := map[string]struct {
		event   event.Event
		want    map[string]interface{}
		wantErr *string
	}{
		"empty struct": {
			event: event.Event{},
			wantErr: func() *string {
				s := "json: error calling MarshalJSON for type event.Event: missing event content type\n"
				return &s
			}(),
		},
		"struct Data v1": {
			event: func() event.Event {
				e := v1.CloudNativeEvent()
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
							"data_type":       "notification",
							"value":           "FREERUN",
							"value_type":      "enumeration"},
						map[string]interface{}{
							"ResourceAddress": resource,
							"data_type":       "metric",
							"value":           "10.7",
							"value_type":      "decimal64.3"}},
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

func mustJSONMarshal(tb testing.TB, body interface{}) []byte {
	b, err := json.Marshal(body)
	require.NoError(tb, err)
	return b
}

func assertJSONEquals(t *testing.T, want map[string]interface{}, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
