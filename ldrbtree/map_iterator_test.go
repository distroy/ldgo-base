/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestMap_Iterator(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		m := testNewMap()

		convey.Convey("next", func() {
			it := m.Begin()
			for _, n := range _nums {
				convey.So(it.Key(), convey.ShouldEqual, n)
				it = it.Next()
			}
			convey.So(it, convey.ShouldResemble, m.End())
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("prev", func() {
			i := len(_nums)
			for begin, it := m.Begin(), m.End(); it != begin; {
				it = it.Prev()
				i--
				n := _nums[i]
				convey.So(it.Key(), convey.ShouldEqual, n)
			}
		})

		convey.Convey("delete ordered", func() {
			for it, end := m.Begin(), m.End(); it != end; {
				it = m.Delete(it)
				convey.So(m.tree.root.checkParent(m.tree.sentinel), convey.ShouldBeTrue)
				convey.So(m.tree.root.checkColor(m.tree.sentinel), convey.ShouldBeTrue)
			}

			convey.So(m.Begin(), convey.ShouldResemble, m.End())
			convey.So(m.Len(), convey.ShouldEqual, 0)
		})

		convey.Convey("delete unordered", func() {
			for _, n := range _numsUnordered {
				it := m.Search(n)
				convey.So(it.Key(), convey.ShouldResemble, n)

				m.Delete(it)
				convey.So(m.tree.root.checkParent(m.tree.sentinel), convey.ShouldBeTrue)
				convey.So(m.tree.root.checkColor(m.tree.sentinel), convey.ShouldBeTrue)

				// it = m.Search(n)
				// convey.So(it, convey.ShouldResemble, m.End())
			}

			convey.So(m.Begin(), convey.ShouldResemble, m.End())
			convey.So(m.Len(), convey.ShouldEqual, 0)
		})

	})
}

func TestMap_ReverseIterator(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		m := testNewMap()

		convey.Convey("next", func() {
			it := m.RBegin()
			for i := len(_nums) - 1; i >= 0; i-- {
				n := _nums[i]
				convey.So(it.Key(), convey.ShouldEqual, n)
				it = it.Next()
			}
			convey.So(it, convey.ShouldResemble, m.REnd())
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("prev", func() {
			i := 0
			for begin, it := m.RBegin(), m.REnd(); it != begin; {
				it = it.Prev()
				n := _nums[i]
				i++
				convey.So(it.Key(), convey.ShouldEqual, n)
			}
		})

		convey.Convey("delete ordered", func() {
			for it, end := m.RBegin(), m.REnd(); it != end; {
				it = m.RDelete(it)
				convey.So(m.tree.root.checkParent(m.tree.sentinel), convey.ShouldBeTrue)
				convey.So(m.tree.root.checkColor(m.tree.sentinel), convey.ShouldBeTrue)
			}

			convey.So(m.Begin(), convey.ShouldResemble, m.End())
			convey.So(m.Len(), convey.ShouldEqual, 0)
		})

		convey.Convey("delete unordered", func() {
			for _, n := range _numsUnordered {
				it := m.RSearch(n)
				convey.So(it.Key(), convey.ShouldResemble, n)

				m.RDelete(it)
				convey.So(m.tree.root.checkParent(m.tree.sentinel), convey.ShouldBeTrue)
				convey.So(m.tree.root.checkColor(m.tree.sentinel), convey.ShouldBeTrue)

				// it = rbtree.Search(n)
				// convey.So(it, convey.ShouldResemble, rbtree.End())
			}

			convey.So(m.Begin(), convey.ShouldResemble, m.End())
			convey.So(m.Len(), convey.ShouldEqual, 0)
		})
	})
}
