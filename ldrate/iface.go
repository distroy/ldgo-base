/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"context"
	"time"

	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
)

var (
	_ LimiterIface = (*Limiter)(nil)
	_ LimiterIface = (*Limiters)(nil)
)

type Reservation interface {
	Cancel()
	CancelAt(t time.Time)
	Delay() time.Duration
	DelayFrom(t time.Time) time.Duration
	Ok() bool
}

type LimiterIface interface {
	// Allow reports whether an event may happen now.
	Allow() bool
	// AllowN reports whether n events may happen at time t.
	AllowN(t time.Time, n int) bool

	// Wait is shorthand for WaitN(ctx, 1).
	Wait(c context.Context) error
	// WaitN blocks until lim permits n events to happen.
	WaitN(c context.Context, n int) error

	// Reserve is shorthand for ReserveN(time.Now(), 1).
	Reserve() Reservation
	// ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
	ReserveN(t time.Time, n int) Reservation
	// ReserveWait returns a Reservation that indicates how long the caller must wait before n events happen.
	ReserveWait(t time.Time, n int, wait time.Duration) Reservation
}

type reserver interface {
	ReserveWait(t time.Time, n int, wait time.Duration) Reservation
}

func wait(c context.Context, l reserver, n int) error {
	select {
	case <-c.Done():
		// return lderr.ErrCtxCanceled
		return ldctx.GetError(c)
	default:
	}

	now := time.Now()

	// Determine wait limit
	waitLimit := infDuration
	if deadline, ok := c.Deadline(); ok {
		waitLimit = deadline.Sub(now)
	}

	r := l.ReserveWait(now, n, waitLimit)
	if !r.Ok() {
		return lderr.ErrCtxDeadlineNotEnough
	}

	// Wait if necessary
	delay := r.DelayFrom(now)
	if delay <= 0 {
		return nil
	}

	// if delay >= waitLimit {
	// 	r.CancelAt(now)
	// 	return lderr.ErrCtxDeadlineNotEnough
	// }

	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-t.C:
		break

	case <-c.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.CancelAt(now)
		return ldctx.GetError(c)
	}

	// We can proceed.
	return nil
}

// func str(i any) string {
// 	switch v := i.(type) {
// 	case *Limiter:
// 		return fmt.Sprintf("{limit:%f, burst:%d, tokens:%f, last:%s, lastEvent:%s}",
// 			v.limit, v.burst, v.tokens, str(v.last), str(v.lastEvent))
//
// 	case *Limiters:
// 		b := make([]string, 0, len(*v))
// 		for _, v := range *v {
// 			b = append(b, str(v))
// 		}
// 		return fmt.Sprintf("[%s]", strings.Join(b, ","))
//
// 	case time.Time:
// 		return ldtime.TimeToStr(v)
// 	}
// 	return fmt.Sprint(i)
// }
