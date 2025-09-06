/*
 * Copyright (C) distroy
 */

package ldhook

import "reflect"

type VariableHook struct {
	Target any
	Double any
}

func (h VariableHook) hook(patches *patches) {
	patches.coreApplyVariable(reflect.ValueOf(h.Target), reflect.ValueOf(h.Double))
}
