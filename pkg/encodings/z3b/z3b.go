// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package z3b

import (
	"fmt"
	"sync"
)

// Constants
const (
	bufferSize = 4096 // Default buffer size

	separator1 = '_' // Set +1 set control character
	separator2 = '-' // Set +2 set control character
	separator3 = '$' // Set +3 set control character
	separator4 = '`' // Set +4 set control character
	separator5 = '!' // Set +5 set control character
	separator6 = ' ' // Set +6 set control character

	invalidCharacter = '"' // Invalid character

	printableSet  = "#%&()*+,./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^abcdefghijklmnopqrstuvwxyz{|}~" // Data characters
	printableSet2 = "-_$`! "                                                                                 // These can be used for control
	printableSet3 = "'\"\\"                                                                                  // We don't use these because JSON escaping

	numSets    = 7   // Number of binary sets
	setSize    = 86  // Number of characters in printableSet
	totalBytes = 256 // Total number of byte values (0-255)
)

// Character Sets Arrays
var (
	printableBytes = []byte(printableSet) // Convert printableSet to []byte
)

// Mapping Tables
var (
	// byteToCharSet maps each set to a byte-to-char mapping array.
	byteToCharSet [numSets][totalBytes]int

	// charToByteSet maps each set to a char-to-byte mapping array.
	// Index 0 corresponds to char 0, up to char 128.
	charToByteSet [numSets][128]byte
)

// Initialize mapping tables
func init() {
	initializeMappings()
}

// initializeMappings populates the byteToCharSet and charToByteSet arrays.
func initializeMappings() {

	if len(printableBytes) < setSize {
		panic(fmt.Sprintf("Not enough printable bytes: %d of %d required", len(printableBytes), setSize))
	}

	// Fill byteToCharSet with 0 char if not mapped
	for set := 0; set < numSets; set++ {
		for b := 0; b < totalBytes; b++ {
			byteToCharSet[set][b] = -1
		}
	}

	// Populate Set1
	for i, b := range binarySet1 {
		r := printableBytes[i]
		byteToCharSet[0][b] = int(r)
		charToByteSet[0][r] = b
	}

	// Populate Set2
	for i, b := range binarySet2 {
		r := printableBytes[i]
		byteToCharSet[1][b] = int(r)
		charToByteSet[1][r] = b
	}

	// Populate Set3
	for i, b := range binarySet3 {
		r := printableBytes[i]
		byteToCharSet[2][b] = int(r)
		charToByteSet[2][r] = b
	}

	// Populate Set4
	for i, b := range binarySet4 {
		r := printableBytes[i]
		byteToCharSet[3][b] = int(r)
		charToByteSet[3][r] = b
	}

	// Populate Set5
	for i, b := range binarySet5 {
		r := printableBytes[i]
		byteToCharSet[4][b] = int(r)
		charToByteSet[4][r] = b
	}

	// Populate Set6
	for i, b := range binarySet6 {
		r := printableBytes[i]
		byteToCharSet[5][b] = int(r)
		charToByteSet[5][r] = b
	}

	// Populate Set7
	for i, b := range binarySet7 {
		r := printableBytes[i]
		byteToCharSet[6][b] = int(r)
		charToByteSet[6][r] = b
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
	for i, b := range data {

		r := byteToCharSet[currentSet][b]
		if r >= 0 && r <= totalBytes {
			result[idx] = byte(r)
			idx++
			continue
		}

		ch := make(chan [2]int)
		for j := 1; j < 7; j++ {
			go findChunk(currentSet, j, data[i:], ch)
		}

		bestSet := -1 // The current best set
		bestSize := 0 // How many characters can be presented in the best set
		received := 0
	events:
		for {
			select {
			case res := <-ch:
				j, size := res[0], res[1]
				//log.Printf("[%d]: Received +%d (%d bytes)", received, j, size)
				if bestSize < size {
					bestSet = j
					bestSize = size
				}

				received++
				if received >= 6 {
					break events
				}
			}
		}

		//log.Printf("Best found +%d (%d bytes)", bestSize, bestSize)

		switch bestSet {
		case 1:
			result[idx] = separator1
			idx++
		case 2:
			result[idx] = separator2
			idx++
		case 3:
			result[idx] = separator3
			idx++
		case 4:
			result[idx] = separator4
			idx++
		case 5:
			result[idx] = separator5
			idx++
		case 6:
			result[idx] = separator6
			idx++
		default:
			return nil, fmt.Errorf("failed to find encoding set for: '%c'", b)
		}

		currentSet = (currentSet + bestSet) % numSets
		r = byteToCharSet[currentSet][b]

		result[idx] = byte(r)
		idx++
	}

	return result[:idx], nil
}

func findChunk(currentSet, j int, next []byte, ch chan [2]int) {
	set := (currentSet + j) % numSets
	idx := 0
	for _, b := range next {
		r := byteToCharSet[set][b]
		if r < 0 || r >= totalBytes {
			break
		}
		idx++
	}
	ch <- [2]int{j, idx}
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

		switch c {
		case separator1:
			currentSet = (currentSet + 1) % numSets
			continue
		case separator2:
			currentSet = (currentSet + 2) % numSets
			continue
		case separator3:
			currentSet = (currentSet + 3) % numSets
			continue
		case separator4:
			currentSet = (currentSet + 4) % numSets
			continue
		case separator5:
			currentSet = (currentSet + 5) % numSets
			continue
		case separator6:
			currentSet = (currentSet + 6) % numSets
			continue
		}

		if c < 32 || c > 126 || c == invalidCharacter {
			return nil, fmt.Errorf("invalid character in encoded string: '%c'", c)
		}

		b := charToByteSet[currentSet][c]
		//if b == 0 && binarySet1[0] != 0 { // Assuming 0 is a valid byte, check if mapped
		//	return nil, fmt.Errorf("invalid character '%c' for set %d", c, currentSet+1)
		//}
		decoded[idx] = b
		idx++
	}
	return decoded[:idx], nil
}
