// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

import (
	"encoding/gob"
)

// Ensure types are registered only once
//var (
//	registerMutex sync.Mutex
//	registerMap   = make(map[reflect.Type]*sync.Once)
//)

// RegisterGobTypeOnce registers a necessary types with gob
func RegisterGobTypeOnce[T interface{}](name string) {

	if name != "" {

		////registerOnce.Do(func() {
		var data T
		gob.RegisterName(name, data)
		////})

		//var data T
		//typ := reflect.TypeOf(data)
		//
		//registerMutex.Lock()
		//once, exists := registerMap[typ]
		//if !exists {
		//	once = &sync.Once{}
		//	registerMap[typ] = once
		//}
		//registerMutex.Unlock()
		//
		//once.Do(func() {
		//	gob.RegisterName(name, data)
		//})

	}
}
