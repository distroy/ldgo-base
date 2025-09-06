/*
 * Copyright (C) distroy
 */

package ldatomic

type (
	Uints    = Any[[]uint]
	Uint8s   = Any[[]uint8]
	Uint16s  = Any[[]uint16]
	Uint32s  = Any[[]uint32]
	Uint64s  = Any[[]uint64]
	Uintptrs = Any[[]uintptr]
)

func NewUints(d []uint) *Uints          { return NewAny[[]uint](d) }
func NewUint8s(d []uint8) *Uint8s       { return NewAny[[]uint8](d) }
func NewUint16s(d []uint16) *Uint16s    { return NewAny[[]uint16](d) }
func NewUint32s(d []uint32) *Uint32s    { return NewAny[[]uint32](d) }
func NewUint64s(d []uint64) *Uint64s    { return NewAny[[]uint64](d) }
func NewUintptrs(d []uintptr) *Uintptrs { return NewAny[[]uintptr](d) }
