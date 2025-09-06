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
