/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"log/slog"

	"github.com/distroy/ldgo-base/ldlog/internal/handler__"
)

var (
	_ logHandler = (*handler__.Handler)(nil)
)

type logHandler interface {
	Handler

	Sync() error
	Close() error

	Level() slog.Level
	Sequence() string
}

func wrapHandler(h slog.Handler) logHandler {
	if h == nil {
		return nil
	}
	if hh, _ := h.(logHandler); hh != nil {
		return hh
	}
	return handlerWrapper{h}
}

type handlerWrapper struct {
	Handler
}

func (h handlerWrapper) Sync() error {
	switch hh := h.Handler.(type) {
	case interface{ Sync() error }:
		return hh.Sync()
	case interface{ Sync() }:
		hh.Sync()
	}
	return nil
}
func (h handlerWrapper) Close() error {
	switch hh := h.Handler.(type) {
	case interface{ Close() error }:
		return hh.Close()
	case interface{ Close() }:
		hh.Close()
	}
	return nil
}

func (h handlerWrapper) Level() slog.Level {
	switch hh := h.Handler.(type) {
	case interface{ Level() slog.Level }:
		return hh.Level()
	case interface{ Level() Level }:
		return hh.Level().Level()
	case interface{ Level() slog.Leveler }:
		return hh.Level().Level()
	}
	return LevelInfo.Level()
}
func (h handlerWrapper) Sequence() string {
	switch hh := h.Handler.(type) {
	case interface{ Sequence() string }:
		return hh.Sequence()
	}
	return ""
}
