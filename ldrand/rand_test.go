/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestRandString(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		a := String(16)
		b := String(16)
		c.So(a, convey.ShouldHaveLength, 16)
		c.So(b, convey.ShouldHaveLength, 16)
		c.So(a, convey.ShouldNotEqual, b)
	})
}
