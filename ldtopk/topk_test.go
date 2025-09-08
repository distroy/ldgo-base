/*
 * Copyright (C) distroy
 */

package ldtopk

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldrand"
	"github.com/distroy/ldgo-base/ldsort"
)

func testTopk(c convey.C, n, k int) {
	// less := func(a, b int) bool { return a < b }

	name := fmt.Sprintf("n:%d-k:%d", n, k)
	c.Convey(name, func(c convey.C) {
		origin := make([]int, 0, n)

		topk := New[int](k, nil)

		for range n {
			x := ldrand.Intn(100)
			origin = append(origin, x)
			topk.Add(x)
		}

		ldsort.SortInts(origin)
		ldsort.SortInts(topk.Data())

		size := min(n, k)
		origin = origin[:size]
		c.So(topk.Data(), convey.ShouldResemble, origin)
	})
}

func TestTopK(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		testTopk(c, 10, 20)
		testTopk(c, 100, 5)
		testTopk(c, 100, 10)
	})
}
