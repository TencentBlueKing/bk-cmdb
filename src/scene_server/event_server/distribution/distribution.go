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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/event_server/types"
)

func (dh *DistHandler) StartDistribute() (err error) {
	defer func() {
		if err == nil {
			syserror := recover()
			if syserror != nil {
				err = fmt.Errorf("system error: %v", syserror)
			}
		}
		blog.Errorf("%s", debug.Stack())
	}()

	rccler := newReconciler(dh.cache, dh.db)
	rccler.loadAll()
	rccler.reconcile()
	subscribers := rccler.persistedSubscribers

	chErr := make(chan error)
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
			rccler.loadAll()
			rccler.reconcile()
			for _, sub := range rccler.persistedSubscribers {
				MsgChan <- "update" + sub
			}
		}
	}()

	go func() {
		blog.Infof("discovering subscriber change")

		defer blog.Warn("discovering subscriber change process stoped")
		for {
			mesg := <-MsgChan
			mesgAction := getChangeAction(mesg)
			mesgBody := getChangeBody(mesg)

			subscriber := metadata.Subscription{}
			blog.Infof("mesg: action:%s ,body:%s", mesgAction, mesgBody)
			if err := json.Unmarshal([]byte(mesgBody), &subscriber); err != nil {
				chErr <- err
				return
			}
			switch mesgAction {
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
				if renewCh, ok := renewMaps[subscriber.SubscriptionID]; ok {
					renewCh <- subscriber
				} else {
					MsgChan <- "create" + mesgBody
				}
			case "delete":
				blog.Infof("stoping subscribers process %d", subscriber.SubscriptionID)
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
		syserror := recover()
		if syserror != nil {
			err = fmt.Errorf("system error: %v", syserror)
		}
		if err != nil {
			blog.Infof("event inst handle process stoped by %v: %s", err, debug.Stack())
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
				blog.Infof("refreshed subcriber %v", sub.GetCacheKey())
			} else {
				blog.Infof("refresh ignore, subcriber cache key not change\nold:%s\nnew:%s ", sub.GetCacheKey(), nsub.GetCacheKey())
			}
		case <-ticker.C:
			count, counterr := dh.db.GetCntByCondition(common.BKTableNameSubscription, condition.CreateCondition().Field(common.BKSubscriptionIDField).Eq(sub.SubscriptionID).ToMapStr())
			if counterr != nil {
				blog.Errorf("get subscription count error %v", counterr)
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
	runningkey := types.EventCacheDistRunningPrefix + subscriberID + "_" + distID
	if err = saveRunning(dh.cache, runningkey, timeout+sub.GetTimeout()); err != nil {
		if ErrProcessExists == err {
			blog.Infof("process exist, continue")
			return nil
		}
		return err
	}

	priviousID := fmt.Sprint(dist.DstbID - 1)
	priviousRunningkey := types.EventCacheDistRunningPrefix + subscriberID + "_" + priviousID
	done, err := checkFromDone(dh.cache, types.EventCacheDistDonePrefix+subscriberID, priviousID)
	if err != nil {
		return err
	}
	if !done {

		running, checkErr := checkFromRunning(dh.cache, priviousRunningkey)
		if checkErr != nil {
			return checkErr
		}
		if !running {

			time.Sleep(time.Second * 5)
			running, checkErr = checkFromRunning(dh.cache, priviousRunningkey)
			if checkErr != nil {
				return checkErr
			}
		}
		if running {

			blog.Infof("waitting previous id: " + priviousID)
			if checkErr = waitPreviousDone(dh.cache, types.EventCacheDistDonePrefix+subscriberID, priviousID, sub.GetTimeout()); checkErr != nil && checkErr != ErrWaitTimeout {
				return checkErr
			}
			if checkErr == ErrWaitTimeout {
				blog.Infof("wait timeout previous id: %v, begin send callback", priviousID)
			}
		}
	}

	defer func() {
		if err = dh.saveDistDone(dist); err != nil {
			return
		}
		blog.Info("done event dist : %v", dist.DstbID)
	}()

	if err = dh.SendCallback(sub, dist.Raw); err != nil {
		blog.Errorf("send callback error: %v", err)
		return
	}

	return
}

func (dh *DistHandler) popDistInst(subID int64) *metadata.DistInstCtx {
	eventslice := dh.cache.BLPop(time.Second*10, types.EventCacheDistQueuePrefix+fmt.Sprint(subID)).Val()

	if len(eventslice) <= 0 {
		return nil
	}

	eventbytes := []byte(eventslice[1])
	event := metadata.DistInst{}
	if err := json.Unmarshal(eventbytes, &event); err != nil {
		blog.Errorf("event distribute fail, unmarshal error: %v, date=[%s]", err, eventbytes)
		return nil
	}

	return &metadata.DistInstCtx{DistInst: event, Raw: eventslice[1]}
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

func getChangeAction(mesg string) string {
	return mesg[:6]
}
func getChangeBody(mesg string) string {
	return mesg[6:]
}
