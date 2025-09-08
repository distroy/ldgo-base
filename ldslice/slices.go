/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"slices"

	"github.com/distroy/ldgo-base/internal/cmp_"
)

func IndexFunc[S ~[]V, V any](s S, f func(v V) bool) int {
	return slices.IndexFunc(s, f)
}
func Index[S ~[]V, V any](s S, x V) int {
	return slices.IndexFunc(s, func(v V) bool {
		return cmp_.Compare(v, x) == 0
	})
}

func ContainsFunc[S ~[]V, V any](s S, f func(v V) bool) bool { return IndexFunc(s, f) >= 0 }
func Contains[S ~[]V, V any](s S, v V) bool                  { return Index(s, v) >= 0 }
