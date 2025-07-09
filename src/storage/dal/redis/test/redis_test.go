/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package redis_test

import (
	"context"
	"testing"
	"time"

	"configcenter/src/common/ssl"
	"configcenter/src/storage/dal/redis"

	"github.com/stretchr/testify/assert"
)

func TestNewFromConfig(t *testing.T) {
	redisAddr := ""     // 默认地址，可以根据实际环境修改
	redisPassword := "" // 用户提供的密码
	caFile := ""        // ca 证书路径

	// 测试普通客户端配置
	t.Run("基本功能测试", func(t *testing.T) {
		cfg := redis.Config{
			Address:      redisAddr,
			Password:     redisPassword,
			Database:     "0",
			MaxOpenConns: 10,
			TLSConfig: &ssl.TLSClientConfig{
				CAFile:             caFile,
				InsecureSkipVerify: true,
			},
		}

		// 连接Redis
		client, err := redis.NewFromConfig(cfg)
		if err != nil {
			t.Fatalf("连接Redis失败: %v", err)
		}
		assert.NotNil(t, client)

		// 测试基本操作
		ctx := context.Background()

		// 1. SET - 设置键值
		testKey := "test_key"
		testValue := "test_value"
		err = client.Set(ctx, testKey, testValue, time.Minute).Err()
		assert.NoError(t, err)

		// 2. GET - 获取键值
		val, err := client.Get(ctx, testKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, testValue, val)

		// 3. DEL - 删除键值
		n, err := client.Del(ctx, testKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), n)

		// 4. 验证删除结果
		_, err = client.Get(ctx, testKey).Result()
		assert.True(t, redis.IsNilErr(err), "键应该已被删除")
	})
}
