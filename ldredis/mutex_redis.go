/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"context"
	"time"
)

type MutexRedis[BoolCmd BoolCmder, IntCmd IntCmder, StringCmd StringCmder] interface {
	SetNX(c context.Context, key string, value any, expiration time.Duration) BoolCmd
	Get(c context.Context, key string) StringCmd
	Expire(c context.Context, key string, expiration time.Duration) BoolCmd
	Del(c context.Context, keys ...string) IntCmd
}

type MutexRedisWithCtx[Redis any, BoolCmd BoolCmder, IntCmd IntCmder, StringCmd StringCmder] interface {
	SetNX(key string, value any, expiration time.Duration) BoolCmd
	Get(key string) StringCmd
	Expire(key string, expiration time.Duration) BoolCmd
	Del(keys ...string) IntCmd

	WithContext(c context.Context) Redis
}

func WrapMutexRedisWithCtx[Redis MutexRedisWithCtx[Redis, BoolCmd, IntCmd, StringCmd],
	BoolCmd BoolCmder, IntCmd IntCmder, StringCmd StringCmder,
](rdb Redis) MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd] {
	return MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd]{
		redis: rdb,
	}
}

type MutexRedisWithCtxWrapper[Redis MutexRedisWithCtx[Redis, BoolCmd, IntCmd, StringCmd],
	BoolCmd BoolCmder, IntCmd IntCmder, StringCmd StringCmder,
] struct {
	redis Redis
}

func (r MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd]) get(c context.Context) Redis {
	return r.redis.WithContext(c)
}
func (r MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd]) SetNX(c context.Context, key string, value any, expiration time.Duration) BoolCmd {
	return r.get(c).SetNX(key, value, expiration)
}

func (r MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd]) Get(c context.Context, key string) StringCmd {
	return r.get(c).Get(key)
}
func (r MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd]) Expire(c context.Context, key string, expiration time.Duration) BoolCmd {
	return r.get(c).Expire(key, expiration)
}
func (r MutexRedisWithCtxWrapper[Redis, BoolCmd, IntCmd, StringCmd]) Del(c context.Context, keys ...string) IntCmd {
	return r.get(c).Del(keys...)
}
