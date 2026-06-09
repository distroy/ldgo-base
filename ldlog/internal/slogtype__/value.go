/*
 * Copyright (C) distroy
 */

package slogtype__

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/distroy/ldgo-base/ldlog/internal/buf__"
)

func init() {
	checkTypeEqual(reflect.TypeFor[Value](), reflect.TypeFor[slog.Value]())
}

func GetValuePtr(v *slog.Value) *Value { return toType[*Value](v) }
func GetValue(v slog.Value) Value      { return *GetValuePtr(&v) }

func GetSValue(v Value) slog.Value { return *toType[*slog.Value](&v) }

func CountEmptyGroups(as []Attr) int {
	n := 0
	for _, a := range as {
		if a.Value.IsEmptyGroup() {
			n++
		}
	}
	return n
}

func NewValue(num uint64, any any) Value {
	return Value{
		num: num,
		any: any,
	}
}

type Value struct {
	_ [0]func() // disallow ==
	// num holds the value for Kinds Int64, Uint64, Float64, Bool and Duration,
	// the string length for KindString, and nanoseconds since the epoch for KindTime.
	num uint64
	// If any is of type Kind, then the value is in num as described above.
	// If any is of type *time.Location, then the Kind is Time and time.Time value
	// can be constructed from the Unix nanos in num and the location (monotonic time
	// is not preserved).
	// If any is of type stringptr, then the Kind is String and the string value
	// consists of the length in num and the pointer in any.
	// Otherwise, the Kind is Any and any is the value.
	// (This implies that Attrs cannot store values of type Kind, *time.Location
	// or stringptr.)
	any any
}

func (v *Value) Get() *slog.Value { return toType[*slog.Value](v) }

func (v *Value) DirectlyNum() uint64 { return v.num }
func (v *Value) DirectlyAny() any    { return v.any }

func (v *Value) Kind() Kind { return v.Get().Kind() }

func (v *Value) Any() any       { return v.Get().Any() }
func (v *Value) String() string { return v.Get().String() }

func (v *Value) Int64() int64     { return v.Get().Int64() }
func (v *Value) Uint64() uint64   { return v.Get().Uint64() }
func (v *Value) Float64() float64 { return v.Get().Float64() }

func (v *Value) Bool() bool              { return v.Get().Bool() }
func (v *Value) Duration() time.Duration { return v.Get().Duration() }
func (v *Value) Time() time.Time         { return v.Get().Time() }

func (v *Value) LogValuer() slog.LogValuer { return v.Get().LogValuer() }
func (v *Value) Group() []Attr             { return GetAttrs(v.Get().Group()) }

func (v *Value) Equal(w Value) bool { return v.Get().Equal(*w.Get()) }

func (v *Value) Resolve() (rv Value) { return toType[Value](v.Get().Resolve()) }

// IsEmptyGroup reports whether v is a group that has no attributes.
func (v *Value) IsEmptyGroup() bool {
	if v.Kind() != KindGroup {
		return false
	}
	// We do not need to recursively examine the group's Attrs for emptiness,
	// because GroupValue removed them when the group was constructed, and
	// groups are immutable.
	return len(v.Group()) == 0
}

func (v *Value) Append(dst []byte) []byte {
	switch v.Kind() {
	case KindString:
		return append(dst, v.String()...)
	case KindInt64:
		return strconv.AppendInt(dst, int64(v.num), 10)
	case KindUint64:
		return strconv.AppendUint(dst, v.num, 10)
	case KindFloat64:
		return strconv.AppendFloat(dst, v.Float64(), 'g', -1, 64)
	case KindBool:
		return strconv.AppendBool(dst, v.Bool())
	case KindDuration:
		return append(dst, v.Duration().String()...)
	case KindTime:
		return append(dst, v.Time().String()...)
	case KindGroup:
		return fmt.Append(dst, v.Group())
	case KindAny, KindLogValuer:
		return fmt.Append(dst, v.any)
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
}

func (v *Value) WriteToBuffer(buf *buf__.Buffer) {
	err := v.writeToBuffer(buf)
	if err != nil {
		buf.AppendString(fmt.Sprintf("!ERROR:%v", err))
	}
}

func (v *Value) writeToBuffer(buf *buf__.Buffer) error {
	switch v.Kind() {
	case KindString:
		buf.AppendQuote(v.String())
	case KindInt64:
		buf.AppendInt(v.Int64())
		// *s.buf = strconv.AppendInt(*s.buf, v.Int64(), 10)
	case KindUint64:
		buf.AppendUint(v.Uint64())
		// *s.buf = strconv.AppendUint(*s.buf, v.Uint64(), 10)
	case KindFloat64:
		// json.Marshal is funny about floats; it doesn't
		// always match strconv.AppendFloat. So just call it.
		// That's expensive, but floats are rare.
		buf.AppendFloat(v.Float64(), 64)
	case KindBool:
		buf.AppendBool(v.Bool())
	case KindDuration:
		// Do what json.Marshal does.
		buf.AppendDuration(v.Duration())
	case KindTime:
		buf.AppendTime(v.Time(), "")
		// s.appendTime(v.Time())
	case KindAny:
		vv := v.Any()
		switch m := vv.(type) {
		case io.WriterTo:
			_, err := m.WriteTo(buf)
			return err

		case json.Marshaler:
			return buf.AppendJson(vv)

		case error:
			buf.AppendQuote(m.Error())

		case unsafe.Pointer:
			buf.AppendString("0x")
			*buf = strconv.AppendUint(*buf, uint64(uintptr(m)), 0x10)

		default:
			return buf.AppendJson(vv)
		}
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
	return nil
}
