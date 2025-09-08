/*
 * Copyright (C) distroy
 */

package ldtopk

type sortable interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func less[T sortable](a, b T) bool {
	return a < b
}

func TopkAdd[T sortable](b []T, k int, x T) ([]T, bool) { return topkAdd(b, k, x, less[T]) }

func topkAdd[T any](b []T, k int, x T, less LessFunc[T]) ([]T, bool) {
	if k <= 0 {
		return b, false
	}

	if pos := len(b); pos < k {
		topkAppendTail(&b, less, x)
		return b, true
	}

	if !less(x, b[0]) {
		return b, false
	}

	b[0] = x
	topkFixupHead(&b, less)
	return b, true
}

func topkAppendTail[T any](heap *[]T, less LessFunc[T], d T) {
	pos := len(*heap)

	*heap = append(*heap, d)
	for parent := 0; pos > 0; pos = parent {
		parent = (pos - 1) / 2
		if !less((*heap)[parent], (*heap)[pos]) {
			break
		}
		(*heap)[parent], (*heap)[pos] = (*heap)[pos], (*heap)[parent]
	}
}

func topkFixupHead[T any](heap *[]T, less LessFunc[T]) {
	size := len(*heap)
	for pos := 0; ; {
		lChild := (pos * 2) + 1
		rChild := (pos * 2) + 2
		if lChild >= size {
			break
		}

		child := lChild
		if rChild < size && less((*heap)[lChild], (*heap)[rChild]) {
			child = rChild
		}

		if !less((*heap)[pos], (*heap)[child]) {
			break
		}

		(*heap)[pos], (*heap)[child] = (*heap)[child], (*heap)[pos]
		pos = child
	}
}
