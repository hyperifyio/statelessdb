// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions_test

import (
	"statelessdb/internal/encryptions"
	"testing"
)

// TestGenerateKey verifies that GenerateKey generates a key of the correct size.
func TestGenerateKey(t *testing.T) {
	key, err := encryptions.GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if len(key) != encryptions.MinimumKeySizeAES256 {
		t.Errorf("Expected key length %d, got %d", encryptions.MinimumKeySizeAES256, len(key))
	}
}

// TestGenerateKeyUniqueness ensures that multiple keys generated are unique.
func TestGenerateKeyUniqueness(t *testing.T) {
	key1, err := encryptions.GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	key2, err := encryptions.GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if string(key1) == string(key2) {
		t.Errorf("Generated keys are not unique")
	}
}
