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

package rest

import (
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var rm *restMetric
var once sync.Once

func init() {
	once = sync.Once{}
	rm = new(restMetric)
}

// initMetric initialize metrics, and must be called for only once.
func initMetric() {
	rm.totalErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "response",
		Name:      "total_response_error_count",
		Help:      "the total error count responded to users",
	}, []string{"app_code", "uri"})
	metrics.Register().MustRegister(rm.totalErrorCount)

}

type restMetric struct {
	// recorded the total responded with error response count
	totalErrorCount *prometheus.CounterVec
}

func (c *Contexts) collectErrorMetric() {
	rm.totalErrorCount.With(prometheus.Labels{
		"app_code": c.Kit.Header.Get(common.BKHTTPRequestAppCode),
		"uri":      c.uri,
	}).Inc()
}
