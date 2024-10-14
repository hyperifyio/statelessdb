// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

import (
	"bytes"

	"github.com/hyperifyio/statelessdb/pkg/encodings"
	"github.com/hyperifyio/statelessdb/pkg/encodings/json"
	"github.com/hyperifyio/statelessdb/pkg/errors"
)

type EncryptedRequestManager[T interface{}, R Request, D interface{}] struct {
	Encryptor  *encodings.Encryptor[T]
	Decryptor  *encodings.Decryptor[T]
	NewState   func() T
	NewRequest func() R
}

func NewEncryptedRequestManager[T interface{}, R Request, D interface{}](
	encryptor *encodings.Encryptor[T],
	decryptor *encodings.Decryptor[T],
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

	// TODO: Use of this pool resulted in random fails. Fallback to new decoder each time.
	//reader := GetJsonReaderState()
	//defer reader.Release()
	//log.Debugf("[EncryptedRequestManager.DecodeRequest]: Resetting as: %v", body)
	//reader.Buffer.Reset(body)

	decoder := json.NewDecoder(bytes.NewReader(body))
	if err = decoder.Decode(&req); err != nil {
		log.Errorf("[EncryptedRequestManager.DecodeRequest]: Bad body error: %v", err)
		log.Debugf("[EncryptedRequestManager.DecodeRequest]: Bad body is: %v", body)
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
