/*
 * Copyright (C) distroy
 */

package _buf

import (
	"sync"
)

func NewBlockingRing[T any](buf []T) *BlockingRing[T] {
	b := &BlockingRing[T]{
		buf: makeRing(buf),
	}
	b.init()
	return b
}

type BlockingRing[T any] struct {
	buf Ring[T]

	mu        sync.Mutex
	readCond  sync.Cond
	writeCond sync.Cond
}

func (b *BlockingRing[T]) init() {
	b.readCond.L = &b.mu
	b.writeCond.L = &b.mu
}

func (b *BlockingRing[T]) Close() error {
	return blockCall(b, func() error {
		err := b.buf.Close()
		if err == nil {
			b.readCond.Broadcast()
			b.writeCond.Broadcast()
		}
		return err
	})
}
func (b *BlockingRing[T]) Closed() bool {
	return blockCall(b, func() bool {
		return b.buf.Closed()
	})
}

func (b *BlockingRing[T]) Cap() int  { return blockCall(b, func() int { return b.buf.Cap() }) }
func (b *BlockingRing[T]) Size() int { return blockCall(b, func() int { return b.buf.Size() }) }

func (b *BlockingRing[T]) Write(d []T) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	{
		n, err := b.buf.Write(d)
		if err != nil {
			return n, err
		}

		if n > 0 {
			b.readCond.Signal()
			return n, nil
		}
	}

	b.writeCond.Wait()

	n, err := b.buf.Write(d)
	if n > 0 {
		b.readCond.Signal()
		if !b.buf.full {
			b.writeCond.Signal()
		}
	}
	return n, err
}

func (b *BlockingRing[T]) Read(d []T) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	{
		n, err := b.buf.Read(d)
		if err != nil {
			return n, err
		}
		if n > 0 {
			b.writeCond.Signal()
			return n, nil
		}
	}

	b.readCond.Wait()

	n, err := b.buf.Read(d)
	if n > 0 {
		b.writeCond.Signal()
		if b.buf.begin != b.buf.end {
			b.readCond.Signal()
		}
	}
	return n, err
}

func (b *BlockingRing[T]) Put(d T) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	{
		ok, err := b.buf.Put(d)
		if err != nil {
			return err
		}
		if ok {
			b.readCond.Signal()
			return nil
		}
	}

	b.writeCond.Wait()
	ok, err := b.buf.Put(d)
	if ok {
		b.readCond.Signal()
		if !b.buf.full {
			b.writeCond.Signal()
		}
	}
	return err
}

func (b *BlockingRing[T]) Pop() (T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	{
		d, ok, err := b.buf.Pop()
		if err != nil {
			return d, err
		}
		if ok {
			b.writeCond.Signal()
			return d, nil
		}
	}

	b.readCond.Wait()

	d, ok, err := b.buf.Pop()
	if ok {
		b.writeCond.Signal()
		if b.buf.begin != b.buf.end {
			b.readCond.Signal()
		}
	}
	return d, err
}

func blockCall[T any, R any](b *BlockingRing[T], f func() R) R {
	b.mu.Lock()
	r := f()
	b.mu.Unlock()
	return r
}
