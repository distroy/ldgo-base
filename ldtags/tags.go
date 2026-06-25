/*
 * Copyright (C) distroy
 */

package ldtags

import (
	"strings"
)

type Tags map[string][]string

func New(size ...int) Tags {
	if len(size) > 0 && size[0] > 0 {
		return make(Tags, size[0])
	}
	return make(Tags)
}

func (m Tags) Add(key, value string) {
	key = strings.ToLower(key)
	m[key] = append(m[key], value)
}

func (m Tags) Set(key, value string) {
	key = strings.ToLower(key)
	m[key] = []string{value}
}

func (m Tags) Has(key string) bool {
	key = strings.ToLower(key)

	_, ok := m[key]
	return ok
}

func (m Tags) Values(key string) []string {
	key = strings.ToLower(key)
	return m[key]
}

func (m Tags) Get(key string, def ...string) string {
	key = strings.ToLower(key)

	v := m[key]
	if len(v) != 0 {
		return v[0]
	}

	if len(def) > 0 {
		return def[0]
	}

	return ""
}

// Parse returns ParseWithSeq(tag, ":", ";")
func Parse(tag string) Tags { return ParseWithSeq(tag, ":", ";") }

// ParseWithSeq(`name:size; default:1; meta:n`, ":", ";") return the follow map:
//
//	name: size
//	default: 1
//	meta: n
func ParseWithSeq(tag string, kvSeq, itemSeq string) Tags {
	tagList := strings.Split(tag, itemSeq)
	m := New(len(tagList))
	for _, v := range tagList {
		if len(v) == 0 {
			continue
		}

		l := strings.SplitN(v, kvSeq, 2)
		k := strings.TrimSpace(l[0])
		if k == "" {
			continue
		}

		v := ""
		if len(l) >= 2 {
			v = l[1]
		}

		m.Add(k, v)
	}
	return m
}
