package pubsub

import (
	"fmt"
	"github.com/redhat-cne/sdk-go/types"
	"net/url"
	"regexp"
	"strings"
)

var _ Writer = (*PubSub)(nil)

// SetSpecVersion implements EventWriter.SetSpecVersion
func (ps *PubSub) SetSpecVersion(v string) error {
	ps.fieldOK("specVersion")
	return nil
}

// SetResource implements EventWriter.SetResource
func (ps *PubSub) SetResource(s string) error {
	matched, err := regexp.MatchString(`([^/]+(/{2,}[^/]+)?)`, s)
	if matched {
		ps.Resource = s
	} else {
		ps.fieldError("Resource", err)
		return err
	}
	return nil
}

// SetID implements EventWriter.SetID
func (ps *PubSub) SetID(id string) error {
	ps.ID = id
	ps.fieldOK("id")
	return nil
}

// SetEndpointURI ...
func (ps *PubSub) SetEndpointURI(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		ps.EndPointURI = nil
		err := fmt.Errorf("uriLocation is given empty string,should be valid url")
		ps.fieldError("endPointURI", err)
		return err
	}
	pu, err := url.Parse(s)
	if err != nil {
		ps.fieldError("endPointURI", err)
		return err
	}
	ps.EndPointURI = &types.URI{URL: *pu}

	ps.fieldOK("type")
	return nil
}

// SetURILocation ...
func (ps *PubSub) SetURILocation(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		ps.URILocation = nil
		err := fmt.Errorf("uriLocation is given empty string,should be valid url")
		ps.fieldError("uriLocation", err)
		return err
	}
	pu, err := url.Parse(s)
	if err != nil {
		ps.fieldError("uriLocation", err)
		return err
	}
	ps.URILocation = &types.URI{URL: *pu}

	ps.fieldOK("uriLocation")
	return nil
}
