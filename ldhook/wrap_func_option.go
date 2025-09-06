/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"reflect"
)

type FuncWrapperOption func(p *funcWrapper)

func SwapInput(a, b int) FuncWrapperOption {
	return func(p *funcWrapper) {
		p.SwapInput(a, b)
	}
}

func SwapOutput(a, b int) FuncWrapperOption {
	return func(p *funcWrapper) {
		p.SwapOutput(a, b)
	}
}

func AppendInput(val any, typ ...reflect.Type) FuncWrapperOption {
	return func(p *funcWrapper) {
		p.AppendInput(val, typ...)
	}
}

func AppendOutput(val any, typ ...reflect.Type) FuncWrapperOption {
	return func(p *funcWrapper) {
		p.AppendOutput(val, typ...)
	}
}

func InsertInput(pos int, val any, typ ...reflect.Type) FuncWrapperOption {
	return func(p *funcWrapper) {
		p.InsertInput(pos, val, typ...)
	}
}

func InsertOutput(pos int, val any, typ ...reflect.Type) FuncWrapperOption {
	return func(p *funcWrapper) {
		p.InsertOutput(pos, val, typ...)
	}
}
