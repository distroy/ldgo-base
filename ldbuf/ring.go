/*
 * Copyright (C) distroy
 */

package ldbuf

import (
	"github.com/distroy/ldgo-base/internal/buf_"
)

func MakeRing[T any](n int) Ring[T] {
	n = max(n, 1)
	buf := make([]T, n)
	return Ring[T]{
		buf: buf_.NewRing(buf),
	}
}

type Ring[T any] struct {
	buf *buf_.Ring[T]
}

func (b Ring[T]) Cap() int  { return b.buf.Cap() }
func (b Ring[T]) Size() int { return b.buf.Size() }

func (b Ring[T]) Close() error { return b.buf.Close() }
func (b Ring[T]) Closed() bool { return b.buf.Closed() }

func (b Ring[T]) Put(d T) bool   { return b.put(d) }
func (b Ring[T]) Pop() (T, bool) { return b.pop() }

func (b Ring[T]) pop() (T, bool) {
	d, ok, _ := b.buf.Pop()
	return d, ok
}

func (b Ring[T]) put(d T) bool {
	ok, _ := b.buf.Put(d)
	return ok
}

func MakeBlockingRing[T any](n int) BlockingRing[T] {
	n = max(n, 1)
	buf := make([]T, n)
	b := BlockingRing[T]{
		buf: buf_.NewBlockingRing(buf),
	}
	return b
}

type BlockingRing[T any] struct {
	buf *buf_.BlockingRing[T]
}

func (b BlockingRing[T]) Cap() int  { return b.buf.Cap() }
func (b BlockingRing[T]) Size() int { return b.buf.Size() }

func (b BlockingRing[T]) Close() error { return b.buf.Close() }
func (b BlockingRing[T]) Closed() bool { return b.buf.Closed() }

func (b BlockingRing[T]) Put(d T) bool   { return b.Put(d) }
func (b BlockingRing[T]) Pop() (T, bool) { return b.Pop() }
