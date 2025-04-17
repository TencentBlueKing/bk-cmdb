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

package util

import (
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common/ssl"

	"github.com/prometheus/client_golang/prometheus"
)

// APIMachineryConfig TODO
type APIMachineryConfig struct {
	// request's qps value
	QPS int64
	// request's burst value
	Burst     int64
	TLSConfig *ssl.TLSClientConfig
	ExtraConf *ExtraClientConfig
}

// Capability TODO
type Capability struct {
	Client     HttpClient
	Discover   discovery.Interface
	Throttle   flowctrl.RateLimiter
	Mock       MockInfo
	MetricOpts MetricOption
	// the max tolerance api request latency time, if exceeded this time, then
	// this request will be logged and warned.
	ToleranceLatencyTime time.Duration
}

// MetricOption TODO
type MetricOption struct {
	// prometheus metric register
	Register prometheus.Registerer
	// if not set, use default buckets value
	DurationBuckets []float64
}

// MockInfo TODO
type MockInfo struct {
	Mocked      bool
	SetMockData bool
	MockData    interface{}
}

// ExtraClientConfig extra http client configuration
type ExtraClientConfig struct {
	// ResponseHeaderTimeout the amount of time to wait for a server's response headers
	ResponseHeaderTimeout time.Duration
}
