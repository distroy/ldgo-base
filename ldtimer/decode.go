/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/distroy/ldgo-base/3rd/yaml"
	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldlog"
)

func decode(c context.Context, task *Task, param any) error {
	paramStr := task.Info.GetParams()
	reqVal := reflect.ValueOf(param)
	if reqVal.Kind() != reflect.Ptr || reqVal.Elem().Kind() != reflect.Struct {
		ldctx.LogE(c, "[ldtimer] input parameter type must be pointer to struct", ldlog.Stringer("type", reqVal.Kind()))
		return lderr.ErrInvalidRequestType
	}

	bytes := ldconv.StrToBytesUnsafe(paramStr)
	if len(bytes) == 0 {
		ldctx.LogI(c, "[ldtimer] decode parameter succ", ldlog.Reflect("param", param))
		return nil
	}

	reqVal = reqVal.Elem()
	reqType := getParamType(reqVal.Type())

	if err := reqType.Func(bytes, param); err != nil {
		ldctx.LogE(c, "[ldtimer] decode parameter fail", ldlog.String("tag", reqType.Tag),
			ldlog.String("param", paramStr), ldlog.Error(err))
		return lderr.ErrParseRequest
	}

	ldctx.LogI(c, "[ldtimer] decode parameter succ", ldlog.Reflect("param", param))
	return nil
}

var paramTypes = &sync.Map{}
var paramParsers = []paramParser{
	{Tag: "json", Func: json.Unmarshal},
	{Tag: "yaml", Func: yaml.Unmarshal},
}

type paramParser struct {
	Tag  string
	Func func(bytes []byte, o any) error
}

type paramType struct {
	Type reflect.Type
	Tag  string
	Func func(bytes []byte, o any) error
}

func getParamType(t reflect.Type) *paramType {
	v, _ := paramTypes.Load(t)
	p, _ := v.(*paramType)
	if p != nil {
		return p
	}

	parser := paramParser{
		Tag:  "json",
		Func: json.Unmarshal,
	}
	for _, v := range paramParsers {
		if !hasTags(t, v.Tag) {
			continue
		}
		parser = v
		break
	}

	p = &paramType{
		Type: t,
		Tag:  parser.Tag,
		Func: parser.Func,
	}

	paramTypes.Store(t, p)
	return p
}

func hasTags(t reflect.Type, tagName string) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagStr, _ := field.Tag.Lookup(tagName)
		if tagStr != "" && tagStr != "-" {
			return false
		}
	}
	return false
}
