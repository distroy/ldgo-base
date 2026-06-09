/*
 * Copyright (C) distroy
 */

package ldlog

import "context"

type Wrapper struct {
	core
}

func (l *Wrapper) Logger() *Logger { return (*Logger)(l) }

func (l *Wrapper) Tracef(fmt string, args ...any)   { l.LogFmt(lvlT, 1, fmt, args...) }
func (l *Wrapper) Trace(args ...any)                { l.Logln(lvlT, 1, args) }
func (l *Wrapper) Traceln(args ...any)              { l.Logln(lvlT, 1, args) }
func (l *Wrapper) Tracex(msg string, attrs ...Attr) { l.LogAttrs(lvlT, 1, msg, attrs...) }

func (l *Wrapper) Debugf(fmt string, args ...any)   { l.LogFmt(lvlD, 1, fmt, args...) }
func (l *Wrapper) Debug(args ...any)                { l.Logln(lvlD, 1, args) }
func (l *Wrapper) Debugln(args ...any)              { l.Logln(lvlD, 1, args) }
func (l *Wrapper) Debugx(msg string, attrs ...Attr) { l.LogAttrs(lvlD, 1, msg, attrs...) }

func (l *Wrapper) Infof(fmt string, args ...any)   { l.LogFmt(lvlI, 1, fmt, args...) }
func (l *Wrapper) Info(args ...any)                { l.Logln(lvlI, 1, args) }
func (l *Wrapper) Infoln(args ...any)              { l.Logln(lvlI, 1, args) }
func (l *Wrapper) Infox(msg string, attrs ...Attr) { l.LogAttrs(lvlI, 1, msg, attrs...) }

func (l *Wrapper) Printf(fmt string, args ...any)   { l.LogFmt(lvlI, 1, fmt, args...) }
func (l *Wrapper) Print(args ...any)                { l.Logln(lvlI, 1, args) }
func (l *Wrapper) Println(args ...any)              { l.Logln(lvlI, 1, args) }
func (l *Wrapper) Printx(msg string, attrs ...Attr) { l.LogAttrs(lvlI, 1, msg, attrs...) }

func (l *Wrapper) Logf(fmt string, args ...any)   { l.LogFmt(lvlI, 1, fmt, args...) }
func (l *Wrapper) Log(args ...any)                { l.Logln(lvlI, 1, args) }
func (l *Wrapper) Logln(args ...any)              { l.Logln(lvlI, 1, args) }
func (l *Wrapper) Logx(msg string, attrs ...Attr) { l.LogAttrs(lvlI, 1, msg, attrs...) }

func (l *Wrapper) Warnf(fmt string, args ...any)   { l.LogFmt(lvlW, 1, fmt, args...) }
func (l *Wrapper) Warn(args ...any)                { l.Logln(lvlW, 1, args) }
func (l *Wrapper) Warnln(args ...any)              { l.Logln(lvlW, 1, args) }
func (l *Wrapper) Warnx(msg string, attrs ...Attr) { l.LogAttrs(lvlW, 1, msg, attrs...) }

func (l *Wrapper) Warningf(fmt string, args ...any)   { l.LogFmt(lvlW, 1, fmt, args...) }
func (l *Wrapper) Warning(args ...any)                { l.Logln(lvlW, 1, args) }
func (l *Wrapper) Warningln(args ...any)              { l.Logln(lvlW, 1, args) }
func (l *Wrapper) Warningx(msg string, attrs ...Attr) { l.LogAttrs(lvlW, 1, msg, attrs...) }

func (l *Wrapper) Errorf(fmt string, args ...any)   { l.LogFmt(lvlE, 1, fmt, args...) }
func (l *Wrapper) Error(args ...any)                { l.Logln(lvlE, 1, args) }
func (l *Wrapper) Errorln(args ...any)              { l.Logln(lvlE, 1, args) }
func (l *Wrapper) Errorx(msg string, attrs ...Attr) { l.LogAttrs(lvlE, 1, msg, attrs...) }

func (l *Wrapper) Panicf(fmt string, args ...any)   { l.LogFmt(lvlP, 1, fmt, args...) }
func (l *Wrapper) Panic(args ...any)                { l.Logln(lvlP, 1, args) }
func (l *Wrapper) Panicln(args ...any)              { l.Logln(lvlP, 1, args) }
func (l *Wrapper) Panicx(msg string, attrs ...Attr) { l.LogAttrs(lvlP, 1, msg, attrs...) }

func (l *Wrapper) V(v int) bool {
	if v <= 0 {
		return !l.Enabled(context.Background(), LevelDebug)
	}
	return true
}
