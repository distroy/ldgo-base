/*
 * Copyright (C) distroy
 */

package ldflag

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldptr"
)

func TestValues(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		type Flags struct {
			Duration time.Duration
			Func     func(string) error `ldflag:"-"`
			String   string
			Bool     bool
			Int      int
			Int64    int64
			Uint     uint
			Uint64   uint64
			Float32  float32
			Float64  float64
			Strings  []string
			Ints     []int
			Int64s   []int64
			Uints    []uint
			Uint64s  []uint64
			Float32s []float32
			Float64s []float64
		}

		flags := &Flags{}
		s := newTestFlagSet()
		s.Model(flags)

		err := s.Parse([]string{
			`-duration`, "1m",
			// `-func`, "xxx",
			`-string`, "abc",
			`-bool`, "true",
			`-int`, "-123",
			`-int64`, "-1234",
			`-uint`, "123",
			`-uint64`, "1234",
			`-float32`, "123.123",
			`-float64`, "1234.123",
			`-strings`, "s1",
			`-strings`, "s2",
			`-ints`, "-101",
			`-ints`, "-102",
			`-int64s`, "-10001",
			`-int64s`, "-10002",
			`-uints`, "101",
			`-uints`, "102",
			`-uint64s`, "10001",
			`-uint64s`, "10002",
			`-float32s`, "101.123",
			`-float32s`, "102.123",
			`-float64s`, "10001.123",
			`-float64s`, "10002.123",
		})
		c.So(err, convey.ShouldBeNil)
		c.So(flags, convey.ShouldResemble, &Flags{
			Duration: time.Minute,
			Func:     nil,
			String:   "abc",
			Bool:     true,
			Int:      -123,
			Int64:    -1234,
			Uint:     123,
			Uint64:   1234,
			Float32:  123.123,
			Float64:  1234.123,
			Strings:  []string{"s1", "s2"},
			Ints:     []int{-101, -102},
			Int64s:   []int64{-10001, -10002},
			Uints:    []uint{101, 102},
			Uint64s:  []uint64{10001, 10002},
			Float32s: []float32{101.123, 102.123},
			Float64s: []float64{10001.123, 10002.123},
		})
	})
}

func TestPtrValues(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		type Flags struct {
			Duration *time.Duration
			String   *string
			Bool     *bool
			Int      *int
			Int64    *int64
			Uint     *uint
			Uint64   *uint64
			Float32  *float32
			Float64  *float64
		}

		flags := &Flags{}
		s := newTestFlagSet()
		s.Model(flags)

		err := s.Parse([]string{
			`-duration`, "1m",
			`-string`, "abc",
			`-bool`, "true",
			`-int`, "-123",
			`-int64`, "-1234",
			`-uint`, "123",
			`-uint64`, "1234",
			`-float32`, "123.123",
			`-float64`, "1234.123",
		})
		c.So(err, convey.ShouldBeNil)
		c.So(flags, convey.ShouldResemble, &Flags{
			Duration: ldptr.New(time.Minute),
			String:   ldptr.New("abc"),
			Bool:     ldptr.New(true),
			Int:      ldptr.New(-123),
			Int64:    ldptr.New[int64](-1234),
			Uint:     ldptr.New[uint](123),
			Uint64:   ldptr.New[uint64](1234),
			Float32:  ldptr.New[float32](123.123),
			Float64:  ldptr.New(1234.123),
		})
	})
}
