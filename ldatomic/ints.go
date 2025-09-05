/*
 * Copyright (C) distroy
 */

package ldatomic

type (
	Ints   = Any[[]int]
	Int8s  = Any[[]int8]
	Int16s = Any[[]int16]
	Int32s = Any[[]int32]
	Int64s = Any[[]int64]
)

func NewInts(d []int) *Ints       { return NewAny[[]int](d) }
func NewInt8s(d []int8) *Int8s    { return NewAny[[]int8](d) }
func NewInt16s(d []int16) *Int16s { return NewAny[[]int16](d) }
func NewInt32s(d []int32) *Int32s { return NewAny[[]int32](d) }
func NewInt64s(d []int64) *Int64s { return NewAny[[]int64](d) }
