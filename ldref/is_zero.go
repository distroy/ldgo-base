/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
)

func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}

	if vv, ok := v.(reflect.Value); ok {
		return IsValZero(vv)
	}

	vv := reflect.ValueOf(v)
	return IsValZero(vv)
}

func IsValZero(v reflect.Value) bool {
	return v.Kind() == reflect.Invalid || v.IsZero()
}
