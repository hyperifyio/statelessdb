// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions

const (
	MinimumKeySizeAES256                   = 32 // AES-256
	DefaultDecryptDataBufferCapacity       = 1024
	DefaultDecryptSerializedBufferCapacity = 1024
	DefaultEncryptBufferCapacity           = 1024
	ByteBufferPoolCapacityFactor           = 512
)
