/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"reflect"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func testFuncForConvert(_ ...any) int { return 1001 }

type testMethodForConvert struct{}

func (_ *testMethodForConvert) Test(_ ...any) any {
	return 2001
}

func TestHookConvert(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := NewPatches()
		defer patches.Reset() // check double reset

		convey.Convey("FuncHook", func() {
			patches.Apply(FuncHook{
				Target: testFuncForConvert,
				Double: Values{int64(1101)},
			})
			convey.So(testFuncForConvert(1, ""), convey.ShouldEqual, 1101)
			convey.So(testFuncForConvert(1, ""), convey.ShouldEqual, 1101)

			patches.Apply(FuncHook{
				Target: testFuncForConvert,
				Double: Values{int64(1201)},
			})
			convey.So(testFuncForConvert(1, ""), convey.ShouldEqual, 1201)

			patches.Reset()
			convey.So(testFuncForConvert(1, ""), convey.ShouldEqual, 1001)
		})
		convey.Convey("MethodHook", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForConvert)(nil),
				Method: "Test",
				Double: Values{1701},
			})
			obj := &testMethodForConvert{}
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 1701)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 1701)
		})
		convey.Convey("MethodHook(reflect.Type)", func() {
			patches.Apply(MethodHook{
				Target: reflect.TypeOf((*testMethodForConvert)(nil)),
				Method: "Test",
				Double: Values{1701},
			})
			obj := &testMethodForConvert{}
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 1701)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 1701)
		})
	})
}
func testFuncForHook(_ ...any) any { return 1001 }

func TestFuncHook(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := NewPatches()
		defer patches.Reset()

		convey.Convey("function", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: func(_ ...any) any {
					return 1101
				},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1101)
		})
		convey.Convey("[]OutputCell", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: []OutputCell{
					{
						Values: Values{(1201)},
					},
					{
						Times:  2,
						Values: Values{1202},
					},
					{
						Times:  1,
						Values: Values{1203},
					},
				},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1201)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1202)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1202)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1203)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1203)
		})
		convey.Convey("[]*OutputCell", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: []*OutputCell{
					{
						Times:  2,
						Values: Values{1301},
					},
					{
						Times:  0,
						Values: Values{1302},
					},
					{
						Times:  1,
						Values: Values{1303},
					},
				},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1301)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1301)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1302)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1303)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1303)
		})
		convey.Convey("OutputCell", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: OutputCell{
					Times:  2,
					Values: Values{1401},
				},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1401)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1401)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1401)
		})
		convey.Convey("*OutputCell", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: &OutputCell{
					Times:  2,
					Values: Values{1501},
				},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1501)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1501)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1501)
		})
		convey.Convey("[]Values", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: []Values{
					{1601},
					{1602},
					{1603},
				},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1601)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1602)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1603)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1603)
		})
		convey.Convey("Values", func() {
			patches.Apply(FuncHook{
				Target: testFuncForHook,
				Double: Values{1701},
			})
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1701)
			convey.So(testFuncForHook(1, ""), convey.ShouldEqual, 1701)
		})
	})
}

type testMethodForHook struct{}

func (_ *testMethodForHook) Test(_ ...any) any {
	return 2001
}

func TestMethodHook(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := NewPatches()
		obj := &testMethodForHook{}
		defer patches.Reset()

		convey.Convey("function", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: func(_ *testMethodForHook, _ ...any) any {
					return 2101
				},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2101)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2101)
		})
		convey.Convey("[]OutputCell", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: []OutputCell{
					{
						Values: Values{2201},
					},
					{
						Times:  2,
						Values: Values{2202},
					},
					{
						Times:  1,
						Values: Values{2203},
					},
				},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2201)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2202)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2202)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2203)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2203)
		})
		convey.Convey("[]*OutputCell", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: []*OutputCell{
					{
						Times:  2,
						Values: Values{2301},
					},
					{
						Times:  0,
						Values: Values{2302},
					},
					{
						Times:  1,
						Values: Values{2303},
					},
				},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2301)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2301)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2302)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2303)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2303)
		})
		convey.Convey("OutputCell", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: OutputCell{
					Times:  2,
					Values: Values{2401},
				},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2401)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2401)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2401)
		})
		convey.Convey("*OutputCell", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: &OutputCell{
					Times:  2,
					Values: Values{2501},
				},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2501)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2501)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2501)
		})
		convey.Convey("[]Values", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: []Values{
					{2601},
					{2602},
					{2603},
				},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2601)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2602)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2603)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 2603)
		})
		convey.Convey("Values", func() {
			patches.Apply(MethodHook{
				Target: (*testMethodForHook)(nil),
				Method: "Test",
				Double: Values{1701},
			})
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 1701)
			convey.So(obj.Test(1, ""), convey.ShouldEqual, 1701)
		})
	})
}

func testFuncForInput(i *int, i64 *int64, s *string) {
	*i = -1000
	*i64 = -1000
	*s = "-1000"
}

func testFuncForInput2(datas ...any) {
}

func TestInput(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := NewPatches()
		defer patches.Reset()

		var i int
		var i64 int64
		var s string

		convey.Convey("auto find ...any", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput2,
				Double: ResultCell{
					Inputs: Values{int64(1234), "abc", 123},
				},
			})

			testFuncForInput2(&i, &i64, &s)
			convey.So(i, convey.ShouldEqual, 123)
			convey.So(i64, convey.ShouldEqual, int64(1234))
			convey.So(s, convey.ShouldEqual, "abc")
		})

		convey.Convey("auto find", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput,
				Double: ResultCell{
					Inputs: Values{int64(1234), "abc", 123},
				},
			})

			testFuncForInput(&i, &i64, &s)
			convey.So(i, convey.ShouldEqual, 123)
			convey.So(i64, convey.ShouldEqual, int64(1234))
			convey.So(s, convey.ShouldEqual, "abc")
		})

		convey.Convey("auto find (func)", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput,
				Double: ResultCell{
					Inputs: Values{func(v *string) { *v = "1234" }},
				},
			})

			testFuncForInput(&i, &i64, &s)
			convey.So(i, convey.ShouldEqual, 0)
			convey.So(i64, convey.ShouldEqual, int64(0))
			convey.So(s, convey.ShouldEqual, "1234")
		})

		convey.Convey("cannot find", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput,
				Double: ResultCell{
					Inputs: Values{uint(123)},
				},
			})
			convey.So(func() { testFuncForInput(&i, &i64, &s) }, convey.ShouldPanic)
		})

		convey.Convey("bind (int) succ", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput,
				Double: ResultCell{
					Inputs: Values{BindInput(0, 123)},
				},
			})

			testFuncForInput(&i, &i64, &s)
			convey.So(i, convey.ShouldEqual, 123)
			convey.So(i64, convey.ShouldEqual, int64(0))
			convey.So(s, convey.ShouldEqual, "")
		})
		convey.Convey("bind (int) panic", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput,
				Double: ResultCell{
					Inputs: Values{BindInput(1, 123)},
				},
			})

			convey.So(func() { testFuncForInput(&i, &i64, &s) }, convey.ShouldPanic)
		})

		convey.Convey("bind (func) succ", func() {
			patches.Apply(FuncHook{
				Target: testFuncForInput,
				Double: ResultCell{
					Inputs: Values{BindInput(2, func(s *string) { *s = "1234" })},
				},
			})

			testFuncForInput(&i, &i64, &s)
			convey.So(i, convey.ShouldEqual, 0)
			convey.So(i64, convey.ShouldEqual, int64(0))
			convey.So(s, convey.ShouldEqual, "1234")
		})
	})
}
