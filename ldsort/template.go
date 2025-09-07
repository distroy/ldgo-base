/*
 * Copyright (C) distroy
 */

package ldsort

import (
	"github.com/distroy/ldgo-base/internal/cmp_"
)

type sortable interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func compare[T sortable](a, b T) int { return cmp_.CompareComparable(a, b) }

type slice[T sortable] []T

func (s slice[T]) Len() int      { return len(s) }
func (s slice[T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s slice[T]) Compare(i, j int) int {
	return compare(s[i], s[j])
}

// templateSearch return the smallest index which a[i] >= x
func templateSearch[T sortable](a []T, x T) int {
	return internalSearch(len(a), func(i int) bool { return compare(a[i], x) >= 0 })
}

func templateIndex[T sortable](a []T, x T) int {
	return internalIndex(len(a), func(i int) int {
		return compare(a[i], x)
	})
}
