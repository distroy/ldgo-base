/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestBeginAndEnd(t *testing.T) {

	convey.Convey(t.Name(), t, func() {
		t0 := time.Unix(1640125769, 0).In(testTz)
		t1 := t0.In(time.UTC)

		convey.So(TimeToStr(t0), convey.ShouldResemble, "2021-12-22T06:29:29+0800")
		convey.So(TimeToStr(t1), convey.ShouldResemble, "2021-12-21T22:29:29+0000")

		convey.Convey("minute", func() {
			convey.So(TimeToStr(MinuteBegin(t0)), convey.ShouldResemble, "2021-12-22T06:29:00+0800")
			convey.So(TimeToStr(MinuteEnd(t0)), convey.ShouldResemble, "2021-12-22T06:29:59+0800")

			convey.So(TimeToStr(MinuteBegin(t1)), convey.ShouldResemble, "2021-12-21T22:29:00+0000")
			convey.So(TimeToStr(MinuteEnd(t1)), convey.ShouldResemble, "2021-12-21T22:29:59+0000")
		})

		convey.Convey("hour", func() {
			convey.So(TimeToStr(HourBegin(t0)), convey.ShouldResemble, "2021-12-22T06:00:00+0800")
			convey.So(TimeToStr(HourEnd(t0)), convey.ShouldResemble, "2021-12-22T06:59:59+0800")

			convey.So(TimeToStr(HourBegin(t1)), convey.ShouldResemble, "2021-12-21T22:00:00+0000")
			convey.So(TimeToStr(HourEnd(t1)), convey.ShouldResemble, "2021-12-21T22:59:59+0000")
		})

		convey.Convey("date", func() {
			convey.So(TimeToStr(DateBegin(t0)), convey.ShouldResemble, "2021-12-22T00:00:00+0800")
			convey.So(TimeToStr(DateEnd(t0)), convey.ShouldResemble, "2021-12-22T23:59:59+0800")

			convey.So(TimeToStr(DateBegin(t1)), convey.ShouldResemble, "2021-12-21T00:00:00+0000")
			convey.So(TimeToStr(DateEnd(t1)), convey.ShouldResemble, "2021-12-21T23:59:59+0000")
		})

		convey.Convey("week", func() {
			convey.So(TimeToStr(WeekBegin(t0)), convey.ShouldResemble, "2021-12-19T00:00:00+0800")
			convey.So(TimeToStr(WeekEnd(t0)), convey.ShouldResemble, "2021-12-25T23:59:59+0800")

			convey.So(TimeToStr(WeekBegin(t1)), convey.ShouldResemble, "2021-12-19T00:00:00+0000")
			convey.So(TimeToStr(WeekEnd(t1)), convey.ShouldResemble, "2021-12-25T23:59:59+0000")
		})

		convey.Convey("month", func() {
			convey.So(TimeToStr(MonthBegin(t0)), convey.ShouldResemble, "2021-12-01T00:00:00+0800")
			convey.So(TimeToStr(MonthEnd(t0)), convey.ShouldResemble, "2021-12-31T23:59:59+0800")

			convey.So(TimeToStr(MonthBegin(t1)), convey.ShouldResemble, "2021-12-01T00:00:00+0000")
			convey.So(TimeToStr(MonthEnd(t1)), convey.ShouldResemble, "2021-12-31T23:59:59+0000")
		})

		convey.Convey("year", func() {
			convey.So(TimeToStr(YearBegin(t0)), convey.ShouldResemble, "2021-01-01T00:00:00+0800")
			convey.So(TimeToStr(YearEnd(t0)), convey.ShouldResemble, "2021-12-31T23:59:59+0800")

			convey.So(TimeToStr(YearBegin(t1)), convey.ShouldResemble, "2021-01-01T00:00:00+0000")
			convey.So(TimeToStr(YearEnd(t1)), convey.ShouldResemble, "2021-12-31T23:59:59+0000")
		})
	})
}
