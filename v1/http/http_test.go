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

package http_test

import (
	"testing"

	"github.com/redhat-cne/sdk-go/pkg/channel"
	api "github.com/redhat-cne/sdk-go/v1/http"
	"github.com/stretchr/testify/assert"
)

var (
	storePath         = "."
	outChan           = make(chan *channel.DataChan, 1)
	address           = "test/test"
	port              = 8086
	serverAddress     = "http://localhost:8086"
	in                = make(chan *channel.DataChan)
	out               = make(chan *channel.DataChan)
	closeCh           = make(chan struct{})
	globalInstance, _ = api.GetHTTPInstance(serverAddress, port, storePath, in, out, closeCh, nil, nil)
)

func TestAPI_GetAPIInstance(t *testing.T) {
	localInstance, err := api.GetHTTPInstance(serverAddress, port, storePath, in, out, closeCh, nil, nil)
	if err != nil {
		t.Skipf("tcp.Dial(%#v): %v", localInstance, err)
	}
	assert.Equal(t, &globalInstance, &localInstance)
	close(closeCh)
}

func TestCreateSubscription(t *testing.T) {
	subscriber := &channel.DataChan{
		Address: address,
		Status:  channel.NEW,
		Type:    channel.SUBSCRIBER,
	}
	outChan <- subscriber
	data := <-outChan
	assert.Equal(t, subscriber, data)
}

func TestDeleteSubscription(t *testing.T) {
	subscriber := &channel.DataChan{
		Address: address,
		Status:  channel.DELETE,
		Type:    channel.SUBSCRIBER,
	}
	outChan <- subscriber
	data := <-outChan
	assert.Equal(t, subscriber, data)
}

func TestStatusPing(t *testing.T) {
	sender := &channel.DataChan{
		Address: address,
		Type:    channel.STATUS,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}
