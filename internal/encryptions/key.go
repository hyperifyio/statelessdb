// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateKey generates a new AES key.
// The key should be at least 32 bytes (AES-256).
func GenerateKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("[GenerateKey(%d)]: %w", keySize, err)
	}
	return key, nil
}
