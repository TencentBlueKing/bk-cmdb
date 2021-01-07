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

package distribution

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal/redis"

	"github.com/prometheus/client_golang/prometheus"
)

var httpCli = httpclient.NewHttpClient()

const (
	// defaultSendTimeout is default timeout for send action.
	defaultSendTimeout = 10 * time.Second

	// defaultEventCacheSubscriberCursorExpire is default expire duration for subscriber cursor.
	defaultEventCacheSubscriberCursorExpire = 6 * time.Hour
)

// EventPusher sends target events to subscribers in callback mode.
type EventPusher struct {
	ctx    context.Context
	engine *backbone.Engine

	// subid is subscription id.
	subid int64

	// cache is cc redis client.
	cache redis.Client

	// distributer handles all events distribution.
	distributer *Distributor

	// metrics.
	// pusherHandleTotal is event pusher handle total stat.
	pusherHandleTotal *prometheus.CounterVec

	// pusherHandleDuration is event pusher cost duration stat.
	pusherHandleDuration *prometheus.HistogramVec
}

// NewEventPusher creates a new EventPusher object.
func NewEventPusher(ctx context.Context, engine *backbone.Engine, subid int64, cache redis.Client, distributer *Distributor,
	pusherHandleTotal *prometheus.CounterVec, pusherHandleDuration *prometheus.HistogramVec) *EventPusher {
	return &EventPusher{
		ctx:                  ctx,
		engine:               engine,
		subid:                subid,
		cache:                cache,
		distributer:          distributer,
		pusherHandleTotal:    pusherHandleTotal,
		pusherHandleDuration: pusherHandleDuration,
	}
}

// Handle add event dist inst to subscriber chan, and pusher would send message to
// subscriber base on callback.
func (s *EventPusher) Handle(dist *metadata.DistInst) error {
	if dist == nil {
		return errors.New("invalid event dist metadata")
	}

	distData, err := json.Marshal(dist)
	if err != nil {
		return err
	}

	subscriberEventQueueKey := types.EventCacheSubscriberEventQueueKeyPrefix + fmt.Sprint(dist.SubscriptionID)
	if err := s.cache.LPush(s.ctx, subscriberEventQueueKey, distData).Err(); err != nil {
		return err
	}
	return nil
}

func (s *EventPusher) increaseTotal(subid int64) error {
	return s.statIncrease(subid, "total")
}

func (s *EventPusher) increaseFailure(subid int64) error {
	return s.statIncrease(subid, "failue")
}

// statIncrease stats callback details by increase in cache.
func (s *EventPusher) statIncrease(subid int64, key string) error {
	return s.cache.HIncrBy(s.ctx, types.EventCacheDistCallBackCountPrefix+fmt.Sprint(subid), key, 1).Err()
}

// cleaning keeps cleaning expire or redundancy events in subscriber event cache queue.
func (s *EventPusher) cleaning() {
	ticker := time.NewTicker(defaultCleanCheckInterval)
	defer ticker.Stop()

	for {
		if !s.engine.ServiceManageInterface.IsMaster() {
			blog.Warnf("not master eventserver node, skip cleaning subscriber[%d] events!", s.subid)
			time.Sleep(defaultMasterCheckInterval)
			continue
		}

		<-ticker.C
		blog.Info("cleaning expire and redundancy subscriber events now")

		// clean.
		subscriberEventQueueKey := types.EventCacheSubscriberEventQueueKeyPrefix + fmt.Sprint(s.subid)
		eventCount, err := s.cache.LLen(s.ctx, subscriberEventQueueKey).Result()
		if err != nil {
			blog.Errorf("fetch expire and redundancy subscriber events failed, %+v", err)
			continue
		}
		blog.Info("cleaning subscriber events, current events count[%d], trim threshold[%d], delete threshold[%d], clean unit count[%d]",
			eventCount, defaultCleanTrimThreshold, defaultCleanDelThreshold, defaultCleanUnit)

		if eventCount < defaultCleanTrimThreshold {
			continue
		}

		if eventCount < defaultCleanDelThreshold {
			if err := s.cache.LTrim(s.ctx, subscriberEventQueueKey, 0, eventCount-defaultCleanUnit).Err(); err != nil {
				blog.Errorf("trim expire and redundancy subscriber events failed, %+v", err)
				continue
			}
		}

		// too many events.
		if err := s.cache.Del(s.ctx, subscriberEventQueueKey).Err(); err != nil {
			blog.Errorf("delete expire and redundancy subscriber queue failed, %+v", err)
			continue
		}
		blog.Info("cleaning expire and redundancy subscriber events done")
	}
}

// push sends new event to target subscriber base on callback url.
func (s *EventPusher) push(dist *metadata.DistInst) error {
	// try to find new subscription data everytime, and send event
	// with newest http callback url.
	subscription := s.distributer.FindSubscription(s.subid)
	if subscription == nil {
		return fmt.Errorf("subscription not found, %+v", s.subid)
	}

	// setups ownerid here.
	dist.OwnerID = subscription.OwnerID

	var errFinal error

	// stats.
	s.increaseTotal(subscription.SubscriptionID)
	defer func() {
		if errFinal != nil {
			s.increaseFailure(subscription.SubscriptionID)
		}
	}()

	// marshal message data.
	distData, err := json.Marshal(dist)
	if err != nil {
		errFinal = err
		return err
	}

	// build http request.
	body := bytes.NewBuffer(distData)
	req, err := http.NewRequest("POST", subscription.CallbackURL, body)
	if err != nil {
		errFinal = err
		return err
	}

	// callback timeout.
	var duration time.Duration
	if subscription.TimeOutSeconds == 0 {
		duration = defaultSendTimeout
	} else {
		duration = subscription.GetTimeout()
	}

	// send now.
	resp, err := httpCli.DoWithTimeout(duration, req)
	if err != nil {
		errFinal = err
		return err
	}
	defer resp.Body.Close()

	// read response.
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errFinal = err
		return err
	}

	// confirm mode.
	if subscription.ConfirmMode == metadata.ConfirmModeHTTPStatus {
		if strconv.Itoa(resp.StatusCode) != subscription.ConfirmPattern {
			errFinal = err
			return fmt.Errorf("not confirm http pattern, received %s", respData)
		}
	} else if subscription.ConfirmMode == metadata.ConfirmModeRegular {
		pattern, err := regexp.Compile(subscription.ConfirmPattern)
		if err != nil {
			errFinal = err
			return fmt.Errorf("build regexp error, %+v", err)
		}

		if !pattern.Match(respData) {
			errFinal = err
			return fmt.Errorf("not confirm regular pattern, received %s", respData)
		}
	} else {
		// do nothing, just let it go.
	}

	// mark resource type and action cursor.
	eventType := dist.EventInst.GetType()
	suberCursorKey := types.EventCacheSubscriberCursorKey(eventType, s.subid)

	if err := s.cache.Set(s.ctx, suberCursorKey, dist.Cursor, defaultEventCacheSubscriberCursorExpire).Err(); err != nil {
		blog.Warnf("save subscriber[%d] cursor for action[%s] failed, %+v", s.subid, eventType, err)
	}
	return nil
}

func (s *EventPusher) run() {
	// keep cleaning.
	go s.cleaning()

	for {
		if !s.engine.ServiceManageInterface.IsMaster() {
			blog.Warnf("not master eventserver node, skip push event for subscriber[%d]", s.subid)
			time.Sleep(defaultMasterCheckInterval)
			continue
		}

		// keep sending.
		cost := time.Now()
		distDatas := s.cache.BRPop(s.ctx, defaultTransTimeout, types.EventCacheSubscriberEventQueueKeyPrefix+fmt.Sprint(s.subid)).Val()
		s.pusherHandleDuration.WithLabelValues("PopSubscriberEvent").Observe(time.Since(cost).Seconds())

		// distDatas is redis brpop results, and you can parse it base on CMD
		// formats, https://redis.io/commands/brpop.
		if len(distDatas) == 0 || distDatas[1] == types.NilStr || len(distDatas[1]) == 0 {
			continue
		}
		distData := distDatas[1]

		dist := &metadata.DistInst{}
		if err := json.Unmarshal([]byte(distData), dist); err != nil {
			blog.Errorf("unmarshal new event dist inst for subscriber[%d] failed, %+v", s.subid, err)
			continue
		}

		if time.Now().Unix()-dist.EventInst.ActionTime.Unix() > defaultFusingEventExpireSec {
			// old event, expire it.
			s.pusherHandleTotal.WithLabelValues("ExpireEventNum").Inc()
			continue
		}

		// send message to subscriber.
		cost = time.Now()
		err := s.push(dist)
		s.pusherHandleDuration.WithLabelValues("SendSubscriberEvent").Observe(time.Since(cost).Seconds())

		if err != nil {
			s.pusherHandleTotal.WithLabelValues("SendCallbackFailed").Inc()
			blog.Errorf("send event failed, err: %+v, data=[%+v]", err, dist)
			continue
		}
		s.pusherHandleTotal.WithLabelValues("Success").Inc()
	}
}

// Run setups pusher and keep handling event dist.
func (s *EventPusher) Run() {
	// run pusher.
	go s.run()
}
