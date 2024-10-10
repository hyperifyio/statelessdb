// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package errors

import "errors"

// NOTE! For each of these, there should be only one place where it is returned.

var (
	ErrFailedToDecryptComputeState                     = errors.New("failed to decrypt compute state")
	ErrFailedToInitializeComputeState                  = errors.New("failed to initialize compute state")
	ErrDecryptorInitializeFailedKeySizeLessThanMinimum = errors.New("initializing decryptor: Key size is not enough")
	ErrDecryptorInitializeFailedNewCipher              = errors.New("initializing decryptor: failed to create cipher")
	ErrDecryptorInitializeFailedNewGCM                 = errors.New("initializing decryptor: failed to create GCM")
	ErrEncryptorInitializeFailedKeySizeLessThanMinimum = errors.New("initializing encryptor: Key size is not enough")
	ErrEncryptorInitializeFailedNewCipher              = errors.New("initializing encryptor: failed to create cipher")
	ErrEncryptorInitializeFailedNewGCM                 = errors.New("initializing encryptor: failed to create GCM")
	ErrEncryptorFailedToInitializeNonce                = errors.New("encryptor: failed to initialize nonce")
	ErrEncryptorFailedToSerializeData                  = errors.New("encryptor: failed to serialize data")
	ErrDecryptBase64StringFailed                       = errors.New("decrypting: Base64 decoding failed")
	ErrDecryptDataLengthLessThanNonceSize              = errors.New("decrypting: Data length less than nonce size")
	ErrDecryptDecodingSerializationFailed              = errors.New("decrypting: Failed to decode serialized data")
	ErrDecryptDecodingGobSerializationFailed           = errors.New("decrypting: Failed to decode GOB serialized data")
	ErrDecryptDecodingJsonSerializationFailed          = errors.New("decrypting: Failed to decode JSON serialized data")
	ErrDecryptFailed                                   = errors.New("decrypting failed")
	ErrFailedToInitializeEncryptor                     = errors.New("encryptor initialization failed")
	ErrFailedToInitializeDecryptor                     = errors.New("decryptor initialization failed")
	ErrBadRequestBodyError                             = errors.New("bad request body error")
	ErrRequestEncodingError                            = errors.New("request encoding error")
	ErrComputeStateEncryptionFailed                    = errors.New("compute state encryption failed")
)
