/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"fmt"
	"reflect"
	"runtime"
)

type Values = []any
type OutputCell struct {
	Values Values
	Times  int
}

type ResultCell struct {
	Inputs  Values
	Outputs Values
	Times   int
}

func BindInput(pos int, val any) any {
	if pos < 0 {
		return val
	}
	return input{
		Position: pos,
		Value:    reflect.ValueOf(val),
	}
}

type input struct {
	Position int
	Value    reflect.Value
	Type     reflect.Type
	DataType reflect.Type
}

type wrapOutputs struct {
	FuncType reflect.Type
	Results  []ResultCell
	Length   int
	Position int
	Times    int
}

func newWrap(funcType reflect.Type, outputs []ResultCell) *wrapOutputs {
	return &wrapOutputs{
		FuncType: funcType,
		Results:  outputs,
		Length:   len(outputs),
		Position: 0,
		Times:    0,
	}
}

func (w *wrapOutputs) getResultCell() *ResultCell {
	if w.Position >= w.Length {
		return &w.Results[w.Length-1]
	}

	cell := w.Results[w.Position]
	w.Times++
	if w.Times >= cell.Times {
		w.Position++
		w.Times = 0
	}
	return &cell
}

func (w *wrapOutputs) fillInputs(res *ResultCell, ins []reflect.Value) {
	if w.FuncType.IsVariadic() {
		l := len(ins)
		last := ins[l-1]
		tmp := make([]reflect.Value, 0, l-1+last.Len())
		tmp = append(tmp, ins[:l-1]...)
		for i := 0; i < last.Len(); i++ {
			tmp = append(tmp, last.Index(i))
		}
		ins = tmp
	}

	// types := make([]string, 0, len(ins))
	// for _, v := range ins {
	// 	types = append(types, v.Type().String())
	// }
	// log.Printf("input types: %s", strings.Join(types, ", "))

	for _, i := range res.Inputs {
		switch v := i.(type) {
		case input:
			w.fillOneToInputs(&v, ins)
		case *input:
			w.fillOneToInputs(v, ins)
		case *reflect.Value:
			w.fillOneToInputs(&input{Position: -1, Value: *v}, ins)
		case reflect.Value:
			w.fillOneToInputs(&input{Position: -1, Value: v}, ins)
		default:
			w.fillOneToInputs(&input{Position: -1, Value: reflect.ValueOf(i)}, ins)
		}
	}
}

func (w *wrapOutputs) fillOneToInput(res *input, in reflect.Value) error {
	resV := res.Value
	resT := res.Type

	inV := in
	inT := inV.Type()

	if resT.Kind() == reflect.Func {
		valT := res.DataType
		if inT.ConvertibleTo(valT) {
			resV.Call([]reflect.Value{inV})
			return nil
		}

		return fmt.Errorf("the parameter of position %d must be convertible to %s",
			res.Position, valT.String())
	}

	if inT.Kind() != reflect.Ptr {
		return fmt.Errorf("the parameter of position %d must be pointer", res.Position)
	}

	if resT.Kind() == reflect.Ptr && resT.Elem().AssignableTo(inT.Elem()) {
		inV.Elem().Set(resV.Elem())
		return nil
	}
	if resT.AssignableTo(inT.Elem()) {
		inV.Elem().Set(resV)
		return nil
	}

	return fmt.Errorf("the parameter of position %d must be %s", res.Position, inT.String())
}

func (w *wrapOutputs) fillOneToInputs(res *input, ins []reflect.Value) {
	if res.Position >= len(ins) {
		panic(fmt.Sprintf("func has only %d parameters, can not get the parameter of position %d",
			len(ins), res.Position))
	}

	if res.Type == nil {
		res.Type = res.Value.Type()
	}
	if res.DataType == nil {
		if res.Type.Kind() != reflect.Func {
			res.DataType = res.Type

		} else if res.Type.NumIn() == 1 {
			res.DataType = res.Type.In(0)

		} else {
			funcName := runtime.FuncForPC(res.Value.Pointer()).Name()
			panic(fmt.Errorf("func must have 1 parameter. func:%s", funcName))
		}
	}

	if res.Position >= 0 {
		in := ins[res.Position]
		inV := reflect.ValueOf(in.Interface())
		err := w.fillOneToInput(res, inV)
		if err != nil {
			panic(err)
		}
		return
	}

	for _, in := range ins {
		inV := reflect.ValueOf(in.Interface())
		err := w.fillOneToInput(res, inV)
		if err != nil {
			continue
		}
		return
	}

	if res.Type.Kind() == reflect.Func {
		panic(fmt.Sprintf("cannot find the parameter can be convertible to %s", res.DataType.String()))
	} else {
		panic(fmt.Sprintf("cannot find the parameter is %s", res.Type.String()))
	}
}

func (w *wrapOutputs) getOutputs(res *ResultCell) []reflect.Value {
	funcType := w.FuncType
	resultValues := make([]reflect.Value, 0, len(res.Outputs))

	for i, r := range res.Outputs {
		var resultValue reflect.Value
		if r == nil {
			resultValue = reflect.Zero(funcType.Out(i))

		} else {
			vf := reflect.New(funcType.Out(i)).Elem()
			vr := reflect.ValueOf(r)
			if vf.Kind() != vr.Kind() && vr.Type().ConvertibleTo(vf.Type()) {
				vr = vr.Convert(vf.Type())
			}
			vf.Set(vr)
			resultValue = vf
		}
		resultValues = append(resultValues, resultValue)
	}

	return resultValues
}

func (w *wrapOutputs) Do(ins []reflect.Value) []reflect.Value {
	res := w.getResultCell()
	w.fillInputs(res, ins)
	return w.getOutputs(res)
}

func (w *wrapOutputs) MakeFunc() reflect.Value {
	funcType := w.FuncType
	outputs := w.Results
	if funcType.NumOut() != len(outputs[0].Outputs) {
		panic(fmt.Sprintf("func type has %v return values, but only %v values provided as double, func type: %s",
			funcType.NumOut(), len(outputs[0].Outputs), funcType))
	}

	return reflect.MakeFunc(funcType, func(ins []reflect.Value) []reflect.Value {
		return w.Do(ins)
	})
}
