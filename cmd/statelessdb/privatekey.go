// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/hyperifyio/statelessdb/pkg/encodings"
)

// parsePrivateKeyString parses AES-256 key, used in --private-key argument
func parsePrivateKeyString(privateKeyString string) ([]byte, error) {
	if privateKeyString == "" {
		key, err := encodings.GenerateKey(32) // AES-256
		if err != nil {
			return nil, fmt.Errorf("parsePrivateKeyString: failed to generate key: %v", err)
		}
		log.Warnf("Initialized with a random private key '%s'. You might want to make this persistent.", hex.EncodeToString(key))
		return key, nil
	}
	serverKey, err := hex.DecodeString(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("parsePrivateKeyString: failed to decode private key: %v", err)
	}
	return serverKey, nil
}
