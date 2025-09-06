/*
 * Copyright (C) distroy
 */

package ldatomic

type Byte = Uint8

type String struct {
	Any[string]
}

func NewString(d string) *String {
	p := &String{}
	p.Store(d)
	return p
}

func (p *String) String() string { return p.Load() }
