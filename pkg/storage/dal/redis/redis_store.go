/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package redis

import (
	"time"

	"github.com/FZambia/sentinel"
	"github.com/boj/redistore"
	"github.com/gin-contrib/sessions"
	"github.com/gomodule/redigo/redis"
	sess "github.com/gorilla/sessions"
)

// RedisStore interface
type RedisStore interface {
	sessions.Store
}

// NewRedisStore create redis store
// size: maximum number of idle connections.
// network: tcp or udp
// address: host:port
// password: redis-password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewRedisStore(size int, network, address, password string, keyPairs ...[]byte) (RedisStore, error) {
	store, err := redistore.NewRediStore(size, network, address, password, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &redisStore{store}, nil
}

type redisStore struct {
	*redistore.RediStore
}

// Options redisStore option
func (c *redisStore) Options(options sessions.Options) {
	c.RediStore.Options = &sess.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

// NewRedisStoreWithSentinel create redis sentinel store
// address: host:port array
// size: maximum number of idle connections.
// masterName: sentinel master name
// network: tcp or udp
// password: redis-password
// sentinelPwd: redis sentinel password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewRedisStoreWithSentinel(address []string, size int, masterName, network, password string, sentinelPwd string,
	keyPairs ...[]byte) (RedisStore, error) {

	sntnl := &sentinel.Sentinel{
		Addrs:      address,
		MasterName: masterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := time.Second
			return redis.Dial(network, addr,
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout),
				redis.DialWriteTimeout(timeout),
				redis.DialPassword(sentinelPwd),
			)
		},
	}

	pool := &redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			return dial(network, masterAddr, password)
		},
	}

	store, err := redistore.NewRediStoreWithPool(pool, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &redisSentinelStore{store}, nil
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

type redisSentinelStore struct {
	*redistore.RediStore
}

// Options redisSentinelStore option
func (c *redisSentinelStore) Options(options sessions.Options) {
	c.RediStore.Options = &sess.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
