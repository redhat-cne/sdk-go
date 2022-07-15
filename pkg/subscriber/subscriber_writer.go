package subscriber

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/redhat-cne/sdk-go/pkg/store"
	"github.com/redhat-cne/sdk-go/pkg/types"
)

var _ Writer = (*Subscriber)(nil)

// SetClientID  ...
func (s *Subscriber) SetClientID(clientID string) {
	s.ClientID = clientID
}

// setSubStore ...
func (s *Subscriber) setSubStore(store store.PubSubStore) {
	s.SubStore = &store
}

// SetURILocation set uri location  (return url)
func (s *Subscriber) SetURILocation(sURLLocation string) error {
	sURLLocation = strings.TrimSpace(sURLLocation)
	if sURLLocation == "" {
		s.HealthEndPoint = nil
		err := fmt.Errorf("uriLocation is given empty string,should be valid url")
		return err
	}
	pu, err := url.Parse(sURLLocation)
	if err != nil {
		return err
	}
	s.HealthEndPoint = &types.URI{URL: *pu}
	return nil
}

// SetStatus set status of the connection
func (s *Subscriber) SetStatus(status Status) {
	s.Status = status
}
