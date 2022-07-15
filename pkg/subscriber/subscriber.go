package subscriber

// Copyright 2022 The Cloud Native Events Authors
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

import (
	"strings"
	"sync"

	"github.com/redhat-cne/sdk-go/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/pkg/store"
	"github.com/redhat-cne/sdk-go/pkg/types"
)

// Status of teh client connections
type Status int64

const (
	// InActive client
	InActive Status = iota
	// Active Client
	Active
)

// Subscriber object holds client connections
type Subscriber struct {
	// ClientID of the sub
	// +required
	ClientID string `json:"clientID" omit:"empty"`
	// +required
	SubStore *store.PubSubStore `json:"subStore" omit:"empty"`
	// HealthEndPoint - A URI describing the
	//producer/subscription get link.
	HealthEndPoint *types.URI `json:"healthEndPoint" omit:"empty"`

	Status Status `json:"status" omit:"empty"`
}

// String returns a pretty-printed representation of the Event.
func (s *Subscriber) String() string {
	b := strings.Builder{}

	b.WriteString("  HealthEndPoint: " + s.GetURILocation() + "\n")
	b.WriteString("  ID: " + s.GetClientID() + "\n")
	b.WriteString("  sub :{")
	for _, v := range s.SubStore.Store {
		b.WriteString(" {")
		b.WriteString(v.String() + "\n")
		b.WriteString(" },")
	}
	b.WriteString(" }")
	return b.String()
}

// New create new subscriber
func New(clientID string, healthEndPoint *types.URI) *Subscriber {

	return &Subscriber{
		ClientID: clientID,
		SubStore: &store.PubSubStore{
			RWMutex: sync.RWMutex{},
			Store:   map[string]*pubsub.PubSub{},
		},
		HealthEndPoint: healthEndPoint,
		Status:         0,
	}
}
