/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"
	"fmt"
	"time"

	"github.com/distroy/ldgo-base/internal/ctx_"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldlog"
)

type manager struct {
	tasks map[string]*cronTask
	// funcs       map[string]func(c context.Context, args map[string]any) error
	funcs map[string]func(*Task) error
}

func (m *manager) Register(taskName string, userFunc func(*Task) error) {
	if m.funcs == nil {
		m.funcs = make(map[string]func(*Task) error)
	}
	m.funcs[taskName] = userFunc
}

func (m *manager) Run(c context.Context) {
	if m.tasks == nil {
		m.tasks = make(map[string]*cronTask)
	}
	ldctx.LogI(c, "[timer] manager main goroutine begin")
	defer func() {
		ldctx.LogI(c, "[timer] manager main goroutine end")
	}()
	for {
		cc := c
		cc = ldctx.WithSequence(cc, "")
		m.main(cc)
	}
}

func (m *manager) main(c context.Context) {
	c, cancel := ldctx.WithCancel(c)

	ldctx.LogI(c, "[timer] manager main begin")
	defer func() {
		ldctx.LogI(c, "[timer] manager main end")
		cancel()
	}()

}

func (m *manager) taskGoroutine(c context.Context, task *cronTask) {
	ldctx.LogI(c, "[timer] task goroutine begin", ldlog.String("task", task.name))
	defer func() {
		ldctx.LogI(c, "[timer] task goroutine end", ldlog.String("task", task.name))
	}()

	// interval := task.cfg.Load().Interval.Duration()
	timer := time.NewTimer(1 * time.Second)
	// timerStart := time.Now()

	fnDo := func() {
		oldLogId := ldctx.GetSequence(c)
		newLogId := ctx_.NewSequence()
		ldctx.LogI(c, "[timer] get new logid for task", ldlog.String("new", newLogId),
			ldlog.String("task", task.name))
		c = ldctx.WithSequence(c, newLogId)
		ldctx.LogI(c, "[timer] get new logid for task", ldlog.String("old", oldLogId),
			ldlog.String("task", task.name))

		m.doTask(c, task)

		// interval = task.cfg.Load().Interval.Duration()
		// timer.Reset(interval)
		// timerStart = time.Now()
	}
	fnDo()
	for {
		select {
		case <-c.Done():
			return

		case <-timer.C:
			fnDo()
		}
	}
}

func (m *manager) doTask(c context.Context, task *cronTask) (_err error) {
	// cfg := task.cfg.Load()
	begin := time.Now()
	// ldctx.LogI(c, "[timer] do task begin", ldlog.String("task", task.name), ldlog.Reflect("cfg", cfg))

	defer func() {
		if e := recover(); e != nil {
			if ee, _ := e.(error); ee != nil {
				_err = ee
			} else {
				_err = fmt.Errorf("%v", ee)
			}
			ldctx.LogE(c, "[timer] do task panic", ldlog.String("task", task.name),
				ldlog.Error(_err))
		}
		ldctx.LogI(c, "[timer] do task end", ldlog.Duration("cost", time.Since(begin)),
			ldlog.String("task", task.name), ldlog.Error(_err))
	}()

	// fn := m.funcs[cfg.Func]
	fn := m.funcs[""]
	if fn == nil {
		// ldctx.LogE(c, "[timer] invalid func name", ldlog.String("func", cfg.Func))
		return lderr.Override(lderr.ErrInvalidParameter, "invalid func name")
	}

	err := fn(nil)
	if err != nil {
		ldctx.LogE(c, "[timer] do task fail", ldlog.Error(err))
		return err
	}

	return nil
}
