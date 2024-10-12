// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package z3b

import (
	"fmt"
	"sync"
)

// Constants
const (
	bufferSize       = 4096
	separator1       = '_' // Set +1 control character
	separator2       = '-' // Set +2 control character
	invalidCharacter = '"' // Invalid character
	printableSet     = "!#$%&()*+,./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^abcdefghijklmnopqrstuvwxyz{|}~"
	numSets          = 3   // Number of binary sets
	setSize          = 86  // Number of characters in printableSet
	totalBytes       = 256 // Total number of byte values (0-255)
)

// Binary Sets: Each set maps to a unique range of byte values.
var (
	binarySet1 = [setSize]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
		0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55,
	}

	binarySet2 = [setSize]byte{
		0x56, 0x57,
		0x58, 0x59,
		0x5A, 0x5B,
		0x5C,
		0x5D, 0x5E, 0x5F, 0x60, 0x61, 0x62, 0x63, 0x64,
		0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C,
		0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74,
		0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C,
		0x7D, 0x7E, 0x7F, 0x80, 0x81, 0x82, 0x83, 0x84, 0x85,
		0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D,
		0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95,
		0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D,
		0x9E, 0x9F, 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5,
		0xA6, 0xA7, 0xA8, 0xA9, 0xAA,
		0xAB,
	}

	binarySet3 = [setSize]byte{
		0xAC, 0xAD, 0xAE, 0xAF,
		0xB0, 0xB1, 0xB2,
		0xB3,
		0xB4, 0xB5, 0xB6, 0xB7,
		0xB8, 0xB9,
		0xBA,
		0xBB, 0xBC, 0xBD, 0xBE, 0xBF, 0xC0, 0xC1, 0xC2,
		0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA,
		0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2,
		0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA,
		0xDB, 0xDC, 0xDD, 0xDE, 0xDF, 0xE0, 0xE1, 0xE2,
		0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA,
		0xEB, 0xEC, 0xED, 0xEE, 0xEF, 0xF0, 0xF1, 0xF2,
		0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA,
		0xFB, 0xFC, 0xFD, 0xFE, 0xFF,
		0x00, 0x01,
	}
)

// Character Sets Arrays
var (
	printableBytes = []byte(printableSet) // Convert printableSet to []byte
)

// Mapping Tables
var (
	// byteToCharSet maps each set to a byte-to-char mapping array.
	byteToCharSet [numSets][totalBytes]byte

	// charToByteSet maps each set to a char-to-byte mapping array.
	// Index 0 corresponds to char 0, up to char 126.
	charToByteSet [numSets][128]byte
)

// Initialize mapping tables
func init() {
	initializeMappings()
}

// initializeMappings populates the byteToCharSet and charToByteSet arrays.
func initializeMappings() {

	// Populate Set1
	for i, b := range binarySet1 {
		if i >= len(printableBytes) {
			break
		}
		r := printableBytes[i]
		byteToCharSet[0][b] = r
		charToByteSet[0][r] = b
	}

	// Populate Set2
	for i, b := range binarySet2 {
		if i >= len(printableBytes) {
			break
		}
		r := printableBytes[i]
		byteToCharSet[1][b] = r
		charToByteSet[1][r] = b
	}

	// Populate Set3
	for i, b := range binarySet3 {
		if i >= len(printableBytes) {
			break
		}
		r := printableBytes[i]
		byteToCharSet[2][b] = r
		charToByteSet[2][r] = b
	}

	// Fill remaining byteToCharSet with 0 char if not mapped
	for set := 0; set < numSets; set++ {
		for b := 0; b < totalBytes; b++ {
			if byteToCharSet[set][b] == 0 {
				byteToCharSet[set][b] = 0
			}
		}
	}

	//// Verify that separator1 is not present in printableBytes
	//for i, r := range printableBytes {
	//	if r == separator1 {
	//		panic(fmt.Sprintf("Separator character is present in index %d", i))
	//	}
	//}
}

// Memory Pools for Encode and Decode
var encodePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, bufferSize)
	},
}

var decodePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, bufferSize)
	},
}

// Encode encodes the input bytes into a z3b-encoded string.
func Encode(data []byte) ([]byte, error) {
	l := len(data)
	estimatedSize := l * 4

	result := encodePool.Get().([]byte)
	defer encodePool.Put(result[:0])
	if cap(result) < estimatedSize {
		encodePool.Put(result[:0])
		result = make([]byte, 0, estimatedSize)
		//result = slices.Grow(result, estimatedSize)
	}
	result = result[:estimatedSize]

	currentSet := 0
	idx := 0
	for _, b := range data {
		r := byteToCharSet[currentSet][b]
		if r == 0 {
			// Switch to next set
			currentSet = (currentSet + 1) % numSets
			r = byteToCharSet[currentSet][b]
			if r == 0 {
				// If still not found, switch to the next set
				currentSet = (currentSet + 1) % numSets
				r = byteToCharSet[currentSet][b]
				if r == 0 {
					return nil, fmt.Errorf("byte value cannot be encoded: %d", b)
				}

				result[idx] = separator2
				idx++
				result[idx] = r
				idx++
				continue
			}

			result[idx] = separator1
			idx++
			result[idx] = r
			idx++
			continue

		}
		result[idx] = r
		idx++
	}

	return result[:idx], nil
}

func ReleaseDecodedBytes(decoded []byte) {
	decodePool.Put(decoded[:0])
}

// Decode decodes a z3b-encoded string back into bytes.
func Decode(encoded []byte) ([]byte, error) {
	l := len(encoded)

	decoded := decodePool.Get().([]byte)
	if cap(decoded) < l {
		decodePool.Put(decoded[:0])
		decoded = make([]byte, 0, l)
		//decoded = slices.Grow(decoded, l)
	}
	decoded = decoded[:l]

	currentSet := 0
	idx := 0
	for i := 0; i < l; i++ {
		c := encoded[i]

		if c == separator1 {
			// Toggle to the next set
			currentSet = (currentSet + 1) % numSets
			continue
		}

		if c == separator2 {
			// Toggle to the next set
			currentSet = (currentSet + 2) % numSets
			continue
		}

		if c < 32 || c > 126 || c == invalidCharacter { // Printable ASCII range check
			return nil, fmt.Errorf("invalid character in encoded string: '%c'", c)
		}

		b := charToByteSet[currentSet][c]
		if b == 0 && binarySet1[0] != 0 { // Assuming 0 is a valid byte, check if mapped
			return nil, fmt.Errorf("invalid character '%c' for set %d", c, currentSet+1)
		}
		decoded[idx] = b
		idx++
	}
	return decoded[:idx], nil
}
