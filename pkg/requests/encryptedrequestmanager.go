// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

import (
	encodings2 "statelessdb/pkg/encodings"
	"statelessdb/pkg/errors"
)

type EncryptedRequestManager[T interface{}, R Request, D interface{}] struct {
	Encryptor  *encodings2.Encryptor[T]
	Decryptor  *encodings2.Decryptor[T]
	NewState   func() T
	NewRequest func() R
}

func NewEncryptedRequestManager[T interface{}, R Request, D interface{}](
	encryptor *encodings2.Encryptor[T],
	decryptor *encodings2.Decryptor[T],
	newState func() T,
	newRequest func() R,
) *EncryptedRequestManager[T, R, D] {
	return &EncryptedRequestManager[T, R, D]{
		Encryptor:  encryptor,
		Decryptor:  decryptor,
		NewState:   newState,
		NewRequest: newRequest,
	}
}

var _ RequestManager[any, Request, any] = &EncryptedRequestManager[any, Request, any]{}

// DecodeRequest will decode request data bytes to request
func (h *EncryptedRequestManager[T, R, D]) DecodeRequest(body []byte) (R, error) {
	var err error
	req := h.NewRequest()
	reader := GetJsonReaderState()
	defer reader.Release()
	reader.Buffer.Reset(body)
	if err = reader.Decoder.Decode(&req); err != nil {
		log.Errorf("[EncryptedRequestManager.DecodeRequest]: Bad body error: %v", err)
		return req, errors.ErrBadRequestBodyError
	}
	return req, nil
}

// DecryptState will decrypt optional private state from the request
func (h *EncryptedRequestManager[T, R, D]) DecryptState(privateData string) (T, error) {
	state := h.NewState()
	if privateData != "" {
		if err := h.Decryptor.Decrypt(privateData, state); err != nil {
			log.Errorf("[EncryptedRequestManager.DecryptState] failed to decrypt state: %v", err)
			return state, errors.ErrFailedToDecryptComputeState
		}
	}
	return state, nil
}

// EncryptState will return the state as encrypted string
func (h *EncryptedRequestManager[T, R, D]) EncryptState(state T) (string, error) {
	var err error
	var private string
	if private, err = h.Encryptor.Encrypt(state); err != nil {
		log.Errorf("[EncryptedRequestManager.EncryptState]: encrypting: error: %v", err)
		return "", errors.ErrComputeStateEncryptionFailed
	}
	return private, nil
}

// HandleWith configures a function to handle specific request path
func (h *EncryptedRequestManager[T, R, D]) HandleWith(handleRequest ApiRequestHandlerFunc[T, R]) *RequestResponseManager[T, R, D] {
	return &RequestResponseManager[T, R, D]{
		h,
		handleRequest,
		nil,
		nil,
	}
}
