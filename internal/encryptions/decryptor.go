// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"statelessdb/internal/errors"
)

// Decryptor helps with providing memory for encryption
type Decryptor[T interface{}] struct {
	unserializer Unserializer[T]
	block        cipher.Block
	gcm          cipher.AEAD
}

// NewDecryptor creates a new encryptor
// - keySize should be at least 32
func NewDecryptor[T interface{}](unserializer Unserializer[T]) *Decryptor[T] {
	return &Decryptor[T]{unserializer: unserializer}
}

// Initialize initializes internal memory
func (e *Decryptor[T]) Initialize(key []byte) error {
	var err error
	if len(key) < MinimumKeySizeAES256 {
		log.Errorf("[Decryptor.Initialize] key size %d less than minimum %d", len(key), MinimumKeySizeAES256)
		return errors.ErrDecryptorInitializeFailedKeySizeLessThanMinimum
	}
	e.block, err = aes.NewCipher(key)
	if err != nil {
		log.Errorf("[Decryptor.Initialize]: NewCipher: %v", err)
		return errors.ErrDecryptorInitializeFailedNewCipher
	}
	e.gcm, err = cipher.NewGCM(e.block)
	if err != nil {
		log.Errorf("[Decryptor.Initialize]: NewGCM: %v", err)
		return errors.ErrDecryptorInitializeFailedNewGCM
	}
	return nil
}

// Decrypt decrypts encrypted Base64 encoded string using AES256 to a plaintext
// string.
func (e *Decryptor[T]) Decrypt(encryptedData string, out T) error {
	var err error

	data, err := base64.StdEncoding.DecodeString(encryptedData)

	if err != nil {
		log.Errorf("[Decrypt]: base64: DecodeString: %v", err)
		return errors.ErrDecryptBase64StringFailed
	}
	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		log.Errorf("[Decrypt]: data length %d is less than nonce size %d", len(data), nonceSize)
		return errors.ErrDecryptDataLengthLessThanNonceSize
	}
	nonce := data[:nonceSize]
	ciphertextBytes := data[nonceSize:]

	serialized, err := e.gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		log.Errorf("[Decrypt]: gcm.Open: %v", err)
		return errors.ErrDecryptFailed
	}
	if err = e.unserializer.Unserialize(serialized, out); err != nil {
		log.Errorf("[Decrypt]: decoding serialized data failed: %v", err)
		return errors.ErrDecryptDecodingSerializationFailed
	}
	return nil
}
