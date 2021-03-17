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
	e.fieldOK("Time")
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
		e.fieldError("DataSchema", err)
		return err
	}
	e.DataSchema = &types.URI{URL: *pu}
	e.fieldOK("DataSchema")
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
	e.fieldOK("DataContentType")
	return nil
}

//SetData ...
func (e *Event) SetData(data Data) error {
	nData := Data{
		Version: data.Version,
	}

	var nValues []DataValue

	for _, v := range data.Values {
		nValue := DataValue{
			Resource:  v.Resource,
			DataType:  v.DataType,
			ValueType: v.ValueType,
			Value:     v.Value,
		}
		nValues = append(nValues, nValue)
	}
	nData.Values = nValues
	e.Data = &nData
	e.fieldOK("Data")
	return nil
}
