package test

import (
	"testing"

	. "github.com/distroy/ldgo-base/3rd/convey"
	. "github.com/distroy/ldgo-base/3rd/gomonkey"
)

var num = 10

func TestApplyGlobalVar(t *testing.T) {
	Convey("TestApplyGlobalVar", t, func() {

		Convey("change", func() {
			patches := ApplyGlobalVar(&num, 150)
			defer patches.Reset()
			So(num, ShouldEqual, 150)
		})

		Convey("recover", func() {
			So(num, ShouldEqual, 10)
		})
	})
}
