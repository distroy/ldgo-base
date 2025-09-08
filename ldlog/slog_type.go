/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"log/slog"

	"github.com/distroy/ldgo-base/ldlog/internal/slogtype__"
)

type (
	Level   = slogtype__.Level
	Attr    = slog.Attr
	Value   = slog.Value
	Record  = slog.Record
	Handler = slog.Handler
)

const (
	LevelTrace Level = slogtype__.LevelTrace
	LevelDebug Level = slogtype__.LevelDebug
	LevelInfo  Level = slogtype__.LevelInfo
	LevelWarn  Level = slogtype__.LevelWarn
	LevelError Level = slogtype__.LevelError
	LevelPanic Level = slogtype__.LevelPanic
)
