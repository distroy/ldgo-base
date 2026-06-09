/*
 * Copyright (C) distroy
 */

package ldlog

const (
	lvlT = LevelTrace
	lvlD = LevelDebug
	lvlI = LevelInfo
	lvlW = LevelWarn
	lvlE = LevelError
	lvlP = LevelPanic
)

func newLogger(log core) *Logger {
	return &Logger{
		core: log,
	}
}

type Logger struct {
	core
}

func (l *Logger) Wrapper() *Wrapper { return (*Wrapper)(l) }

func (l *Logger) clone() *Logger {
	cp := *l
	return &cp
}

func (l *Logger) With(attrs ...Attr) *Logger {
	if len(attrs) == 0 || l == nil || l.handler == nil {
		return l
	}
	l = l.clone()
	l.handler = wrapHandler(l.handler.WithAttrs(attrs))
	return l
}

func (l *Logger) WithOptions(opts ...Option) *Logger {
	if len(opts) == 0 || l == nil || l.handler == nil {
		return l
	}
	c := l.core
	for _, opt := range opts {
		opt(&c)
	}
	if c == l.core {
		return l
	}
	l = l.clone()
	l.core = c
	return l
}

func (l *Logger) Trace(msg string, attrs ...Attr) { l.LogAttrs(lvlT, 1, msg, attrs...) }
func (l *Logger) Debug(msg string, attrs ...Attr) { l.LogAttrs(lvlD, 1, msg, attrs...) }
func (l *Logger) Info(msg string, attrs ...Attr)  { l.LogAttrs(lvlI, 1, msg, attrs...) }
func (l *Logger) Warn(msg string, attrs ...Attr)  { l.LogAttrs(lvlW, 1, msg, attrs...) }
func (l *Logger) Error(msg string, attrs ...Attr) { l.LogAttrs(lvlE, 1, msg, attrs...) }
func (l *Logger) Panic(msg string, attrs ...Attr) { l.LogAttrs(lvlP, 1, msg, attrs...) }

func (l *Logger) Tracef(fmt string, args ...any) { l.LogFmt(lvlT, 1, fmt, args...) }
func (l *Logger) Debugf(fmt string, args ...any) { l.LogFmt(lvlD, 1, fmt, args...) }
func (l *Logger) Infof(fmt string, args ...any)  { l.LogFmt(lvlI, 1, fmt, args...) }
func (l *Logger) Warnf(fmt string, args ...any)  { l.LogFmt(lvlW, 1, fmt, args...) }
func (l *Logger) Errorf(fmt string, args ...any) { l.LogFmt(lvlE, 1, fmt, args...) }
func (l *Logger) Panicf(fmt string, args ...any) { l.LogFmt(lvlP, 1, fmt, args...) }

func (l *Logger) Traceln(args ...any) { l.Logln(lvlT, 1, args...) }
func (l *Logger) Debugln(args ...any) { l.Logln(lvlD, 1, args...) }
func (l *Logger) Infoln(args ...any)  { l.Logln(lvlI, 1, args...) }
func (l *Logger) Warnln(args ...any)  { l.Logln(lvlW, 1, args...) }
func (l *Logger) Errorln(args ...any) { l.Logln(lvlE, 1, args...) }
func (l *Logger) Panicln(args ...any) { l.Logln(lvlP, 1, args...) }
