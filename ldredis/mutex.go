/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"context"
	"time"

	"github.com/distroy/ldgo-base/ldatomic"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldlog"
	"github.com/distroy/ldgo-base/ldrand"
	"github.com/distroy/ldgo-base/ldslice"
)

type MutexEvent int

const (
	MutexEventDeleted MutexEvent = iota + 1
)

var (
	closedMutextEventsChan chan MutexEvent
)

const (
	mutexMinHeartbeatInterval = 1 * time.Second
	mutexMinHeartbeatTimeout  = 10 * time.Second
	mutexMinLockForceInterval = 1 * time.Millisecond
)

type mutexContext struct {
	ctx    context.Context // control goroutine to exit
	cancel context.CancelFunc

	key           string
	token         string
	lastHeartbeat time.Time
	lockTime      ldatomic.Int64 // if equal 0, has not locked
	events        chan MutexEvent
}

func NewMutex[Redis MutexRedis[BoolCmd, IntCmd, StringCmd], BoolCmd BoolCmder, IntCmd IntCmder,
	StringCmd StringCmder,
](redis Redis) *Mutex[Redis, BoolCmd, IntCmd, StringCmd] {
	m := &Mutex[Redis, BoolCmd, IntCmd, StringCmd]{
		redis:             redis,
		heartbeatInterval: 5 * time.Second,
		heartbeatTimeout:  2 * time.Minute,
		lockForceInterval: 1 * time.Second,
	}
	return m
}

type Mutex[Redis MutexRedis[BoolCmd, IntCmd, StringCmd], BoolCmd BoolCmder, IntCmd IntCmder,
	StringCmd StringCmder,
] struct {
	redis             Redis
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration
	unlockDelay       time.Duration
	lockForceInterval time.Duration
	lockForceTimeout  time.Duration

	// if equal nil, has not locked
	// but when mutex has been cloned. maybe the lockTime is equal 0, the ctx is not equal nil
	ctx ldatomic.Ptr[mutexContext]
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) clone() *Mutex[Redis, BoolCmd, IntCmd, StringCmd] {
	c := *m
	return &c
}

// WithLockForce returns the redis mutex with lock force
// but if the context is cancelled, the lock force is not available
//
// WithLockForce should be called like these:
//
//	WithLockForce(false)
//	WithLockForce(true, {interval})
//	WithLockForce(true, {interval}, {timeout})
func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) WithLockForce(force bool, intervalAndTimeout ...time.Duration) *Mutex[Redis, BoolCmd, IntCmd, StringCmd] {
	m = m.clone()

	if !force {
		m.lockForceInterval = 0
		m.lockForceTimeout = 0
		return m
	}

	m.lockForceInterval = ldslice.Get(intervalAndTimeout, 0)
	m.lockForceInterval = max(m.lockForceInterval, mutexMinLockForceInterval)

	m.lockForceTimeout = ldslice.Get(intervalAndTimeout, 1)

	return m
}

// WithUnlockDelay returns the redis mutex with unlock delay
func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) WithUnlockDelay(delay ...time.Duration) *Mutex[Redis, BoolCmd, IntCmd, StringCmd] {
	m = m.clone()
	m.unlockDelay = 0
	if len(delay) > 0 && delay[0] > 0 {
		m.unlockDelay = delay[0]
	}
	return m
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) WithInterval(d time.Duration) *Mutex[Redis, BoolCmd, IntCmd, StringCmd] {
	m = m.clone()

	m.heartbeatInterval = max(d, mutexMinHeartbeatInterval)
	if timeout := m.getMinTimeout(d); m.heartbeatTimeout < timeout {
		m.heartbeatTimeout = timeout
	}
	return m
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) WithTimeout(d time.Duration) *Mutex[Redis, BoolCmd, IntCmd, StringCmd] {
	m = m.clone()
	m.heartbeatTimeout = m.getMinTimeout(d)
	return m
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) getMinTimeout(d time.Duration) time.Duration {
	d = max(d, mutexMinHeartbeatTimeout)
	if t := m.heartbeatInterval * 3; d < t {
		d = t
	}

	return d
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) mustGetContext() *mutexContext {
	mc := m.ctx.Load()
	if mc == nil {
		mc = &mutexContext{
			events: closedMutextEventsChan,
		}
	}
	return mc
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) getExpiration() time.Duration {
	return m.heartbeatInterval + m.heartbeatTimeout
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) Key() string { return m.mustGetContext().key }
func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) Events() <-chan MutexEvent {
	return m.mustGetContext().events
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) Lock(ctx context.Context, key string) error {
	mc := m.ctx.Load()
	if mc != nil && mc.lockTime.Load() != 0 {
		ldctx.LogE(ctx, "[ldredis] mutex has been locked", ldlog.String("key", key), ldlog.String("old", mc.key),
			getCallerField())
		return lderr.ErrCacheMutexLocked
	}
	if mc != nil {
		m.ctx.CompareAndSwap(mc, nil)
	}

	token := ldrand.String(16)

	ctx, cancel := ldctx.WithCancel(ctx)
	ctx = ldctx.WithLogger(ctx, nil, ldlog.String("key", key), ldlog.String("token", token))

	if err := m.internalLockOrLockForce(ctx, key, token); err != nil {
		return err
	}

	now := time.Now().UnixNano()
	mc = &mutexContext{
		ctx:      ctx,
		cancel:   cancel,
		key:      key,
		token:    token,
		lockTime: ldatomic.Int64(now),
		events:   make(chan MutexEvent, 1),
	}
	if ok := m.ctx.CompareAndSwap(nil, mc); !ok {
		// cli := m.redis
		// cli.Del(key)
		return lderr.ErrCacheMutexLocked
	}

	go m.goroutine(mc, now)
	ldctx.LogD(ctx, "[ldredis] mutex lock succ", getCallerField())
	return nil
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) internalLockOrLockForce(ctx context.Context, key, token string) error {
	err := m.internalLock(ctx, key, token)
	if err == nil {
		return nil
	}

	if m.lockForceInterval <= 0 {
		return err
	}

	timeout := m.lockForceTimeout
	if timeout > 0 {
		ctx, _ = ldctx.WithTimeout(ctx, timeout)
	}

	ticker := time.NewTicker(m.lockForceInterval)
	defer ticker.Stop()
	for err != nil {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return err

		case <-ticker.C:
			err = m.internalLock(ctx, key, token)
		}
	}

	return nil
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) internalLock(ctx context.Context, key, token string) error {
	cli := m.redis

	cmd := cli.SetNX(ctx, key, token, m.getExpiration())
	if err := cmd.Err(); err != nil {
		ldctx.LogE(ctx, "[ldredis] mutex setnx fail", ldlog.Error(err), getCallerField())
		return lderr.Wrap(err, lderr.ErrCacheRead)
	}

	if ok := cmd.Val(); !ok {
		ldctx.LogW(ctx, "[ldredis] mutex has been locked by another", getCallerField())
		return lderr.ErrCacheMutexLocked
	}

	return nil
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) Unlock(ctx context.Context) error {
	d := m.unlockDelay
	if d <= 0 {
		return m.unlock(ctx)
	}

	go func() {
		time.Sleep(d)
		m.unlock(ctx)
	}()

	return nil
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) unlock(ctx context.Context) error {
	mc := m.ctx.Load()
	if mc == nil {
		ldctx.LogW(ctx, "[ldredis] mutex has not been locked", getCallerField())
		return nil
	}

	lockTime := mc.lockTime.Load()
	if lockTime == 0 {
		ldctx.LogW(ctx, "[ldredis] mutex has not been locked", getCallerField())
		return nil
	}

	// ctx = mc.ctx
	cli := m.redis
	key := mc.key
	val := mc.token

	ctx = ldctx.WithLogger(ctx, nil, ldlog.String("key", key), ldlog.String("token", val))

	if ok := mc.lockTime.CompareAndSwap(lockTime, 0); !ok {
		ldctx.LogW(ctx, "[ldredis] mutex has been unlocked by another goroutine", getCallerField())
		return nil
	}
	m.ctx.CompareAndSwap(mc, nil)

	ldctx.LogD(ctx, "[ldredis] mutex will be unlocked", getCallerField())

	// ldctx.TryCancel(mc.ctx)
	mc.cancel()
	if err := m.checkToken(ctx, mc); err != nil {
		return err
	}

	cmd := cli.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		ldctx.LogW(ctx, "[ldredis] mutex del fail", ldlog.Error(err), getCallerField())
		return lderr.Wrap(err, lderr.ErrCacheWrite)
	}

	return nil
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) goroutine(mc *mutexContext, lockTime int64) {
	ctx := mc.ctx
	ldctx.LogD(ctx, "[ldredis] mutex goroutine start")
	ticker := time.NewTicker(m.heartbeatInterval)

	defer func() {
		ldctx.LogD(ctx, "[ldredis] mutex goroutine stop")

		// ldctx.TryCancel(ctx)
		mc.cancel()
		mc.lockTime.CompareAndSwap(lockTime, 0)

		close(mc.events)
		m.ctx.CompareAndSwap(mc, nil)

		ticker.Stop()
	}()

	mc.lastHeartbeat = time.Now()
	for running := true; running; {
		select {
		case now := <-ticker.C:
			running = m.heartbeat(ctx, mc, now)

		case <-ctx.Done():
			c := ldctx.WithLogger(ldctx.Default(), ldctx.GetLogger(ctx))
			m.unlock(c)
			return
		}
	}
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) heartbeat(ctx context.Context, mc *mutexContext, now time.Time) bool {
	switch err := m.checkToken(mc.ctx, mc); err {
	case nil:
		mc.lastHeartbeat = now

	case lderr.ErrCacheMutexNotExists, lderr.ErrCacheMutexNotMatch:
		m.doHeartbeatError(ctx, mc)
		return false

	default:
		return m.checkHeartbeatTime(ctx, mc)
	}

	return true
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) checkToken(ctx context.Context, mc *mutexContext) error {
	cli := m.redis
	key := mc.key
	val := mc.token
	{
		cmd := cli.Expire(ctx, key, m.getExpiration())
		if err := cmd.Err(); err != nil {
			ldctx.LogE(ctx, "[ldredis] mutex expire fail", ldlog.Error(err))
			return lderr.Wrap(err, lderr.ErrCacheWrite)
		}

		if ok := cmd.Val(); !ok {
			ldctx.LogE(ctx, "[ldredis] mutex is not exists")
			return lderr.ErrCacheMutexNotExists
		}
	}

	{
		cmd := cli.Get(ctx, key)
		if err := cmd.Err(); err != nil {
			ldctx.LogE(ctx, "[ldredis] mutex get fail", ldlog.Error(err))
			return lderr.Wrap(err, lderr.ErrCacheRead)
		}

		if val != cmd.Val() {
			ldctx.LogE(ctx, "[ldredis] mutex token is not match", ldlog.String("old", val),
				ldlog.String("new", cmd.Val()))
			return lderr.ErrCacheMutexNotMatch
		}
	}

	return nil
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) checkHeartbeatTime(ctx context.Context, mc *mutexContext) bool {
	if time.Since(mc.lastHeartbeat) < m.heartbeatTimeout {
		return true
	}

	ldctx.LogW(ctx, "[ldredis] mutex heartbeat timeout")
	m.doHeartbeatError(ctx, mc)
	return false
}

func (m *Mutex[Redis, BoolCmd, IntCmd, StringCmd]) doHeartbeatError(ctx context.Context, mc *mutexContext) {
	select {
	case <-ctx.Done():
	case mc.events <- MutexEventDeleted:
	}
}

func init() {
	closedMutextEventsChan = make(chan MutexEvent)
	close(closedMutextEventsChan)
}
