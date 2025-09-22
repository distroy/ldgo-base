/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sort"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestMap(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("new map by nil", func(c convey.C) {
			m := NewMap[int, int](nil)

			c.So(m.Size(), convey.ShouldEqual, 0)
			c.So(m.Map(), convey.ShouldResemble, map[int]int{})

			v, loaded := m.Load(123)
			c.So(loaded, convey.ShouldBeFalse)
			c.So(v, convey.ShouldEqual, 0)

			v, loaded = m.Swap(123, 100)
			c.So(loaded, convey.ShouldBeFalse)
			c.So(v, convey.ShouldEqual, 0)

			c.So(m.Size(), convey.ShouldEqual, 1)
			c.So(m.Get(123), convey.ShouldEqual, 100)

			m.Set(123, 200)
			c.So(m.Size(), convey.ShouldEqual, 1)
			c.So(m.Get(123), convey.ShouldEqual, 200)

			m.Del(234)
			c.So(m.Size(), convey.ShouldEqual, 1)
		})

		c.Convey("new map by map", func(c convey.C) {
			m := NewMap(map[int]int{
				1000: 1234,
				2000: 1234,
			})

			c.So(m.Size(), convey.ShouldEqual, 2)
			c.So(m.Map(), convey.ShouldResemble, map[int]int{
				1000: 1234,
				2000: 1234,
			})

			c.Convey("common", func(c convey.C) {
				v, loaded := m.Load(123)
				c.So(loaded, convey.ShouldBeFalse)
				c.So(v, convey.ShouldEqual, 0)
				c.So(m.Has(123), convey.ShouldBeFalse)

				v, loaded = m.Swap(123, 100)
				c.So(loaded, convey.ShouldBeFalse)
				c.So(v, convey.ShouldEqual, 0)
				c.So(m.Has(123), convey.ShouldBeTrue)

				c.So(m.Size(), convey.ShouldEqual, 3)
				c.So(m.Get(123), convey.ShouldEqual, 100)

				v, loaded = m.Swap(123, 200)
				c.So(loaded, convey.ShouldBeTrue)
				c.So(v, convey.ShouldEqual, 100)

				c.So(m.Size(), convey.ShouldEqual, 3)
				c.So(m.Get(123), convey.ShouldEqual, 200)

				m.Del(234)
				c.So(m.Size(), convey.ShouldEqual, 3)

				m.Set(234, 100)
				c.So(m.Size(), convey.ShouldEqual, 4)
				c.So(m.Get(234), convey.ShouldEqual, 100)

				m.Del(234)
				c.So(m.Size(), convey.ShouldEqual, 3)

				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					1000: 1234,
					2000: 1234,
					123:  200,
				})

				keys := m.Keys()
				sort.Ints(keys)
				c.So(keys, convey.ShouldResemble, []int{
					123, 1000, 2000,
				})
			})

			c.Convey("load and deleted", func(c convey.C) {
				val, deleted := m.LoadAndDelete(1000)
				c.So(deleted, convey.ShouldBeTrue)
				c.So(val, convey.ShouldEqual, 1234)
				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					2000: 1234,
				})
			})

			c.Convey("load or store", func(c convey.C) {
				m.Set(123, 200)
				v, loaded := m.LoadOrStore(123, 100)
				c.So(loaded, convey.ShouldBeTrue)
				c.So(v, convey.ShouldEqual, 200)
				c.So(m.Size(), convey.ShouldEqual, 3)
				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					1000: 1234,
					2000: 1234,
					123:  200,
				})

				v, loaded = m.LoadOrStore(234, 100)
				c.So(loaded, convey.ShouldBeFalse)
				c.So(v, convey.ShouldEqual, 100)
				c.So(m.Size(), convey.ShouldEqual, 4)
				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					1000: 1234,
					2000: 1234,
					123:  200,
					234:  100,
				})
			})

			c.Convey("compare and swap", func(c convey.C) {
				swapped := m.CompareAndSwap(100, 200, 300)
				c.So(swapped, convey.ShouldBeFalse)
				c.So(m.Size(), convey.ShouldEqual, 2)

				swapped = m.CompareAndSwap(100, 0, 300)
				c.So(swapped, convey.ShouldBeTrue)
				c.So(m.Size(), convey.ShouldEqual, 3)
				c.So(m.Get(100), convey.ShouldEqual, 300)

				swapped = m.CompareAndSwap(1000, 1234, 123)
				c.So(swapped, convey.ShouldBeTrue)
				c.So(m.Size(), convey.ShouldEqual, 3)
				c.So(m.Get(1000), convey.ShouldEqual, 123)

				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					1000: 123,
					2000: 1234,
					100:  300,
				})
			})

			c.Convey("compare and delete", func(c convey.C) {
				deleted := m.CompareAndDelete(1123, 123)
				c.So(deleted, convey.ShouldBeFalse)
				c.So(m.Size(), convey.ShouldEqual, 2)

				deleted = m.CompareAndDelete(1000, 1234)
				c.So(deleted, convey.ShouldBeTrue)
				c.So(m.Size(), convey.ShouldEqual, 1)

				c.So(m.Map(), convey.ShouldResemble, map[int]int{
					2000: 1234,
				})
			})
		})
	})
}
