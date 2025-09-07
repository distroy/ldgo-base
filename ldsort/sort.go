/*
 * Copyright (C) distroy
 */

package ldsort

import (
	"sort"
)

type sortSliceAndCompare[T any] struct {
	slice   []T
	compare func(a, b T) int
}

func (s sortSliceAndCompare[T]) Len() int             { return len(s.slice) }
func (s sortSliceAndCompare[T]) Swap(i, j int)        { s.slice[i], s.slice[j] = s.slice[j], s.slice[i] }
func (s sortSliceAndCompare[T]) Compare(i, j int) int { return s.compare(s.slice[i], s.slice[j]) }

// Sort sorts slice with compare
func Sort[T any](slice []T, compare func(a, b T) int) {
	if slice == nil {
		return
	}
	// i := getSortInterface(slice, compare)
	i := &sortSliceAndCompare[T]{
		slice:   slice,
		compare: compare,
	}
	internalSort(i)
}

// IsSorted reports whether data is sorted.
func IsSorted[T any](slice []T, compare func(a, b T) int) bool {
	if slice == nil {
		return true
	}
	// i := getSortInterface(slice, compare)
	i := sortSliceAndCompare[T]{
		slice:   slice,
		compare: compare,
	}
	return internalIsSorted(i)
}

type iface interface {
	Len() int
	Swap(i, j int)
	Compare(i, j int) int
}

type sortIface struct {
	iface
}

func (s sortIface) Less(i, j int) bool {
	return s.Compare(i, j) < 0
}

func internalSort(s iface)          { sort.Sort(sortIface{iface: s}) }
func internalIsSorted(s iface) bool { return sort.IsSorted(sortIface{iface: s}) }
