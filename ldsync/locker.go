/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"
)

type TryLocker interface {
	sync.Locker

	TryLock() bool
}
