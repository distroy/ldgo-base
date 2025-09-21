/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/lditer"
	"github.com/distroy/ldgo-base/ldmath"
)

/*
 * pkg: github.com/distroy/ldgo/v2/ldrand
 * cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
 * BenchmarkRandGo
 * BenchmarkRandGo-12      18746797                63.71 ns/op
 * BenchmarkRand
 * BenchmarkRand-12        68977040                17.42 ns/op
 */
func BenchmarkRandGo(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Intn(100)
		}
	})
}

func BenchmarkRand(b *testing.B) {
	r := New(NewFastSource(rand.Int63()))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Intn(100)
		}
	})
}

type testFastSource struct {
	Mod, Scale, Diff int
}

func maxInt(a []int) int { return ldmath.Max(a[0], a[1:]...) }
func minInt(a []int) int { return ldmath.Min(a[0], a[1:]...) }

func diffRatio(a []int) float64 {
	sum := ldmath.Sum2[int64](a...)
	cnt := int64(len(a))
	avg := (sum + cnt/2) / cnt

	diff := int64(0)
	for _, n := range a {
		diff += ldmath.Abs(avg - int64(n))
	}

	return float64(diff) / float64(sum)
}

func (t *testFastSource) Test() {
	var (
		mod   = t.Mod
		scale = t.Scale
		diff  = t.Diff
	)
	name := fmt.Sprintf("mod=%d,scale=%d,diff=%d", mod, scale, diff)
	convey.Convey(name, func(c convey.C) {
		r := New(NewFastSource(time.Now().UnixNano()))

		counts := make([]int, mod)
		for i := 0; i < mod*scale; i++ {
			// x := r.Int() % mod
			x := r.Intn(mod)
			counts[x]++
		}

		min := minInt(counts)
		max := maxInt(counts)
		ratio := diffRatio(counts)

		log.Printf("diff:%d, ratio:%.04g, min:%d, max:%d", max-min, ratio, min, max)
		convey.So(max-min, convey.ShouldBeLessThan, diff)
	})
}

func Test_fastSource_ProbabilityOfOverall(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		(&testFastSource{
			Mod:   100,
			Scale: 1000 * 200,
			Diff:  1000 * 4,
		}).Test()
		(&testFastSource{
			Mod:   16,
			Scale: 1000 * 200,
			Diff:  1000 * 4,
		}).Test()
		(&testFastSource{
			Mod:   256,
			Scale: 1000 * 200,
			Diff:  1000 * 4,
		}).Test()
	})
}

func Test_fastSource_ProbabilityOfVery4Bits(t *testing.T) {
	r := New(NewFastSource(time.Now().UnixNano()))

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("check the probability of very 4 bits", func(c convey.C) {
			const (
				scale = 1000 * 100
				diff  = 1000 * 4
			)

			countsPer4Bits := [16][16]int{}
			for range scale * 16 {
				v := r.Uint64()
				for i := range countsPer4Bits {
					countsPer4Bits[i][v&0xf]++
					v = v >> 4
				}
			}

			log.Printf("")
			for i, v := range countsPer4Bits {
				min := minInt(v[:])
				max := maxInt(v[:])
				ratio := diffRatio(v[:])
				log.Printf("postion:%d, diff:%d, ratio:%.04g, min:%d, max:%d",
					i, max-min, ratio, min, max)
				convey.So(max-min, convey.ShouldBeLessThan, diff)
			}
		})
	})
}
func testGetDiffThresholdBySliceFunc(diff1, diff2 int) func(idx int, slice any) int {
	return func(idx int, slice any) int {
		sliceV := reflect.ValueOf(slice)
		if idx == sliceV.Len()-1 {
			return diff2
		}
		return diff1
	}
}

func Test_fastSource_ProbabilityOfVery4BitsWithPreviousNumber(t *testing.T) {
	r := New(NewFastSource(time.Now().UnixNano()))
	const (
		scale = 1000 * 100 * 16 * 16
		diff1 = 5000
		diff2 = 5000
		// diff2 = 6500
	)

	getDiffThreshold := testGetDiffThresholdBySliceFunc(diff1, diff2)

	countsPer4BitsWithPrev := [16][16][16]int{}
	runCount := 0
	run := func() {
		runCount++
		prevNum := r.Uint64()
		for range scale {
			v := r.Uint64()
			p := prevNum
			prevNum = v
			for i := range countsPer4BitsWithPrev {
				countsPer4BitsWithPrev[i][p&0xf][v&0xf]++
				v = v >> 4
				p = p >> 4
			}
		}
	}

	check := func() bool {
		return testFor2(lditer.Slice(countsPer4BitsWithPrev[:]), func(i int, v [16][16]int) bool {
			return testFor2(lditer.Slice(v[:]), func(_ int, w [16]int) bool {
				min := minInt(w[:])
				max := maxInt(w[:])

				diff := getDiffThreshold(i, countsPer4BitsWithPrev[:])
				if (max - min) >= (diff * runCount) {
					return false
				}
				return true
			})
		})
	}

	convey.Convey(t.Name(), t, func(c convey.C) {
		run()
		if !check() {
			run()
			run()
		}

		c.Printf("\nrun times:%d, rand times:%d", runCount, runCount*scale)
		testFor2(lditer.Slice(countsPer4BitsWithPrev[:]), func(i int, v [16][16]int) bool {
			// return false
			testFor2(lditer.Slice(v[:]), func(j int, w [16]int) bool {
				min := minInt(w[:])
				max := maxInt(w[:])
				ratio := diffRatio(w[:])
				c.Printf("\npostion:%d, prev:%d, diff:%d, ratio:%.04g, min:%d, max:%d, v:%v",
					i, j, max-min, ratio, min, max, w[:])

				diff := getDiffThreshold(i, countsPer4BitsWithPrev[:])
				// sometimes it will be failed
				c.So(max-min, convey.ShouldBeLessThan, diff*runCount)
				return true
			})
			return true
		})
	})
}

func Test_fastSource_ProbabilityOfVeryByte(t *testing.T) {
	r := New(NewFastSource(time.Now().UnixNano()))
	convey.Convey(t.Name(), t, func(c convey.C) {
		const (
			scale = 1000 * 100
			diff  = 1000 * 5
		)

		countsPer4Bits := [8][256]int{}
		for range scale * 256 {
			v := r.Uint64()
			for i := range countsPer4Bits {
				countsPer4Bits[i][v&0xff]++
				v = v >> 8
			}
		}

		log.Printf("")
		for i, v := range countsPer4Bits {
			min := minInt(v[:])
			max := maxInt(v[:])
			ratio := diffRatio(v[:])
			log.Printf("postion:%d, diff:%d, ratio:%.04g, min:%d, max:%d",
				i, max-min, ratio, min, max)
			convey.So(max-min, convey.ShouldBeLessThan, diff)
		}
	})
}

func Test_fastSource_Repeated(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		r := New(NewFastSource(time.Now().UnixNano()))

		c.Convey("check result if repeated", func(c convey.C) {
			const times = 100 * 10000

			m := make(map[uint64]struct{}, times)
			for range times {
				x := r.Uint64()
				if _, ok := m[x]; ok {
					t.Fatalf("number repeated. %d", x)
				}
				m[x] = struct{}{}
			}
		})
	})
}

func testGreatestCommonDivisor(a, b uint64) uint64 {
	for b != 0 {
		a %= b
		a, b = b, a
	}

	return a
}

func Test_testGreatestCommonDivisor(t *testing.T) {
	type Case struct {
		A, B, R uint64
	}
	cases := []*Case{
		{A: 0, B: 2, R: 2},
		{A: 2, B: 0, R: 2},
		{A: 2, B: 100, R: 2},
		{A: 2, B: 3, R: 1},
		{A: 10, B: 6, R: 2},
		{A: 6, B: 10, R: 2},
	}
	convey.Convey(t.Name(), t, func(c convey.C) {
		for _, v := range cases {
			name := fmt.Sprintf("a:%d, b:%d", v.A, v.B)
			c.Convey(name, func(c convey.C) {
				r := testGreatestCommonDivisor(v.A, v.B)
				c.So(r, convey.ShouldEqual, v.R)
			})
		}
	})
}

func Test__fastSourceStep(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		r := testGreatestCommonDivisor(math.MaxUint64, fastSourceStep)
		c.So(r, convey.ShouldEqual, 1)
	})
}
