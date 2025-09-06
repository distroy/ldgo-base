/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"reflect"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestGetMethod(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		obj := &testMethodForHook{}

		c.Convey("panic by obj", func(c convey.C) {
			c.So(func() {
				GetMethod(obj, "1111")
			}, convey.ShouldPanic)
		})
		c.Convey("panic by type", func(c convey.C) {
			c.So(func() {
				GetMethod(reflect.TypeOf(obj), "1111")
			}, convey.ShouldPanic)
		})

		c.Convey("succ by obj", func(c convey.C) {
			f := GetMethod(obj, "Test")
			c.So(f, convey.ShouldEqual, (*testMethodForHook).Test)
		})
		c.Convey("succ by type", func(c convey.C) {
			f := GetMethod(reflect.TypeOf(obj), "Test")
			c.So(f, convey.ShouldEqual, (*testMethodForHook).Test)
		})
	})
}
