/*
 * Copyright (C) distroy
 */

package ldflag

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/distroy/ldgo-base/internal/time_"
	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/ldptr"
)

type Value interface {
	flag.Value
}

type valueWithDefault interface {
	Value

	Default() string
}

type valueWithMeta interface {
	Value

	Meta() string
}

func mustMarshalJson(v any) string {
	b := bytes.NewBuffer(nil)
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	e.Encode(v)

	s := b.Bytes()
	if l := len(s) - 1; l >= 0 && s[l] == '\n' {
		s = s[:l]
	}

	return ldconv.BytesToStrUnsafe(s)
}

type typeIface[T any] interface {
	String(v T) string
	Parse(s string) (T, error)
}

func getMetaByIface[I typeIface[T], T any](i I) string {
	switch ii := any(i).(type) {
	case interface{ Meta() string }:
		return ii.Meta()
	}
	typ := typeFor[T]().String()
	return fmt.Sprintf("<%s>", typ)
}

type anyVal[I typeIface[T], T any] struct{ V T }

func newAnyVal[I typeIface[T], T any](p *T) *anyVal[I, T] { return (*anyVal[I, T])(unsafe.Pointer(p)) }
func (p *anyVal[I, T]) iface() (i I)                      { return i }
func (p *anyVal[I, T]) Meta() string                      { return getMetaByIface(p.iface()) }
func (p *anyVal[I, T]) String() string                    { return p.iface().String(p.V) }
func (p *anyVal[I, T]) Set(s string) error {
	v, err := p.iface().Parse(s)
	if err != nil {
		return err
	}
	p.V = v
	return nil
}

type sliceVal[I typeIface[T], T any] []T

func newSliceVal[I typeIface[T], T any](p *[]T) *sliceVal[I, T] { return (*sliceVal[I, T])(p) }
func (p *sliceVal[I, T]) iface() (i I)                          { return i }
func (p *sliceVal[I, T]) Meta() string                          { return getMetaByIface(p.iface()) }
func (p *sliceVal[I, T]) String() string                        { return mustMarshalJson(*p) }
func (p *sliceVal[I, T]) Set(s string) error {
	v, err := p.iface().Parse(s)
	if err != nil {
		return err
	}
	*p = append(*p, v)
	return nil
}

type ptrVal[I typeIface[T], T any] struct{ V **T }

func newPtrVal[I typeIface[T], T any](p **T) ptrVal[I, T] { return ptrVal[I, T]{V: p} }
func (p ptrVal[I, T]) iface() (i I)                       { return i }
func (p ptrVal[I, T]) Meta() string                       { return getMetaByIface(p.iface()) }
func (p ptrVal[I, T]) String() string                     { return p.iface().String((ldptr.Get(*p.V))) }
func (p ptrVal[I, T]) Set(s string) error {
	v, err := p.iface().Parse(s)
	if err != nil {
		return err
	}
	*p.V = &v
	return nil
}

type stringIface struct{}

func (stringIface) Meta() string                   { return "<string>" }
func (stringIface) String(v string) string         { return v }
func (stringIface) Parse(s string) (string, error) { return s, nil }

type boolIface struct{}

func (boolIface) Meta() string                 { return "<bool>" }
func (boolIface) String(v bool) string         { return strconv.FormatBool(v) }
func (boolIface) Parse(s string) (bool, error) { return strconv.ParseBool(s) }

type sinteger interface {
	int | int8 | int16 | int32 | int64
}

type sintIface[T sinteger] struct{}

func (sintIface[T]) Meta() string      { return "<int>" }
func (sintIface[T]) String(v T) string { return strconv.FormatInt(int64(v), 10) }
func (sintIface[T]) Parse(s string) (T, error) {
	byteSize := sizeFor[T]()
	v, err := strconv.ParseInt(s, 0, int(byteSize)*8)
	return T(v), err
}

type uinteger interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type uintIface[T uinteger] struct{}

func (uintIface[T]) Meta() string      { return "<uint>" }
func (uintIface[T]) String(v T) string { return strconv.FormatUint(uint64(v), 10) }
func (uintIface[T]) Parse(s string) (T, error) {
	byteSize := sizeFor[T]()
	v, err := strconv.ParseUint(s, 0, int(byteSize)*8)
	return T(v), err
}

type floatIface[T float32 | float64] struct{}

func (floatIface[T]) Meta() string      { return "<float>" }
func (floatIface[T]) String(v T) string { return strconv.FormatFloat(float64(v), 'g', -1, 64) }
func (floatIface[T]) Parse(s string) (T, error) {
	byteSize := sizeFor[T]()
	v, err := strconv.ParseFloat(s, byteSize*8)
	return T(v), err
}

type durationIface struct{}

func (durationIface) Meta() string                  { return "<duration>" }
func (durationIface) String(v time.Duration) string { return v.String() }
func (durationIface) Parse(s string) (time.Duration, error) {
	b := ldconv.StrToBytesUnsafe(s)
	if d, err := time_.DurationUnmarshalJsonByNumber(b); err == nil {
		return d, nil
	}
	return time.ParseDuration(s)
}

// func
type funcValue func(string) error

func newFuncValue(f func(string) error) funcValue { return funcValue(f) }
func (f funcValue) Set(s string) error            { return f(s) }
func (f funcValue) String() string                { return "" }

// bool
type boolValue = anyVal[boolIface, bool]

// type boolFlag bool
type boolFlag struct{ boolValue }

func newBoolFlag(p *boolValue) *boolFlag { return (*boolFlag)(unsafe.Pointer(p)) }
func (p *boolFlag) Meta() string         { return "" }
func (p *boolFlag) String() string       { return p.boolValue.String() }
func (p *boolFlag) Set(s string) error   { return p.boolValue.Set(s) }
func (p *boolFlag) IsBoolFlag() bool     { return true }

// bool ptr
type boolPtrValue = ptrVal[boolIface, bool]

type boolPtrFlag struct{ boolPtrValue }

func newBoolPtrFlag(p boolPtrValue) boolPtrFlag { return boolPtrFlag{p} }
func (p boolPtrFlag) Meta() string              { return "" }
func (p boolPtrFlag) String() string            { return p.boolPtrValue.String() }
func (p boolPtrFlag) Set(s string) error        { return p.boolPtrValue.Set(s) }
func (p boolPtrFlag) IsBoolFlag() bool          { return true }

// func (p *boolPtrValue) Get() interface{} { return bool(*p) }
// func (p *boolPtrValue) IsBoolFlag() bool { return true }
