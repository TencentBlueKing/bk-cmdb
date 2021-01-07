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
	"fmt"
	"math"
	"sort"
	"time"
)

type Statistic struct {
	// the total seconds of this statistic
	SustainSecond float64
	// the count of goroutines to run
	Concurrent int
	// total request that has been send.
	TotalRequest int64
	// success request number
	TotalSucceed int64
	// failed request number
	TotalFailed int64

	// in millisecond
	// only record success metrics
	TotalCostDuration float64
	Values            []float64
	// Request cost duration
	MinCostDuration float64
	MaxCostDuration float64
}

func (s *Statistic) CollectStatus(status *Status) {
	if status.Error != nil {
		s.IncreaseFailed()
	} else {
		s.IncreaseSuccess()
		s.UpdateCostDuration(status.CostDuration)
	}
}

func (s *Statistic) IncreaseRequest() {
	s.TotalRequest += 1
}

func (s *Statistic) IncreaseSuccess() {
	s.TotalSucceed += 1
}

func (s *Statistic) IncreaseFailed() {
	s.TotalFailed += 1
}

// only update success request data
func (s *Statistic) UpdateCostDuration(t time.Duration) {
	mt := t.Seconds() * 1000
	s.TotalCostDuration += mt
	if s.Values == nil {
		s.Values = make([]float64, 0)
	}
	s.Values = append(s.Values, mt)

	// update max
	if mt > s.MaxCostDuration {
		s.MaxCostDuration = mt
	}

	// update min
	if s.MinCostDuration == 0 {
		s.MinCostDuration = mt
	}

	if s.MinCostDuration > mt {
		s.MinCostDuration = mt
	}
}

type Metrics struct {
	SustainSeconds float64
	Concurrent     int
	QPS            float64
	// in millisecond
	MaxDuration       float64
	MinDuration       float64
	MedianDuration    float64
	AverageDuration   float64
	Percent85Duration float64
	Percent95Duration float64

	TotalRequest    int64
	SucceedRequest  int64
	FailedRequest   int64
	OnTheFlyRequest int64
}

func (m *Metrics) Format() string {
	var f string
	f = SetYellow("Load test metrics:\n--------------------\n")
	f += fmt.Sprintf("  Sustain: %ss\n", SetGreen(m.SustainSeconds))
	f += fmt.Sprintf("  Cocurrent: %s\n", SetGreen(m.Concurrent))
	f += fmt.Sprintf("  Total:   %s\n", SetGreen(m.TotalRequest))
	f += fmt.Sprintf("  Succeed: %s\n", SetGreen(m.SucceedRequest))
	f += fmt.Sprintf("  Failed:  %s\n", SetRed(m.FailedRequest))
	f += fmt.Sprintf("  Fly:     %s\n", SetYellow(m.OnTheFlyRequest))
	f += fmt.Sprintf("  QPS:     %s\n", SetGreen(fmt.Sprintf("%.1f", m.QPS)))
	f += fmt.Sprintf("  Max:     %sms\n", SetRed(fmt.Sprintf("%.1f", m.MaxDuration)))
	f += fmt.Sprintf("  Min:     %sms\n", SetGreen(fmt.Sprintf("%.1f", m.MinDuration)))
	f += fmt.Sprintf("  Med:     %sms\n", SetGreen(fmt.Sprintf("%.1f", m.MedianDuration)))
	f += fmt.Sprintf("  Avg:     %sms\n", SetGreen(fmt.Sprintf("%.1f", m.AverageDuration)))
	f += fmt.Sprintf("  P(85):   %sms\n", SetGreen(fmt.Sprintf("%.1f", m.Percent85Duration)))
	f += fmt.Sprintf("  P(95):   %sms\n", SetGreen(fmt.Sprintf("%.1f", m.Percent95Duration)))

	return f + "\n"
}

func (s *Statistic) CalculateMetrics() Metrics {

	var m Metrics
	m.SustainSeconds = s.SustainSecond
	m.Concurrent = s.Concurrent
	m.QPS = float64(float64(len(s.Values)) / s.SustainSecond)
	m.MaxDuration = s.MaxCostDuration
	m.MinDuration = s.MinCostDuration

	// sort the data.
	sort.Float64s(s.Values)

	if s.TotalSucceed == 0 {
		m.MedianDuration = 0
		m.AverageDuration = 0
	} else {
		// The median of an even number of values is the average of the middle two.
		if (s.TotalSucceed & 0x01) == 0 {
			m.MedianDuration = (s.Values[s.TotalSucceed/2-1] + s.Values[s.TotalSucceed/2]) / 2 / float64(time.Millisecond)
		} else {
			m.MedianDuration = s.Values[s.TotalSucceed/2]
		}
		m.AverageDuration = s.TotalCostDuration / float64(s.TotalSucceed)
	}
	m.Percent85Duration = s.percent(0.85)
	m.Percent95Duration = s.percent(0.95)

	m.TotalRequest = s.TotalRequest
	m.SucceedRequest = s.TotalSucceed
	m.FailedRequest = s.TotalFailed
	m.OnTheFlyRequest = s.TotalRequest - s.TotalSucceed - s.TotalFailed

	return m
}

func (s *Statistic) percent(p float64) float64 {
	num := len(s.Values)
	switch num {
	case 0:
		return 0
	case 1:
		return s.Values[0]
	default:
		sort.Float64s(s.Values)
		i := p * (float64(num) - 1.0)
		j := s.Values[int(math.Floor(i))]
		k := s.Values[int(math.Ceil(i))]
		f := i - math.Floor(i)
		return j + (k-j)*f
	}
}
