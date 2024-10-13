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

	groupSet = printableSet + printableSet2 // Data characters

	numGroups  = len(groupSet)     // Number of encoding groups to select the best from (this is the first byte)
	numSets    = 7                 // Number of binary sets per a encoding group
	setSize    = len(printableSet) // Number of characters in printableSet
	totalBytes = 256               // Total number of byte values (0-255)
)

// Character Sets Arrays
var (
	printableBytes = []byte(printableSet) // Convert printableSet to []byte
	groupBytes     = []byte(groupSet)     // Convert groupSet to []byte
)

// Mapping Tables
var (
	// byteToCharSet maps each set to a byte-to-char mapping array.
	byteToCharSet [numGroups][numSets][totalBytes]int

	// charToByteSet maps each set to a char-to-byte mapping array.
	// Index 0 corresponds to char 0, up to char 128.
	charToByteSet [numGroups][numSets][128]byte
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
	for g := 0; g < numGroups; g++ {
		for set := 0; set < numSets; set++ {
			for b := 0; b < totalBytes; b++ {
				byteToCharSet[g][set][b] = -1
			}
		}
	}

	for g := 0; g < numGroups; g++ {
		for set := 0; set < numSets; set++ {
			for i, b := range binaryGroups[g][set] {
				r := printableBytes[i]
				byteToCharSet[g][set][b] = int(r)
				charToByteSet[g][set][r] = b
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

type groupResult struct {
	err     error
	group   int
	encoded []byte
}

// Encode encodes the input bytes into a z3b-encoded string.
func Encode(data []byte) ([]byte, error) {

	//log.Printf("Encode: data = %v", data)

	if len(data) == 0 {
		return nil, nil
	}

	ch := make(chan groupResult)
	for g := 0; g < numGroups; g++ {
		go findGroup(g, data, ch)
	}

	bestGroup := -1
	bestSize := -1
	var bestEncoding []byte

	for g := 0; g < numGroups; g++ {
		res := <-ch
		if res.err != nil {
			return nil, fmt.Errorf("error in encoding group %d: %w", res.group, res.err)
		}
		size := len(res.encoded)
		if bestSize < 0 || size < bestSize {
			//if cap(bestEncoding) >= 1 {
			//	encodePool.Put(bestEncoding[:0])
			//}
			bestSize = size
			bestGroup = res.group
			bestEncoding = res.encoded
			//} else {
			//	if cap(bestEncoding) >= 1 {
			//		encodePool.Put(bestEncoding[:0])
			//	}

			if bestSize <= 2 {
				break
			}

		}
	}

	if bestGroup < 0 {
		return nil, fmt.Errorf("failed to find best encoding group")
	}

	//log.Printf("Encode: Best group %d '%c' (%d bytes): \"%s\"",
	//	bestGroup, groupBytes[bestGroup], bestSize, string(bestEncoding))

	return bestEncoding, nil
}

func findGroup(group int, data []byte, ch chan groupResult) {
	encoded, err := encodeGroup(group, data)
	if err != nil {
		ch <- groupResult{
			group: group,
			err:   err,
		}
	} else {
		ch <- groupResult{
			group:   group,
			encoded: encoded,
		}
	}
}

// encodeGroup encodes the input bytes into a z3b-encoded string.
func encodeGroup(currentGroup int, data []byte) ([]byte, error) {

	l := len(data)
	estimatedSize := l * 4

	//result := encodePool.Get().([]byte)
	//if cap(result) < estimatedSize {
	//	encodePool.Put(result[:0])
	//	result = make([]byte, 0, estimatedSize)
	//	//result = slices.Grow(result, estimatedSize)
	//}
	//result = result[:estimatedSize]

	result := make([]byte, estimatedSize)
	//log.Printf("encodeGroup: start: group %d '%c'", currentGroup, groupBytes[currentGroup])
	idx := 0
	result[idx] = groupBytes[currentGroup]
	idx++

	currentSet := 0
	for i, b := range data {

		r := byteToCharSet[currentGroup][currentSet][b]
		if r >= 0 && r < totalBytes {
			result[idx] = byte(r)
			idx++
			continue
		}

		ch := make(chan [2]int)
		for j := 1; j < numSets; j++ {
			go findChunk(currentGroup, currentSet, j, data[i:], ch)
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
			//encodePool.Put(result[:0])
			return nil, fmt.Errorf("failed to find encoding set for: '%c'", b)
		}

		currentSet = (currentSet + bestSet) % numSets
		r = byteToCharSet[currentGroup][currentSet][b]

		result[idx] = byte(r)
		idx++
	}

	result = result[:idx]

	//log.Printf("encodeGroup: Group %d '%c': \"%s\"\n", currentGroup, groupBytes[currentGroup], result)

	return result, nil
}

func findChunk(currentGroup, currentSet, j int, next []byte, ch chan [2]int) {
	set := (currentSet + j) % numSets
	idx := 0
	for _, b := range next {
		r := byteToCharSet[currentGroup][set][b]
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

	//log.Printf("Decode: encoded = \"%s\"", string(encoded))

	l := len(encoded)
	if l <= 0 {
		return nil, nil
	}
	if l == 1 {
		return nil, fmt.Errorf("invalid single character string to decode")
	}

	gb := encoded[0]
	currentGroup := -1
	for i, b := range groupBytes {
		if gb == b {
			currentGroup = i
			break
		}
	}
	if currentGroup < 0 {
		return nil, fmt.Errorf("invalid character in encoded string: '%c'", gb)
	}
	//log.Printf("Decode: currentGroup = %d", currentGroup)

	decoded := decodePool.Get().([]byte)
	if cap(decoded) < l {
		decodePool.Put(decoded[:0])
		decoded = make([]byte, 0, l)
		//decoded = slices.Grow(decoded, l)
	}
	decoded = decoded[:l]

	currentSet := 0
	idx := 0
	for i := 1; i < l; i++ {
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

		b := charToByteSet[currentGroup][currentSet][c]
		//if b == 0 && binarySet1[0] != 0 { // Assuming 0 is a valid byte, check if mapped
		//	return nil, fmt.Errorf("invalid character '%c' for set %d", c, currentSet+1)
		//}
		decoded[idx] = b
		idx++
	}

	//log.Printf("Decode: decoded = %v", decoded[:idx])

	return decoded[:idx], nil
}
