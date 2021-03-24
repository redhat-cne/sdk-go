package rest

import (
	"encoding/json"
	"fmt"
	"github.com/aneeshkp/cloudevents-amqp/pkg/protocol"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/redhat-cne/sdk-go/channel"
	"github.com/redhat-cne/sdk-go/pubsub"
	"io/ioutil"
	"log"
	"net/http"
)

//createSubscription The POST method creates a subscription resource for the (Event) API consumer.
// SubscriptionInfo  status 201
// Shall be returned when the subscription resource created successfully.
/*Request
   {
	"ResourceType": "ptp",
	"SourceAddress":"/cluster-x/worker-1/SYNC/ptp",
    "EndpointURI ": "http://localhost:9090/resourcestatus/ptp", /// daemon
	"ResourceQualifier": {
			"NodeName":"worker-1"
		}
	}
Response:
		{
		//"SubscriptionID": "789be75d-7ac3-472e-bbbc-6d62878aad4a",
        "PublisherId": "789be75d-7ac3-472e-bbbc-6d62878aad4a",
        "SourceAddress":"/cluster-x/worker-1/SYNC/ptp",
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
func (s *Server) createSubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sub := pubsub.PubSub{}

	if err := json.Unmarshal(bodyBytes, &sub); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "marshalling error")
		return
	}
	if exists, err := s.GetFromPubStore(sub.GetResource()); err == nil {
		log.Printf("There was already subscription,skipping creation %v", exists)
		s.sendOut(channel.LISTENER, &sub)
		s.respondWithJSON(w, http.StatusCreated, exists)
		return
	}

	if sub.GetEndpointURI() != "" {
		response, err := s.HTTPClient.Post(sub.GetEndpointURI(), cloudevents.ApplicationJSON, nil)
		if err != nil {
			log.Printf("There was error validating endpointurl %v", err)
			s.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusNoContent {
			log.Printf("There was error validating endpointurl returned status code %d", response.StatusCode)
			s.respondWithError(w, http.StatusBadRequest, "Return url validation check failed for create subscription.check endpointURI")
			return
		}
	}

	//check sub.EndpointURI by get
	_ = sub.SetID(uuid.New().String())                                                                              //noling:errcheck
	_ = sub.SetURILocation(fmt.Sprintf("http://localhost:%d/%s/%s/%s", s.port, s.apiPath, "subscriptions", sub.ID)) //noling:errcheck

	// persist the subscription -
	//TODO:might want to use PVC to live beyond pod crash
	err = s.writeToFile(sub, fmt.Sprintf("%s/%s", s.storePath, subFile))
	if err != nil {
		log.Printf("Error writing to store %v\n", err)
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Stored in a file")
	//store the publisher
	s.publisher.Set(sub.ID, &sub)
	// go ahead and create QDR to this address
	s.sendOut(channel.LISTENER, &sub)
	s.respondWithJSON(w, http.StatusCreated, sub)
}

func (s *Server) createPublisher(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pub := pubsub.PubSub{}

	if err := json.Unmarshal(bodyBytes, &pub); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "marshalling error")
		return
	}
	if exists, err := s.GetFromPubStore(pub.GetResource()); err == nil {
		log.Printf("There was already publication,skipping creation %v", exists)
		s.sendOut(channel.SENDER, &pub)
		s.respondWithJSON(w, http.StatusCreated, exists)
		return
	}

	if pub.GetEndpointURI() != "" {
		response, err := s.HTTPClient.Post(pub.GetEndpointURI(), cloudevents.ApplicationJSON, nil)
		if err != nil {
			log.Printf("There was error validating endpointurl %v", err)
			s.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusNoContent {
			log.Printf("There was error validating endpointurl returned status code %d", response.StatusCode)
			s.respondWithError(w, http.StatusBadRequest, "Return url validation check failed for create subscription.check endpointURI")
			return
		}
	}

	//check sub.EndpointURI by get
	pub.SetID(uuid.New().String())                                                                           //nolint:errcheck                                                                      //noling:errcheck
	pub.SetURILocation(fmt.Sprintf("http://localhost:%d/%s/%s/%s", s.port, s.apiPath, "publishers", pub.ID)) //nolint:errcheck

	// persist the subscription -
	//TODO:might want to use PVC to live beyond pod crash
	err = s.writeToFile(pub, fmt.Sprintf("%s/%s", s.storePath, pubFile))
	if err != nil {
		log.Printf("Error writing to store %v\n", err)
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Stored in a file")
	//store the publisher
	s.publisher.Set(pub.ID, &pub)
	// go ahead and create QDR to this address
	s.sendOut(channel.SENDER, &pub)
	s.respondWithJSON(w, http.StatusCreated, pub)
}

func (s *Server) sendOut(eType channel.EventType, sub *pubsub.PubSub) {
	// go ahead and create QDR to this address
	s.dataOut <- channel.DataEvent{
		Address:     sub.GetResource(),
		Data:        &event.Event{},
		EventType:   eType,
		EndPointURI: sub.GetEndpointURI(),
		EventStatus: channel.NEW,
	}
}
func (s *Server) getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)
	subscriptionID, ok := queries["subscriptionid"]
	if ok {
		log.Printf("getting subscription by id %s", subscriptionID)
		if sub, ok := s.subscription.Store[subscriptionID]; ok {
			s.respondWithJSON(w, http.StatusOK, sub)
			return
		}
	}
	s.respondWithError(w, http.StatusBadRequest, "Subscriptions not found")
}

func (s *Server) getPublisherByID(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)
	PublisherID, ok := queries["publisherid"]
	if ok {
		log.Printf("Getting subscription by id %s", PublisherID)
		if pub, ok := s.publisher.Store[PublisherID]; ok {
			s.respondWithJSON(w, http.StatusOK, pub)
			return
		}
	}
	s.respondWithError(w, http.StatusBadRequest, "Publisher not found")
}
func (s *Server) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	s.getPubSub(w, r, fmt.Sprintf("%s/%s", s.storePath, subFile))
}

func (s *Server) getPublishers(w http.ResponseWriter, r *http.Request) {
	s.getPubSub(w, r, fmt.Sprintf("%s/%s", s.storePath, pubFile))
}

func (s *Server) getPubSub(w http.ResponseWriter, r *http.Request, filepath string) {
	var pubSub pubsub.PubSub
	b, err := pubSub.ReadFromFile(filepath)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "error reading file")
		return
	}
	s.respondWithByte(w, http.StatusOK, b)
}

func (s *Server) deletePublisher(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)
	PublisherID, ok := queries["publisherid"]
	if ok {

		if pub, ok := s.publisher.Store[PublisherID]; ok {
			if err := s.deleteFromFile(*pub, fmt.Sprintf("%s/%s", s.storePath, pubFile)); err != nil {
				s.respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			s.publisher.Delete(PublisherID)
			s.respondWithMessage(w, http.StatusOK, "OK")
			return
		}
	}
	//TODO: close QDR connection for this --> use same method as create
	s.respondWithError(w, http.StatusBadRequest, "publisherid param is missing")
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)
	subscriptionID, ok := queries["subscriotionid"]

	if ok {

		if sub, ok := s.subscription.Store[subscriptionID]; ok {
			if err := s.deleteFromFile(*sub, fmt.Sprintf("%s/%s", s.storePath, subFile)); err != nil {
				s.respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			s.publisher.Delete(subscriptionID)
			s.respondWithMessage(w, http.StatusOK, "Deleted")
			return
		}
	}
	//TODO: close QDR connection for this subscription --> use same method as create
	s.respondWithError(w, http.StatusBadRequest, "subscriotionid param is missing")
}
func (s *Server) deleteAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	err := s.deleteAllFromFile(fmt.Sprintf("%s/%s", s.storePath, subFile))
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	//empty the store
	s.subscription.Store = make(map[string]*pubsub.PubSub)
	//TODO: close QDR connection for this --> use same method as create
	s.respondWithMessage(w, http.StatusOK, "deleted all subscriptions")
}

func (s *Server) deleteAllPublishers(w http.ResponseWriter, r *http.Request) {
	err := s.deleteAllFromFile(fmt.Sprintf("%s/%s", s.storePath, pubFile))
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	//empty the store
	s.publisher.Store = make(map[string]*pubsub.PubSub)
	//TODO: close QDR connection for this --> use same method as create
	s.respondWithMessage(w, http.StatusOK, "deleted all publishers")
}

func (s *Server) publishEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	pub := pubsub.PubSub{}
	if err := json.Unmarshal(bodyBytes, &pub); err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	pub.ID = uuid.New().String()
	var eventData []byte
	if eventData, err = json.Marshal(&pub); err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	event, err := protocol.GetCloudEvent(eventData)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		s.dataOut <- channel.DataEvent{
			EventType:   channel.EVENT,
			Data:        &event,
			Address:     pub.GetResource(),
			EndPointURI: pub.GetEndpointURI()}
		s.respondWithMessage(w, http.StatusAccepted, "Event published")
	}
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", cloudevents.ApplicationJSON)
	s.respondWithJSON(w, code, map[string]string{"error": message})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", cloudevents.ApplicationJSON)
	w.WriteHeader(code)
	w.Write(response) //nolint:errcheck
}
func (s *Server) respondWithMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", cloudevents.ApplicationJSON)
	s.respondWithJSON(w, code, map[string]string{"status": message})
}

func (s *Server) respondWithByte(w http.ResponseWriter, code int, message []byte) {
	w.Header().Set("Content-Type", cloudevents.ApplicationJSON)
	w.WriteHeader(code)
	w.Write(message) //nolint:errcheck
}
