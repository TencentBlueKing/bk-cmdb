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
	"context"
	"fmt"
	"net"
	"time"

	localRedis "configcenter/src/storage/dal/redis"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-redis/redis/v7"
)

// these test cases extract from cases in following files at git commit bfb4c92
// https://github.com/go-redis/redis/blob/master/commands_test.go
// https://github.com/go-redis/redis/blob/master/pubsub_test.go
// little modification on some test cases to make them available to test the encapsulated redis Client
var _ = Describe("Commands", func() {
	var client localRedis.Client
	ctx := context.Background()

	BeforeEach(func() {
		client = localRedis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		//client = localRedis.NewFailoverClient(&redis.FailoverOptions{
		//	MasterName:       "mymaster",
		//	SentinelAddrs:    []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		//	SentinelPassword: "",
		//	Password:         "",
		//})

		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	It("should PSubscribe", func() {
		pubsub := client.PSubscribe(ctx, "mychannel*")
		defer pubsub.Close()

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Subscription)
			Expect(subscr.Kind).To(Equal("psubscribe"))
			Expect(subscr.Channel).To(Equal("mychannel*"))
			Expect(subscr.Count).To(Equal(1))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err.(net.Error).Timeout()).To(Equal(true))
			Expect(msgi).To(BeNil())
		}

		n, err := client.Publish(ctx, "mychannel1", "hello").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(1)))

		Expect(pubsub.PUnsubscribe("mychannel*")).NotTo(HaveOccurred())

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Message)
			Expect(subscr.Channel).To(Equal("mychannel1"))
			Expect(subscr.Pattern).To(Equal("mychannel*"))
			Expect(subscr.Payload).To(Equal("hello"))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Subscription)
			Expect(subscr.Kind).To(Equal("punsubscribe"))
			Expect(subscr.Channel).To(Equal("mychannel*"))
			Expect(subscr.Count).To(Equal(0))
		}

	})

	It("should Subscribe", func() {
		pubsub := client.Subscribe(ctx, "mychannel", "mychannel2")
		defer pubsub.Close()

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Subscription)
			Expect(subscr.Kind).To(Equal("subscribe"))
			Expect(subscr.Channel).To(Equal("mychannel"))
			Expect(subscr.Count).To(Equal(1))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Subscription)
			Expect(subscr.Kind).To(Equal("subscribe"))
			Expect(subscr.Channel).To(Equal("mychannel2"))
			Expect(subscr.Count).To(Equal(2))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err.(net.Error).Timeout()).To(Equal(true))
			Expect(msgi).NotTo(HaveOccurred())
		}

		n, err := client.Publish(ctx, "mychannel", "hello").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(1)))

		n, err = client.Publish(ctx, "mychannel2", "hello2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(1)))

		Expect(pubsub.Unsubscribe("mychannel", "mychannel2")).NotTo(HaveOccurred())

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			msg := msgi.(*redis.Message)
			Expect(msg.Channel).To(Equal("mychannel"))
			Expect(msg.Payload).To(Equal("hello"))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			msg := msgi.(*redis.Message)
			Expect(msg.Channel).To(Equal("mychannel2"))
			Expect(msg.Payload).To(Equal("hello2"))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Subscription)
			Expect(subscr.Kind).To(Equal("unsubscribe"))
			Expect(subscr.Channel).To(Equal("mychannel"))
			Expect(subscr.Count).To(Equal(1))
		}

		{
			msgi, err := pubsub.ReceiveTimeout(time.Second)
			Expect(err).NotTo(HaveOccurred())
			subscr := msgi.(*redis.Subscription)
			Expect(subscr.Kind).To(Equal("unsubscribe"))
			Expect(subscr.Channel).To(Equal("mychannel2"))
			Expect(subscr.Count).To(Equal(0))
		}
	})

	It("should Pipeline", func() {
		pipe := client.Pipeline()
		ping := pipe.Ping()
		set := pipe.Set("aaa", 99, 0)
		get := pipe.Get("aaa")
		_, err := pipe.Exec()
		Expect(err).NotTo(HaveOccurred())

		Expect(ping.Err()).NotTo(HaveOccurred())
		Expect(ping.Val()).To(Equal("PONG"))
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))
		Expect(get.Err()).NotTo(HaveOccurred())
		Expect(get.Val()).To(Equal("99"))
	})

	It("should BRPop", func() {
		rPush := client.RPush(ctx, "list1", "a", "b", "c")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		bRPop := client.BRPop(ctx, 0, "list1", "list2")
		Expect(bRPop.Err()).NotTo(HaveOccurred())
		Expect(bRPop.Val()).To(Equal([]string{"list1", "c"}))
	})

	It("should BRPop blocks", func() {
		started := make(chan bool)
		done := make(chan bool)
		go func() {
			defer GinkgoRecover()

			started <- true
			brpop := client.BRPop(ctx, 0, "list")
			Expect(brpop.Err()).NotTo(HaveOccurred())
			Expect(brpop.Val()).To(Equal([]string{"list", "a"}))
			done <- true
		}()
		<-started

		select {
		case <-done:
			Fail("BRPop is not blocked")
		case <-time.After(time.Second):
			// ok
		}

		rPush := client.RPush(ctx, "list", "a")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		select {
		case <-done:
			// ok
		case <-time.After(time.Second):
			Fail("BRPop is still blocked")
			// ok
		}
	})

	It("should BRPopLPush", func() {
		_, err := client.BRPopLPush(ctx, "list1", "list2", time.Second).Result()
		Expect(err).To(Equal(redis.Nil))

		err = client.RPush(ctx, "list1", "a", "b", "c").Err()
		Expect(err).NotTo(HaveOccurred())

		v, err := client.BRPopLPush(ctx, "list1", "list2", 0).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal("c"))
	})

	It("should Del", func() {
		err := client.Set(ctx, "key1", "Hello", 0).Err()
		Expect(err).NotTo(HaveOccurred())
		err = client.Set(ctx, "key2", "World", 0).Err()
		Expect(err).NotTo(HaveOccurred())

		n, err := client.Del(ctx, "key1", "key2", "key3").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(2)))
	})

	It("should Eval returns keys and values", func() {
		vals, err := client.Eval(
			ctx,
			"return {KEYS[1],ARGV[1]}",
			[]string{"key"},
			"hello",
		).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(vals).To(Equal([]interface{}{"key", "hello"}))
	})

	It("should Eval returns all values after an error", func() {
		vals, err := client.Eval(
			ctx,
			`return {12, {err="error"}, "abc"}`,
			nil,
		).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(vals.([]interface{})[0]).To(Equal(int64(12)))
		Expect(vals.([]interface{})[1].(error).Error()).To(Equal("error"))
		Expect(vals.([]interface{})[2]).To(Equal("abc"))
	})

	It("should Get", func() {
		get := client.Get(ctx, "_")
		Expect(get.Err()).To(Equal(redis.Nil))
		Expect(get.Val()).To(Equal(""))

		set := client.Set(ctx, "key", "hello", 0)
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))

		get = client.Get(ctx, "key")
		Expect(get.Err()).NotTo(HaveOccurred())
		Expect(get.Val()).To(Equal("hello"))
	})

	It("should Exists", func() {
		set := client.Set(ctx, "key1", "Hello", 0)
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))

		n, err := client.Exists(ctx, "key1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(1)))

		n, err = client.Exists(ctx, "key2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(0)))

		n, err = client.Exists(ctx, "key1", "key2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(1)))

		n, err = client.Exists(ctx, "key1", "key1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(2)))
	})

	It("should HDel", func() {
		hSet := client.HSet(ctx, "hash", "key", "hello")
		Expect(hSet.Err()).NotTo(HaveOccurred())

		hDel := client.HDel(ctx, "hash", "key")
		Expect(hDel.Err()).NotTo(HaveOccurred())
		Expect(hDel.Val()).To(Equal(int64(1)))

		hDel = client.HDel(ctx, "hash", "key")
		Expect(hDel.Err()).NotTo(HaveOccurred())
		Expect(hDel.Val()).To(Equal(int64(0)))
	})

	It("should HGet", func() {
		hSet := client.HSet(ctx, "hash", "key", "hello")
		Expect(hSet.Err()).NotTo(HaveOccurred())

		hGet := client.HGet(ctx, "hash", "key")
		Expect(hGet.Err()).NotTo(HaveOccurred())
		Expect(hGet.Val()).To(Equal("hello"))

		hGet = client.HGet(ctx, "hash", "key1")
		Expect(hGet.Err()).To(Equal(redis.Nil))
		Expect(hGet.Val()).To(Equal(""))
	})

	It("should HGetAll", func() {
		err := client.HSet(ctx, "hash", "key1", "hello1").Err()
		Expect(err).NotTo(HaveOccurred())
		err = client.HSet(ctx, "hash", "key2", "hello2").Err()
		Expect(err).NotTo(HaveOccurred())

		m, err := client.HGetAll(ctx, "hash").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(m).To(Equal(map[string]string{"key1": "hello1", "key2": "hello2"}))
	})

	It("should HIncrBy", func() {
		hSet := client.HSet(ctx, "hash", "key", "5")
		Expect(hSet.Err()).NotTo(HaveOccurred())

		hIncrBy := client.HIncrBy(ctx, "hash", "key", 1)
		Expect(hIncrBy.Err()).NotTo(HaveOccurred())
		Expect(hIncrBy.Val()).To(Equal(int64(6)))

		hIncrBy = client.HIncrBy(ctx, "hash", "key", -1)
		Expect(hIncrBy.Err()).NotTo(HaveOccurred())
		Expect(hIncrBy.Val()).To(Equal(int64(5)))

		hIncrBy = client.HIncrBy(ctx, "hash", "key", -10)
		Expect(hIncrBy.Err()).NotTo(HaveOccurred())
		Expect(hIncrBy.Val()).To(Equal(int64(-5)))
	})

	It("should HScan", func() {
		for i := 0; i < 1000; i++ {
			sadd := client.HSet(ctx, "myhash", fmt.Sprintf("key%d", i), "hello")
			Expect(sadd.Err()).NotTo(HaveOccurred())
		}

		keys, cursor, err := client.HScan(ctx, "myhash", 0, "", 0).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(keys).NotTo(BeEmpty())
		Expect(cursor).NotTo(BeZero())
	})

	It("should HSet", func() {
		ok, err := client.HSet(ctx, "hash", map[string]interface{}{
			"key1": "hello1",
			"key2": "hello2",
		}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(Equal(int64(2)))

		v, err := client.HGet(ctx, "hash", "key1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal("hello1"))

		v, err = client.HGet(ctx, "hash", "key2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal("hello2"))

		keys, err := client.HKeys(ctx, "hash").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(keys).To(ConsistOf([]string{"key1", "key2"}))
	})

	It("should Incr", func() {
		set := client.Set(ctx, "key", "10", 0)
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))

		incr := client.Incr(ctx, "key")
		Expect(incr.Err()).NotTo(HaveOccurred())
		Expect(incr.Val()).To(Equal(int64(11)))

		get := client.Get(ctx, "key")
		Expect(get.Err()).NotTo(HaveOccurred())
		Expect(get.Val()).To(Equal("11"))
	})

	It("should Keys", func() {
		mset := client.MSet(ctx, "one", "1", "two", "2", "three", "3", "four", "4")
		Expect(mset.Err()).NotTo(HaveOccurred())
		Expect(mset.Val()).To(Equal("OK"))

		keys := client.Keys(ctx, "*o*")
		Expect(keys.Err()).NotTo(HaveOccurred())
		Expect(keys.Val()).To(ConsistOf([]string{"four", "one", "two"}))

		keys = client.Keys(ctx, "t??")
		Expect(keys.Err()).NotTo(HaveOccurred())
		Expect(keys.Val()).To(Equal([]string{"two"}))

		keys = client.Keys(ctx, "*")
		Expect(keys.Err()).NotTo(HaveOccurred())
		Expect(keys.Val()).To(ConsistOf([]string{"four", "one", "three", "two"}))
	})

	It("should LLen", func() {
		lPush := client.LPush(ctx, "list", "World")
		Expect(lPush.Err()).NotTo(HaveOccurred())
		lPush = client.LPush(ctx, "list", "Hello")
		Expect(lPush.Err()).NotTo(HaveOccurred())

		lLen := client.LLen(ctx, "list")
		Expect(lLen.Err()).NotTo(HaveOccurred())
		Expect(lLen.Val()).To(Equal(int64(2)))
	})

	It("should LPush", func() {
		lPush := client.LPush(ctx, "list", "World")
		Expect(lPush.Err()).NotTo(HaveOccurred())
		lPush = client.LPush(ctx, "list", "Hello")
		Expect(lPush.Err()).NotTo(HaveOccurred())

		lRange := client.LRange(ctx, "list", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"Hello", "World"}))
	})

	It("should LRange", func() {
		rPush := client.RPush(ctx, "list", "one")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "two")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "three")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		lRange := client.LRange(ctx, "list", 0, 0)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"one"}))

		lRange = client.LRange(ctx, "list", -3, 2)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"one", "two", "three"}))

		lRange = client.LRange(ctx, "list", -100, 100)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"one", "two", "three"}))

		lRange = client.LRange(ctx, "list", 5, 10)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{}))
	})

	It("should LRem", func() {
		rPush := client.RPush(ctx, "list", "hello")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "hello")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "key")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "hello")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		lRem := client.LRem(ctx, "list", -2, "hello")
		Expect(lRem.Err()).NotTo(HaveOccurred())
		Expect(lRem.Val()).To(Equal(int64(2)))

		lRange := client.LRange(ctx, "list", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"hello", "key"}))
	})

	It("should LTrim", func() {
		rPush := client.RPush(ctx, "list", "one")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "two")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "three")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		lTrim := client.LTrim(ctx, "list", 1, -1)
		Expect(lTrim.Err()).NotTo(HaveOccurred())
		Expect(lTrim.Val()).To(Equal("OK"))

		lRange := client.LRange(ctx, "list", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"two", "three"}))
	})

	It("should MSet MGet", func() {
		mSet := client.MSet(ctx, "key1", "hello1", "key2", "hello2")
		Expect(mSet.Err()).NotTo(HaveOccurred())
		Expect(mSet.Val()).To(Equal("OK"))

		mGet := client.MGet(ctx, "key1", "key2", "_")
		Expect(mGet.Err()).NotTo(HaveOccurred())
		Expect(mGet.Val()).To(Equal([]interface{}{"hello1", "hello2", nil}))
	})

	It("should Ping", func() {
		ping := client.Ping(ctx)
		Expect(ping.Err()).NotTo(HaveOccurred())
		Expect(ping.Val()).To(Equal("PONG"))
	})

	It("should Rename", func() {
		set := client.Set(ctx, "key", "hello", 0)
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))

		status := client.Rename(ctx, "key", "key1")
		Expect(status.Err()).NotTo(HaveOccurred())
		Expect(status.Val()).To(Equal("OK"))

		get := client.Get(ctx, "key1")
		Expect(get.Err()).NotTo(HaveOccurred())
		Expect(get.Val()).To(Equal("hello"))
	})

	It("should RenameNX", func() {
		set := client.Set(ctx, "key", "hello", 0)
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))

		renameNX := client.RenameNX(ctx, "key", "key1")
		Expect(renameNX.Err()).NotTo(HaveOccurred())
		Expect(renameNX.Val()).To(Equal(true))

		get := client.Get(ctx, "key1")
		Expect(get.Err()).NotTo(HaveOccurred())
		Expect(get.Val()).To(Equal("hello"))
	})

	It("should RPop", func() {
		rPush := client.RPush(ctx, "list", "one")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "two")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "three")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		rPop := client.RPop(ctx, "list")
		Expect(rPop.Err()).NotTo(HaveOccurred())
		Expect(rPop.Val()).To(Equal("three"))

		lRange := client.LRange(ctx, "list", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"one", "two"}))
	})

	It("should RPopLPush", func() {
		rPush := client.RPush(ctx, "list", "one")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "two")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		rPush = client.RPush(ctx, "list", "three")
		Expect(rPush.Err()).NotTo(HaveOccurred())

		rPopLPush := client.RPopLPush(ctx, "list", "list2")
		Expect(rPopLPush.Err()).NotTo(HaveOccurred())
		Expect(rPopLPush.Val()).To(Equal("three"))

		lRange := client.LRange(ctx, "list", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"one", "two"}))

		lRange = client.LRange(ctx, "list2", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"three"}))
	})

	It("should RPush", func() {
		rPush := client.RPush(ctx, "list", "Hello")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		Expect(rPush.Val()).To(Equal(int64(1)))

		rPush = client.RPush(ctx, "list", "World")
		Expect(rPush.Err()).NotTo(HaveOccurred())
		Expect(rPush.Val()).To(Equal(int64(2)))

		lRange := client.LRange(ctx, "list", 0, -1)
		Expect(lRange.Err()).NotTo(HaveOccurred())
		Expect(lRange.Val()).To(Equal([]string{"Hello", "World"}))
	})

	It("should SAdd", func() {
		sAdd := client.SAdd(ctx, "set", "Hello")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "set", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "set", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(0)))

		sMembers := client.SMembers(ctx, "set")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
	})

	It("should SAdd strings", func() {
		set := []string{"Hello", "World", "World"}
		sAdd := client.SAdd(ctx, "set", set)
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(2)))

		sMembers := client.SMembers(ctx, "set")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
	})

	It("should Set with expiration", func() {
		err := client.Set(ctx, "key", "hello", 100*time.Millisecond).Err()
		Expect(err).NotTo(HaveOccurred())

		val, err := client.Get(ctx, "key").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("hello"))

		Eventually(func() error {
			return client.Get(ctx, "foo").Err()
		}, "1s", "100ms").Should(Equal(redis.Nil))
	})

	It("should SetNX", func() {
		setNX := client.SetNX(ctx, "key", "hello", 0)
		Expect(setNX.Err()).NotTo(HaveOccurred())
		Expect(setNX.Val()).To(Equal(true))

		setNX = client.SetNX(ctx, "key", "hello2", 0)
		Expect(setNX.Err()).NotTo(HaveOccurred())
		Expect(setNX.Val()).To(Equal(false))

		get := client.Get(ctx, "key")
		Expect(get.Err()).NotTo(HaveOccurred())
		Expect(get.Val()).To(Equal("hello"))
	})

	It("should SetNX with expiration", func() {
		isSet, err := client.SetNX(ctx, "key", "hello", time.Second).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(isSet).To(Equal(true))

		isSet, err = client.SetNX(ctx, "key", "hello2", time.Second).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(isSet).To(Equal(false))

		val, err := client.Get(ctx, "key").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("hello"))
	})

	It("should SetNX with no expiration", func() {
		isSet, err := client.SetNX(ctx, "key", "hello1", 0).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(isSet).To(Equal(true))

		ttl := client.TTL(ctx, "key")
		Expect(ttl.Err()).NotTo(HaveOccurred())
		Expect(ttl.Val().Nanoseconds()).To(Equal(int64(-1)))
	})

	It("should SMembers", func() {
		sAdd := client.SAdd(ctx, "set", "Hello")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sMembers := client.SMembers(ctx, "set")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
	})

	It("should SRem", func() {
		sAdd := client.SAdd(ctx, "set", "one")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set", "two")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set", "three")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sRem := client.SRem(ctx, "set", "one")
		Expect(sRem.Err()).NotTo(HaveOccurred())
		Expect(sRem.Val()).To(Equal(int64(1)))

		sRem = client.SRem(ctx, "set", "four")
		Expect(sRem.Err()).NotTo(HaveOccurred())
		Expect(sRem.Val()).To(Equal(int64(0)))

		sMembers := client.SMembers(ctx, "set")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(ConsistOf([]string{"three", "two"}))
	})

	It("should TTL", func() {
		ttl := client.TTL(ctx, "key")
		Expect(ttl.Err()).NotTo(HaveOccurred())
		Expect(ttl.Val() < 0).To(Equal(true))

		set := client.Set(ctx, "key", "hello", 0)
		Expect(set.Err()).NotTo(HaveOccurred())
		Expect(set.Val()).To(Equal("OK"))

		expire := client.Expire(ctx, "key", 60*time.Second)
		Expect(expire.Err()).NotTo(HaveOccurred())
		Expect(expire.Val()).To(Equal(true))

		ttl = client.TTL(ctx, "key")
		Expect(ttl.Err()).NotTo(HaveOccurred())
		Expect(ttl.Val()).To(Equal(60 * time.Second))
	})

})
