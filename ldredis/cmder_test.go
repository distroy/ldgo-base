/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/ldslice"
)

type testBaseCmd struct {
	args []any
	err  error
	val  any
}

func (c *testBaseCmd) Name() string     { return ldconv.AsString(ldslice.Get(c.args, 0)) }
func (c *testBaseCmd) FullName() string { return c.Name() }
func (c *testBaseCmd) Args() []any      { return c.args }
func (c *testBaseCmd) SetErr(err error) { c.err = err }
func (c *testBaseCmd) Err() error       { return c.err }

func (c *testBaseCmd) String() string {
	b := make([]byte, 0, 64)
	for i, arg := range c.args {
		if i > 0 {
			b = append(b, ' ')
		}
		b = append(b, ldconv.AsString(arg)...)
	}

	if err := c.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	} else if val := c.val; val != nil {
		b = append(b, ": "...)
		b = append(b, ldconv.AsString(val)...)
	}

	return ldconv.BytesToStrUnsafe(b)
}

type testValueCmd[V any] struct {
	testBaseCmd
}

func (c *testValueCmd[V]) Result() (V, error) { return c.Val(), c.Err() }
func (c *testValueCmd[V]) SetVal(val V)       { c.val = val }
func (c *testValueCmd[V]) Val() V {
	v, _ := c.val.(V)
	return v
}

func newTestBoolCmd(args ...any) *testBoolCmd {
	return &testBoolCmd{
		testValueCmd: testValueCmd[bool]{
			testBaseCmd: testBaseCmd{
				args: args,
				err:  nil,
				val:  nil,
			},
		},
	}
}

type testBoolCmd struct {
	testValueCmd[bool]
}

func newTestIntCmd(args ...any) *testIntCmd {
	return &testIntCmd{
		testValueCmd: testValueCmd[int64]{
			testBaseCmd: testBaseCmd{
				args: args,
				err:  nil,
				val:  nil,
			},
		},
	}
}

type testIntCmd struct {
	testValueCmd[int64]
}

func newTestStringCmd(args ...any) *testStringCmd {
	return &testStringCmd{
		testValueCmd: testValueCmd[string]{
			testBaseCmd: testBaseCmd{
				args: args,
				err:  nil,
				val:  nil,
			},
		},
	}
}

type testStringCmd struct {
	testValueCmd[string]
}
