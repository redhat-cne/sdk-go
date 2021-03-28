package protocol

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/redhat-cne/sdk-go/pkg/channel"
)

//Binder ...protocol binder base struct
type Binder struct {
	ID            string
	Ctx           context.Context
	ParentContext context.Context
	CancelFn      context.CancelFunc
	Client        cloudevents.Client
	// Address of the protocol
	Address string
	//DataIn data coming in to this protocol
	DataIn <-chan channel.DataChan
	//DataOut data coming out of this protocol
	DataOut chan<- channel.DataChan
	//close on true
	Close    <-chan bool
	Protocol interface{}
}
