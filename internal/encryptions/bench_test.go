// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions_test

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"statelessdb/internal/encryptions"
	"statelessdb/internal/states"
	"testing"
)

func BenchmarkEncryptorDecryptor(b *testing.B) {

	dtoName := "ComputeState"
	dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)

	//dtoName := "states.ComputeState"
	//dto := &states.ComputeState{
	//	ID:      1,
	//	Name:    "Unserialize Test",
	//	Numbers: []int{4, 5, 6},
	//	Details: map[string]string{
	//		"info1": "data1",
	//		"info2": "data2",
	//	},
	//}

	key, err := encryptions.GenerateKey(32)
	if err != nil {
		b.Fatalf("Failed to generate key: %v", err)
	}

	///// GOB tests /////

	b.Run("GOB", func(b *testing.B) {

		gobSerializer := encryptions.NewGobSerializer[*states.ComputeState](dtoName)
		gobUnserializer := encryptions.NewGobUnserializer[*states.ComputeState](dtoName)

		gobEncryptor := encryptions.NewEncryptor[*states.ComputeState](gobSerializer)
		err = gobEncryptor.Initialize(key)
		if err != nil {
			b.Fatalf("Failed to initialize dobEncryptor: %v", err)
		}

		gobDecryptor := encryptions.NewDecryptor[*states.ComputeState](gobUnserializer)
		err = gobDecryptor.Initialize(key)
		if err != nil {
			b.Fatalf("Failed to initialize gobDecryptor: %v", err)
		}

		// Pre-encrypt the plaintext to use in decryption benchmark
		gobCiphertext, err := gobEncryptor.Encrypt(dto)
		if err != nil {
			b.Fatalf("GOB encryption failed: %v", err)
		}

		b.Run("Encrypt_Encode", func(b *testing.B) {
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				_, err := gobEncryptor.Encrypt(dto)
				b.StopTimer()
				if err != nil {
					b.Fatalf("Encryption failed: %v", err)
				}
			}
		})

		b.Run("Decrypt_Decode", func(b *testing.B) {
			out := states.ComputeState{}
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				err := gobDecryptor.Decrypt(gobCiphertext, &out)
				b.StopTimer()
				if err != nil {
					b.Fatalf("Decryption failed: %v", err)
				}
			}
		})

	})

	///// JSON tests /////

	b.Run("JSON", func(b *testing.B) {

		jsonSerializer := encryptions.NewJsonSerializer[*states.ComputeState](dtoName)
		jsonUnserializer := encryptions.NewJsonUnserializer[*states.ComputeState](dtoName)

		jsonEncryptor := encryptions.NewEncryptor[*states.ComputeState](jsonSerializer)
		err = jsonEncryptor.Initialize(key)
		if err != nil {
			b.Fatalf("Failed to initialize jsonEncryptor: %v", err)
		}

		jsonDecryptor := encryptions.NewDecryptor[*states.ComputeState](jsonUnserializer)
		err = jsonDecryptor.Initialize(key)
		if err != nil {
			b.Fatalf("Failed to initialize jsonDecryptor: %v", err)
		}

		// Pre-encrypt the plaintext to use in decryption benchmark
		jsonCiphertext, err := jsonEncryptor.Encrypt(dto)
		if err != nil {
			b.Fatalf("JSON encryption failed: %v", err)
		}

		b.Run("Encrypt_Encode", func(b *testing.B) {
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				_, err := jsonEncryptor.Encrypt(dto)
				b.StopTimer()
				if err != nil {
					b.Fatalf("Encryption failed: %v", err)
				}
			}
		})

		b.Run("Decrypt_Decode", func(b *testing.B) {
			out := states.ComputeState{}
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				err := jsonDecryptor.Decrypt(jsonCiphertext, &out)
				b.StopTimer()
				if err != nil {
					b.Fatalf("Decryption failed: %v", err)
				}
			}
		})

	})

}

func Benchmark_GOB(b *testing.B) {

	b.Run("Encode", func(b *testing.B) {

		b.Run("with_NewEncoder", func(b *testing.B) {
			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer
				b.StartTimer()
				encoder := gob.NewEncoder(&buf)
				err := encoder.Encode(dto)
				if err != nil {
					b.Fatal(err)
				}
				b.StopTimer()
			}
		})

	})

	b.Run("Decode", func(b *testing.B) {

		b.Run("with_NewDecoder", func(b *testing.B) {
			b.StopTimer()

			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)

			// Serialize the DTO once outside the loop to use in decoding
			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)
			if err := encoder.Encode(dto); err != nil {
				b.Fatalf("Failed to encode dto for GOBDecode benchmark: %v", err)
			}
			serializedData := buf.Bytes()

			// Reset the buffer to reuse it in the loop if necessary
			buf.Reset()

			var decoded states.ComputeState

			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				b.StartTimer()
				decoder := gob.NewDecoder(bytes.NewReader(serializedData))
				if err := decoder.Decode(&decoded); err != nil {
					b.StopTimer()
					b.Fatalf("Failed to decode dto in GOBDecode benchmark: %v", err)
				}
				b.StopTimer()

				if !decoded.Equals(dto) {
					b.Fatalf("Results did not match")
				}

			}
		})

	})

	b.Run("Encode_and_Decode", func(b *testing.B) {

		b.Run("with_NewEncoder_and_NewDecoder", func(b *testing.B) {
			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				var buf bytes.Buffer
				var decoded states.ComputeState

				// Change each dto a bit

				b.StartTimer()

				// Serialize using GOB
				encoder := gob.NewEncoder(&buf)
				if err := encoder.Encode(dto); err != nil {
					b.StopTimer()
					b.Fatalf("GOB Encode failed: %v", err)
				}

				// Deserialize using GOB
				decoder := gob.NewDecoder(&buf)
				if err := decoder.Decode(&decoded); err != nil {
					b.StopTimer()
					b.Fatalf("GOB Decode failed: %v", err)
				}
				b.StopTimer()

				if !decoded.Equals(dto) {
					b.Fatalf("Results did not match")
				}

			}
		})

		b.Run("shared_NewEncoder_and_NewDecoder", func(b *testing.B) {
			b.StopTimer()

			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)
			decoder := gob.NewDecoder(&buf)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)

				// Serialize using GOB
				var decoded states.ComputeState

				b.StartTimer()
				buf.Reset()
				if err := encoder.Encode(dto); err != nil {
					b.StopTimer()
					b.Fatalf("GOB Encode failed: %v", err)
				}

				// Deserialize using GOB
				if err := decoder.Decode(&decoded); err != nil {
					b.StopTimer()
					b.Fatalf("GOB Decode failed: %v", err)
				}
				b.StopTimer()

				if !decoded.Equals(dto) {
					b.Fatalf("Results did not match")
				}

			}
		})

	})

}

func Benchmark_JSON(b *testing.B) {

	b.Run("Decode", func(b *testing.B) {

		b.Run("with_NewDecoder", func(b *testing.B) {
			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)

			// Serialize the DTO once outside the loop to use in decoding
			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			if err := encoder.Encode(dto); err != nil {
				b.Fatalf("Failed to encode dto for JSONDecode benchmark: %v", err)
			}
			serializedData := buf.Bytes()

			// Reset the buffer to reuse it in the loop if necessary
			buf.Reset()

			var decoded states.ComputeState
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				decoder := json.NewDecoder(bytes.NewReader(serializedData))
				if err := decoder.Decode(&decoded); err != nil {
					b.StopTimer()
					b.Fatalf("Failed to decode dto in JSONDecode benchmark: %v", err)
				}
				b.StopTimer()

				if !decoded.Equals(dto) {
					b.Fatalf("Results did not match")
				}

			}
		})

	})

	b.Run("Encode", func(b *testing.B) {

		b.Run("with_NewEncoder", func(b *testing.B) {
			b.StopTimer()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				var buf bytes.Buffer
				dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)

				b.StartTimer()
				encoder := json.NewEncoder(&buf)
				err := encoder.Encode(dto)
				if err != nil {
					b.StopTimer()
					b.Fatal(err)
				}
				b.StopTimer()

			}
		})

	})

	b.Run("Marshal_and_Unmarshal", func(b *testing.B) {
		b.StopTimer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			// Change each dto a bit
			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)
			dto.Created = int64(i)
			var decoded states.ComputeState

			b.StartTimer()

			// Serialize using JSON
			serializedData, err := json.Marshal(dto)
			if err != nil {
				b.StopTimer()
				b.Fatalf("JSON Marshal failed: %v", err)
			}

			// Deserialize using JSON
			if err := json.Unmarshal(serializedData, &decoded); err != nil {
				b.StopTimer()
				b.Fatalf("JSON Unmarshal failed: %v", err)
			}

			b.StopTimer()

			if !decoded.Equals(dto) {
				b.Fatalf("Results did not match")
			}

		}
	})

	b.Run("Encode_and_Decode", func(b *testing.B) {

		b.Run("shared_NewEncoder_and_NewDecoder", func(b *testing.B) {
			b.StopTimer()

			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			decoder := json.NewDecoder(&buf)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				var decoded states.ComputeState
				dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)

				b.StartTimer()

				// Serialize using GOB
				buf.Reset()
				if err := encoder.Encode(dto); err != nil {
					b.StopTimer()
					b.Fatalf("GOB Encode failed: %v", err)
				}

				// Deserialize using GOB
				if err := decoder.Decode(&decoded); err != nil {
					b.StopTimer()
					b.Fatalf("GOB Decode failed: %v", err)
				}

				b.StopTimer()

				if !decoded.Equals(dto) {
					b.Fatalf("Results did not match")
				}

			}
		})

	})

	b.Run("Marshal", func(b *testing.B) {
		b.StopTimer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)
			dto.Created = int64(i)

			b.StartTimer()
			_, err := json.Marshal(dto)
			if err != nil {
				b.StopTimer()
				b.Fatal(err)
			}
			b.StopTimer()
		}
	})

	b.Run("Unmarshal", func(b *testing.B) {
		b.StopTimer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dto := states.New(uuid.New(), uuid.New(), 0, 0, nil, nil)
			dto.Created = int64(i)

			// Serialize the DTO once outside the loop to use in unmarshalling
			serializedData, err := json.Marshal(dto)
			if err != nil {
				b.Fatalf("Failed to marshal dto for JSONUnmarshal benchmark: %v", err)
			}

			var decoded states.ComputeState

			b.StartTimer()
			if err := json.Unmarshal(serializedData, &decoded); err != nil {
				b.StopTimer()
				b.Fatalf("Failed to unmarshal dto in JSONUnmarshal benchmark: %v", err)
			}
			b.StopTimer()

			if !decoded.Equals(dto) {
				b.Fatalf("Results did not match")
			}

		}
	})

}
