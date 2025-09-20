/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldasync"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
)

func TestLimiter_Wait(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		ctx := ldctx.Default()
		gp := &ldasync.GoPool{}

		begin := time.Now()
		interval := time.Second

		l := NewLimiter(Every(interval), 1)
		// l.SetBurst(1)
		// l.SetLimit(Every(interval))

		c.Convey("wait times", func() {
			sleep := interval / 10
			for i := range 3 {
				n := time.Duration(i) * interval
				gp.Go(func() error {
					err := l.Wait(ctx)

					c.So(err, convey.ShouldBeNil)
					c.So(time.Now(), convey.ShouldHappenOnOrAfter, begin.Add(n))
					return nil
				})
				time.Sleep(sleep)
			}

			gp.Wait()
		})

		c.Convey("context has cancelled", func() {
			ctx, cancel := ldctx.WithCancel(ctx)
			cancel()

			err := l.Wait(ctx)

			c.So(err, convey.ShouldResemble, lderr.ErrCtxCanceled)
		})

		c.Convey("deedline not enough", func() {
			ctx, _ := ldctx.WithTimeout(ctx, interval/2)

			err := l.Wait(ctx)
			c.So(err, convey.ShouldBeNil)

			err = l.Wait(ctx)
			c.So(err, convey.ShouldResemble, lderr.ErrCtxDeadlineNotEnough)
		})

		// c.Convey("no wait time", func() {
		// 	l.refresh(ctx, begin.Add(-interval))
		// 	// time.Sleep(interval)
		//
		// 	err := l.Wait(ctx)
		// 	end := time.Now()
		// 	c.So(err, convey.ShouldBeNil)
		// 	c.So(end, convey.ShouldHappenBefore, begin.Add(interval))
		// 	c.So(end, convey.ShouldHappenBefore, begin.Add(1*time.Millisecond))
		// })
	})
}

func TestLimiter_Allow(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		ctx := ldctx.Default()
		ctx, cancel := ldctx.WithCancel(ctx)
		cancel()

		interval := time.Second

		l := NewLimiter(1, 1)
		l.SetBurst(1)
		l.SetLimit(Every(interval))

		c.So(l.Allow(), convey.ShouldBeTrue)
		c.So(l.Allow(), convey.ShouldBeFalse)

		time.Sleep(interval * 2)
		c.So(l.Allow(), convey.ShouldBeTrue)
		c.So(l.Allow(), convey.ShouldBeFalse)

		time.Sleep(interval)
		c.So(l.Allow(), convey.ShouldBeTrue)
		c.So(l.Allow(), convey.ShouldBeFalse)

		time.Sleep(interval)
		c.So(l.Allow(), convey.ShouldBeTrue)
		c.So(l.Allow(), convey.ShouldBeFalse)
	})
}
