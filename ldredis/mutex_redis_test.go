/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"context"
	"sync"
	"time"

	"github.com/distroy/ldgo-base/ldconv"
)

func testMutexRedis() *memMutexRedis {
	return &memMutexRedis{
		data: make(map[string]*memMutexRedisData, 16),
	}
}

type memMutexRedisData struct {
	Deadline time.Time
	Value    any
}

type memMutexRedis struct {
	data map[string]*memMutexRedisData
	lock sync.Mutex
}

func (r *memMutexRedis) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data = make(map[string]*memMutexRedisData)
}

func (r *memMutexRedis) isDataExpire(d *memMutexRedisData, now time.Time) bool {
	if d == nil {
		return true
	}
	if d.Deadline.IsZero() {
		return false
	}
	return d.Deadline.Before(now)
}

func (r *memMutexRedis) WithContext(c context.Context) *memMutexRedis { return r }

func (r *memMutexRedis) SetNX(key string, value any, expiration time.Duration) BoolCmder {
	newCmd := func(val bool) BoolCmder {
		cmd := newTestBoolCmd("setnx", key, value, expiration)
		cmd.SetVal(val)
		return cmd
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	d := r.data[key]
	now := time.Now()
	if !r.isDataExpire(d, now) {
		return newCmd(false)
	}

	ddl := time.Time{}
	if expiration > 0 {
		ddl = now.Add(expiration)
	}
	r.data[key] = &memMutexRedisData{
		Deadline: ddl,
		Value:    value,
	}
	return newCmd(true)
}

func (r *memMutexRedis) Get(key string) StringCmder {
	newCmd := func(val string) StringCmder {
		cmd := newTestStringCmd("get", key)
		cmd.SetVal(val)
		return cmd
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	d := r.data[key]
	now := time.Now()
	if r.isDataExpire(d, now) {
		return newCmd("")
	}
	return newCmd(ldconv.AsString(d.Value))
}
func (r *memMutexRedis) Expire(key string, expiration time.Duration) BoolCmder {
	newCmd := func(val bool) BoolCmder {
		cmd := newTestBoolCmd("get", key)
		cmd.SetVal(val)
		return cmd
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	d := r.data[key]
	now := time.Now()
	if r.isDataExpire(d, now) {
		return newCmd(false)
	}
	d.Deadline = now.Add(expiration)
	return newCmd(true)
}
func (r *memMutexRedis) Del(keys ...string) IntCmder {
	newCmd := func(val int64) IntCmder {
		args := make([]any, 0, len(keys))
		args = append(args, "del")
		for _, key := range keys {
			args = append(args, key)
		}
		cmd := newTestIntCmd(args...)
		cmd.SetVal(val)
		return cmd
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	cnt := int64(0)
	for _, key := range keys {
		cnt += r.del(key)
	}
	return newCmd(cnt)
}

func (r *memMutexRedis) del(key string) int64 {
	d := r.data[key]
	now := time.Now()
	if r.isDataExpire(d, now) {
		return 0
	}
	delete(r.data, key)
	return 1
}
