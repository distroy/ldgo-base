/*
 * Copyright (C) distroy
 */

package ldrcfg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/distroy/ldgo-base/3rd/yaml"
	"github.com/distroy/ldgo-base/ldatomic"
	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/ldlog"
	"github.com/distroy/ldgo-base/ldstr"
	"github.com/distroy/ldgo-base/ldtags"
)

var (
	_ ldatomic.StoreLoader[int] = (*Event[int])(nil)

	typeOfCtx = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type Parser interface {
	Parse(ctx context.Context, v []byte) error
}

type Event[T any] struct {
	value  ldatomic.Any[T]
	events []func(c context.Context, v T)
}

func (e *Event[T]) Store(v T)                              { e.value.Store(v) }
func (e *Event[T]) Load() T                                { return e.value.Load() }
func (e *Event[T]) Value() *ldatomic.Any[T]                { return &e.value }
func (e *Event[T]) Events() []func(c context.Context, v T) { return e.events }

func (e *Event[T]) OnChange(c context.Context, onChanges ...func(c context.Context, v T)) *Event[T] {
	v := e.Load()
	for _, onChange := range onChanges {
		onChange(c, v)
	}
	e.events = append(e.events, onChanges...)
	return e
}

func (c *Client) Register(ns string, cfg any) {
	val := reflect.ValueOf(cfg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("[ldrcfg] config should be pointer of struct. %s", val.Type()))
	}

	ctx := c.getContextWithoutSequence()
	typ := val.Type()
	for i := range val.NumField() {
		fv := val.Field(i)
		sf := typ.Field(i)

		if !sf.IsExported() {
			continue
		}

		cf := parseConfigField(fv, sf)
		if cf == nil || cf.Key == "" || cf.Callback == nil {
			continue
		}

		setConfigFieldDefaultData(cf)
		c.RegisterKey(ns, cf.Key, cf.Callback)
		ldctx.LogI(ctx, "[ldrcfg] register config field succ", ldlog.String("ns", ns),
			ldlog.String("key", cf.Key), ldlog.String("type", cf.Type))
	}
}

type configField struct {
	Field       reflect.Value
	StructField reflect.StructField
	DataType    reflect.Type
	Key         string
	Type        string
	Callback    func(context.Context, *KeyChangeEvent)
	Store       reflect.Value
	Load        reflect.Value
	Events      reflect.Value
}

func parseConfigField(fv reflect.Value, sf reflect.StructField) *configField {
	tagStr := sf.Tag.Get("ldrcfg")
	if tagStr == "-" {
		return nil
	}

	tags := ldtags.Parse(tagStr)
	key := tags.Get("key")
	typ := tags.Get("type") // yaml/json/parser/string

	if key == "" {
		key = ldstr.ToSnakeCase(sf.Name)
	}

	if fv.Kind() == reflect.Ptr && fv.Pointer() == 0 {
		fv.Set(reflect.New(fv.Type().Elem()))

	} else {
		fv = fv.Addr()
	}

	store := fv.MethodByName("Store")
	if !store.IsValid() || store.Type().NumIn() != 1 {
		return nil
	}
	dataType := store.Type().In(0)

	load := fv.MethodByName("Load")
	if !checkConfigFieldLoadType(load, dataType) {
		return nil
	}

	events := fv.MethodByName("Events")
	if !checkConfigFieldEventType(events, dataType) {
		events = reflect.Value{}
	}

	cf := &configField{
		Field:       fv,
		StructField: sf,
		DataType:    dataType,
		Key:         key,
		Type:        typ,
		Store:       store,
		Load:        load,
		Events:      events,
	}
	cf.Callback = getConfigFieldCallback(cf)
	return cf
}

func checkConfigFieldLoadType(load reflect.Value, dateType reflect.Type) bool {
	lt := load.Type()
	if lt.NumOut() != 1 && lt.NumIn() != 0 {
		return false
	}

	if o0 := lt.Out(0); o0 != dateType {
		return false
	}

	return true
}

func checkConfigFieldEventType(events reflect.Value, dataType reflect.Type) bool {
	et := events.Type()
	if et.NumOut() != 1 && et.NumIn() != 0 {
		return false
	}

	out := et.Out(0)
	if out.Kind() != reflect.Slice {
		return false
	}

	out = out.Elem()
	if out.Kind() != reflect.Func || out.NumIn() != 2 {
		return false
	}

	if i0 := out.In(0); i0 != typeOfCtx && !i0.Implements(typeOfCtx) {
		return false
	}

	if i1 := out.In(1); i1 != dataType {
		return false
	}

	return true
}

func getConfigFieldCallback(cf *configField) func(context.Context, *KeyChangeEvent) {
	p := reflect.New(cf.DataType)
	if cf.DataType.Kind() == reflect.Ptr {
		p = p.Elem()
	}

	if cf.Type == "" {
		_, isParser := p.Interface().(Parser)
		dt := cf.DataType
		if isParser {
			cf.Type = "parser"

		} else if dt.Kind() == reflect.String || (dt.Kind() == reflect.Ptr && dt.Elem().Kind() == reflect.String) {
			cf.Type = "string"
		}
	}

	switch cf.Type {
	case "string":
		return getConfigFieldCallbackByDecode(cf, func(c context.Context, in []byte, out any) error {
			*out.(*string) = ldconv.BytesToStrUnsafe(in)
			return nil
		})

	case "json":
		return getConfigFieldCallbackByDecode(cf, func(c context.Context, in []byte, out any) error {
			return json.Unmarshal(in, out)
		})

	case "yaml":
		return getConfigFieldCallbackByDecode(cf, func(c context.Context, in []byte, out any) error {
			return yaml.Unmarshal(in, out)
		})

	case "parser":
		_, isParser := p.Interface().(Parser)
		if !isParser {
			panic(fmt.Errorf("[ldrcfg] data type should be parser. field:%s, data type:%s",
				cf.StructField.Name, cf.DataType))
		}

		return getConfigFieldCallbackByDecode(cf, func(c context.Context, in []byte, out any) error {
			v := out.(Parser)
			return v.Parse(c, in)
		})
	}

	return getConfigFieldCallbackByDecode(cf, func(c context.Context, in []byte, out any) error {
		in = bytes.TrimSpace(in)
		if len(in) == 0 || in[0] == '[' || in[0] == '{' {
			return json.Unmarshal(in, out)
		}
		return yaml.Unmarshal(in, out)
	})
}

func getConfigFieldCallbackByDecode(cf *configField, decode func(c context.Context, in []byte, out any) error) func(context.Context, *KeyChangeEvent) {
	return func(ctx context.Context, ev *KeyChangeEvent) {
		var (
			ns  = ev.Namespace
			key = ev.Key
		)

		var p, v reflect.Value
		if cf.DataType.Kind() == reflect.Ptr {
			p = reflect.New(cf.DataType.Elem())
			v = p

		} else {
			p = reflect.New(cf.DataType)
			v = p.Elem()
		}

		raw := ldconv.StrToBytesUnsafe(ev.Change.NewValue)
		if err := decode(ctx, raw, p.Interface()); err != nil {
			ldctx.LogE(ctx, "[config center] parse object fail", ldlog.String("method", cf.Type),
				ldlog.String("ns", ns), ldlog.String("key", key), ldlog.String("str", ev.Change.NewValue),
				ldlog.Error(err))
			return
		}

		x := v.Interface()
		ldctx.LogI(ctx, "[config center] update parser object succ", ldlog.String("method", cf.Type),
			ldlog.String("ns", ns), ldlog.String("key", key), ldlog.Reflect("value", x))

		cf.Store.Call([]reflect.Value{v})

		if !cf.Events.IsValid() {
			return
		}

		ctxVal := reflect.ValueOf(ctx)

		outs := cf.Events.Call(nil)[0]
		for i := range outs.Len() {
			f := outs.Index(i)
			f.Call([]reflect.Value{ctxVal, v})
		}
		// e.trigger(ctx, x)
	}
}

func setConfigFieldDefaultData(cf *configField) {
	dt := cf.DataType
	if dt.Kind() != reflect.Ptr {
		return
	}

	out := cf.Load.Call(nil)[0]
	if out.IsValid() && !out.IsNil() {
		return
	}

	v := reflect.New(dt.Elem())
	cf.Store.Call([]reflect.Value{v})
}
