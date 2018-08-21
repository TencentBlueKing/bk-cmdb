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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"configcenter/src/storage/dal/redis"
)

func TestRedis(t *testing.T) {
	conf := redis.NewConfigFromKV("redis", map[string]string{
		"address":  "127.0.0.1:6379",
		"pwd":      "cc",
		"database": "0",
	})
	cache, err := redis.NewFromConfig(*conf)
	require.NoError(t, err)

	err = cache.LPush("test_queue", "values1").Err()
	require.NoError(t, err)

	var value string
	for value != "nil" {
		err = cache.RPopLPush("test_queue", "test_queue2").Err()
		fmt.Fprintf(os.Stdout, "%s", value)
		require.NoError(t, err)
	}

	err = cache.Del("test_queue", "test_queue2").Err()
	require.NoError(t, err)
}
