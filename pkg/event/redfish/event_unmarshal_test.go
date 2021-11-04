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

	"github.com/redhat-cne/sdk-go/pkg/event"
	"github.com/redhat-cne/sdk-go/pkg/event/redfish"
	"github.com/stretchr/testify/require"

	"github.com/google/go-cmp/cmp"
	"github.com/redhat-cne/sdk-go/pkg/types"
)

func TestUnMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	_type := string(redfish.Alert)
	version := "v1"
	id := "ABC-1234"

	testCases := map[string]struct {
		body    []byte
		want    *event.Event
		wantErr error
	}{

		"struct Data fan": {
			body: mustJSONMarshal(t, map[string]interface{}{
				"data": map[string]interface{}{
					"resource": resource,
					"values": []interface{}{
						map[string]interface{}{
							"resource":  resource,
							"dataType":  event.NOTIFICATION,
							"value":     JSON_EVENT_TMP0100,
							"valueType": event.REDFISH_EVENT}},
					"version": version,
				},
				"id":         id,
				"time":       now.Format(time.RFC3339Nano),
				"type":       _type,
				"dataSchema": nil,
			}),
			want: &event.Event{
				ID:         id,
				Type:       _type,
				Time:       &now,
				DataSchema: nil,
				Data: &event.Data{
					Version: version,
					Values: []event.DataValue{
						{
							Resource:  resource,
							DataType:  event.NOTIFICATION,
							ValueType: event.REDFISH_EVENT,
							Value:     REDFISH_EVENT_TMP0100,
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			got := &event.Event{}
			err := json.Unmarshal(tc.body, got)

			if tc.wantErr != nil || err != nil {
				if diff := cmp.Diff(tc.wantErr, err); diff != "" {
					t.Errorf("unexpected error (-want, +got) = %v", diff)
				}
				return
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected event (-want, +got) = %v", diff)
			}
		})
	}
}

func mustJSONMarshal(tb testing.TB, body interface{}) []byte {
	b, err := json.Marshal(body)
	require.NoError(tb, err)
	return b
}
