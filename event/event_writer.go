package event

import (
	"github.com/redhat-cne/sdk-go/types"
	"net/url"
	"strings"
	"time"
)

var _ Writer = (*Event)(nil)

// SetType implements Writer.SetType
func (e *Event) SetType(t string) error {
	e.Type = t
	e.fieldOK("type")
	return nil
}

// SetID implements Writer.SetID
func (e *Event) SetID(id string) error {
	e.ID = id
	e.fieldOK("ID")
	return nil
}

// SetTime implements Writer.SetTime
func (e *Event) SetTime(t time.Time) error {
	if t.IsZero() {
		e.Time = nil
	} else {
		e.Time = &types.Timestamp{Time: t}
	}
	e.fieldOK("time")
	return nil

}

// SetDataSchema implements Writer.SetDataSchema
func (e *Event) SetDataSchema(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		e.DataSchema = nil
	}
	pu, err := url.Parse(s)
	if err != nil {
		e.fieldError("dataSchema", err)
		return err
	}
	e.DataSchema = &types.URI{URL: *pu}
	e.fieldOK("dataSchema")
	return nil
}

// SetDataContentType implements Writer.SetDataContentType
func (e *Event) SetDataContentType(ct string) error {
	ct = strings.TrimSpace(ct)
	if ct == "" {
		e.DataContentType = nil
	} else {
		e.DataContentType = &ct
	}
	e.fieldOK("dataContentType")
	return nil
}

//SetData ...
func (e *Event) SetData(data Data) error {
	e.Data = data
	e.fieldOK("data")
	return nil
}
