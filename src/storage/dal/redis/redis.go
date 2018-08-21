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
	"strconv"
	"strings"

	redis "gopkg.in/redis.v5"
)

// Config define redis config
type Config struct {
	Address    string
	Password   string
	Database   string
	MasterName string
}

// NewConfigFromKV returns new config
func NewConfigFromKV(profix string, conifgmap map[string]string) *Config {
	prefix := "redis"
	return &Config{
		Address:    conifgmap[prefix+".address"],
		Password:   conifgmap[prefix+".pwd"],
		Database:   conifgmap[prefix+".database"],
		MasterName: conifgmap[prefix+".mastername"],
	}
}

// NewFromConfig returns new redis client from config
func NewFromConfig(cfg Config) (*redis.Client, error) {
	dbNum, err := strconv.Atoi(cfg.Database)
	if nil != err {
		return nil, err
	}
	var client *redis.Client
	if cfg.MasterName == "" {
		option := &redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       dbNum,
			PoolSize: 100,
		}
		client = redis.NewClient(option)
	} else {
		hosts := strings.Split(cfg.Address, ",")
		option := &redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			SentinelAddrs: hosts,
			Password:      cfg.Password,
			DB:            dbNum,
			PoolSize:      100,
		}
		client = redis.NewFailoverClient(option)
	}

	err = client.Ping().Err()
	if err != nil {
		return nil, err
	}

	return client, err
}
