package pubsub

import (
	"github.com/redhat-cne/sdk-go/types"
	"io/ioutil"
	"os"
	"strings"
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

	FieldErrors map[string]error
}

func (ps *PubSub) fieldError(field string, err error) {
	if ps.FieldErrors == nil {
		ps.FieldErrors = make(map[string]error)
	}
	ps.FieldErrors[field] = err
}

func (ps *PubSub) fieldOK(field string) {
	if ps.FieldErrors != nil {
		delete(ps.FieldErrors, field)
	}
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

// ReadFromFile is used to read subscription from the file system
func (ps *PubSub) ReadFromFile(filePath string) (b []byte, err error) {
	//open file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//read file and unmarshall json file to slice of users
	b, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}
