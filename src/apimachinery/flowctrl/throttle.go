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

package flowctrl

import "github.com/juju/ratelimit"

// RateLimiter TODO
type RateLimiter interface {
	// TryAccept returns true if a token is taken immediately. Otherwise,
	// it returns false.
	TryAccept() bool

	// Accept will wait and not return unless a token becomes available.
	Accept()

	// QPS returns QPS of this rate limiter
	QPS() int64

	// Burst returns the burst of this rate limiter
	Burst() int64

	// AcceptMany will wait and not return unless the token becomes available.
	AcceptMany(count int64)
}

// NewRateLimiter TODO
func NewRateLimiter(qps, burst int64) RateLimiter {
	limiter := ratelimit.NewBucketWithRate(float64(qps), burst)
	return &tokenBucket{
		limiter: limiter,
		qps:     qps,
		burst:   burst,
	}
}

type tokenBucket struct {
	limiter *ratelimit.Bucket
	qps     int64
	burst   int64
}

// TryAccept TODO
func (t *tokenBucket) TryAccept() bool {
	return t.limiter.TakeAvailable(1) == 1
}

// Accept TODO
func (t *tokenBucket) Accept() {
	t.limiter.Wait(1)
}

// QPS TODO
func (t *tokenBucket) QPS() int64 {
	return t.qps
}

// Burst TODO
func (t *tokenBucket) Burst() int64 {
	return t.burst
}

// AcceptMany accept many token
func (t *tokenBucket) AcceptMany(count int64) {
	t.limiter.Wait(count)
}

// NewMockRateLimiter TODO
func NewMockRateLimiter() RateLimiter {
	return &mockRatelimiter{}
}

type mockRatelimiter struct{}

// TryAccept TODO
func (*mockRatelimiter) TryAccept() bool {
	return true
}

// Accept TODO
func (*mockRatelimiter) Accept() {

}

// QPS TODO
func (*mockRatelimiter) QPS() int64 {
	return 0
}

// Burst TODO
func (*mockRatelimiter) Burst() int64 {
	return 0
}

// AcceptMany accept many token
func (*mockRatelimiter) AcceptMany(count int64) {
}
