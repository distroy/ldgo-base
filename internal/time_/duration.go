/*
 * Copyright (C) distroy
 */

package time_

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

func DurationMarshalJson(d time.Duration) ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+4)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

func DurationUnmarshalJson(b []byte) (time.Duration, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("unexpected end of JSON input")
	}

	switch b[0] {
	case '"', '\'', '`':
		return durationUnmarshalJsonByStr(b)
	}

	return durationUnmarshalJsonByNumber(b)
}

func durationUnmarshalJSONError(b []byte) error {
	return fmt.Errorf("invalid duration: %s", b)
}

func durationUnmarshalJsonByStr(b []byte) (time.Duration, error) {
	str := unsafe.String(unsafe.SliceData(b), len(b))
	vv, err := strconv.Unquote(str)
	if err != nil {
		return 0, durationUnmarshalJSONError(b)
	}
	dur, err := time.ParseDuration(vv)
	if err != nil {
		return 0, durationUnmarshalJSONError(b)
	}
	return dur, nil
}

func DurationUnmarshalYaml(b []byte) (time.Duration, error) {
	// log.Printf(" === UnmarshalYAML: %s", n.Value)
	if dur, err := ParseDurationByNumber(b); err == nil {
		return dur, nil
	}
	s := unsafe.String(unsafe.SliceData(b), len(b))
	return time.ParseDuration(s)
}

func ParseDurationByNumber(b []byte) (time.Duration, error) {
	return durationUnmarshalJsonByNumber(b)
}
func durationUnmarshalJsonByNumber(b []byte) (time.Duration, error) {
	if b[0] != '0' {
		str := unsafe.String(unsafe.SliceData(b), len(b))
		i64, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return time.Second * time.Duration(i64), nil
		}

		f, err := strconv.ParseFloat(str, 64)
		if err == nil {
			// r := time.Second * time.Duration(f)
			// r += time.Second * time.Duration(float64(time.Second)*(f-float64(time.Duration(f))))
			return time.Duration(f * float64(time.Second)), nil
		}
		return 0, durationUnmarshalJSONError(b)
	}

	if len(b) == 1 {
		return 0, nil
	}

	switch b[1] {
	case 'x', 'X':
		bb := b[2:]
		str := unsafe.String(unsafe.SliceData(bb), len(bb))
		if i64, err := strconv.ParseInt(str, 16, 64); err == nil {
			return time.Second * time.Duration(i64), nil
		}

	case 'o', 'O':
		bb := b[2:]
		str := unsafe.String(unsafe.SliceData(bb), len(bb))
		if i64, err := strconv.ParseInt(str, 8, 64); err == nil {
			return time.Second * time.Duration(i64), nil
		}

	default:
		bb := b[1:]
		str := unsafe.String(unsafe.SliceData(bb), len(bb))
		if i64, err := strconv.ParseInt(str, 8, 64); err == nil {
			return time.Second * time.Duration(i64), nil
		}
	}

	return 0, durationUnmarshalJSONError(b)
}
