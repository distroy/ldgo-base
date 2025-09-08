/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestRBTree_Iterator(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		rbtree := testNewRBTree()

		convey.Convey("next", func() {
			it := rbtree.Begin()
			for _, n := range _nums {
				convey.So(it.Get(), convey.ShouldEqual, n)
				it = it.Next()
			}
			convey.So(it, convey.ShouldResemble, rbtree.End())
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("prev", func() {
			i := len(_nums)
			for begin, it := rbtree.Begin(), rbtree.End(); it != begin; {
				it = it.Prev()
				i--
				n := _nums[i]
				convey.So(it.Get(), convey.ShouldEqual, n)
			}
		})

		convey.Convey("delete ordered", func() {
			for it, end := rbtree.Begin(), rbtree.End(); it != end; {
				it = rbtree.Delete(it)
				convey.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
				convey.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)
			}

			convey.So(rbtree.Begin(), convey.ShouldResemble, rbtree.End())
			convey.So(rbtree.Len(), convey.ShouldEqual, 0)
		})

		convey.Convey("delete unordered", func() {
			for _, n := range _numsUnordered {
				it := rbtree.Search(n)
				convey.So(it.Get(), convey.ShouldResemble, n)

				rbtree.Delete(it)
				convey.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
				convey.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)

				// it = rbtree.Search(n)
				// convey.So(it, convey.ShouldResemble, rbtree.End())
			}

			convey.So(rbtree.Begin(), convey.ShouldResemble, rbtree.End())
			convey.So(rbtree.Len(), convey.ShouldEqual, 0)
		})

	})
}

func TestRBTree_ReverseIterator(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		rbtree := testNewRBTree()

		convey.Convey("next", func() {
			it := rbtree.RBegin()
			for i := len(_nums) - 1; i >= 0; i-- {
				n := _nums[i]
				convey.So(it.Get(), convey.ShouldEqual, n)
				it = it.Next()
			}
			convey.So(it, convey.ShouldResemble, rbtree.REnd())
			convey.So(func() { it.Next() }, convey.ShouldPanic)
		})

		convey.Convey("prev", func() {
			i := 0
			for begin, it := rbtree.RBegin(), rbtree.REnd(); it != begin; {
				it = it.Prev()
				n := _nums[i]
				i++
				convey.So(it.Get(), convey.ShouldEqual, n)
			}
		})

		convey.Convey("delete ordered", func() {
			for it, end := rbtree.RBegin(), rbtree.REnd(); it != end; {
				it = rbtree.RDelete(it)
				convey.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
				convey.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)
			}

			convey.So(rbtree.Begin(), convey.ShouldResemble, rbtree.End())
			convey.So(rbtree.Len(), convey.ShouldEqual, 0)
		})

		convey.Convey("delete unordered", func() {
			for _, n := range _numsUnordered {
				it := rbtree.RSearch(n)
				convey.So(it.Get(), convey.ShouldResemble, n)

				rbtree.RDelete(it)
				convey.So(rbtree.root.checkParent(rbtree.sentinel), convey.ShouldBeTrue)
				convey.So(rbtree.root.checkColor(rbtree.sentinel), convey.ShouldBeTrue)

				// it = rbtree.Search(n)
				// convey.So(it, convey.ShouldResemble, rbtree.End())
			}

			convey.So(rbtree.Begin(), convey.ShouldResemble, rbtree.End())
			convey.So(rbtree.Len(), convey.ShouldEqual, 0)
		})
	})
}
