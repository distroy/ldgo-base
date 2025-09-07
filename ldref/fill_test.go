/*
 * Copyright (C) distroy
 */

package ldref

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

type testFillStructField struct {
	IntPtrPtr **int
}

type testFillStruct struct {
	unexported *int
	Slice      []int
	SlicePtr   *[]int
	Map        map[string]string
	MapPtr     *map[string]string
	StringPtr  *string
	FieldPtr   *testFillStructField
	// Chan       chan uint
	// ChanPtr    *chan uint
}

func TestFill(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("struct", func() {
			v := &testFillStruct{}
			Fill(v)

			r := &testFillStruct{
				unexported: new(int),
				Slice:      make([]int, 0),
				SlicePtr:   new([]int),
				Map:        make(map[string]string),
				MapPtr:     new(map[string]string),
				StringPtr:  new(string),
				FieldPtr:   new(testFillStructField),
				// Chan:       make(chan uint),
				// ChanPtr:    new(chan uint),
			}
			convey.So(v, convey.ShouldResemble, r)
		})
		convey.Convey("chan", func() {
			var v chan int
			Fill(&v)
			convey.So(v, convey.ShouldNotBeNil)
		})
		convey.Convey("*chan", func() {
			var v *chan int
			Fill(&v)
			convey.So(v, convey.ShouldNotBeNil)
			convey.So(*v, convey.ShouldBeNil)
		})
	})
}

func TestDeepFill(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("struct", func() {
			v := &testFillStruct{}
			DeepFill(v)

			r := &testFillStruct{
				unexported: new(int),
				Slice:      make([]int, 0),
				SlicePtr:   new([]int),
				Map:        make(map[string]string),
				MapPtr:     new(map[string]string),
				StringPtr:  new(string),
				FieldPtr:   new(testFillStructField),
				// Chan:       make(chan uint),
				// ChanPtr:    new(chan uint),
			}
			*r.SlicePtr = make([]int, 0)
			*r.MapPtr = make(map[string]string)
			// *r.ChanPtr = make(chan uint)

			r.FieldPtr.IntPtrPtr = new(*int)
			*r.FieldPtr.IntPtrPtr = new(int)

			convey.So(v, convey.ShouldResemble, r)
		})
		convey.Convey("chan", func() {
			var v chan int
			DeepFill(&v)
			convey.So(v, convey.ShouldNotBeNil)
		})
		convey.Convey("*chan", func() {
			var v *chan int
			DeepFill(&v)
			convey.So(v, convey.ShouldNotBeNil)
			convey.So(*v, convey.ShouldNotBeNil)
		})
	})
}
