// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis

import (
	"bytes"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var jsonReaderPoolState = sync.Pool{
	New: func() interface{} {
		reader := bytes.NewReader(nil)
		return &JsonReaderState{
			reader,
			json.NewDecoder(reader),
		}
	},
}

func GetJsonReaderState() *JsonReaderState {
	return jsonReaderPoolState.Get().(*JsonReaderState)
}

type JsonReaderState struct {
	Buffer  *bytes.Reader
	Decoder *jsoniter.Decoder
}

func (e *JsonReaderState) Release() {
	e.Buffer.Reset(nil)
	jsonReaderPoolState.Put(e)
}