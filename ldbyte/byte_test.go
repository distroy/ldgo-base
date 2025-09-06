/*
 * Copyright (C) distroy
 */

package ldbyte

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestToUpper(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		b0 := byte('a')
		b1 := byte('A')
		for i := byte(0); i < 'z'-'a'; i++ {
			c0 := b0 + i
			c1 := b1 + i
			r0 := ToUpper(c0)
			c.So(r0, convey.ShouldEqual, c1)
		}
		c.So(ToUpper(' '), convey.ShouldEqual, ' ')
	})
}

func TestToLower(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		b0 := byte('A')
		b1 := byte('a')
		for i := byte(0); i < 'z'-'a'; i++ {
			c0 := b0 + i
			c1 := b1 + i
			r0 := ToLower(c0)
			c.So(r0, convey.ShouldEqual, c1)
		}
		c.So(ToLower(' '), convey.ShouldEqual, ' ')
	})
}
