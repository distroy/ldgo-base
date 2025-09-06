/*
 * Copyright (C) distroy
 */

package ldconv

import "math/big"

type bigFloat big.Float

func (p *bigFloat) get() *big.Float { return (*big.Float)(p) }

func (p *bigFloat) String() string {
	if p == nil {
		return ""
	}

	return p.get().String()
}

func (p *bigFloat) Bytes() []byte {
	if p == nil {
		return nil
	}

	return StrToBytes(p.get().String())
}

func (p *bigFloat) IsZero() bool {
	if p == nil {
		return true
	}
	return p.get().Cmp(newBigFloatZero()) == 0
}

func (p *bigFloat) Int64() int64 {
	if p == nil {
		return 0
	}
	r, _ := p.get().Int64()
	return r
}

func (p *bigFloat) Uint64() uint64 {
	if p == nil {
		return 0
	}
	r, _ := p.get().Uint64()
	return r
}

func (p *bigFloat) Float32() float32 {
	if p == nil {
		return 0
	}
	r, _ := p.get().Float32()
	return r
}

func (p *bigFloat) Float64() float64 {
	if p == nil {
		return 0
	}
	r, _ := p.get().Float64()
	return r
}
