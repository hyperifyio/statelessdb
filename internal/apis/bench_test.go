// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis_test

import (
	"bytes"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/http/httptest"
	"statelessdb/internal/dtos"
	"time"

	"testing"

	"statelessdb/internal/apis"
	"statelessdb/internal/encryptions"
	"statelessdb/internal/states"
)

func BenchmarkHandleComputeStateRequest(b *testing.B) {

	now := time.Now().UnixMilli()

	serverKey, err := encryptions.GenerateKey(32)
	if err != nil {
		b.Fatalf("Could not create server key: %v", err)
	}

	serializer := encryptions.NewJsonSerializer[*states.ComputeState]("ComputeState")
	unserializer := encryptions.NewJsonUnserializer[*states.ComputeState]("ComputeState")

	encryptor := encryptions.NewEncryptor[*states.ComputeState](serializer)
	if err = encryptor.Initialize(serverKey); err != nil {
		b.Fatalf("Could not create encryptor: %v", err)
	}

	decryptor := encryptions.NewDecryptor[*states.ComputeState](unserializer)
	if err = decryptor.Initialize(serverKey); err != nil {
		b.Fatalf("Could not create decryptor: %v", err)
	}

	server := &apis.Server{
		Encryptor: encryptor,
		Decryptor: decryptor,
	}

	b.Run("with_Private", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {

			// Create a sample request body
			public := make(map[string]interface{})
			private := make(map[string]interface{})

			state := states.New(
				uuid.New(),
				uuid.New(),
				now,
				now,
				public,
				private,
			)

			var privateString string
			if privateString, err = state.Encrypt(server.Encryptor); err != nil {
				b.Fatalf("Could not create privateString: %v", err)
				return
			}

			request := dtos.ComputeRequestDTO{
				Private: privateString,
			}

			// Marshal the struct into JSON
			requestBodyBytes, err := jsoniter.Marshal(request)
			if err != nil {
				b.Fatalf("Could not marshal privateString: %v", err)
				return
			}

			// Convert to string if necessary
			requestBody := string(requestBodyBytes)

			req, err2 := http.NewRequest("POST", "/api/v1", bytes.NewBufferString(requestBody))
			if err2 != nil {
				b.Fatalf("Could not create request: %v", err2)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			b.StartTimer()
			server.HandleComputeStateRequest(rr, req)
			b.StopTimer()
		}

	})

	b.Run("without_Private", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			// Create a sample request body
			request := dtos.ComputeRequestDTO{}

			// Marshal the struct into JSON
			requestBodyBytes, err := jsoniter.Marshal(request)
			if err != nil {
				b.Fatalf("Could not marshal privateString: %v", err)
				return
			}

			// Convert to string if necessary
			requestBody := string(requestBodyBytes)

			req, err2 := http.NewRequest("POST", "/api/v1", bytes.NewBufferString(requestBody))
			if err2 != nil {
				b.Fatalf("Could not create request: %v", err2)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			b.StartTimer()
			server.HandleComputeStateRequest(rr, req)
			b.StopTimer()

		}

	})

}
