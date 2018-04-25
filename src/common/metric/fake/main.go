package main

import (
	"fmt"
	"math/rand"
	"time"

	"bcs/bcs-common/common/metric"
	"bcs/bcs-common/common/metric/plugin"
)

func main() {

	// 生成metric配置数据信息
	conf := metric.Config{
		ModuleName: "DemoModule",
		RunMode:    metric.Master_Master_Mode,
		MetricPort: 9980,
		IP:         "127.0.0.1",
		ClusterID:  "bcs-test-cluster-id",
		Labels: map[string]string{
			"demoLabelKey": "demoLabelValue",
		},
	}

	healthFunc := func() metric.HealthMeta {
		return metric.HealthMeta{
			IsHealthy:   true,
			CurrentRole: metric.MasterRole,
			Message:     "is healthy now.",
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
		Money:    cm,
		Password: pw,
	}

	go cc.Sync()
	return metric.NewCollector("counter_metric_demo", cc)
}

type CounterController struct {
	Money    plugin.CounterInterface
	Password plugin.CounterInterface
}

func (cc *CounterController) Sync() {
	cc.Money.Set(10000)
	for {
		time.Sleep(time.Second * time.Duration(1))
		cc.Money.Increase(rand.Float64() * 100)

		time.Sleep(time.Second * time.Duration(1))
		cc.Money.Decrease(rand.Float64() * 100)
		cc.Password.Set(rand.Float64() * 100000000000000000)
	}
}

func (cc CounterController) Collect() []metric.MetricInterf {
	m := []metric.MetricInterf{cc.Money, cc.Password}
	return m
}

func NewIntStringDemo() *IntStringMetricDemo {
	return &IntStringMetricDemo{
		metrics: []DemoMetric{
			{
				Name:    "IntMetric",
				Help:    "this metric describe a int metric like number of cpu cost.",
				GetFunc: func() (*metric.FloatOrString, error) { return metric.FormFloatOrString(8888) },
			},
			{
				Name:    "StringMetric",
				Help:    "this metric gives you more ways to describe you process, like it's runing well.",
				GetFunc: func() (*metric.FloatOrString, error) { return metric.FormFloatOrString("bcs is running well.") },
			},
		},
	}
}

type IntStringMetricDemo struct {
	metrics []DemoMetric
}

func (d IntStringMetricDemo) Collect() []metric.MetricInterf {
	m := make([]metric.MetricInterf, 0)
	for idx := range d.metrics {
		m = append(m, &d.metrics[idx])
	}
	return m
}

type DemoMetric struct {
	Name    string
	Help    string
	GetFunc func() (*metric.FloatOrString, error)
}

func (d DemoMetric) GetMeta() metric.MetricMeta {
	return metric.MetricMeta{
		Name: d.Name,
		Help: d.Help,
	}
}

func (d DemoMetric) GetValue() (*metric.FloatOrString, error) {
	return d.GetFunc()
}

func (d DemoMetric) GetExtension() (*metric.MetricExtension, error) {
	return nil, nil
}
