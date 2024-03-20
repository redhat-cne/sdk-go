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

package amqp_test

import (
	"testing"
	"time"

	"github.com/redhat-cne/sdk-go/pkg/channel"
	api "github.com/redhat-cne/sdk-go/v1/amqp"
	"github.com/stretchr/testify/assert"
)

var (
	outChan           = make(chan *channel.DataChan, 1)
	address           = "test/test"
	s                 = "amqp://localhost:5672"
	in                = make(chan *channel.DataChan)
	out               = make(chan *channel.DataChan)
	closeCh           = make(chan struct{})
	timeout           = 1 * time.Second
	globalInstance, _ = api.GetAMQPInstance(s, in, out, closeCh, timeout)
)

func TestAPI_GetAPIInstance(t *testing.T) {
	localInstance, err := api.GetAMQPInstance(s, in, out, closeCh, timeout)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", localInstance, err)
	}
	assert.Equal(t, &globalInstance, &localInstance)
}

func TestCreateSender(t *testing.T) {
	sender := &channel.DataChan{
		Address: address,
		Status:  channel.NEW,
		Type:    channel.PUBLISHER,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}

func TestDeleteSender(t *testing.T) {
	sender := &channel.DataChan{
		Address: address,
		Status:  channel.DELETE,
		Type:    channel.PUBLISHER,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}

func TestDeleteListener(t *testing.T) {
	sender := &channel.DataChan{
		Address: address,
		Status:  channel.DELETE,
		Type:    channel.SUBSCRIBER,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}
