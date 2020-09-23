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
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	redisMock, err := miniredis.Run()
	require.NoError(t, err)
	defer redisMock.Close()
	redisMock.RequireAuth("bk-cmdb")
	addrs := strings.Split(redisMock.Addr(), ":")

	prefix := "test_redis"
	configMap := map[string]string{
		prefix + ".host":       addrs[0],
		prefix + ".port":       addrs[1],
		prefix + ".pwd":        "bk-cmdb",
		prefix + ".database":   "0",
		prefix + ".mastername": "",
	}

	config, err := ParseConfig(prefix, configMap)
	require.NoError(t, err)

	err = InitClient(prefix, config)
	require.NoError(t, err)

	ctx := context.Background()
	testRedisKey := "test_redis_default_redis_key"
	// test default redis
	cacheErr := Client().Set(ctx, testRedisKey, "aa", time.Minute*10).Err()
	require.NoError(t, cacheErr)
	val, cacheErr := Client().Get(ctx, testRedisKey).Result()
	require.NoError(t, cacheErr)
	require.Equal(t, "aa", val)
	cacheErr = Client().Del(ctx, testRedisKey).Err()
	require.NoError(t, cacheErr)

	// test prefix
	cacheErr = ClientInstance(prefix).Set(ctx, testRedisKey, "aa", time.Minute*10).Err()
	require.NoError(t, cacheErr)
	val, cacheErr = ClientInstance(prefix).Get(ctx, testRedisKey).Result()
	require.NoError(t, cacheErr)
	require.Equal(t, "aa", val)
	cacheErr = ClientInstance(prefix).Del(ctx, testRedisKey).Err()
	require.NoError(t, cacheErr)

}
