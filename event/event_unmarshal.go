package event

import (
	"io"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"

	"github.com/redhat-cne/sdk-go/types"
)

var iterPool = sync.Pool{
	New: func() interface{} {
		return jsoniter.Parse(jsoniter.ConfigFastest, nil, 1024)
	},
}

func borrowIterator(reader io.Reader) *jsoniter.Iterator {
	iter := iterPool.Get().(*jsoniter.Iterator)
	iter.Reset(reader)
	return iter
}

func returnIterator(iter *jsoniter.Iterator) {
	iter.Error = nil
	iter.Attachment = nil
	iterPool.Put(iter)
}

//ReadJSON ...
func ReadJSON(out *Event, reader io.Reader) error {
	iterator := borrowIterator(reader)
	defer returnIterator(iterator)

	return readJSONFromIterator(out, iterator)
}

// readJSONFromIterator allows you to read the bytes reader as an event
func readJSONFromIterator(out *Event, iterator *jsoniter.Iterator) error {
	// Parsing dependency graph:
	//         SpecVersion
	//          ^     ^
	//          |     +--------------+
	//          +                    +
	//  All Attributes           datacontenttype (and datacontentencoding for v0.3)
	//  (except datacontenttype)     ^
	//                               |
	//                               |
	//                               +
	//                              Data

	var (
		// Universally parseable fields.
		id   string
		typ  string
		time *types.Timestamp
		data *Data

		// These fields require knowledge about the specversion to be parsed.
		//schemaurl jsoniter.Any
	)

	for key := iterator.ReadObject(); key != ""; key = iterator.ReadObject() {
		// Check if we have some error in our error cache
		if iterator.Error != nil {
			return iterator.Error
		}

		// If no specversion ...
		switch key {
		case "id":
			id = iterator.ReadString()
		case "type":
			typ = iterator.ReadString()
		case "time":
			time = readTimestamp(iterator)
		case "data":
			data, _ = readData(iterator)
		//case "dataSchema":
		//schemaurl = iterator.ReadAny()
		default:
			iterator.Skip()
		}

	}

	if iterator.Error != nil {
		return iterator.Error
	}
	out.Time = time
	out.ID = id
	out.Type = typ
	out.Data = *data
	return nil
}

func readTimestamp(iter *jsoniter.Iterator) *types.Timestamp {
	t, err := types.ParseTimestamp(iter.ReadString())
	if err != nil {
		iter.Error = err
	}
	return t
}

func readData(iter *jsoniter.Iterator) (*Data, error) {
	data := &Data{
		Resource: "",
		Version:  "",
		Values:   []DataValue{},
	}

	var values []DataValue
	for key := iter.ReadObject(); key != ""; key = iter.ReadObject() {
		// Check if we have some error in our error cache
		if iter.Error != nil {
			return data, iter.Error
		}

		switch key {
		case "resource":
			data.Resource = iter.ReadString()
		case "version":
			data.Version = iter.ReadString()
		case "values":

			for iter.ReadArray() {
				var cacheValue string
				dv := DataValue{}
				for dvField := iter.ReadObject(); dvField != ""; dvField = iter.ReadObject() {
					switch dvField {
					case "dataType":
						dv.DataType = DataType(iter.ReadString())
					case "valueType":
						dv.ValueType = ValueType(iter.ReadString())
					case "value":
						if dv.ValueType == DECIMAL {
							dv.Value = iter.ReadFloat64()
						} else {
							cacheValue = iter.ReadString()
						}
					default:
						iter.Skip()
					}
				}
				if dv.ValueType == DECIMAL {
					dv.Value, _ = strconv.ParseFloat(cacheValue, 3)
				} else {
					dv.Value = SyncState(cacheValue)
				}
				values = append(values, dv)
			}

			data.Values = values
		default:
			iter.Skip()
		}
	}

	return data, nil
}

// UnmarshalJSON implements the json unmarshal method used when this type is
// unmarshaled using json.Unmarshal.
func (e *Event) UnmarshalJSON(b []byte) error {
	iterator := jsoniter.ConfigFastest.BorrowIterator(b)
	defer jsoniter.ConfigFastest.ReturnIterator(iterator)
	return readJSONFromIterator(e, iterator)
}
