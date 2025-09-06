/*
 * Copyright (C) distroy
 */

package ldatomic

import "math"

type Float32 struct {
	d Uint32
}

func NewFloat32(d float32) *Float32 {
	return &Float32{d: Uint32(math.Float32bits(d))}
}

func (p *Float32) Store(d float32) { p.d.Store(math.Float32bits(d)) }
func (p *Float32) Load() float32   { return math.Float32frombits(p.d.Load()) }
func (p *Float32) Swap(new float32) (old float32) {
	return math.Float32frombits(p.d.Swap(math.Float32bits(new)))
}
func (p *Float32) CompareAndSwap(old, new float32) bool {
	return p.d.CompareAndSwap(math.Float32bits(old), math.Float32bits(new))
}
func (p *Float32) Add(delta float32) (new float32) {
	for {
		old := p.Load()
		new := old + delta
		if p.CompareAndSwap(old, new) {
			return new
		}
	}
}
func (p *Float32) Sub(delta float32) (new float32) {
	return p.Add(-delta)
}

type Float64 struct {
	d Uint64
}

func NewFloat64(d float64) *Float64 {
	return &Float64{d: Uint64(math.Float64bits(d))}
}

func (p *Float64) Store(d float64) { p.d.Store(math.Float64bits(d)) }
func (p *Float64) Load() float64   { return math.Float64frombits(p.d.Load()) }
func (p *Float64) Swap(new float64) (old float64) {
	return math.Float64frombits(p.d.Swap(math.Float64bits(new)))
}
func (p *Float64) CompareAndSwap(old, new float64) bool {
	return p.d.CompareAndSwap(math.Float64bits(old), math.Float64bits(new))
}
func (p *Float64) Add(delta float64) (new float64) {
	for {
		old := p.Load()
		new := old + delta
		if p.CompareAndSwap(old, new) {
			return new
		}
	}
}
func (p *Float64) Sub(delta float64) (new float64) {
	return p.Add(-delta)
}
