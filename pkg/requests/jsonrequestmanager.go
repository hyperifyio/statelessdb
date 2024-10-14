// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

import (
	encodings2 "github.com/hyperifyio/statelessdb/pkg/encodings"
	"github.com/hyperifyio/statelessdb/pkg/errors"
)

func NewJsonRequestManager[T interface{}, R Request, D interface{}](
	name string,
	serverKey []byte,
	newState func() T,
	newRequest func() R,
) (*EncryptedRequestManager[T, R, D], error) {

	serializer := encodings2.NewJsonSerializer[T](name)
	encryptor := encodings2.NewEncryptor[T](serializer)
	if err := encryptor.Initialize(serverKey); err != nil {
		log.Errorf("Failed to initialize encryptor: %v", err)
		return nil, errors.ErrFailedToInitializeEncryptor
	}

	unserializer := encodings2.NewJsonUnserializer[T](name)
	decryptor := encodings2.NewDecryptor[T](unserializer)
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
