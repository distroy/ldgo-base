/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestWrapFunc(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		fn := func(a int, b string) string { return fmt.Sprintf("a:%d, b:%s", a, b) }

		errType := reflect.TypeOf((*error)(nil)).Elem()

		c.Convey("append input/output", func(c convey.C) {
			opts := []FuncWrapperOption{
				AppendInput(true),

				AppendOutput(nil, errType),
			}

			fn0 := WrapFunc(fn, opts...)
			fn1, ok := fn0.(func(a int, b string, c bool) (string, error))
			c.So(ok, convey.ShouldBeTrue)

			c.Convey("direct call", func(c convey.C) {
				// c.So(func() { fn1(1, "abc", false) }, convey.ShouldPanic)
				r0, r1 := fn1(1, "abc", false)
				c.So(r0, convey.ShouldEqual, "a:1, b:abc")
				c.So(r1, convey.ShouldBeNil)
			})

			c.Convey("reflect call", func(c convey.C) {
				fnV := reflect.ValueOf(fn0)
				ins := []reflect.Value{reflect.ValueOf(1), reflect.ValueOf("abc"), reflect.ValueOf(false)}

				outs := fnV.Call(ins)
				c.So(outs, convey.ShouldHaveLength, 2)
				c.So(outs[0].Interface(), convey.ShouldEqual, "a:1, b:abc")
				c.So(outs[1].Interface(), convey.ShouldBeNil)
			})
		})

		c.Convey("append/swap input/output", func(c convey.C) {
			opts := []FuncWrapperOption{
				AppendInput(true),
				SwapInput(0, 2),

				AppendOutput(nil, errType),
				AppendOutput(true),
				SwapOutput(1, 2),
			}

			fn0 := WrapFunc(fn, opts...)
			fn1, ok := fn0.(func(a bool, b string, c int) (string, bool, error))
			c.So(ok, convey.ShouldBeTrue)

			r0, r1, r2 := fn1(false, "abc", 1)
			c.So(r0, convey.ShouldEqual, "a:1, b:abc")
			c.So(r1, convey.ShouldEqual, true)
			c.So(r2, convey.ShouldBeNil)
		})

		c.Convey("append/swap/add input/output", func(c convey.C) {
			opts := []FuncWrapperOption{
				AppendInput(true),
				SwapInput(0, 2),
				InsertInput(3, "xyz"),

				AppendOutput(nil, errType),
				AppendOutput(true),
				SwapOutput(1, 2),
				InsertOutput(2, "xyz"),
			}

			fn0 := WrapFunc(fn, opts...)
			fn1 := fn0.(func(a bool, b string, c int, d string) (string, bool, string, error))
			fn1, ok := fn0.(func(a bool, b string, c int, d string) (string, bool, string, error))
			c.So(ok, convey.ShouldBeTrue)

			r0, r1, r2, r3 := fn1(false, "abc", 1, "zzz")
			c.So(r0, convey.ShouldEqual, "a:1, b:abc")
			c.So(r1, convey.ShouldEqual, true)
			c.So(r2, convey.ShouldEqual, "xyz")
			c.So(r3, convey.ShouldBeNil)
		})

		c.Convey("variadic", func(c convey.C) {
			fn := func(a int, b string, c ...int) string {
				return fmt.Sprintf("a:%d, b:%s, c:%+v", a, b, c)
			}

			c.Convey("swap last input parameter", func(c convey.C) {
				c.So(func() {
					WrapFunc(fn, SwapInput(0, 2))
				}, convey.ShouldPanic)
			})

			c.Convey("append input", func(c convey.C) {
				c.So(func() {
					WrapFunc(fn, AppendInput(true))
				}, convey.ShouldPanic)
			})

			c.Convey("add input", func(c convey.C) {
				c.Convey("out of range", func(c convey.C) {
					c.So(func() {
						WrapFunc(fn, InsertInput(4, "xyz"))
					}, convey.ShouldPanic)
				})
				c.Convey("after last", func(c convey.C) {
					c.So(func() {
						WrapFunc(fn, InsertInput(3, "xyz"))
					}, convey.ShouldPanic)
				})
			})

			c.Convey("succ", func(c convey.C) {
				opts := []FuncWrapperOption{
					InsertInput(2, true),
					AppendOutput(nil, errType),
				}

				fn0 := WrapFunc(fn, opts...)
				fn1, ok := fn0.(func(a int, b string, c bool, d ...int) (string, error))
				c.So(ok, convey.ShouldBeTrue)

				r0, r1 := fn1(1, "abc", false, 2, 3, 4)
				c.So(r0, convey.ShouldEqual, "a:1, b:abc, c:[2 3 4]")
				c.So(r1, convey.ShouldBeNil)

			})
		})
	})
}
