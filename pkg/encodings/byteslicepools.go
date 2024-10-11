// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

//var byteSlicePoolManager *types.MemoryPoolManager[[]byte]
//
//func init () {
//	byteSlicePoolManager = types.NewMemoryPoolManager[[]byte]( func(size int) func () []byte {
//		return func () []byte {
//			return make([]byte, 0, size)
//		}
//	})
//}
//
//// There is no proof using []byte pool helps, so it is not used right now.
//// Here for benchmarking purposes only.
//// @deprecated
//func getByteSlice (size, capacity int) []byte {
//	p := capacity
//	s := byteSlicePoolManager.Pool(p).Get()
//	//if cap(s) < capacity {
//	//	return make([]byte, size, capacity)
//	//}
//	s = s[:size]
//	return s
//}
//
//func releaseByteSlice (s []byte) {
//	p := cap(s)
//	byteSlicePoolManager.Pool(p).Put(s[:0])
//}
