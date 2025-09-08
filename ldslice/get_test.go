/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestGet(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		s := []string{"a", "b", "c"}
		c.So(Get(s, 0), convey.ShouldEqual, "a")
		c.So(Get(s, 1), convey.ShouldEqual, "b")
		c.So(Get(s, 2), convey.ShouldEqual, "c")
		c.So(Get(s, 3), convey.ShouldEqual, "")
		c.So(Get(s, 3, `x`), convey.ShouldEqual, "x")
	})
}
