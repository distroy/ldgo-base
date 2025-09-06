/*
 * Copyright (C) distroy
 */

package ldptr

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestGet(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(Get[time.Duration](nil), convey.ShouldEqual, time.Duration(0))
		c.So(Get[time.Duration](nil, 1), convey.ShouldEqual, time.Duration(1))
		c.So(Get[time.Duration](New[time.Duration](100), time.Duration(0)), convey.ShouldEqual, time.Duration(100))
	})
}

func TestNewByPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(NewByPtr[time.Duration](nil), convey.ShouldBeNil)
		c.So(NewByPtr[time.Duration](nil, 0), convey.ShouldResemble, New[time.Duration](0))
		c.So(NewByPtr[time.Duration](nil, 1), convey.ShouldResemble, New[time.Duration](1))
		c.So(NewByPtr[time.Duration](New[time.Duration](1)), convey.ShouldResemble, New[time.Duration](1))
		c.So(NewByPtr[time.Duration](New[time.Duration](100)), convey.ShouldResemble, New[time.Duration](100))
		c.So(NewByPtr[time.Duration](New[time.Duration](-100)), convey.ShouldResemble, New[time.Duration](-100))
	})
}
