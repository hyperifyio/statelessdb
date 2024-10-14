// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests_test

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"testing"

	"github.com/hyperifyio/statelessdb/pkg/encodings"
	"github.com/hyperifyio/statelessdb/pkg/requests"
	"github.com/hyperifyio/statelessdb/pkg/states"
)

// MockSerializer and MockUnserializer can be implemented if needed.
// For this test, we are using the actual JsonSerializer and JsonUnserializer.

// TestEncryptedRequestManager_DecodeRequest tests the DecodeRequest method
func TestEncryptedRequestManager_DecodeRequest(t *testing.T) {

	// Define a 32-byte key for AES-256
	key := "a1ee74883d70fa9c4b5c9e5856ca58f99b26176be805d20d9c43fc4dbf880b91"

	serverKey, err := hex.DecodeString(key)
	if err != nil {
		t.Fatalf("failed to decode private key: %v", err)
	}

	// Initialize JsonSerializer and JsonUnserializer for ComputeState
	serializer := encodings.NewJsonSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewJsonUnserializer[*states.ComputeState]("ComputeState")

	// Initialize Encryptor
	encryptor := encodings.NewEncryptor[*states.ComputeState](serializer)
	if err := encryptor.Initialize(serverKey); err != nil {
		t.Fatalf("Failed to initialize Encryptor: %v", err)
	}

	// Initialize Decryptor
	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	if err := decryptor.Initialize(serverKey); err != nil {
		t.Fatalf("Failed to initialize Decryptor: %v", err)
	}

	// Define NewState function
	newState := func() *states.ComputeState {
		return states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil)
	}

	// Define NewRequest function
	newRequest := func() *requests.ComputeRequest {
		return requests.NewComputeRequest(0, nil, "")
	}

	// Create EncryptedRequestManager
	manager := requests.NewEncryptedRequestManager[*states.ComputeState, *requests.ComputeRequest, interface{}](
		encryptor,
		decryptor,
		newState,
		newRequest,
	)

	// Define the byte array as provided
	body := []byte{
		123, 34, 112, 114, 105, 118, 97, 116, 101, 34, 58, 34, 56, 88, 109, 103, 102, 106, 117, 57,
		110, 73, 72, 53, 116, 67, 84, 67, 83, 71, 114, 113, 71, 49, 114, 111, 118, 51, 102, 74,
		108, 85, 50, 67, 107, 53, 53, 121, 68, 84, 81, 108, 121, 86, 98, 119, 113, 72, 48, 99, 88,
		53, 115, 68, 111, 77, 107, 89, 74, 52, 86, 97, 90, 76, 84, 73, 76, 85, 82, 76, 110, 116,
		86, 87, 113, 86, 43, 106, 69, 47, 121, 50, 118, 69, 53, 87, 49, 122, 68, 52, 120, 103, 47,
		106, 90, 109, 52, 80, 69, 97, 105, 103, 86, 69, 51, 43, 51, 107, 69, 84, 57, 54, 85, 69,
		82, 47, 89, 52, 82, 120, 79, 87, 85, 82, 88, 112, 112, 83, 99, 48, 84, 69, 90, 48, 108,
		106, 72, 53, 116, 114, 84, 101, 106, 50, 87, 102, 104, 70, 122, 100, 90, 103, 76, 55, 70,
		85, 114, 71, 68, 72, 116, 114, 67, 81, 77, 82, 73, 90, 68, 122, 106, 75, 85, 86, 105, 110,
		99, 87, 80, 119, 101, 66, 84, 113, 48, 66, 54, 52, 70, 108, 78, 65, 86, 47, 74, 82, 116,
		101, 108, 103, 73, 72, 89, 102, 82, 106, 55, 104, 97, 77, 109, 118, 115, 81, 69, 82, 84,
		78, 109, 104, 76, 107, 69, 102, 109, 79, 102, 120, 74, 57, 81, 85, 49, 75, 70, 106, 43,
		106, 122, 103, 43, 56, 57, 83, 70, 119, 88, 110, 121, 108, 109, 103, 77, 84, 111, 98, 49,
		47, 69, 101, 84, 75, 117, 103, 61, 61, 34, 125,
	}

	// Expected PrivateData value after decoding
	expectedPrivate := "8Xmgfju9nIH5tCTCSGrqG1rov3fJlU2Ck55yDTQlyVbwqH0cX5sDoMkYJ4VaZLTILURLntVWqV+jE/y2vE5W1zD4xg/jZm4PEaigVE3+3kET96UER/Y4RxOWURXppSc0TEZ0ljH5trTej2WfhFzdZgL7FUrGDHtrCQMRIZDzjKUVincWPweBTq0B64FlNAV/JRtelgIHYfRj7haMmvsQERTNmhLkEfmOfxJ9QU1KFj+jzg+89SFwXnylmgMTob1/EeTKug=="

	// Call DecodeRequest
	decodedRequest, err := manager.DecodeRequest(body)
	if err != nil {
		t.Fatalf("DecodeRequest failed: %v", err)
	}

	// Verify that the PrivateData field matches the expected value
	if decodedRequest.PrivateData != expectedPrivate {
		t.Errorf("Expected PrivateData to be '%s', but got '%s'", expectedPrivate, decodedRequest.PrivateData)
	}

}

// TestEncryptedRequestManager_DecodeRequest_Concurrent tests the DecodeRequest method under concurrent usage
func TestEncryptedRequestManager_DecodeRequest_Concurrent(t *testing.T) {

	// Define a 32-byte key for AES-256 in hexadecimal
	keyHex := "a1ee74883d70fa9c4b5c9e5856ca58f99b26176be805d20d9c43fc4dbf880b91"

	// Decode the hexadecimal key to bytes
	serverKey, err := hex.DecodeString(keyHex)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	// Initialize JsonSerializer and JsonUnserializer for ComputeState
	serializer := encodings.NewJsonSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewJsonUnserializer[*states.ComputeState]("ComputeState")

	// Initialize Encryptor
	encryptor := encodings.NewEncryptor[*states.ComputeState](serializer)
	if err := encryptor.Initialize(serverKey); err != nil {
		t.Fatalf("Failed to initialize Encryptor: %v", err)
	}

	// Initialize Decryptor
	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	if err := decryptor.Initialize(serverKey); err != nil {
		t.Fatalf("Failed to initialize Decryptor: %v", err)
	}

	// Define NewState function
	newState := func() *states.ComputeState {
		return states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil)
	}

	// Define NewRequest function
	newRequest := func() *requests.ComputeRequest {
		return requests.NewComputeRequest(0, nil, "")
	}

	// Create EncryptedRequestManager
	manager := requests.NewEncryptedRequestManager[*states.ComputeState, *requests.ComputeRequest, interface{}](
		encryptor,
		decryptor,
		newState,
		newRequest,
	)

	// Define the byte array as provided
	body := []byte{
		123, 34, 112, 114, 105, 118, 97, 116, 101, 34, 58, 34, 56, 88, 109, 103, 102, 106, 117, 57,
		110, 73, 72, 53, 116, 67, 84, 67, 83, 71, 114, 113, 71, 49, 114, 111, 118, 51, 102, 74,
		108, 85, 50, 67, 107, 53, 53, 121, 68, 84, 81, 108, 121, 86, 98, 119, 113, 72, 48, 99, 88,
		53, 115, 68, 111, 77, 107, 89, 74, 52, 86, 97, 90, 76, 84, 73, 76, 85, 82, 76, 110, 116,
		86, 87, 113, 86, 43, 106, 69, 47, 121, 50, 118, 69, 53, 87, 49, 122, 68, 52, 120, 103, 47,
		106, 90, 109, 52, 80, 69, 97, 105, 103, 86, 69, 51, 43, 51, 107, 69, 84, 57, 54, 85, 69,
		82, 47, 89, 52, 82, 120, 79, 87, 85, 82, 88, 112, 112, 83, 99, 48, 84, 69, 90, 48, 108,
		106, 72, 53, 116, 114, 84, 101, 106, 50, 87, 102, 104, 70, 122, 100, 90, 103, 76, 55, 70,
		85, 114, 71, 68, 72, 116, 114, 67, 81, 77, 82, 73, 90, 68, 122, 106, 75, 85, 86, 105, 110,
		99, 87, 80, 119, 101, 66, 84, 113, 48, 66, 54, 52, 70, 108, 78, 65, 86, 47, 74, 82, 116,
		101, 108, 103, 73, 72, 89, 102, 82, 106, 55, 104, 97, 77, 109, 118, 115, 81, 69, 82, 84,
		78, 109, 104, 76, 107, 69, 102, 109, 79, 102, 120, 74, 57, 81, 85, 49, 75, 70, 106, 43,
		106, 122, 103, 43, 56, 57, 83, 70, 119, 88, 110, 121, 108, 109, 103, 77, 84, 111, 98, 49,
		47, 69, 101, 84, 75, 117, 103, 61, 61, 34, 125,
	}

	// Expected PrivateData value after decoding
	expectedPrivate := "8Xmgfju9nIH5tCTCSGrqG1rov3fJlU2Ck55yDTQlyVbwqH0cX5sDoMkYJ4VaZLTILURLntVWqV+jE/y2vE5W1zD4xg/jZm4PEaigVE3+3kET96UER/Y4RxOWURXppSc0TEZ0ljH5trTej2WfhFzdZgL7FUrGDHtrCQMRIZDzjKUVincWPweBTq0B64FlNAV/JRtelgIHYfRj7haMmvsQERTNmhLkEfmOfxJ9QU1KFj+jzg+89SFwXnylmgMTob1/EeTKug=="

	// Number of concurrent goroutines
	const goroutines = 100

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Channel to collect errors
	errCh := make(chan error, goroutines)

	// Launch concurrent goroutines
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Call DecodeRequest
			decodedRequest, err := manager.DecodeRequest(body)
			if err != nil {
				errCh <- err
				return
			}

			// Verify that the PrivateData field matches the expected value
			if decodedRequest.PrivateData != expectedPrivate {
				errCh <- &DecodeError{
					GoroutineID: id,
					Expected:    expectedPrivate,
					Actual:      decodedRequest.PrivateData,
				}
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errCh)

	// Check for errors
	for err := range errCh {
		if err != nil {
			t.Errorf("Concurrent DecodeRequest failed: %v", err)
		}
	}
}

// TestEncryptedRequestManager_DecodeRequest_Repeating tests the DecodeRequest method under repeating usage
func TestEncryptedRequestManager_DecodeRequest_Repeating(t *testing.T) {

	// Define a 32-byte key for AES-256 in hexadecimal
	keyHex := "a1ee74883d70fa9c4b5c9e5856ca58f99b26176be805d20d9c43fc4dbf880b91"

	// Decode the hexadecimal key to bytes
	serverKey, err := hex.DecodeString(keyHex)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	// Initialize JsonSerializer and JsonUnserializer for ComputeState
	serializer := encodings.NewJsonSerializer[*states.ComputeState]("ComputeState")
	unserializer := encodings.NewJsonUnserializer[*states.ComputeState]("ComputeState")

	// Initialize Encryptor
	encryptor := encodings.NewEncryptor[*states.ComputeState](serializer)
	if err := encryptor.Initialize(serverKey); err != nil {
		t.Fatalf("Failed to initialize Encryptor: %v", err)
	}

	// Initialize Decryptor
	decryptor := encodings.NewDecryptor[*states.ComputeState](unserializer)
	if err := decryptor.Initialize(serverKey); err != nil {
		t.Fatalf("Failed to initialize Decryptor: %v", err)
	}

	// Define NewState function
	newState := func() *states.ComputeState {
		return states.NewComputeState(uuid.New(), uuid.New(), 0, 0, nil, nil, nil)
	}

	// Define NewRequest function
	newRequest := func() *requests.ComputeRequest {
		return requests.NewComputeRequest(0, nil, "")
	}

	// Create EncryptedRequestManager
	manager := requests.NewEncryptedRequestManager[*states.ComputeState, *requests.ComputeRequest, interface{}](
		encryptor,
		decryptor,
		newState,
		newRequest,
	)

	// Expected PrivateData value after decoding
	expectedPrivate := "8Xmgfju9nIH5tCTCSGrqG1rov3fJlU2Ck55yDTQlyVbwqH0cX5sDoMkYJ4VaZLTILURLntVWqV+jE/y2vE5W1zD4xg/jZm4PEaigVE3+3kET96UER/Y4RxOWURXppSc0TEZ0ljH5trTej2WfhFzdZgL7FUrGDHtrCQMRIZDzjKUVincWPweBTq0B64FlNAV/JRtelgIHYfRj7haMmvsQERTNmhLkEfmOfxJ9QU1KFj+jzg+89SFwXnylmgMTob1/EeTKug=="

	// Number of concurrent goroutines
	const times = 100

	// Launch concurrent times
	for i := 0; i < times; i++ {

		// Define the byte array as provided
		body := []byte{
			123, 34, 112, 114, 105, 118, 97, 116, 101, 34, 58, 34, 56, 88, 109, 103, 102, 106, 117, 57,
			110, 73, 72, 53, 116, 67, 84, 67, 83, 71, 114, 113, 71, 49, 114, 111, 118, 51, 102, 74,
			108, 85, 50, 67, 107, 53, 53, 121, 68, 84, 81, 108, 121, 86, 98, 119, 113, 72, 48, 99, 88,
			53, 115, 68, 111, 77, 107, 89, 74, 52, 86, 97, 90, 76, 84, 73, 76, 85, 82, 76, 110, 116,
			86, 87, 113, 86, 43, 106, 69, 47, 121, 50, 118, 69, 53, 87, 49, 122, 68, 52, 120, 103, 47,
			106, 90, 109, 52, 80, 69, 97, 105, 103, 86, 69, 51, 43, 51, 107, 69, 84, 57, 54, 85, 69,
			82, 47, 89, 52, 82, 120, 79, 87, 85, 82, 88, 112, 112, 83, 99, 48, 84, 69, 90, 48, 108,
			106, 72, 53, 116, 114, 84, 101, 106, 50, 87, 102, 104, 70, 122, 100, 90, 103, 76, 55, 70,
			85, 114, 71, 68, 72, 116, 114, 67, 81, 77, 82, 73, 90, 68, 122, 106, 75, 85, 86, 105, 110,
			99, 87, 80, 119, 101, 66, 84, 113, 48, 66, 54, 52, 70, 108, 78, 65, 86, 47, 74, 82, 116,
			101, 108, 103, 73, 72, 89, 102, 82, 106, 55, 104, 97, 77, 109, 118, 115, 81, 69, 82, 84,
			78, 109, 104, 76, 107, 69, 102, 109, 79, 102, 120, 74, 57, 81, 85, 49, 75, 70, 106, 43,
			106, 122, 103, 43, 56, 57, 83, 70, 119, 88, 110, 121, 108, 109, 103, 77, 84, 111, 98, 49,
			47, 69, 101, 84, 75, 117, 103, 61, 61, 34, 125,
		}

		// Call DecodeRequest
		decodedRequest, err := manager.DecodeRequest(body)
		if err != nil {
			t.Errorf("Repeating DecodeRequest failed: Output not expected: %v", err)
		}

		// Verify that the PrivateData field matches the expected value
		if decodedRequest.PrivateData != expectedPrivate {
			t.Errorf("Repeating DecodeRequest failed: Output not expected")
		}
	}

}

// DecodeError represents a custom error for decoding mismatches
type DecodeError struct {
	GoroutineID int
	Expected    string
	Actual      string
}

func (e *DecodeError) Error() string {
	return "DecodeRequest mismatch in goroutine " +
		fmt.Sprintf("%d", e.GoroutineID) + ": expected '" + e.Expected + "', got '" + e.Actual + "'"
}
