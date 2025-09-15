/*
 * Copyright (C) distroy
 */

package logref__

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestCheckTypeEqual(t *testing.T) {
	type args struct {
		this reflect.Type
		that reflect.Type
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckTypeEqual(tt.args.this, tt.args.that)
		})
	}
}

func Test_isTypeEqual(t *testing.T) {
	type (
		Any     any
		Int     int
		Uint    uint
		String  string
		Complex complex64
		Slice   []unsafe.Pointer
		Map     map[int64]*string
		Func    func(uint64) error
	)
	type T0 struct {
		Any     any
		Int     int
		Uint    uint
		Str     string
		Complex complex64
		Slice   []unsafe.Pointer
		Map     map[int64]*string
		Func    func(uint64) error
	}

	type T1 struct {
		Any     Any
		Int     Int
		Uint    Uint
		Str     String
		Complex Complex
		Slice   Slice
		Map     Map
		Func    Func
	}

	convey.Convey(t.Name(), t, func(c convey.C) {
		t0 := reflect.TypeOf((*T0)(nil))
		t1 := reflect.TypeOf((*T1)(nil))
		c.So(isTypeEqual(t0, t1), convey.ShouldBeTrue)
	})
}
