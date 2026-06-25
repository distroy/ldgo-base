/*
 * Copyright (C) distroy
 */

package ldenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/distroy/ldgo-base/internal/time_"
	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/ldstr"
	"github.com/distroy/ldgo-base/ldtags"
	"github.com/distroy/ldgo-base/ldtime"
)

const (
	tag = "ldenv"
)

var (
	errInvalidValueType = fmt.Errorf(`invalid value type`)
)

func Parse(p any) error {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input parameter should be pointer to struct. type:%s", v.Type().String())
	}

	v = v.Elem()
	t := v.Type()

	var resErr error
	for i, l := 0, v.NumField(); i < l; i++ {
		fv := v.Field(i)
		sf := t.Field(i)
		err := parseEnvField(fv, sf)
		if err != nil && resErr == nil {
			resErr = err
		}
	}
	return resErr
}

func parseEnvField(fv reflect.Value, sf reflect.StructField) error {
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

	vv := fv
	if vv.Kind() == reflect.Ptr {
		if vv.IsNil() {
			vv.Set(reflect.New(sf.Type.Elem()))
		}
		vv = vv.Elem()
	}
	if err := parseEnvFieldValue(vv, val, tags); err != nil {
		if err == errInvalidValueType {
			err = fmt.Errorf(`invalid value type(%s)`, fv.Type().String())
		}
		return fmt.Errorf(`parse env fail. env key:%s, env value:%s, err:%v`,
			name, val, err)
	}
	return nil
}

func parseEnvFieldValue(fv reflect.Value, val string, tags ldtags.Tags) error {
	layout := tags.Get("layout")
	lower := tags.Has("lower")
	upper := tags.Has("upper")
	type Parser interface {
		Parse(raw []byte) error
	}

	switch v := fv.Addr().Interface().(type) {
	case *time.Duration:
		x, err := time_.DurationUnmarshalYaml(ldconv.StrToBytesUnsafe(val))
		*v = x
		return err

	case *ldtime.Duration:
		x, err := time_.DurationUnmarshalYaml(ldconv.StrToBytesUnsafe(val))
		*v = ldtime.Duration(x)
		return err

	case *time.Time:
		x, err := time.Parse(if_func(layout != "", layout, time_.TimeLayout), val)
		*v = x
		return err

	case interface{ Parse(raw []byte) error }:
		return v.Parse(ldconv.StrToBytesUnsafe(val))
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := ldconv.ToInt64(val)
		fv.SetInt(x)
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x, err := ldconv.ToUint64(val)
		fv.SetUint(x)
		return err

	case reflect.String:
		switch {
		case lower:
			val = strings.ToLower(val)
		case upper:
			val = strings.ToUpper(val)
		}
		fv.SetString(val)
		return nil

	case reflect.Bool:
		x, err := ldconv.ToBool(val)
		fv.SetBool(x)
		return err

	case reflect.Float32, reflect.Float64:
		x, err := ldconv.ToFloat64(val)
		fv.SetFloat(x)
		return err
	}
	return errInvalidValueType
}

func if_func[T any](cond bool, trueRes T, falseRes ...T) T {
	if cond {
		return trueRes
	}
	if len(falseRes) > 0 {
		return falseRes[0]
	}
	var x T
	return x
}
