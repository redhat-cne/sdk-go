package amqp_test

import (
	"github.com/redhat-cne/sdk-go/pkg/channel"
	api "github.com/redhat-cne/sdk-go/v1/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	outChan           = make(chan channel.DataChan, 1)
	address           = "test/test"
	s                 = "amqp://localhost:5672"
	in                = make(chan channel.DataChan)
	out               = make(chan channel.DataChan)
	close             = make(chan bool)
	globalInstance, _ = api.GetAMQPInstance(s, in, out, close)
)

func TestAPI_GetAPIInstance(t *testing.T) {

	localInstance, err := api.GetAMQPInstance(s, in, out, close)
	if err != nil {
		t.Skipf("ampq.Dial(%#v): %v", localInstance, err)
	}

	assert.Equal(t, &globalInstance, &localInstance)
}

func TestCreateSender(t *testing.T) {
	sender := channel.DataChan{
		Address: address,
		Status:  channel.NEW,
		Type:    channel.SENDER,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}

func TestDeleteSender(t *testing.T) {
	sender := channel.DataChan{
		Address: address,
		Status:  channel.DELETE,
		Type:    channel.SENDER,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}

func TestDeleteListener(t *testing.T) {
	sender := channel.DataChan{
		Address: address,
		Status:  channel.DELETE,
		Type:    channel.LISTENER,
	}
	outChan <- sender
	data := <-outChan
	assert.Equal(t, sender, data)
}
