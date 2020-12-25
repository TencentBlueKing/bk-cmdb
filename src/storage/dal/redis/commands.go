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
)

// Commands is the interface for redis commands
type Commands interface {
	Pipeline() Pipeliner

	BRPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceResult
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) StringResult
	Close() error
	Del(ctx context.Context, keys ...string) IntResult
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) Result
	Exists(ctx context.Context, keys ...string) IntResult
	Expire(ctx context.Context, key string, expiration time.Duration) BoolResult
	FlushDB(ctx context.Context) StatusResult
	Get(ctx context.Context, key string) StringResult
	HDel(ctx context.Context, key string, fields ...string) IntResult
	HGet(ctx context.Context, key, field string) StringResult
	HGetAll(ctx context.Context, key string) StringStringMapResult
	HIncrBy(ctx context.Context, key, field string, incr int64) IntResult
	HKeys(ctx context.Context, key string) StringSliceResult
	HMGet(ctx context.Context, key string, fields ...string) SliceResult
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ScanResult
	HSet(ctx context.Context, key string, values ...interface{}) IntResult
	Incr(ctx context.Context, key string) IntResult
	Keys(ctx context.Context, pattern string) StringSliceResult
	LLen(ctx context.Context, key string) IntResult
	LPush(ctx context.Context, key string, values ...interface{}) IntResult
	LRange(ctx context.Context, key string, start, stop int64) StringSliceResult
	LRem(ctx context.Context, key string, count int64, value interface{}) IntResult
	LTrim(ctx context.Context, key string, start, stop int64) StatusResult
	MGet(ctx context.Context, keys ...string) SliceResult
	MSet(ctx context.Context, values ...interface{}) StatusResult
	Ping(ctx context.Context) StatusResult
	Publish(ctx context.Context, channel string, message interface{}) IntResult
	Rename(ctx context.Context, key, newkey string) StatusResult
	RenameNX(ctx context.Context, key, newkey string) BoolResult
	RPop(ctx context.Context, key string) StringResult
	RPopLPush(ctx context.Context, source, destination string) StringResult
	RPush(ctx context.Context, key string, values ...interface{}) IntResult
	SAdd(ctx context.Context, key string, members ...interface{}) IntResult
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) StatusResult
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) BoolResult
	SMembers(ctx context.Context, key string) StringSliceResult
	SRem(ctx context.Context, key string, members ...interface{}) IntResult
	TTL(ctx context.Context, key string) DurationResult
}
