/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"os"

	"github.com/distroy/ldgo-base/ldatomic"
)

var (
	console   = New(NewHandler(os.Stderr, nil))
	discard   = newDiscard()
	defLogger = ldatomic.NewAny(console)
)

func SetDefault(l *Logger) { defLogger.Store(l) }

func Default() *Logger { return defLogger.Load() }
func Console() *Logger { return console }
func Discard() *Logger { return discard }
