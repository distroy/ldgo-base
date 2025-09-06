/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"encoding/json"
	"math/big"
)

func AsByte(val interface{}, def ...byte) byte { return AsUint8(val, def...) }

func AsBool(val interface{}, def ...bool) bool {
	v, err := ToBool(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func AsInt(val interface{}, def ...int) int {
	v, err := ToInt(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsInt8(val interface{}, def ...int8) int8 {
	v, err := ToInt8(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsInt16(val interface{}, def ...int16) int16 {
	v, err := ToInt16(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsInt32(val interface{}, def ...int32) int32 {
	v, err := ToInt32(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsInt64(val interface{}, def ...int64) int64 {
	v, err := ToInt64(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func AsUint(val interface{}, def ...uint) uint {
	v, err := ToUint(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsUint8(val interface{}, def ...uint8) uint8 {
	v, err := ToUint8(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsUint16(val interface{}, def ...uint16) uint16 {
	v, err := ToUint16(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsUint32(val interface{}, def ...uint32) uint32 {
	v, err := ToUint32(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsUint64(val interface{}, def ...uint64) uint64 {
	v, err := ToUint64(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsUintptr(val interface{}, def ...uintptr) uintptr {
	v, err := ToUintptr(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func AsFloat32(val interface{}, def ...float32) float32 {
	v, err := ToFloat32(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func AsFloat64(val interface{}, def ...float64) float64 {
	v, err := ToFloat64(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func AsString(val interface{}, def ...string) string {
	v, err := ToString(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func asBigFloat(val interface{}, def ...*big.Float) *big.Float {
	v, err := toBigFloat(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return newBigFloatZero()
}

func asDecimal(val interface{}, def ...decimalNumber) decimalNumber {
	v, err := toDecimal(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return newDecimalZero()
}

func AsJsonNumber(val interface{}, def ...json.Number) json.Number {
	v, err := ToJsonNumber(val)
	if err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return "0"
}
