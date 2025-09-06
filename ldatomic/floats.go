/*
 * Copyright (C) distroy
 */

package ldatomic

type (
	Float32s = Any[[]float32]
	Float64s = Any[[]float64]
)

func NewFloat32s(d []float32) *Float32s {
	p := &Float32s{}
	if d != nil {
		p.Store(d)
	}
	return p
}

func NewFloat64s(d []float64) *Float64s {
	p := &Float64s{}
	if d != nil {
		p.Store(d)
	}
	return p
}
