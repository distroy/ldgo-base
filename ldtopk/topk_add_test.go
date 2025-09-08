/*
 * Copyright (C) distroy
 */

package ldtopk

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/internal/cmp_"
	"github.com/distroy/ldgo-base/ldrand"
	"github.com/distroy/ldgo-base/ldsort"
)

func testTopkAdd[T cmp_.Comparable](c convey.C, n, k int, rand func() T) {
	name := fmt.Sprintf("n:%d-k:%d", n, k)
	c.Convey(name, func() {
		origin := make([]T, 0, n)
		topk := make([]T, 0, k)

		for range n {
			x := rand()
			origin = append(origin, x)
			topk, _ = TopkAdd(topk, k, x)
		}

		ldsort.Sort(origin, cmp_.CompareComparable[T])
		ldsort.Sort(topk, cmp_.CompareComparable[T])

		size := min(n, k)
		origin = origin[:size]
		convey.So(topk, convey.ShouldResemble, origin)
	})
}

func TestTopk(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("int", func(c convey.C) {
			rand := func() int { return ldrand.Intn(256) }
			testTopkAdd(c, 10, 20, rand)
			testTopkAdd(c, 100, 5, rand)
			testTopkAdd(c, 100, 10, rand)
		})

		c.Convey("int64", func(c convey.C) {
			rand := func() int64 { return int64(ldrand.Intn(256)) }
			testTopkAdd(c, 10, 20, rand)
			testTopkAdd(c, 100, 5, rand)
			testTopkAdd(c, 100, 10, rand)
		})

		c.Convey("float64", func(c convey.C) {
			rand := func() float64 { return ldrand.ExpFloat64() }
			testTopkAdd(c, 10, 20, rand)
			testTopkAdd(c, 100, 5, rand)
			testTopkAdd(c, 100, 10, rand)
		})

		c.Convey("string", func(c convey.C) {
			rand := func() string { return ldrand.String(8) }
			testTopkAdd(c, 10, 20, rand)
			testTopkAdd(c, 100, 5, rand)
			testTopkAdd(c, 100, 10, rand)
		})
	})
}
