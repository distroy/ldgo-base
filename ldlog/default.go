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
	defLogger = ldatomic.NewPtr(console)
)

func SetDefault(l *Logger) (old *Logger) { return defLogger.Swap(l) }
func SetDefaultWithClose(l *Logger) {
	old := SetDefault(l)
	old.Close()
}

func Default() *Logger { return defLogger.Load() }
func Console() *Logger { return console }
func Discard() *Logger { return discard }
