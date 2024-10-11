// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"statelessdb/internal/helpers"
	"statelessdb/pkg/encodings"
	"sync"
	"testing"
)

// TestGobSerializer_SerializeBasicTypes tests serialization of basic types.
func TestGobSerializer_SerializeBasicTypes_String(t *testing.T) {
	serializer := encodings.NewGobSerializer[string]("string")

	// Test string
	strData := "Hello, World!"
	state, err := serializer.Serialize(strData)
	if err != nil {
		t.Fatalf("Failed to serialize string: %v", err)
	}
	defer state.Release()
	serializedStr := state.Bytes()

	var decodedStr string
	decoder := gob.NewDecoder(bytes.NewReader(serializedStr))
	if err := decoder.Decode(&decodedStr); err != nil {
		t.Fatalf("Failed to deserialize string: %v", err)
	}

	if decodedStr != strData {
		t.Errorf("Expected %s, got %s", strData, decodedStr)
	}
}

// TestGobSerializer_SerializeBasicTypes tests serialization of basic types.
func TestGobSerializer_SerializeBasicTypes_Int(t *testing.T) {
	serializer := encodings.NewGobSerializer[int]("int")

	// Test integer
	intData := 42
	state, err := serializer.Serialize(intData)
	if err != nil {
		t.Fatalf("Failed to serialize int: %v", err)
	}
	defer state.Release()
	serializedInt := state.Bytes()

	var decodedInt int
	decoder := gob.NewDecoder(bytes.NewReader(serializedInt))
	if err := decoder.Decode(&decodedInt); err != nil {
		t.Fatalf("Failed to deserialize int: %v", err)
	}

	if decodedInt != intData {
		t.Errorf("Expected %d, got %d", intData, decodedInt)
	}

}

// TestGobSerializer_SerializeStruct tests serialization of a struct.
func TestGobSerializer_SerializeStruct(t *testing.T) {
	serializer := encodings.NewGobSerializer[*SampleStruct]("SampleStruct")

	original := &SampleStruct{
		ID:      1,
		Name:    "Test",
		Numbers: []int{1, 2, 3},
		Details: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	state, err := serializer.Serialize(original)
	if err != nil {
		t.Fatalf("Failed to serialize struct: %v", err)
	}
	defer state.Release()
	serialized := state.Bytes()

	var decoded SampleStruct
	decoder := gob.NewDecoder(bytes.NewReader(serialized))
	if err := decoder.Decode(&decoded); err != nil {
		t.Fatalf("Failed to deserialize struct: %v", err)
	}

	if decoded.ID != original.ID || decoded.Name != original.Name {
		t.Errorf("Decoded struct fields do not match original")
	}

	if !helpers.CompareSlices(decoded.Numbers, original.Numbers) {
		t.Errorf("Decoded Numbers do not match original")
	}

	if !helpers.CompareMaps(decoded.Details, original.Details) {
		t.Errorf("Decoded Details do not match original")
	}
}

//// TestGobSerializer_SerializeUnsupportedType tests serialization of an unsupported type.
//func TestGobSerializer_SerializeUnsupportedType(t *testing.T) {
//	serializer := encodings.NewGobSerializer[chan int]()
//
//	// Channels are not supported by gob
//	ch := make(chan int)
//
//	_, err := serializer.Serialize(ch)
//	if err == nil {
//		t.Error("Expected error when serializing unsupported type, got nil")
//	}
//}

// TestGobSerializer_Concurrency tests serialization under concurrent access.
func TestGobSerializer_Concurrency(t *testing.T) {
	serializer := encodings.NewGobSerializer[*SampleStruct]("SampleStruct")

	wg := sync.WaitGroup{}
	numGoroutines := 100
	numIterations := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numIterations; j++ {
				data := &SampleStruct{
					ID:      id,
					Name:    "Concurrent Test",
					Numbers: []int{j, j + 1, j + 2},
					Details: map[string]string{
						"iteration": fmt.Sprintf("%d", j),
					},
				}
				state, err := serializer.Serialize(data)
				if err != nil {
					t.Errorf("Failed to serialize data in goroutine %d: %v", id, err)
					state.Release()
					return
				}
				serialized := state.Bytes()

				decoded := &SampleStruct{}
				decoder := gob.NewDecoder(bytes.NewReader(serialized))
				if err := decoder.Decode(decoded); err != nil {
					t.Errorf("Failed to deserialize data in goroutine %d: %v", id, err)
					state.Release()
					return
				}

				if !decoded.Equals(data) {
					t.Errorf("Decoded data does not match original in goroutine %d: %v vs %v", id, decoded, data)
				}

				state.Release()
			}
		}(i)
	}

	wg.Wait()
}

// TestGobSerializer_Once tests only once
func TestGobSerializer_Once(t *testing.T) {
	serializer := encodings.NewGobSerializer[*SampleStruct]("SampleStruct")
	data := &SampleStruct{
		ID:      1,
		Name:    "Once Test",
		Numbers: []int{1, 2, 3, 4},
		Details: map[string]string{
			"iteration": "1",
		},
	}
	state, err := serializer.Serialize(data)
	if err != nil {
		t.Errorf("Failed to serialize data: %v", err)
		return
	}
	defer state.Release()
	serialized := state.Bytes()

	var decoded SampleStruct
	decoder := gob.NewDecoder(bytes.NewReader(serialized))
	if err := decoder.Decode(&decoded); err != nil {
		t.Errorf("Failed to deserialize data: %v", err)
		return
	}

	if decoded.ID != data.ID || decoded.Name != data.Name {
		t.Errorf("Decoded data does not match original")
	}
}
