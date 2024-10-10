// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

import (
	"statelessdb/internal/encryptions"
	"statelessdb/pkg/errors"
)

func NewJsonRequestManager[T interface{}, R Request, D interface{}](
	name string,
	serverKey []byte,
	newState func() T,
	newRequest func() R,
) (*EncryptedRequestManager[T, R, D], error) {

	serializer := encryptions.NewJsonSerializer[T](name)
	encryptor := encryptions.NewEncryptor[T](serializer)
	if err := encryptor.Initialize(serverKey); err != nil {
		log.Errorf("Failed to initialize encryptor: %v", err)
		return nil, errors.ErrFailedToInitializeEncryptor
	}

	unserializer := encryptions.NewJsonUnserializer[T](name)
	decryptor := encryptions.NewDecryptor[T](unserializer)
	if err := decryptor.Initialize(serverKey); err != nil {
		log.Errorf("Failed to initialize decryptor: %v", err)
		return nil, errors.ErrFailedToInitializeDecryptor
	}

	return NewEncryptedRequestManager[T, R, D](
		encryptor,
		decryptor,
		newState,
		newRequest,
	), nil
}
