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

package plugin

import (
	"sync"

	"configcenter/src/common/metric"
)

type CounterInterface interface {
	metric.MetricInterf
	// increase the counter with value inc, and returned the increased value.
	Increase(inc float64) float64
	// decrease the counter with value dec, and returned the decreased value.
	Decrease(dec float64) float64
	// get the current counter value.
	GetCounter() float64
	// set the counter with val number.
	Set(val float64)
	// reset the counter with 0.
	Reset()
}

func NewCounterMetric(name, help string) CounterInterface {
	return &CounterMetric{
		name: name,
		help: help,
	}
}

var _ metric.MetricInterf = &CounterMetric{}

type CounterMetric struct {
	name string
	help string
	counter
}

func (cm *CounterMetric) GetMeta() *metric.MetricMeta {
	return &metric.MetricMeta{
		Name: cm.name,
		Help: cm.help,
	}
}

func (cm *CounterMetric) GetValue() (*metric.FloatOrString, error) {
	return metric.FormFloatOrString(cm.counter.GetCounter())
}

func (cm *CounterMetric) GetExtension() (*metric.MetricExtension, error) {
	return nil, nil
}

type counter struct {
	value  float64
	locker sync.RWMutex
}

func (c *counter) Increase(inc float64) float64 {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.value = c.value + inc
	return c.value
}

func (c *counter) Decrease(dec float64) float64 {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.value = c.value - dec
	return c.value
}

func (c *counter) GetCounter() float64 {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.value
}

func (c *counter) Set(val float64) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.value = val
	return
}

func (c *counter) Reset() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.value = 0
	return
}
