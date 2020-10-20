/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
)

// Client is the interface for redis client
type Client interface {
	Subscribe(ctx context.Context, channels ...string) PubSub
	PSubscribe(ctx context.Context, channels ...string) PubSub

	Commands
}

type client struct {
	cli *redis.Client
}

// NewClient returns a client to the Redis Server specified by Options
func NewClient(opt *redis.Options) Client {
	return &client{
		cli: redis.NewClient(opt),
	}
}

// NewFailoverClient returns a Redis client that uses Redis Sentinel for automatic failover
func NewFailoverClient(failoverOpt *redis.FailoverOptions) Client {
	return &client{
		cli: redis.NewFailoverClient(failoverOpt),
	}
}

func (c *client) Subscribe(ctx context.Context, channels ...string) PubSub {
	return c.cli.Subscribe(channels...)
}

func (c *client) PSubscribe(ctx context.Context, channels ...string) PubSub {
	return c.cli.PSubscribe(channels...)
}

func (c *client) Pipeline() Pipeliner {
	return c.cli.Pipeline()
}

func (c *client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceResult {
	return c.cli.BRPop(timeout, keys...)
}

func (c *client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) StringResult {
	return c.cli.BRPopLPush(source, destination, timeout)
}

func (c *client) Close() error {
	return c.cli.Close()
}

func (c *client) Del(ctx context.Context, keys ...string) IntResult {
	return c.cli.Del(keys...)
}

func (c *client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) Result {
	return c.cli.Eval(script, keys, args...)
}

func (c *client) Exists(ctx context.Context, keys ...string) IntResult {
	return c.cli.Exists(keys...)
}

func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) BoolResult {
	return c.cli.Expire(key, expiration)
}

func (c *client) FlushDB(ctx context.Context) StatusResult {
	return c.cli.FlushDB()
}

func (c *client) Get(ctx context.Context, key string) StringResult {
	return c.cli.Get(key)
}

func (c *client) HDel(ctx context.Context, key string, fields ...string) IntResult {
	return c.cli.HDel(key, fields...)
}

func (c *client) HGet(ctx context.Context, key, field string) StringResult {
	return c.cli.HGet(key, field)
}

func (c *client) HGetAll(ctx context.Context, key string) StringStringMapResult {
	return c.cli.HGetAll(key)
}

func (c *client) HIncrBy(ctx context.Context, key, field string, incr int64) IntResult {
	return c.cli.HIncrBy(key, field, incr)
}

func (c *client) HKeys(ctx context.Context, key string) StringSliceResult {
	return c.cli.HKeys(key)
}

func (c *client) HMGet(ctx context.Context, key string, fields ...string) SliceResult {
	return c.cli.HMGet(key, fields...)
}

func (c *client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ScanResult {
	return c.cli.HScan(key, cursor, match, count)
}

func (c *client) HSet(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.HSet(key, values...)
}

func (c *client) Incr(ctx context.Context, key string) IntResult {
	return c.cli.Incr(key)
}

func (c *client) Keys(ctx context.Context, pattern string) StringSliceResult {
	return c.cli.Keys(pattern)
}

func (c *client) LLen(ctx context.Context, key string) IntResult {
	return c.cli.LLen(key)
}

func (c *client) LPush(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.LPush(key, values...)
}

func (c *client) LRange(ctx context.Context, key string, start, stop int64) StringSliceResult {
	return c.cli.LRange(key, start, stop)
}

func (c *client) LRem(ctx context.Context, key string, count int64, value interface{}) IntResult {
	return c.cli.LRem(key, count, value)
}

func (c *client) LTrim(ctx context.Context, key string, start, stop int64) StatusResult {
	return c.cli.LTrim(key, start, stop)
}

func (c *client) MGet(ctx context.Context, keys ...string) SliceResult {
	return c.cli.MGet(keys...)
}

func (c *client) MSet(ctx context.Context, values ...interface{}) StatusResult {
	return c.cli.MSet(values...)
}

func (c *client) Ping(ctx context.Context) StatusResult {
	return c.cli.Ping()
}

func (c *client) Publish(ctx context.Context, channel string, message interface{}) IntResult {
	return c.cli.Publish(channel, message)
}

func (c *client) Rename(ctx context.Context, key, newkey string) StatusResult {
	return c.cli.Rename(key, newkey)
}

func (c *client) RenameNX(ctx context.Context, key, newkey string) BoolResult {
	return c.cli.RenameNX(key, newkey)
}

func (c *client) RPop(ctx context.Context, key string) StringResult {
	return c.cli.RPop(key)
}

func (c *client) RPopLPush(ctx context.Context, source, destination string) StringResult {
	return c.cli.RPopLPush(source, destination)
}

func (c *client) RPush(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.RPush(key, values...)
}

func (c *client) SAdd(ctx context.Context, key string, members ...interface{}) IntResult {
	return c.cli.SAdd(key, members...)
}

func (c *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) StatusResult {
	return c.cli.Set(key, value, expiration)
}

func (c *client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) BoolResult {
	return c.cli.SetNX(key, value, expiration)
}

func (c *client) SMembers(ctx context.Context, key string) StringSliceResult {
	return c.cli.SMembers(key)
}

func (c *client) SRem(ctx context.Context, key string, members ...interface{}) IntResult {
	return c.cli.SRem(key, members...)
}

func (c *client) TTL(ctx context.Context, key string) DurationResult {
	return c.cli.TTL(key)
}
