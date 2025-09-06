/*
 * Copyright (C) distroy
 */

package ldsync

type (
	doneWaitWaitGroup = WaitGroup
	doneWaitDone      = Done
)

type DoneWait struct {
	doneWaitWaitGroup
	doneWaitDone
}
