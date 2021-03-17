package pubsub

// Reader is the interface for reading through an event from attributes.
type Reader interface {
	// GetResource returns event.GetResource()
	GetResource() string
	// GetEndpointUri returns event.GetEndpointUri()
	GetEndpointURI() string
	// URILocation returns event.GetUriLocation()
	GetURILocation() string
	GetID() string
	// String returns a pretty-printed representation of the PubSub.
	String() string
}

// Writer is the interface for writing through an event onto attributes.
// If an error is thrown by a sub-component, Writer caches the error
// internally and exposes errors with a call to Writer.Validate().
type Writer interface {
	// Resource performs event.SetResource()
	SetResource(string) error
	// EndpointURI [erforms] event.SetEndpointURI()
	SetEndpointURI(string) error
	// URILocation performs event.SetURILocation()
	SetURILocation(string) error
	// SetID performs event.SetID.
	SetID(string) error
}
