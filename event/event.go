package event

import (
	"fmt"
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
	// DataContentType - the Data content type
	DataContentType *string `json:"dataContentType"`
	// Time - A Timestamp when the event happened.
	// +optional
	Time *types.Timestamp `json:"time,omitempty"`
	// DataSchema - A link to the schema that the `Data` attribute adheres to.
	// +optional
	DataSchema *types.URI `json:"dataSchema,omitempty"`

	Data *Data `json:"data,omitempty"`
}

// New returns a new Event, an optional version can be passed to change the
// default spec version from 1.0 to the provided version.
func New() Event {
	/*specVersion := defaultEventVersion
	if len(version) >= 1 {
		specVersion = version[0]
	}*/
	e := Event{}

	return e
}

// NewCloudNativeEvent returns a new Event, an optional version can be passed to change the
// default spec version from 1.0 to the provided version.
func NewCloudNativeEvent() *Event {

	e := &Event{}

	return e
}

// String returns a pretty-printed representation of the Event.
func (e Event) String() string {
	b := strings.Builder{}
	b.WriteString("  id: " + e.ID + "\n")
	b.WriteString("  type: " + e.Type + "\n")
	if e.Time != nil {
		b.WriteString("  time: " + e.Time.String() + "\n")
	}

	b.WriteString("  data: \n")
	b.WriteString("  version: " + e.Data.Version + "\n")
	b.WriteString("  values: \n")
	for _, v := range e.Data.Values {
		b.WriteString("  value type : " + string(v.ValueType) + "\n")
		b.WriteString("  data type : " + string(v.DataType) + "\n")
		b.WriteString("  value : " + fmt.Sprintf("%v", v.Value) + "\n")
		b.WriteString("  resource: " + v.GetResource() + "\n")
	}

	return b.String()
}

//Clone ...
func (e Event) Clone() Event {
	out := Event{}
	out.SetData(*e.Data) //nolint:errcheck
	return out
}
