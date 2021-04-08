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

// Data ... cloud native events data
// Data Json payload is as follows,
//{
//	"version": "v1.0",
//	"values": [{
//		"resource": "/cluster/node/ptp",
//		"data_type": "notification",
//		"value_type": "enumeration",
//		"value": "ACQUIRING-SYNC"
//		}, {
//		"resource": "/cluster/node/clock",
//		"data_type": "metric",
// 		"value_type": "decimal64.3",
//		"value": 100.3
//		}]
//}
type Data struct {
	Version string      `json:"version" example:"v1"`
	Values  []DataValue `json:"values"`
}

// DataValue ...
// DataValue Json payload is as follows,
//{
//	"resource": "/cluster/node/ptp",
//	"data_type": "notification",
//	"value_type": "enumeration",
//	"value": "ACQUIRING-SYNC"
//}
type DataValue struct {
	Resource  string      `json:"resource" example:"/cluster/node/clock"`
	DataType  DataType    `json:"dataType" example:"metric"`
	ValueType ValueType   `json:"valueType" example:"decimal64.3"`
	Value     interface{} `json:"value" example:"100.3"`
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
