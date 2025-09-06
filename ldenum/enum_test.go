/*
 * Copyright (C) distroy
 */

package ldenum

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldptr"
)

type RetCodeToString struct{}

func (RetCodeToString) EnumToString(n int) string {
	switch n {
	case 0:
		return "ok"
	case 1:
		return "panic"
	}
	return "unkown error"
}

type RetCode = Enum[RetCodeToString]

func TestEnum(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		n := RetCode(1)
		c.Convey("enum", func(c convey.C) {
			c.So(n, convey.ShouldResemble, RetCode(1))
			c.So(n.Int(), convey.ShouldResemble, 1)
			c.So(n.Str(), convey.ShouldResemble, "panic")
		})

		p := n.Ptr()
		c.Convey("ptr", func(c convey.C) {
			c.So(p, convey.ShouldResemble, ldptr.New(n))
			c.So(p.Get(), convey.ShouldResemble, n)
			c.So(p.GetInt(), convey.ShouldResemble, int(n))
			c.So(p.GetStr(), convey.ShouldResemble, "panic")

			c.So(p.New(), convey.ShouldResemble, ldptr.New(n))
			c.So(p.NewInt(), convey.ShouldResemble, ldptr.New(int(n)))
			c.So(p.NewStr(), convey.ShouldResemble, ldptr.New("panic"))
		})

		p = nil
		c.Convey("nil", func(c convey.C) {
			c.So(p.Get(), convey.ShouldResemble, RetCode(0))
			c.So(p.GetInt(), convey.ShouldResemble, 0)
			c.So(p.GetStr(), convey.ShouldResemble, "ok")

			c.So(p.New(), convey.ShouldResemble, (*RetCode)(nil))
			c.So(p.NewInt(), convey.ShouldResemble, (*int)(nil))
			c.So(p.NewStr(), convey.ShouldResemble, (*string)(nil))
		})
	})
}
