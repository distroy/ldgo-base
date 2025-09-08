/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"github.com/distroy/ldgo-base/internal/cmp_"
)

func DefaultCompare[T any](a, b T) int {
	return cmp_.Compare(a, b)
}
