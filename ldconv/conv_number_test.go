/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

const (
	_TEST_FLOAT_MUL_BASE = 100000
)

func testConvInt(s string, r int64) {
	convey.Convey(s, func() {
		convey.So(AsInt64(s), convey.ShouldEqual, r)
	})
}

func testConvUint(s string, r uint64) {
	convey.Convey(s, func() {
		convey.So(AsUint64(s), convey.ShouldEqual, r)
	})
}

func testConvFloat(s string, rf32 float32, rf64 float64, ri64 int64) {
	name := fmt.Sprintf("%s|%g|%g|%d", s, rf32, rf64, ri64)
	convey.Convey(name, func() {

		convey.So(AsFloat32(s), convey.ShouldEqual, rf32)

		convey.So(AsFloat64(s), convey.ShouldEqual, rf64)

		n0 := asDecimal(s)
		baseFloat := asDecimal(_TEST_FLOAT_MUL_BASE)
		n1 := n0.Mul(baseFloat)
		vi64, _ := n1.BigFloat().Int64()
		convey.So(vi64, convey.ShouldEqual, ri64)
	})
}

func TestConvInt(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		testConvInt("0", 0)
		testConvInt("123", 123)
		testConvInt("123.1", 123)
		testConvInt(".1", 0)
		testConvInt("0123", 123)
		testConvInt("0o123", 0123)
		testConvInt("0O123", 0123)
		testConvInt("+-0o123", -0123)
		testConvInt("123e+1", 1230)
		testConvInt("123e1", 1230)
		testConvInt("123e-1", 12)
		testConvInt("0x123", 0x123)
		testConvInt("--0X123", 0x123)
		testConvInt("-+-0x123", 0x123)
		testConvInt("-+-0xeFaB", 0xeFaB)
	})
}

func TestConvUint(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		testConvUint("0", 0)
		testConvUint("123", 123)
		testConvUint("123.1", 123)
		testConvUint(".1", 0)
		testConvUint("0123", 123)
		testConvUint("0o123", 0123)
		// testConvUint("+-0123", uint64(int64(-0123)))
		testConvUint("123e+1", 1230)
		testConvUint("123e1", 1230)
		testConvUint("123e-1", 12)
		testConvUint("0x123", 0x123)
		testConvUint("--0X123", 0x123)
		testConvUint("-+-0x123", 0x123)
		testConvUint("-+-0xeFaB", 0xeFaB)
	})
}

func TestConvFloat(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		testConvFloat("0", 0, 0, 0*_TEST_FLOAT_MUL_BASE)
		testConvFloat("9.9", 9.9, 9.9, 9.9*_TEST_FLOAT_MUL_BASE)
		testConvFloat("20.09", 20.09, 20.09, 20.09*_TEST_FLOAT_MUL_BASE)
		testConvFloat("19.99", 19.99, 19.99, 19.99*_TEST_FLOAT_MUL_BASE)
		testConvFloat("0.99", 0.99, 0.99, 0.99*_TEST_FLOAT_MUL_BASE)
		testConvFloat("2.32", 2.32, 2.32, 2.32*_TEST_FLOAT_MUL_BASE)
		testConvFloat("123", 123, 123, 123*_TEST_FLOAT_MUL_BASE)
		testConvFloat("123.1", 123.1, 123.1, 123.1*_TEST_FLOAT_MUL_BASE)
		testConvFloat(".1", 0.1, 0.1, 0.1*_TEST_FLOAT_MUL_BASE)
		testConvFloat("0123", 123, 123, 123*_TEST_FLOAT_MUL_BASE)
		testConvFloat("0o123", 0123, 0123, 0123*_TEST_FLOAT_MUL_BASE)
		// testConvFloat("+-0123", uint64(int64(-0123)))
		testConvFloat("123e+1", 1230, 1230, 1230*_TEST_FLOAT_MUL_BASE)
		testConvFloat("123e1", 1230, 1230, 1230*_TEST_FLOAT_MUL_BASE)
		testConvFloat("123e-1", 12.3, 12.3, 12.3*_TEST_FLOAT_MUL_BASE)
		testConvFloat("0x123", 0x123, 0x123, 0x123*_TEST_FLOAT_MUL_BASE)
		testConvFloat("--0X123", 0x123, 0x123, 0x123*_TEST_FLOAT_MUL_BASE)
		testConvFloat("-+-0x123", 0x123, 0x123, 0x123*_TEST_FLOAT_MUL_BASE)
		testConvFloat("-+-0xeFaB", 0xeFaB, 0xeFaB, 0xeFaB*_TEST_FLOAT_MUL_BASE)
		testConvFloat("---0xeFaB", -0xeFaB, -0xeFaB, -0xeFaB*_TEST_FLOAT_MUL_BASE)
		testConvFloat("0.30129", 0.30129, 0.30129, 0.30129*_TEST_FLOAT_MUL_BASE)
		testConvFloat("0.30129e3", 301.29, 301.29, 301.29*_TEST_FLOAT_MUL_BASE)
		testConvFloat("-0.30129E-3", -0.00030129, -0.00030129, -30)
	})
}
