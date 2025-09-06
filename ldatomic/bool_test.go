/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"fmt"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func testStoreLoadCompareAndSwapper[I StoreLoadCompareAndSwapper[T], T any](t *testing.T, fNew func(T) I, values ...T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(len(values), convey.ShouldBeGreaterThanOrEqualTo, 2)

		for i, l := 0, len(values); i < l; i++ {
			for j := i + 1; j < l; j++ {
				v0 := values[i]
				v1 := values[j]
				name := fmt.Sprintf("%v-%v", v0, v1)
				c.Convey(name, func(c convey.C) {
					p := fNew(v0)

					c.So(p.Load(), convey.ShouldResemble, v0)

					p.Store(v1)
					c.So(p.Load(), convey.ShouldResemble, v1)

					c.So(p.Swap(v0), convey.ShouldResemble, v1)
					c.So(p.Load(), convey.ShouldResemble, v0)

					c.So(p.CompareAndSwap(v1, v0), convey.ShouldResemble, false)
					c.So(p.Load(), convey.ShouldResemble, v0)

					c.So(p.CompareAndSwap(v0, v1), convey.ShouldResemble, true)
					c.So(p.Load(), convey.ShouldResemble, v1)
				})
			}
		}
	})
}

func TestNewBool(t *testing.T) {
	testStoreLoadCompareAndSwapper(t, NewBool, false, true)
}
