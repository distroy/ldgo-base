/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"context"
	"slices"
	"time"

	"github.com/distroy/ldgo-base/ldptr"
)

func NewLimiters(limiters ...LimiterIface) *Limiters {
	l := &Limiters{}
	l.AddLimiter(limiters...)
	return l
}

type Limiters []LimiterIface

func (l *Limiters) AddLimiter(limiters ...LimiterIface) *Limiters {
	count := 0
	for _, limiter := range limiters {
		n := 1
		if v, ok := limiter.(*Limiters); ok {
			n = len(*v)
		}
		count += n
	}

	buf := *l
	// buf = make([]LimiterIface, 0, len(buf)+count)
	// buf = append(buf, l.limiters...)
	buf = slices.Grow(buf, count)

	for _, limiter := range limiters {
		if v, ok := limiter.(*Limiters); ok {
			buf = append(buf, *v...)
		} else {
			buf = append(buf, limiter)
		}
	}

	*l = buf
	return l
}

func (l *Limiters) Allow() bool                    { return l.AllowN(time.Now(), 1) }
func (l *Limiters) AllowN(t time.Time, n int) bool { return l.reserveWait(t, n, 0).ok }

func (l *Limiters) Wait(c context.Context) error         { return l.WaitN(c, 1) }
func (l *Limiters) WaitN(c context.Context, n int) error { return wait(c, l, n) }

func (l *Limiters) Reserve() Reservation                    { return l.ReserveN(time.Now(), 1) }
func (l *Limiters) ReserveN(t time.Time, n int) Reservation { return l.ReserveWait(t, n, infDuration) }
func (l *Limiters) ReserveWait(t time.Time, n int, wait time.Duration) Reservation {
	return ldptr.New(l.reserveWait(t, n, wait))
}

func (l *Limiters) reserveWait(t time.Time, n int, wait time.Duration) reservations {
	if len(*l) == 0 {
		return reservations{ok: true}
	}

	r := reservations{
		list: make([]Reservation, 0, len(*l)),
		ok:   true,
	}

	for _, ll := range *l {
		rr := ll.ReserveWait(t, n, wait)
		if !rr.Ok() {
			// log.Printf("ReserveN. ok:%v, time:%s, detail:%s", rr.Ok(), str(t), str(ll))
			r.ok = false
			r.cancelAt(t)
			return r
		}

		// log.Printf("ReserveN. ok:%v, time:%s, detail:%s", rr.Ok(), str(t), str(ll))
		r.list = append(r.list, rr)
	}

	// log.Printf("ReserveN finaly result. ok:%v", r.ok)
	return r
}

type reservations struct {
	list []Reservation
	ok   bool
}

func (r *reservations) Ok() bool { return r.ok }

func (r *reservations) Cancel() { r.CancelAt(time.Now()) }
func (r *reservations) CancelAt(t time.Time) {
	if !r.ok {
		return
	}
	r.cancelAt(t)
}
func (r *reservations) cancelAt(t time.Time) {
	for _, v := range r.list {
		// log.Printf("CancelAt 1111. time:%s, limiter:%s", str(t), str(v.Limiter))
		v.CancelAt(t)
		// log.Printf("CancelAt 2222. time:%s, limiter:%s", str(t), str(v.Limiter))
	}
}

func (r *reservations) Delay() time.Duration { return r.DelayFrom(time.Now()) }
func (r *reservations) DelayFrom(t time.Time) time.Duration {
	if !r.ok {
		return infDuration
	}

	delay := time.Duration(0)
	for _, v := range r.list {
		d := v.DelayFrom(t)
		delay = max(delay, d)
	}
	return delay
}
