/*
 * Copyright (C) distroy
 */

package ldatomic

import "sync/atomic"

type Uint32 uint32

func NewUint32(d uint32) *Uint32 {
	return (*Uint32)(&d)
}

func (p *Uint32) get() *uint32 { return (*uint32)(p) }

func (p *Uint32) Add(delta uint32) (new uint32) { return atomic.AddUint32(p.get(), delta) }
func (p *Uint32) Sub(delta uint32) (new uint32) { return atomic.AddUint32(p.get(), ^(delta - 1)) }
func (p *Uint32) Store(d uint32)                { atomic.StoreUint32(p.get(), d) }
func (p *Uint32) Load() uint32                  { return atomic.LoadUint32(p.get()) }
func (p *Uint32) Swap(new uint32) (old uint32)  { return atomic.SwapUint32(p.get(), new) }
func (p *Uint32) CompareAndSwap(old, new uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32(p.get(), old, new)
}

type Uint64 uint64

func NewUint64(d uint64) *Uint64 {
	return (*Uint64)(&d)
}

func (p *Uint64) get() *uint64 { return (*uint64)(p) }

func (p *Uint64) Add(delta uint64) (new uint64) { return atomic.AddUint64(p.get(), delta) }
func (p *Uint64) Sub(delta uint64) (new uint64) { return atomic.AddUint64(p.get(), ^(delta - 1)) }
func (p *Uint64) Store(d uint64)                { atomic.StoreUint64(p.get(), d) }
func (p *Uint64) Load() uint64                  { return atomic.LoadUint64(p.get()) }
func (p *Uint64) Swap(new uint64) (old uint64)  { return atomic.SwapUint64(p.get(), new) }
func (p *Uint64) CompareAndSwap(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(p.get(), old, new)
}

type Uintptr uintptr

func NewUintptr(d uintptr) *Uintptr {
	return (*Uintptr)(&d)
}

func (p *Uintptr) get() *uintptr { return (*uintptr)(p) }

func (p *Uintptr) Add(delta uintptr) (new uintptr) { return atomic.AddUintptr(p.get(), delta) }
func (p *Uintptr) Sub(delta uintptr) (new uintptr) { return atomic.AddUintptr(p.get(), ^(delta - 1)) }
func (p *Uintptr) Store(d uintptr)                 { atomic.StoreUintptr(p.get(), d) }
func (p *Uintptr) Load() uintptr                   { return atomic.LoadUintptr(p.get()) }
func (p *Uintptr) CompareAndSwap(old, new uintptr) (swapped bool) {
	return atomic.CompareAndSwapUintptr(p.get(), old, new)
}

type Uint uint64

func NewUint(d uint) *Uint {
	v := Uint(d)
	return &v
}

func (p *Uint) get() *Uint64 { return (*Uint64)(p) }

func (p *Uint) Add(delta uint) (new uint) { return uint(p.get().Add(uint64(delta))) }
func (p *Uint) Sub(delta uint) (new uint) { return uint(p.get().Sub(uint64(delta))) }
func (p *Uint) Store(d uint)              { p.get().Store(uint64(d)) }
func (p *Uint) Load() uint                { return uint(p.get().Load()) }
func (p *Uint) Swap(new int) (old uint)   { return uint(p.get().Swap(uint64(new))) }
func (p *Uint) CompareAndSwap(old, new uint) (swapped bool) {
	return p.get().CompareAndSwap(uint64(old), uint64(new))
}

type Uint8 uint32

func NewUint8(d uint8) *Uint8 {
	v := Uint8(d)
	return &v
}

func (p *Uint8) get() *Uint32 { return (*Uint32)(p) }

func (p *Uint8) Add(delta uint8) (new uint8) { return uint8(p.get().Add(uint32(delta))) }
func (p *Uint8) Sub(delta uint8) (new uint8) { return uint8(p.get().Sub(uint32(delta))) }
func (p *Uint8) Store(d uint8)               { p.get().Store(uint32(d)) }
func (p *Uint8) Load() uint8                 { return uint8(p.get().Load()) }
func (p *Uint8) Swap(new int8) (old uint8)   { return uint8(p.get().Swap(uint32(new))) }
func (p *Uint8) CompareAndSwap(old, new uint8) (swapped bool) {
	return p.get().CompareAndSwap(uint32(old), uint32(new))
}

type Uint16 uint32

func NewUint16(d uint16) *Uint16 {
	v := Uint16(d)
	return &v
}

func (p *Uint16) get() *Uint32 { return (*Uint32)(p) }

func (p *Uint16) Add(delta uint16) (new uint16) { return uint16(p.get().Add(uint32(delta))) }
func (p *Uint16) Sub(delta uint16) (new uint16) { return uint16(p.get().Sub(uint32(delta))) }
func (p *Uint16) Store(d uint16)                { p.get().Store(uint32(d)) }
func (p *Uint16) Load() uint16                  { return uint16(p.get().Load()) }
func (p *Uint16) Swap(new int16) (old uint16)   { return uint16(p.get().Swap(uint32(new))) }
func (p *Uint16) CompareAndSwap(old, new uint16) (swapped bool) {
	return p.get().CompareAndSwap(uint32(old), uint32(new))
}
