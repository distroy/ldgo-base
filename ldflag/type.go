/*
 * Copyright (C) distroy
 */

package ldflag

import (
	"reflect"
	"time"
	"unsafe"
)

var (
	typeFunc = reflect.TypeOf((func(string) error)(nil))
)

func typeFor[T any]() reflect.Type { return reflect.TypeFor[T]() }
func sizeFor[T any]() int {
	var x T
	return int(unsafe.Sizeof(x))
}

type newValFuncType = func(val reflect.Value) Value

func wrapNewValFunc[T any, V valueWithMeta](fn func(p T) V) newValFuncType {
	return func(val reflect.Value) Value {
		return fn(val.Addr().Interface().(T))
	}
}

var newValFuncMap = map[reflect.Type]newValFuncType{
	typeFunc: func(val reflect.Value) Value { return newFuncValue(val.Interface().(func(s string) error)) },

	typeFor[time.Duration]():  wrapNewValFunc(newAnyVal[durationIface]),
	typeFor[*time.Duration](): wrapNewValFunc(newPtrVal[durationIface]),

	typeFor[bool]():  wrapNewValFunc(newAnyVal[boolIface, bool]),
	typeFor[*bool](): wrapNewValFunc(newPtrVal[boolIface, bool]),

	typeFor[string]():   wrapNewValFunc(newAnyVal[stringIface]),
	typeFor[*string]():  wrapNewValFunc(newPtrVal[stringIface]),
	typeFor[[]string](): wrapNewValFunc(newSliceVal[stringIface]),

	typeFor[int]():     wrapNewValFunc(newAnyVal[sintIface[int]]),
	typeFor[int8]():    wrapNewValFunc(newAnyVal[sintIface[int8]]),
	typeFor[int16]():   wrapNewValFunc(newAnyVal[sintIface[int16]]),
	typeFor[int32]():   wrapNewValFunc(newAnyVal[sintIface[int32]]),
	typeFor[int64]():   wrapNewValFunc(newAnyVal[sintIface[int64]]),
	typeFor[uint]():    wrapNewValFunc(newAnyVal[uintIface[uint]]),
	typeFor[uint8]():   wrapNewValFunc(newAnyVal[uintIface[uint8]]),
	typeFor[uint16]():  wrapNewValFunc(newAnyVal[uintIface[uint16]]),
	typeFor[uint32]():  wrapNewValFunc(newAnyVal[uintIface[uint32]]),
	typeFor[uint64]():  wrapNewValFunc(newAnyVal[uintIface[uint64]]),
	typeFor[uintptr](): wrapNewValFunc(newAnyVal[uintIface[uintptr]]),
	typeFor[float32](): wrapNewValFunc(newAnyVal[floatIface[float32]]),
	typeFor[float64](): wrapNewValFunc(newAnyVal[floatIface[float64]]),

	typeFor[[]int]():     wrapNewValFunc(newSliceVal[sintIface[int]]),
	typeFor[[]int8]():    wrapNewValFunc(newSliceVal[sintIface[int8]]),
	typeFor[[]int16]():   wrapNewValFunc(newSliceVal[sintIface[int16]]),
	typeFor[[]int32]():   wrapNewValFunc(newSliceVal[sintIface[int32]]),
	typeFor[[]int64]():   wrapNewValFunc(newSliceVal[sintIface[int64]]),
	typeFor[[]uint]():    wrapNewValFunc(newSliceVal[uintIface[uint]]),
	typeFor[[]uint8]():   wrapNewValFunc(newSliceVal[uintIface[uint8]]),
	typeFor[[]uint16]():  wrapNewValFunc(newSliceVal[uintIface[uint16]]),
	typeFor[[]uint32]():  wrapNewValFunc(newSliceVal[uintIface[uint32]]),
	typeFor[[]uint64]():  wrapNewValFunc(newSliceVal[uintIface[uint64]]),
	typeFor[[]uintptr](): wrapNewValFunc(newSliceVal[uintIface[uintptr]]),
	typeFor[[]float32](): wrapNewValFunc(newSliceVal[floatIface[float32]]),
	typeFor[[]float64](): wrapNewValFunc(newSliceVal[floatIface[float64]]),

	typeFor[*int]():     wrapNewValFunc(newPtrVal[sintIface[int]]),
	typeFor[*int8]():    wrapNewValFunc(newPtrVal[sintIface[int8]]),
	typeFor[*int16]():   wrapNewValFunc(newPtrVal[sintIface[int16]]),
	typeFor[*int32]():   wrapNewValFunc(newPtrVal[sintIface[int32]]),
	typeFor[*int64]():   wrapNewValFunc(newPtrVal[sintIface[int64]]),
	typeFor[*uint]():    wrapNewValFunc(newPtrVal[uintIface[uint]]),
	typeFor[*uint8]():   wrapNewValFunc(newPtrVal[uintIface[uint8]]),
	typeFor[*uint16]():  wrapNewValFunc(newPtrVal[uintIface[uint16]]),
	typeFor[*uint32]():  wrapNewValFunc(newPtrVal[uintIface[uint32]]),
	typeFor[*uint64]():  wrapNewValFunc(newPtrVal[uintIface[uint64]]),
	typeFor[*uintptr](): wrapNewValFunc(newPtrVal[uintIface[uintptr]]),
	typeFor[*float32](): wrapNewValFunc(newPtrVal[floatIface[float32]]),
	typeFor[*float64](): wrapNewValFunc(newPtrVal[floatIface[float64]]),
}
