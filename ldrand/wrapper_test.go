/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"math"
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/lditer"
)

func testFor[T any](seq lditer.Seq[T], do func(v T)) {
	for v := range seq {
		do(v)
	}
}

func testPrintGroup(c convey.C, group, idx int, value any) {
	if (idx % group) == 0 {
		c.Print("\n      ", value)
	} else {
		c.Print(" ", value)
	}
}

func TestNew(t *testing.T) {
	times := 1024
	retry := func(retry int, f func(i int)) {
		testFor(lditer.Int(retry), f)
	}
	convey.Convey(t.Name(), t, func(c convey.C) {
		r := New(NewFastSource(time.Now().UnixNano()))

		c.Convey(`int`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Int(), convey.ShouldBeBetweenOrEqual, 0, math.MaxInt) })
		})
		c.Convey(`int31`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Int31(), convey.ShouldBeBetweenOrEqual, 0, math.MaxInt32) })
		})
		c.Convey(`int63`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Int63(), convey.ShouldBeBetweenOrEqual, 0, math.MaxInt64) })
		})

		c.Convey(`uint`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Uint(), convey.ShouldBeBetweenOrEqual, 0, uint(math.MaxUint)) })
		})
		c.Convey(`uint32`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Uint32(), convey.ShouldBeBetweenOrEqual, 0, uint32(math.MaxUint32)) })
		})
		c.Convey(`uint64`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Uint64(), convey.ShouldBeBetweenOrEqual, 0, uint64(math.MaxUint64)) })
		})

		c.Convey(`intn`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Intn(100), convey.ShouldBeBetweenOrEqual, 0, 100) })
		})
		c.Convey(`int31n`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Int31n(100), convey.ShouldBeBetweenOrEqual, 0, 100) })
		})
		c.Convey(`int63n`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Int63n(100), convey.ShouldBeBetweenOrEqual, 0, 100) })
		})

		c.Convey(`float32`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Float32(), convey.ShouldBeBetweenOrEqual, 0.0, 1.0) })
		})
		c.Convey(`float64`, func(c convey.C) {
			retry(times, func(i int) { c.So(r.Float64(), convey.ShouldBeBetweenOrEqual, 0.0, 1.0) })
		})

		c.Convey(`perm`, func(c convey.C) {
			n := 16
			l := r.Perm(n)
			c.Print(l)
			testFor(lditer.Int(len(l)), func(i int) {
				c.So(l, convey.ShouldContain, i)
			})
		})
		c.Convey(`shuffle`, func(c convey.C) {
			n := 16
			l := make([]int, 0, n)
			testFor(lditer.Int(len(l)), func(i int) {
				l = append(l, i)
			})
			r.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })
			c.Print(l)
			testFor(lditer.Int(len(l)), func(i int) {
				c.So(l, convey.ShouldContain, i)
			})
		})

		c.Convey(`read`, func(c convey.C) {
			b := make([]byte, 1024)
			r.Read(b)
			c.Print(b)
			testFor(lditer.ToSeqByValue(lditer.Slice(b)), func(v byte) {
				c.So(v, convey.ShouldBeBetweenOrEqual, 0, 255)
			})
		})
		c.Convey(`bytes`, func(c convey.C) {
			b := r.Bytes(1024)
			c.Print(b)
			testFor(lditer.ToSeqByValue(lditer.Slice(b)), func(v byte) {
				c.So(v, convey.ShouldBeBetweenOrEqual, 0, 255)
			})
		})
		c.Convey(`string`, func(c convey.C) {
			s := r.String(1024)
			c.Print(s)
			testFor(lditer.ToSeqByValue(lditer.String(s)), func(v rune) {
				c.So(v, convey.ShouldBeBetweenOrEqual, 0, 255)
			})
		})

		c.Convey(`norm float64`, func(c convey.C) {
			retry(times, func(i int) {
				// 生成一个均值为 0，标准差为 1 的正态分布随机数
				n := r.NormFloat64()*1 + 0
				testPrintGroup(c, 6, i, n)
				c.So(n, convey.ShouldBeBetweenOrEqual, -128, 128)
			})
		})
		c.Convey(`exp float64`, func(c convey.C) {
			retry(times, func(i int) {
				// 生成一个均值为 0，标准差为 1 的正态分布随机数
				n := r.ExpFloat64()*1 + 0
				testPrintGroup(c, 6, i, n)
				c.So(n, convey.ShouldBeBetweenOrEqual, -128, 128)
			})
		})
	})
}
