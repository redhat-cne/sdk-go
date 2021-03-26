package api

import (
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/redhat-cne/sdk-go/pubsub"
	"log"

	"github.com/redhat-cne/sdk-go/event"
)

//PublishEventToLog .. publish event data to a log
func PublishEventToLog(event event.Event) {
	log.Printf("Publishing event to log %#v", event)

}

//NewCloudEvents create new cloud event from cloud native events and pubsub
func NewCloudEvents(e event.Event, ps *pubsub.PubSub) (*cloudevents.Event, error) {
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
func GetCloudNativeEvents(ce *cloudevents.Event) (err error) {
	e := event.Event{}
	if ce.Data() == nil {
		return fmt.Errorf("event data is empty")
	}
	data := event.Data{}
	if err = json.Unmarshal(ce.Data(), &data); err != nil {
		return
	}
	e.SetDataContentType(event.ApplicationJSON)
	e.SetTime(ce.Time())
	e.SetType(ce.Type())
	e.SetData(data)
	return
}
