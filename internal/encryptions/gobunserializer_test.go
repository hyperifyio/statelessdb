// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions_test

import (
	"statelessdb/internal/encryptions"
	"sync"
	"testing"
)

// TestGobUnserializer_UnserializeBasic tests basic unserialization.
func TestGobUnserializer_UnserializeBasic(t *testing.T) {

	data := []byte("Hello, Unserialize!")

	unserializer := encryptions.NewGobUnserializer[*[]byte]("")
	serializer := encryptions.NewGobSerializer[[]byte]("")

	state, err := serializer.Serialize(data)
	defer state.Release()
	serialized := state.Bytes()

	if err != nil {
		t.Fatalf("Failed to serialize data: %v", err)
	}

	var decoded []byte
	if err := unserializer.Unserialize(serialized, &decoded); err != nil {
		t.Fatalf("Failed to unserialize data: %v", err)
	}

	if string(decoded) != string(data) {
		t.Errorf("Expected %s, got %s", data, string(decoded))
	}
}

// TestGobUnserializer_UnserializeStruct tests unserialization of a struct.
func TestGobUnserializer_UnserializeStruct(t *testing.T) {
	unserializer := encryptions.NewGobUnserializer[*SampleStruct]("SampleStruct")
	data := &SampleStruct{
		ID:      1,
		Name:    "Unserialize Test",
		Numbers: []int{4, 5, 6},
		Details: map[string]string{
			"info1": "data1",
			"info2": "data2",
		},
	}
	serializer := encryptions.NewGobSerializer[*SampleStruct]("SampleStruct")

	state, err := serializer.Serialize(data)
	defer state.Release()
	serialized := state.Bytes()
	if err != nil {
		t.Fatalf("Failed to serialize data: %v", err)
	}

	decoded := &SampleStruct{}
	if err := unserializer.Unserialize(serialized, decoded); err != nil {
		t.Fatalf("Failed to unserialize data: %v", err)
	}

	if !data.Equals(decoded) {
		t.Errorf("Decoded data does not match original")
	}
}

//// TestGobUnserializer_UnserializeUnsupportedType tests unserialization of unsupported type.
//func TestGobUnserializer_UnserializeUnsupportedType(t *testing.T) {
//	unserializer := encryptions.NewGobUnserializer[chan int]()
//
//	// Channels are not supported by gob, so serialization would fail before unserialization
//	// Thus, this test is redundant. Instead, you can ensure that Serialize fails.
//	// Alternatively, test unserialization with corrupted data.
//}

// TestGobUnserializer_Concurrency tests unserialization under concurrent access.
func TestGobUnserializer_Concurrency(t *testing.T) {
	unserializer := encryptions.NewGobUnserializer[*SampleStruct]("SampleStruct")
	wg := sync.WaitGroup{}
	numGoroutines := 50
	numIterations := 100

	serializer := encryptions.NewGobSerializer[*SampleStruct]("SampleStruct")

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
				serialized := state.Bytes()
				if err != nil {
					t.Errorf("Goroutine %d: Failed to serialize data: %v", id, err)
					return
				}

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

// TestGobUnserializer_UnserializeInvalidData tests unserialization of invalid data.
func TestGobUnserializer_UnserializeInvalidData(t *testing.T) {
	unserializer := encryptions.NewGobUnserializer[*SampleStruct]("SampleStruct")

	// Create invalid serialized data
	invalidData := []byte("invalid gob data")

	decoded := &SampleStruct{}
	err := unserializer.Unserialize(invalidData, decoded)
	if err == nil {
		t.Error("Expected error when unserializing invalid data, got nil")
	}
}
