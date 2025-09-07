/*
 * Copyright (C) distroy
 */

package ldsort

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestSortInts(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []int{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortInts(l)

		convey.So(IsSortedInts(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []int{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchInts(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchInts(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchInts(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchInts(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchInts(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqInts(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqInts(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []int{123}

			SortInts(l)
			convey.So(l, convey.ShouldResemble, []int{123})

			l = UniqInts(l)
			convey.So(l, convey.ShouldResemble, []int{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []int{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortInts(l)
			convey.So(l, convey.ShouldResemble, []int{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqInts(l)
			convey.So(l, convey.ShouldResemble, []int{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexInts(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexInts(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInts(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexInts(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInts(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortInt8s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []int8{13, 125, 57, 83, 73, 111, 13, 41, 7, 23}

		SortInt8s(l)

		convey.So(IsSortedInt8s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []int8{
			7, 13, 13, 23, 41, 57, 73, 83, 111, 125,
		})

		convey.So(SearchInt8s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchInt8s(l, 7), convey.ShouldEqual, 0)
		convey.So(SearchInt8s(l, 10), convey.ShouldEqual, 1)
		convey.So(SearchInt8s(l, 23), convey.ShouldEqual, 3)
		convey.So(SearchInt8s(l, 100), convey.ShouldEqual, 8)
		convey.So(SearchInt8s(l, 127), convey.ShouldEqual, 10)
	})
}

func TestUniqInt8s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqInt8s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []int8{123}

			SortInt8s(l)
			convey.So(l, convey.ShouldResemble, []int8{123})

			l = UniqInt8s(l)
			convey.So(l, convey.ShouldResemble, []int8{123})
		})

		convey.Convey("[23, 75, 24, 123, 42, 23, 123, 25, 23, 42, 25, 75]", func() {
			l := []int8{23, 75, 24, 123, 42, 23, 123, 25, 23, 42, 25, 75}

			SortInt8s(l)
			convey.So(l, convey.ShouldResemble, []int8{
				23, 23, 23, 24, 25, 25, 42, 42, 75, 75, 123, 123,
			})

			l = UniqInt8s(l)
			convey.So(l, convey.ShouldResemble, []int8{
				23, 24, 25, 42, 75, 123,
			})

			convey.So(IndexInt8s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexInt8s(l, 23), convey.ShouldEqual, 0)
			convey.So(IndexInt8s(l, 42), convey.ShouldEqual, 3)
			convey.So(IndexInt8s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt8s(l, 123), convey.ShouldEqual, 5)
			convey.So(IndexInt8s(l, 125), convey.ShouldEqual, -1)
		})
	})
}

func TestSortInt16s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []int16{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortInt16s(l)

		convey.So(IsSortedInt16s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []int16{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchInt16s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchInt16s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchInt16s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchInt16s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchInt16s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqInt16s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqInt16s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []int16{123}

			SortInt16s(l)
			convey.So(l, convey.ShouldResemble, []int16{123})

			l = UniqInt16s(l)
			convey.So(l, convey.ShouldResemble, []int16{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []int16{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortInt16s(l)
			convey.So(l, convey.ShouldResemble, []int16{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqInt16s(l)
			convey.So(l, convey.ShouldResemble, []int16{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexInt16s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexInt16s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt16s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexInt16s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt16s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortInt32s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []int32{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortInt32s(l)

		convey.So(IsSortedInt32s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []int32{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchInt32s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchInt32s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchInt32s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchInt32s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchInt32s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqInt32s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqInt32s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []int32{123}

			SortInt32s(l)
			convey.So(l, convey.ShouldResemble, []int32{123})

			l = UniqInt32s(l)
			convey.So(l, convey.ShouldResemble, []int32{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []int32{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortInt32s(l)
			convey.So(l, convey.ShouldResemble, []int32{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqInt32s(l)
			convey.So(l, convey.ShouldResemble, []int32{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexInt32s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexInt32s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt32s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexInt32s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt32s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortInt64s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []int64{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortInt64s(l)

		convey.So(IsSortedInt64s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []int64{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchInt64s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchInt64s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchInt64s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchInt64s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchInt64s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqInt64s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqInt64s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []int64{123}

			SortInt64s(l)
			convey.So(l, convey.ShouldResemble, []int64{123})

			l = UniqInt64s(l)
			convey.So(l, convey.ShouldResemble, []int64{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []int64{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortInt64s(l)
			convey.So(l, convey.ShouldResemble, []int64{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqInt64s(l)
			convey.So(l, convey.ShouldResemble, []int64{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexInt64s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexInt64s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt64s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexInt64s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexInt64s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortUints(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []uint{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortUints(l)

		convey.So(IsSortedUints(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []uint{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchUints(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchUints(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchUints(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchUints(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchUints(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqUints(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqUints(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []uint{123}

			SortUints(l)
			convey.So(l, convey.ShouldResemble, []uint{123})

			l = UniqUints(l)
			convey.So(l, convey.ShouldResemble, []uint{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []uint{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortUints(l)
			convey.So(l, convey.ShouldResemble, []uint{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqUints(l)
			convey.So(l, convey.ShouldResemble, []uint{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexUints(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexUints(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUints(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexUints(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUints(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortUint8s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []uint8{13, 125, 57, 83, 73, 111, 13, 41, 7, 23}

		SortUint8s(l)

		convey.So(IsSortedUint8s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []uint8{
			7, 13, 13, 23, 41, 57, 73, 83, 111, 125,
		})

		convey.So(SearchUint8s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchUint8s(l, 7), convey.ShouldEqual, 0)
		convey.So(SearchUint8s(l, 10), convey.ShouldEqual, 1)
		convey.So(SearchUint8s(l, 23), convey.ShouldEqual, 3)
		convey.So(SearchUint8s(l, 100), convey.ShouldEqual, 8)
		convey.So(SearchUint8s(l, 127), convey.ShouldEqual, 10)
	})
}

func TestUniqUint8s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqUint8s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []uint8{123}

			SortUint8s(l)
			convey.So(l, convey.ShouldResemble, []uint8{123})

			l = UniqUint8s(l)
			convey.So(l, convey.ShouldResemble, []uint8{123})
		})

		convey.Convey("[23, 75, 24, 123, 42, 23, 123, 25, 23, 42, 25, 75]", func() {
			l := []uint8{23, 75, 24, 123, 42, 23, 123, 25, 23, 42, 25, 75}

			SortUint8s(l)
			convey.So(l, convey.ShouldResemble, []uint8{
				23, 23, 23, 24, 25, 25, 42, 42, 75, 75, 123, 123,
			})

			l = UniqUint8s(l)
			convey.So(l, convey.ShouldResemble, []uint8{
				23, 24, 25, 42, 75, 123,
			})

			convey.So(IndexUint8s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexUint8s(l, 23), convey.ShouldEqual, 0)
			convey.So(IndexUint8s(l, 42), convey.ShouldEqual, 3)
			convey.So(IndexUint8s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint8s(l, 123), convey.ShouldEqual, 5)
			convey.So(IndexUint8s(l, 125), convey.ShouldEqual, -1)
		})
	})
}

func TestSortUint16s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []uint16{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortUint16s(l)

		convey.So(IsSortedUint16s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []uint16{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchUint16s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchUint16s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchUint16s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchUint16s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchUint16s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqUint16s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqUint16s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []uint16{123}

			SortUint16s(l)
			convey.So(l, convey.ShouldResemble, []uint16{123})

			l = UniqUint16s(l)
			convey.So(l, convey.ShouldResemble, []uint16{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []uint16{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortUint16s(l)
			convey.So(l, convey.ShouldResemble, []uint16{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqUint16s(l)
			convey.So(l, convey.ShouldResemble, []uint16{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexUint16s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexUint16s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint16s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexUint16s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint16s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortUint32s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []uint32{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortUint32s(l)

		convey.So(IsSortedUint32s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []uint32{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchUint32s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchUint32s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchUint32s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchUint32s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchUint32s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqUint32s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqUint32s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []uint32{123}

			SortUint32s(l)
			convey.So(l, convey.ShouldResemble, []uint32{123})

			l = UniqUint32s(l)
			convey.So(l, convey.ShouldResemble, []uint32{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []uint32{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortUint32s(l)
			convey.So(l, convey.ShouldResemble, []uint32{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqUint32s(l)
			convey.So(l, convey.ShouldResemble, []uint32{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexUint32s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexUint32s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint32s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexUint32s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint32s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortUint64s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []uint64{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortUint64s(l)

		convey.So(IsSortedUint64s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []uint64{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchUint64s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchUint64s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchUint64s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchUint64s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchUint64s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqUint64s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqUint64s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []uint64{123}

			SortUint64s(l)
			convey.So(l, convey.ShouldResemble, []uint64{123})

			l = UniqUint64s(l)
			convey.So(l, convey.ShouldResemble, []uint64{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []uint64{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortUint64s(l)
			convey.So(l, convey.ShouldResemble, []uint64{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqUint64s(l)
			convey.So(l, convey.ShouldResemble, []uint64{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexUint64s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexUint64s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint64s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexUint64s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUint64s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortUintptrs(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []uintptr{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortUintptrs(l)

		convey.So(IsSortedUintptrs(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []uintptr{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchUintptrs(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchUintptrs(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchUintptrs(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchUintptrs(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchUintptrs(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqUintptrs(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqUintptrs(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []uintptr{123}

			SortUintptrs(l)
			convey.So(l, convey.ShouldResemble, []uintptr{123})

			l = UniqUintptrs(l)
			convey.So(l, convey.ShouldResemble, []uintptr{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []uintptr{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortUintptrs(l)
			convey.So(l, convey.ShouldResemble, []uintptr{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqUintptrs(l)
			convey.So(l, convey.ShouldResemble, []uintptr{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexUintptrs(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexUintptrs(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUintptrs(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexUintptrs(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexUintptrs(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortFloat32s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []float32{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortFloat32s(l)

		convey.So(IsSortedFloat32s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []float32{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchFloat32s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchFloat32s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchFloat32s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchFloat32s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchFloat32s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqFloat32s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqFloat32s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []float32{123}

			SortFloat32s(l)
			convey.So(l, convey.ShouldResemble, []float32{123})

			l = UniqFloat32s(l)
			convey.So(l, convey.ShouldResemble, []float32{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []float32{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortFloat32s(l)
			convey.So(l, convey.ShouldResemble, []float32{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqFloat32s(l)
			convey.So(l, convey.ShouldResemble, []float32{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexFloat32s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexFloat32s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexFloat32s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexFloat32s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexFloat32s(l, 223), convey.ShouldEqual, 1)
		})
	})
}

func TestSortFloat64s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		l := []float64{223, 562, 424, 642, 223, 123, 496, 623, 845, 375}

		SortFloat64s(l)

		convey.So(IsSortedFloat64s(l), convey.ShouldBeTrue)
		convey.So(l, convey.ShouldResemble, []float64{
			123, 223, 223, 375, 424, 496, 562, 623, 642, 845,
		})

		convey.So(SearchFloat64s(l, 0), convey.ShouldEqual, 0)
		convey.So(SearchFloat64s(l, 123), convey.ShouldEqual, 0)
		convey.So(SearchFloat64s(l, 223), convey.ShouldEqual, 1)
		convey.So(SearchFloat64s(l, 300), convey.ShouldEqual, 3)
		convey.So(SearchFloat64s(l, 10000), convey.ShouldEqual, 10)
	})
}

func TestUniqFloat64s(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("nil", func() {
			l := UniqFloat64s(nil)
			convey.So(l, convey.ShouldBeNil)
		})

		convey.Convey("[123]", func() {
			l := []float64{123}

			SortFloat64s(l)
			convey.So(l, convey.ShouldResemble, []float64{123})

			l = UniqFloat64s(l)
			convey.So(l, convey.ShouldResemble, []float64{123})
		})

		convey.Convey("[223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375]", func() {
			l := []float64{223, 375, 424, 123, 642, 223, 123, 725, 223, 642, 725, 375}

			SortFloat64s(l)
			convey.So(l, convey.ShouldResemble, []float64{
				123, 123, 223, 223, 223, 375, 375, 424, 642, 642, 725, 725,
			})

			l = UniqFloat64s(l)
			convey.So(l, convey.ShouldResemble, []float64{
				123, 223, 375, 424, 642, 725,
			})

			convey.So(IndexFloat64s(l, 0), convey.ShouldEqual, -1)
			convey.So(IndexFloat64s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexFloat64s(l, 123), convey.ShouldEqual, 0)
			convey.So(IndexFloat64s(l, 100), convey.ShouldEqual, -1)
			convey.So(IndexFloat64s(l, 223), convey.ShouldEqual, 1)
		})
	})
}
