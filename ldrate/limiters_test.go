/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"log"
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldasync"
	"github.com/distroy/ldgo-base/ldctx"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func TestLimiters_Wait(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		gp := &ldasync.GoPool{}

		var (
			interval0 = time.Millisecond * 500
			interval1 = time.Millisecond * 700
			interval2 = time.Millisecond * 1000
		)

		var (
			l0 = NewLimiter(Every(interval0), 1)
			l1 = NewLimiter(Every(interval1), 1)
			l2 = NewLimiter(Every(interval2), 1)
		)

		c.So(l0.Allow(), convey.ShouldBeTrue)
		c.So(l1.Allow(), convey.ShouldBeTrue)
		c.So(l2.Allow(), convey.ShouldBeTrue)
		c.So(l0.Allow(), convey.ShouldBeFalse)
		c.So(l1.Allow(), convey.ShouldBeFalse)
		c.So(l2.Allow(), convey.ShouldBeFalse)

		var (
			ctx   = ldctx.Default()
			begin = time.Now()
			sleep = time.Millisecond * 10
		)

		var (
			l01 = NewLimiters(l0, l1)
			l02 = NewLimiters(l0, l2)
			l12 = NewLimiters(l1, l2)
		)

		gp.Go(func() error {
			err := l01.Wait(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(interval1))
			return nil
		})

		time.Sleep(sleep)
		gp.Go(func() error {
			err := l12.Wait(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(interval1*2))
			return nil
		})

		time.Sleep(sleep)
		gp.Go(func() error {
			err := l02.Wait(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(interval2*2))
			return nil
		})

		gp.Wait()
	})
}

func TestLimiters_Allow(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		var (
			interval0 = time.Millisecond * 500
			interval1 = time.Millisecond * 700
			interval2 = time.Millisecond * 1000
		)

		var (
			// now = time.Now()
			l0 = NewLimiter(Every(interval0), 1)
			l1 = NewLimiter(Every(interval1), 1)
			l2 = NewLimiter(Every(interval2), 1)
		)

		var (
			l01 = NewLimiters(l0, l1)
			l02 = NewLimiters(l0, l2)
			l12 = NewLimiters(l1, l2)
		)

		time.Sleep(700 * time.Millisecond)
		c.So(l01.Allow(), convey.ShouldBeTrue)
		c.So(l01.Allow(), convey.ShouldBeFalse)
		c.So(l02.Allow(), convey.ShouldBeFalse)
		c.So(l12.Allow(), convey.ShouldBeFalse)

		time.Sleep(500 * time.Millisecond)
		c.So(l02.Allow(), convey.ShouldBeTrue)
		c.So(l01.Allow(), convey.ShouldBeFalse)
		c.So(l02.Allow(), convey.ShouldBeFalse)
		c.So(l12.Allow(), convey.ShouldBeFalse)

		time.Sleep(1000 * time.Millisecond)
		c.So(l12.Allow(), convey.ShouldBeTrue)
		c.So(l01.Allow(), convey.ShouldBeFalse)
		c.So(l02.Allow(), convey.ShouldBeFalse)
		c.So(l12.Allow(), convey.ShouldBeFalse)
	})
}
