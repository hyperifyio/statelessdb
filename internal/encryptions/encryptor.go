// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"statelessdb/pkg/errors"
)

// Encryptor helps with providing memory for encryption
type Encryptor[T interface{}] struct {
	serializer Serializer[T]
	block      cipher.Block
	gcm        cipher.AEAD
	buf        bytes.Buffer
}

// NewEncryptor creates a new encryptor
// - keySize should be at least 32
// - nonceSize should be at least 12
func NewEncryptor[T interface{}](serializer Serializer[T]) *Encryptor[T] {
	return &Encryptor[T]{serializer: serializer}
}

// Initialize initializes internal memory
func (e *Encryptor[T]) Initialize(key []byte) error {
	var err error
	if len(key) < MinimumKeySizeAES256 {
		log.Errorf("[Encryptor.Initialize] key size %d less than minimum %d", len(key), MinimumKeySizeAES256)
		return errors.ErrEncryptorInitializeFailedKeySizeLessThanMinimum
	}
	e.block, err = aes.NewCipher(key)
	if err != nil {
		log.Errorf("[Encryptor.Initialize]: NewCipher: %v", err)
		return errors.ErrEncryptorInitializeFailedNewCipher
	}
	e.gcm, err = cipher.NewGCM(e.block)
	if err != nil {
		log.Errorf("[Encryptor.Initialize]: NewGCM: %v", err)
		return errors.ErrEncryptorInitializeFailedNewGCM
	}
	return nil
}

// Encrypt encrypts plaintext string using AES.
//   - key should be at least 32 bytes.
//   - nonce should be at least 12 bytes.
//
// Returns Base64 encoded encrypted string.
func (e *Encryptor[T]) Encrypt(data T) (string, error) {
	var err error

	state, err := e.serializer.Serialize(data)
	defer state.Release()

	serialized := state.Bytes()

	if err != nil {
		log.Errorf("[Encrypt]: GobSerializer failed: %v", err)
		return "", errors.ErrEncryptorFailedToSerializeData
	}
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Errorf("[Encrypt]: Nonce generation failed: %v", err)
		return "", errors.ErrEncryptorFailedToInitializeNonce
	}
	ciphertext := e.gcm.Seal(nonce, nonce, serialized, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
