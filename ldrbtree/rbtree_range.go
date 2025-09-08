/*
 * Copyright (C) distroy
 */

package ldrbtree

type RBTreeRange[T any] struct {
	Begin RBTreeIterator[T] // [begin, end)
	End   RBTreeIterator[T] // [begin, end)
}

func (p *RBTreeRange[T]) HasNext() bool {
	return p.Begin.tree != nil && p.Begin != p.End
}

func (p *RBTreeRange[T]) Next() {
	p.Begin = p.Begin.Next()
}

func (p *RBTreeRange[T]) Get() T {
	return p.Begin.Get()
}

type RBTreeReverseRange[T any] struct {
	Begin RBTreeReverseIterator[T]
	End   RBTreeReverseIterator[T]
}

func (p *RBTreeReverseRange[T]) HasNext() bool {
	return p.Begin.tree != nil && p.Begin != p.End
}

func (p *RBTreeReverseRange[T]) Next() {
	p.Begin = p.Begin.Next()
}

func (p *RBTreeReverseRange[T]) Get() T {
	return p.Begin.Get()
}
