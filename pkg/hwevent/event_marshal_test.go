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

package hwevent_test

import (
	"encoding/json"
	"github.com/redhat-cne/sdk-go/pkg/hwevent"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/redhat-cne/sdk-go/pkg/event"
	"github.com/redhat-cne/sdk-go/pkg/types"
)

func TestMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	schemaURL := "http://example.com/schema"
	_type := "hw_fan_type"
	version := "v1"
	data := hwevent.Data{}
	data.SetVersion(version) //nolint:errcheck
	data.Data = []byte(`{"resource": "cluster/node/hw","dataType": "notification","value": "ACQUIRING-SYNC", "valueType": "enumeration" }`)

	testCases := map[string]struct {
		event   hwevent.Event
		want    map[string]interface{}
		wantErr *string
	}{
		"struct Data v1": {
			event: func() hwevent.Event {
				e := hwevent.Event{Type: "Event"}
				e.SetDataContentType(event.ApplicationJSON)
				_ = e.SetDataSchema(schemaURL)
				e.Time = &now
				e.SetType(_type)
				e.SetData(data)
				return e
			}(),
			want: map[string]interface{}{
				"dataContentType": "application/json",
				"data": map[string]interface{}{
					"data": []byte(`{"resource":
					"cluster/node/hw",
					"dataType": "notification",
					"value":     "ACQUIRING-SYNC",
					"valueType": "enumeration"}`),
					"version" :"v1",

				},
				"id":              "",
				"time":            now.Format(time.RFC3339Nano),
				"type":            _type,
				"dataSchema":      schemaURL,
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
	//var gotToCompare map[string]interface{}
	gotToCompare:= hwevent.Event{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	 wantToCompare:=  hwevent.Event{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
