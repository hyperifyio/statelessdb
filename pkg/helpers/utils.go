// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package helpers

import (
	"math/bits"
	"time"
)

// Abs function to calculate absolute value of an integer
func Abs(a int) int {
	mask := a >> (bits.UintSize - 1)
	return (a ^ mask) - mask
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// ContainsZero checks if the slice contains a 0.
func ContainsZero(s []int) bool {
	for _, value := range s {
		if value == 0 {
			return true
		}
	}
	return false
}

func CountZero(s []int) int {
	count := 0
	for _, value := range s {
		if value == 0 {
			count++
		}
	}
	return count
}

func MillisToISO(timestampMs int64) string {
	t := time.Unix(0, timestampMs*int64(time.Millisecond)).UTC()
	return t.Format(time.RFC3339)
}

// CompareSlices compares two slices of comparable items
func CompareSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// CompareMaps compares two slices of comparable items
func CompareMaps[T, I comparable](map1, map2 map[I]T) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if val, ok := map2[key]; !ok || val != value {
			return false
		}
	}
	return true
}
