package amqp

import (
	"fmt"
	"github.com/Azure/go-amqp"
	"github.com/redhat-cne/sdk-go/pkg/channel"
	amqp1 "github.com/redhat-cne/sdk-go/pkg/protocol/amqp"
	"log"
	"sync"
)

var (
	instance *AMQP
	once     sync.Once
)

//AMQP exposes amqp api methods
type AMQP struct {
	Router *amqp1.Router
}

//GetAMQPInstance get event instance
func GetAMQPInstance(amqpHost string, DataIn <-chan *channel.DataChan, DataOut chan<- *channel.DataChan, close <-chan bool) (*AMQP, error) {
	once.Do(func() {
		router, err := amqp1.InitServer(amqpHost, DataIn, DataOut, close)
		if err == nil {
			instance = &AMQP{
				Router: router,
			}
		} else {
			log.Printf("error connecting to amqp %v", err)
		}
	})
	if instance == nil || instance.Router == nil {
		return nil, fmt.Errorf("error conecting to amqp")

	}
	if instance.Router.Client == nil {
		client, err := instance.Router.NewClient(amqpHost, []amqp.ConnOption{})
		if err != nil {
			log.Printf("error creating client %v", err)
			return nil, err
		}
		instance.Router.Client = client
	}
	return instance, nil
}

//Start start amqp processors
func (a *AMQP) Start(wg *sync.WaitGroup) {
	go instance.Router.QDRRouter(wg)
}

//NewSender - create new sender independent of the framework
func NewSender(hostName string, port int, address string) (*amqp1.Protocol, error) {
	return amqp1.NewSender(hostName, port, address)
}

// NewReceiver create new receiver independent of the framework
func NewReceiver(hostName string, port int, address string) (*amqp1.Protocol, error) {
	return amqp1.NewReceiver(hostName, port, address)
}

//DeleteSender send publisher address information  on a channel to delete its sender object
func DeleteSender(inChan chan<- *channel.DataChan, address string) {
	// go ahead and create QDR to this address
	inChan <- &channel.DataChan{
		Address: address,
		Type:    channel.SENDER,
		Status:  channel.DELETE,
	}
}

//CreateSender send publisher address information  on a channel to create it's sender object
func CreateSender(inChan chan<- *channel.DataChan, address string) {
	// go ahead and create QDR to this address
	inChan <- &channel.DataChan{
		Address: address,
		Type:    channel.SENDER,
		Status:  channel.NEW,
	}
}

//DeleteListener send subscription address information  on a channel to delete its listener object
func DeleteListener(inChan chan<- *channel.DataChan, address string) {
	// go ahead and create QDR listener to this address
	inChan <- &channel.DataChan{
		Address: address,
		Type:    channel.LISTENER,
		Status:  channel.DELETE,
	}
}

//CreateListener send subscription address information  on a channel to create its listener object
func CreateListener(inChan chan<- *channel.DataChan, address string) {
	// go ahead and create QDR listener to this address
	inChan <- &channel.DataChan{
		Address: address,
		Type:    channel.LISTENER,
		Status:  channel.NEW,
	}
}
