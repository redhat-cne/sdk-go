package event_test

import (
	"encoding/json"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	cneevent "github.com/redhat-cne/sdk-go/pkg/event"
	cnepubsub "github.com/redhat-cne/sdk-go/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/pkg/types"
	cneeventv1 "github.com/redhat-cne/sdk-go/v1/event"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
)

var (
	now         = types.Timestamp{Time: time.Now().UTC()}
	uriLocation = "http://localhost:9090/event/subscription/1234"
	endPointURI = "http://localhost:8080/event/ack/event"
	resource    = "/cluster/node/ptp"
	_type       = "ptp_status_type"
	version     = "v1"
	id          = uuid.New().String()
	data        cneevent.Data
	pubsub      cnepubsub.PubSub
)

func setup() {
	data = cneevent.Data{}
	value := cneevent.DataValue{
		Resource:  resource,
		DataType:  cneevent.NOTIFICATION,
		ValueType: cneevent.ENUMERATION,
		Value:     cneevent.ACQUIRING_SYNC,
	}
	data.SetVersion(version) //nolint:errcheck
	data.AppendValues(value) //nolint:errcheck
	pubsub = cnepubsub.PubSub{}
	_ = pubsub.SetResource(resource)
	_ = pubsub.SetURILocation(uriLocation)
	_ = pubsub.SetEndpointURI(endPointURI)
	pubsub.SetID(id)

}

func TestEvent_NewCloudEvent(t *testing.T) {
	setup()
	testCases := map[string]struct {
		cneEvent  *cneevent.Event
		cnePubsub *cnepubsub.PubSub
		want      *ce.Event
		wantErr   *string
	}{
		"struct Data v1": {
			cneEvent: func() *cneevent.Event {
				e := cneeventv1.CloudNativeEvent()

				e.SetDataContentType(cneevent.ApplicationJSON)
				e.SetTime(now.Time)
				e.SetType(_type)
				e.SetData(data)
				return &e
			}(),
			cnePubsub: &pubsub,
			want: func() *ce.Event {
				e := ce.NewEvent()
				e.SetSpecVersion(ce.VersionV03)
				e.SetType(_type)
				_ = e.SetData(ce.ApplicationJSON, data)
				e.SetTime(now.Time)
				e.SetSource(pubsub.GetResource())
				e.SetID(id)
				return &e
			}(),
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			event := tc.cneEvent
			cEvent, err := event.NewCloudEvent(tc.cnePubsub)
			assert.Nil(t, err)
			tc.want.SetID(cEvent.ID())
			gotBytes, err := json.Marshal(cEvent)
			log.Printf("cloud events %s\n", string(gotBytes))
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
		ceEvent *ce.Event
		want    *cneevent.Event
		wantErr *string
	}{
		"struct Data v1": {
			ceEvent: func() *ce.Event {
				e := ce.NewEvent()
				e.SetType(_type)
				_ = e.SetData(ce.ApplicationJSON, data)
				e.SetTime(now.Time)
				e.SetSource(pubsub.GetResource())
				e.SetID(id)
				return &e
			}(),
			want: func() *cneevent.Event {
				e := cneeventv1.CloudNativeEvent()
				e.SetDataContentType(cneevent.ApplicationJSON)
				e.SetTime(now.Time)
				e.SetType(_type)
				e.SetData(data)
				return &e
			}(),
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			event := cneeventv1.CloudNativeEvent()
			err := event.GetCloudNativeEvents(tc.ceEvent)
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

func assertCEJsonEquals(t *testing.T, want *ce.Event, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}

func assertCNEJsonEquals(t *testing.T, want *cneevent.Event, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
