/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

type testInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type testIntIface[T testInteger] interface {
	StoreLoadCompareAndSwapper[T]

	Add(delta T) (new T)
	Sub(delta T) (new T)
	// Store(d T)
	// Load() T
	// Swap(new T) (old T)
	// CompareAndSwap(old, new T) (swapped bool)
}

func testInt[I testIntIface[T], T testInteger](t *testing.T, fNew func(T) I) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		p := fNew(123)
		c.So(p.Load(), convey.ShouldEqual, T(123))

		p.Store(12)
		c.So(p.Load(), convey.ShouldEqual, T(12))

		c.So(p.Add(5), convey.ShouldEqual, T(17))
		c.So(p.Sub(4), convey.ShouldEqual, T(13))

		c.So(p.Swap(7), convey.ShouldEqual, T(13))

		c.So(p.Load(), convey.ShouldEqual, T(7))
		c.So(p.CompareAndSwap(3, 17), convey.ShouldEqual, false)
		c.So(p.Load(), convey.ShouldEqual, T(7))
		c.So(p.CompareAndSwap(7, 17), convey.ShouldEqual, true)
		c.So(p.Load(), convey.ShouldEqual, T(17))
	})
}

func TestInt(t *testing.T)   { testInt(t, NewInt) }
func TestInt8(t *testing.T)  { testInt(t, NewInt8) }
func TestInt16(t *testing.T) { testInt(t, NewInt16) }
func TestInt32(t *testing.T) { testInt(t, NewInt32) }
func TestInt64(t *testing.T) { testInt(t, NewInt64) }
