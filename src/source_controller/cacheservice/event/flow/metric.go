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

package flow

import (
	"time"

	"configcenter/src/common/metrics"
	"configcenter/src/storage/stream/types"
	"github.com/prometheus/client_golang/prometheus"
)

func initialMetrics(collection string) *eventMetrics {
	labels := prometheus.Labels{"collection": collection}

	m := new(eventMetrics)
	m.totalEventCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   "watch",
		Name:        "total_event_count",
		Help:        "the total event count which we handled with this resources",
		ConstLabels: labels,
	}, []string{"action"})
	metrics.Register().MustRegister(m.totalEventCount)

	m.lastEventTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   "watch",
		Name:        "last_event_unix_time_seconds",
		Help:        "records the time that event occurs at unix time seconds",
		ConstLabels: labels,
	}, []string{})
	metrics.Register().MustRegister(m.lastEventTime)

	m.eventLagDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   "watch",
		Name:        "event_lag_seconds",
		Help:        "the lags(seconds) of the event between it occurs and we received it",
		ConstLabels: labels,
		Buckets:     []float64{0.02, 0.04, 0.06, 0.08, 0.1, 0.3, 0.5, 0.7, 1, 5, 10, 20, 30, 60, 120},
	}, []string{"action"})
	metrics.Register().MustRegister(m.eventLagDurations)

	m.lastEventLagDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   "watch",
		Name:        "last_event_lag_seconds",
		Help:        "record the last event's lag duration in seconds",
		ConstLabels: labels,
	}, []string{})
	metrics.Register().MustRegister(m.lastEventLagDuration)

	m.cycleDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   "watch",
		Name:        "event_cycle_seconds",
		Help:        "the total duration(seconds) of each event being handled",
		ConstLabels: labels,
		Buckets:     []float64{0.02, 0.04, 0.06, 0.08, 0.1, 0.3, 0.5, 0.7, 1, 5, 10, 15, 20, 40, 60, 120},
	}, []string{})
	metrics.Register().MustRegister(m.cycleDurations)

	m.totalErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: "watch",
		Name:      "total_error_count",
		Help: "the total error event count which we handled with this resources, " +
			"including invalid event and re-watch operations",
		ConstLabels: labels,
	}, []string{"error_type"})
	metrics.Register().MustRegister(m.totalErrorCount)

	return m
}

type eventMetrics struct {
	// record the total event count be handled.
	totalEventCount *prometheus.CounterVec

	// last event time records when the current last event occurs and when we received with unix time seconds.
	// we can use this metric to know which time point we have handled the event for now.
	lastEventTime *prometheus.GaugeVec

	// record the all event lag which is the duration time difference
	// between the event is occur and the time we received the with watch.
	// unit is seconds.
	eventLagDurations *prometheus.HistogramVec

	// record the last event's lag duration
	lastEventLagDuration *prometheus.GaugeVec

	// record the cost time of every cycle of handling a event chain.
	// unit is seconds
	cycleDurations *prometheus.HistogramVec

	// record the total errors when we watch the event.
	// which contains the number of as follows:
	// 1. watch failed count
	// 2. invalid event count
	totalErrorCount *prometheus.CounterVec
}

// collectBasic collect the basic event's metrics
func (em *eventMetrics) collectBasic(e *types.Event) {
	// increase event's total count with operate type
	em.totalEventCount.With(prometheus.Labels{"action": string(e.OperationType)}).Inc()

	// the time when the event is really happens.
	at := time.Unix(int64(e.ClusterTime.Sec), int64(e.ClusterTime.Nano))

	// unix time in seconds
	em.lastEventTime.With(prometheus.Labels{}).Set(float64(at.Unix()))

	// calculate event lags, in seconds
	lags := time.Since(at).Seconds()

	// set last lag duration
	em.lastEventLagDuration.With(prometheus.Labels{}).Set(lags)

	// add to lags durations
	em.eventLagDurations.With(prometheus.Labels{"action": string(e.OperationType)}).Observe(lags)
}

// the total duration(seconds) of each event being handled
func (em *eventMetrics) collectCycleDuration(d time.Duration) {
	em.cycleDurations.With(prometheus.Labels{}).Observe(d.Seconds())
}

// collect retry operation for any reason
func (em *eventMetrics) collectRetryError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "retry"}).Inc()
}

// redis lock related error
func (em *eventMetrics) collectLockError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "lock"}).Inc()
}

// collect redis operation related errors
func (em *eventMetrics) collectRedisError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "redis_command"}).Inc()
}

// collect mongodb related errors, such as get info from table cc_DelArchive
func (em *eventMetrics) collectMongoError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "mongo_command"}).Inc()
}
