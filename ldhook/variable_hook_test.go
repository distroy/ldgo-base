/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestVariableHook(t *testing.T) {
	var _TestVariableHook *int
	convey.Convey(t.Name(), t, func() {
		v0 := 100
		v1 := 200

		patches := NewPatches()
		defer patches.Reset()

		patches.Apply(VariableHook{
			Target: &_TestVariableHook,
			Double: &v0,
		})
		convey.So(_TestVariableHook, convey.ShouldNotBeNil)
		convey.So(*_TestVariableHook, convey.ShouldEqual, v0)

		patches.Apply(VariableHook{
			Target: &_TestVariableHook,
			Double: &v1,
		})
		convey.So(_TestVariableHook, convey.ShouldNotBeNil)
		convey.So(*_TestVariableHook, convey.ShouldEqual, v1)

		patches.Reset()
		convey.So(_TestVariableHook, convey.ShouldBeNil)
	})
}
