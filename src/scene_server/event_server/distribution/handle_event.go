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
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/event_server/types"

	"gopkg.in/redis.v5"
)

var (
	timeout    = time.Second * 10
	waitPeriod = time.Second
)

// Err define
var (
	ErrWaitTimeout   = fmt.Errorf("wait timeout")
	ErrProcessExists = fmt.Errorf("process exists")
)

func (eh *EventHandler) Run() (err error) {
	defer func() {
		sysError := recover()
		if sysError != nil {
			err = fmt.Errorf("system error: %v", sysError)
		}
		if err != nil {
			blog.Infof("event inst handle process stopped by %v", err)
			blog.Errorf("%s", debug.Stack())
		}
	}()

	blog.Info("event inst handle process started")
	for {
		event := eh.popEvent()
		if event == nil {
			time.Sleep(time.Second * 2)
			continue
		}
		if err := eh.handleEvent(event); err != nil {
			blog.Errorf("handle event failed, err: %+v, event: %+v", err, event)
		}
	}
}

func (eh *EventHandler) handleEvent(event *metadata.EventInstCtx) (err error) {
	blog.Infof("handling event inst : %v", event.Raw)
	defer blog.Infof("done event inst : %d", event.ID)

	// check and set running status on, if event is being deal, then simple ignore this event
	if err = saveRunning(eh.cache, types.EventCacheEventRunningPrefix+fmt.Sprint(event.ID), timeout); err != nil {
		if ErrProcessExists == err {
			blog.Infof("%v process exist, continue", event.ID)
			return nil
		}
		blog.Infof("save runtime error: %v, raw = %s", err, event.Raw)
		return err
	}

	// check and wait previous event finished handling
	previousID := strconv.FormatInt(event.ID-1, 10)
	done, err := checkFromDone(eh.cache, types.EventCacheEventDoneKey, previousID)
	if err != nil {
		return err
	}
	if !done {
		previousRunningKey := types.EventCacheEventRunningPrefix + previousID
		running, checkErr := checkFromRunning(eh.cache, previousRunningKey)
		if checkErr != nil {
			return checkErr
		}
		if !running {
			time.Sleep(time.Second * 3)
			running, checkErr = checkFromRunning(eh.cache, previousRunningKey)
			if checkErr != nil {
				return checkErr
			}
		}
		if running {
			if checkErr = waitPreviousDone(eh.cache, types.EventCacheEventDoneKey, previousID, timeout); checkErr != nil {
				if checkErr == ErrWaitTimeout {
					return nil
				}
				return checkErr
			}
		}
	}

	defer func() {
		if err != nil {
			blog.Errorf("prepare dist event error:%v", err)
		}
		err = eh.SaveEventDone(event)
	}()

	originDists := eh.GetDistInst(&event.EventInst)

	for _, originDist := range originDists {
		subscribers := eh.findEventTypeSubscribers(originDist.GetType(), event.OwnerID)
		if len(subscribers) <= 0 || nilStr == subscribers[0] {
			blog.Infof("%v no subscriber，continue", originDist.GetType())
			return eh.SaveEventDone(event)
		}

		for _, subscriber := range subscribers {
			var dstbID, subscribeID int64
			distInst := originDist
			dstbID, err = eh.nextDistID(subscriber)
			if err != nil {
				return err
			}
			subscribeID, err = strconv.ParseInt(subscriber, 10, 64)
			if err != nil {
				return err
			}
			distInst.DstbID = dstbID
			distInst.SubscriptionID = subscribeID
			distByte, _ := json.Marshal(distInst)
			eh.pushToQueue(types.EventCacheDistQueuePrefix+subscriber, string(distByte))
		}
	}

	return
}

func (eh *EventHandler) GetDistInst(e *metadata.EventInst) []metadata.DistInst {
	distInst := metadata.DistInst{
		EventInst: *e,
	}
	distInst.ID = 0
	var ds []metadata.DistInst
	if e.EventType == metadata.EventTypeInstData && e.ObjType == common.BKInnerObjIDObject {
		if len(e.Data) <= 0 {
			return nil
		}

		var ok bool
		var m map[string]interface{}
		if e.Action == metadata.EventActionDelete {
			m, ok = e.Data[0].PreData.(map[string]interface{})
		} else {
			m, ok = e.Data[0].CurData.(map[string]interface{})
		}
		if !ok {
			// TODO: just ignore this event without any warnings?
			return nil
		}

		// TODO: panic if key not exist?
		if m[common.BKObjIDField] != nil {
			distInst.ObjType = m[common.BKObjIDField].(string)
		}

	}
	ds = append(ds, distInst)

	return ds
}

func (eh *EventHandler) pushToQueue(key, value string) (err error) {
	err = eh.cache.RPush(key, value).Err()
	blog.Infof("pushed to queue:%v", key)
	return
}

func (eh *EventHandler) nextDistID(eventType string) (nextid int64, err error) {
	var id int64
	id, err = eh.cache.Incr(types.EventCacheDistIDPrefix + eventType).Result()
	if err != nil {
		return
	}
	return id, nil
}

func (eh *EventHandler) SaveEventDone(event *metadata.EventInstCtx) (err error) {
	if err = eh.cache.HSet(types.EventCacheEventDoneKey, fmt.Sprint(event.ID), event.Raw).Err(); err != nil {
		return
	}
	if err = eh.cache.Del(types.EventCacheEventRunningPrefix + fmt.Sprint(event.ID)).Err(); err != nil {
		return
	}
	return
}

func waitPreviousDone(cache *redis.Client, key string, id string, timeout time.Duration) (err error) {
	var done bool
	timer := time.NewTimer(timeout)
	for !done {
		select {
		case <-timer.C:
			timer.Stop()
			return ErrWaitTimeout
		default:
			done, err = checkFromDone(cache, key, id)
			if err != nil {
				return
			}
		}
		time.Sleep(waitPeriod)
	}
	return
}

func checkFromDone(cache *redis.Client, key string, id string) (bool, error) {
	if id == "0" {
		return true, nil
	}
	return cache.HExists(key, fmt.Sprint(id)).Result()
}

func checkFromRunning(cache *redis.Client, key string) (bool, error) {
	return cache.Exists(key).Result()
}

func saveRunning(cache *redis.Client, key string, timeout time.Duration) (err error) {
	set, err := cache.SetNX(key, time.Now().UTC().Format(time.RFC3339), timeout).Result()
	if !set {
		return ErrProcessExists
	}
	return err
}

func (eh *EventHandler) findEventTypeSubscribers(eventType, ownerID string) []string {
	return eh.cache.SMembers(types.EventSubscriberCacheKey(ownerID, eventType)).Val()
}

func (eh *EventHandler) popEvent() *metadata.EventInstCtx {
	var eventStr string

	// push event into types.EventCacheEventQueueDuplicateKey queue so that identifierHandler could deal with it
	eh.cache.BRPopLPush(types.EventCacheEventQueueKey, types.EventCacheEventQueueDuplicateKey, time.Second*60).Scan(&eventStr)

	if eventStr == "" || eventStr == nilStr {
		return nil
	}
	eventBytes := []byte(eventStr)
	event := metadata.EventInst{}
	if err := json.Unmarshal(eventBytes, &event); err != nil {
		blog.Errorf("event distribute fail, unmarshal error: %v, data=[%s]", err, eventBytes)
		return nil
	}
	return &metadata.EventInstCtx{EventInst: event, Raw: eventStr}
}

const nilStr = "nil"
