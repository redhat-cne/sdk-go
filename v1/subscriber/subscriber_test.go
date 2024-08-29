package subscriber_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/redhat-cne/sdk-go/v2/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/v2/pkg/store"
	"github.com/redhat-cne/sdk-go/v2/pkg/subscriber"
	"github.com/redhat-cne/sdk-go/v2/pkg/types"
	api "github.com/redhat-cne/sdk-go/v2/v1/subscriber"
	"github.com/stretchr/testify/assert"
)

var (
	storePath         = "./subscribers"
	clientID, _       = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	subscriptionOneID = "123e4567-e89b-12d3-a456-426614174001"
	subscriptionTwoID = "123e4567-e89b-12d3-a456-426614174002"

	subscriptionOne = &pubsub.PubSub{
		Version:     "1.0",
		ID:          subscriptionOneID,
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/event"}},
		Resource:    "test/test/1",
		URILocation: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/subscription"}},
	}

	subscriptionTwo = &pubsub.PubSub{
		Version:     "1.0",
		ID:          subscriptionTwoID,
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/event"}},
		Resource:    "test/test/2",
		URILocation: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/subscription"}},
	}

	subscriberWithOneEventCheck = subscriber.Subscriber{
		ClientID: clientID,
		SubStore: &store.PubSubStore{
			RWMutex: sync.RWMutex{},
			Store:   map[string]*pubsub.PubSub{subscriptionOneID: subscriptionOne},
		},
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/health"}},
		Status:      1,
	}

	subscriberWithManyEventCheck = subscriber.Subscriber{
		ClientID: clientID,
		SubStore: &store.PubSubStore{
			RWMutex: sync.RWMutex{},
			Store:   map[string]*pubsub.PubSub{subscriptionOneID: subscriptionOne, subscriptionTwoID: subscriptionTwo},
		},
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/health"}},
		Status:      1,
	}

	globalInstance = api.GetAPIInstance(storePath)
)

func TestAPI_GetAPIInstance(t *testing.T) {
	localInstance := api.GetAPIInstance(storePath)
	assert.Equal(t, &globalInstance, &localInstance)
}

func TestAPI_CreateSubscription(t *testing.T) {
	defer clean()
	s, e := globalInstance.CreateSubscription(clientID, subscriberWithOneEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	assert.Equal(t, 1, len(s.SubStore.Store))
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].URILocation, subscriptionOne.URILocation)
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].Resource, subscriptionOne.Resource)
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].EndPointURI, subscriptionOne.EndPointURI)
	b, e := globalInstance.GetSubscriptionsFromFile(clientID)
	assert.Nil(t, e)
	assert.NotNil(t, b)
	assert.NotEmpty(t, b)
	var subscriptionClient subscriber.Subscriber
	e = json.Unmarshal(b, &subscriptionClient)
	assert.NotNil(t, subscriptionClient)
	assert.Nil(t, e)
	assert.NotEmpty(t, s, subscriptionClient)
	assert.Equal(t, *s, subscriptionClient)
	assert.NotNil(t, subscriptionClient.SubStore)
	assert.Equal(t, len(subscriptionClient.SubStore.Store), len(s.SubStore.Store))
	assert.Equal(t, subscriptionOne, subscriptionClient.SubStore.Store[subscriptionOne.ID])
}

func TestAPI_CreateTwoSubscription(t *testing.T) {
	defer clean()
	s, e := globalInstance.CreateSubscription(clientID, subscriberWithOneEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	assert.Equal(t, 1, len(s.SubStore.Store))
	s, e = globalInstance.CreateSubscription(clientID, subscriberWithManyEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	assert.Equal(t, 2, len(s.SubStore.Store))
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].URILocation, subscriptionOne.URILocation)
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].Resource, subscriptionOne.Resource)
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].EndPointURI, subscriptionOne.EndPointURI)
	b, e := globalInstance.GetSubscriptionsFromFile(clientID)
	assert.Nil(t, e)
	assert.NotNil(t, b)
	assert.NotEmpty(t, b)
	var subscriptionClient subscriber.Subscriber
	e = json.Unmarshal(b, &subscriptionClient)
	assert.NotNil(t, subscriptionClient)
	assert.Nil(t, e)
	assert.NotEmpty(t, s, subscriptionClient)
	assert.Equal(t, *s, subscriptionClient)
	assert.NotNil(t, subscriptionClient.SubStore)
	assert.Equal(t, len(s.SubStore.Store), len(subscriptionClient.SubStore.Store))
	//assert.NotEmpty(t, subscriber[0].SubStore.)
	assert.Equal(t, subscriptionOne, subscriptionClient.SubStore.Store[subscriptionOne.ID])
}

func TestAPI_DeleteAllSubscriptions(t *testing.T) {
	s, e := globalInstance.CreateSubscription(clientID, subscriberWithOneEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	_, e = globalInstance.DeleteAllSubscriptionsForClient(clientID)
	assert.Nil(t, e)
	b, e := globalInstance.GetSubscriptionsFromFile(clientID)
	assert.Nil(t, e)
	assert.Len(t, b, 0)
	assert.Len(t, globalInstance.GetSubscriptionsFromClientID(clientID), 0)
}

func TestAPI_DeleteSubscription(t *testing.T) {
	defer clean()
	s, e := globalInstance.CreateSubscription(clientID, subscriberWithOneEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	e = globalInstance.DeleteSubscription(clientID, subscriptionOne.ID)
	assert.Nil(t, e)
	delSub, _ := globalInstance.GetSubscription(clientID, subscriptionOne.ID)
	assert.Equal(t, delSub, pubsub.PubSub{})
}

func TestAPI_HasSubscription(t *testing.T) {
	defer clean()
	s, e := globalInstance.CreateSubscription(clientID, subscriberWithOneEventCheck)
	assert.Nil(t, e)
	fs, found := globalInstance.HasSubscription(clientID, subscriptionOne.Resource)
	assert.True(t, found)
	assert.Equal(t, *s.SubStore.Store[fs.ID], fs)
}

func Test_Concurrency(t *testing.T) {
	defer clean()
	s, e := globalInstance.CreateSubscription(clientID, subscriberWithOneEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	assert.Equal(t, 1, len(s.SubStore.Store))
	s, e = globalInstance.CreateSubscription(clientID, subscriberWithManyEventCheck)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ClientID)
	assert.NotNil(t, s.SubStore.Store)
	assert.Equal(t, 2, len(s.SubStore.Store))
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].URILocation, subscriptionOne.URILocation)
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].Resource, subscriptionOne.Resource)
	assert.Equal(t, s.SubStore.Store[subscriptionOne.ID].EndPointURI, subscriptionOne.EndPointURI)
	b, e := globalInstance.GetSubscriptionsFromFile(clientID)
	assert.Nil(t, e)
	assert.NotNil(t, b)
	assert.NotEmpty(t, b)
	var subscriptionClient subscriber.Subscriber
	e = json.Unmarshal(b, &subscriptionClient)
	assert.NotNil(t, subscriptionClient)
	assert.Nil(t, e)
	assert.NotEmpty(t, s, subscriptionClient)
	assert.Equal(t, *s, subscriptionClient)
	assert.NotNil(t, subscriptionClient.SubStore)
	assert.Equal(t, len(s.SubStore.Store), len(subscriptionClient.SubStore.Store))
	//assert.NotEmpty(t, subscriber[0].SubStore.)
	assert.Equal(t, subscriptionOne, subscriptionClient.SubStore.Store[subscriptionOne.ID])

	var wg sync.WaitGroup
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			cID := uuid.New()
			s, e = globalInstance.CreateSubscription(cID, subscriberWithManyEventCheck)
			assert.Nil(t, e)
			assert.NotEmpty(t, s.ClientID)
			assert.NotNil(t, globalInstance.GetSubscriptionsFromClientID(cID))
			store, ok := globalInstance.SubscriberStore.Get(cID)
			assert.True(t, ok)
			assert.Equal(t, 2, len(store.SubStore.Store))

			assert.Equal(t, store.SubStore.Store[subscriptionOne.ID].URILocation, subscriptionOne.URILocation)
			assert.Equal(t, store.SubStore.Store[subscriptionOne.ID].Resource, subscriptionOne.Resource)
			assert.Equal(t, store.SubStore.Store[subscriptionOne.ID].EndPointURI, subscriptionOne.EndPointURI)
			_, T := globalInstance.HasClient(subscriptionClient.ClientID)
			assert.True(t, true, T)
			_ = globalInstance.UpdateStatus(subscriptionClient.ClientID, subscriber.InActive)
			if s, ok := globalInstance.SubscriberStore.Get(subscriptionClient.ClientID); ok {
				assert.NotEmpty(t, s.EndPointURI)
			}
		}(&wg)
	}
	wg.Wait()
}

func clean() {
	globalInstance.DeleteAllSubscriptionsForClient(clientID) //nolint
}

func TestTeardown(*testing.T) {
	_ = os.Remove(fmt.Sprintf("%s/%s.json", storePath, clientID))
}
