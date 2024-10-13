// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build generate_sets
// +build generate_sets

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

// Initialize mapping tables
func init() {

	setSize := 86
	totalBytes := 256
	allBytes := make([]byte, totalBytes)
	for i := 0; i < totalBytes; i++ {
		allBytes[i] = byte(i)
	}

	RandomShuffle(allBytes)

	set1 := allBytes[0:setSize]
	set2 := allBytes[setSize : setSize*2]
	set3 := allBytes[setSize*2-2:]

	fmt.Printf("Set 1 (%d):\n%s\n\n", len(set1), printSet(set1))
	fmt.Printf("Set 2 (%d):\n%s\n\n", len(set2), printSet(set2))
	fmt.Printf("Set 3 (%d):\n%s\n\n", len(set3), printSet(set3))

}

func printSet(data []byte) string {
	d := make([]byte, len(data))
	copy(d, data)
	sort.Slice(d, func(i, j int) bool {
		return d[i] < d[j]
	})
	str := ""
	idx := 0
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
