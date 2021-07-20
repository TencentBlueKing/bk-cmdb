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

package main

import (
	"context"
	"fmt"
	"time"

	localRedis "configcenter/src/storage/dal/redis"

	"github.com/go-redis/redis/v7"
)

func main() {
	MyClient()
	MySentinelClient()
}

func MyClient() {
	client := localRedis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "cc",
		DB:       0,
	})

	DBOps(client)
}

func MySentinelClient() {
	sentinelClient := localRedis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       "mymaster",
		SentinelAddrs:    []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		SentinelPassword: "ss",
		Password:         "cc",
	})
	DBOps(sentinelClient)
}

func DBOps(cli localRedis.Client) {
	ctx := context.Background()
	key := "mykey"
	listName := "mylist"
	listName2 := "mylist2"
	hashKey := "myHashKey"
	setKey := "mySetKey"

	pipe := cli.Pipeline()
	pipe.Set("aaa", 99, 0)
	pipe.Get("aaa")
	vals, err := pipe.Exec()
	checkErr(err)
	fmt.Println("Pipeline", vals)

	err = cli.Set(ctx, key, "Hello,man!", 0).Err()
	checkErr(err)

	intVal, err := cli.Exists(ctx, key).Result()
	checkErr(err)
	fmt.Println("Exists", intVal)

	strVal, err := cli.Get(ctx, key).Result()
	checkErr(err)
	fmt.Println("Get", key, strVal)

	interfVal, err := cli.Eval(ctx, "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}", []string{"key1", "key2"}, "arg1", "arg2").Result()
	checkErr(err)
	fmt.Println("Eval:", interfVal)

	statusVal, err := cli.Ping(ctx).Result()
	checkErr(err)
	fmt.Println("Ping:", statusVal)

	cli.Set(ctx, "key1", "value111", 0)
	cli.Set(ctx, "key2", "value222", 0)
	interfSliVal, err := cli.MGet(ctx, "key1", "key2").Result()
	checkErr(err)
	fmt.Println("MGet:", interfSliVal)

	intVal, err = cli.Del(ctx, "key1", "key2").Result()
	checkErr(err)
	fmt.Println("Del:", intVal)

	sub := cli.Subscribe(ctx, "channels")
	checkErr(err)

	go func() {
		time.Sleep(time.Second)
		cli.Publish(ctx, "channels", "hello,a subscribe test")
	}()
	msg, err := sub.ReceiveMessage()
	checkErr(err)
	fmt.Println("ReceiveMessage:", msg)

	err = sub.Unsubscribe("channels")
	checkErr(err)

	err = sub.Close()
	checkErr(err)

	intVal, err = cli.Incr(ctx, "key").Result()
	checkErr(err)
	fmt.Println("Incr:", "key", intVal)

	intVal, err = cli.LPush(ctx, listName, "111").Result()
	checkErr(err)
	fmt.Println("LPush", listName, intVal)

	strSliVal, err := cli.BRPop(ctx, time.Second*30, listName).Result()
	checkErr(err)
	fmt.Println("BRPop", listName, strSliVal)

	cli.LPush(ctx, listName, "333")
	strVal, err = cli.BRPopLPush(ctx, listName, listName2, time.Second).Result()
	checkErr(err)
	fmt.Println("BRPopLPush", strVal)

	intVal, err = cli.LLen(ctx, listName2).Result()
	checkErr(err)
	fmt.Println("LLen", strVal)

	statusVal, err = cli.LTrim(ctx, listName2, 0, 100).Result()
	checkErr(err)
	fmt.Println("LTrim", strVal)

	intVal, err = cli.HIncrBy(ctx, hashKey, "field", 5).Result()
	checkErr(err)
	fmt.Println("HIncrBy", intVal)

	intVal, err = cli.SAdd(ctx, setKey, "m1", "m2", "m3").Result()
	checkErr(err)
	fmt.Println("SAdd", intVal)

	intVal, err = cli.SRem(ctx, setKey, "m2").Result()
	checkErr(err)
	fmt.Println("SRem", intVal)

	strSliVal, err = cli.SMembers(ctx, setKey).Result()
	checkErr(err)
	fmt.Println("SAdd", strSliVal)
}

func checkErr(err error) {
	if err != nil {
		if err == redis.Nil {
			return
		}
		panic(err)
	}
}
