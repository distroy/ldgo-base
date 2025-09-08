/*
 * Copyright (C) distroy
 */

package ldrbtree

type rbtreeInterface[T any] interface {
	Tree() *RBTree[T]
	Root() *rbtreeNode[T]
	Sentinel() *rbtreeNode[T]

	Compare(a, b T) int
	Left(n *rbtreeNode[T]) *rbtreeNode[T]
	Right(n *rbtreeNode[T]) *rbtreeNode[T]
	Reverse() rbtreeInterface[T]
}

func forward[T any](tree *RBTree[T]) rbtreeInterface[T] { return rbtreeForward[T]{tree: tree} }
func reverse[T any](tree *RBTree[T]) rbtreeInterface[T] { return rbtreeReverse[T]{tree: tree} }

type rbtreeForward[T any] struct {
	tree *RBTree[T]
}

func (p rbtreeForward[T]) Tree() *RBTree[T]                      { return p.tree }
func (p rbtreeForward[T]) Root() *rbtreeNode[T]                  { return p.tree.root }
func (p rbtreeForward[T]) Sentinel() *rbtreeNode[T]              { return p.tree.sentinel }
func (p rbtreeForward[T]) Compare(a, b T) int                    { return p.tree.Compare(a, b) }
func (p rbtreeForward[T]) Left(n *rbtreeNode[T]) *rbtreeNode[T]  { return n.Left }
func (p rbtreeForward[T]) Right(n *rbtreeNode[T]) *rbtreeNode[T] { return n.Right }
func (p rbtreeForward[T]) Reverse() rbtreeInterface[T]           { return reverse(p.tree) }

type rbtreeReverse[T any] struct {
	tree *RBTree[T]
}

func (p rbtreeReverse[T]) Tree() *RBTree[T]                      { return p.tree }
func (p rbtreeReverse[T]) Root() *rbtreeNode[T]                  { return p.tree.root }
func (p rbtreeReverse[T]) Sentinel() *rbtreeNode[T]              { return p.tree.sentinel }
func (p rbtreeReverse[T]) Compare(a, b T) int                    { return p.tree.Compare(b, a) }
func (p rbtreeReverse[T]) Left(n *rbtreeNode[T]) *rbtreeNode[T]  { return n.Right }
func (p rbtreeReverse[T]) Right(n *rbtreeNode[T]) *rbtreeNode[T] { return n.Left }
func (p rbtreeReverse[T]) Reverse() rbtreeInterface[T]           { return forward(p.tree) }
