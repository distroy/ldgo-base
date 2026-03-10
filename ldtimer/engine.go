/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"

	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/ldlog"
)

func NewEngine[A Adaptor](adaptor A) Engine[A] {
	return Engine[A]{
		adaptor: adaptor,
	}
}

type Engine[A Adaptor] struct {
	adaptor A
}

func (e Engine[A]) Adaptor() Adaptor      { return e.adaptor }
func (e Engine[A]) Name() string          { return e.adaptor.Name() }
func (e Engine[A]) Run(c context.Context) { e.adaptor.Run(c) }

func (e Engine[A]) SetProgress(c context.Context, task *Task, progress string) {
	e.adaptor.SetProgress(c, task, progress)
}

// RegService registers timer task handlers
// Parameters:
//   - service: service instance containing task methods
//   - tasks: task mapping where key is task name and value is method name in service
//
// method should be:
//   - func()
//   - func() error
//   - func(c context.Context)
//   - func(c context.Context) error
//   - func(c context.Context, task *ldtimer.Task)
//   - func(c context.Context, task *ldtimer.Task) error
//   - func(c context.Context, params string)
//   - func(c context.Context, params string) error
//   - func(c context.Context, args *T)
//   - func(c context.Context, args *T) error
func (e Engine[A]) RegService(c context.Context, service any, tasks map[string]string) {
	for name, method := range tasks {
		h := NewHandlerByMethod(service, method)
		e.adaptor.Register(c, name, h.Do)
		ldctx.LogI(c, "[ldtimer] register handler succ", ldlog.String("adaptor", e.adaptor.Name()),
			ldlog.String("task", name), ldlog.String("func", h.Name))
	}
}

// RegHandler registers timer task handlers
// Parameters:
//   - name: task name
//   - handler: task handler
//
// handler should be:
//   - func()
//   - func() error
//   - func(c context.Context)
//   - func(c context.Context) error
//   - func(c context.Context, task *ldtimer.Task)
//   - func(c context.Context, task *ldtimer.Task) error
//   - func(c context.Context, params string)
//   - func(c context.Context, params string) error
//   - func(c context.Context, args *T)
//   - func(c context.Context, args *T) error
func (e Engine[A]) RegHandler(c context.Context, name string, handler any) {
	h := NewHandler(handler)
	e.adaptor.Register(c, name, h.Do)
	ldctx.LogI(c, "[ldtimer] register handler succ", ldlog.String("adaptor", e.adaptor.Name()),
		ldlog.String("task", name), ldlog.String("func", h.Name))
}
