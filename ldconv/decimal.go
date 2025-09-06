/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"math/big"

	"github.com/distroy/ldgo-base/3rd/decimal"
)

type decimalNumber = decimal.Decimal

var decimalZero = newDecimalFromInt(0)

func newDecimalZero() decimalNumber {
	return decimalZero
}

func newDecimalFromFloat(n float64) decimalNumber {
	return decimal.NewFromFloat(n)
}

func newDecimalFromFloatStr(s string) (decimalNumber, error) {
	return decimal.NewFromString(s)
}

func mustNewDecimalFromStr(s string) decimalNumber {
	n, _ := decimal.NewFromString(s)
	return n
}

func newDecimalFromBigFloat(n *big.Float) decimalNumber {
	r, _ := decimal.NewFromString(n.String())
	return r
}

func newDecimalFromInt(n int64) decimalNumber {
	return decimal.NewFromInt(n)
}

func newDecimalFromUint(n uint64) decimalNumber {
	return decimal.NewFromBigInt(new(big.Int).SetUint64(uint64(n)), 0)
}

func newDecimalFromIntStr(s string, negative bool, base int) (decimalNumber, error) {
	if len(s) == 0 {
		return newDecimalZero(), nil
	}
	i, _ := new(big.Int).SetString(s, base)
	if negative {
		i = i.Neg(i)
	}
	return decimal.NewFromBigInt(i, 0), nil
}

type internalDecimal decimalNumber

func (p *internalDecimal) get() *decimalNumber { return (*decimalNumber)(p) }

func (p *internalDecimal) String() string {
	if p == nil {
		return ""
	}
	return p.get().String()
}

func (p *internalDecimal) Bytes() []byte {
	if p == nil {
		return nil
	}
	return StrToBytes(p.get().String())
}

func (p *internalDecimal) IsZero() bool {
	if p == nil {
		return true
	}
	return p.get().IsZero()
}

func (p *internalDecimal) Int64() int64 {
	if p == nil {
		return 0
	}
	return p.get().BigInt().Int64()
}

func (p *internalDecimal) Uint64() uint64 {
	if p == nil {
		return 0
	}
	return p.get().BigInt().Uint64()
}

func (p *internalDecimal) Float32() float32 {
	if p == nil {
		return 0
	}

	f, _ := p.get().Rat().Float32()
	return f
}

func (p *internalDecimal) Float64() float64 {
	if p == nil {
		return 0
	}

	f, _ := p.get().Rat().Float64()
	return f
}

func (p *internalDecimal) BigFloat() *big.Float {
	if p == nil {
		return nil
	}

	return p.get().BigFloat()
}
