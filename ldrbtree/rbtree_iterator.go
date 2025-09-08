/*
 * Copyright (C) distroy
 */

package ldrbtree

import (
	"fmt"
)

type rbtreeIterator[T any] struct {
	tree *RBTree[T]
	node *rbtreeNode[T]
}

func (it rbtreeIterator[T]) next(name string, iface rbtreeInterface[T]) rbtreeIterator[T] {
	if it.tree == nil {
		panic(fmt.Errorf("the %s not bind any red-black tree, can not move next", name))
	}

	sentinel := it.tree.sentinel
	if it.node == sentinel {
		panic(fmt.Errorf("the %s is already at the end of the red-black tree, can not move next", name))
	}

	node := it.node.next(iface)
	return rbtreeIterator[T]{
		tree: it.tree,
		node: node,
	}
}

func (it rbtreeIterator[T]) prev(name string, iface rbtreeInterface[T]) rbtreeIterator[T] {
	if it.tree == nil {
		panic(fmt.Errorf("the %s does not bind any red-black tree, can not move prev", name))
	}

	sentinel := it.tree.sentinel
	root := it.tree.root

	var node *rbtreeNode[T]
	if it.node == sentinel {
		node = root.max(iface)
	} else {
		node = it.node.prev(iface)
	}

	if node == sentinel {
		panic(fmt.Errorf("the %s is already at the begin of the red-black tree, can not move prev", name))
	}

	return rbtreeIterator[T]{
		tree: it.tree,
		node: node,
	}
}

type RBTreeIterator[T any] rbtreeIterator[T]

func (i RBTreeIterator[T]) Get() T {
	return i.node.Data
}

func (i RBTreeIterator[T]) base() rbtreeIterator[T] {
	return rbtreeIterator[T](i)
}

func (i RBTreeIterator[T]) Next() RBTreeIterator[T] {
	return RBTreeIterator[T](i.base().next("rbtree iterator", forward(i.tree)))
}

func (i RBTreeIterator[T]) Prev() RBTreeIterator[T] {
	return RBTreeIterator[T](i.base().prev("rbtree iterator", forward(i.tree)))
}

type RBTreeReverseIterator[T any] rbtreeIterator[T]

func (i RBTreeReverseIterator[T]) Get() T {
	return i.node.Data
}

func (i RBTreeReverseIterator[T]) base() rbtreeIterator[T] {
	return rbtreeIterator[T](i)
}

func (i RBTreeReverseIterator[T]) Next() RBTreeReverseIterator[T] {
	return RBTreeReverseIterator[T](i.base().next("rbtree reverse iterator", reverse(i.tree)))
}

func (i RBTreeReverseIterator[T]) Prev() RBTreeReverseIterator[T] {
	return RBTreeReverseIterator[T](i.base().prev("rbtree reverse iterator", reverse(i.tree)))
}
