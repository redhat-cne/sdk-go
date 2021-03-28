package pubsub

var _ Reader = (*PubSub)(nil)

// GetResource implements EventReader.Resource
func (ps *PubSub) GetResource() string {
	return ps.Resource
}

// GetID implements EventReader.id
func (ps *PubSub) GetID() string {
	return ps.ID
}

//GetEndpointURI ...
func (ps *PubSub) GetEndpointURI() string {
	return ps.EndPointURI.String()
}

// GetURILocation ...
func (ps *PubSub) GetURILocation() string {
	return ps.URILocation.String()
}
