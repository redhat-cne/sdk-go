package event_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/redhat-cne/sdk-go/pkg/event"
	"github.com/redhat-cne/sdk-go/pkg/types"
	v1 "github.com/redhat-cne/sdk-go/v1/event"
)

func TestMarshal(t *testing.T) {
	now := types.Timestamp{Time: time.Now().UTC()}
	schemaUrl := "http://example.com/schema"
	resource := "/cluster/node/ptp"
	_type := "ptp_status_type"
	version := "v1"
	data := event.Data{}
	value := event.DataValue{
		Resource:  resource,
		DataType:  event.NOTIFICATION,
		ValueType: event.ENUMERATION,
		Value:     event.GNSS_ACQUIRING_SYNC,
	}
	data.SetVersion(version) //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck

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
				_ = e.SetDataSchema(schemaUrl)
				e.Time = &now
				e.SetType(_type)
				e.SetData(data)

				return e
			}(),
			want: map[string]interface{}{
				"dataContentType": "application/json",
				"data": map[string]interface{}{
					"values":  []interface{}{map[string]interface{}{"resource": resource, "dataType": "notification", "value": "ACQUIRING-SYNC", "valueType": "enumeration"}},
					"version": "v1",
				},
				"id":         "",
				"time":       now.Format(time.RFC3339Nano),
				"type":       _type,
				"dataSchema": schemaUrl,
			},
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			event := tc.event
			gotBytes, err := json.Marshal(event)
			if tc.wantErr != nil {
				require.Error(t, err, *tc.wantErr)
				return
			}
			assertJsonEquals(t, tc.want, gotBytes)
		})
	}
}

func mustJsonMarshal(tb testing.TB, body interface{}) []byte {
	b, err := json.Marshal(body)
	require.NoError(tb, err)
	return b
}

func assertJsonEquals(t *testing.T, want map[string]interface{}, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
