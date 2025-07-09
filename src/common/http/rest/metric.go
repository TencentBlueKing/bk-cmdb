/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package rest

import (
	"sync"

	httpheader "configcenter/src/common/http/header"
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

	rm.noPermissionRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cmdb_no_permission_request_total",
			Help: "total number of request without permission.",
		},
		[]string{metrics.LabelHandler, metrics.LabelAppCode},
	)
	metrics.Register().MustRegister(rm.noPermissionRequestTotal)
}

type restMetric struct {
	// recorded the total responded with error response count
	totalErrorCount *prometheus.CounterVec
	// noPermissionRequestTotal is the total number of request without permission
	noPermissionRequestTotal *prometheus.CounterVec
}

func (c *Contexts) collectErrorMetric() {
	rm.totalErrorCount.With(prometheus.Labels{
		"app_code": httpheader.GetAppCode(c.Kit.Header),
		"uri":      c.uri,
	}).Inc()
}

func (c *Contexts) collectNoAuthMetric() {
	rm.noPermissionRequestTotal.With(
		prometheus.Labels{
			metrics.LabelHandler: c.Request.Request.URL.Path,
			metrics.LabelAppCode: httpheader.GetAppCode(c.Kit.Header),
		},
	).Inc()
}
