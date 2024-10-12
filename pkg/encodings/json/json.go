// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package json

import (
	jsoniter "github.com/json-iterator/go"
	"io"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Decoder interface {
	Decode(obj interface{}) error
	More() bool
	Buffered() io.Reader
}

type Encoder interface {
	Encode(val interface{}) error
	SetEscapeHTML(escapeHTML bool)
	SetIndent(prefix, indent string)
}

func NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func NewEncoder(reader io.Writer) Encoder {
	return json.NewEncoder(reader)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
