/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

var _testTimezone = testFixedTimezone("Asia/Bejing", +(8 * time.Hour))

func testFixedTimezone(name string, tm time.Duration) *time.Location {
	sec := int(tm / time.Second)
	return time.FixedZone(name, sec)
}

func TestTime(t *testing.T) {
	const format = "2006-01-02T15:04:05-0700"
	now := time.Unix(1644479966, 0).In(_testTimezone)
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(NewTime(time.Time{}).Load().Format(format), convey.ShouldEqual, _zeroTime.Format(format))

		p := Time{}
		c.So(p.Load().Format(format), convey.ShouldEqual, _zeroTime.Format(format))

		p.Store(now)
		c.So(p.Load().Format(format), convey.ShouldEqual, "2022-02-10T15:59:26+0800")

		// p.MustAdd(60 * time.Second)
		p.MustChange(func(old time.Time) (new time.Time) { return old.Add(60 * time.Second) })
		c.So(p.Load().Format(format), convey.ShouldEqual, "2022-02-10T16:00:26+0800")

		// p.MustAddDate(-1, 1, 3)
		p.MustChange(func(old time.Time) (new time.Time) { return old.AddDate(-1, 1, 3) })
		c.So(p.Load().Format(format), convey.ShouldEqual, "2021-03-13T16:00:26+0800")

		c.So(p.Swap(now).Format(format), convey.ShouldEqual, "2021-03-13T16:00:26+0800")
		c.So(p.Load().Format(format), convey.ShouldEqual, "2022-02-10T15:59:26+0800")

		c.So(p.CompareAndSwap(time.Time{}, now), convey.ShouldBeFalse)
		c.So(p.CompareAndSwap(now, time.Time{}), convey.ShouldBeTrue)
		c.So(p.CompareAndSwap(time.Time{}, now), convey.ShouldBeTrue)
	})
}
