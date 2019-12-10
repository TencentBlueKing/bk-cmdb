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
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/event_server/types"
)

func (dh *DistHandler) StartDistribute() (err error) {
	defer func() {
		sysError := recover()
		if sysError != nil {
			err = fmt.Errorf("system error: %v", sysError)
		}
		if err != nil {
			blog.Errorf("distribute process stop with error: %v, stack:%s", err, debug.Stack())
		}
	}()

	blog.Info("distribution handle process started")

	reconciler := newReconciler(dh.ctx, dh.cache, dh.db)
	reconciler.loadAll()
	reconciler.reconcile()
	subscribers := reconciler.persistedSubscribers

	chErr := make(chan error, 1)
	routines := map[int64]chan struct{}{}
	renewMaps := map[int64]chan metadata.Subscription{}
	for _, str := range subscribers {
		subscriber := metadata.Subscription{}
		if err := json.Unmarshal([]byte(str), &subscriber); err != nil {
			return err
		}

		done := make(chan struct{})
		renewCh := make(chan metadata.Subscription)
		go func() {
			err := dh.distToSubscribe(subscriber, renewCh, done)
			if err != nil {
				chErr <- err
			}
		}()
		renewMaps[subscriber.SubscriptionID] = renewCh
		routines[subscriber.SubscriptionID] = done
	}

	go func() {
		for range time.Tick(time.Second * 60) {
			reconciler.loadAll()
			reconciler.reconcile()
			for _, sub := range reconciler.persistedSubscribers {
				MsgChan <- "update" + sub
			}
		}
	}()

	go func() {
		blog.Infof("discovering subscriber change")

		defer blog.Warn("discovering subscriber change process stopped")
		for {
			msg := <-MsgChan
			msgAction := extractChangeAction(msg)
			msgBody := extractChangeBody(msg)

			subscriber := metadata.Subscription{}
			blog.Infof("msg: action:%s, body:%s", msgAction, msgBody)
			if err := json.Unmarshal([]byte(msgBody), &subscriber); err != nil {
				chErr <- err
				return
			}
			switch msgAction {
			case "create":
				blog.Infof("starting subscribers process %d", subscriber.SubscriptionID)
				if renewCh, ok := renewMaps[subscriber.SubscriptionID]; ok {
					renewCh <- subscriber
					continue
				}
				done := make(chan struct{})
				renewCh := make(chan metadata.Subscription)
				go func() {
					err := dh.distToSubscribe(subscriber, renewCh, done)
					if err != nil {
						chErr <- err
					}
				}()
				routines[subscriber.SubscriptionID] = done
				renewMaps[subscriber.SubscriptionID] = renewCh
			case "update":
				blog.Infof("renew subscribers process %d", subscriber.SubscriptionID)
				if renewCh, exist := renewMaps[subscriber.SubscriptionID]; exist == true {
					renewCh <- subscriber
				} else {
					MsgChan <- "create" + msgBody
				}
			case "delete":
				blog.Infof("subscriber has been deleted, now stopping subscribe process, subscriberID: %d", subscriber.SubscriptionID)
				if routines[subscriber.SubscriptionID] != nil {
					close(routines[subscriber.SubscriptionID])
					delete(routines, subscriber.SubscriptionID)
					delete(renewMaps, subscriber.SubscriptionID)
				}
			}
		}
	}()

	return <-chErr

}

func (dh *DistHandler) distToSubscribe(param metadata.Subscription, chNew chan metadata.Subscription, done chan struct{}) (err error) {
	blog.Infof("start handle dist %v", param.SubscriptionID)
	defer func() {
		sysError := recover()
		if sysError != nil {
			err = fmt.Errorf("system error: %v", sysError)
		}
		if err != nil {
			blog.Infof("event inst handle process stopped by %v: %s", err, debug.Stack())
		}
	}()
	sub := param
	ticker := time.NewTicker(time.Minute)
	defer blog.Infof("ended handle dist %v", sub.SubscriptionID)
	for {
		select {
		case nsub := <-chNew:
			if nsub.GetCacheKey() != sub.GetCacheKey() {
				sub = nsub
				blog.Infof("refreshed subscriber %v", sub.GetCacheKey())
			} else {
				blog.Infof("refresh ignore, subscriber cache key not change\nold:%s\nnew:%s ", sub.GetCacheKey(), nsub.GetCacheKey())
			}
		case <-ticker.C:
			filter := map[string]interface{}{
				common.BKSubscriptionIDField: sub.SubscriptionID,
			}
			count, countErr := dh.db.Table(common.BKTableNameSubscription).Find(filter).Count(context.Background())
			if countErr != nil {
				blog.Errorf("get subscription count error %v", countErr)
				continue
			}
			if count <= 0 {
				ticker.Stop()
				return
			}
		case <-done:
			return
		default:
			dist := dh.popDistInst(sub.SubscriptionID)
			if dist == nil {
				continue
			}
			if err = dh.handleDist(&sub, dist); err != nil {
				blog.Errorf("error handle dist: %v, %v", err, dist)
			}
		}
	}
}

func (dh *DistHandler) handleDist(sub *metadata.Subscription, dist *metadata.DistInstCtx) (err error) {
	blog.Infof("handling dist %s", dist.Raw)
	distID := fmt.Sprint(dist.DstbID - 1)
	subscriberID := fmt.Sprint(dist.SubscriptionID)
	runningKey := types.EventCacheDistRunningPrefix + subscriberID + "_" + distID
	if err = saveRunning(dh.cache, runningKey, timeout+sub.GetTimeout()); err != nil {
		if ErrProcessExists == err {
			blog.Infof("process exist, continue")
			return nil
		}
		return err
	}

	previousID := fmt.Sprint(dist.DstbID - 1)
	previousRunningKey := types.EventCacheDistRunningPrefix + subscriberID + "_" + previousID
	done, err := checkFromDone(dh.cache, types.EventCacheDistDonePrefix+subscriberID, previousID)
	if err != nil {
		return err
	}
	if !done {

		running, checkErr := checkFromRunning(dh.cache, previousRunningKey)
		if checkErr != nil {
			return checkErr
		}
		if !running {

			time.Sleep(time.Second * 5)
			running, checkErr = checkFromRunning(dh.cache, previousRunningKey)
			if checkErr != nil {
				return checkErr
			}
		}
		if running {

			blog.Infof("waiting previous id: " + previousID)
			if checkErr = waitPreviousDone(dh.cache, types.EventCacheDistDonePrefix+subscriberID, previousID, sub.GetTimeout()); checkErr != nil && checkErr != ErrWaitTimeout {
				return checkErr
			}
			if checkErr == ErrWaitTimeout {
				blog.Infof("wait timeout previous id: %v, begin send callback", previousID)
			}
		}
	}

	defer func() {
		if err = dh.saveDistDone(dist); err != nil {
			return
		}
		blog.Infof("done event dist : %v", dist.DstbID)
	}()

	if err = dh.SendCallback(sub, dist.Raw); err != nil {
		blog.Errorf("send callback error: %v", err)
		return
	}

	return
}

func (dh *DistHandler) popDistInst(subID int64) *metadata.DistInstCtx {
	eventSlice := dh.cache.BLPop(time.Second*10, types.EventCacheDistQueuePrefix+fmt.Sprint(subID)).Val()

	if len(eventSlice) <= 0 {
		return nil
	}

	eventBytes := []byte(eventSlice[1])
	event := metadata.DistInst{}
	if err := json.Unmarshal(eventBytes, &event); err != nil {
		blog.Errorf("event distribute fail, unmarshal error: %v, data=[%s]", err, eventBytes)
		return nil
	}

	return &metadata.DistInstCtx{DistInst: event, Raw: eventSlice[1]}
}

func (dh *DistHandler) saveDistDone(dist *metadata.DistInstCtx) (err error) {
	if err = dh.cache.HSet(types.EventCacheDistDonePrefix+fmt.Sprint(dist.SubscriptionID), fmt.Sprint(dist.DstbID), dist.Raw).Err(); err != nil {
		return
	}
	if err = dh.cache.Del(types.EventCacheDistRunningPrefix + fmt.Sprintf("%d_%d", dist.SubscriptionID, dist.DstbID)).Err(); err != nil {
		return
	}
	return
}

func extractChangeAction(msg string) string {
	return msg[:6]
}
func extractChangeBody(msg string) string {
	return msg[6:]
}
