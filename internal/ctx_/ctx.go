/*
 * Copyright (C) distroy
 */

package ctx_

import "github.com/distroy/ldgo-base/ldrand"

func NewSequence() string { return ldrand.String(16) }
