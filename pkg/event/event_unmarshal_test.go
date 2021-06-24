package event_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/redhat-cne/sdk-go/pkg/event"
	"github.com/redhat-cne/sdk-go/pkg/types"
)

func TestUnMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	resource := "/cluster/node/ptp"
	_type := "ptp_status_type"
	version := "v1"
	id := "ABC-1234"

	testCases := map[string]struct {
		body    []byte
		want    *event.Event
		wantErr error
	}{

		"struct Data notification": {
			body: mustJSONMarshal(t, map[string]interface{}{
				"data": map[string]interface{}{
					"resource": resource,
					"values": []interface{}{
						map[string]interface{}{
							"resource":  resource,
							"dataType":  "notification",
							"value":     "FREERUN",
							"valueType": "enumeration"},
						map[string]interface{}{
							"resource":  resource,
							"dataType":  "metric",
							"value":     "10.63",
							"valueType": "decimal64.3"}},
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
							ValueType: event.ENUMERATION,
							Value:     event.FREERUN,
						},
						{
							Resource:  resource,
							DataType:  event.METRIC,
							ValueType: event.DECIMAL,
							Value:     10.63,
						},
					},
				},
			},
			wantErr: nil,
		},
		"struct Data metric": {
			body: mustJSONMarshal(t, map[string]interface{}{
				"data": map[string]interface{}{
					"values": []interface{}{
						map[string]interface{}{
							"resource":  resource,
							"dataType":  "notification",
							"value":     "FREERUN",
							"valueType": "enumeration"},
						map[string]interface{}{
							"resource":  resource,
							"dataType":  "metric",
							"value":     "64.3",
							"valueType": "decimal64.3"}},
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
							ValueType: event.ENUMERATION,
							Value:     event.FREERUN,
						},
						{
							Resource:  resource,
							DataType:  event.METRIC,
							ValueType: event.DECIMAL,
							Value:     64.3,
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
