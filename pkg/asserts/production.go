// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build !asserts
// +build !asserts

package asserts

func Index(i, maxIndex int) {
}

func Coordinate(x, y, w, h int) {
}

func NotEqual[T comparable](value, expected T) {
}

func Equal[T comparable](value, expected T) {
}

func Capacity[T interface{ ~[]E | chan E }, E any](value T, expected int) {
}

func MinCapacity[T interface{ ~[]E | chan E }, E any](value T, expected int) {
}

func Length[T interface{ ~[]E | chan E }, E any](value T, expected int) {
}

func MinLength[T interface{ ~[]E | chan E }, E any](value T, expected int) {
}

//func NotNil[T comparable](value T) {
//}
//
//func Nil[T comparable](value T) {
//}
