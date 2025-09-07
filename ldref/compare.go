/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"

	"github.com/distroy/ldgo-base/internal/cmp_"
)

func Compare(a, b interface{}) int {
	return cmp_.Compare(a, b)
}

func CompareReflect(a, b reflect.Value) int {
	return cmp_.CompareReflect(a, b)
}
