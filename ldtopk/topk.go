/*
 * Copyright (C) distroy
 */

package ldtopk

import (
	"sync"

	"github.com/distroy/ldgo-base/internal/cmp_"
)

type LessFunc[T any] = func(a, b T) bool

func defaultLess[T any](a, b T) bool {
	return cmp_.Compare(a, b) < 0
}

func New[T any](k int, less func(a, b T) bool) *Topk[T] {
	p := &Topk[T]{
		Size: k,
		Less: less,
	}
	return p
}

// Topk will keep at most size elements for which less returns true
type Topk[T any] struct {
	Size   int         // K
	Less   LessFunc[T] // less func
	Locker sync.Locker // locker
	data   []T         //
}

func (p *Topk[T]) Data() []T { return p.data }

func (p *Topk[T]) SetLocker(l sync.Locker) { p.Locker = l }

func (p *Topk[T]) Add(d T) bool {
	locker := p.Locker
	if locker == nil {
		locker = nullLocker{}
	}

	locker.Lock()
	defer locker.Unlock()

	p.init()

	less := p.Less
	if less == nil {
		less = defaultLess[T]
	}
	b, ok := topkAdd(p.data, p.Size, d, less)
	p.data = b
	return ok
}

func (p *Topk[T]) init() {
	if p.data == nil && p.Size > 0 {
		p.data = make([]T, 0, p.Size)
	}
}

type nullLocker struct{}

func (_ nullLocker) Lock()   {}
func (_ nullLocker) Unlock() {}
