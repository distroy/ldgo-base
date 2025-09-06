/*
 * Copyright (C) distroy
 */

package ldmath

import "math"

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Number interface {
	Integer | Float
}

type Absable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Abs[T Absable](n T) T {
	if n < 0 {
		n = -n
	}
	return n
}

func Max[T Number](n T, args ...T) T {
	for _, v := range args {
		if n < v {
			n = v
		}
	}
	return n
}

func Min[T Number](n T, args ...T) T {
	for _, v := range args {
		if n > v {
			n = v
		}
	}
	return n
}

func NaN32() float32 { return float32(NaN64()) }
func NaN64() float64 { return math.NaN() }

func Inf32(sign int) float32 { return float32(Inf64(sign)) }
func Inf64(sign int) float64 { return math.Inf(sign) }

func IsNaN[T Float](n T) bool           { return n != n }
func IsInf[T Float](f T, sign int) bool { return math.IsInf(float64(f), sign) }

func Sum[T Number](args ...T) T     { return sum[T](args...) }
func Sum2[R, T Number](args ...T) R { return sum[R](args...) }

func sum[R Number, T Number](args ...T) R {
	var sum R = 0
	for _, v := range args {
		sum += R(v)
	}
	return sum
}
