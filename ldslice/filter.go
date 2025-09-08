/*
 * Copyright (C) distroy
 */

package ldslice

// Filter returns the slice with filter(v) == true
func Filter[T any](slice []T, filter func(v T) bool) (trueSlice, falseSlice []T) {
	n := filterSlice(slice, filter)
	return slice[:n], slice[n:]
}

func filterSlice[T any](slice []T, filter func(v T) bool) int {
	i := 0
	j := len(slice)

	for i < j {
		var vi, vj *T

		for ; i < j; i++ {
			vi = &slice[i]
			if !filter(*vi) {
				break
			}
		}

		for ; i < j; j-- {
			vj = &slice[j-1]
			if filter(*vj) {
				break
			}
		}

		if i < j-1 {
			*vi, *vj = *vj, *vi
			i++
		}
	}

	return i
}
