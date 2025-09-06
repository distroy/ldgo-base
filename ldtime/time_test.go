/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

var testTz = fixedTimezone("Asia/Bejing", +(8 * time.Hour))

func paresTimeStr(s string) time.Time {
	const timeFormat = "2006-01-02T15:04:05-0700"
	tm, _ := time.Parse(timeFormat, s)
	return tm
}

func fixedTimezone(name string, tm time.Duration) *time.Location {
	sec := int(tm / time.Second)
	return time.FixedZone(name, sec)
}

func TestTimeToDateNum(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("2020-12-24 16:59:54", func(c convey.C) {
			t := paresTimeStr("2020-12-24T16:59:54+0800")
			n := TimeToDateNum(t)
			c.So(n, convey.ShouldEqual, int64(20201224))
		})
	})
}

func TestDateNumToTime(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("20201224 +0800", func(c convey.C) {
			n := int64(20201224)
			t := DateNumToTime(n, testTz)
			c.So(t, convey.ShouldResemble, paresTimeStr("2020-12-24T00:00:00+0800").In(testTz))
		})
	})
}

func TestTimeToDateStr(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		n := int64(20201224)
		t := DateNumToTime(n, testTz)
		c.So(TimeToDateStr(t), convey.ShouldEqual, "2020-12-24")
	})
}

func TestDateStrToTime(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("2020-12-24", func(c convey.C) {
			s := "2020-12-24"
			t, err := DateStrToTime(s)
			c.So(err, convey.ShouldBeNil)
			c.So(t, convey.ShouldResemble, DateNumToTime(20201224))
		})
		c.Convey("2020-12-24 (UTC)", func(c convey.C) {
			s := "2020-12-24"
			t, err := DateStrToTime(s, time.UTC)
			c.So(err, convey.ShouldBeNil)
			c.So(t, convey.ShouldResemble, DateNumToTime(20201224, time.UTC))
		})
	})
}

func TestTimeToNum(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("2020-12-24T16:59:54+0800", func(c convey.C) {
			t := paresTimeStr("2020-12-24T16:59:54+0800")
			n := TimeToNum(t)
			c.So(n, convey.ShouldEqual, int64(20201224165954))
		})
	})
}

func TestNumToTime(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("20201224165954", func(c convey.C) {
			n := int64(20201224165954)
			t := NumToTime(n, testTz)
			c.So(t, convey.ShouldResemble, paresTimeStr("2020-12-24T16:59:54+0800").In(testTz))
		})
	})
}
