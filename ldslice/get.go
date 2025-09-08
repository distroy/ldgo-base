/*
 * Copyright (C) distroy
 */

package ldslice

func Get[T any](s []T, idx int, def ...T) T {
	if len(s) > idx {
		return s[idx]
	}
	if len(def) > 0 {
		return def[0]
	}
	var v T
	return v
}
