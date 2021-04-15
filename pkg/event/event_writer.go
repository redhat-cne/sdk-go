package event

import (
	"net/url"
	"strings"
	"time"

	"github.com/redhat-cne/sdk-go/pkg/types"
)

var _ Writer = (*Event)(nil)

// SetType implements Writer.SetType
func (e *Event) SetType(t string) {
	e.Type = t
}

// SetID implements Writer.SetID
func (e *Event) SetID(id string) {
	e.ID = id

}

// SetTime implements Writer.SetTime
func (e *Event) SetTime(t time.Time) {
	if t.IsZero() {
		e.Time = nil
	} else {
		e.Time = &types.Timestamp{Time: t}
	}

}

// SetDataSchema implements Writer.SetDataSchema
func (e *Event) SetDataSchema(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		e.DataSchema = nil
	}
	pu, err := url.Parse(s)
	if err != nil {
		return err
	}
	e.DataSchema = &types.URI{URL: *pu}
	return nil
}

// SetDataContentType implements Writer.SetDataContentType
func (e *Event) SetDataContentType(ct string) {
	ct = strings.TrimSpace(ct)
	if ct == "" {
		e.DataContentType = nil
	} else {
		e.DataContentType = &ct
	}
}

//SetData ...
func (e *Event) SetData(data Data) {
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
}
