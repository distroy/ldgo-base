/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"time"
)

type CronEvery struct {
	d time.Duration
}

func (c CronEvery) Match(t time.Time) bool                { return true }
func (c CronEvery) Next(from time.Time) (time.Time, bool) { return from.Add(c.d), true }
