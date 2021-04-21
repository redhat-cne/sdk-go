package pubsub_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/redhat-cne/sdk-go/pkg/pubsub"
	"github.com/redhat-cne/sdk-go/pkg/types"
	api "github.com/redhat-cne/sdk-go/v1/pubsub"
	"github.com/stretchr/testify/assert"
)

var (
	storePath = "."

	publisher = pubsub.PubSub{
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/post/event"}},
		Resource:    "test/test",
	}
	subscription = pubsub.PubSub{
		EndPointURI: &types.URI{URL: url.URL{Scheme: "http", Host: "localhost:8080", Path: "/get/event"}},
		Resource:    "test/test",
	}
	globalInstance = api.GetAPIInstance(storePath)
)

func TestAPI_GetAPIInstance(t *testing.T) {

	localInstance := api.GetAPIInstance(storePath)

	assert.Equal(t, &globalInstance, &localInstance)
}

func TestAPI_CreatePublisher(t *testing.T) {
	p, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, p.URILocation, publisher.URILocation)
	assert.Equal(t, p.Resource, publisher.Resource)
	assert.Equal(t, p.EndPointURI, publisher.EndPointURI)
}
func TestAPI_CreateSubscription(t *testing.T) {
	s, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ID)
	assert.Equal(t, s.URILocation, subscription.URILocation)
	assert.Equal(t, s.Resource, subscription.Resource)
	assert.Equal(t, s.EndPointURI, subscription.EndPointURI)
}

func TestAPI_DeleteAllPublishers(t *testing.T) {
	p, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	assert.NotEmpty(t, p.ID)
	e = globalInstance.DeleteAllPublishers()
	assert.Nil(t, e)
	b, e := globalInstance.GetPublishersFromFile()
	assert.Nil(t, e)
	assert.Len(t, b, 0)
	assert.Len(t, globalInstance.GetPublishers(), 0)
}

func TestAPI_DeleteAllSubscriptions(t *testing.T) {
	p, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	assert.NotEmpty(t, p.ID)
	e = globalInstance.DeleteAllSubscriptions()
	assert.Nil(t, e)
	b, e := globalInstance.GetSubscriptionsFromFile()
	assert.Nil(t, e)
	assert.Len(t, b, 0)
	assert.Len(t, globalInstance.GetSubscriptions(), 0)
}

func TestAPI_DeletePublisher(t *testing.T) {
	p, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	assert.NotEmpty(t, p.ID)
	e = globalInstance.DeletePublisher(p.ID)
	assert.Nil(t, e)
	delPub, e := globalInstance.GetPublisher(p.ID)
	assert.NotNil(t, e)
	assert.Equal(t, delPub, pubsub.PubSub{})
}
func TestAPI_DeleteSubscription(t *testing.T) {
	s, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	assert.NotEmpty(t, s.ID)
	e = globalInstance.DeleteSubscription(s.ID)
	assert.Nil(t, e)
	delSub, e := globalInstance.GetSubscription(s.ID)
	assert.NotNil(t, e)
	assert.Equal(t, delSub, pubsub.PubSub{})
}
func TestAPI_GetFromPubStore(t *testing.T) {
	p, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	storePub, e := globalInstance.GetFromPubStore(p.Resource)
	assert.Nil(t, e)
	assert.Equal(t, p, storePub)
}
func TestAPI_GetFromSubStore(t *testing.T) {
	s, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	storeSub, e := globalInstance.GetFromSubStore(s.Resource)
	assert.Nil(t, e)
	assert.Equal(t, s, storeSub)

}
func TestAPI_GetPublisher(t *testing.T) {
	p, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	storeP, e := globalInstance.GetPublisher(p.ID)
	assert.Nil(t, e)
	assert.Equal(t, p, storeP)
}

func TestAPI_GetPublishers(t *testing.T) {
	_, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	pubs := globalInstance.GetPublishers()
	assert.Greater(t, len(pubs), 0)
}

func TestAPI_GetSubscription(t *testing.T) {
	s, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	storeS, e := globalInstance.GetSubscription(s.ID)
	assert.Nil(t, e)
	assert.Equal(t, s, storeS)
}

func TestAPI_GetSubscriptions(t *testing.T) {
	_, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	subs := globalInstance.GetSubscriptions()
	assert.Greater(t, len(subs), 0)
}

func TestAPI_HasPublisher(t *testing.T) {
	p, e := globalInstance.CreatePublisher(publisher)
	assert.Nil(t, e)
	fp, found := globalInstance.HasPublisher(p.Resource)
	assert.True(t, found)
	assert.Equal(t, p, fp)

}

func TestAPI_HasSubscription(t *testing.T) {
	s, e := globalInstance.CreateSubscription(subscription)
	assert.Nil(t, e)
	fs, found := globalInstance.HasSubscription(s.Resource)
	assert.True(t, found)
	assert.Equal(t, s, fs)
}

func TestTeardown(t *testing.T) {
	_ = globalInstance.DeleteAllSubscriptions()
	_ = globalInstance.DeleteAllPublishers()
	_ = os.Remove("./pub.json")
	_ = os.Remove("./sub.json")
}
