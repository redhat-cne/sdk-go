package hwevent_test

import (
	"encoding/json"
	"github.com/redhat-cne/sdk-go/pkg/hwevent"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/redhat-cne/sdk-go/pkg/types"
)

func TestUnMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	_type := "hw_fan_event"
	version := "v1"
	id := "ABC-1234"

	testCases := map[string]struct {
		body    []byte
		want    *hwevent.Event
		wantErr error
	}{

		"struct Data fan": {
			body: mustJSONMarshal(t, map[string]interface{}{
				"data": map[string]interface{}{
					"data":  []byte(`{"resource": "/cluster/node/hw", "dataType": "fan", "value": ""}`),
					"version":  version,
				},
				"id":         id,
				"time":       now.Format(time.RFC3339Nano),
				"type":       _type,
				"dataSchema": nil,
			}),
			want: &hwevent.Event{
				ID:         id,
				Type:       _type,
				Time:       &now,
				DataSchema: nil,
				Data: &hwevent.Data{
					Version: version,
					Data: []byte(`{"resource": "/cluster/node/hw", "dataType": "fan", "value": ""}`),
				},
			},
			wantErr: nil,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			got := &hwevent.Event{}
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
