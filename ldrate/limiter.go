/*
 * Copyright (C) distroy
 */

package ldrate

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/distroy/ldgo-base/ldptr"
)

type Limit float64

const (
	// Inf is the infinite rate limit; it allows all events (even if burst is zero).
	Inf = Limit(math.MaxFloat64)

	// infDuration is the duration returned by Delay when a Reservation is not Ok.
	infDuration = time.Duration(math.MaxInt64)
)

// Every converts a minimum time interval between events to a Limit.
func Every(interval time.Duration) Limit {
	if interval <= 0 {
		return Inf
	}
	return 1 / Limit(interval.Seconds())
}

// NewLimiter returns a new Limiter that allows events up to rate limit and permits
// bursts of at most burst tokens.
func NewLimiter(limit Limit, burst int) *Limiter {
	return &Limiter{
		limit: limit,
		burst: burst,
	}
}

type Limiter struct {
	mu        sync.Mutex
	limit     Limit
	burst     int
	tokens    float64
	last      time.Time // last is the last time the limiter's tokens field was updated
	lastEvent time.Time // lastEvent is the latest time of a rate-limited event (past or future)
}

// Limit returns the maximum overall event rate.
func (l *Limiter) Limit() Limit {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}

// Burst returns the maximum burst size. Burst is the maximum number of tokens
// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
// Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (l *Limiter) Burst() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.burst
}

// SetLimit is shorthand for SetLimitAt(time.Now(), newLimit).
func (l *Limiter) SetLimit(limit Limit) { l.SetLimitAt(time.Now(), limit) }

// SetLimitAt sets a new Limit for the limiter. The new Limit, and Burst, may be violated
// or underutilized by those which reserved (using Reserve or Wait) but did not yet act
// before SetLimitAt was called.
func (l *Limiter) SetLimitAt(t time.Time, newLimit Limit) {
	l.mu.Lock()
	defer l.mu.Unlock()

	t, tokens := l.calcTokensTo(t)

	l.last = t
	l.tokens = tokens
	l.limit = newLimit
}

// SetBurst is shorthand for SetBurstAt(time.Now(), newBurst).
func (l *Limiter) SetBurst(burst int) { l.SetBurstAt(time.Now(), burst) }

// SetBurstAt sets a new burst size for the limiter.
func (l *Limiter) SetBurstAt(t time.Time, newBurst int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	t, tokens := l.calcTokensTo(t)

	l.last = t
	l.tokens = tokens
	l.burst = newBurst
}

func (l *Limiter) Allow() bool                    { return l.AllowN(time.Now(), 1) }
func (l *Limiter) AllowN(t time.Time, n int) bool { return l.reserveWait(t, n, 0).ok }

func (l *Limiter) Wait(c context.Context) error         { return l.WaitN(c, 1) }
func (l *Limiter) WaitN(c context.Context, n int) error { return wait(c, l, n) }

func (l *Limiter) Reserve() Reservation                    { return l.ReserveN(time.Now(), 1) }
func (l *Limiter) ReserveN(t time.Time, n int) Reservation { return l.ReserveWait(t, n, infDuration) }
func (l *Limiter) ReserveWait(t time.Time, n int, wait time.Duration) Reservation {
	return ldptr.New(l.reserveWait(t, n, wait))
}

func (l *Limiter) reserveWait(t time.Time, n int, maxFutureReserve time.Duration) reservation {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.limit == Inf {
		return reservation{
			ok:        true,
			limiter:   l,
			tokens:    n,
			timeToAct: t,
		}
	}
	if l.limit == 0 {
		var ok bool
		if l.burst >= n {
			ok = true
			l.burst -= n
		}
		return reservation{
			ok:        ok,
			limiter:   l,
			tokens:    n,
			timeToAct: t,
		}
	}

	t, tokens := l.calcTokensTo(t)

	// Calculate the remaining number of tokens resulting from the request.
	tokens -= float64(n)

	// Calculate the wait duration
	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = l.limit.durationFromTokens(-tokens)
	}

	// Decide result
	ok := n <= l.burst && waitDuration <= maxFutureReserve
	// log.Printf("waitDuration:%s, maxFutureReserve:%s", waitDuration, maxFutureReserve)

	// Prepare reservation
	r := reservation{
		ok:      ok,
		limiter: l,
		limit:   l.limit,
	}
	if ok {
		r.tokens = n
		r.timeToAct = t.Add(waitDuration)

		// Update state
		l.last = t
		l.tokens = tokens
		l.lastEvent = r.timeToAct
	}

	return r
}

// calcTokensTo calculates and returns an updated state for lim resulting from the passage of time.
// lim is not changed.
// calcTokensTo requires that lim.mu is held.
func (l *Limiter) calcTokensTo(t time.Time) (newT time.Time, newTokens float64) {
	last := l.last
	if t.Before(last) {
		last = t
	}

	// Calculate the new number of tokens, due to time that passed.
	elapsed := t.Sub(last)
	delta := l.limit.tokensFromDuration(elapsed)
	tokens := l.tokens + delta
	if burst := float64(l.burst); tokens > burst {
		tokens = burst
	}
	return t, tokens
}

// durationFromTokens is a unit conversion function from the number of tokens to the duration
// of time it takes to accumulate them at a rate of limit tokens per second.
func (l Limit) durationFromTokens(tokens float64) time.Duration {
	if l <= 0 {
		return infDuration
	}
	seconds := tokens / float64(l)
	return time.Duration(float64(time.Second) * seconds)
}

// tokensFromDuration is a unit conversion function from a time duration to the number of tokens
// which could be accumulated during that duration at a rate of limit tokens per second.
func (l Limit) tokensFromDuration(d time.Duration) float64 {
	if l <= 0 {
		return 0
	}
	return d.Seconds() * float64(l)
}

// A reservation holds information about events that are permitted by a Limiter to happen after a delay.
// A reservation may be canceled, which may enable the Limiter to permit additional events.
type reservation struct {
	ok        bool
	limiter   *Limiter
	tokens    int
	timeToAct time.Time
	limit     Limit // This is the Limit at reservation time, it can change later.
}

// Ok returns whether the limiter can provide the requested number of tokens
// within the maximum wait time.  If Ok is false, Delay returns InfDuration, and
// Cancel does nothing.
func (r *reservation) Ok() bool {
	return r.ok
}

// Delay is shorthand for DelayFrom(time.Now()).
func (r *reservation) Delay() time.Duration { return r.DelayFrom(time.Now()) }

// DelayFrom returns the duration for which the reservation holder must wait
// before taking the reserved action.  Zero duration means act immediately.
// InfDuration means the limiter cannot grant the tokens requested in this
// reservation within the maximum wait time.
func (r *reservation) DelayFrom(t time.Time) time.Duration {
	if !r.ok {
		return infDuration
	}
	delay := r.timeToAct.Sub(t)
	if delay < 0 {
		return 0
	}
	// log.Printf("DelayFrom. time:%s, delay:%s, limiter:%s", str(t), delay, str(r.limiter))
	return delay
}

// Cancel is shorthand for CancelAt(time.Now()).
func (r *reservation) Cancel() {
	r.CancelAt(time.Now())
}

// CancelAt indicates that the reservation holder will not perform the reserved action
// and reverses the effects of this reservation on the rate limit as much as possible,
// considering that other reservations may have already been made.
func (r *reservation) CancelAt(t time.Time) {
	if !r.ok {
		return
	}

	r.limiter.mu.Lock()
	defer r.limiter.mu.Unlock()

	if r.limiter.limit == Inf || r.tokens == 0 || r.timeToAct.Before(t) {
		return
	}

	// calculate tokens to restore
	// The duration between lim.lastEvent and r.timeToAct tells us how many tokens were reserved
	// after r was obtained. These tokens should not be restored.
	restoreTokens := float64(r.tokens) - r.limit.tokensFromDuration(r.limiter.lastEvent.Sub(r.timeToAct))
	if restoreTokens <= 0 {
		return
	}
	// advance time to now
	t, tokens := r.limiter.calcTokensTo(t)
	// calculate new number of tokens
	tokens += restoreTokens
	if burst := float64(r.limiter.burst); tokens > burst {
		tokens = burst
	}
	// update state
	r.limiter.last = t
	r.limiter.tokens = tokens
	if r.timeToAct.Equal(r.limiter.lastEvent) {
		prevEvent := r.timeToAct.Add(r.limit.durationFromTokens(float64(-r.tokens)))
		if !prevEvent.Before(t) {
			r.limiter.lastEvent = prevEvent
		}
	}
	// log.Printf("CancelAt. time:%s, limiter:%s", str(t), str(r.limiter))
}
