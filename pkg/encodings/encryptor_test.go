// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings_test

import (
	"github.com/google/uuid"
	encodings2 "statelessdb/pkg/encodings"
	"statelessdb/pkg/errors"
	"statelessdb/pkg/states"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	validKey, err := encodings2.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	serializer := encodings2.NewGobSerializer[*states.ComputeState]("ComputeState")

	// Test successful initialization
	encryptor := encodings2.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(validKey)
	if err != nil {
		t.Fatalf("NewEncryptor failed with valid key: %v", err)
	}

	// Test initialization with invalid key size
	invalidKey := []byte("shortkey")
	encryptor = encodings2.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(invalidKey)
	if err == nil {
		t.Errorf("NewEncryptor should fail with invalid key size")
	}

	if err != errors.ErrEncryptorInitializeFailedKeySizeLessThanMinimum {
		t.Errorf("Expected error %v, got %v", errors.ErrEncryptorInitializeFailedKeySizeLessThanMinimum, err)
	}
}

func TestEncrypt(t *testing.T) {
	key, err := encodings2.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	serializer := encodings2.NewGobSerializer[*states.ComputeState]("ComputeState")

	encryptor := encodings2.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(key)
	if err != nil {
		t.Fatalf("Failed to initialize Encryptor: %v", err)
	}

	testCases := []struct {
		name string
		data *states.ComputeState
	}{
		{
			name: "Normal Board",
			data: states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext, err := encryptor.Encrypt(tc.data)
			if err != nil {
				t.Errorf("Encrypt failed: %v", err)
			}

			if ciphertext == "" {
				t.Errorf("Encrypt returned empty ciphertext")
			}
		})
	}
}

func TestEncryptorEncryptNonceUniqueness(t *testing.T) {
	key, err := encodings2.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	serializer := encodings2.NewGobSerializer[*states.ComputeState]("ComputeState")

	encryptor := encodings2.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(key)
	if err != nil {
		t.Fatalf("Failed to initialize Encryptor: %v", err)
	}

	data := states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil)

	ciphertext1, err := encryptor.Encrypt(data)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	ciphertext2, err := encryptor.Encrypt(data)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	if ciphertext1 == ciphertext2 {
		t.Errorf("Ciphertexts should differ due to unique nonces")
	}
}
