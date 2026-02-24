/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"

	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/ldlog"
)

func NewCron() *Cron {
	cron := &Cron{}
	engine := NewEngine(cronAdaptor{Cron: cron})
	cron.engine = engine
	return cron
}

type cronAdaptor struct {
	*Cron
}

func (c cronAdaptor) Name() string { return "cron" }

func (c cronAdaptor) SetProgress(ctx context.Context, task *Task, progress string) {
	ldctx.LogI(ctx, "[ldtimer] cron task progress", ldlog.String("progress", progress))
}

func (c cronAdaptor) Register(ctx context.Context, taskName string, taskFunc func(*Task) error) {
	c.Cron.register(ctx, taskName, taskFunc)
}

func (c cronAdaptor) Run(ctx context.Context) { c.Cron.run(ctx) }

type Cron struct {
	engine Engine[cronAdaptor]
	tasks  map[string]*cronTask
	funcs  map[string]func(*Task) error
}

func (c *Cron) register(ctx context.Context, taskName string, taskFunc func(*Task) error) {
	if c.funcs == nil {
		c.funcs = make(map[string]func(*Task) error)
	}
	c.funcs[taskName] = taskFunc
}

func (c *Cron) Run(ctx context.Context) {
	c.engine.Run(ctx)
}

func (c *Cron) run(ctx context.Context) {
	if c.tasks == nil {
		c.tasks = make(map[string]*cronTask)
	}
	ldctx.LogI(ctx, "[ldtimer] cron run begin")
	defer func() {
		ldctx.LogI(ctx, "[ldtimer] cron run end")
	}()
	for {
		cc := ctx
		cc = ldctx.WithSequence(cc, "")
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
