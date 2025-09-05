/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestNewStrings(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(NewStrings(nil).Load(), convey.ShouldBeNil)
		c.So(NewStrings([]string{}).Load(), convey.ShouldResemble, []string{})

		p := Strings{}
		c.So(p.Load(), convey.ShouldBeNil)

		c.So(p.CompareAndSwap(nil, []string{}), convey.ShouldEqual, true)
		c.So(p.Load(), convey.ShouldResemble, []string{})

		c.So(p.CompareAndSwap(nil, []string{"abc"}), convey.ShouldEqual, false)
		c.So(p.Load(), convey.ShouldResemble, []string{})

		// c.So(p.CompareAndSwap([]string{}, []string{"abc"}), convey.ShouldEqual, false)
		// c.So(p.Load(), convey.ShouldResemble, []string{})
		//
		// c.So(p.CompareAndSwap(p.Load(), []string{"abc"}), convey.ShouldEqual, true)
		// c.So(p.Load(), convey.ShouldResemble, []string{"abc"})
		c.So(p.CompareAndSwap([]string{}, []string{"abc"}), convey.ShouldEqual, true)
		c.So(p.Load(), convey.ShouldResemble, []string{"abc"})

		c.So(p.CompareAndSwap([]string{"abc"}, []string{"a"}), convey.ShouldEqual, false)
		c.So(p.Load(), convey.ShouldResemble, []string{"abc"})

		c.So(p.CompareAndSwap(p.Load(), []string{"a"}), convey.ShouldEqual, true)
		c.So(p.Load(), convey.ShouldResemble, []string{"a"})
	})
}
