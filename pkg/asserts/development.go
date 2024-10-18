// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build asserts
// +build asserts

package asserts

import (
	"cmp"
	"log"
	"runtime"
)

var stack []byte

func init() {
	stack = make([]byte, 512)
}

func Index[T cmp.Ordered](i, maxIndex T) {
	if i < T(0) || i >= maxIndex {
		runtime.Stack(stack, false)
		log.Fatalf("Index out of boundaries: %d (0..%d)\n\nStack is:\n%s", i, maxIndex, string(stack))
	}
}

func Coordinate[T cmp.Ordered](x, y, w, h T) {
	if x < T(0) || x >= w || y < T(0) || y >= h {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Coordinate out of boundaries: %dx%d (%d x %d)\n\nStack is:\n%s", x, y, w, h, string(stack))
	}
}

func NotEqual[T comparable](value, expected T) {
	if value == expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Not equal assert failed: Got %v, expected not to be %v. \n\nStack is:\n%s", value, expected, string(stack))
	}
}

func Equal[T comparable](value, expected T) {
	if value != expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Equal assert failed: Got %v, expected to be %v. \n\nStack is:\n%s", value, expected, string(stack))
	}
}

func GreaterOrEqual[T cmp.Ordered](value, expected T) {
	if value < expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Assert failed: Expected %v >= %v. \n\nStack is:\n%s", value, expected, string(stack))
	}
}

func Greater[T cmp.Ordered](value, expected T) {
	if value <= expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Assert failed: Expected %v > %v. \n\nStack is:\n%s", value, expected, string(stack))
	}
}

func Less[T cmp.Ordered](value, expected T) {
	if value >= expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Assert failed: Expected %v < %v. \n\nStack is:\n%s", value, expected, string(stack))
	}
}

func LessOrEqual[T cmp.Ordered](value, expected T) {
	if value > expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Assert failed: Expected %v <= %v. \n\nStack is:\n%s", value, expected, string(stack))
	}
}

func Capacity[T interface{ ~[]E | chan E }, E any](value T, expected int) {
	if cap(value) != expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Capasity of slice incorrect: Got %d, expected to be %d. \n\nStack is:\n%s", cap(value), expected, string(stack))
	}
}

func MinCapacity[T interface{ ~[]E | chan E }, E any](value T, expected int) {
	if cap(value) < expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Capasity of slice too low: Got %d, expected to be at least %d. \n\nStack is:\n%s", cap(value), expected, string(stack))
	}
}

func Length[T interface{ ~[]E | chan E }, E any](value T, expected int) {
	if len(value) != expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Length of slice incorrect: Got %d, expected to be %d. \n\nStack is:\n%s", len(value), expected, string(stack))
	}
}

func MinLength[T interface{ ~[]E | chan E }, E any](value T, expected int) {
	if len(value) < expected {
		runtime.Stack(stack, false)
		log.Fatalf("FATAL ERROR: Length of slice too low: Got %d, expected to be at least %d. \n\nStack is:\n%s", len(value), expected, string(stack))
	}
}

//func NotNil[T any](value T) {
//	if value == nil {
//		runtime.Stack(stack, false)
//		log.Fatalf("FATAL ERROR: Not nil assert failed: Got %v, expected not to be nil. \n\nStack is:\n%s", value, string(stack))
//	}
//}
//
//func Nil[T any](value T) {
//	if value != expected {
//		runtime.Stack(stack, false)
//		log.Fatalf("FATAL ERROR: Nil assert failed: Got %v, expected to be nil. \n\nStack is:\n%s", value, string(stack))
//	}
//}
