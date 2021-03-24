package amqp

import (
	"context"
	"fmt"
	"github.com/Azure/go-amqp"
	amqp1 "github.com/cloudevents/sdk-go/protocol/amqp/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	channel "github.com/redhat-cne/sdk-go/channel"
	"github.com/redhat-cne/sdk-go/protocol"
	"log"
	"sync"
	"time"
)

/*var (
	_ protocol.Protocol = (*Router)(nil)
)*/

//Protocol ...
type Protocol struct {
	protocol.Binder
	Protocol *amqp1.Protocol
}

//Router defines QDR router object
type Router struct {
	Listeners map[string]*Protocol
	Senders   map[string]*Protocol
	AMQPHost  string
	DataIn    <-chan channel.DataEvent
	DataOut   chan<- channel.DataEvent
	Client    *amqp.Client
	//close on true
	Close <-chan bool
}

//InitServer initialize QDR configurations
func InitServer(amqpHost string, DataIn <-chan channel.DataEvent, DataOut chan<- channel.DataEvent, close <-chan bool) (*Router, error) {
	server := Router{
		Listeners: map[string]*Protocol{},
		Senders:   map[string]*Protocol{},
		DataIn:    DataIn,
		AMQPHost:  amqpHost,
		DataOut:   DataOut,
		Close:     close,
	}
	client, err := server.NewClient(amqpHost, []amqp.ConnOption{})
	if err != nil {
		return nil, err
	}
	server.Client = client
	return &server, nil
}

// NewClient ...
func (q *Router) NewClient(server string, connOption []amqp.ConnOption) (*amqp.Client, error) {
	client, err := amqp.Dial(server, connOption...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

//NewSession Open a session
func (q *Router) NewSession(sessionOption []amqp.SessionOption) (*amqp.Session, error) {
	session, err := q.Client.NewSession(sessionOption...)
	if err != nil {
		return session, err
	}
	return session, nil
}

//NewSender creates new QDR ptp
func (q *Router) NewSender(address string) error {
	var opts []amqp1.Option
	//p, err := amqp1.NewSenderProtocol(q.AMQPHost, address, []amqp.ConnOption{}, []amqp.SessionOption{}, opts...)
	session, err := q.NewSession([]amqp.SessionOption{})
	if err != nil {
		log.Printf("failed to create an amqp session for a sender : %v", err)
		return err
	}
	p, err := amqp1.NewSenderProtocolFromClient(q.Client, session, address, opts...)
	if err != nil {
		log.Printf("failed to create an amqp sender protocol: %v", err)
		return err
	}

	l := Protocol{}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create an amqp sender client: %v", err)
	}
	l.Protocol = p
	l.Client = c
	q.Senders[address] = &l

	return nil
}

//NewReceiver creates new QDR receiver
func (q *Router) NewReceiver(address string) error {
	var opts []amqp1.Option
	opts = append(opts, amqp1.WithReceiverLinkOption(amqp.LinkCredit(50)))

	session, err := q.NewSession([]amqp.SessionOption{})

	if err != nil {
		log.Printf("failed to create an amqp session for a sender : %v", err)
		return err
	}

	p, err := amqp1.NewReceiverProtocolFromClient(q.Client, session, address, opts...)
	if err != nil {
		log.Printf("failed to create an amqp protocol for a receiver: %v", err)
		return err
	}
	log.Printf("(mew receiver) router connection established %s\n", address)

	l := Protocol{}
	parent, cancelParent := context.WithCancel(context.Background())
	l.CancelFn = cancelParent
	l.ParentContext = parent
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create a receiver client: %v", err)
	}
	l.Protocol = p
	l.Client = c
	q.Listeners[address] = &l
	return nil
}

//Receive is a QDR receiver listening to a address specified
func (q *Router) Receive(wg *sync.WaitGroup, address string, fn func(e cloudevents.Event)) {
	var err error
	defer wg.Done()
	if val, ok := q.Listeners[address]; ok {
		log.Printf("waiting and listening at  %s\n", address)
		err = val.Client.StartReceiver(val.ParentContext, fn)
		if err != nil {
			log.Printf("amqp receiver error: %v", err)
		}
	} else {
		log.Printf("amqp receiver not found in the list\n")
	}
}

//ReceiveAll creates receiver to all address and receives events for all address
func (q *Router) ReceiveAll(wg *sync.WaitGroup, fn func(e cloudevents.Event)) {
	defer wg.Done()
	var err error
	for _, l := range q.Listeners {
		wg.Add(1)
		go func(l *Protocol, wg *sync.WaitGroup) {
			fmt.Printf("listenining to queue %s by %s\n", l.Address, l.ID)
			defer wg.Done()
			err = l.Client.StartReceiver(context.Background(), fn)
			if err != nil {
				log.Printf("amqp receiver error: %v", err)
			}
		}(l, wg)
	}

}

//QDRRouter the QDR Server listens  on data and do action either create sender or receivers
//QDRRouter is qpid router object configured to create publishers and  consumers
func (q *Router) QDRRouter(wg *sync.WaitGroup) {
	wg.Add(1)
	go func(q *Router, wg *sync.WaitGroup) {
		defer wg.Done()
		for { //nolint:gosimple
			select {
			case d := <-q.DataIn:
				if d.EventType == channel.STATUS {
					if d.Address != "" && d.Data.Data() == nil { //create new address protocol for checking status
						if _, ok := q.Listeners[d.Address]; !ok {
							log.Printf("(1)Listener not found for the following address %s , creating listener", d.Address)
							err := q.NewReceiver(d.Address)
							if err != nil {
								log.Printf("Error creating Receiver %v", err)
							} else {
								wg.Add(1)
								go q.Receive(wg, d.Address, func(e cloudevents.Event) { // just spawn and forget
									q.DataOut <- channel.DataEvent{
										Address:     d.Address,
										Data:        &e,
										EventStatus: d.EventStatus,
										EndPointURI: "",
										StatusCh:    d.StatusCh,
										EventType:   channel.STATUS,
									}
								})
							}
						}
					} else {
						log.Printf("Got empty %s", string(d.Data.Data()))
					}
				} else if d.EventType == channel.LISTENER {
					// create receiver and let it run
					if _, ok := q.Listeners[d.Address]; !ok {
						log.Printf("(1)Listener not found for the following address %s , creating listener", d.Address)
						err := q.NewReceiver(d.Address)
						if err != nil {
							log.Printf("Error creating Receiver %v", err)
						} else {
							wg.Add(1)
							go q.Receive(wg, d.Address, func(e cloudevents.Event) {
								log.Printf("Got data ... processing")
								q.DataOut <- channel.DataEvent{
									Address:     d.Address,
									Data:        &e,
									EventStatus: channel.NEW,
									EventType:   channel.EVENT,
								}
							})
							log.Printf("Done setting up receiver for consumer")
						}
					} else {
						log.Printf("(1)Listener already found so not creating again %s\n", d.Address)
					}
				} else if d.EventType == channel.SENDER {
					//log.Printf("reading from data %s", d.Address)
					if _, ok := q.Senders[d.Address]; !ok {
						log.Printf("(1)Sender not found for the following address %s", d.Address)
						err := q.NewSender(d.Address)
						if err != nil {
							log.Printf("(1)error creating sender %v for address %s", err, d.Address)
						}
					} else {
						log.Printf("(1)Sender already found so not creating again %s\n", d.Address)
					}
				} else if d.EventType == channel.EVENT {
					if _, ok := q.Senders[d.Address]; ok {
						q.SendTo(wg, d.Address, d.Data)
					}
				}
			case <-q.Close:
				log.Printf("Recieved close order")
				return
			}
		}
	}(q, wg)
}

//SendTo sends events to the address specified
func (q *Router) SendTo(wg *sync.WaitGroup, address string, event *cloudevents.Event) {
	if sender, ok := q.Senders[address]; ok {
		wg.Add(1) //for each ptp you send a message since its
		go func(s *Protocol, e *cloudevents.Event, wg *sync.WaitGroup) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
			defer cancel()
			if result := sender.Client.Send(ctx, *event); cloudevents.IsUndelivered(result) {
				log.Printf("Failed to send(TO): %s result %v, reason: no listeners", address, result)
				q.DataOut <- channel.DataEvent{
					Address:     address,
					Data:        e,
					EventStatus: channel.FAILED,
					EventType:   channel.SENDER,
				}
			} else if cloudevents.IsNACK(result) {
				log.Printf("Event not accepted: %v", result)
				q.DataOut <- channel.DataEvent{
					Address:     address,
					Data:        e,
					EventStatus: channel.SUCCEED,
					EventType:   channel.SENDER,
				}
			}
		}(sender, event, wg)
	}
}

// SendToAll ... sends events to all registered amqp address
func (q *Router) SendToAll(wg *sync.WaitGroup, event cloudevents.Event) {
	for k, s := range q.Senders {
		wg.Add(1) //for each ptp you send a message since its
		go func(s *Protocol, address string, e cloudevents.Event, wg *sync.WaitGroup) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
			defer cancel()
			if result := s.Client.Send(ctx, event); cloudevents.IsUndelivered(result) {
				log.Printf("Failed to send(TOALL): %v", result)
				q.DataOut <- channel.DataEvent{
					Address:     address,
					Data:        &e,
					EventStatus: channel.FAILED,
					EventType:   channel.SENDER,
				} // Not the clean way of doing , revisit
			} else if cloudevents.IsNACK(result) {
				log.Printf("Event not accepted: %v", result)
				q.DataOut <- channel.DataEvent{
					Address:     address,
					Data:        &e,
					EventStatus: channel.SUCCEED,
					EventType:   channel.SENDER,
				} // Not the clean way of doing , revisit
			}
		}(s, k, event, wg)
	}
}

// NewSenderReceiver created New Sender and Receiver object
func NewSenderReceiver(hostName string, port int, senderAddress string, receiverAddress string) (sender *Protocol, receiver *Protocol, err error) {
	sender, err = NewReceiver(hostName, port, senderAddress)
	if err == nil {
		receiver, err = NewSender(hostName, port, receiverAddress)
	}
	return
}

//NewReceiver creates new receiver object
func NewReceiver(hostName string, port int, receiverAddress string) (receiver *Protocol, err error) {
	receiver = &Protocol{}
	var opts []amqp1.Option
	opts = append(opts, amqp1.WithReceiverLinkOption(amqp.LinkCredit(50)))

	p, err := amqp1.NewReceiverProtocol(fmt.Sprintf("%s:%d", hostName, port), receiverAddress, []amqp.ConnOption{}, []amqp.SessionOption{}, opts...)
	if err != nil {
		log.Printf("Failed to create amqp protocol for a Receiver: %v", err)
		return
	}
	log.Printf("(New Receiver) Connection established %s\n", receiverAddress)

	parent, cancelParent := context.WithCancel(context.Background())
	receiver.CancelFn = cancelParent
	receiver.ParentContext = parent
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	receiver.Protocol = p
	receiver.Client = c
	return
}

//NewSender creates new QDR ptp
func NewSender(hostName string, port int, address string) (sender *Protocol, err error) {
	sender = &Protocol{}
	var opts []amqp1.Option
	p, err := amqp1.NewSenderProtocol(fmt.Sprintf("%s:%d", hostName, port), address, []amqp.ConnOption{}, []amqp.SessionOption{}, opts...)
	if err != nil {
		log.Printf("Failed to create amqp protocol: %v", err)
		return
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	sender.Protocol = p
	sender.Client = c
	return
}
