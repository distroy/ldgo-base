/*
 * Copyright (C) distroy
 */

package ldtimer

import "context"

type TaskInfo interface {
	GetParams() string
	GetSequence() string
}

type Task struct {
	Info    TaskInfo `json:"info"`
	Adaptor Adaptor  `json:"-"`
}

type Adaptor interface {
	Name() string

	Register(c context.Context, taskName string, taskFunc func(*Task) error)

	Run(c context.Context)

	SetProgress(c context.Context, task *Task, progress string)
}
