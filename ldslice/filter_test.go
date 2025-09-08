/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestFilter(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("part elements have been filtered", func(c convey.C) {
			s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
			r0, r1 := Filter(s, func(d int) bool { return d%2 == 0 })
			c.So(r0, convey.ShouldResemble, []int{0, 2, 8, 4, 6})
			c.So(r1, convey.ShouldResemble, []int{5, 7, 3, 9, 1})
		})
		c.Convey("no element has been filtered", func(c convey.C) {
			s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
			r0, r1 := Filter(s, func(d int) bool { return true })
			c.So(r0, convey.ShouldResemble, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
			c.So(r1, convey.ShouldResemble, []int{})
		})
		c.Convey("all elements have been filtered", func(c convey.C) {
			s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
			r0, r1 := Filter(s, func(d int) bool { return false })
			c.So(r0, convey.ShouldResemble, []int{})
			c.So(r1, convey.ShouldResemble, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
		})
	})
}
