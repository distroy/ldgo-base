/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/distroy/ldgo-base/internal/time_"
)

const (
	timeFormat = time_.TimeLayout
)

func ToByte(v any) (byte, error) { return ToUint8(v) }

func ToBool(v any) (bool, error) {
	switch vv := v.(type) {
	case bool:
		return vv, nil

	case int:
		return vv != 0, nil
	case int8:
		return vv != 0, nil
	case int16:
		return vv != 0, nil
	case int32:
		return vv != 0, nil
	case int64:
		return vv != 0, nil

	case uint:
		return vv != 0, nil
	case uint8:
		return vv != 0, nil
	case uint16:
		return vv != 0, nil
	case uint32:
		return vv != 0, nil
	case uint64:
		return vv != 0, nil

	case float32:
		return vv != 0, nil
	case float64:
		return vv != 0, nil

	case big.Float:
		return !(*bigFloat)(&vv).IsZero(), nil
	case *big.Float:
		return !(*bigFloat)(vv).IsZero(), nil

	// case decimalNumber:
	// 	return !(*internalDecimal)(&vv).IsZero(), nil
	// case *decimalNumber:
	// 	return !(*internalDecimal)(vv).IsZero(), nil

	case json.Number:
		return (*jsonNumber)(&vv).Bool()
	case *json.Number:
		return (*jsonNumber)(vv).Bool()

	case []byte:
		return convBool(vv)
	case string:
		return convBool(StrToBytesUnsafe(vv))
	}
	return false, _ERR_INVALID_TYPE
}

func ToInt(v any) (int, error) {
	n, err := ToInt64(v)
	return int(n), err
}
func ToInt8(v any) (int8, error) {
	n, err := ToInt64(v)
	return int8(n), err
}
func ToInt16(v any) (int16, error) {
	n, err := ToInt64(v)
	return int16(n), err
}
func ToInt32(v any) (int32, error) {
	n, err := ToInt64(v)
	return int32(n), err
}

func ToInt64(v any) (int64, error) {
	switch vv := v.(type) {
	case bool:
		if vv {
			return 1, nil
		}
		return 0, nil

	case int:
		return int64(vv), nil
	case int8:
		return int64(vv), nil
	case int16:
		return int64(vv), nil
	case int32:
		return int64(vv), nil
	case int64:
		return int64(vv), nil

	case uint:
		return int64(vv), nil
	case uint8:
		return int64(vv), nil
	case uint16:
		return int64(vv), nil
	case uint32:
		return int64(vv), nil
	case uint64:
		return int64(vv), nil

	case float32:
		return int64(vv), nil
	case float64:
		return int64(vv), nil

	case big.Float:
		return (*bigFloat)(&vv).Int64(), nil
	case *big.Float:
		return (*bigFloat)(vv).Int64(), nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).Int64(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).Int64(), nil

	case json.Number:
		return (*jsonNumber)(&vv).Int64()
	case *json.Number:
		return (*jsonNumber)(vv).Int64()

	case []byte:
		return convInt(vv)
	case string:
		return convInt(StrToBytesUnsafe(vv))
	}
	return 0, _ERR_INVALID_TYPE
}

func ToUint(v any) (uint, error) {
	n, err := ToUint64(v)
	return uint(n), err
}
func ToUint8(v any) (uint8, error) {
	n, err := ToUint64(v)
	return uint8(n), err
}
func ToUint16(v any) (uint16, error) {
	n, err := ToUint64(v)
	return uint16(n), err
}
func ToUint32(v any) (uint32, error) {
	n, err := ToUint64(v)
	return uint32(n), err
}

func ToUint64(v any) (uint64, error) {
	switch vv := v.(type) {
	case bool:
		if vv {
			return 1, nil
		}
		return 0, nil

	case int:
		return uint64(vv), nil
	case int8:
		return uint64(vv), nil
	case int16:
		return uint64(vv), nil
	case int32:
		return uint64(vv), nil
	case int64:
		return uint64(vv), nil

	case uint:
		return uint64(vv), nil
	case uint8:
		return uint64(vv), nil
	case uint16:
		return uint64(vv), nil
	case uint32:
		return uint64(vv), nil
	case uint64:
		return uint64(vv), nil

	case float32:
		return uint64(vv), nil
	case float64:
		return uint64(vv), nil

	case big.Float:
		return (*bigFloat)(&vv).Uint64(), nil
	case *big.Float:
		return (*bigFloat)(vv).Uint64(), nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).Uint64(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).Uint64(), nil

	case json.Number:
		return (*jsonNumber)(&vv).Uint64()
	case *json.Number:
		return (*jsonNumber)(vv).Uint64()

	case []byte:
		return convUint(vv)
	case string:
		return convUint(StrToBytesUnsafe(vv))
	}
	return 0, _ERR_INVALID_TYPE
}

func ToUintptr(v any) (uintptr, error) {
	n, err := ToUint64(v)
	return uintptr(n), err
}

func ToFloat32(v any) (float32, error) {
	switch vv := v.(type) {
	case bool:
		if vv {
			return 1, nil
		}
		return 0, nil

	case int:
		return float32(vv), nil
	case int8:
		return float32(vv), nil
	case int16:
		return float32(vv), nil
	case int32:
		return float32(vv), nil
	case int64:
		return float32(vv), nil

	case uint:
		return float32(vv), nil
	case uint8:
		return float32(vv), nil
	case uint16:
		return float32(vv), nil
	case uint32:
		return float32(vv), nil
	case uint64:
		return float32(vv), nil

	case float32:
		return vv, nil
	case float64:
		return float32(vv), nil

	case big.Float:
		return (*bigFloat)(&vv).Float32(), nil
	case *big.Float:
		return (*bigFloat)(vv).Float32(), nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).Float32(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).Float32(), nil

	case json.Number:
		return (*jsonNumber)(&vv).Float32()
	case *json.Number:
		return (*jsonNumber)(vv).Float32()

	case []byte:
		f, err := convFloat(vv)
		r, _ := f.Rat().Float32()
		return r, err
	case string:
		f, err := convFloat(StrToBytesUnsafe(vv))
		r, _ := f.Rat().Float32()
		return r, err
	}
	return 0, _ERR_INVALID_TYPE
}

func ToFloat64(v any) (float64, error) {
	switch vv := v.(type) {
	case bool:
		if vv {
			return 1, nil
		}
		return 0, nil

	case int:
		return float64(vv), nil
	case int8:
		return float64(vv), nil
	case int16:
		return float64(vv), nil
	case int32:
		return float64(vv), nil
	case int64:
		return float64(vv), nil

	case uint:
		return float64(vv), nil
	case uint8:
		return float64(vv), nil
	case uint16:
		return float64(vv), nil
	case uint32:
		return float64(vv), nil
	case uint64:
		return float64(vv), nil

	case float32:
		return float64(vv), nil
	case float64:
		return float64(vv), nil

	case big.Float:
		return (*bigFloat)(&vv).Float64(), nil
	case *big.Float:
		return (*bigFloat)(vv).Float64(), nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).Float64(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).Float64(), nil

	case json.Number:
		return (*jsonNumber)(&vv).Float64()
	case *json.Number:
		return (*jsonNumber)(vv).Float64()

	case []byte:
		f, err := convFloat(vv)
		r, _ := f.Rat().Float64()
		return r, err
	case string:
		f, err := convFloat(StrToBytesUnsafe(vv))
		r, _ := f.Rat().Float64()
		return r, err
	}
	return 0, _ERR_INVALID_TYPE
}

func ToString(v any) (string, error) {
	switch vv := v.(type) {
	case bool:
		return strconv.FormatBool(vv), nil

	case int:
		return strconv.FormatInt(int64(vv), 10), nil
	case int8:
		return strconv.FormatInt(int64(vv), 10), nil
	case int16:
		return strconv.FormatInt(int64(vv), 10), nil
	case int32:
		return strconv.FormatInt(int64(vv), 10), nil
	case int64:
		return strconv.FormatInt(int64(vv), 10), nil

	case uint:
		return strconv.FormatUint(uint64(vv), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(vv), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(vv), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(vv), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(vv), 10), nil

	case float32:
		return strconv.FormatFloat(float64(vv), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(float64(vv), 'f', -1, 64), nil

	case big.Float:
		return (*bigFloat)(&vv).String(), nil
	case *big.Float:
		return (*bigFloat)(vv).String(), nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).String(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).String(), nil

	case []byte:
		return BytesToStrUnsafe(vv), nil
	case string:
		return vv, nil

	case json.Number:
		return (*jsonNumber)(&vv).String(), nil
	case *json.Number:
		return (*jsonNumber)(vv).String(), nil

	case time.Time:
		return vv.Format(timeFormat), nil
	case *time.Time:
		return vv.Format(timeFormat), nil

	case time.Duration:
		return vv.String(), nil
	case *time.Duration:
		return vv.String(), nil
	}
	return "", _ERR_INVALID_TYPE
}

func ToBytes(v any) ([]byte, error) {
	switch vv := v.(type) {
	case big.Float:
		return (*bigFloat)(&vv).Bytes(), nil
	case *big.Float:
		return (*bigFloat)(vv).Bytes(), nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).Bytes(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).Bytes(), nil

	case []byte:
		return vv, nil
	case string:
		return StrToBytes(vv), nil

	case json.Number:
		return (*jsonNumber)(&vv).Bytes(), nil
	case *json.Number:
		return (*jsonNumber)(vv).Bytes(), nil
	}

	s, err := ToString(v)
	if err != nil {
		return nil, err
	}

	return StrToBytes(s), nil
}

func toBigFloat(v any) (*big.Float, error) {
	f := &big.Float{}
	switch vv := v.(type) {
	case bool:
		if vv {
			return f.SetInt64(1), nil
		} else {
			return f.SetInt64(0), nil
		}

	case int:
		return f.SetInt64(int64(vv)), nil
	case int8:
		return f.SetInt64(int64(vv)), nil
	case int16:
		return f.SetInt64(int64(vv)), nil
	case int32:
		return f.SetInt64(int64(vv)), nil
	case int64:
		return f.SetInt64(vv), nil

	case uint:
		return f.SetUint64(uint64(vv)), nil
	case uint8:
		return f.SetUint64(uint64(vv)), nil
	case uint16:
		return f.SetUint64(uint64(vv)), nil
	case uint32:
		return f.SetUint64(uint64(vv)), nil
	case uint64:
		return f.SetUint64(vv), nil

	case float32:
		return f.SetFloat64(float64(vv)), nil
	case float64:
		return f.SetFloat64(vv), nil

	case big.Float:
		return &vv, nil
	case *big.Float:
		return vv, nil

	// case decimalNumber:
	// 	return (*internalDecimal)(&vv).BigFloat(), nil
	// case *decimalNumber:
	// 	return (*internalDecimal)(vv).BigFloat(), nil

	case json.Number:
		return (*jsonNumber)(&vv).BigFloat()
	case *json.Number:
		return (*jsonNumber)(vv).BigFloat()

	case []byte:
		d, err := convFloat(vv)
		return d.BigFloat(), err
	case string:
		d, err := convFloat(StrToBytesUnsafe(vv))
		return d.BigFloat(), err
	}

	return newBigFloatZero(), _ERR_INVALID_TYPE
}

func toDecimal(v any) (decimalNumber, error) {
	switch vv := v.(type) {
	case bool:
		if vv {
			return newDecimalFromInt(1), nil
		} else {
			return newDecimalZero(), nil
		}

	case int:
		return newDecimalFromInt(int64(vv)), nil
	case int8:
		return newDecimalFromInt(int64(vv)), nil
	case int16:
		return newDecimalFromInt(int64(vv)), nil
	case int32:
		return newDecimalFromInt(int64(vv)), nil
	case int64:
		return newDecimalFromInt(vv), nil

	case uint:
		return newDecimalFromUint(uint64(vv)), nil
	case uint8:
		return newDecimalFromUint(uint64(vv)), nil
	case uint16:
		return newDecimalFromUint(uint64(vv)), nil
	case uint32:
		return newDecimalFromUint(uint64(vv)), nil
	case uint64:
		return newDecimalFromUint(vv), nil

	case float32:
		return newDecimalFromFloat(float64(vv)), nil
	case float64:
		return newDecimalFromFloat(vv), nil

	case big.Float:
		return newDecimalFromBigFloat(&vv), nil
	case *big.Float:
		return newDecimalFromBigFloat(vv), nil

	case decimalNumber:
		return vv, nil
	case *decimalNumber:
		return *vv, nil

	case json.Number:
		return (*jsonNumber)(&vv).Decimal()
	case *json.Number:
		return (*jsonNumber)(vv).Decimal()

	case []byte:
		return convFloat(vv)
	case string:
		return convFloat(StrToBytesUnsafe(vv))
	}
	return newDecimalZero(), _ERR_INVALID_TYPE
}

func ToJsonNumber(v any) (json.Number, error) {
	switch vv := v.(type) {
	case bool:
		if vv {
			return "1", nil
		}
		return "0", nil

	case json.Number:
		if vv == "" {
			return "0", nil
		}
		return vv, nil

	case *json.Number:
		if vv == nil || *vv == "" {
			return "0", nil
		}
		return *vv, nil
	}

	s, err := ToString(v)
	return json.Number(s), err
}
