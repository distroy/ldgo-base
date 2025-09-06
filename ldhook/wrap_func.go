/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"fmt"
	"reflect"
)

func WrapFunc(fn any, options ...FuncWrapperOption) any {
	p := NewFuncWrapper(fn, options...)
	return p.Interface()
}

func NewFuncWrapper(fn any, options ...FuncWrapperOption) *funcWrapper {
	val := reflect.ValueOf(fn)
	typ := val.Type()
	if val.Kind() != reflect.Func {
		panic(fmt.Sprintf("should be func, but be `%s`", typ.String()))
	}

	p := &funcWrapper{}
	p.init(val, typ)

	for _, opt := range options {
		opt(p)
	}

	return p
}

type funcWrapperParameter struct {
	typ   reflect.Type
	value reflect.Value
	pos   int
}

type funcWrapper struct {
	value    reflect.Value
	typ      reflect.Type
	variadic bool
	inputs   []*funcWrapperParameter
	outputs  []*funcWrapperParameter
}

func (p *funcWrapper) init(val reflect.Value, typ reflect.Type) {
	inputs := make([]*funcWrapperParameter, 0, typ.NumIn())
	for i, n := 0, typ.NumIn(); i < n; i++ {
		inT := typ.In(i)

		inputs = append(inputs, &funcWrapperParameter{
			typ: inT,
			pos: i,
		})
	}

	outputs := make([]*funcWrapperParameter, 0, typ.NumOut())
	for i, n := 0, typ.NumOut(); i < n; i++ {
		inT := typ.Out(i)

		outputs = append(outputs, &funcWrapperParameter{
			typ: inT,
			pos: i,
		})
	}

	*p = funcWrapper{
		value:    val,
		typ:      typ,
		variadic: typ.IsVariadic(),
		inputs:   inputs,
		outputs:  outputs,
	}
}

func (p *funcWrapper) getParamterTypes(l []*funcWrapperParameter) []reflect.Type {
	res := make([]reflect.Type, 0, len(l))
	for _, v := range l {
		res = append(res, v.typ)
	}

	return res
}

func (p *funcWrapper) Interface() any {
	return p.MakeFunc().Interface()
}

func (p *funcWrapper) MakeFunc() reflect.Value {
	inTypes := p.getParamterTypes(p.inputs)
	outTypes := p.getParamterTypes(p.outputs)
	funcType := reflect.FuncOf(inTypes, outTypes, p.variadic)

	target := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
		return p.do(args)
	})

	return target
}

func (p *funcWrapper) do(ins1 []reflect.Value) []reflect.Value {
	ins0 := make([]reflect.Value, p.typ.NumIn())
	for i, v := range p.inputs {
		if v.pos < 0 {
			continue
		}

		ins0[v.pos] = ins1[i]
	}

	if p.variadic {
		// log.Println()
		// for i, v := range ins1 {
		// 	log.Printf(" *** %d: %v \n", i, v.Interface())
		// }

		ins0 = ins0[:p.typ.NumIn()-1]
		last := ins1[len(ins1)-1]
		for i, n := 0, last.Len(); i < n; i++ {
			ins0 = append(ins0, last.Index(i))
		}
	}

	outs0 := p.value.Call(ins0)
	outs1 := make([]reflect.Value, len(p.outputs))
	for i, v := range p.outputs {
		if v.pos < 0 {
			outs1[i] = v.value
			continue
		}

		outs1[i] = outs0[v.pos]
	}

	return outs1
}

func (p *funcWrapper) reflectOf(val any, typ ...reflect.Type) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(val)

	if len(typ) == 0 {
		t := v.Type()
		return v, t
	}

	t := typ[0]
	if !v.IsValid() {
		v = reflect.Zero(t)

	} else if !v.Type().Implements(t) {
		panic(fmt.Sprintf("`%s` should implement `%s`", v.Type().String(), t.String()))
	}

	return v, t
}

func (p *funcWrapper) SwapInput(a, b int) *funcWrapper {
	n := len(p.inputs)
	if p.variadic && (a >= n-1 || b >= n-1) {
		panic(fmt.Sprintf("variadic func must not swap the last input parameter"))
	}

	p.inputs[a], p.inputs[b] = p.inputs[b], p.inputs[a]

	return p
}

func (p *funcWrapper) SwapOutput(a, b int) *funcWrapper {
	p.outputs[a], p.outputs[b] = p.outputs[b], p.outputs[a]
	return p
}

func (p *funcWrapper) AppendInput(val any, typ ...reflect.Type) *funcWrapper {
	if p.variadic {
		panic(fmt.Sprintf("variadic func must not append input parameter"))
	}

	v, t := p.reflectOf(val, typ...)

	p.inputs = append(p.inputs, &funcWrapperParameter{
		value: v,
		typ:   t,
		pos:   -1,
	})
	return p
}

func (p *funcWrapper) AppendOutput(val any, typ ...reflect.Type) *funcWrapper {
	v, t := p.reflectOf(val, typ...)

	p.outputs = append(p.outputs, &funcWrapperParameter{
		value: v,
		typ:   t,
		pos:   -1,
	})
	return p
}

func (p *funcWrapper) InsertInput(pos int, val any, typ ...reflect.Type) *funcWrapper {
	n := len(p.inputs)
	if pos > n {
		panic(fmt.Sprintf("the input parameter position for add is out of range. pos:%d, len:%d", pos, n))
	}

	if pos == n && p.variadic {
		panic(fmt.Sprintf("variadic func must not add input parameter after last"))
	}

	v, t := p.reflectOf(val, typ...)

	p.inputs = append(p.inputs[:pos+1], p.inputs[pos:]...)
	p.inputs[pos] = &funcWrapperParameter{
		value: v,
		typ:   t,
		pos:   -1,
	}
	return p
}

func (p *funcWrapper) InsertOutput(pos int, val any, typ ...reflect.Type) *funcWrapper {
	n := len(p.outputs)
	if pos > n {
		panic(fmt.Sprintf("the output parameter position for add is out of range. pos:%d, len:%d", pos, n))
	}

	v, t := p.reflectOf(val, typ...)

	p.outputs = append(p.outputs[:pos+1], p.outputs[pos:]...)
	p.outputs[pos] = &funcWrapperParameter{
		value: v,
		typ:   t,
		pos:   -1,
	}
	return p
}
