/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"

	"github.com/distroy/ldgo-base/ldlog/internal/attr__"
)

func SetBriefStringLen(n int) { attr__.SetBriefStringLen(n) }
func SetBriefArrayLen(n int)  { attr__.SetBriefArrayLen(n) }
func SetBriefMapLen(n int)    { attr__.SetBriefMapLen(n) }

func BriefString(key string, val string) Attr         { return attr__.BriefString(key, val) }
func BriefByteString(key string, val []byte) Attr     { return attr__.BriefByteString(key, val) }
func BriefStringer(key string, val fmt.Stringer) Attr { return attr__.BriefStringer(key, val) }

func BriefStringp(key string, val *string) Attr               { return attr__.BriefStringp(key, val) }
func BriefStrings(key string, val []string) Attr              { return attr__.BriefStrings(key, val) }
func BriefByteStrings(key string, val [][]byte) Attr          { return attr__.BriefByteStrings(key, val) }
func BriefStringers[T fmt.Stringer](key string, val []T) Attr { return attr__.BriefStringers(key, val) }

func BriefReflect(key string, val any) Attr { return attr__.BriefReflect(key, val) }
