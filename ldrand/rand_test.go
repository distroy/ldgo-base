/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/lditer"
)

func testFor[T any](seq lditer.Seq[T], do func(v T) bool) bool {
	for v := range seq {
		if !do(v) {
			return false
		}
	}
	return true
}

func testFor2[K, V any](seq lditer.Seq2[K, V], do func(k K, v V) bool) bool {
	for k, v := range seq {
		if !do(k, v) {
			return false
		}
	}
	return true
}

func TestRandString(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		a := String(16)
		b := String(16)
		c.So(a, convey.ShouldHaveLength, 16)
		c.So(b, convey.ShouldHaveLength, 16)
		c.So(a, convey.ShouldNotEqual, b)
	})
}
