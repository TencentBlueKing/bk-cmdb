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

// Subscribe TODO
func (c *client) Subscribe(ctx context.Context, channels ...string) PubSub {
	return c.cli.Subscribe(channels...)
}

// PSubscribe TODO
func (c *client) PSubscribe(ctx context.Context, channels ...string) PubSub {
	return c.cli.PSubscribe(channels...)
}

// Pipeline TODO
func (c *client) Pipeline() Pipeliner {
	return c.cli.Pipeline()
}

// BRPop TODO
func (c *client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceResult {
	return c.cli.BRPop(timeout, keys...)
}

// BRPopLPush TODO
func (c *client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) StringResult {
	return c.cli.BRPopLPush(source, destination, timeout)
}

// Close TODO
func (c *client) Close() error {
	return c.cli.Close()
}

// Del TODO
func (c *client) Del(ctx context.Context, keys ...string) IntResult {
	return c.cli.Del(keys...)
}

// Eval TODO
func (c *client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) Result {
	return c.cli.Eval(script, keys, args...)
}

// Exists TODO
func (c *client) Exists(ctx context.Context, keys ...string) IntResult {
	return c.cli.Exists(keys...)
}

// Expire TODO
func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) BoolResult {
	return c.cli.Expire(key, expiration)
}

// FlushDB TODO
func (c *client) FlushDB(ctx context.Context) StatusResult {
	return c.cli.FlushDB()
}

// Get TODO
func (c *client) Get(ctx context.Context, key string) StringResult {
	return c.cli.Get(key)
}

// HDel TODO
func (c *client) HDel(ctx context.Context, key string, fields ...string) IntResult {
	return c.cli.HDel(key, fields...)
}

// HGet TODO
func (c *client) HGet(ctx context.Context, key, field string) StringResult {
	return c.cli.HGet(key, field)
}

// HGetAll TODO
func (c *client) HGetAll(ctx context.Context, key string) StringStringMapResult {
	return c.cli.HGetAll(key)
}

// HIncrBy TODO
func (c *client) HIncrBy(ctx context.Context, key, field string, incr int64) IntResult {
	return c.cli.HIncrBy(key, field, incr)
}

// HKeys TODO
func (c *client) HKeys(ctx context.Context, key string) StringSliceResult {
	return c.cli.HKeys(key)
}

// HMGet TODO
func (c *client) HMGet(ctx context.Context, key string, fields ...string) SliceResult {
	return c.cli.HMGet(key, fields...)
}

// HScan TODO
func (c *client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ScanResult {
	return c.cli.HScan(key, cursor, match, count)
}

// Scan TODO
func (c *client) Scan(ctx context.Context, cursor uint64, match string, count int64) ScanResult {
	return c.cli.Scan(cursor, match, count)
}

// HSet TODO
func (c *client) HSet(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.HSet(key, values...)
}

// Incr TODO
func (c *client) Incr(ctx context.Context, key string) IntResult {
	return c.cli.Incr(key)
}

// Keys TODO
func (c *client) Keys(ctx context.Context, pattern string) StringSliceResult {
	return c.cli.Keys(pattern)
}

// LLen TODO
func (c *client) LLen(ctx context.Context, key string) IntResult {
	return c.cli.LLen(key)
}

// LPush TODO
func (c *client) LPush(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.LPush(key, values...)
}

// LRange TODO
func (c *client) LRange(ctx context.Context, key string, start, stop int64) StringSliceResult {
	return c.cli.LRange(key, start, stop)
}

// LRem TODO
func (c *client) LRem(ctx context.Context, key string, count int64, value interface{}) IntResult {
	return c.cli.LRem(key, count, value)
}

// LTrim TODO
func (c *client) LTrim(ctx context.Context, key string, start, stop int64) StatusResult {
	return c.cli.LTrim(key, start, stop)
}

// MGet TODO
func (c *client) MGet(ctx context.Context, keys ...string) SliceResult {
	return c.cli.MGet(keys...)
}

// MSet TODO
func (c *client) MSet(ctx context.Context, values ...interface{}) StatusResult {
	return c.cli.MSet(values...)
}

// Ping TODO
func (c *client) Ping(ctx context.Context) StatusResult {
	return c.cli.Ping()
}

// Publish TODO
func (c *client) Publish(ctx context.Context, channel string, message interface{}) IntResult {
	return c.cli.Publish(channel, message)
}

// Rename TODO
func (c *client) Rename(ctx context.Context, key, newkey string) StatusResult {
	return c.cli.Rename(key, newkey)
}

// RenameNX TODO
func (c *client) RenameNX(ctx context.Context, key, newkey string) BoolResult {
	return c.cli.RenameNX(key, newkey)
}

// RPop TODO
func (c *client) RPop(ctx context.Context, key string) StringResult {
	return c.cli.RPop(key)
}

// RPopLPush TODO
func (c *client) RPopLPush(ctx context.Context, source, destination string) StringResult {
	return c.cli.RPopLPush(source, destination)
}

// RPush TODO
func (c *client) RPush(ctx context.Context, key string, values ...interface{}) IntResult {
	return c.cli.RPush(key, values...)
}

// SAdd TODO
func (c *client) SAdd(ctx context.Context, key string, members ...interface{}) IntResult {
	return c.cli.SAdd(key, members...)
}

// Set TODO
func (c *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) StatusResult {
	return c.cli.Set(key, value, expiration)
}

// SetNX TODO
func (c *client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) BoolResult {
	return c.cli.SetNX(key, value, expiration)
}

// TxPipeline TODO
func (c *client) TxPipeline(ctx context.Context) Pipeliner {
	return c.cli.TxPipeline()
}

// Discard TODO
func (c *client) Discard(ctx context.Context, pipe Pipeliner) error {
	return pipe.Discard()
}

// MSetNX TODO
func (c *client) MSetNX(ctx context.Context, values ...interface{}) BoolResult {
	return c.cli.MSetNX(values...)
}

// SMembers TODO
func (c *client) SMembers(ctx context.Context, key string) StringSliceResult {
	return c.cli.SMembers(key)
}

// SRem TODO
func (c *client) SRem(ctx context.Context, key string, members ...interface{}) IntResult {
	return c.cli.SRem(key, members...)
}

// TTL TODO
func (c *client) TTL(ctx context.Context, key string) DurationResult {
	return c.cli.TTL(key)
}

// BLPop TODO
func (c *client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceResult {
	return c.cli.BLPop(timeout, keys...)
}

// ZRemRangeByRank TODO
func (c *client) ZRemRangeByRank(key string, start, stop int64) IntResult {
	return c.cli.ZRemRangeByRank(key, start, stop)
}
