/*
 * Copyright (C) distroy
 */

package ldrbtree

func (rbt *RBTree[T]) init() {
	if rbt.sentinel == nil {
		rbt.sentinel = getRBTreeNode[T](nil)
		rbt.root = rbt.sentinel
	}
	if rbt.Compare == nil {
		rbt.Compare = DefaultCompare[T]
	}
}

func (rbt *RBTree[T]) toDebugMap() map[string]interface{} {
	return map[string]interface{}{
		"count": rbt.count,
		"root":  rbt.root.toDebugMap(rbt.sentinel),
	}
}
