/*
 * Copyright (C) distroy
 */

// remote config
package ldrcfg

import (
	"context"
	"sync"

	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/ldlog"
	"github.com/distroy/ldgo-base/ldslice"
	"github.com/distroy/ldgo-base/ldsync"
)

type Adaptor interface {
	Type() string

	Start() error
	Stop()

	GetAllKeys(ns string) ([]string, error)
	GetKey(ns, key string) (string, error)
}

type ListenerAdaptor interface {
	Adaptor

	AddListener(ns string, cb func(*NamespaceChangeEvent))
}

func NewClient(cli Adaptor) *Client {
	_, ok := cli.(ListenerAdaptor)
	return &Client{
		client:     cli,
		autoListen: ok,
		listenChan: make(chan *listenData, 16),
	}
}

type Client struct {
	client     Adaptor
	logger     *ldlog.Logger
	callbacks  map[string]*clientNsCallback
	mu         sync.Mutex
	doneWait   ldsync.DoneWait
	autoListen bool
	listenChan chan *listenData
}

func (c *Client) GetAllKeys(ns string) []string {
	keys, err := c.client.GetAllKeys(ns)
	if err != nil {
		return nil
	}
	if c.autoListen {
		return keys
	}
	d := &listenData{
		Namespace: ns,
		Keys:      keys,
	}
	select {
	case c.listenChan <- d:
	default:
	}
	return keys
}
func (c *Client) GetKey(ns, key string, def ...string) string {
	val, err := c.client.GetKey(ns, key)
	if err != nil {
		return ldslice.Get(def, 0)
	}
	if c.autoListen {
		return val
	}
	d := &listenData{
		Namespace: ns,
		Key:       key,
		Value:     val,
	}
	select {
	case c.listenChan <- d:
	default:
	}
	return val
}

func (c *Client) Client() Adaptor                { return c.client }
func (c *Client) SetLogger(logger *ldlog.Logger) { c.logger = logger }

func (c *Client) Start(ctx context.Context) error {
	cli := c.client

	if err := cli.Start(); err != nil {
		return err
	}

	// trigger after start
	c.triggerAllCallback(ctx)

	if !c.autoListen {
		c.doneWait.Add(1)
		go c.listenGoroutine()
	}

	ldctx.LogI(ctx, "[ldrcfg] start succ", ldlog.String("type", cli.Type()))
	return nil
}

func (c *Client) Stop(ctx context.Context) {
	cli := c.client
	cli.Stop()

	if !c.autoListen {
		c.doneWait.Stop()
		c.doneWait.Wait()
	}

	ldctx.LogI(ctx, "[ldrcfg] stop succ", ldlog.String("type", cli.Type()))
}
