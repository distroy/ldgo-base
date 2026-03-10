/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"
	"sync"
	"time"

	"github.com/distroy/ldgo-base/internal/ctx_"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldlog"
)

func NewCron() *Cron {
	cron := &Cron{}
	return cron
}

type cronAdaptor struct {
	*Cron
}

func (c cronAdaptor) Name() string { return "cron" }

func (c cronAdaptor) SetProgress(ctx context.Context, task *Task, progress string) {
	panic(lderr.Override(lderr.ErrInternalServerError, "unimplemented interface"))
}

func (c cronAdaptor) Register(ctx context.Context, taskName string, taskFunc func(*Task) error) {
	panic(lderr.Override(lderr.ErrInternalServerError, "unimplemented interface"))
}

func (c cronAdaptor) Run(ctx context.Context) {
	panic(lderr.Override(lderr.ErrInternalServerError, "unimplemented interface"))
}

type Cron struct {
	lock  sync.Mutex
	tasks map[string]*cronTask
}

func (c *Cron) Run(ctx context.Context) {
	if c.tasks == nil {
		c.tasks = make(map[string]*cronTask)
	}
	ldctx.LogI(ctx, "[ldtimer] cron run begin")
	defer func() {
		ldctx.LogI(ctx, "[ldtimer] cron run end")
	}()
	for {
		cc := ctx
		cc = ldctx.WithSequence(cc, ctx_.NewSequence())
		c.main(cc)
	}
}

func (c *Cron) main(ctx context.Context) {
	ctx, cancel := ldctx.WithCancel(ctx)

	ldctx.LogI(ctx, "[ldtimer] cron main begin")
	defer func() {
		ldctx.LogI(ctx, "[ldtimer] cron main end")
		cancel()
	}()
}

func (c *Cron) Every(ctx context.Context, name string, d time.Duration, handler any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if task := c.tasks[name]; task != nil {
		ldctx.LogW(ctx, "[ldtimer] cron task had been registered", ldlog.String("name", name))
		return
	}

	h := NewHandler(handler)
	task := &cronTask{
		cron:    c,
		Name:    name,
		Handler: h,
		Timer:   CronEvery{d},
	}
	c.tasks[name] = task
}
