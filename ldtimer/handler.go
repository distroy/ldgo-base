/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldlog"
)

type (
	inConvType  = func(context.Context, *Task) (reflect.Value, error)
	outConvType = func(context.Context, []reflect.Value) error
)

var (
	typeOfTask    = reflect.TypeFor[*Task]()
	typeOfContext = reflect.TypeFor[context.Context]()
	typeOfError   = reflect.TypeFor[error]()
	typeOfString  = reflect.TypeFor[string]()
)

func NewHandler(f any) *Handler {
	handlerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	h := &Handler{
		Service: nil,
		Func:    f,
		Name:    handlerName,
		Value:   reflect.ValueOf(f),
	}
	h.init()
	return h
}

func NewHandlerByMethod(service any, method string) *Handler {
	v := reflect.ValueOf(service)
	m := v.MethodByName(method)
	if !m.IsValid() {
		panic(fmt.Errorf("cannot find the method `%s` from type `%s`", method, v.Type().String()))
	}
	h := &Handler{
		Service: nil,
		Func:    m.Interface(),
		Name:    method,
		Value:   m,
	}
	h.init()
	return h
}

type Handler struct {
	Service any
	Func    any
	Name    string
	Value   reflect.Value

	Type    reflect.Type
	InConvs []inConvType
	OutConv outConvType
}

func (h *Handler) init() {
	// h.Value = reflect.ValueOf(h.Func)
	h.Type = h.Value.Type()

	h.InConvs = h.getAllInConvs(h.Type)
	h.OutConv = h.getOutConv(h.Type)
}

func (h *Handler) getAllInConvs(funcType reflect.Type) []inConvType {
	offset := 0
	if h.Service != nil {
		receiver := reflect.TypeOf(h.Service)
		offset = 1
		if funcType.NumIn() == 0 {
			panic(fmt.Sprintf("[ldtimer] the receiver type should be `%s`. method:%s, func:%s",
				receiver.String(), h.Name, funcType.String()))

		} else if inType := funcType.In(0); inType != receiver {
			panic(fmt.Sprintf("[ldtimer] the receiver type should be `%s`. method:%s, actual:%s",
				receiver.String(), h.Name, inType.String()))
		}
	}

	inCount := funcType.NumIn() - offset
	switch inCount {
	case 0, 1, 2:
		break

	default:
		panic(fmt.Sprintf("[ldtimer] %s input parameter count should be 0 or 1 or 2", h.Name))
	}

	if inCount >= 1 {
		ctxType := funcType.In(0 + offset)
		if ctxType != typeOfContext && !h.isType(typeOfContext, ctxType) {
			panic(fmt.Sprintf("[ldtimer] %s input parameter type should be `context.Context`", h.Name))
		}
	}

	// if inCount >= 2 {
	// 	paramType := funcType.In(1 + offset)
	// }

	inConvs := make([]inConvType, 0, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		inType := funcType.In(i)
		inConvs = append(inConvs, h.getOneInConv(inType))
	}

	return inConvs
}

func (h *Handler) getOneInConv(t reflect.Type) inConvType {
	if h.Service != nil && reflect.TypeOf(h.Service) == t {
		return func(ctx context.Context, task *Task) (reflect.Value, error) {
			return reflect.ValueOf(h.Service), nil
		}
	}

	if t == typeOfContext || h.isType(typeOfContext, t) {
		return func(ctx context.Context, task *Task) (reflect.Value, error) {
			return reflect.ValueOf(ctx), nil
		}
	}

	if t == typeOfTask {
		return func(ctx context.Context, task *Task) (reflect.Value, error) {
			return reflect.ValueOf(task), nil
		}
	}

	if t == typeOfString {
		return func(ctx ldctx.Context, task *Task) (reflect.Value, error) {
			return reflect.ValueOf(task.Info.GetParams()), nil
		}
	}

	if t.Kind() != reflect.Ptr {
		return func(ctx context.Context, task *Task) (reflect.Value, error) {
			v := reflect.New(t)
			err := decode(ctx, task, v.Interface())
			return v.Elem(), err
		}
	}

	return func(ctx context.Context, task *Task) (reflect.Value, error) {
		v := reflect.New(t.Elem())
		err := decode(ctx, task, v.Interface())
		return v, err
	}
}

func (h *Handler) getOutConv(funcType reflect.Type) outConvType {
	switch funcType.NumOut() {
	case 0:
		return func(ctx context.Context, v []reflect.Value) error {
			return nil
		}

	case 1:
		break

	default:
		panic(fmt.Sprintf("[ldtimer] %s output parameter count should be 0 or 1", h.Name))
	}

	outType := funcType.Out(0)
	errType := outType
	if errType != typeOfError && !h.isType(typeOfError, errType) {
		panic(fmt.Sprintf("[ldtimer] %s output parameter type should be `error`", h.Name))
	}

	return func(c context.Context, outs []reflect.Value) error {
		out0 := outs[0].Interface()
		if out0 == nil {
			return nil
		}

		err := out0.(error)
		if !lderr.IsSuccess(err) {
			return err
		}
		return nil
	}
}

func (h *Handler) isType(child, parent reflect.Type) bool {
	if child == parent {
		return true
	}
	if parent.Kind() == reflect.Interface && child.Implements(parent) {
		return true
	}
	return false
}

func (h *Handler) Do(task *Task) (_err error) {
	c := newContext(task)
	ldctx.LogI(c, "[ldtimer] do timer handler begin", ldlog.String("adaptor", task.Adaptor.Name()),
		ldlog.String("name", h.Name), ldlog.Reflect("task", task))
	defer func() {
		if e := recover(); e != nil {
			if ee, _ := e.(error); ee != nil {
				_err = ee
			} else {
				_err = fmt.Errorf("%v", ee)
			}
			ldctx.LogE(c, "[ldtimer] do timer handler panic", ldlog.String("adaptor", task.Adaptor.Name()),
				ldlog.String("name", h.Name), ldlog.Error(_err), ldlog.Stack("stack"))
		}
		ldctx.LogI(c, "[ldtimer] do timer handler end", ldlog.String("adaptor", task.Adaptor.Name()),
			ldlog.String("name", h.Name), ldlog.Error(_err))
	}()

	ins := make([]reflect.Value, 0, len(h.InConvs))
	for _, inConv := range h.InConvs {
		in, err := inConv(c, task)
		if err != nil {
			return err
		}
		ins = append(ins, in)
	}

	outs := h.Value.Call(ins)
	err := h.OutConv(c, outs)

	if err != nil {
		return err
	}
	return nil
}
