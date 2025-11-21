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

package orm

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

type options struct {
	// ingressLimiter db request limiter.
	ingressLimiter *rate.Limiter
	// mc db request metrics.
	mc *metric
	// slowRequestTime db slow request time, beyond this time, the db request will be logged.
	slowRequestTime time.Duration

	// debug will set logger level to info.
	debug bool
}

// Option orm option func defines.
type Option func(opt *options)

// IngressLimiter set db request limiter related params.
func IngressLimiter(qps, burst uint) Option {
	return func(opt *options) {
		opt.ingressLimiter = rate.NewLimiter(rate.Limit(qps), int(burst))
	}
}

// MetricsRegisterer set metrics registerer.
func MetricsRegisterer(register prometheus.Registerer) Option {
	return func(opt *options) {
		opt.mc = initMetric(register)
	}
}

// SlowRequest set db slow request time.
func SlowRequest(duration time.Duration) Option {
	return func(opt *options) {
		opt.slowRequestTime = duration
	}
}

// Debug set debug mode.
func Debug() Option {
	return func(opt *options) {
		opt.debug = true
	}
}
