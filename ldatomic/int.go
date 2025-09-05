/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"sync/atomic"
)

type Int32 int32

func NewInt32(d int32) *Int32 {
	return (*Int32)(&d)
}

func (p *Int32) get() *int32 { return (*int32)(p) }

func (p *Int32) Add(delta int32) (new int32) { return atomic.AddInt32(p.get(), delta) }
func (p *Int32) Sub(delta int32) (new int32) { return atomic.AddInt32(p.get(), -delta) }
func (p *Int32) Store(d int32)               { atomic.StoreInt32(p.get(), d) }
func (p *Int32) Load() int32                 { return atomic.LoadInt32(p.get()) }
func (p *Int32) Swap(new int32) (old int32)  { return atomic.SwapInt32(p.get(), new) }
func (p *Int32) CompareAndSwap(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(p.get(), old, new)
}

type Int64 int64

func NewInt64(d int64) *Int64 {
	return (*Int64)(&d)
}

func (p *Int64) get() *int64 { return (*int64)(p) }

func (p *Int64) Add(delta int64) (new int64) { return atomic.AddInt64(p.get(), delta) }
func (p *Int64) Sub(delta int64) (new int64) { return atomic.AddInt64(p.get(), -delta) }
func (p *Int64) Store(d int64)               { atomic.StoreInt64(p.get(), d) }
func (p *Int64) Load() int64                 { return atomic.LoadInt64(p.get()) }
func (p *Int64) Swap(new int64) (old int64)  { return atomic.SwapInt64(p.get(), new) }
func (p *Int64) CompareAndSwap(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(p.get(), old, new)
}

type Int int64

func NewInt(d int) *Int {
	v := Int(d)
	return &v
}

func (p *Int) get() *Int64 { return (*Int64)(p) }

func (p *Int) Add(delta int) (new int) { return int(p.get().Add(int64(delta))) }
func (p *Int) Sub(delta int) (new int) { return int(p.get().Sub(int64(delta))) }
func (p *Int) Store(d int)             { p.get().Store(int64(d)) }
func (p *Int) Load() int               { return int(p.get().Load()) }
func (p *Int) Swap(new int) (old int)  { return int(p.get().Swap(int64(new))) }
func (p *Int) CompareAndSwap(old, new int) (swapped bool) {
	return p.get().CompareAndSwap(int64(old), int64(new))
}

type Int8 int32

func NewInt8(d int8) *Int8 {
	v := Int8(d)
	return &v
}

func (p *Int8) get() *Int32 { return (*Int32)(p) }

func (p *Int8) Add(delta int8) (new int8) { return int8(p.get().Add(int32(delta))) }
func (p *Int8) Sub(delta int8) (new int8) { return int8(p.get().Sub(int32(delta))) }
func (p *Int8) Store(d int8)              { p.get().Store(int32(d)) }
func (p *Int8) Load() int8                { return int8(p.get().Load()) }
func (p *Int8) Swap(new int8) (old int8)  { return int8(p.get().Swap(int32(new))) }
func (p *Int8) CompareAndSwap(old, new int8) (swapped bool) {
	return p.get().CompareAndSwap(int32(old), int32(new))
}

type Int16 int32

func NewInt16(d int16) *Int16 {
	v := Int16(d)
	return &v
}

func (p *Int16) get() *Int32 { return (*Int32)(p) }

func (p *Int16) Add(delta int16) (new int16) { return int16(p.get().Add(int32(delta))) }
func (p *Int16) Sub(delta int16) (new int16) { return int16(p.get().Sub(int32(delta))) }
func (p *Int16) Store(d int16)               { p.get().Store(int32(d)) }
func (p *Int16) Load() int16                 { return int16(p.get().Load()) }
func (p *Int16) Swap(new int16) (old int16)  { return int16(p.get().Swap(int32(new))) }
func (p *Int16) CompareAndSwap(old, new int16) (swapped bool) {
	return p.get().CompareAndSwap(int32(old), int32(new))
}
