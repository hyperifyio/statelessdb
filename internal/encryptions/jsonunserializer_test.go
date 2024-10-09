// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions_test

import (
	"sync"
	"testing"

	"statelessdb/internal/encryptions"
)

// TestJsonUnserializer_UnserializeBasic tests basic unserialization.
func TestJsonUnserializer_UnserializeBasic(t *testing.T) {

	data := []byte("Hello, Unserialize!")

	unserializer := encryptions.NewJsonUnserializer[*[]byte]("")
	serializer := encryptions.NewJsonSerializer[[]byte]("")

	state, err := serializer.Serialize(data)
	if err != nil {
		t.Fatalf("Failed to serialize data: %v", err)
	}
	defer state.Release()
	serialized := state.Bytes()

	var decoded []byte
	if err := unserializer.Unserialize(serialized, &decoded); err != nil {
		t.Fatalf("Failed to unserialize data: %v", err)
	}

	if string(decoded) != string(data) {
		t.Errorf("Expected %s, got %s", data, string(decoded))
	}
}

// TestJsonUnserializer_UnserializeStruct tests unserialization of a struct.
func TestJsonUnserializer_UnserializeStruct(t *testing.T) {
	unserializer := encryptions.NewJsonUnserializer[*SampleStruct]("SampleStruct")
	data := &SampleStruct{
		ID:      1,
		Name:    "Unserialize Test",
		Numbers: []int{4, 5, 6},
		Details: map[string]string{
			"info1": "data1",
			"info2": "data2",
		},
	}
	serializer := encryptions.NewJsonSerializer[*SampleStruct]("SampleStruct")

	state, err := serializer.Serialize(data)
	if err != nil {
		t.Fatalf("Failed to serialize data: %v", err)
	}
	defer state.Release()
	serialized := state.Bytes()

	decoded := &SampleStruct{}
	if err := unserializer.Unserialize(serialized, decoded); err != nil {
		t.Fatalf("Failed to unserialize data: %v", err)
	}

	if !data.Equals(decoded) {
		t.Errorf("Decoded data does not match original")
	}
}

//// TestJsonUnserializer_UnserializeUnsupportedType tests unserialization of unsupported type.
//func TestJsonUnserializer_UnserializeUnsupportedType(t *testing.T) {
//	unserializer := encryptions.NewJsonUnserializer[chan int]()
//
//	// Channels are not supported by gob, so serialization would fail before unserialization
//	// Thus, this test is redundant. Instead, you can ensure that Serialize fails.
//	// Alternatively, test unserialization with corrupted data.
//}

// TestJsonUnserializer_Concurrency tests unserialization under concurrent access.
func TestJsonUnserializer_Concurrency(t *testing.T) {
	unserializer := encryptions.NewJsonUnserializer[*SampleStruct]("SampleStruct")
	wg := sync.WaitGroup{}
	numGoroutines := 50
	numIterations := 100

	serializer := encryptions.NewJsonSerializer[*SampleStruct]("SampleStruct")

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				data := &SampleStruct{
					ID:      id,
					Name:    "Concurrent Unserialize Test",
					Numbers: []int{j, j + 1, j + 2},
					Details: map[string]string{
						"iteration": "value",
					},
				}
				state, err := serializer.Serialize(data)
				if err != nil {
					t.Errorf("Goroutine %d: Failed to serialize data: %v", id, err)
					return
				}
				serialized := state.Bytes()

				decoded := &SampleStruct{}
				if err := unserializer.Unserialize(serialized, decoded); err != nil {
					t.Errorf("Goroutine %d: Failed to unserialize data: %v", id, err)
					state.Release()
					return
				}

				if !data.Equals(decoded) {
					t.Errorf("Goroutine %d: Decoded data does not match original", id)
				}
				state.Release()
			}
		}(i)
	}

	wg.Wait()
}

// TestJsonUnserializer_UnserializeInvalidData tests unserialization of invalid data.
func TestJsonUnserializer_UnserializeInvalidData(t *testing.T) {
	unserializer := encryptions.NewJsonUnserializer[*SampleStruct]("SampleStruct")

	// Create invalid serialized data
	invalidData := []byte("invalid json data")

	decoded := &SampleStruct{}
	err := unserializer.Unserialize(invalidData, decoded)
	if err == nil {
		t.Error("Expected error when unserializing invalid data, got nil")
	}
}
