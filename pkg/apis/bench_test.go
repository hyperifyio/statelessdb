// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis_test

import (
	"testing"

	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"

	"statelessdb/pkg/apis"
	"statelessdb/pkg/dtos"
	"statelessdb/pkg/encodings"
	"statelessdb/pkg/encodings/json"
	"statelessdb/pkg/requests"
	"statelessdb/pkg/states"
)

func BenchmarkHandleComputeStateRequest(b *testing.B) {
	var err error

	now := time.Now().UnixMilli()

	serverKey, err := encodings.GenerateKey(32)
	if err != nil {
		b.Fatalf("Could not create server key: %v", err)
	}

	serializer := encodings.NewJsonSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewJsonUnserializer[*states.ComputeState]("ComputeState")

	encryptor := encodings.NewEncryptor[*states.ComputeState](serializer)
	if err = encryptor.Initialize(serverKey); err != nil {
		b.Fatalf("Could not create encryptor: %v", err)
	}

	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	if err = decryptor.Initialize(serverKey); err != nil {
		b.Fatalf("Could not create decryptor: %v", err)
	}

	newState := func() *states.ComputeState {
		return &states.ComputeState{}
	}

	newRequestDTO := func() *requests.ComputeRequest {
		return &requests.ComputeRequest{}
	}

	requestHandler := func(r *requests.ComputeRequest, state *states.ComputeState) (*states.ComputeState, error) {
		if state == nil {
			return states.NewComputeState(uuid.New(), uuid.New(), r.Received, r.Received, nil, nil, nil), nil
		}
		return state, nil
	}

	responseHandler := func(state *states.ComputeState, private string) interface{} {
		return dtos.NewComputeResponseDTO(
			state.Id,
			state.Owner,
			state.Created,
			state.Updated,
			state.Public,
			private,
		)
	}

	server := apis.NewServer()

	// With previous private data
	b.Run("with_Previous_Private_Data", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {

			// Create a sample request body
			public := make(map[string]interface{})
			private := make(map[string]interface{})

			state := states.NewComputeState(uuid.New(), uuid.New(), now, now, public, private, nil)

			var privateString string
			if privateString, err = encryptor.Encrypt(state); err != nil {
				b.Fatalf("Could not create privateString: %v", err)
				return
			}

			request := requests.ComputeRequest{
				PrivateData: privateString,
			}

			// Marshal the struct into JSON
			requestBodyBytes, err := json.Marshal(request)
			if err != nil {
				b.Fatalf("Could not marshal privateString: %v", err)
				return
			}

			// Convert to string if necessary
			requestBody := string(requestBodyBytes)

			req, err := http.NewRequest("POST", "/api/v1", bytes.NewBufferString(requestBody))
			if err != nil {
				b.Fatalf("Could not create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			requestManager, err := requests.NewJsonRequestManager[*states.ComputeState, *requests.ComputeRequest, *dtos.ComputeResponseDTO](
				"ComputeState", serverKey, newState, newRequestDTO,
			)
			if err != nil {
				b.Fatalf("Could not create request manager: %v", err)
			}

			requestResponseManager := requestManager.HandleWith(requestHandler).WithResponse(responseHandler)

			handler := server.BuildHandler(requestResponseManager)

			b.StartTimer()
			handler(rr, req)
			b.StopTimer()

			response := rr.Result()
			if response.StatusCode != 200 {
				b.Fatalf("Request failed with status %d %s", response.StatusCode, response.Status)
			}

		}

	})

	// Without private data
	b.Run("without_Private_Data", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			// Create a sample request body
			request := requests.ComputeRequest{}

			// Marshal the struct into JSON
			requestBodyBytes, err := json.Marshal(request)
			if err != nil {
				b.Fatalf("Could not marshal privateString: %v", err)
				return
			}

			// Convert to string if necessary
			requestBody := string(requestBodyBytes)

			req, err := http.NewRequest("POST", "/api/v1", bytes.NewBufferString(requestBody))
			if err != nil {
				b.Fatalf("Could not create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			requestManager, err := requests.NewJsonRequestManager[*states.ComputeState, *requests.ComputeRequest, *dtos.ComputeResponseDTO](
				"ComputeState", serverKey, newState, newRequestDTO,
			)
			if err != nil {
				b.Fatalf("Could not create request manager: %v", err)
			}

			requestResponseManager := requestManager.HandleWith(requestHandler).WithResponse(responseHandler)
			handler := server.BuildHandler(requestResponseManager)

			b.StartTimer()
			handler(rr, req)
			b.StopTimer()

			response := rr.Result()
			if response.StatusCode != 200 {
				b.Fatalf("Request failed with status %d %s", response.StatusCode, response.Status)
			}

		}

	})

}
