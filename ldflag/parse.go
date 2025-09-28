/*
 * Copyright (C) distroy
 */

package ldflag

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/distroy/ldgo-base/ldstr"
)

const (
	tagName = "ldflag"
)

var (
	defaultFlagSet *FlagSet
	defaultOptions = ([]func(s *FlagSet))(nil)
)

func Default() *FlagSet {
	s := defaultFlagSet
	if s == nil {
		s = newDefaultFlagSet()
	}
	return s
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func newDefaultFlagSet() *FlagSet {
	s := NewFlagSet()

	for _, opt := range defaultOptions {
		opt(s)
	}

	defaultFlagSet = s
	return s
}

func MustParse(v any, args ...[]string) {
	s := newDefaultFlagSet()
	s.Model(v)
	s.MustParse(args...)
}

func Parse(v any, args ...[]string) error {
	s := newDefaultFlagSet()
	s.Model(v)
	return s.Parse(args...)
}

func PrintUsage()            { defaultFlagSet.PrintUsage() }
func WriteUsage(w io.Writer) { defaultFlagSet.WriteUsage(w) }

func parseFlagName(f reflect.StructField) string {
	name := f.Name
	name = ldstr.ToSnakeCase(name, '-')
	name = strings.ToLower(name)
	return name
}

func packMeta(meta string) string {
	if !strings.HasPrefix(meta, "<") && !strings.HasPrefix(meta, ">") {
		meta = fmt.Sprintf("<%s>", meta)
	}
	return meta
}

func unquoteUsage(f *Flag) (meta string, usage string) {
	usage = f.Usage
	meta = f.Meta
	if meta != "" {
		meta = packMeta(meta)
		return meta, usage
	}

	// // Look for a back-quoted name, but avoid the strings package.
	// for i := 0; i < len(usage); i++ {
	// 	if usage[i] == '`' {
	// 		for j := i + 1; j < len(usage); j++ {
	// 			if usage[j] == '`' {
	// 				meta = usage[i+1 : j]
	// 				meta = packMeta(meta)
	//
	// 				usage = usage[:i] + meta + usage[j+1:]
	// 				return meta, usage
	// 			}
	// 		}
	// 		break // Only one back quote; use type name.
	// 	}
	// }

	// No explicit name, so use type if we can find one.
	meta = "<value>"
	switch v := f.Value.(type) {
	case valueWithMeta:
		meta = v.Meta()
	}
	return meta, usage
}

func getAddrValue(v reflect.Value) Value {
	if !v.CanAddr() {
		return nil
	}

	vv, _ := v.Addr().Interface().(Value)
	return vv
}

func isFlagDefaultZero(f *Flag) bool {
	value := f.Value
	defaultValue := f.Default

	if defaultValue == "" {
		return true
	}

	if f.val.Kind() == reflect.Slice {
		return defaultValue == "null" || defaultValue == "[]"
	}

	val := f.val
	if getAddrValue(val) != nil {
		typ := val.Type()
		v := reflect.New(typ)
		switch typ.Kind() {
		case reflect.Interface:
			v.Elem().Set(reflect.New(typ.Elem()))
		}

		z, _ := v.Interface().(Value)
		return defaultValue == z.String()
	}

	if val.Kind() == reflect.Ptr {
		typ := val.Type()
		v := reflect.New(typ).Elem()
		v.Set(reflect.New(typ.Elem()))
		// v := reflect.New(typ.Elem())

		z, _ := v.Interface().(Value)
		if z == nil {
			z = newValFuncMap[typ](v)
		}
		return defaultValue == z.String()
	}

	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return defaultValue == z.Interface().(Value).String()
}
