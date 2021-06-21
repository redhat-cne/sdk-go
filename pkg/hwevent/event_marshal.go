// Copyright 2020 The Cloud Native Events Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hwevent

import (
	"bytes"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"

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
		stream.WriteObjectField("data")
		n, err := stream.Write(data.Data)
		if err != nil {
			return fmt.Errorf("error writing data: %w", err)
		}
		if n < len(data.Data) {
			return fmt.Errorf("failed to write data: %v of %v bytes written", n, len(data.Data))
		}
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
	log.Debugf("DZK UnmarshalJSON for event is called")
	var buf bytes.Buffer
	err := WriteJSON(&e, &buf)
	return buf.Bytes(), err
}

// MarshalJSON implements a custom json marshal method used when this type is
// marshaled using json.Marshal.
func (d Data) MarshalJSON() ([]byte, error) {
	log.Debugf("DZK UnmarshalJSON for data is called")
	var buf bytes.Buffer
	err := WriteDataJSON(&d, &buf)
	return buf.Bytes(), err
}
