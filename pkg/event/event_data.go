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
	Version string      `json:"version"`
	Values  []DataValue `json:"values"`
}

// DataValue ...
type DataValue struct {
	Resource  string      `json:"resource"`
	DataType  DataType    `json:"dataType"`
	ValueType ValueType   `json:"valueType"`
	Value     interface{} `json:"value"`
}

// SetVersion  ...
func (d *Data) SetVersion(s string) error {
	d.Version = s
	if s == "" {
		err := fmt.Errorf("version cannot be empty")
		return err
	}
	return nil
}

// SetValues ...
func (d *Data) SetValues(v []DataValue) {
	d.Values = v
}

// AppendValues ...
func (d *Data) AppendValues(v DataValue) {
	d.Values = append(d.Values, v)

}

// GetVersion ...
func (d *Data) GetVersion() string {
	return d.Version
}

// GetValues ...
func (d *Data) GetValues() []DataValue {
	return d.Values
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
	} else {
		return err
	}
	return nil
}
