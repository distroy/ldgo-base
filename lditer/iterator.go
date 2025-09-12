/*
 * Copyright (C) distroy
 */

package lditer

import (
	"iter"
)

type (
	Seq[V any]     = iter.Seq[V]
	Seq2[K, V any] = iter.Seq2[K, V]
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func ToSeq2[Seq ~func(yield func(V) bool), V any](fn Seq) Seq2[int, V] {
	return func(yield func(i int, v V) bool) {
		x := 0
		fn(func(vv V) bool {
			i := x
			x++
			return yield(i, vv)
		})
	}
}

func ToSeqByKey[Seq2 ~func(yield func(k K, v V) bool), K, V any](fn Seq2) Seq[K] {
	return func(yield func(K) bool) {
		fn(func(k K, v V) bool { return yield(k) })
	}
}

func ToSeqByValue[Seq2 ~func(yield func(k K, v V) bool), K, V any](fn Seq2) Seq[V] {
	return func(yield func(V) bool) {
		fn(func(k K, v V) bool { return yield(v) })
	}
}

func Chan[C interface{ ~<-chan V | ~chan V }, V any](ch C) Seq[V] {
	return func(yield func(V) bool) {
		for v := range ch {
			if !yield(v) {
				break
			}
		}
	}
}

// Int(n) return [0, n)
// Int(begin, end) return [begin, end)
// Int(begin, end, step) return [begin, end) by step
func Int[T Integer](n T, ns ...T) Seq[T] {
	var (
		begin = T(0)
		end   = n
		step  = T(1)
	)
	if len(ns) > 0 {
		begin = n
		end = ns[0]
	}
	if len(ns) > 1 {
		step = ns[1]
	}
	return func(yield func(T) bool) {
		for v := begin; v < end; v += step {
			if !yield(v) {
				break
			}
		}
	}
}

func Slice[S ~[]V, V any](s S) Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range s {
			if !yield(i, v) {
				break
			}
		}
	}
}

func SliceBackward[S ~[]V, V any](s S) Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				break
			}
		}
	}
}

func Map[M ~map[K]V, K comparable, V any](m M) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				break
			}
		}
	}
}

func String[S ~string](s S) Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		for i, v := range s {
			if !yield(i, v) {
				break
			}
		}
	}
}
