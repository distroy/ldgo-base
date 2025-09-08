/*
 * Copyright (C) distroy
 */

package ldrbtree

type MapNode[K any, V any] rbtreeIterator[Pair[K, V]]

func (n MapNode[K, V]) Key() K { return n.node.Data.Key }
func (n MapNode[K, V]) Value(new ...V) (old V) {
	old = n.node.Data.Value
	if len(new) > 0 {
		n.node.Data.Value = new[0]
	}
	return
}

type mapIterator[K any, V any] rbtreeIterator[Pair[K, V]]

func (it mapIterator[K, V]) next(name string, iface rbtreeInterface[Pair[K, V]]) mapIterator[K, V] {
	return (mapIterator[K, V])((*rbtreeIterator[Pair[K, V]])(&it).next(name, iface))
}

func (it mapIterator[K, V]) prev(name string, iface rbtreeInterface[Pair[K, V]]) mapIterator[K, V] {
	return (mapIterator[K, V])((*rbtreeIterator[Pair[K, V]])(&it).prev(name, iface))
}

type MapIterator[K any, V any] mapIterator[K, V]

func (i MapIterator[K, V]) base() mapIterator[K, V] { return mapIterator[K, V](i) }
func (i MapIterator[K, V]) Get() MapNode[K, V]      { return MapNode[K, V](i) }
func (i MapIterator[K, V]) Key() K                  { return i.Get().Key() }
func (i MapIterator[K, V]) Value(new ...V) (old V)  { return i.Get().Value(new...) }

func (i MapIterator[K, V]) Next() MapIterator[K, V] {
	return MapIterator[K, V](i.base().next("map iterator", forward(i.tree)))
}

func (i MapIterator[K, V]) Prev() MapIterator[K, V] {
	return MapIterator[K, V](i.base().prev("map iterator", forward(i.tree)))
}

type MapReverseIterator[K any, V any] mapIterator[K, V]

func (i MapReverseIterator[K, V]) base() mapIterator[K, V] { return mapIterator[K, V](i) }
func (i MapReverseIterator[K, V]) Get() MapNode[K, V]      { return MapNode[K, V](i) }
func (i MapReverseIterator[K, V]) Key() K                  { return i.Get().Key() }
func (i MapReverseIterator[K, V]) Value(new ...V) (old V)  { return i.Get().Value(new...) }

func (i MapReverseIterator[K, V]) Next() MapReverseIterator[K, V] {
	return MapReverseIterator[K, V](i.base().next("map reverse iterator", reverse(i.tree)))
}

func (i MapReverseIterator[K, V]) Prev() MapReverseIterator[K, V] {
	return MapReverseIterator[K, V](i.base().prev("map reverse iterator", reverse(i.tree)))
}
