package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redhat-cne/sdk-go/channel"
	"github.com/redhat-cne/sdk-go/pkg/store"
	"github.com/redhat-cne/sdk-go/pubsub"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	pubFile = "pub.json"
	subFile = "sub.json"
)

// Address of teh QDR
type Address struct {
	// Name of teh QDR address
	Name string `json:"name"`
}

// Server defines rest api server object
type Server struct {
	storePath  string
	port       int
	apiPath    string
	dataOut    chan<- channel.DataEvent
	HTTPClient *http.Client
	// PublisherStore stores publishers in a map
	publisher *store.PubStore
	// SubscriptionStore stores subscription in a map
	subscription *store.SubStore
	// StatusListenerQueue listens to any request coming for status check
	StatusListenerQueue *channel.ListenerChannel
}

// InitServer is used to supply configurations for rest api server
func InitServer(port int, apiPath, storePath string, dataOut chan<- channel.DataEvent) *Server {

	server := Server{
		storePath: storePath,
		port:      port,
		apiPath:   apiPath,
		dataOut:   dataOut,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
			Timeout: 10 * time.Second,
		},
		publisher: &store.PubStore{
			RWMutex: sync.RWMutex{},
			Store:   map[string]*pubsub.PubSub{},
		},
		subscription: &store.SubStore{
			RWMutex: sync.RWMutex{},
			Store:   map[string]*pubsub.PubSub{},
		},
		StatusListenerQueue: nil,
	}

	return &server
}

//StorePath path to store  pub/sub info
func (s *Server) StorePath() string {
	return s.storePath
}

// Port port id
func (s *Server) Port() int {
	return s.port
}

//GetFromPubStore get data from pub store
func (s *Server) GetFromPubStore(address string) (pubsub.PubSub, error) {
	for _, pub := range s.publisher.Store {
		if pub.GetResource() == address {
			return *pub, nil
		}
	}
	return pubsub.PubSub{}, fmt.Errorf("publisher not found for address %s", address)
}

//GetFromSubStore get data from sub store
func (s *Server) GetFromSubStore(address string) (pubsub.PubSub, error) {
	for _, sub := range s.subscription.Store {
		if sub.GetResource() == address {
			return *sub, nil
		}
	}
	return pubsub.PubSub{}, fmt.Errorf("subscription not found for address %s and %v", address, s.subscription.Store)
}

//WriteToFile writes subscription data to a file
func (s *Server) writeToFile(sub pubsub.PubSub, filePath string) error {
	//open file
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	//read file and unmarshall json file to slice of users
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var allSubs []pubsub.PubSub
	if len(b) > 0 {
		err = json.Unmarshal(b, &allSubs)
		if err != nil {
			return err
		}
	}
	allSubs = append(allSubs, sub)
	newBytes, err := json.MarshalIndent(&allSubs, "", " ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, newBytes, 0666); err != nil {
		return err
	}
	return nil

}

//DeleteAllFromFile deletes  publisher and subscription information from the file system
func (s *Server) deleteAllFromFile(filePath string) error {
	//open file
	if err := ioutil.WriteFile(filePath, []byte{}, 0666); err != nil {
		return err
	}
	return nil
}

//DeleteFromFile is used to delete subscription from the file system
func (s *Server) deleteFromFile(sub pubsub.PubSub, filePath string) error {
	//open file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	//read file and unmarshall json file to slice of users
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var allSubs []pubsub.PubSub
	if len(b) > 0 {
		err = json.Unmarshal(b, &allSubs)
		if err != nil {
			return err
		}
	}
	for k := range allSubs {
		// Remove the element at index i from a.
		if allSubs[k].ID == sub.ID {
			allSubs[k] = allSubs[len(allSubs)-1]      // Copy last element to index i.
			allSubs[len(allSubs)-1] = pubsub.PubSub{} // Erase last element (write zero value).
			allSubs = allSubs[:len(allSubs)-1]        // Truncate slice.
			break
		}
	}
	newBytes, err := json.MarshalIndent(&allSubs, "", " ")
	if err != nil {
		log.Printf("error deleting sub %v", err)
		return err
	}
	if err := ioutil.WriteFile(filePath, newBytes, 0666); err != nil {
		return err
	}
	return nil

}

// Start will start res api service
func (s *Server) Start() {
	pub := pubsub.PubSub{}
	b, err := pub.ReadFromFile(fmt.Sprintf("%s/%s", s.storePath, pubFile))
	if err != nil {
		panic(err)
	}
	var pubs []pubsub.PubSub
	if len(b) > 0 {
		if err := json.Unmarshal(b, &pubs); err != nil {
			panic(err)
		}
	}
	for _, pub := range pubs {
		s.publisher.Set(pub.ID, &pub)
	}

	//load subscription store
	sub := pubsub.PubSub{}
	b, err = sub.ReadFromFile(fmt.Sprintf("%s/%s", s.storePath, subFile))
	if err != nil {
		panic(err)
	}
	var subs []pubsub.PubSub
	if len(b) > 0 {
		if err := json.Unmarshal(b, &subs); err != nil {
			panic(err)
		}
	}
	for _, sub := range subs {
		s.subscription.Set(sub.ID, &sub)
	}

	r := mux.NewRouter()
	api := r.PathPrefix(s.apiPath).Subrouter()

	//The POST method creates a subscription resource for the (Event) API consumer.
	// SubscriptionInfo  status 201
	// Shall be returned when the subscription resource created successfully.
	/*Request
	   {
		"ResourceType": "ptp",
	    "EndpointURI ": "http://localhost:9090/resourcestatus/ptp", /// daemon
		"ResourceQualifier": {
				"NodeName":"worker-1"
				"Source":"/cluster-x/worker-1/SYNC/ptp"
			}
		}
	Response:
			{
			//"SubscriptionID": "789be75d-7ac3-472e-bbbc-6d62878aad4a",
	        "PublisherId": "789be75d-7ac3-472e-bbbc-6d62878aad4a",
			"URILocation": "http://localhost:8080/ocloudNotifications/v1/subsciptions/789be75d-7ac3-472e-bbbc-6d62878aad4a",
			"ResourceType": "ptp",
	         "EndpointURI ": "http://localhost:9090/resourcestatus/ptp", // address where the event
				"ResourceQualifier": {
				"NodeName":"worker-1"
	              "Source":"/cluster-x/worker-1/SYNC/ptp"
			}
		}*/

	/*201 Shall be returned when the subscription resource created successfully.
		See note below.
	400 Bad request by the API consumer. For example, the endpoint URI does not include ‘localhost’.
	404 Subscription resource is not available. For example, ptp is not supported by the node.
	409 The subscription resource already exists.
	*/
	api.HandleFunc("/subscriptions", s.createSubscription).Methods(http.MethodPost)
	api.HandleFunc("/publishers", s.createPublisher).Methods(http.MethodPost)
	/*
		 this method a list of subscription object(s) and their associated properties
		200  Returns the subscription resources and their associated properties that already exist.
			See note below.
		404 Subscription resources are not available (not created).
	*/
	api.HandleFunc("/subscriptions", s.getSubscriptions).Methods(http.MethodGet)
	api.HandleFunc("/publishers", s.getPublishers).Methods(http.MethodGet)
	// 200 and 404
	api.HandleFunc("/subscriptions/{subscriptionid}", s.getSubscriptionByID).Methods(http.MethodGet)
	api.HandleFunc("/publishers/{publisherid}", s.getPublisherByID).Methods(http.MethodGet)
	// 204 on success or 404
	api.HandleFunc("/subscriptions/{subscriptionid}", s.deleteSubscription).Methods(http.MethodDelete)
	api.HandleFunc("/publishers/{publisherid}", s.deletePublisher).Methods(http.MethodDelete)

	api.HandleFunc("/subscriptions", s.deleteAllSubscriptions).Methods(http.MethodDelete)
	api.HandleFunc("/publishers", s.deleteAllPublishers).Methods(http.MethodDelete)

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK") //nolint:errcheck
	}).Methods(http.MethodGet)

	api.HandleFunc("/dummy", s.dummy).Methods(http.MethodPost)

	api.HandleFunc("/create/event", s.publishEvent).Methods(http.MethodPost)

	err = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r)
	})

	log.Print("Started Rest API Server")
	log.Printf("endpoint %s", s.apiPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), api))

}
