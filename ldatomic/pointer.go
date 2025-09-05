/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"sync/atomic"
	"unsafe"
)

type Pointer struct {
	d unsafe.Pointer
}

func NewPointer(d unsafe.Pointer) *Pointer {
	return &Pointer{d: d}
}

func (p *Pointer) get() *unsafe.Pointer { return &p.d }

func (p *Pointer) Store(d unsafe.Pointer) { atomic.StorePointer(p.get(), d) }
func (p *Pointer) Load() unsafe.Pointer   { return atomic.LoadPointer(p.get()) }
func (p *Pointer) Swap(new unsafe.Pointer) (old unsafe.Pointer) {
	return atomic.SwapPointer(p.get(), new)
}
func (p *Pointer) CompareAndSwap(old, new unsafe.Pointer) (swapped bool) {
	return atomic.CompareAndSwapPointer(p.get(), old, new)
}
