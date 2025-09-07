/*
 * Copyright (C) distroy
 */

package ldsort

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestSearch(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		arr := []int64{7, 11, 20, 25, 38, 50, 55, 76, 76, 84}
		convey.So(Search(arr, func(v int64) bool { return v >= 19 }), convey.ShouldEqual, 2)
		convey.So(Search(arr, func(v int64) bool { return v >= 20 }), convey.ShouldEqual, 2)
		convey.So(Search(arr, func(v int64) bool { return v >= 21 }), convey.ShouldEqual, 3)
	})
}

func TestIndex(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		arr := []int64{7, 11, 20, 25, 38, 50, 55, 76, 76, 84}
		convey.So(Index(arr, func(v int64) int { return int(v - 19) }), convey.ShouldEqual, -1)
		convey.So(Index(arr, func(v int64) int { return int(v - 20) }), convey.ShouldEqual, 2)
		convey.So(Index(arr, func(v int64) int { return int(v - 21) }), convey.ShouldEqual, -1)
	})
}
