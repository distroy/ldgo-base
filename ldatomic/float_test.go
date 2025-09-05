/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"math"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestFloat32(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewFloat32(1).Load(), convey.ShouldResemble, float32(1))

		p := Float32{}
		convey.So(p.Load(), convey.ShouldResemble, float32(0))

		p.Store(1)
		convey.So(p.Load(), convey.ShouldResemble, float32(1))
		p.Store(2)
		convey.So(p.Load(), convey.ShouldResemble, float32(2))
		p.Store(0)
		convey.So(p.Load(), convey.ShouldResemble, float32(0))

		convey.So(p.Swap(1), convey.ShouldResemble, float32(0))
		convey.So(p.Swap(2), convey.ShouldResemble, float32(1))
		convey.So(p.Swap(0), convey.ShouldResemble, float32(2))

		convey.So(p.CompareAndSwap(0, 1), convey.ShouldResemble, true)
		convey.So(p.Load(), convey.ShouldResemble, float32(1))

		convey.So(p.CompareAndSwap(0, 2), convey.ShouldResemble, false)
		convey.So(p.Load(), convey.ShouldResemble, float32(1))
		convey.So(p.CompareAndSwap(1, 2), convey.ShouldResemble, true)
		convey.So(p.Load(), convey.ShouldResemble, float32(2))

		convey.So(p.CompareAndSwap(1, 0), convey.ShouldResemble, false)
		convey.So(p.Load(), convey.ShouldResemble, float32(2))
		convey.So(p.CompareAndSwap(2, 0), convey.ShouldResemble, true)
		convey.So(p.Load(), convey.ShouldResemble, float32(0))
		convey.So(p.Sub(-1), convey.ShouldResemble, float32(1))
	})
}

func TestFloat64(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewFloat64(1).Load(), convey.ShouldResemble, float64(1))

		p := Float64{}
		convey.So(p.Load(), convey.ShouldResemble, float64(0))

		p.Store(1)
		convey.So(p.Load(), convey.ShouldResemble, float64(1))
		p.Store(2)
		convey.So(p.Load(), convey.ShouldResemble, float64(2))
		p.Store(0)
		convey.So(p.Load(), convey.ShouldResemble, float64(0))

		convey.So(p.Swap(1), convey.ShouldResemble, float64(0))
		convey.So(p.Swap(2), convey.ShouldResemble, float64(1))
		convey.So(p.Swap(0), convey.ShouldResemble, float64(2))

		convey.So(p.CompareAndSwap(0, 1), convey.ShouldResemble, true)
		convey.So(p.Load(), convey.ShouldResemble, float64(1))

		convey.So(p.CompareAndSwap(0, 2), convey.ShouldResemble, false)
		convey.So(p.Load(), convey.ShouldResemble, float64(1))
		convey.So(p.CompareAndSwap(1, 2), convey.ShouldResemble, true)
		convey.So(p.Load(), convey.ShouldResemble, float64(2))

		convey.So(p.CompareAndSwap(1, 0), convey.ShouldResemble, false)
		convey.So(p.Load(), convey.ShouldResemble, float64(2))
		convey.So(p.CompareAndSwap(2, 0), convey.ShouldResemble, true)
		convey.So(p.Load(), convey.ShouldResemble, float64(0))
		convey.So(p.Sub(-1), convey.ShouldResemble, float64(1))

		p.Store(math.NaN())
		convey.So(math.IsNaN(p.Load()), convey.ShouldResemble, true)

		p.Store(math.Inf(1))
		convey.So(math.IsInf(p.Load(), 1), convey.ShouldResemble, true)
		p.Store(math.Inf(-1))
		convey.So(math.IsInf(p.Load(), -1), convey.ShouldResemble, true)
	})
}
