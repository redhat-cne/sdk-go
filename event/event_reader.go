package event

import (
	"time"
)

var _ Reader = (*Event)(nil)

// GetType implements Reader.Type
func (e Event) GetType() string {
	return e.Type
}

// GetID implements Reader.ID
func (e Event) GetID() string {
	return e.ID
}

// GetTime implements Reader.Time
func (e Event) GetTime() time.Time {
	if e.Time != nil {
		return e.Time.Time
	}
	return time.Time{}
}

// GetDataSchema implements Reader.DataSchema
func (e Event) GetDataSchema() string {
	return e.DataSchema.String()
}

// GetDataContentType implements Reader.DataContentType
func (e Event) GetDataContentType() string {
	return *e.DataContentType
}

// GetData ...
func (e Event) GetData() Data {
	return e.Data
}
