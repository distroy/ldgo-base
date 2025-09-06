/*
 * Copyright (C) distroy
 */

package cmp_

func isNaN[T comparable](a T) bool {
	return a != a
}

type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func CompareComparable[T Comparable](a, b T) int {
	aNaN := isNaN(a)
	bNaN := isNaN(b)
	if aNaN && bNaN {
		return 0
	}

	switch {
	case aNaN || a < b:
		return -1
	case bNaN || a > b:
		return +1
	default:
		return 0
	}
}

type Comparer[T any] interface {
	Compare(b T) int
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func CompareInteger[T Integer](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return +1
	default:
		return 0
	}
}
