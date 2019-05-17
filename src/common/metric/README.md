# 简介
Metric SDK 包致力于为bcs体系内的各组件提供灵活、方便、可插拔式的运行时指标导出服务。指标导出服务以http服务的方式对外提供。Metric SDK支持基础的数值型（仅包括int8, int, int16, int32, int64, float32, float64）和字符串两种类型的key-value指标导出服务。对于数值型的指标，为了兼顾不同场景，SDK统一以float64对外展示。

各组件可根据自身需求定制自身需要导出的指标信息，Metric SDK 的主要包含以下特性：
 - 目前仅支持`golang`语言，其它语言暂不支持。
 - 默认提供该组件的runtime metric。如：使用的CPU数量、goroutine数量和内存等信息。
 - 非侵入式接口设计。核心指标数据通过回调函数进行采集，在满足接口间隔离、解耦的同时，方便后续的功能扩展，尽量避免SDK的更新、升级、迭代会影响到相关的使用组件。
 - `非阻塞`式接口设计。由于无法保证各组件在使用接口时的具体行为，metric在处理每一个metric时会默认设置一个`超时时间`为`5s`，保证所有的metric信息不会因为某一个metric的阻塞而造成拉取metric信息失败。所以使用SDK的组件在实现相关接口时尽量采用`非阻塞`式设计。
 - 提供metric`分组管理`机制，各组件可根据自身情况进行metric的分类管理和导出。
 - 对于同一组件，metric SDK 也提供各实例间的`个性化标识`功能，方便指标数据汇聚后的过滤、清洗等。该功能通过label(key, value)实现。
 　　考虑以下场景，bcs-health在深圳、上海、成都均部署有，那么当指标数据汇聚以后如何区分不同的bcs-health集群和实例呢，我们可以在启动bcs-health时配置不同的label，如("set", "shenzhen")，("set", "shanghai")等。

 - 对于常见的数值型metric，提供基础的`数值型metric接入插件`, 方便各组件快速构建。

# Golang Metric

|指标                                 |     意义                                                                 |
|-------------------------------------|------------------------------------------------------------------------|
|go_goroutines                        |Number of goroutines that currently exist.                           |
|go_threads                           |Number of OS threads created                                         |
|go_cpu_used                          |The number of logical CPUs usable by the current process            |
|go_memstats_alloc_bytes              |Number of bytes allocated and still in use.                          |
|go_memstats_alloc_bytes_total        |Total number of bytes allocated, even if freed.                      |
|go_memstats_sys_bytes                |Number of bytes obtained from system.                                |
|go_memstats_mallocs_total            |Total number of mallocs.                                             |
|go_memstats_frees_total              |Total number of frees.                                               |
|go_memstats_lookups_total            |Total number of pointer lookups.                                     |
|go_memstats_heap_alloc_bytes         |Number of heap bytes allocated and still in use.                     |
|go_memstats_heap_sys_bytes           |Number of heap bytes obtained from system.                           |
|go_memstats_heap_idle_bytes          |Number of heap bytes waiting to be used.                             |
|go_memstats_heap_inuse_bytes         |Number of heap bytes that are in use.                                |
|go_memstats_heap_released_bytes      |Number of heap bytes released to OS.                                 |
|go_memstats_heap_objects             |Number of allocated objects.                                         |
|go_memstats_stack_inuse_bytes        |Number of bytes in use by the stack allocator.                       |
|go_memstats_stack_sys_bytes          |Number of bytes obtained from system for stack allocator.            |
|go_memstats_mspan_inuse_bytes        |Number of bytes in use by mspan structures.                          |
|go_memstats_mspan_sys_bytes          |Number of bytes used for mspan structures obtained from s            |
|go_memstats_mcache_inuse_bytes       |Number of bytes in use by mcache structures.                         |
|go_memstats_mcache_sys_bytes         |Number of bytes used for mcache structures obtained from             |
|go_memstats_buck_hash_sys_bytes      |Number of bytes used by the profiling bucket hash table.             |
|go_memstats_gc_sys_bytes             |Number of bytes used for garbage collection system metada            |
|go_memstats_other_sys_bytes          |Number of bytes used for other system allocations.                   |
|go_memstats_next_gc_bytes            |Number of heap bytes when next garbage collection will take place    |
|go_memstats_last_gc_time_seconds     |Number of seconds since 1970 of last garbage collection.             |
|go_memstats_gc_cpu_fraction          |The fraction of this program's available CPU time used by the GC since the program started.           |

# 设计原理
使用metric SDK的组件在启用metric导出服务时，仅调用函数[metric.NewMetricController()] (./api.go#3)。调用NewMetricController()会完成以下动作：

 - 注册该组件的`元数据`信息，如下。具体含义如下：
    ```golang
    type Config struct {
    	Module     string            
    	IP         string            
    	MetricPort uint16            
    	ClusterID  string            
    	Labels     map[string]string 
    }
    ```

   - `Module`即为使用SDK模块的名称（如bcs-health），属必填字段；
   - `IP`为该组件所在主机的物理IP地址，同时也是metric提供http服务时所绑定的IP地址，属必填字段。
   - `ClusterID`为该组件的获取的bcs Cluster ID值，属选填字段，建议有该属性的组件必填，没有的不填。方便数据的后续过滤、清洗。
   - `MetricPort`为该组件监听的http端口。
   - `Labels`为组件的不同实例提供个性化定义的需求。属选填字段。

 - 注册用户自定义的Collector，即为用户的`metric组`，每个组又包含不同的具体metric信息。如下:
     ```golang
     type CollectorName string
     type Collector struct {
        Name      CollectorName
        Collector CollectInter
     }
     ```
   - `Name`: 即为**metric 组**的名称，同一组件内不可重复。
   - `Collector`: 即为获取该**metric**指标的接口。具体接口为：
    ```golang
    type CollectInter interface {
        Collect()[]MetricInterf
    }
    ```
   - `指标获取`是通过调用用户实现的函数`Collect()[]MetricInterf`来实现。
   - 考虑不同指标的获取方法不同，同时为用户提供灵活的实现方法，具体每个指标数据的获取也通过接口（interface）的方式来获取。每个指标只要实现了`MetricInterf接口`, Metric SDK即可抓取数据。目前`GetExtension()`接口为预留接口，返回的数据为空。
   ```golang
   type MetricInterf interface {
   	    GetMeta() MetricMeta
   	    GetValue() (*FloatOrString, error)
   	    GetExtension() (*MetricExtension, error)
   }
   ```

# Demo
　　下面给大家展示一下具体的实用方法，源码在[这里] (fake/main.go)。

　　该Demo为大家展示以下信息：
- 如何使用FloatOrString接口。
- 如何实例化Metric SDK。
- 如何使用Metric SDK提供的插件。

```golang
package main

import (
	"fmt"
	"time"
	"math/rand"
	
	"configcenter/common/metric/plugin"
	"configcenter/common/metric"
)

func main() {

	// 生成metric配置数据信息
   conf := metric.Config{
		ModuleName: "DemoModule",
		RunMode: metric.Master_Master_Mode,
		MetricPort: 9980,
		IP: "127.0.0.1",
		ClusterID: "bcs-test-cluster-id",
		Labels: map[string]string {
			"demoLabelKey": "demoLabelValue",
		},
	}
	
	healthFunc := func() metric.HealthMeta {
		return metric.HealthMeta{
			IsHealthy: true,
			CurrentRole: metric.MasterRole,
			Message: "is healthy now.",
		}
	}
	
	// 注册相关元数据信息、指标组，并拉起metric SDK。
	if err := metric.NewMetricController(
		conf,
		healthFunc,
		// 注册DemoCollector指标组，该指标组包含数据类型指标和字符类型指标的
		// 具体使用方法。
		metric.NewCollector("DemoCollector", NewIntStringDemo()),
		// 注册plugin demo指标组。描述计数类metric的插件使用方法。
		NewCounterCollector()); nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("start metric controller success.")
	select {}
}


func NewCounterCollector() *metric.Collector {
	cm := plugin.NewCounterMetric("money_in_bank", "this metric describe the money saved in bank.")
	pw := plugin.NewCounterMetric("password", "this metric describe the current password.")
	cc := &CounterController{
		Money: cm,
		Password: pw,
	}

	go cc.Sync()
	return metric.NewCollector("counter_metric_demo", cc)
}

type CounterController struct {
	Money plugin.CounterInterface
	Password plugin.CounterInterface
}

func (cc *CounterController)Sync()  {
	cc.Money.Set(10000)
	for {
		time.Sleep(time.Second * time.Duration(1))
		cc.Money.Increase(rand.Float64() * 100)

		time.Sleep(time.Second * time.Duration(1))
		cc.Money.Decrease(rand.Float64() * 100)
		cc.Password.Set(rand.Float64() * 100000000000000000)
	}
}

func (cc CounterController) Collect()[]metric.MetricInterf {
	m := []metric.MetricInterf{cc.Money, cc.Password}
	return m
}



func NewIntStringDemo() *IntStringMetricDemo {
	return &IntStringMetricDemo{
		metrics: []DemoMetric{
			{
				Name: "IntMetric",
				Help: "this metric describe a int metric like number of cpu cost.",
				GetFunc: func()(*metric.FloatOrString, error) {return metric.FormFloatOrString(8888)},
			},
			{
				Name: "StringMetric",
				Help: "this metric gives you more ways to describe you process, like it's runing well.",
				GetFunc: func()(*metric.FloatOrString, error){return metric.FormFloatOrString("bcs is running well.")},
			},
		},
	}
}




type IntStringMetricDemo struct{
	metrics []DemoMetric
}

func(d IntStringMetricDemo)Collect()[]metric.MetricInterf {
	m := make([]metric.MetricInterf, 0)
	for idx := range d.metrics {
		m = append(m, &d.metrics[idx])
	}
	return m
}


type DemoMetric struct{
	Name string
	Help string
	GetFunc func() (*metric.FloatOrString, error)
}

func(d DemoMetric) GetMeta() metric.MetricMeta {
	return metric.MetricMeta{
		Name: d.Name,
		Help: d.Help,
	}
}

func(d DemoMetric) GetValue() (*metric.FloatOrString, error) {
	return d.GetFunc()
}

func(d DemoMetric) GetExtension() (*metric.MetricExtension, error){
	return nil, nil
}

```

　　返回数据数据示例：
```json
{
    "metaData": {
        "module": "DemoModule",
        "ipAddr": "127.0.0.1",
        "clusterID": "bcs-test-cluster-id",
        "label": {
            "demoLabelKey": "demoLabelValue"
        }
    },
    "metricBundle": {
        "DemoCollector": [
            {
                "name": "IntMetric",
                "help": "this metric describe a int metric like number of cpu cost.",
                "value": 8888,
                "extension": null
            },
            {
                "name": "StringMetric",
                "help": "this metric gives you more ways to describe you process, like it's runing well.",
                "value": "bcs is running well.",
                "extension": null
            }
        ],
        "counter_metric_demo": [
            {
                "name": "money_in_bank",
                "help": "this metric describe the money saved in bank.",
                "value": 10010.186538712158,
                "extension": null
            },
            {
                "name": "password",
                "help": "this metric describe the current password.",
                "value": 66456005321849040,
                "extension": null
            }
        ],
        "golang_metrics": [
            {
                "name": "go_goroutines",
                "help": "Number of goroutines that currently exist.",
                "value": 6,
                "extension": null
            },
            {
                "name": "go_threads",
                "help": "Number of OS threads created",
                "value": 7,
                "extension": null
            },
            {
                "name": "go_memstats_alloc_bytes",
                "help": "Number of bytes allocated and still in use.",
                "value": 397024,
                "extension": null
            },
            {
                "name": "go_memstats_alloc_bytes_total",
                "help": "Total number of bytes allocated, even if freed.",
                "value": 397024,
                "extension": null
            },
            {
                "name": "go_memstats_sys_bytes",
                "help": "Number of bytes obtained from system.",
                "value": 3084288,
                "extension": null
            },
            {
                "name": "go_memstats_mallocs_total",
                "help": "Total number of mallocs.",
                "value": 5047,
                "extension": null
            },
            {
                "name": "go_memstats_frees_total",
                "help": "Total number of frees.",
                "value": 111,
                "extension": null
            }
        ]
    },
    "reportTimeMs": 1509099991
}
```


