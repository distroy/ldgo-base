/*
 * Copyright (C) distroy
 */

package slogtype__

import (
	"reflect"

	"github.com/distroy/ldgo-base/ldlog/internal/logref__"
)

func toType[T, S any](v S) T                 { return logref__.ToType[T](v) }
func asType[T any](v any, def ...T) T        { return logref__.AsType(v, def...) }
func checkTypeEqual(this, that reflect.Type) { logref__.CheckTypeEqual(this, that) }
