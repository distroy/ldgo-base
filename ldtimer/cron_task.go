/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"
	"time"

	"github.com/distroy/ldgo-base/internal/ctx_"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/ldlog"
	"github.com/distroy/ldgo-base/ldsync"
)

var (
	_ taskTimer = (*CronSpec)(nil)
	_ taskTimer = (*CronEvery)(nil)
)

type taskTimer interface {
	Match(t time.Time) bool
	Next(from time.Time) (time.Time, bool)
}

type cronTask struct {
	cron    *Cron
	Name    string
	Handler *Handler
	Timer   taskTimer
	Done    ldsync.Done
}

func (t *cronTask) Main() {
	c := ldctx.Default()
	var timer *time.Timer

	ldctx.LogI(c, "[ldtimer] cron task goroutine begin", ldlog.String("name", t.Name))
	defer func() {
		ldctx.LogI(c, "[ldtimer] cron task goroutine end", ldlog.String("name", t.Name))
		if timer != nil {
			timer.Stop()
		}
	}()

	run := func(ok bool) {
		if ok {
			t.run()
		}
	}

	last := time.Now()
	run(t.Timer.Match(last))

	for {
		next, ok := t.Timer.Next(last)
		now := time.Now()
		dur := next.Sub(now)
		if dur <= 0 {
			last = now
			run(ok)
			continue
		}
		if timer == nil {
			timer = time.NewTimer(dur)
		} else {
			timer.Reset(dur)
		}
		select {
		case <-t.Done.Chan():
			return

		case now := <-timer.C:
			last = now
			run(ok)
			continue
		}
	}
}

func (t *cronTask) run() (_err error) {
	c := ldctx.Default()
	c = ldctx.WithSequence(c, ctx_.NewSequence())
	return t.Handler.Do(&Task{
		Info:    &cronTaskInfo{},
		Adaptor: cronAdaptor{t.cron},
	})
}

type cronTaskInfo struct {
}

func (_ cronTaskInfo) GetParams() string                              { return "" }
func (_ cronTaskInfo) WithSequence(c context.Context) context.Context { return c }
