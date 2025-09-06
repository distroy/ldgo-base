/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestWaitGroup(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		wg := &WaitGroup{}
		sleepTime := time.Millisecond * 10

		var t0 time.Time
		wg.Add(1)
		go func() {
			time.Sleep(sleepTime)
			t0 = time.Now()
			wg.Done()
		}()

		c.So(wg.Count(), convey.ShouldEqual, 1)

		wg.Wait()
		t1 := time.Now()
		c.So(t1, convey.ShouldHappenAfter, t0)
		wg.Wait()
	})
}
