## 实现方式
 - 通过对底层redis client的封装，实现接口Client来提供redis常用的一些操作
 - 当前封装的底层redis client是[go-redis](https://github.com/go-redis/redis/tree/v7)，git地址为<https://github.com/go-redis/redis/tree/v7>

## 实现目的
- 在底层redis client之上增加一层接口封装，该接口供上层调用，有利于代码的扩展性
- 在底层redis client发生替换或变化时，无需对Client接口协议协议进行变更，只需改变接口方法内的具体实现，增强了代码的可维护性


## 快速上手

``` go
package main

import (
	"context"
	"fmt"

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
		Password: "",
		DB:       0,
	})

	DBOps(client)
}

func MySentinelClient() {
	sentinelClient := localRedis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       "mymaster",
		SentinelAddrs:    []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		SentinelPassword: "",
		Password:         "",
	})
	DBOps(sentinelClient)
}

func DBOps(cli localRedis.Client) {
	ctx := context.Background()
	key := "mykey"

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
}

func checkErr(err error) {
	if err != nil {
		if err == redis.Nil {
			return
		}
		panic(err)
	}
}
```

## 样例

更多例子见 [example.go](./example/example.go)