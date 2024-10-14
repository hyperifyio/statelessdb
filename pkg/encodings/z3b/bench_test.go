// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package z3b_test

import (
	"crypto/rand"
	"io"
	"testing"

	"bytes"

	"encoding/base64"

	"github.com/hyperifyio/statelessdb/pkg/encodings/z3b"
)

func BenchmarkEncodeZ3b_Small(b *testing.B) {
	data := generateRandomBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := z3b.Encode(data)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkDecodeZ3b_Small(b *testing.B) {
	data := generateRandomBytes(32)
	encoded, err := z3b.Encode(data)
	if err != nil {
		b.Fatalf("Encode failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := z3b.Decode(encoded)
		if err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
		if !bytes.Equal(decoded, data) {
			b.Fatalf("Decoded data does not match original data")
		}
	}
}

func BenchmarkEncodeZ3b_Medium(b *testing.B) {
	data := generateRandomBytes(1024) // 1 KB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := z3b.Encode(data)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkDecodeZ3b_Medium(b *testing.B) {
	data := generateRandomBytes(1024) // 1 KB
	encoded, err := z3b.Encode(data)
	if err != nil {
		b.Fatalf("Encode failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := z3b.Decode(encoded)
		if err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
		if !bytes.Equal(decoded, data) {
			b.Fatalf("Decoded data does not match original data")
		}
	}
}

func BenchmarkEncodeZ3b_Large(b *testing.B) {
	data := generateRandomBytes(1024 * 1024) // 1 MB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := z3b.Encode(data)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkDecodeZ3b_Large(b *testing.B) {
	data := generateRandomBytes(1024 * 1024) // 1 MB
	encoded, err := z3b.Encode(data)
	if err != nil {
		b.Fatalf("Encode failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := z3b.Decode(encoded)
		if err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
		if !bytes.Equal(decoded, data) {
			b.Fatalf("Decoded data does not match original data")
		}
	}
}

func BenchmarkEncodeBase64_Small(b *testing.B) {
	data := generateRandomBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = base64.StdEncoding.EncodeToString(data)
	}
}

func BenchmarkDecodeBase64_Small(b *testing.B) {
	data := generateRandomBytes(32)
	encoded := base64.StdEncoding.EncodeToString(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
		if !bytes.Equal(decoded, data) {
			b.Fatalf("Decoded data does not match original data")
		}
	}
}

func BenchmarkEncodeBase64_Medium(b *testing.B) {
	data := generateRandomBytes(1024) // 1 KB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = base64.StdEncoding.EncodeToString(data)
	}
}

func BenchmarkDecodeBase64_Medium(b *testing.B) {
	data := generateRandomBytes(1024) // 1 KB
	encoded := base64.StdEncoding.EncodeToString(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
		if !bytes.Equal(decoded, data) {
			b.Fatalf("Decoded data does not match original data")
		}
	}
}

func BenchmarkEncodeBase64_Large(b *testing.B) {
	data := generateRandomBytes(1024 * 1024) // 1 MB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = base64.StdEncoding.EncodeToString(data)
	}
}

func BenchmarkDecodeBase64_Large(b *testing.B) {
	data := generateRandomBytes(1024 * 1024) // 1 MB
	encoded := base64.StdEncoding.EncodeToString(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
		if !bytes.Equal(decoded, data) {
			b.Fatalf("Decoded data does not match original data")
		}
	}
}

// Helper function to generate random bytes
func generateRandomBytes(size int) []byte {
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic("Failed to generate random data")
	}
	return data
}
