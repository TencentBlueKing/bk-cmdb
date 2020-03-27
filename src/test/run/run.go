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

package run

import (
	"time"
)

var Concurrent int
var SustainSeconds int

type Status struct {
	CostDuration time.Duration
	Error        error
}

func FireLoadTest(f func() error) Metrics {
	limiter := NewStreamLimiter(Concurrent)
	timeout := make(chan bool)
	go func() {
		select {
		case <-time.After(time.Duration(SustainSeconds) * time.Second):
			close(timeout)
		}
	}()

	stats := new(Statistic)
	stats.SustainSecond = SustainSeconds
	ch := make(chan *Status, 3000)
	done := make(chan bool)
	go func() {
		// start collect request metrics

		for {
			select {
			case <-timeout:
				goto outer

			case s := <-ch:
				stats.CollectStatus(s)
			}
		}
	outer:
		// delay 5 seconds to wait for the requests on the fly.
		delay := time.After(5 * time.Second)
		// fmt.Println("wait for request on the fly.")
		for {
			select {
			case <-delay:
				done <- true
				return
			case s := <-ch:
				stats.CollectStatus(s)
			}
		}

		return
	}()

exitFor:
	for {
		select {
		case <-timeout:
			break exitFor
		default:
			limiter.Execute(ch, f)
			stats.IncreaseRequest()
			// time.Sleep(2 * time.Millisecond)
		}
	}

	<-done
	// it's time to calculate the metrics
	return stats.CalculateMetrics()
}
