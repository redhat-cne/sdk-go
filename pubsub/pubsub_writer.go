package pubsub

import (
	"fmt"
	"github.com/redhat-cne/sdk-go/types"
	"net/url"
	"regexp"
	"strings"
)

var _ Writer = (*PubSub)(nil)

// SetResource implements EventWriter.SetResource
func (ps *PubSub) SetResource(s string) error {
	matched, err := regexp.MatchString(`([^/]+(/{2,}[^/]+)?)`, s)
	if matched {
		ps.Resource = s
	} else {
		return err
	}
	return nil
}

// SetID implements EventWriter.SetID
func (ps *PubSub) SetID(id string) {
	ps.ID = id
}

// SetEndpointURI ...
func (ps *PubSub) SetEndpointURI(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		ps.EndPointURI = nil
		err := fmt.Errorf("uriLocation is given empty string,should be valid url")
		return err
	}
	pu, err := url.Parse(s)
	if err != nil {
		return err
	}
	ps.EndPointURI = &types.URI{URL: *pu}
	return nil
}

// SetURILocation ...
func (ps *PubSub) SetURILocation(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		ps.URILocation = nil
		err := fmt.Errorf("uriLocation is given empty string,should be valid url")
		return err
	}
	pu, err := url.Parse(s)
	if err != nil {
		return err
	}
	ps.URILocation = &types.URI{URL: *pu}
	return nil
}
