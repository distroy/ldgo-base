/*
 * Copyright (C) distroy
 */

package ldsort

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestUniq(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		f := func(a, b int64) int { return int(a - b) }
		s := []int64{
			223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375,
		}

		Sort(s, f)
		c.So(s, convey.ShouldResemble, []int64{
			123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
		})

		s = Uniq(s, f)
		c.So(s, convey.ShouldResemble, []int64{
			123, 223, 375, 424, 642, 725,
		})
	})
}
