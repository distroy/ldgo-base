/*
 * Copyright (C) distroy
 */

package ldctx

import (
	"context"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestWithMap(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		ctx0 := context.Background()

		c.Convey("no map", func(c convey.C) {
			ctx := ctx0
			m := GetMap(ctx)
			c.So(m, convey.ShouldBeNil)
			c.So(m.Get("xxxx"), convey.ShouldBeNil)
			c.So(m.Set("xxxx", "xxxx"), convey.ShouldEqual, false)
		})
		c.Convey("with map", func(c convey.C) {
			ctx1 := WithMap(ctx0)
			ctx := ctx1

			m := GetMap(ctx)
			c.So(m, convey.ShouldNotBeNil)
			c.So(m.Get("xxxx"), convey.ShouldBeNil)

			m.Set("a", 1)
			c.So(m.Get("a"), convey.ShouldEqual, 1)
			c.So(GetMap(ctx).Get("a"), convey.ShouldEqual, 1)

			m.Set("b", 2)
			c.So(m.Get("b"), convey.ShouldEqual, 2)
			c.So(GetMap(ctx).Get("b"), convey.ShouldEqual, 2)

			c.Convey("with map again", func(c convey.C) {
				ctx2 := WithMap(ctx1)
				ctx := ctx2

				m := GetMap(ctx)
				c.So(m, convey.ShouldNotBeNil)
				c.So(m.Get("a"), convey.ShouldEqual, 1)
				c.So(m.Get("b"), convey.ShouldEqual, 2)

				m.Set("a", 100)
				c.So(m.Get("a"), convey.ShouldEqual, 100)
				c.So(GetMap(ctx).Get("a"), convey.ShouldEqual, 100)
				c.So(GetMap(ctx).Get("b"), convey.ShouldEqual, 2)

				c.So(GetMap(ctx1).Get("a"), convey.ShouldEqual, 1)
				c.So(GetMap(ctx2).Get("b"), convey.ShouldEqual, 2)

				c.Convey("clear", func(c convey.C) {
					ctx := ctx2

					m := GetMap(ctx)
					c.So(m.Get("a"), convey.ShouldEqual, 100)
					c.So(m.Get("b"), convey.ShouldEqual, 2)

					m.Clear()
					c.So(m.Get("a"), convey.ShouldBeNil)
					c.So(m.Get("b"), convey.ShouldBeNil)
				})

				c.Convey("with map again twice", func(c convey.C) {
					ctx3 := WithMap(ctx2)
					ctx := ctx3

					m := GetMap(ctx)
					c.So(m, convey.ShouldNotBeNil)
					c.So(m.Get("a"), convey.ShouldEqual, 100)
					c.So(m.Get("b"), convey.ShouldEqual, 2)
				})
			})
		})
	})
}
