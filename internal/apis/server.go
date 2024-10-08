// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis

import (
	"statelessdb/internal/encryptions"
	"statelessdb/internal/states"
)

type Server struct {
	Encryptor *encryptions.Encryptor[*states.ComputeState]
	Decryptor *encryptions.Decryptor[*states.ComputeState]
}

func NewServer(
	encryptor *encryptions.Encryptor[*states.ComputeState],
	decryptor *encryptions.Decryptor[*states.ComputeState],
) *Server {
	return &Server{
		encryptor,
		decryptor,
	}
}
