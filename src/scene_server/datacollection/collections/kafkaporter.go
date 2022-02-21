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

package collections

import (
	"context"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

const (
	msgCount = 100
)

type consumerGroupHandler struct {
	// name of this porter.
	name string

	// analyzer analyzes message from collectors.
	analyzer Analyzer

	// collectors message channels kafka topics.
	topics []string

	// metrics.
	// receiveTotal is message received total stat.
	receiveTotal prometheus.Counter

	// receiveInvalidTotal is message received but it's not valid data total stat.
	receiveInvalidTotal prometheus.Counter

	// analyzeTotal is message analyzed total stat with target labels.
	analyzeTotal *prometheus.CounterVec

	// analyzeDuration is analyze cost duration stat.
	analyzeDuration prometheus.Histogram

	// registry is prometheus register.
	registry prometheus.Registerer
}

// init inits a new consumerGroupHandler.
func (c *consumerGroupHandler) init() {
	// register metrics.
	c.registerMetrics()
	blog.Infof("KafkaPorter[%s]| init metrics success!", c.name)
}

// registerMetrics registers prometheus metrics.
func (c *consumerGroupHandler) registerMetrics() {
	c.receiveTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_receive_total", metricsNamespacePrefix, c.name),
			Help: "total number of received message.",
		},
	)
	c.registry.MustRegister(c.receiveTotal)

	c.receiveInvalidTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_receive_invalid_total", metricsNamespacePrefix, c.name),
			Help: "total number of received invalid message.",
		},
	)
	c.registry.MustRegister(c.receiveInvalidTotal)

	c.analyzeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_analyze_total", metricsNamespacePrefix, c.name),
			Help: "total number of analyzed message.",
		},
		[]string{"status"},
	)
	c.registry.MustRegister(c.analyzeTotal)

	c.analyzeDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_%s_analyze_duration", metricsNamespacePrefix, c.name),
			Help: "analyze duration of each batch message.",
		},
	)
	c.registry.MustRegister(c.analyzeDuration)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim consumption logic
func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	blog.Infof("KafkaPorter[%s]| start a new analyze loop now!", c.name)
	var msgArray []string
	retry := false
	// record the last message for commit
	var lastMessage *sarama.ConsumerMessage
	for {
		if !retry {
			msgArray, lastMessage = c.readMessage(claim)
		}

		retry = false

		// open goroutine to analyse message until consumption is complete
		var wg sync.WaitGroup
		for _, msg := range msgArray {
			if retry {
				break
			}

			wg.Add(1)
			go func(msg string) {
				defer wg.Done()
				result, err := c.analyzer.Analyze(&msg)
				if err != nil {
					blog.Errorf("KafkaPorter[%s]| analyze message failed, %v", c.name, err)
					// metrics stats for analyze failed.
					c.analyzeTotal.WithLabelValues("failed").Inc()
				} else {
					// metrics stats for analyze success.
					c.analyzeTotal.WithLabelValues("success").Inc()
				}

				if result {
					retry = true
				}
			}(msg)
		}
		wg.Wait()

		if retry {
			continue
		}

		// metrics stats for analyze duration.
		c.analyzeDuration.Observe(time.Since(time.Now()).Seconds())

		// commit offset to kafka
		sess.MarkMessage(lastMessage, "")
		sess.Commit()
	}
	return nil
}

func (c *consumerGroupHandler) readMessage(claim sarama.ConsumerGroupClaim) ([]string, *sarama.ConsumerMessage) {
	// record the last message for commit
	var lastMessage *sarama.ConsumerMessage
	// the key is unique id for message, the value is message timestamp
	timeMap := make(map[string]int64)
	// the key is unique id for message, the value is message
	msgMap := make(map[string]string)

	startTime := time.Now()
	// 如果 没有消息 或者 在2秒内读取的消息没超过100条，那么这两种情况会进行消息的读取操作
	for len(msgMap) == 0 || (len(msgMap) <= msgCount && time.Now().Sub(startTime) < 2*time.Second) {
		var message *sarama.ConsumerMessage
		select {
		case message = <-claim.Messages():
			break
		case <-time.After(time.Second):
			break
		}

		if message == nil {
			continue
		}

		// metrics stats for message receiving.
		c.receiveTotal.Inc()

		if message.Value == nil {
			blog.Errorf("KafkaPorter[%s]| receive a message with empty value!", c.name)

			// metrics stats for invalid message.
			c.receiveInvalidTotal.Inc()
			continue
		}

		val := gjson.Parse(string(message.Value))
		cloudID := val.Get("cloudid").String()
		ip := val.Get("ip").String()
		uniqueKey := cloudID + ":" + ip
		timestamp := val.Get("data.timestamp").Int()
		oldTimestamp, exist := timeMap[uniqueKey]
		if exist && oldTimestamp > timestamp {
			continue
		}

		msgMap[uniqueKey] = string(message.Value)
		timeMap[uniqueKey] = timestamp

		lastMessage = message
	}
	var msgArray []string
	for _, msg := range msgMap {
		msgArray = append(msgArray, msg)
	}
	return msgArray, lastMessage
}

// KafkaPorter is a porter to handle message from kafka.
type KafkaPorter struct {
	engine *backbone.Engine

	// name of this porter.
	name string

	// collectors message channels kafka topics.
	topics []string

	// analyzer analyzes message from collectors.
	analyzer Analyzer

	ctx context.Context

	// consumerGroupHandler consumer group handle
	consumerGroupHandler sarama.ConsumerGroupHandler

	// consumerGroup consumer group
	consumerGroup sarama.ConsumerGroup
}

// NewKafkaPorter new kafka porter
func NewKafkaPorter(name string, engine *backbone.Engine, ctx context.Context, analyzer Analyzer,
	consumerGroup sarama.ConsumerGroup, topics []string, registry prometheus.Registerer) *KafkaPorter {
	consumerGroupHandler := &consumerGroupHandler{
		name:     name,
		analyzer: analyzer,
		topics:   topics,
		registry: registry,
	}
	consumerGroupHandler.init()
	return &KafkaPorter{
		name:                 name,
		engine:               engine,
		analyzer:             analyzer,
		topics:               topics,
		ctx:                  ctx,
		consumerGroupHandler: consumerGroupHandler,
		consumerGroup:        consumerGroup,
	}
}

// Name porter name
func (k *KafkaPorter) Name() string {
	return k.name
}

// Run run porter
func (k *KafkaPorter) Run() error {
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := k.consumerGroup.Consume(k.ctx, k.topics, k.consumerGroupHandler)
		if err != nil {
			blog.Errorf("KafkaPorter[%s]| have a error, err: %v", k.name, err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if k.ctx.Err() != nil {
			return err
		}
	}
	return nil
}

// Mock mock analyzer
func (k *KafkaPorter) Mock() error {
	mock := k.analyzer.Mock()
	if _, err := k.analyzer.Analyze(&mock); err != nil {
		return fmt.Errorf("mock failed, %+v", err)
	}
	return nil
}
