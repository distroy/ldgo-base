/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"

	"github.com/distroy/ldgo-base/internal/ctx_"
	"github.com/distroy/ldgo-base/ldctx"
)

type ctxKey int

const (
	ctxKeyTask ctxKey = iota + 1
)

func newContext(task *Task) context.Context {
	seq := task.Info.GetSequence()
	if seq == "" {
		seq = ctx_.NewSequence()
	}
	ctx := ldctx.Default()
	ctx = ldctx.WithSequence(ctx, seq)
	ctx = ldctx.WithValue(ctx, ctxKeyTask, task)

	return ctx
}

func GetTask(ctx context.Context) *TaskInfo {
	v, _ := ctx.Value(ctxKeyTask).(*TaskInfo)
	return v
}
