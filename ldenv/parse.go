/*
 * Copyright (C) distroy
 */

package ldenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/ldstr"
	"github.com/distroy/ldgo-base/ldtags"
)

const (
	tag = "ldenv"
)

func Parse(p any) {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("input parameter should be pointer to struct. type:%s", v.Type().String()))
	}

	v = v.Elem()
	t := v.Type()

	for i, l := 0, v.NumField(); i < l; i++ {
		fv := v.Field(i)
		sf := t.Field(i)
		parseEnvField(fv, sf)
	}
}

func parseEnvField(fv reflect.Value, sf reflect.StructField) {
	tagStr := sf.Tag.Get(tag)
	tags := ldtags.Parse(tagStr)

	name := tags.Get("name")
	def := tags.Get("default")
	if name == "" {
		name = ldstr.ToSnakeCase(sf.Name)
		name = strings.ToUpper(name)
	}

	val := os.Getenv(name)
	if val == "" {
		val = def
	}

	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv.Set(reflect.New(sf.Type.Elem()))
		}
		fv = fv.Elem()
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fv.SetInt(ldconv.AsInt64(val))
		// v, _ := strconv.ParseInt(val, 10, 64)
		// fv.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fv.SetUint(ldconv.AsUint64(val))
		// v, _ := strconv.ParseUint(val, 10, 64)
		// fv.SetUint(v)

	case reflect.String:
		fv.SetString(val)

	case reflect.Float32, reflect.Float64:
		fv.SetFloat(ldconv.AsFloat64(val))
		// v, _ := strconv.ParseFloat(val, 64)
		// fv.SetFloat(v)
	}
}
