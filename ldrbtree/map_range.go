/*
 * Copyright (C) distroy
 */

package ldrbtree

type MapRange[K any, V any] struct {
	Begin MapIterator[K, V] // [begin, end)
	End   MapIterator[K, V] // [begin, end)
}

func (p *MapRange[K, V]) HasNext() bool          { return p.Begin.tree != nil && p.Begin != p.End }
func (p *MapRange[K, V]) Next()                  { p.Begin = p.Begin.Next() }
func (p *MapRange[K, V]) Get() MapNode[K, V]     { return p.Begin.Get() }
func (p *MapRange[K, V]) Key() K                 { return p.Begin.Key() }
func (p *MapRange[K, V]) Value(new ...V) (old V) { return p.Begin.Value(new...) }

type MapReverseRange[K any, V any] struct {
	Begin MapReverseIterator[K, V] // [begin, end)
	End   MapReverseIterator[K, V] // [begin, end)
}

func (p *MapReverseRange[K, V]) HasNext() bool          { return p.Begin.tree != nil && p.Begin != p.End }
func (p *MapReverseRange[K, V]) Next()                  { p.Begin = p.Begin.Next() }
func (p *MapReverseRange[K, V]) Get() MapNode[K, V]     { return p.Begin.Get() }
func (p *MapReverseRange[K, V]) Key() K                 { return p.Begin.Key() }
func (p *MapReverseRange[K, V]) Value(new ...V) (old V) { return p.Begin.Value(new...) }
