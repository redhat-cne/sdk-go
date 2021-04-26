package pubsub

import (
	"strings"

	"github.com/redhat-cne/sdk-go/pkg/types"
)

// PubSub represents the canonical representation of a Cloud Native Event Publisher and Subscribers .
// PubSub Json request payload is as follows,
// {
//  "id": "789be75d-7ac3-472e-bbbc-6d62878aad4a",
//  "endpointUri": "http://localhost:9090/ack/event",
//  "uriLocation":  "http://localhost:8080/api/ocloudNotifications/v1/publishers/{publisherid}",
//  "resource":  "/east-edge-10/vdu3/o-ran-sync/sync-group/sync-status/sync-state"
// }
// PubSub request model
type PubSub struct {
	// ID of the pub/sub; is updated on successful creation of publisher/subscription.
	ID string `json:"id" omit:"empty"`
	// EndPointURI - A URI describing the event action link.
	// +required
	EndPointURI *types.URI `json:"endpointUri" example:"http://localhost:9090/ack/event"`

	// URILocation - A URI describing the producer/subscription get link.
	URILocation *types.URI `json:"uriLocation" omit:"empty"`
	// Resource - The type of the Resource.
	// +required
	Resource string `json:"resource" example:"/east-edge-10/vdu3/o-ran-sync/sync-group/sync-status/sync-state"`
}

// String returns a pretty-printed representation of the Event.
func (ps *PubSub) String() string {
	b := strings.Builder{}
	b.WriteString("  endpointURI: " + ps.GetEndpointURI() + "\n")
	b.WriteString("  URILocation: " + ps.GetURILocation() + "\n")
	b.WriteString("  ID: " + ps.GetID() + "\n")
	b.WriteString("  Resource: " + ps.GetResource() + "\n")
	return b.String()
}
