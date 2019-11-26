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

package metrics

import (
	"math"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	valGauge prometheus.GaugeFunc
	maxGauge prometheus.GaugeFunc

	descC    chan<- *prometheus.Desc
	collectC chan<- *prometheus.Metric

	val uint64
	max uint64
}

func NewGauge(opt prometheus.GaugeOpts) *Gauge {
	g := Gauge{}
	g.valGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: opt.Name,
			Help: opt.Help,
		},
		func() float64 { return math.Float64frombits(atomic.LoadUint64(&g.val)) },
	)

	g.maxGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: opt.Name + "_max",
			Help: "max number of " + opt.Name,
		},
		func() float64 { return math.Float64frombits(atomic.LoadUint64(&g.max)) },
	)

	return &g
}

func (g *Gauge) Describe(ch chan<- *prometheus.Desc) {
	g.valGauge.Describe(ch)
	g.maxGauge.Describe(ch)
}

func (g *Gauge) Collect(ch chan<- prometheus.Metric) {
	g.valGauge.Collect(ch)
	g.maxGauge.Collect(ch)
}

func (g *Gauge) Inc() {
	new := g.Add(1)
	old := atomic.LoadUint64(&g.max)
	if new > old {
		atomic.CompareAndSwapUint64(&g.max, old, new)
	}
}
func (g *Gauge) Dec() {
	g.Add(-1)
}

func (g *Gauge) Add(val float64) uint64 {
	for {
		oldBits := atomic.LoadUint64(&g.val)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&g.val, oldBits, newBits) {
			return newBits
		}
	}
}
