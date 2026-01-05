/*
 * Copyright (C) distroy
 */

package ldmetric

import (
	"reflect"
)

func Report[T any](reporter T) {
	v := reflect.ValueOf(reporter)
	typ := getTypeInfo(v.Type())
	typ.Report(v)
}

func ResetReporter[T any](reporter T) {
	v := reflect.ValueOf(reporter)
	typ := getTypeInfo(v.Type())
	typ.Reset()
}
