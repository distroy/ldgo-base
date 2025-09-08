/*
 * Copyright (C) distroy
 */

package ldlog

import "github.com/distroy/ldgo-base/ldlog/internal/handler__"

func GetLevelKey() string  { return handler__.LevelKey }
func GetCallerKey() string { return handler__.CallerKey }

func SetSequenceKey(key string) { handler__.SequenceKey = key }
func GetSequenceKey() string    { return handler__.SequenceKey }

type Option func(l *core)

func SetLevel(lvl Level) Option   { return func(l *core) { l.withAttrs(Any(GetLevelKey(), lvl)) } }
func SetEnabler(e Enabler) Option { return func(l *core) { l.enabler = e } }
func SetSequence(s string) Option { return func(l *core) { l.withAttrs(String(GetSequenceKey(), s)) } }

func EnableCaller(e bool) Option     { return func(l *core) { l.withAttrs(Bool(GetCallerKey(), e)) } }
func AddCallerSkip(delta int) Option { return func(l *core) { l.callerSkip += delta } }
