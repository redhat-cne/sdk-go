package event

import (
	"fmt"
	"regexp"
)

// DataType ...
type DataType string

const (
	// NOTIFICATION ...
	NOTIFICATION DataType = "notification"
	// METRIC ...
	METRIC DataType = "metric"
)

// ValueType ...
type ValueType string

const (
	// ENUMERATION ...
	ENUMERATION ValueType = "enumeration"
	// DECIMAL ...
	DECIMAL ValueType = "decimal64.3"
)

// Data ...
type Data struct {
	Version     string      `json:"version"`
	Values      []DataValue `json:"values"`
	FieldErrors map[string]error
}

// DataValue ...
type DataValue struct {
	Resource    string      `json:"resource"`
	DataType    DataType    `json:"dataType"`
	ValueType   ValueType   `json:"valueType"`
	Value       interface{} `json:"value"`
	FieldErrors map[string]error
}

func (d *Data) fieldError(field string, err error) {
	if d.FieldErrors == nil {
		d.FieldErrors = make(map[string]error)
	}
	d.FieldErrors[field] = err
}

func (d *Data) fieldOK(field string) {
	if d.FieldErrors != nil {
		delete(d.FieldErrors, field)
	}
}

// SetVersion  ...
func (d *Data) SetVersion(s string) error {
	d.Version = s
	if s == "" {
		err := fmt.Errorf("version cannot be empty")
		d.fieldError("version", err)
		return err
	}
	d.fieldOK("version")
	return nil
}

// SetValues ...
func (d *Data) SetValues(v []DataValue) error {
	d.Values = v
	d.fieldOK("value")
	return nil
}

// AppendValues ...
func (d *Data) AppendValues(v DataValue) error {
	d.Values = append(d.Values, v)
	d.fieldOK("value")
	return nil
}

// GetVersion ...
func (d *Data) GetVersion() string {
	return d.Version
}

// GetValues ...
func (d *Data) GetValues() []DataValue {
	return d.Values
}

func (v *DataValue) fieldError(field string, err error) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]error)
	}
	v.FieldErrors[field] = err
}

func (v *DataValue) fieldOK(field string) {
	if v.FieldErrors != nil {
		delete(v.FieldErrors, field)
	}
}

// GetResource ...
func (v *DataValue) GetResource() string {
	return v.Resource
}

// SetResource ...
func (v *DataValue) SetResource(r string) error {
	matched, err := regexp.MatchString(`([^/]+(/{2,}[^/]+)?)`, r)
	if matched {
		v.Resource = r
		v.fieldOK("resource")
	} else {
		v.fieldError("resource", err)
		return err
	}
	return nil
}
