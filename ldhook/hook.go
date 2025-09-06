/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"reflect"
)

type Hook interface {
	hook(patches *patches)
}

type FuncHook struct {
	Target any
	Double any // must be func/[]OutputCell/OutputCell/[]Values/Values
}

func (h FuncHook) hook(patches *patches) {
	target := h.Target
	funcType := reflect.TypeOf(h.Target)
	double := getDoubleInterface(funcType, h.Double)

	patches.coreApplyFunc(reflect.ValueOf(target), reflect.ValueOf(double))
}

type MethodHook struct {
	Target any    // object
	Method string // method name
	Double any    // must be func/[]OutputCell/OutputCell/[]Values/Values
}

func (h MethodHook) hook(patches *patches) {
	method := getMethod(h.Target, h.Method)

	mType := method.Type
	double := getDoubleInterface(mType, h.Double)

	patches.coreApplyFunc(method.Func, reflect.ValueOf(double))
}

func getDoubleInterface(funcType reflect.Type, double any) any {
	switch v := double.(type) {
	default:
		return double

	case []OutputCell:
		outs := make([]ResultCell, 0, len(v))
		for _, out := range v {
			outs = append(outs, ResultCell{
				Outputs: out.Values,
				Times:   out.Times,
			})
		}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case []*OutputCell:
		outs := make([]ResultCell, 0, len(v))
		for _, out := range v {
			outs = append(outs, ResultCell{
				Outputs: out.Values,
				Times:   out.Times,
			})
		}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case OutputCell:
		outs := []ResultCell{{
			Outputs: v.Values,
			Times:   v.Times,
		}}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case *OutputCell:
		outs := []ResultCell{{
			Outputs: v.Values,
			Times:   v.Times,
		}}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case []ResultCell:
		return newWrap(funcType, v).MakeFunc().Interface()

	case []*ResultCell:
		outs := make([]ResultCell, 0, len(v))
		for _, out := range v {
			outs = append(outs, *out)
		}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case ResultCell:
		outs := []ResultCell{v}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case *ResultCell:
		outs := []ResultCell{*v}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case []Values:
		outs := make([]ResultCell, 0, len(v))
		for _, out := range v {
			outs = append(outs, ResultCell{Outputs: out})
		}
		return newWrap(funcType, outs).MakeFunc().Interface()

	case Values:
		outs := []ResultCell{
			{Outputs: v},
		}
		return newWrap(funcType, outs).MakeFunc().Interface()
	}
}
