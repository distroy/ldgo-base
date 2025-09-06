/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"testing"
	"unsafe"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldptr"
)

func TestPointer(t *testing.T) {
	vals := []unsafe.Pointer{
		nil, unsafe.Pointer(new(int)), unsafe.Pointer(new(uint)),
	}
	testStoreLoadCompareAndSwapper(t, NewPointer, vals...)
}

func TestPtr(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("int", func(c convey.C) {
			vals := []*int{
				nil, ldptr.New(123),
			}
			conveyStoreLoadCompareAndSwapper(c, NewPtr[int], vals...)
		})
	})
}
