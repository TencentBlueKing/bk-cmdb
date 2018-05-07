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
 
package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFromBytes(t *testing.T) {
	data := []byte(`[mongodb]
host=127.0.0.1
usr=user
pwd=pwd
database=cmdb
port=27107
maxOpenConns=3000
maxIDleConns=1000
[redis]
host=127.0.0.1
pwd=redisauth
database=0
port=6379
maxOpenConns=3000
maxIDleConns=1000
[errors]
res=conf/errors
`)
	expect := Config{
		"mongodb.host":         "127.0.0.1",
		"mongodb.usr":          "user",
		"mongodb.pwd":          "pwd",
		"mongodb.database":     "cmdb",
		"mongodb.port":         "27107",
		"mongodb.maxOpenConns": "3000",
		"mongodb.maxIDleConns": "1000",
		"redis.host":           "127.0.0.1",
		"redis.pwd":            "redisauth",
		"redis.database":       "0",
		"redis.port":           "6379",
		"redis.maxOpenConns":   "3000",
		"redis.maxIDleConns":   "1000",
		"errors.res":           "conf/errors",
	}
	err := ParseFromBytes(data)
	assert.NoError(t, err)
	assert.Equal(t, expect, Get())
}
func TestParseFromFile(t *testing.T) {
	err := ParseFromFile(`testdata/server.conf`)
	expect := Config{
		"mongodb.host":         "127.0.0.1",
		"mongodb.usr":          "user",
		"mongodb.pwd":          "pwd",
		"mongodb.database":     "cmdb",
		"mongodb.port":         "27107",
		"mongodb.maxOpenConns": "3000",
		"mongodb.maxIDleConns": "1000",
		"redis.host":           "127.0.0.1",
		"redis.pwd":            "redisauth",
		"redis.database":       "0",
		"redis.port":           "6379",
		"redis.maxOpenConns":   "3000",
		"redis.maxIDleConns":   "1000",
		"errors.res":           "conf/errors",
	}
	assert.NoError(t, err)
	assert.Equal(t, expect, Get())
}
