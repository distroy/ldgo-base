/*
 * Copyright (C) distroy
 */

package ldrbtree

func (rbt *RBTree[T]) rotateLeft(node *rbtreeNode[T]) {
	//    N             T
	//   / \           / \
	//  a   T   -->   N   c
	//     / \       / \
	//    b   c     a   b

	root := &rbt.root
	sentinel := rbt.sentinel

	temp := node.Right
	node.Right = temp.Left

	if temp.Left != sentinel {
		temp.Left.Parent = node
	}

	temp.Parent = node.Parent

	if node == *root {
		*root = temp
	} else if node == node.Parent.Left {
		node.Parent.Left = temp
	} else {
		node.Parent.Right = temp
	}

	temp.Left = node
	node.Parent = temp
}

func (rbt *RBTree[T]) rotateRight(node *rbtreeNode[T]) {
	//      N         T
	//     / \       / \
	//    T   c --> a   N
	//   / \           / \
	//  a   b         b   c

	root := &rbt.root
	sentinel := rbt.sentinel

	temp := node.Left
	node.Left = temp.Right

	if temp.Right != sentinel {
		temp.Right.Parent = node
	}

	temp.Parent = node.Parent

	if node == *root {
		*root = temp
	} else if node == node.Parent.Right {
		node.Parent.Right = temp
	} else {
		node.Parent.Left = temp
	}

	temp.Right = node
	node.Parent = temp
}
