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

type Ptr[T any] struct {
	d unsafe.Pointer
}

func NewPtr[T any](d *T) *Ptr[T] {
	return &Ptr[T]{d: unsafe.Pointer(d)}
}

func (p *Ptr[T]) get() *unsafe.Pointer { return &p.d }

func (p *Ptr[T]) Store(d *T) { atomic.StorePointer(p.get(), p.pack(d)) }
func (p *Ptr[T]) Load() *T   { return p.unpack(atomic.LoadPointer(p.get())) }
func (p *Ptr[T]) Swap(new *T) (old *T) {
	return p.unpack(atomic.SwapPointer(p.get(), p.pack(new)))
}
func (p *Ptr[T]) CompareAndSwap(old, new *T) (swapped bool) {
	return atomic.CompareAndSwapPointer(p.get(), p.pack(old), p.pack(new))
}
func (p *Ptr[T]) pack(d *T) unsafe.Pointer   { return unsafe.Pointer(d) }
func (p *Ptr[T]) unpack(d unsafe.Pointer) *T { return (*T)(d) }
