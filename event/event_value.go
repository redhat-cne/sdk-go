package event

import (
	"fmt"
)

// SetDataValue encodes the given payload
func (e *Event) SetDataValue(dataType DataType, obj interface{}) (err error) {
	var data DataValue

	switch dataType {
	case NOTIFICATION:
		{
			data.DataType = dataType
			data.ValueType = ENUMERATION
			data.Value = obj

		}
	case METRIC:
		{
			data.DataType = dataType
			data.ValueType = DECIMAL
			data.Value = obj

		}
	default:
		err = fmt.Errorf("error setting data %s - %v", dataType, obj)
	}

	e.Data.Values = append(e.Data.Values, data)
	e.Data.Values = append(e.Data.Values, data)
	return
}

// GetDataValue encodes the given payload
func (e *Event) GetDataValue() (data []DataValue, err error) {
	return e.Data.Values, nil
}
