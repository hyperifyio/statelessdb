// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/hyperifyio/statelessdb/pkg/errors"
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

	//// Decode base62 (good but REALLY slow!)
	//data, err := base62.DecodeString(encryptedData)

	// Decode base64 (FAST!)
	data, err := base64.StdEncoding.DecodeString(encryptedData)

	//// Decode ascii85 (not nice with JSON!)
	//maxDecodedLen := len(encryptedData) * 4 / 5
	//decoded := make([]byte, maxDecodedLen)
	//n, _, err := ascii85.Decode(decoded, []byte(encryptedData), true)
	//if err != nil {
	//	log.Errorf("[Decryptor.Decrypt]: base64: DecodeString: %v", err)
	//	return errors.ErrDecryptBase64StringFailed
	//}
	//data := decoded[:n]

	// Decrypt
	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		log.Errorf("[Decryptor.Decrypt]: data length %d is less than nonce size %d", len(data), nonceSize)
		return errors.ErrDecryptDataLengthLessThanNonceSize
	}
	nonce := data[:nonceSize]
	ciphertextBytes := data[nonceSize:]

	serialized, err := e.gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		log.Errorf("[Decryptor.Decrypt]: gcm.Open: %v", err)
		return errors.ErrDecryptFailed
	}
	if err = e.unserializer.Unserialize(serialized, out); err != nil {
		log.Errorf("[Decryptor.Decrypt]: decoding serialized data failed: %v", err)
		return errors.ErrDecryptDecodingSerializationFailed
	}
	return nil
}
