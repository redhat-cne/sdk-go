package pubsub

import (
	"strings"

	"github.com/redhat-cne/sdk-go/pkg/types"
)

// PubSub represents the canonical representation of a Cloud Native Event Publisher and Subscribers .
type PubSub struct {
	// ID of the event; must be non-empty and unique within the scope of the producer.
	// +required
	ID string `json:"id"`
	// endPointURI - A URI describing the event producer.
	// +required
	EndPointURI *types.URI `json:"endpointUri"`
	// Source - A URI describing the event producer.
	// +required
	URILocation *types.URI `json:"uriLocation"`

	// Resource - The type of the Resource.
	// +required
	Resource string `json:"resource"`
}

// String returns a pretty-printed representation of the Event.
func (ps PubSub) String() string {
	b := strings.Builder{}
	b.WriteString("  endpointURI: " + ps.GetEndpointURI() + "\n")
	b.WriteString("  URILocation: " + ps.GetURILocation() + "\n")
	b.WriteString("  id: " + ps.GetID() + "\n")
	b.WriteString("  Resource: " + ps.GetResource() + "\n")
	return b.String()
}
