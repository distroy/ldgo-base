/*
 * Copyright (C) distroy
 */

package ldslice

import "github.com/distroy/ldgo-base/lditer"

func SplitByCount[S ~[]V, V any](s S, n int) []S {
	l := len(s)
	if l == 0 {
		return nil
	}
	if n >= l || n <= 0 {
		return []S{s}
	}
	r := make([]S, 0, (l+n-1)/n)
	for i := 0; i < l; i += n {
		b := i
		e := i + n
		e = min(e, l)
		r = append(r, s[b:e])
		// log.Printf("b:%d, e:%d, r:%v", b, e, r)
	}
	return r
}

func SplitByCountFunc[S ~[]V, V any](s S, n int, f func(idx int, val S) bool) int {
	l := len(s)
	if l == 0 {
		return 0
	}
	// if n >= l || n <= 0 {
	// 	if f != nil {
	// 		f(0, s)
	// 	}
	// 	return 1
	// }
	count := (l + n - 1) / n
	if f == nil {
		return count
	}
	idx := 0
	for i := 0; i < l; i += n {
		b := i
		e := i + n
		e = min(e, l)
		ss := s[b:e]
		if ok := f(idx, ss); !ok {
			break
		}
		idx++
		// log.Printf("b:%d, e:%d, r:%v", b, e, r)
	}
	return count
}

func SplitByCountIter[S ~[]V, V any](s S, n int) lditer.Seq2[int, S] {
	return func(yield func(int, S) bool) {
		l := len(s)
		if l == 0 {
			return
		}
		idx := 0
		for i := 0; i < l; i += n {
			b := i
			e := i + n
			e = min(e, l)
			ss := s[b:e]
			if ok := yield(idx, ss); !ok {
				break
			}
			idx++
		}
	}
}
