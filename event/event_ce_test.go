package event_test

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	cne_event "github.com/redhat-cne/sdk-go/event"
	cne_pubsub "github.com/redhat-cne/sdk-go/pubsub"
	"github.com/redhat-cne/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

var (
	now         = types.Timestamp{Time: time.Now().UTC()}
	uriLocation = "http://localhost:9090/api/subscription/1234"
	endPointUri = "http://localhost:8080/api/ack/event"
	resource    = "/cluster/node/ptp"
	_type       = "ptp_status_type"
	version     = "v1"
	id          = uuid.New().String()
	data        cne_event.Data
	pubsub      cne_pubsub.PubSub
)

func setup() {
	data = cne_event.Data{}
	value := cne_event.DataValue{
		Resource:  resource,
		DataType:  cne_event.NOTIFICATION,
		ValueType: cne_event.ENUMERATION,
		Value:     cne_event.GNSS_ACQUIRING_SYNC,
	}
	data.SetVersion(version) //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck
	pubsub = cne_pubsub.PubSub{}
	_ = pubsub.SetResource(resource)
	_ = pubsub.SetURILocation(uriLocation)
	_ = pubsub.SetEndpointURI(endPointUri)
	_ = pubsub.SetID(id)

}
func TestEvent_NewCloudEvent(t *testing.T) {
	setup()
	testCases := map[string]struct {
		cne_event  *cne_event.Event
		cne_pubsub *cne_pubsub.PubSub
		want       *cloudevents.Event
		wantErr    *string
	}{
		"struct Data v1": {
			cne_event: func() *cne_event.Event {
				e := cne_event.NewCloudNativeEvent()
				_ = e.SetDataContentType(cne_event.ApplicationJSON)
				_ = e.SetTime(now.Time)
				_ = e.SetType(_type)
				_ = e.SetData(data)
				return e
			}(),
			cne_pubsub: &pubsub,
			want: func() *cloudevents.Event {
				e := cloudevents.NewEvent()
				e.SetSpecVersion(cloudevents.VersionV03)
				e.SetType(_type)
				_ = e.SetData(cloudevents.ApplicationJSON, data)
				e.SetTime(now.Time)
				e.SetSource(pubsub.GetResource())
				e.SetID(id)
				return &e
			}(),
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			event := tc.cne_event
			ce, err := event.NewCloudEvent(tc.cne_pubsub)
			assert.Nil(t, err)
			tc.want.SetID(ce.ID())
			gotBytes, err := json.Marshal(ce)
			log.Printf("cloud events %s\n",string(gotBytes))
			if tc.wantErr != nil {
				require.Error(t, err, *tc.wantErr)
				return
			}
			assertCEJsonEquals(t, tc.want, gotBytes)
		})
	}
}
func TestEvent_GetCloudNativeEvents(t *testing.T) {
	setup()
	testCases := map[string]struct {
		ce_event *cloudevents.Event
		want     *cne_event.Event
		wantErr  *string
	}{
		"struct Data v1": {
			ce_event: func() *cloudevents.Event {
				e := cloudevents.NewEvent()
				e.SetType(_type)
				_ = e.SetData(cloudevents.ApplicationJSON, data)
				e.SetTime(now.Time)
				e.SetSource(pubsub.GetResource())
				e.SetID(id)
				return &e
			}(),
			want: func() *cne_event.Event {
				e := cne_event.NewCloudNativeEvent()
				_ = e.SetDataContentType(cne_event.ApplicationJSON)
				_ = e.SetTime(now.Time)
				_ = e.SetType(_type)
				_ = e.SetData(data)
				return e
			}(),
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			event := cne_event.NewCloudNativeEvent()
			err := event.GetCloudNativeEvents(tc.ce_event)
			assert.Nil(t, err)
			gotBytes, err := json.Marshal(event)
			if tc.wantErr != nil {
				require.Error(t, err, *tc.wantErr)
				return
			}
			assertCNEJsonEquals(t, tc.want, gotBytes)

		})
	}
}

func assertCEJsonEquals(t *testing.T, want *cloudevents.Event, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}

func assertCNEJsonEquals(t *testing.T, want *cne_event.Event, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
