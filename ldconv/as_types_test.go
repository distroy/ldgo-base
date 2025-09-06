/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestAsBool(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(AsBool("abc"), convey.ShouldEqual, false)
		convey.So(AsBool("abc", true), convey.ShouldEqual, true)

		convey.So(AsBool("true", false), convey.ShouldEqual, true)
		convey.So(AsBool("false", true), convey.ShouldEqual, false)
		convey.So(AsBool([]byte("null"), true), convey.ShouldEqual, false)

		convey.So(AsBool(0, true), convey.ShouldEqual, false)

		convey.So(AsBool(int(-1)), convey.ShouldEqual, true)
		convey.So(AsBool(uint(1)), convey.ShouldEqual, true)
	})
}

func TestAsInt(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(AsInt("abc"), convey.ShouldEqual, 0)
		convey.So(AsInt("abc", 1), convey.ShouldEqual, 1)

		convey.So(AsInt("true", 0), convey.ShouldEqual, 1)
		convey.So(AsInt("false", 1), convey.ShouldEqual, 0)
		convey.So(AsInt([]byte("null"), 1), convey.ShouldEqual, 0)

		convey.So(AsInt(0, 100), convey.ShouldEqual, 0)

		convey.So(AsInt(int(-1)), convey.ShouldEqual, -1)
		convey.So(AsInt(uint(1)), convey.ShouldEqual, 1)

		// convey.So(AsInt(mustNewDecimalFromStr("-123.234")), convey.ShouldEqual, -123)
	})
}

func TestAsFloat32(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(AsFloat32("0"), convey.ShouldEqual, 0)
		convey.So(AsFloat32("abc"), convey.ShouldEqual, 0)
		convey.So(AsFloat32("abc", 1), convey.ShouldEqual, 1)

		convey.So(AsFloat32("true", 0), convey.ShouldEqual, 1)
		convey.So(AsFloat32("false", 1), convey.ShouldEqual, 0)
		convey.So(AsFloat32([]byte("null"), 1), convey.ShouldEqual, 0)

		convey.So(AsFloat32(0, 100), convey.ShouldEqual, 0)

		convey.So(AsFloat32(int(-1)), convey.ShouldEqual, -1)
		convey.So(AsFloat32(uint(1)), convey.ShouldEqual, 1)

		convey.So(AsFloat32("-123.234"), convey.ShouldEqual, float32(-123.234))
		// convey.So(AsFloat32(mustNewDecimalFromStr("-123.234")), convey.ShouldEqual, float32(-123.234))
	})
}

func TestAsFloat64(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(AsFloat64("0"), convey.ShouldEqual, 0)
		convey.So(AsFloat64("abc"), convey.ShouldEqual, 0)
		convey.So(AsFloat64("abc", 1), convey.ShouldEqual, 1)

		convey.So(AsFloat64("true", 0), convey.ShouldEqual, 1)
		convey.So(AsFloat64("false", 1), convey.ShouldEqual, 0)
		convey.So(AsFloat64([]byte("null"), 1), convey.ShouldEqual, 0)

		convey.So(AsFloat64(0, 100), convey.ShouldEqual, 0)

		convey.So(AsFloat64(int(-1)), convey.ShouldEqual, -1)
		convey.So(AsFloat64(uint(1)), convey.ShouldEqual, 1)

		// convey.So(AsFloat64(mustNewDecimalFromStr("-123.234")), convey.ShouldEqual, -123.234)
	})
}
