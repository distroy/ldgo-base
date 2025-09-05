/*
 * Copyright (C) distroy
 */

package ldatomic

// Interface provides an atomic load and store of an any typed value.
// The zero value for a Value returns nil from Load.
type Interface = Any[interface{}]

func NewInterface(d interface{}) *Interface {
	p := &Interface{}
	p.Store(d)
	return p
}

type Error = Any[error]

func NewError(d error) *Error {
	p := &Error{}
	p.Store(d)
	return p
}
