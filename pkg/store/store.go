package store

import (
	"github.com/redhat-cne/sdk-go/pubsub"
	"sync"
)

//PubStore ...
type PubStore struct {
	sync.RWMutex
	// PublisherStore stores publishers in a map
	Store map[string]*pubsub.PubSub
}

//SubStore ...
type SubStore struct {
	sync.RWMutex
	// PublisherStore stores publishers in a map
	Store map[string]*pubsub.PubSub
}

// Set is a wrapper for setting the value of a key in the underlying map
func (ps *PubStore) Set(key string, val *pubsub.PubSub) {
	ps.Lock()
	defer ps.Unlock()
	ps.Store[key] = val
}

//Delete ... delete from store
func (ps *PubStore) Delete(key string) {
	ps.Lock()
	defer ps.Unlock()
	delete(ps.Store, key)
}

// Set is a wrapper for setting the value of a key in the underlying map
func (ss *SubStore) Set(key string, val *pubsub.PubSub) {
	ss.Lock()
	defer ss.Unlock()
	ss.Store[key] = val
}

//Delete from subscription
func (ss *SubStore) Delete(key string) {
	ss.Lock()
	defer ss.Unlock()
	delete(ss.Store, key)
}
