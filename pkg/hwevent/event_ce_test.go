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

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	hwevent "github.com/redhat-cne/sdk-go/pkg/hwevent"
	cnepubsub "github.com/redhat-cne/sdk-go/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/pkg/types"
	hweventv1 "github.com/redhat-cne/sdk-go/v1/hwevent"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
)

var (
	now    = types.Timestamp{Time: time.Now().UTC()}
	_type  = "HW_EVENT"
	id     = uuid.New().String()
	data   hwevent.Data
	pubsub cnepubsub.PubSub
)

func setup() {
	data = hwevent.Data{}
	data.SetData([]byte(`{"EventId":"TestEventId","EventTimestamp":"2019-07-29T15:13:49Z","EventType":"Alert","Message":"Test Event","MessageArgs":["NoAMS","Busy","Cached"],"MessageId":"iLOEvents.2.1.ServerPoweredOff","OriginOfCondition":"/redfish/v1/Systems/1/","Severity":"OK"}`)) //nolint:errcheck
	pubsub.SetID(id)

}

func TestEvent_NewCloudEvent(t *testing.T) {
	setup()
	testCases := map[string]struct {
		hwEvent   *hwevent.Event
		cnePubsub *cnepubsub.PubSub
		want      *ce.Event
		wantErr   *string
	}{
		"struct Data v1": {
			hwEvent: func() *hwevent.Event {
				e := hweventv1.CloudNativeEvent()
				e.SetDataContentType(hwevent.ApplicationJSON)
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
			event := tc.hwEvent
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
		want    *hwevent.Event
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
			want: func() *hwevent.Event {
				e := hweventv1.CloudNativeEvent()
				e.SetDataContentType(hwevent.ApplicationJSON)
				e.SetTime(now.Time)
				e.SetType(_type)
				e.SetData(data)
				return &e
			}(),
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			event := hweventv1.CloudNativeEvent()
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

func assertCNEJsonEquals(t *testing.T, want *hwevent.Event, got []byte) {
	var gotToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &gotToCompare))

	// Marshal and unmarshal want to make sure the types are correct
	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)
	var wantToCompare map[string]interface{}
	require.NoError(t, json.Unmarshal(wantBytes, &wantToCompare))

	require.Equal(t, wantToCompare, gotToCompare)
}
