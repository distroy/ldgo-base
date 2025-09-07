/*
 * Copyright (C) distroy
 */

package ldsort

func uniq(s iface) int {
	l := s.Len()
	if l <= 0 {
		return 0
	}

	i := 0
	for j := 1; j < l; j++ {
		if s.Compare(i, j) == 0 {
			continue
		}
		i++
		s.Swap(i, j)
	}
	return i + 1
}

// Uniq uniqs slice with compareFunc
func Uniq[T any](slice []T, compareFunc func(a, b T) int) []T {
	if slice == nil {
		return nil
	}
	i := &sortSliceAndCompare[T]{
		slice:   slice,
		compare: compareFunc,
	}
	n := uniq(i)
	return slice[:n]
}
