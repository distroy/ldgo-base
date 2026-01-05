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

	func() {
		callback.Mu.Lock()
		defer callback.Mu.Unlock()
		callback.Callbacks = append(callback.Callbacks, cb)
	}()

	if !c.started.Done() {
		return
	}

	changes := func() map[string]*ChangeData {
		callback.Mu.Lock()
		defer callback.Mu.Unlock()

		cache := callback.Cache
		if len(cache) == 0 {
			return nil
		}

		res := make(map[string]*ChangeData, len(cache))
		for k, v := range cache {
			res[k] = &ChangeData{
				ChangeType: ChangeTypeAdd,
				NewValue:   v,
			}
		}
		return res
	}()

	ctx := c.getContext()
	cb(ctx, &NamespaceChangeEvent{
		Namespace: ns,
		Changes:   changes,
	})
}

func (c *Client) RegisterKey(ns, key string, cb func(context.Context, *KeyChangeEvent)) {
	callback := c.getKeyCallbacks(ns, key)

	func() {
		callback.Mu.Lock()
		defer callback.Mu.Unlock()
		callback.Callbacks = append(callback.Callbacks, cb)
	}()

	if !c.started.Done() {
		return
	}

	val, err := c.client.GetKey(ns, key)
	if err != nil {
		return
	}

	ctx := c.getContext()
	cb(ctx, &KeyChangeEvent{
		Namespace: ns,
		Key:       key,
		Change: &ChangeData{
			ChangeType: ChangeTypeAdd,
			NewValue:   val,
		},
	})
}
