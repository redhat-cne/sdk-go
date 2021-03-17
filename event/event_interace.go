package event

import (
	"time"
)

// Reader is the interface for reading through an event from attributes.
type Reader interface {
	// GetType returns event.GetType().
	GetType() string
	// GetTime returns event.GetTime().
	GetTime() time.Time
	// GetID returns event.GetID().
	GetID() string
	// GetDataSchema returns event.GetDataSchema().
	GetDataSchema() string
	// GetDataContentType returns event.GetDataContentType().
	GetDataContentType() string
	// GetData returns event.GetData()
	GetData() *Data
	// Clone clones the event .
	Clone() Event
	// String returns a pretty-printed representation of the EventContext.
	String() string
}

// Writer is the interface for writing through an event onto attributes.
// If an error is thrown by a sub-component, Writer caches the error
// internally and exposes errors with a call to event.Validate().
type Writer interface {
	// SetType performs event.SetType.
	SetType(string) error
	// SetID performs event.SetID.
	SetID(string) error
	// SetTime performs event.SetTime.
	SetTime(time.Time) error
	// SetDataSchema performs event.SetDataSchema.
	SetDataSchema(string) error
	// SetDataContentType performs event.SetDataContentType.
	SetDataContentType(string) error
	// SetData
	SetData(Data) error
}
