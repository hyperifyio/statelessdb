// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings_test

import (
	"encoding/base64"

	"github.com/google/uuid"

	"testing"

	"github.com/hyperifyio/statelessdb/pkg/encodings"
	"github.com/hyperifyio/statelessdb/pkg/errors"
	"github.com/hyperifyio/statelessdb/pkg/states"
)

func TestNewDecryptor(t *testing.T) {
	validKey, err := encodings.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test successful initialization
	unserializer := encodings.NewGobUnserializer[*states.ComputeState]("ComputeState")

	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(validKey)
	if err != nil {
		t.Fatalf("NewDecryptor failed with valid key: %v", err)
	}

	// Test initialization with invalid key size
	invalidKey := []byte("shortkey")
	decryptor = encodings.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(invalidKey)
	if err == nil {
		t.Errorf("NewDecryptor should fail with invalid key size")
	}

	if err != errors.ErrDecryptorInitializeFailedKeySizeLessThanMinimum {
		t.Errorf("Expected error %v, got %v", errors.ErrDecryptorInitializeFailedKeySizeLessThanMinimum, err)
	}
}

func TestDecrypt(t *testing.T) {
	key, err := encodings.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	serializer := encodings.NewGobSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewGobUnserializer[*states.ComputeState]("ComputeState")
	encryptor := encodings.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(key)
	if err != nil {
		t.Fatalf("Failed to initialize Encryptor: %v", err)
	}

	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(key)
	if err != nil {
		t.Fatalf("Failed to initialize Decryptor: %v", err)
	}

	testCases := []struct {
		name        string
		data        *states.ComputeState
		expectError bool
	}{
		{
			name:        "Normal 8x8 board",
			data:        states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ciphertext, err := encryptor.Encrypt(tc.data)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			if tc.name == "TamperedCiphertext" {
				// Tamper with the ciphertext by altering a character
				tampered := ciphertext[:len(ciphertext)-1] + "A"
				out := &states.ComputeState{}
				err := decryptor.Decrypt(tampered, out)
				if err == nil {
					t.Errorf("Decrypt should fail for tampered ciphertext")
				}
				return
			}

			var decrypted states.ComputeState
			err = decryptor.Decrypt(ciphertext, &decrypted)
			if tc.expectError {
				if err == nil {
					t.Errorf("Decrypt should have failed but did not")
				}
			} else {
				if err != nil {
					t.Errorf("Decrypt failed: %v", err)
				}
				if !decrypted.Equals(tc.data) {
					t.Errorf("Decrypted data does not match original.\nOriginal: %v\nDecrypted: %v", tc.data, decrypted)
				}
			}
		})
	}
}

func TestDecryptorWithWrongKey(t *testing.T) {
	key1, err := encodings.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key1: %v", err)
	}

	key2, err := encodings.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key2: %v", err)
	}

	serializer := encodings.NewGobSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewGobUnserializer[*states.ComputeState]("ComputeState")

	encryptor := encodings.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(key1)
	if err != nil {
		t.Fatalf("Failed to initialize Encryptor with key1: %v", err)
	}

	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(key2)
	if err != nil {
		t.Fatalf("Failed to initialize Decryptor with key2: %v", err)
	}

	data := states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil)

	ciphertext, err := encryptor.Encrypt(data)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	out := &states.ComputeState{}
	err = decryptor.Decrypt(ciphertext, out)
	if err == nil {
		t.Errorf("Decrypt should fail with wrong key, but got decrypted data: %v", out)
	}
}

func TestDecryptorDecryptDataLengthLessThanNonceSize(t *testing.T) {
	key, err := encodings.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	unserializer := encodings.NewGobUnserializer[*states.ComputeState]("ComputeState")

	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(key)
	if err != nil {
		t.Fatalf("Failed to initialize Decryptor: %v", err)
	}

	// Create a ciphertext shorter than nonce size
	shortCiphertext := base64.StdEncoding.EncodeToString([]byte("short"))

	out := &states.ComputeState{}
	err = decryptor.Decrypt(shortCiphertext, out)
	if err == nil {
		t.Errorf("Decrypt should fail for data length less than nonce size")
	}

	if err != errors.ErrDecryptDataLengthLessThanNonceSize {
		t.Errorf("Expected error %v, got %v", errors.ErrDecryptDataLengthLessThanNonceSize, err)
	}
}

func TestDecryptorDecryptEmptyString(t *testing.T) {
	key, err := encodings.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	serializer := encodings.NewGobSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewGobUnserializer[*states.ComputeState]("ComputeState")

	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(key)
	if err != nil {
		t.Fatalf("Failed to initialize Decryptor: %v", err)
	}

	data := states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil)
	encrypter := encodings.NewEncryptor[*states.ComputeState](serializer)
	err = encrypter.Initialize(key)
	if err != nil {
		t.Fatalf("Encryption initialize failed: %v", err)
	}
	ciphertext, err := encrypter.Encrypt(data)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted := &states.ComputeState{}

	err = decryptor.Decrypt(ciphertext, decrypted)
	if err != nil {
		t.Errorf("Decrypt failed: %v", err)
	}

	if !decrypted.Equals(data) {
		t.Errorf("Decrypted data does not match original.\nOriginal: '%v'\nDecrypted: '%v'", data, decrypted)
	}
}
