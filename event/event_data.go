package event

import "regexp"

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
	Resource    string      `json:"resource"`
	Version     string      `json:"version"`
	Values      []DataValue `json:"values"`
	FieldErrors map[string]error
}

// DataValue ...
type DataValue struct {
	DataType  DataType    `json:"dataType"`
	ValueType ValueType   `json:"valueType"`
	Value     interface{} `json:"value"`
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

// SetResource ...
func (d *Data) SetResource(s string) error {
	matched, err := regexp.MatchString(`([^/]+(/{2,}[^/]+)?)`, s)
	if matched {
		d.Resource = s
		d.fieldOK("resource")
	} else {
		d.fieldError("resource", err)
		return err
	}
	return nil
}

// SetVersion  ...
func (d *Data) SetVersion(s string) error {
	d.Version = s
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
	return nil
}

// GetResource ...
func (d *Data) GetResource() string {
	return d.Resource
}

// GetVersion ...
func (d *Data) GetVersion() string {
	return d.Version
}

// GetValues ...
func (d *Data) GetValues() []DataValue {
	return d.Values
}
