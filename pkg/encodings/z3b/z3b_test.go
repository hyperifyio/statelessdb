// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package z3b_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	"statelessdb/pkg/encodings/z3b"
)

// TestEncodeDecode tests the Encode and Decode functions with various inputs
func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"Empty", []byte{}},
		{"SingleByte", []byte{0x00}},
		{"Hello World", []byte("Hello World")},
		{"AllBytes", func() []byte {
			data := make([]byte, 256)
			for i := 0; i < 256; i++ {
				data[i] = byte(i)
			}
			return data
		}()},
		{"StaticRandomData", []byte{0x01, 0x02, 0x3C, 0x3D, 0x7E, 0x20, 0xFF}},
		{"RandomData", generateRandomBytes(256)},
		{"EdgeValues", []byte{0x00, 0x7F, 0x80, 0xFF}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			encoded, err := z3b.Encode(tc.data)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			b64 := base64.StdEncoding.EncodeToString(tc.data)

			fmt.Printf("\n%s: %d bytes encoded:\n\n    z3b (%d bytes): \"%s\"\n\n    base64 (%d bytes): %s\n",
				tc.name, len(tc.data), len(encoded), string(encoded), len(b64), b64)

			decoded, err := z3b.Decode(encoded)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if !bytes.Equal(decoded, tc.data) {
				t.Fatalf("Decoded data does not match original data.\nOriginal: %v\nDecoded: %v", tc.data, decoded)
			}
		})
	}
}

// TestInvalidDecoding tests decoding of invalid strings
func TestInvalidDecoding(t *testing.T) {
	testCases := []struct {
		name    string
		encoded []byte
	}{
		{"Single byte", []byte("A")},
		{"InvalidCharacter", []byte("ABC\"")},
		//{"InvalidToggle", "A   B"},
		//{"EndsWithSeparator", "ABC "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := z3b.Decode(tc.encoded)
			if err == nil {
				t.Fatalf("Expected error when decoding invalid string '%s', but got none", tc.encoded)
			}
		})
	}
}
