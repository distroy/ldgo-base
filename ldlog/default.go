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

func Tracex(msg string, attrs ...Attr) { Default().LogAttrs(lvlT, 1, msg, attrs...) }
func Debugx(msg string, attrs ...Attr) { Default().LogAttrs(lvlD, 1, msg, attrs...) }
func Infox(msg string, attrs ...Attr)  { Default().LogAttrs(lvlI, 1, msg, attrs...) }
func Warnx(msg string, attrs ...Attr)  { Default().LogAttrs(lvlW, 1, msg, attrs...) }
func Errorx(msg string, attrs ...Attr) { Default().LogAttrs(lvlE, 1, msg, attrs...) }
func Panicx(msg string, attrs ...Attr) { Default().LogAttrs(lvlP, 1, msg, attrs...) }

func Tracef(fmt string, args ...any) { Default().LogFmt(lvlT, 1, fmt, args...) }
func Debugf(fmt string, args ...any) { Default().LogFmt(lvlD, 1, fmt, args...) }
func Infof(fmt string, args ...any)  { Default().LogFmt(lvlI, 1, fmt, args...) }
func Warnf(fmt string, args ...any)  { Default().LogFmt(lvlW, 1, fmt, args...) }
func Errorf(fmt string, args ...any) { Default().LogFmt(lvlE, 1, fmt, args...) }
func Panicf(fmt string, args ...any) { Default().LogFmt(lvlP, 1, fmt, args...) }

func Traceln(args ...any) { Default().Logln(lvlT, 1, args...) }
func Debugln(args ...any) { Default().Logln(lvlD, 1, args...) }
func Infoln(args ...any)  { Default().Logln(lvlI, 1, args...) }
func Warnln(args ...any)  { Default().Logln(lvlW, 1, args...) }
func Errorln(args ...any) { Default().Logln(lvlE, 1, args...) }
func Panicln(args ...any) { Default().Logln(lvlP, 1, args...) }
