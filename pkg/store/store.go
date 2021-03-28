package store

import (
	"github.com/redhat-cne/sdk-go/pkg/pubsub"
	"sync"
)

//PubSubStore ...
type PubSubStore struct {
	sync.RWMutex
	// PublisherStore stores publishers in a map
	Store map[string]*pubsub.PubSub
}

// Set is a wrapper for setting the value of a key in the underlying map
func (ps *PubSubStore) Set(key string, val *pubsub.PubSub) {
	ps.Lock()
	defer ps.Unlock()
	ps.Store[key] = val
}

//Delete ... delete from store
func (ps *PubSubStore) Delete(key string) {
	ps.Lock()
	defer ps.Unlock()
	delete(ps.Store, key)
}
