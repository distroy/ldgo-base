/*
 * Copyright (C) distroy
 */

package ldrbtree

// rbtreeLowerBound returns the first node >= d
func rbtreeLowerBound[T any](d T, iface rbtreeInterface[T]) *rbtreeNode[T] {
	sentinel := iface.Sentinel()
	node := iface.Root()

	res := sentinel
	for node != sentinel {
		r := iface.Compare(node.Data, d)
		if r == 0 {
			// node == d
			res = node
			node = iface.Left(node)

		} else if r > 0 {
			// node > d
			res = node
			node = iface.Left(node)

		} else {
			// node < d
			node = iface.Right(node)
		}
	}
	return res
}

// rbtreeUpperBound returns the first node > d
func rbtreeUpperBound[T any](d T, iface rbtreeInterface[T]) *rbtreeNode[T] {
	sentinel := iface.Sentinel()
	node := iface.Root()

	res := sentinel
	for node != sentinel {
		r := iface.Compare(node.Data, d)
		if r == 0 {
			// node == 0
			node = iface.Right(node)

		} else if r > 0 {
			// node > d
			res = node
			node = iface.Left(node)

		} else {
			// node < d
			node = iface.Right(node)
		}
	}
	return res
}

func rbtreeBeginIterator[T any](iface rbtreeInterface[T]) rbtreeIterator[T] {
	return rbtreeIterator[T]{
		tree: iface.Tree(),
		node: iface.Root().min(iface),
	}
}

func rbtreeEndIterator[T any](iface rbtreeInterface[T]) rbtreeIterator[T] {
	return rbtreeIterator[T]{
		tree: iface.Tree(),
		node: iface.Sentinel(),
	}
}
