package event

import (
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/redhat-cne/sdk-go/pubsub"
)

//NewCloudEvent create new cloud event from cloud native events and pubsub
func (e *Event) NewCloudEvent(ps *pubsub.PubSub) (*cloudevents.Event, error) {
	ce := cloudevents.NewEvent(cloudevents.VersionV03)
	ce.SetTime(e.GetTime())
	ce.SetType(e.Type)
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	ce.SetSource(ps.Resource) // bus address
	ce.SetSpecVersion(cloudevents.VersionV03)
	ce.SetID(uuid.New().String())
	if err := ce.SetData(cloudevents.ApplicationJSON, e.GetData()); err != nil {
		return nil, err
	}
	return &ce, nil
}

// GetCloudNativeEvents  get event data from cloud events object if its valid else return error
func (e *Event) GetCloudNativeEvents(ce *cloudevents.Event) (err error) {
	if ce.Data() == nil {
		return fmt.Errorf("event data is empty")
	}
	data := Data{}
	if err = json.Unmarshal(ce.Data(), &data); err != nil {
		return
	}
	e.SetDataContentType(ApplicationJSON)
	e.SetTime(ce.Time())
	e.SetType(ce.Type())
	e.SetData(data)
	return
}
