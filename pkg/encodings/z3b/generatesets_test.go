// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build disabled
// +build disabled

package z3b_test

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// RNG is a random number generator used for move selection randomness.
var RNG *rand.Rand

func init() {
	// Initialize the random number generator with a seed based on the current time
	RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomShuffle shuffles a slice of any type in place.
// Parameters:
// - list: The slice to shuffle.
func RandomShuffle[T any](list []T) {
	RNG.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}

const (
	numGroups  = 92
	numSets    = 7
	setSize    = 86
	totalBytes = 256
)

func createSet(in []byte) [setSize]byte {
	ret := [setSize]byte{}
	for i, v := range in {
		ret[i] = v
	}
	return ret
}

func compareSet(a, b [setSize]byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func createThreeRandomSets() ([setSize]byte, [setSize]byte, [setSize]byte) {
	allBytes := make([]byte, totalBytes)
	for i := 0; i < totalBytes; i++ {
		allBytes[i] = byte(i)
	}
	RandomShuffle(allBytes)
	set1 := createSet(allBytes[0:setSize])
	set2 := createSet(allBytes[setSize : setSize*2])
	set3 := createSet(allBytes[setSize*2-2:])
	return set1, set2, set3
}

func createSets() ([setSize]byte, [setSize]byte, [setSize]byte, [setSize]byte, [setSize]byte, [setSize]byte, [setSize]byte) {
	for {
		set1, set2, set3 := createThreeRandomSets()
		for {

			set4, set5, set6 := createThreeRandomSets()
			if compareSet(set4, set1) || compareSet(set4, set2) || compareSet(set4, set3) ||
				compareSet(set5, set1) || compareSet(set5, set2) || compareSet(set5, set3) ||
				compareSet(set6, set1) || compareSet(set6, set2) || compareSet(set6, set3) {
				continue
			}

			for {
				set7, _, _ := createThreeRandomSets()
				if compareSet(set7, set1) || compareSet(set7, set2) || compareSet(set7, set3) ||
					compareSet(set7, set4) || compareSet(set7, set5) || compareSet(set7, set6) {
					continue
				}
				return set1, set2, set3, set4, set5, set6, set7
			}
		}
	}
}

func compareBytes[T byte](a, b [setSize]T) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Returns true if these sets are identical
func compareSets(a, b [numSets][setSize]byte) bool {
	return (compareBytes(a[0], b[0]) ||
		compareBytes(a[0], b[1]) ||
		compareBytes(a[0], b[2]) ||
		compareBytes(a[0], b[3]) ||
		compareBytes(a[0], b[4]) ||
		compareBytes(a[0], b[5]) ||
		compareBytes(a[0], b[6])) &&
		(compareBytes(a[1], b[0]) ||
			compareBytes(a[1], b[1]) ||
			compareBytes(a[1], b[2]) ||
			compareBytes(a[1], b[3]) ||
			compareBytes(a[1], b[4]) ||
			compareBytes(a[1], b[5]) ||
			compareBytes(a[1], b[6])) &&
		(compareBytes(a[2], b[0]) ||
			compareBytes(a[2], b[1]) ||
			compareBytes(a[2], b[2]) ||
			compareBytes(a[2], b[3]) ||
			compareBytes(a[2], b[4]) ||
			compareBytes(a[2], b[5]) ||
			compareBytes(a[2], b[6])) &&
		(compareBytes(a[3], b[0]) ||
			compareBytes(a[3], b[1]) ||
			compareBytes(a[3], b[2]) ||
			compareBytes(a[3], b[3]) ||
			compareBytes(a[3], b[4]) ||
			compareBytes(a[3], b[5]) ||
			compareBytes(a[3], b[6])) &&
		(compareBytes(a[4], b[0]) ||
			compareBytes(a[4], b[1]) ||
			compareBytes(a[4], b[2]) ||
			compareBytes(a[4], b[3]) ||
			compareBytes(a[4], b[4]) ||
			compareBytes(a[4], b[5]) ||
			compareBytes(a[4], b[6])) &&
		(compareBytes(a[5], b[0]) ||
			compareBytes(a[5], b[1]) ||
			compareBytes(a[5], b[2]) ||
			compareBytes(a[5], b[3]) ||
			compareBytes(a[5], b[4]) ||
			compareBytes(a[5], b[5]) ||
			compareBytes(a[5], b[6])) &&
		(compareBytes(a[6], b[0]) ||
			compareBytes(a[6], b[1]) ||
			compareBytes(a[6], b[2]) ||
			compareBytes(a[6], b[3]) ||
			compareBytes(a[6], b[4]) ||
			compareBytes(a[6], b[5]) ||
			compareBytes(a[6], b[6]))
}

// Initialize mapping tables by random
func createRandomMaps() {

	str := "var (\n    binaryGroups = [numGroups][numSets][setSize]byte{\n"

	allSets := make([][numSets][setSize]byte, 0, 92)
	for g := 0; g < 92; g++ {
	loop:
		for {
			set1, set2, set3, set4, set5, set6, set7 := createSets()
			set := [numSets][setSize]byte{set1, set2, set3, set4, set5, set6, set7}
			for _, s := range allSets {
				if compareSets(set, s) {
					continue loop
				}
			}
			allSets = append(allSets, set)
			break
		}
	}

	for i, g := range allSets {
		str += fmt.Sprintf(`
        // Group %d
        {
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
        },
`,
			i,
			printSet(g[0]),
			printSet(g[1]),
			printSet(g[2]),
			printSet(g[3]),
			printSet(g[4]),
			printSet(g[5]),
			printSet(g[6]),
		)

	}

	str += "    }\n)\n"

	fmt.Println(str)

}

func printSet(data [86]byte) string {
	d := make([]byte, len(data))
	idx := 0
	for _, b := range data {
		d[idx] = b
		idx++
	}
	d = d[:idx]
	sort.Slice(d, func(i, j int) bool {
		return d[i] < d[j]
	})
	str := ""
	idx = 0
	for _, v := range d {
		if str == "" {
			str = fmt.Sprintf("0x%02x, ", v)
		} else {
			str = fmt.Sprintf("%s0x%02x, ", str, v)
		}
		idx++
		if idx == 8 {
			str = str + "\n"
			idx = 0
		}
	}
	return str
}

// GetBinaryGroups generates statistically and logically optimized mappings
// for encoding random data. It returns a three-dimensional byte array
// representing groups, sets, and their content.
func GetBinaryGroups() [][numSets][setSize]byte {
	totalPositions := numSets * setSize
	if totalPositions < 256 {
		panic("Total positions per group must be at least 256 to cover all byte values.")
	}

	// Initialize the groups slice
	groups := make([][numSets][setSize]byte, numGroups)

	for g := 0; g < numGroups; g++ {

		// Initialize the sets slice for this group
		sets := [numSets][setSize]byte{
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
		}

		// Generate a permutation of byte values for this group
		byteValues := make([]byte, 256)
		for i := 0; i < 256; i++ {
			j := i + g
			if j >= 256 {
				j -= 256
			}
			byteValues[i] = byte(j)
		}
		// Shuffle the byte values to create a unique permutation per group
		RandomShuffle(byteValues)

		// Calculate how many times each byte value should repeat
		timesPerByte := totalPositions / 256
		remaining := totalPositions % 256

		// Build a slice containing all byte values repeated accordingly
		bytesPerGroup := make([]byte, 0, totalPositions)
		for i := 0; i < timesPerByte; i++ {
			bytesPerGroup = append(bytesPerGroup, byteValues...)
		}
		bytesPerGroup = append(bytesPerGroup, byteValues[:remaining]...)

		// Shuffle the bytesPerGroup to distribute duplicates
		RandomShuffle(bytesPerGroup)

		// Distribute bytesPerGroup into sets
		pos := 0
		for s := 0; s < numSets; s++ {
			from := bytesPerGroup[pos : pos+setSize]
			for i := 0; i < setSize; i++ {
				sets[s][i] = from[i]
			}
			pos += setSize
		}

		//fmt.Printf("sets[0] = %v\n", sets[0])

		// Assign the sets to the group
		groups[g] = sets
	}

	return groups
}

// Initialize mapping tables by using GetBinaryGroups()
func createMaps() {

	allSets := GetBinaryGroups()

	str := "var (\n    binaryGroups = [numGroups][numSets][setSize]byte{\n"
	for i, g := range allSets {
		str += fmt.Sprintf(`
        // Group %d
        {
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
            {
              %s
            },
        },
`,
			i,
			printSet(g[0]),
			printSet(g[1]),
			printSet(g[2]),
			printSet(g[3]),
			printSet(g[4]),
			printSet(g[5]),
			printSet(g[6]),
		)
	}
	str += "    }\n)\n"
	fmt.Println(str)

}

func init() {
	createMaps()
}
