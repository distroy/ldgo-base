/*
 * Copyright (C) distroy
 */

package ldtimer

type cronTask struct {
	cancel     func()
	name       string
	Minute     []int
	Hour       []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int
}
