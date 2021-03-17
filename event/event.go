package event

import (
	"github.com/redhat-cne/sdk-go/types"
	"strings"
)

// Event represents the canonical representation of a Cloud Native Event.
type Event struct {
	// ID of the event; must be non-empty and unique within the scope of the producer.
	// +required
	ID string `json:""`
	// Type - The type of the occurrence which has happened.
	// +required
	Type string `json:"type"`
	// DataContentType - the data content type
	DataContentType *string `json:"dataContentType"`
	// Time - A Timestamp when the event happened.
	// +optional
	Time *types.Timestamp `json:"time,omitempty"`
	// DataSchema - A link to the schema that the `data` attribute adheres to.
	// +optional
	DataSchema *types.URI `json:"dataschema,omitempty"`

	Data Data

	FieldErrors map[string]error
}

func (e *Event) fieldError(field string, err error) {
	if e.FieldErrors == nil {
		e.FieldErrors = make(map[string]error)
	}
	e.FieldErrors[field] = err
}

func (e *Event) fieldOK(field string) {
	if e.FieldErrors != nil {
		delete(e.FieldErrors, field)
	}
}

// New returns a new Event, an optional version can be passed to change the
// default spec version from 1.0 to the provided version.
func New() Event {
	/*specVersion := defaultEventVersion
	if len(version) >= 1 {
		specVersion = version[0]
	}*/
	e := &Event{}

	return *e
}

// String returns a pretty-printed representation of the Event.
func (e Event) String() string {
	b := strings.Builder{}

	b.WriteString("TODO")

	return b.String()
}

//Clone ...
func (e Event) Clone() Event {
	out := Event{}
	out.SetData(e.Data) //nolint:errcheck
	out.FieldErrors = e.cloneFieldErrors()
	return out
}

func (e Event) cloneFieldErrors() map[string]error {
	if e.FieldErrors == nil {
		return nil
	}
	newFE := make(map[string]error, len(e.FieldErrors))
	for k, v := range e.FieldErrors {
		newFE[k] = v
	}
	return newFE
}
