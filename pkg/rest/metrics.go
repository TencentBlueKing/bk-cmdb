/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"handler", "method", "code"},
	)
	responseTimeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for HTTP requests.",
			Buckets: []float64{0.1, 0.2, 0.5, 1, 2, 5, 10, 30, 60},
		},
		[]string{"handler", "method", "code"},
	)
)

// getHandleName 获取FuncHandle/StreamHandle函数名
func getHandleName(fn any) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if fullName == "" {
		panic("get func name is empty")
	}

	parts := strings.Split(fullName, ".")
	lastPart := parts[len(parts)-1]
	name := strings.TrimSuffix(lastPart, "-fm")
	return name
}

// collectHandleMetrics api指标数据
func collectHandleMetrics(funcName, method string, st time.Time, err error) {
	code := http.StatusOK
	if err != nil {
		var respErr *cerr.RespError
		if errors.As(err, &respErr) {
			code = cerr.GetHTTPStatus(respErr.Code)
		} else {
			code = http.StatusInternalServerError
		}
	}

	codeStr := strconv.Itoa(code)
	requestCounter.WithLabelValues(funcName, method, codeStr).Inc()
	duration := time.Since(st).Seconds()
	responseTimeDuration.WithLabelValues(funcName, method, codeStr).Observe(duration)
}

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(responseTimeDuration)
}
