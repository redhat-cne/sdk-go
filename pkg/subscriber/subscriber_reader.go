package subscriber

import (
	"github.com/redhat-cne/sdk-go/pkg/store"
)

var _ Reader = (*Subscriber)(nil)

// GetClientID ... Get subscriber ClientID
func (s *Subscriber) GetClientID() string {
	return s.ClientID
}

// GetURILocation returns uri location
func (s *Subscriber) GetURILocation() string {
	return s.HealthEndPoint.String()
}

// GetStatus of the client connection
func (s *Subscriber) GetStatus() Status {
	return s.Status
}

// GetSubStore get subscription store
func (s *Subscriber) GetSubStore() *store.PubSubStore {
	return s.SubStore
}
