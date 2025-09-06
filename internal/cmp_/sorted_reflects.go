/*
 * Copyright (C) distroy
 */

package cmp_

import "reflect"

type sortedReflects []reflect.Value

func (o sortedReflects) Len() int           { return len(o) }
func (o sortedReflects) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o sortedReflects) Less(i, j int) bool { return CompareReflect(o[i], o[j]) <= 0 }
