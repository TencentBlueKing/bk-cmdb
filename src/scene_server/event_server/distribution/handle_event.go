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

	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/event_server/types"
)

var (
	timeout    = time.Second * 10
	waitperiod = time.Second
)

// Err define
var (
	ErrWaitTimeout   = fmt.Errorf("wait timeout")
	ErrProcessExists = fmt.Errorf("process exists")
)

func (eh *EventHandler) StartHandleInsts() (err error) {
	defer func() {
		if err == nil {
			syserror := recover()
			if syserror != nil {
				err = fmt.Errorf("system error: %v", syserror)
			}
		}
		if err != nil {
			blog.Info("event inst handle process stoped by %v", err)
			blog.Errorf("%s", debug.Stack())
		}
	}()

	blog.Info("event inst handle process started")
	for {

		event := eh.popEventInst()
		if event == nil {
			time.Sleep(time.Second * 2)
			continue
		}
		if err := eh.handleInst(event); err != nil {
			blog.Errorf("error handle dist: %v, %v", err, event)
		}
	}
}

func (eh *EventHandler) handleInst(event *metadata.EventInstCtx) (err error) {
	blog.Info("handling event inst : %v", event.Raw)
	defer blog.Info("done event inst : %v", event.ID)
	if err = saveRunning(eh.cache, types.EventCacheEventRunningPrefix+fmt.Sprint(event.ID), timeout); err != nil {
		if ErrProcessExists == err {
			blog.Infof("%v process exist, continue", event.ID)
			return nil
		}
		blog.Infof("save runtime error: %v, raw = %s", err, event.Raw)
		return err
	}

	previousID := fmt.Sprint(event.ID - 1)
	priviousRunningkey := types.EventCacheEventRunningPrefix + previousID
	done, err := checkFromDone(eh.cache, types.EventCacheEventDoneKey, previousID)
	if err != nil {
		return err
	}
	if !done {
		running, checkErr := checkFromRunning(eh.cache, priviousRunningkey)
		if checkErr != nil {
			return checkErr
		}
		if !running {

			time.Sleep(time.Second * 3)
			running, checkErr = checkFromRunning(eh.cache, priviousRunningkey)
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

	origindists := eh.GetDistInst(&event.EventInst)

	for _, origindist := range origindists {
		subscribers := eh.findEventTypeSubscribers(origindist.GetType(), event.OwnerID)
		if len(subscribers) <= 0 || "nil" == subscribers[0] {
			blog.Infof("%v no subscriber，continue", origindist.GetType())
			return eh.SaveEventDone(event)
		}

		for _, subscriber := range subscribers {
			var dstbID, subscribeID int64
			distinst := origindist
			dstbID, err = eh.nextDistID(subscriber)
			if err != nil {
				return err
			}
			subscribeID, err = strconv.ParseInt(subscriber, 10, 64)
			if err != nil {
				return err
			}
			distinst.DstbID = dstbID
			distinst.SubscriptionID = subscribeID
			distByte, _ := json.Marshal(distinst)
			eh.pushToQueue(types.EventCacheDistQueuePrefix+subscriber, string(distByte))
		}
	}

	return
}

func (eh *EventHandler) GetDistInst(e *metadata.EventInst) []metadata.DistInst {
	distinst := metadata.DistInst{
		EventInst: *e,
	}
	distinst.ID = 0
	var ds []metadata.DistInst
	var m map[string]interface{}
	if e.EventType == metadata.EventTypeInstData && e.ObjType == common.BKINnerObjIDObject {
		var ok bool

		if len(e.Data) <= 0 {
			return nil
		}
		if e.Action == metadata.EventActionDelete {
			m, ok = e.Data[0].PreData.(map[string]interface{})
		} else {
			m, ok = e.Data[0].CurData.(map[string]interface{})
		}

		if !ok {
			return nil
		}

		if m[common.BKObjIDField] != nil {
			distinst.ObjType = m[common.BKObjIDField].(string)
		}

	}
	ds = append(ds, distinst)

	return ds
}

func (eh *EventHandler) pushToQueue(key, value string) (err error) {
	err = eh.cache.RPush(key, value).Err()
	blog.Infof("pushed to queue:%v", key)
	return
}

func (eh *EventHandler) nextDistID(eventtype string) (nextid int64, err error) {
	var id int64
	id, err = eh.cache.Incr(types.EventCacheDistIDPrefix + eventtype).Result()
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
		time.Sleep(waitperiod)
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

func (eh *EventHandler) findEventTypeSubscribers(eventtype, ownerID string) []string {
	return eh.cache.SMembers(types.EventSubscriberCacheKey(ownerID, eventtype)).Val()
}

func (eh *EventHandler) popEventInst() *metadata.EventInstCtx {
	var eventstr string

	eh.cache.BRPopLPush(types.EventCacheEventQueueKey, types.EventCacheEventQueueDuplicateKey, time.Second*60).Scan(&eventstr)

	if eventstr == "" || eventstr == "nil" {
		return nil
	}
	eventbytes := []byte(eventstr)
	event := metadata.EventInst{}
	if err := json.Unmarshal(eventbytes, &event); err != nil {
		blog.Errorf("event distribute fail, unmarshal error: %v, date=[%s]", err, eventbytes)
		return nil
	}
	return &metadata.EventInstCtx{EventInst: event, Raw: eventstr}
}
