/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"fmt"
	"runtime/debug"

	"github.com/distroy/ldgo-base/ldatomic"
)

func NewAsyncPool(concurrency int) *AsyncPool {
	ap := &AsyncPool{}
	ap.Start(concurrency)
	return ap
}

type AsyncPool struct {
	asyncBase[func() error]

	onErr ldatomic.Any[func(err error)]
	err   ldatomic.Error
}

func (p *AsyncPool) Start(concurrency int) { p.asyncBase.start(concurrency, p.doWithRecover) }
func (p *AsyncPool) Reset(concurrency int) { p.asyncBase.reset(concurrency, p.doWithRecover) }

func (p *AsyncPool) Capacity() int { return p.asyncBase.getCap() }
func (p *AsyncPool) Running() int  { return p.asyncBase.getLen() }

func (p *AsyncPool) SetLogger(l Logger) { p.asyncBase.setLogger(l) }

func (p *AsyncPool) init() { p.asyncBase.init(p.doWithRecover) }

func (p *AsyncPool) Wait() error {
	p.asyncBase.wait()
	return p.err.Load()
}

func (p *AsyncPool) Close() { p.asyncBase.close() }

func (p *AsyncPool) Async() chan<- func() error {
	p.init()
	return p.asyncBase.async()
}

func (p *AsyncPool) OnError(f func(err error)) { p.onErr.Store(f) }
func (p *AsyncPool) setError(err error) {
	if err == nil {
		return
	}
	ok := p.err.CompareAndSwap(nil, err)
	if !ok {
		return
	}
	if f := p.onErr.Load(); f != nil {
		f(err)
	}
}

func (p *AsyncPool) doWithRecover(fn func() error) {
	defer func() {
		if err := recover(); err != nil {
			buf := debug.Stack()
			p.getLogger().Printf("[async pool] do async func panic. err:%v, stack:\n%s", err, buf)
			err := fmt.Errorf("async func panic. err:%v", err)
			p.setError(err)
		}
	}()

	err := fn()
	p.setError(err)
}
