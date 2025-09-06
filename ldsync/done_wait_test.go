/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestDoneWait(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		p := &DoneWait{}

		var t0 time.Time
		go func() {
			time.Sleep(time.Second)
			t0 = time.Now()
			p.Stop()
		}()

		for i := 0; i < 10; i++ {
			p.Add(1)
			go func() {
				defer p.Done()
				<-p.Chan()

				t := time.Now()
				c.So(t, convey.ShouldHappenAfter, t0)
			}()
		}
		c.So(p.Count(), convey.ShouldEqual, 10)

		p.Wait()
	})
}
