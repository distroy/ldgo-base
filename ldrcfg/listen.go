/*
 * Copyright (C) distroy
 */

package ldrcfg

import (
	"context"
	"sync"
	"time"

	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/ldlog"
	"github.com/distroy/ldgo-base/ldrand"
)

type listenData struct {
	Namespace string
	Key       string
	Keys      []string
	Value     string
}

type clientNsCallback struct {
	Namespace string
	Callbacks []func(context.Context, *NamespaceChangeEvent)
	Keys      map[string]*clientKeyCallback
	Cache     map[string]string // only r/w in goroutine
	Mu        sync.Mutex
}

type clientKeyCallback struct {
	Namespace string
	Key       string
	Callbacks []func(context.Context, *KeyChangeEvent)

	Mu sync.Mutex
}

func (c *Client) newSequence() string {
	a, _ := c.client.(interface{ NewSequence() string })
	if a != nil {
		return a.NewSequence()
	}
	return ldrand.String(16)
}

func (c *Client) getContext() context.Context {
	ctx := c.getContextWithoutSequence()
	seq := c.newSequence()
	if seq != "" {
		ctx = ldctx.WithSequence(ctx, seq)
	}
	return ctx
}

func (c *Client) getContextWithoutSequence() context.Context {
	ctx := ldctx.Default()
	l := c.logger
	if l != nil {
		ctx = ldctx.WithLogger(ctx, l)
	}
	return ctx
}

func (c *Client) getNamespaceCallbacks(ns string) *clientNsCallback {
	c.mu.Lock()
	defer c.mu.Unlock()

	nsCallback := c.callbacks[ns]
	if nsCallback != nil {
		return nsCallback
	}

	cli, _ := c.client.(ListenerAdaptor)
	if cli != nil {
		cli.AddListener(ns, c.onChange)
	}

	nsCallback = &clientNsCallback{
		Namespace: ns,
	}

	if c.callbacks == nil {
		c.callbacks = make(map[string]*clientNsCallback)
	}

	c.callbacks[ns] = nsCallback
	return nsCallback
}

func (c *Client) getKeyCallbacks(ns, key string) *clientKeyCallback {
	nsCallback := c.getNamespaceCallbacks(ns)
	nsCallback.Mu.Lock()
	defer nsCallback.Mu.Unlock()

	keyCallback := nsCallback.Keys[key]
	if keyCallback == nil {
		keyCallback = &clientKeyCallback{
			Namespace: ns,
			Key:       key,
		}

		if nsCallback.Keys == nil {
			nsCallback.Keys = make(map[string]*clientKeyCallback)
		}

		nsCallback.Keys[key] = keyCallback
	}

	return keyCallback
}

func (c *Client) onChange(ev *NamespaceChangeEvent) {
	ctx := c.getContext()
	c.onChangeWitContext(ctx, ev)
}

func (c *Client) onChangeWitContext(ctx context.Context, ev *NamespaceChangeEvent) {
	if seq := ldctx.GetSequence(ctx); seq == "" {
		ctx = ldctx.WithSequence(ctx, c.newSequence())
	}

	nsCallback := c.callbacks[ev.Namespace]
	if nsCallback == nil {
		ldctx.LogW(ctx, "[config center] the namespace has not been registered",
			ldlog.String("ns", ev.Namespace))
		return
	}

	nsCbs := func() []func(context.Context, *NamespaceChangeEvent) {
		nsCallback.Mu.Lock()
		defer nsCallback.Mu.Unlock()

		return nsCallback.Callbacks
	}()
	for _, cb := range nsCbs {
		cb(ctx, ev)
	}

	for key, change := range ev.Changes {
		keyCallback := func() *clientKeyCallback {
			nsCallback.Mu.Lock()
			defer nsCallback.Mu.Unlock()

			return nsCallback.Keys[key]
		}()

		if keyCallback == nil {
			continue
		}

		keyCbs := func() []func(context.Context, *KeyChangeEvent) {
			keyCallback.Mu.Lock()
			defer keyCallback.Mu.Unlock()

			return keyCallback.Callbacks
		}()
		for _, cb := range keyCbs {
			cb(ctx, &KeyChangeEvent{
				Namespace: ev.Namespace,
				Key:       key,
				Change:    change,
			})
		}
	}
}

func (c *Client) listenGoroutine() {
	ticker := time.NewTicker(time.Second * 60)
	defer func() {
		ticker.Stop()

		c.doneWait.Done()
	}()

	for {
		select {
		case <-c.doneWait.Chan():
			return

		case d := <-c.listenChan:
			c.triggerByListenData(d)

		case <-ticker.C:
			c.triggerAllCallback(c.getContext())
		}
	}
}

func (c *Client) triggerByListenData(d *listenData) {
	ns := c.getNamespaceCallbacks(d.Namespace)
	if d.Key != "" {
		oldVal, ok := ns.Cache[d.Key]
		if oldVal == d.Value {
			return
		}

		changeType := ChangeTypeAdd
		if ok {
			changeType = ChangeTypeModify
		}

		c.onChange(&NamespaceChangeEvent{
			Namespace: d.Namespace,
			Changes: map[string]*ChangeData{
				d.Key: {
					ChangeType: changeType,
					OldValue:   oldVal,
					NewValue:   d.Value,
				},
			},
		})
	}
	if ns == nil {
		return
	}

	if len(d.Keys) != len(ns.Cache) {
		c.triggerNsCallback(c.getContext(), ns)
		return
	}

	for _, key := range d.Keys {
		_, ok := ns.Cache[key]
		if !ok {
			c.triggerNsCallback(c.getContext(), ns)
			return
		}
	}
}

func (c *Client) triggerAllCallback(ctx context.Context) {
	nsCaches := func() []*clientNsCallback {
		c.mu.Lock()
		defer c.mu.Unlock()

		res := make([]*clientNsCallback, 0, len(c.callbacks))
		for _, ns := range c.callbacks {
			res = append(res, ns)
		}
		return res
	}()

	for _, ns := range nsCaches {
		c.triggerNsCallback(ctx, ns)
	}
}

func (c *Client) triggerNsCallback(ctx context.Context, ns *clientNsCallback) {
	newMap, err := c.getNsMap(ns.Namespace)
	if err != nil {
		return
	}

	oldMap := ns.Cache
	if !c.isNsMapChanged(ns.Cache, newMap) {
		return
	}

	ev := &NamespaceChangeEvent{
		Namespace: ns.Namespace,
		Changes:   make(map[string]*ChangeData, len(newMap)),
	}

	for key, newVal := range newMap {
		oldVal, ok := oldMap[key]
		if !ok {
			// add
			ev.Changes[key] = &ChangeData{
				ChangeType: ChangeTypeAdd,
				OldValue:   "",
				NewValue:   newVal,
			}

		} else if newVal != oldVal {
			// modify
			ev.Changes[key] = &ChangeData{
				ChangeType: ChangeTypeModify,
				OldValue:   oldVal,
				NewValue:   newVal,
			}
		}
	}

	for key, oldVal := range oldMap {
		_, ok := newMap[key]
		if !ok {
			// delete
			ev.Changes[key] = &ChangeData{
				ChangeType: ChangeTypeDelete,
				OldValue:   oldVal,
			}
		}
	}

	ldctx.LogT(ctx, "[config center] namespace update", ldlog.Reflect("event", ev))
	ns.Cache = newMap
	c.onChangeWitContext(ctx, ev)
}

func (c *Client) isNsMapChanged(m0, m1 map[string]string) bool {
	if len(m0) != len(m1) {
		return true
	}

	for k, v0 := range m0 {
		v1, ok := m1[k]
		if !ok || v0 != v1 {
			return true
		}
	}

	return false
}

func (c *Client) getNsMap(ns string) (map[string]string, error) {
	cli := c.client
	keys, err := cli.GetAllKeys(ns)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(keys))
	for _, key := range keys {
		val, err := cli.GetKey(ns, key)
		if err != nil {
			return nil, err
		}
		res[key] = val
	}
	return res, nil
}
