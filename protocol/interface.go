package protocol

import (
	"context"
	amqp1 "github.com/cloudevents/sdk-go/protocol/amqp/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/redhat-cne/sdk-go/channel"
)

//Protocol ...
type Protocol struct {
	ID            string
	MsgCount      int
	Protocol      *amqp1.Protocol
	Ctx           context.Context
	ParentContext context.Context
	CancelFn      context.CancelFunc
	Client        cloudevents.Client
	Queue         string
	DataIn        <-chan channel.DataEvent
	DataOut       chan<- channel.DataEvent
}
