/*
 * Copyright (C) distroy
 */

package ldctx

import (
	"context"
	"sync"
)

type Map struct {
	m sync.Map
}

func (m *Map) Set(k, v any) bool {
	if m == nil {
		return false
	}
	m.m.Store(k, v)
	return true
}
func (m *Map) Get(k any) any {
	if m == nil {
		return nil
	}

	if v, ok := m.m.Load(k); ok {
		return v
	}

	for {
		p := m.getParent()
		if p == nil {
			return nil

		} else if v, ok := p.m.Load(k); ok {
			return v
		}

		m = p
	}
}
func (m *Map) Clear() {
	if m == nil {
		return
	}
	m.m.Range(func(key, value any) bool {
		m.m.Delete(key)
		return true
	})
}
func (m *Map) setParent(p *Map) { m.m.Store(ctxKeyMap, p) }
func (m *Map) getParent() *Map {
	i, _ := m.m.Load(ctxKeyMap)
	v, _ := i.(*Map)
	return v
}

func WithMap(c context.Context) context.Context {
	p := GetMap(c)
	m := &Map{}
	if p != nil {
		m.setParent(p)
	}
	return WithValue(c, ctxKeyMap, m)
}
func GetMap(c context.Context) *Map {
	m, _ := c.Value(ctxKeyMap).(*Map)
	return m
}
