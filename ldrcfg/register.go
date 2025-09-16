/*
 * Copyright (C) distroy
 */

package ldrcfg

import (
	"context"

	"github.com/distroy/ldgo-base/ldenum"
)

type ChangeType = ldenum.Enum[changeTypeToString]

const (
	ChangeTypeAdd ChangeType = iota
	ChangeTypeModify
	ChangeTypeDelete
)

type changeTypeToString struct{}

func (changeTypeToString) EnumToString(n int) string {
	switch ChangeType(n) {
	case ChangeTypeAdd:
		return "add"
	case ChangeTypeModify:
		return "modify"
	case ChangeTypeDelete:
		return "delete"
	}
	return ""
}

type ChangeData struct {
	ChangeType ChangeType
	OldValue   string
	NewValue   string
}

type NamespaceChangeEvent struct {
	Namespace string
	Changes   map[string]*ChangeData
}

type KeyChangeEvent struct {
	Namespace string
	Key       string
	Change    *ChangeData
}

func (c *Client) RegisterNamespace(ns string, cb func(context.Context, *NamespaceChangeEvent)) {
	callback := c.getNamespaceCallbacks(ns)

	callback.Mu.Lock()
	defer callback.Mu.Unlock()
	callback.Callbacks = append(callback.Callbacks, cb)
}

func (c *Client) RegisterKey(ns, key string, cb func(context.Context, *KeyChangeEvent)) {
	callback := c.getKeyCallbacks(ns, key)

	callback.Mu.Lock()
	defer callback.Mu.Unlock()
	callback.Callbacks = append(callback.Callbacks, cb)
}
