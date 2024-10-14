// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

import (
	"bytes"
	"github.com/hyperifyio/statelessdb/pkg/types"
)

var bytesBufferPoolManager *types.MemoryPoolManager[*bytes.Buffer]

func init() {
	bytesBufferPoolManager = types.NewMemoryPoolManager[*bytes.Buffer](func(size int) func() *bytes.Buffer {
		return func() *bytes.Buffer {
			return new(bytes.Buffer)
		}
	})
}

func getBytesBuffer() *bytes.Buffer {
	return bytesBufferPoolManager.Pool(0).Get()
}

func releaseBytesBuffer(s *bytes.Buffer) {
	s.Reset()
	bytesBufferPoolManager.Pool(0).Put(s)
}
