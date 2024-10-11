// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings_test

import (
	encodings2 "statelessdb/pkg/encodings"
	"testing"
)

// TestGenerateKey verifies that GenerateKey generates a key of the correct size.
func TestGenerateKey(t *testing.T) {
	key, err := encodings2.GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if len(key) != encodings2.MinimumKeySizeAES256 {
		t.Errorf("Expected key length %d, got %d", encodings2.MinimumKeySizeAES256, len(key))
	}
}

// TestGenerateKeyUniqueness ensures that multiple keys generated are unique.
func TestGenerateKeyUniqueness(t *testing.T) {
	key1, err := encodings2.GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	key2, err := encodings2.GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if string(key1) == string(key2) {
		t.Errorf("Generated keys are not unique")
	}
}
