/*
 * Copyright (C) distroy
 */

package ldrbtree

type Map[K any, V any] struct {
	tree       RBTree[Pair[K, V]]
	KeyCompare CompareFunc[K]
}

func (m *Map[K, V]) init() {
	if m.KeyCompare == nil {
		m.KeyCompare = DefaultCompare[K]
	}
	// if m.ValueCompare == nil {
	// 	m.ValueCompare = DefaultCompare
	// }
	if m.tree.Compare == nil {
		m.tree.Compare = wrapPairCompare(m)
	}

	m.tree.init()
}

func (m *Map[K, V]) Len() int {
	m.init()
	return m.tree.Len()
}

func (m *Map[K, V]) Clear() {
	m.init()
	m.tree.Clear()
}

func (m *Map[K, V]) Insert(key K, value V) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](m.tree.Insert(Pair[K, V]{Key: key, Value: value}))
}

func (m *Map[K, V]) InsertOrSearch(key K, value V) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](m.tree.InsertOrSearch(Pair[K, V]{Key: key, Value: value}))
}

func (m *Map[K, V]) InsertOrAssign(key K, value V) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](m.tree.InsertOrAssign(Pair[K, V]{Key: key, Value: value}))
}

func (m *Map[K, V]) Count(key K) int {
	m.init()
	return m.tree.Count(Pair[K, V]{Key: key})
}

func (m *Map[K, V]) Delete(it MapIterator[K, V]) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](m.tree.Delete(RBTreeIterator[Pair[K, V]](it)))
}

// Search returns the first pair.key == key
func (m *Map[K, V]) Search(key K) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](m.tree.Search(Pair[K, V]{Key: key}))
}

// SearchRange returns the node range [first pair.key == key, first pair.key > key)
func (m *Map[K, V]) SearchRange(key K) *MapRange[K, V] {
	m.init()
	it := m.tree.SearchRange(Pair[K, V]{Key: key})
	return &MapRange[K, V]{
		Begin: MapIterator[K, V](it.Begin),
		End:   MapIterator[K, V](it.End),
	}
}

// LowerBound returns the first pair.key >= key
func (m *Map[K, V]) LowerBound(key K) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V]{
		tree: &m.tree,
		node: rbtreeLowerBound(Pair[K, V]{Key: key}, forward(&m.tree)),
	}
}

// UpperBound returns the first pair.key > key
func (m *Map[K, V]) UpperBound(key K) MapIterator[K, V] {
	m.init()
	return MapIterator[K, V]{
		tree: &m.tree,
		node: rbtreeUpperBound(Pair[K, V]{Key: key}, forward(&m.tree)),
	}
}

func (m *Map[K, V]) Range() *MapRange[K, V] {
	m.init()
	return &MapRange[K, V]{
		Begin: MapIterator[K, V](rbtreeBeginIterator(forward(&m.tree))),
		End:   MapIterator[K, V](rbtreeEndIterator(forward(&m.tree))),
	}
}

func (m *Map[K, V]) Begin() MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](rbtreeBeginIterator(forward(&m.tree)))
}

func (m *Map[K, V]) End() MapIterator[K, V] {
	m.init()
	return MapIterator[K, V](rbtreeEndIterator(forward(&m.tree)))
}

// RDelete is reverse delete
func (m *Map[K, V]) RDelete(it MapReverseIterator[K, V]) MapReverseIterator[K, V] {
	m.init()
	return MapReverseIterator[K, V](m.tree.RDelete(RBTreeReverseIterator[Pair[K, V]](it)))
}

// RSearch is reverse search
// RSearch returns the last pair.key == key
func (m *Map[K, V]) RSearch(key K) MapReverseIterator[K, V] {
	m.init()
	return MapReverseIterator[K, V](m.tree.RSearch(Pair[K, V]{Key: key}))
}

// RSearchRange is reverse search range
// RSearchRange returns the node range [last pair.key == key, last pair.key < key)
func (m *Map[K, V]) RSearchRange(key K) *MapReverseRange[K, V] {
	m.init()
	it := m.tree.RSearchRange(Pair[K, V]{Key: key})
	return &MapReverseRange[K, V]{
		Begin: MapReverseIterator[K, V](it.Begin),
		End:   MapReverseIterator[K, V](it.End),
	}
}

// RLowerBound is reverse lower bound
// RLowerBound returns the last pair.key <= key
func (m *Map[K, V]) RLowerBound(key K) MapReverseIterator[K, V] {
	m.init()
	return MapReverseIterator[K, V]{
		tree: &m.tree,
		node: rbtreeLowerBound(Pair[K, V]{Key: key}, reverse(&m.tree)),
	}
}

// RUpperBound is reverse upper bound
// RUpperBound returns the last pair.key < key
func (m *Map[K, V]) RUpperBound(key K) MapReverseIterator[K, V] {
	m.init()
	return MapReverseIterator[K, V]{
		tree: &m.tree,
		node: rbtreeUpperBound(Pair[K, V]{Key: key}, reverse(&m.tree)),
	}
}

// RRange is reverse range
func (m *Map[K, V]) RRange() *MapReverseRange[K, V] {
	m.init()
	return &MapReverseRange[K, V]{
		Begin: MapReverseIterator[K, V](rbtreeBeginIterator(reverse(&m.tree))),
		End:   MapReverseIterator[K, V](rbtreeEndIterator(reverse(&m.tree))),
	}
}

// RBegin is reverse begin
func (m *Map[K, V]) RBegin() MapReverseIterator[K, V] {
	m.init()
	return MapReverseIterator[K, V](rbtreeBeginIterator(reverse(&m.tree)))
}

// REnd is reverse end
func (m *Map[K, V]) REnd() MapReverseIterator[K, V] {
	m.init()
	return MapReverseIterator[K, V](rbtreeEndIterator(reverse(&m.tree)))
}
