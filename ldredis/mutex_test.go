/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldasync"
	"github.com/distroy/ldgo-base/ldctx"
)

func TestMutex_Lock(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		rr := testMutexRedis()
		r := WrapMutexRedisWithCtx(rr)
		defer rr.Close()

		lockKey := "test-key"
		timeout := 1 * time.Second

		ctx, _ := ldctx.WithTimeout(ldctx.Discard(), timeout)

		m0 := NewMutex(r)
		m1 := NewMutex(r)

		convey.So(m0.Lock(ctx, lockKey), convey.ShouldBeNil)
		convey.So(m1.Lock(ctx, lockKey), convey.ShouldNotBeNil)
		convey.So(m0.Unlock(ctx), convey.ShouldBeNil)
	})
}

func TestMutex_Unlock(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		rr := testMutexRedis()
		r := WrapMutexRedisWithCtx(rr)
		defer rr.Close()

		lockKey := "test-key"
		timeout := 1 * time.Second
		ctx, _ := ldctx.WithTimeout(ldctx.Discard(), timeout)

		convey.Convey("unlock after context not timeout", func() {
			m := NewMutex(r)

			convey.So(m.Lock(ctx, lockKey), convey.ShouldBeNil)

			// time.Sleep(timeout - 1*time.Second)

			convey.So(m.Unlock(ctx), convey.ShouldBeNil)
		})

		convey.Convey("unlock after context timeout", func() {
			m := NewMutex(r)

			convey.So(m.Lock(ctx, lockKey), convey.ShouldBeNil)

			<-m.Events()
			// time.Sleep(timeout + 1*time.Second)

			convey.So(m.Unlock(ctx), convey.ShouldBeNil)
		})
	})
}

func TestMutex_WithLockForce(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		rr := testMutexRedis()
		r := WrapMutexRedisWithCtx(rr)
		defer rr.Close()

		ctx := ldctx.Discard()

		lockKey := "test-key"
		interval := 100 * time.Millisecond
		timeout0 := 2 * time.Second
		timeout1 := 1 * time.Second

		convey.Convey("lock force without timeout", func(c convey.C) {
			var t0, t1 time.Time
			m0 := NewMutex(r)
			m0 = m0.WithInterval(time.Second)

			m1 := NewMutex(r)
			m1 = m1.WithLockForce(true, interval)

			gp := &ldasync.GoPool{}

			c.So(m0.Lock(ctx, lockKey), convey.ShouldBeNil)
			gp.Go(func() error {
				time.Sleep(timeout0)
				c.So(m0.Unlock(ctx), convey.ShouldBeNil)

				t0 = time.Now()
				return nil
			})

			gp.Go(func() error {
				c.So(m1.Lock(ctx, lockKey), convey.ShouldBeNil)
				t1 = time.Now()

				// time.Sleep(timeout)
				c.So(m1.Unlock(ctx), convey.ShouldBeNil)
				return nil
			})

			gp.Wait()
			c.So(t0, convey.ShouldHappenBefore, t1)
		})

		convey.Convey("lock force with timeout", func(c convey.C) {
			convey.Convey("lock succ", func(c convey.C) {
				var t0, t1 time.Time
				m0 := NewMutex(r)

				m1 := NewMutex(r)
				m1 = m1.WithLockForce(true, interval, timeout1)

				gos := &ldasync.GoPool{}

				c.So(m0.Lock(ctx, lockKey), convey.ShouldBeNil)
				gos.Go(func() error {

					time.Sleep(timeout0)
					c.So(m0.Unlock(ctx), convey.ShouldBeNil)

					t0 = time.Now()
					return nil
				})

				gos.Go(func() error {
					m := m1
					c.So(m.Lock(ctx, lockKey), convey.ShouldNotBeNil)
					t1 = time.Now()

					// c.So(m.Unlock(), convey.ShouldNotBeNil)
					return nil
				})

				gos.Wait()
				c.So(t0, convey.ShouldHappenAfter, t1)
			})
			convey.Convey("lock timeout", func(c convey.C) {
				var t0, t1 time.Time
				m0 := NewMutex(r)

				m1 := NewMutex(r)
				m1 = m1.WithLockForce(true, interval, timeout0-time.Second)

				gos := &ldasync.GoPool{}

				c.So(m0.Lock(ctx, lockKey), convey.ShouldBeNil)
				gos.Go(func() error {
					time.Sleep(timeout0)
					c.So(m0.Unlock(ctx), convey.ShouldBeNil)

					t0 = time.Now()
					return nil
				})

				gos.Go(func() error {
					m := m1
					c.So(m.Lock(ctx, lockKey), convey.ShouldNotBeNil)
					t1 = time.Now()

					// c.So(m.Unlock(), convey.ShouldNotBeNil)
					return nil
				})

				gos.Wait()
				c.So(t0, convey.ShouldHappenAfter, t1)
			})
		})
	})
}
