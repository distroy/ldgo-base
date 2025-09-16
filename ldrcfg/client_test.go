/*
 * Copyright (C) distroy
 */

package ldrcfg

import (
	"log"
	"os"
	"testing"

	"github.com/distroy/ldgo-base/ldhook"
)

var (
	_ (Adaptor) = (*testAdaptor)(nil)
)

func testTriggerEvent(cli *Client, ns, key, value string) {
	cli.onChange(&NamespaceChangeEvent{
		Namespace: ns,
		Changes: map[string]*ChangeData{
			key: {
				OldValue: "",
				NewValue: value,
			},
		},
	})
}

type testAdaptor struct{}

func (testAdaptor) Type() string                           { return "test" }
func (testAdaptor) Start() error                           { return nil }
func (testAdaptor) Stop()                                  {}
func (testAdaptor) GetAllKeys(ns string) ([]string, error) { return nil, nil }
func (testAdaptor) GetKey(ns, key string) (string, error)  { return "", nil }
func (testAdaptor) AddListener(ns string, cb func(*NamespaceChangeEvent)) {
}

func TestMain(m *testing.M) {
	patches := ldhook.NewPatches()
	defer patches.Reset()

	patches.Applys([]ldhook.Hook{
		// ldhook.FuncHook{
		// 	Target: (*Client).getContext,
		// 	Double: ldhook.Values{ldctx.Discard()},
		// },
	})

	log.SetFlags(log.Flags() | log.Lshortfile)
	os.Exit(m.Run())
}
