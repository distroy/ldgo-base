/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"context"
	"fmt"
	"reflect"

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
func (e Engine[A]) Run(c context.Context) { e.adaptor.Run(c) }

func (e Engine[A]) SetProgress(c context.Context, task *Task, progress string) {
	e.adaptor.SetProgress(c, task, progress)
}

// RegService registers timer task handlers
// Parameters:
//   - service: service instance containing task methods
//   - tasks: task mapping where key is task name and value is method name in service
//
// RegService 注册定时任务处理器
// 参数说明：
//   - service: 包含任务处理方法的服务实例
//   - tasks: 任务映射表，key为任务名称，value为服务中的方法名
func (e Engine[A]) RegService(c context.Context, service any, tasks map[string]string) {
	for name, method := range tasks {
		h := e.newHandler(service, method)
		e.adaptor.Register(c, name, h.Do)
		ldctx.LogI(c, "[ldtimer] register handler succ", ldlog.String("adaptor", e.adaptor.Name()),
			ldlog.String("task", name), ldlog.String("func", h.Name))
	}
}

func (e Engine[A]) newHandler(service any, method string) *handler {
	v := reflect.ValueOf(service)
	m := v.MethodByName(method)
	if !m.IsValid() {
		panic(fmt.Errorf("cannot find the method `%s` from type `%s`", method, v.Type().String()))
	}

	h := &handler{
		Service: nil,
		Func:    m.Interface(),
		Name:    method,
		Value:   m,
	}

	h.Init()
	return h
}
