package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"configcenter/src/common/ssl"
	localRedis "configcenter/src/storage/dal/redis"
)

func TestNewFromConfig(t *testing.T) {
	// 使用真实Redis服务器地址和认证信息
	redisAddr := "127.0.0.1:6379" // 默认地址，可以根据实际环境修改
	redisPassword := "cmdb"       // 用户提供的密码

	// 测试普通客户端配置
	t.Run("基本功能测试", func(t *testing.T) {
		cfg := localRedis.Config{
			Address:      redisAddr,
			Password:     redisPassword,
			Database:     "0",
			MaxOpenConns: 10,
			TLSConfig: ssl.TLSClientConfig{
				CAFile:             "/Users/yuyudeqiu/Desktop/canway/cmdb/redis_cert/ca-cert.pem",
				InsecureSkipVerify: true,
			},
		}

		// 连接Redis
		client, err := localRedis.NewFromConfig(cfg)
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
		assert.True(t, localRedis.IsNilErr(err), "键应该已被删除")
	})
}
