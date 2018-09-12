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

package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	// goredis "gopkg.in/redis.v5"
	"configcenter/src/storage/dal/redis"
)

func TestRedis(t *testing.T) {
	conf := redis.ParseConfigFromKV("redis", map[string]string{
		"redis.host":     "127.0.0.1:6379",
		"redis.pwd":      "cc",
		"redis.database": "0",
	})
	cache, err := redis.NewFromConfig(*conf)
	require.NoError(t, err)

	err = cache.LPush("test_queue", "values1", "values2").Err()
	require.NoError(t, err)

	for {
		var value string
		err = cache.RPopLPush("test_queue", "test_queue2").Scan(&value)
		if redis.IsNil(err) {
			break
		}
		t.Logf("value : %s", value)
	}

	err = cache.Del("test_queue", "test_queue2").Err()
	require.NoError(t, err)
}
