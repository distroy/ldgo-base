/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"container/list"
	"fmt"
	"reflect"
)

type Patches interface {
	Apply(hooks ...Hook)
	Applys(hooks []Hook)

	ApplyFunc(target, double any)
	ApplyMethod(target reflect.Type, methodName string, double any)
	ApplyGlobalVar(target, double any)

	Reset()
}

func NewPatches() Patches {
	return &patches{}
}

type patches struct {
	funcs     list.List
	variables list.List
}

func (p *patches) Reset() {
	for i := p.funcs.Back(); i != nil; i = i.Prev() {
		v := i.Value.(patchFunc)
		resetFunc(v)
	}
	p.funcs.Init()

	for i := p.variables.Back(); i != nil; i = i.Prev() {
		v := i.Value.(patchVariable)
		resetVariable(v)
	}
	p.variables.Init()

}

func (p *patches) Apply(hooks ...Hook) {
	p.Applys(hooks)
}

func (p *patches) Applys(hooks []Hook) {
	if len(hooks) == 0 {
		return
	}

	for _, h := range hooks {
		h.hook(p)
	}
}

func (p *patches) coreApplyVariable(target, double reflect.Value) {
	if target.Type().Kind() != reflect.Ptr {
		panic("target is not a pointer")
	}

	value := applyVariable(target, double)
	p.variables.PushBack(value)
}

func (p *patches) coreApplyFunc(target, double reflect.Value) {
	p.checkFunc(target, double)

	value := applyFunc(target, double)
	p.funcs.PushBack(value)
}

func (p *patches) checkFunc(target, double reflect.Value) {
	if target.Kind() != reflect.Func {
		panic("target is not a func")
	}

	if double.Kind() != reflect.Func {
		panic("double is not a func")
	}

	if target.Type() != double.Type() {
		panic(fmt.Sprintf("target type(%s) and double type(%s) are different", target.Type(), double.Type()))
	}
}

func (p *patches) ApplyFunc(target, double any) {
	p.Apply(FuncHook{
		Target: target,
		Double: double,
	})
}

func (p *patches) ApplyMethod(target reflect.Type, method string, double any) {
	p.Apply(MethodHook{
		Target: target,
		Method: method,
		Double: double,
	})
}

func (p *patches) ApplyGlobalVar(target, double any) {
	p.Apply(VariableHook{
		Target: target,
		Double: double,
	})
}
