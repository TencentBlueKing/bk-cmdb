## 性能测试
> 提供针对模块服务性能测试的工具

## 压测方式

* `redis源压测`：通过工具向redis压入数据，DataCollection从redis SubPub MQ接收并处理数据，该压测方式与真实部署场景一致;
* `mock压测`: DataCollection提供了MockServer，可以通过发送批量的HTTPMock指令进行压测，该压测方式只可考察模块内部队列数据处理能力，与redis之间数据接收环境未在其中;

## 工具介绍

```shell
> benchmark --help
Usage:
  benchmark [command]

  Available Commands:
    mq          mq benchmark command

Use "benchmark [command] --help" for more information about a command.


> benchmark mq --help
mq benchmark command

Usage:
  benchmark mq [flags]

  Flags:
        --conn=1: Num of clients.
        --data="{}": Message for benchmark.
        --file="": File path of meessage for benchmark, eg ./data.json.
        --host="localhost:6379": Host of redis server.
        --num=1: Num of requests.
        --topic="benchmark": Topic of redis SubPub MQ.
```

* 工具中现提供了`redis源压测`的方式，即mq命令, 通过该命令可以进行有效的性能压力测试
