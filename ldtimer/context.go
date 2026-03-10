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
	ctx := ldctx.Default()
	ctx = task.Info.WithSequence(ctx)
	if ldctx.GetSequence(ctx) == "" {
		ctx = ldctx.WithSequence(ctx, ctx_.NewSequence())
	}
	ctx = ldctx.WithValue(ctx, ctxKeyTask, task)

	return ctx
}

func GetTask(ctx context.Context) *TaskInfo {
	v, _ := ctx.Value(ctxKeyTask).(*TaskInfo)
	return v
}
