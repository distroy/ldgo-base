/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/distroy/ldgo-base/ldatomic"
	"github.com/distroy/ldgo-base/ldsync"
)

type Logger interface {
	Printf(fmt string, args ...any)
}

func Go(fn func() error) *GoPool { return GoN(1, fn) }
func GoN(n int, fn func() error) *GoPool {
	p := &GoPool{}
	p.GoN(n, fn)
	return p
}

type GoPool struct {
	log   ldatomic.Any[Logger]
	wg    ldsync.WaitGroup
	onErr ldatomic.Any[func(err error)]
	err   ldatomic.Error
}

func (p *GoPool) SetLogger(l Logger) { p.log.Store(l) }
func (p *GoPool) getLogger() Logger {
	l := p.log.Load()
	if l == nil {
		l = log.Default()
	}
	return l
}

func (p *GoPool) Count() int         { return p.wg.Count() }
func (p *GoPool) Go(fn func() error) { p.GoN(1, fn) }
func (p *GoPool) GoN(n int, fn func() error) {
	n = max(n, 1)

	fnGo := func() {
		defer func() {
			if err := recover(); err != nil {
				buf := debug.Stack()

				// log.Println(err, ldconv.BytesToStrUnsafe(buf))
				p.getLogger().Printf("[go pool] go func panic. err:%v, stack:\n%s", err, buf)
				err := fmt.Errorf("go func panic. err:%v", err)
				p.setError(err)
			}
			p.wg.Done()
		}()

		err := fn()
		p.setError(err)
	}

	p.wg.Add(n)
	for i := 0; i < n; i++ {
		go fnGo()
	}
}

func (p *GoPool) Wait() error {
	p.wg.Wait()
	return p.err.Load()
}

func (p *GoPool) OnError(f func(err error)) { p.onErr.Store(f) }
func (p *GoPool) setError(err error) {
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
