package event

import (
	"bytes"
	"fmt"
	"io"

	jsoniter "github.com/json-iterator/go"
)

// WriteJSON writes the in event in the provided writer.
// Note: this function assumes the input event is valid.
func WriteJSON(in *Event, writer io.Writer) error {
	stream := jsoniter.ConfigFastest.BorrowStream(writer)
	defer jsoniter.ConfigFastest.ReturnStream(stream)
	stream.WriteObjectStart()

	if in.DataContentType != nil {
		switch in.GetDataContentType() {
		case ApplicationJSON:
			stream.WriteObjectField("id")
			stream.WriteString(in.ID)
			stream.WriteMore()

			stream.WriteObjectField("type")
			stream.WriteString(in.GetType())

			if in.GetDataContentType() != "" {
				stream.WriteMore()
				stream.WriteObjectField("dataContentType")
				stream.WriteString(in.GetDataContentType())
			}

			if in.Time != nil {
				stream.WriteMore()
				stream.WriteObjectField("time")
				stream.WriteString(in.Time.String())
			}

			if in.GetDataSchema() != "" {
				stream.WriteMore()
				stream.WriteObjectField("dataSchema")
				stream.WriteString(in.GetDataSchema())
			}
		default:
			return fmt.Errorf("missing event content type")
		}
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event attributes: %w", stream.Error)
	}

	// Let's write the body
	data := in.GetData()

	if data != nil {
		stream.WriteMore()
		stream.WriteObjectField("data")
		if err := writeJSONData(data, writer, stream); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("data is not set")
	}
	stream.WriteObjectEnd()
	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event Data: %w", stream.Error)
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event extensions: %w", stream.Error)
	}
	return stream.Flush()
}

// WriteDataJSON writes the in data in the provided writer.
// Note: this function assumes the input event is valid.
func WriteDataJSON(in *Data, writer io.Writer) error {
	stream := jsoniter.ConfigFastest.BorrowStream(writer)
	defer jsoniter.ConfigFastest.ReturnStream(stream)
	if err := writeJSONData(in, writer, stream); err != nil {
		return err
	}
	return stream.Flush()
}
func writeJSONData(in *Data, writer io.Writer, stream *jsoniter.Stream) error {
	stream.WriteObjectStart()

	// Let's write the body
	if in != nil {
		data := in
		stream.WriteObjectField("version")
		stream.WriteString(data.GetVersion())
		stream.WriteMore()
		stream.WriteObjectField("values")
		stream.WriteArrayStart()
		for _, v := range data.Values {
			stream.WriteObjectStart()
			stream.WriteObjectField("resource")
			stream.WriteString(v.GetResource())
			stream.WriteMore()
			stream.WriteObjectField("dataType")
			stream.WriteString(string(v.DataType))
			stream.WriteMore()
			stream.WriteObjectField("valueType")
			stream.WriteString(string(v.ValueType))
			stream.WriteMore()
			stream.WriteObjectField("value")
			switch v.ValueType {
			case ENUMERATION:
				// if type is a string
				stream.WriteString(fmt.Sprintf("%v", v.Value))

			case DECIMAL:
				stream.WriteFloat64(v.Value.(float64))

			default:
				// if type is other than above
				return fmt.Errorf("error while writing the value attributes: unknown type")
			}
			stream.WriteObjectEnd()
		}
		stream.WriteArrayEnd()
		stream.WriteObjectEnd()
	} else {
		return fmt.Errorf("data version is not set")
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event Data: %w", stream.Error)
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event extensions: %w", stream.Error)
	}
	return nil
}

// MarshalJSON implements a custom json marshal method used when this type is
// marshaled using json.Marshal.
func (e Event) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := WriteJSON(&e, &buf)
	return buf.Bytes(), err
}

// MarshalJSON implements a custom json marshal method used when this type is
// marshaled using json.Marshal.
func (d Data) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := WriteDataJSON(&d, &buf)
	return buf.Bytes(), err
}
