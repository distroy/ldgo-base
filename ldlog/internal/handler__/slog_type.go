/*
 * Copyright (C) distroy
 */

package handler__

import (
	"log/slog"

	"github.com/distroy/ldgo-base/ldlog/internal/logref__"
	"github.com/distroy/ldgo-base/ldlog/internal/slogtype__"
)

func asType[T any](v any, def ...T) T { return logref__.AsType(v, def...) }

type (
	Level  = slogtype__.Level
	Kind   = slogtype__.Kind
	Attr   = slogtype__.Attr
	Value  = slogtype__.Value
	Source = slogtype__.Source
	Record = slogtype__.Record
)

// *** level begin ****

const (
	LevelTrace Level = slogtype__.LevelTrace
	LevelDebug Level = slogtype__.LevelDebug
	LevelInfo  Level = slogtype__.LevelInfo
	LevelWarn  Level = slogtype__.LevelWarn
	LevelError Level = slogtype__.LevelError
	LevelPanic Level = slogtype__.LevelPanic
)

// *** level end ****

// *** kind begin ****

const (
	KindAny       = slogtype__.KindAny
	KindBool      = slogtype__.KindBool
	KindDuration  = slogtype__.KindDuration
	KindFloat64   = slogtype__.KindFloat64
	KindInt64     = slogtype__.KindInt64
	KindString    = slogtype__.KindString
	KindTime      = slogtype__.KindTime
	KindUint64    = slogtype__.KindUint64
	KindGroup     = slogtype__.KindGroup
	KindLogValuer = slogtype__.KindLogValuer
)

// *** kind end ****

// *** source begin ****

func GetSourcePtr(v *slog.Source) *Source { return slogtype__.GetSourcePtr(v) }
func GetSource(v slog.Source) Source      { return slogtype__.GetSource(v) }

// *** source end ****

// *** value begin ****

func GetValuePtr(v *slog.Value) *Value { return slogtype__.GetValuePtr(v) }
func GetValue(v slog.Value) Value      { return slogtype__.GetValue(v) }

func CountEmptyGroups(as []Attr) int { return slogtype__.CountEmptyGroups(as) }

// *** value end ****

// *** attr begin ****

func GetAttrPtr(v *slog.Attr) *Attr { return slogtype__.GetAttrPtr(v) }
func GetAttr(v slog.Attr) Attr      { return slogtype__.GetAttr(v) }
func GetAttrs(v []slog.Attr) []Attr { return slogtype__.GetAttrs(v) }

func GetSAttrs(v []Attr) []slog.Attr { return slogtype__.GetSAttrs(v) }

// *** attr end ****

// *** record begin ****

func GetRecordPtr(v *slog.Record) *Record { return slogtype__.GetRecordPtr(v) }
func GetRecord(v slog.Record) Record      { return slogtype__.GetRecord(v) }

// *** record end ****
