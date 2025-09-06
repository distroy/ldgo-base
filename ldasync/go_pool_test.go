/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestGoPool(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("nomarl", func(c convey.C) {
			n := int32(0)

			p := GoN(10, func() error {
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&n, 1)
				return nil
			})

			c.So(atomic.LoadInt32(&n), convey.ShouldEqual, 0)
			c.So(p.Count(), convey.ShouldEqual, 10)

			err := p.Wait()
			c.So(err, convey.ShouldBeNil)

			c.So(atomic.LoadInt32(&n), convey.ShouldEqual, 10)
			c.So(p.Count(), convey.ShouldEqual, 0)
		})

		c.Convey("panic", func(c convey.C) {
			fn := func() error {
				panic(11)
			}

			p := GoN(2, fn)
			err := p.Wait()
			c.So(err, convey.ShouldNotBeNil)
		})
	})
}
