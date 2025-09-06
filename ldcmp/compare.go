/*
 * Copyright (C) distroy
 */

package ldcmp

import "github.com/distroy/ldgo-base/internal/cmp_"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

type Complex interface {
	~complex64 | ~complex128
}

type Comparable interface {
	Number | ~string
}

type Comparer[T any] interface {
	Compare(b T) int
}

func Compare(a, b any) int { return cmp_.Compare(a, b) }

func CompareBool[T ~bool](a, b T) int { return cmp_.CompareBool(a, b) }

func CompareString[T ~string](a, b T) int { return cmp_.CompareString(a, b) }
func CompareBytes[T ~[]byte](a, b T) int  { return cmp_.CompareBytes(a, b) }

func CompareInteger[T Integer](a, b T) int { return cmp_.CompareInteger(a, b) }
func CompareFloat[T Float](a, b T) int     { return cmp_.CompareComparable(a, b) }
func CompareNumber[T Number](a, b T) int   { return cmp_.CompareComparable(a, b) }
func CompareComplex[T Complex](a, b T) int { return cmp_.CompareComplex(a, b) }

func CompareComparer[T Comparer[T]](a, b T) int  { return cmp_.CompareComparer(a, b) }
func CompareComparable[T Comparable](a, b T) int { return cmp_.CompareComparable(a, b) }
