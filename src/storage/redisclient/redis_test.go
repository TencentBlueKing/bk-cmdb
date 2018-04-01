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

package redisclient

import (
	"configcenter/src/common"
	"fmt"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {

	redis, err := NewRedis("127.0.0.1", "6379", "", "redisauth", "5")
	if nil != err {
		t.Errorf(err.Error())
		return
	}
	redis.Open()
	expire := time.Minute * 120
	var result interface{}
	//value, ok := mapData["value"]
	//exp, ok := mapData["expire"]
	fmt.Println("set ")
	data := common.KvMap{"key": "set-key", "value": "test-key"}
	redis.Insert("Set", data)
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "set-key"}, &result)
	fmt.Println(result)

	fmt.Println("\n===\nsetNx ")

	redis.Insert("SetNx", common.KvMap{"key": "set-keynx", "value": "test-key-nx", "expire": expire})
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "set-key"}, &result)
	fmt.Println(result)
	redis.Insert("SetNx", common.KvMap{"key": "set-keynx", "value": "test-key-setnx1", "expire": expire})
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "set-key"}, &result)
	fmt.Println(result)

	//decr
	fmt.Println("\n===\n decr ")
	redis.Insert("Set", common.KvMap{"key": "int-key", "value": "1000"})
	redis.Insert("Decr", common.KvMap{"key": "int-key"})
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "int-key"}, &result)
	fmt.Println(result)

	fmt.Println("\n===\n decr bys ")
	redis.Insert("DecrBy", common.KvMap{"key": "int-key", "decr": 50})
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "int-key"}, &result)
	fmt.Println(result)

	//incr
	fmt.Println("\n===\n incr ")
	redis.Insert("incr", common.KvMap{"key": "int-key"})
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "int-key"}, &result)
	fmt.Println(result)

	fmt.Println("\n===\n incr bys ")
	redis.Insert("incrby", common.KvMap{"key": "int-key", "incr": 50})
	redis.GetOneByCondition("get", nil, common.KvMap{"key": "int-key"}, &result)
	fmt.Println(result)
	//expire
	data = common.KvMap{"key": "expire-key", "value": "+======"}

	redis.Insert("Set", data)
	redis.Insert("expire", common.KvMap{"key": "expire-key", "expire": expire})
	redis.GetOneByCondition("ttl", nil, common.KvMap{"key": "expire-key"}, &result)

	//expireat
	fmt.Println("\n===\n expire at ")
	_, err = redis.Insert("Set", common.KvMap{"key": "expire-key-at", "value": 10001110000000})
	_, err = redis.Insert("expireat", common.KvMap{"key": "expire-key-at", "expire": time.Now().Add(10000 * time.Second)})
	redis.GetOneByCondition("ttl", nil, common.KvMap{"key": "expire-key-at"}, &result)
	fmt.Println(result)
	//hset
	fmt.Println("\n===\n hset ")
	redis.Insert("hset", common.KvMap{"key": "hset-key", "field": "hset-key1", "value": "key1"})
	err = redis.GetOneByCondition("hget", nil, common.KvMap{"key": "hset-key", "field": "hset-key1"}, &result)
	fmt.Println("hget", result)
	redis.GetOneByCondition("hlen", nil, common.KvMap{"key": "hset-key"}, &result)
	fmt.Println("hlen", result)
	redis.GetOneByCondition("hexists", nil, common.KvMap{"key": "hset-key", "field": "hset-key1"}, &result)
	fmt.Println(result)
	redis.GetOneByCondition("hexists", nil, common.KvMap{"key": "hset-key"}, &result)
	fmt.Println(result)

	//hmset
	fmt.Println("\n===\n hmset ")
	redis.Insert("hmset", common.KvMap{"key": "hset-key", "fields": map[string]string{"hmset-key2": "key2", "hmset-key3": "key3"}})
	redis.GetOneByCondition("hgetall", nil, common.KvMap{"key": "hset-key", "field": "hset-key1"}, &result)
	fmt.Println(result)
	//hsetnx
	fmt.Println("\n===\n hsetnx ")
	redis.Insert("hsetnx", common.KvMap{"key": "hsetnx-key", "field": "hsetnx", "value": "value"})
	redis.Insert("hsetnx", common.KvMap{"key": "hsetnx-key", "field": "hsetnx", "value": "value-hsetnx"})
	redis.GetOneByCondition("hgetall", nil, common.KvMap{"key": "hset-key"}, &result)
	fmt.Println(result)
	redis.GetOneByCondition("hgetall", nil, common.KvMap{"key": "hset-key"}, &result)
	fmt.Println(result)
	//rpush
	fmt.Println("\n===\n rpush ")
	redis.Insert("rpush", common.KvMap{"key": "rpush-key", "values": []interface{}{1, 2, 3, 4, 5}})
	redis.GetOneByCondition("lrange", nil, common.KvMap{"key": "rpush-key", "start": 0, "end": 200}, &result)
	fmt.Println(result)
	//lset
	//index, _ := mapData["index"].(int64)
	//value, _ := mapData["value"]
	fmt.Println("\n===\n lset ")
	redis.Insert("lset", common.KvMap{"key": "rpush-key", "value": 1000})
	redis.GetOneByCondition("lrange", nil, common.KvMap{"key": "rpush-key", "start": 0, "end": 200}, &result)
	fmt.Println(result)

	fmt.Println("\n===\n blpop ")
	redis.GetOneByCondition("blpop", nil, common.KvMap{"key": []string{"rpush-key"}, "expire": time.Duration(0)}, &result)
	fmt.Println(result)
	fmt.Println("\n===\n exists ")
	redis.GetOneByCondition("exists", nil, common.KvMap{"key": "rpush-key"}, &result)
	fmt.Println(result)
	redis.GetOneByCondition("exists", nil, common.KvMap{"key": "rpush-key1"}, &result)
	fmt.Println(result)

	redis.GetOneByCondition("exists", nil, common.KvMap{"key": "rpush-key"}, &result)
	fmt.Println(result)
	redis.GetOneByCondition("del", nil, common.KvMap{"key": []string{"rpush-key"}}, &result)
	fmt.Println(result)
	redis.GetOneByCondition("exists", nil, common.KvMap{"key": "rpush-key"}, &result)
	fmt.Println(result)

	fmt.Println("\n expire ", expire)
}
