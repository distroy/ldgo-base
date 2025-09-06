/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"reflect"
)

func equal(a, b reflect.Value) bool {
	if a.Kind() == reflect.Interface {
		a = a.Elem()
	}
	if b.Kind() == reflect.Interface {
		b = b.Elem()
	}

	if !a.IsValid() || !b.IsValid() {
		return a.IsValid() == b.IsValid()
	}

	if a.Kind() != b.Kind() || a.Type() != b.Type() {
		return false
	}

	switch a.Kind() {
	case reflect.Invalid:
		return true

	case reflect.Bool:
		return a.Bool() == b.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() == b.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() == b.Uint()

	case reflect.Float32, reflect.Float64:
		return a.Float() == b.Float()

	case reflect.Complex64, reflect.Complex128:
		return a.Complex() == b.Complex()

	case reflect.String:
		return a.String() == b.String()

	case reflect.Chan, reflect.Pointer, reflect.UnsafePointer:
		return a.Pointer() == b.Pointer()

	case reflect.Func:
		return a.Pointer() == b.Pointer()

	case reflect.Array:
		return equalArray(a, b)

	case reflect.Struct:
		return equalStruct(a, b)

	case reflect.Map:
		return equalMap(a, b)

	case reflect.Slice:
		return equalSlice(a, b)
	}

	return false
}

func equalArray(a, b reflect.Value) bool {
	al := a.Len()
	if al == 0 {
		return true
	}
	for i := 0; i < al; i++ {
		if !equal(a.Index(i), b.Index(i)) {
			return false
		}
	}
	return true

}

func equalStruct(a, b reflect.Value) bool {
	// a and a have the same type so they have the same fields
	nf := a.NumField()
	for i := 0; i < nf; i++ {
		if !equal(a.Field(i), b.Field(i)) {
			return false
		}
	}
	return true
}

func equalSlice(a, b reflect.Value) bool {
	if a.IsNil() != b.IsNil() {
		return false
	}

	al := a.Len()
	bl := b.Len()
	if al != bl {
		return false
	}
	if al == 0 {
		return true
	}

	// If v's Kind is Slice, the returned pointer is to the first element of the slice
	return a.Pointer() == b.Pointer()
}

func equalMap(a, b reflect.Value) bool {
	if a.IsNil() != b.IsNil() {
		return false
	}

	al := a.Len()
	bl := b.Len()
	if al != bl {
		return false
	}
	if al == 0 {
		return true
	}

	return a.Pointer() == b.Pointer()

	// ai := a.Interface()
	// bi := b.Interface()
	//
	// // log.Printf(" *** ai: %#v", ai)
	// // log.Printf(" *** bi: %#v", bi)
	//
	// ap := (*efaceWords)(unsafe.Pointer(&ai))
	// bp := (*efaceWords)(unsafe.Pointer(&bi))
	//
	// // log.Printf(" *** ap: %#v", ap)
	// // log.Printf(" *** bp: %#v", bp)
	//
	// return *ap == *bp

	// return a.Index(0).Addr() == b.Index(0).Addr()
}
