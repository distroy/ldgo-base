/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"
	"time"

	"github.com/distroy/ldgo-base/ldlog/internal/attr__"
)

func Skip() Attr { return attr__.Skip() }

func Bool(key string, val bool) Attr             { return attr__.Bool(key, val) }
func String(key string, val string) Attr         { return attr__.String(key, val) }
func Stringer(key string, val fmt.Stringer) Attr { return attr__.Stringer(key, val) }
func Binary(key string, val []byte) Attr         { return attr__.Binary(key, val) }
func ByteString(key string, val []byte) Attr     { return attr__.ByteString(key, val) }

func Int(key string, val int) Attr     { return attr__.Int(key, val) }
func Int8(key string, val int8) Attr   { return attr__.Int64(key, int64(val)) }
func Int16(key string, val int16) Attr { return attr__.Int64(key, int64(val)) }
func Int32(key string, val int32) Attr { return attr__.Int64(key, int64(val)) }
func Int64(key string, val int64) Attr { return attr__.Int64(key, val) }

func Uint(key string, val uint) Attr       { return attr__.Uint64(key, uint64(val)) }
func Uint8(key string, val uint8) Attr     { return attr__.Uint64(key, uint64(val)) }
func Uint16(key string, val uint16) Attr   { return attr__.Uint64(key, uint64(val)) }
func Uint32(key string, val uint32) Attr   { return attr__.Uint64(key, uint64(val)) }
func Uint64(key string, val uint64) Attr   { return attr__.Uint64(key, val) }
func Uintptr(key string, val uintptr) Attr { return attr__.Uint64(key, uint64(val)) }

func Float32(key string, val float32) Attr { return attr__.Float64(key, float64(val)) }
func Float64(key string, val float64) Attr { return attr__.Float64(key, val) }

func Complex64(key string, val complex64) Attr   { return attr__.Complex64(key, val) }
func Complex128(key string, val complex128) Attr { return attr__.Complex128(key, val) }

func Boolp(key string, val *bool) Attr                   { return attr__.Boolp(key, val) }
func Bools(key string, val []bool) Attr                  { return attr__.Bools(key, val) }
func Stringp(key string, val *string) Attr               { return attr__.Stringp(key, val) }
func Strings(key string, val []string) Attr              { return attr__.Strings(key, val) }
func Stringers[T fmt.Stringer](key string, val []T) Attr { return attr__.Stringers(key, val) }
func ByteStrings(key string, val [][]byte) Attr          { return attr__.ByteStrings(key, val) }

func Intp(key string, val *int) Attr     { return attr__.Intp(key, val) }
func Int8p(key string, val *int8) Attr   { return attr__.Int8p(key, val) }
func Int16p(key string, val *int16) Attr { return attr__.Int16p(key, val) }
func Int32p(key string, val *int32) Attr { return attr__.Int32p(key, val) }
func Int64p(key string, val *int64) Attr { return attr__.Int64p(key, val) }

func Uintp(key string, val *uint) Attr       { return attr__.Uintp(key, val) }
func Uint8p(key string, val *uint8) Attr     { return attr__.Uint8p(key, val) }
func Uint16p(key string, val *uint16) Attr   { return attr__.Uint16p(key, val) }
func Uint32p(key string, val *uint32) Attr   { return attr__.Uint32p(key, val) }
func Uint64p(key string, val *uint64) Attr   { return attr__.Uint64p(key, val) }
func Uintptrp(key string, val *uintptr) Attr { return attr__.Uintptrp(key, val) }

func Float32p(key string, val *float32) Attr { return attr__.Float32p(key, val) }
func Float64p(key string, val *float64) Attr { return attr__.Float64p(key, val) }

func Complex64p(key string, val *complex64) Attr   { return attr__.Complex64p(key, val) }
func Complex128p(key string, val *complex128) Attr { return attr__.Complex128p(key, val) }

func Ints(key string, val []int) Attr     { return attr__.Ints(key, val) }
func Int8s(key string, val []int8) Attr   { return attr__.Int8s(key, val) }
func Int16s(key string, val []int16) Attr { return attr__.Int16s(key, val) }
func Int32s(key string, val []int32) Attr { return attr__.Int32s(key, val) }
func Int64s(key string, val []int64) Attr { return attr__.Int64s(key, val) }

func Uints(key string, val []uint) Attr       { return attr__.Uints(key, val) }
func Uint8s(key string, val []uint8) Attr     { return attr__.Uint8s(key, val) }
func Uint16s(key string, val []uint16) Attr   { return attr__.Uint16s(key, val) }
func Uint32s(key string, val []uint32) Attr   { return attr__.Uint32s(key, val) }
func Uint64s(key string, val []uint64) Attr   { return attr__.Uint64s(key, val) }
func Uintptrs(key string, val []uintptr) Attr { return attr__.Uintptrs(key, val) }

func Float32s(key string, val []float32) Attr { return attr__.Float32s(key, val) }
func Float64s(key string, val []float64) Attr { return attr__.Float64s(key, val) }

func Complex64s(key string, val []complex64) Attr   { return attr__.Complex64s(key, val) }
func Complex128s(key string, val []complex128) Attr { return attr__.Complex128s(key, val) }

func Time(key string, val time.Time) Attr         { return attr__.Time(key, val) }
func Duration(key string, val time.Duration) Attr { return attr__.Duration(key, val) }

func Timep(key string, val *time.Time) Attr         { return attr__.Timep(key, val) }
func Durationp(key string, val *time.Duration) Attr { return attr__.Durationp(key, val) }

func Times(key string, val []time.Time) Attr         { return attr__.Times(key, val) }
func Durations(key string, val []time.Duration) Attr { return attr__.Durations(key, val) }

func Error(err error) Attr                  { return attr__.Error(err) }
func Errors(key string, err []error) Attr   { return attr__.Errors(key, err) }
func NamedError(key string, err error) Attr { return attr__.NamedError(key, err) }

func Stack(key string) Attr               { return attr__.StackSkip(key, 1) }
func StackSkip(key string, skip int) Attr { return attr__.StackSkip(key, skip+1) }

func Any(key string, val any) Attr     { return attr__.Any(key, val) }
func Reflect(key string, val any) Attr { return attr__.Reflect(key, val) }

func Integer[Int ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr,
](
	key string, val Int,
) Attr {
	if val >= 0 {
		return Uint64(key, uint64(val))
	}
	return Int64(key, int64(val))
}

func Float[T ~float32 | ~float64](key string, val T) Attr { return Float64(key, float64(val)) }
