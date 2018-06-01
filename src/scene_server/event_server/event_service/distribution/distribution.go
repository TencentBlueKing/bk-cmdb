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
	"configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/event_server/types"
	"encoding/json"
	"fmt"
	redis "gopkg.in/redis.v5"
	"runtime/debug"
	"time"
)

func StartDistribute() (err error) {
	defer func() {
		if err == nil {
			syserror := recover()
			if syserror != nil {
				err = fmt.Errorf("system error: %v", syserror)
			}
		}
		debug.PrintStack()
	}()
	// reconcil cache from persistent store
	rccler := newReconciler()
	rccler.loadAll()
	rccler.reconcile()
	subscribers := rccler.persistedSubscribers

	// start distribution goroutines
	chErr := make(chan error)
	routines := map[int64]chan struct{}{}
	renewMaps := map[int64]chan types.Subscription{}
	for _, str := range subscribers {
		subscriber := types.Subscription{}
		if err := json.Unmarshal([]byte(str), &subscriber); err != nil {
			return err
		}

		done := make(chan struct{})
		renewCh := make(chan types.Subscription)
		go func() {
			err := distToSubscribe(subscriber, renewCh, done)
			if err != nil {
				chErr <- err
			}
		}()
		renewMaps[subscriber.SubscriptionID] = renewCh
		routines[subscriber.SubscriptionID] = done
	}

	// discovering subscriber change in order to stop or start distribution goroutine
	go func() {
		blog.Infof("discovering subscriber change")

		defer blog.Warn("discovering subscriber change process stoped")
		for {
			mesg := <-MsgChan
			action := mesg[:6]
			subscriber := types.Subscription{}
			blog.Infof("mesg: action:%s ,body:%s", mesg[:6], mesg[6:])
			if err := json.Unmarshal([]byte(mesg[6:]), &subscriber); err != nil {
				chErr <- err
				return
			}
			switch action {
			case "create":
				blog.Infof("starting subscribers process %d", subscriber.SubscriptionID)
				done := make(chan struct{})
				renewCh := make(chan types.Subscription)
				go func() {
					err := distToSubscribe(subscriber, renewCh, done)
					if err != nil {
						chErr <- err
					}
				}()
				routines[subscriber.SubscriptionID] = done
				renewMaps[subscriber.SubscriptionID] = renewCh
			case "update":
				blog.Infof("renew subscribers process %d", subscriber.SubscriptionID)
				renewMaps[subscriber.SubscriptionID] <- subscriber
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

func distToSubscribe(param types.Subscription, chNew chan types.Subscription, done chan struct{}) (err error) {
	blog.Infof("start handle dist %v", param.SubscriptionID)
	defer func() {
		if err == nil {
			syserror := recover()
			if syserror != nil {
				err = fmt.Errorf("system error: %v", syserror)
			}
		}
		if err != nil {
			blog.Info("event inst handle process stoped by %v", err)
			debug.PrintStack()
		}
	}()
	sub := param
	go func() {
		for {
			sub = <-chNew
			blog.Infof("refreshed subcriber %v", sub.SubscriptionID)
		}
	}()
	defer blog.Infof("ended handle dist %v", sub.SubscriptionID)
	for {
		select {
		case sub = <-chNew:
		case <-done:
			return
		default:
			dist := popDistInst(sub.SubscriptionID)
			if dist == nil {
				continue
			}
			if err = handleDist(&sub, dist); err != nil {
				blog.Errorf("error handle dist: %v, %v", err, dist)
			}
		}
	}
}

func handleDist(sub *types.Subscription, dist *types.DistInstCtx) (err error) {
	blog.Infof("handling dist %s", dist.Raw)
	distID := fmt.Sprint(dist.DstbID - 1)
	subscriberID := fmt.Sprint(dist.SubscriptionID)
	runningkey := types.EventCacheDistRunningPrefix + subscriberID + "_" + distID
	if err = saveRunning(runningkey, timeout+sub.GetTimeout()); err != nil {
		if ERR_PROCESS_EXISTS == err {
			blog.Infof("process exist, continue")
			return nil
		}
		return err
	}

	// check previous done
	priviousID := fmt.Sprint(dist.DstbID - 1)
	priviousRunningkey := types.EventCacheDistRunningPrefix + subscriberID + "_" + priviousID
	done, err := checkFromDone(types.EventCacheDistDonePrefix+subscriberID, priviousID)
	if err != nil {
		return err
	}
	if !done {
		// check previous running
		running, checkErr := checkFromRunning(priviousRunningkey)
		if checkErr != nil {
			return checkErr
		}
		if !running {
			// wait a second to ensure previous not in trouble
			time.Sleep(time.Second * 5)
			running, checkErr = checkFromRunning(priviousRunningkey)
			if checkErr != nil {
				return checkErr
			}
		}
		if running {
			// if previous running, wait it
			blog.Infof("waitting previous id: " + priviousID)
			if checkErr = waitPreviousDone(types.EventCacheDistDonePrefix+subscriberID, priviousID, sub.GetTimeout()); checkErr != nil && checkErr != ERR_WAIT_TIMEOUT {
				return checkErr
			}
			if checkErr == ERR_WAIT_TIMEOUT {
				blog.Infof("wait timeout previous id: %v, begin send callback", priviousID)
			}
		}
	}

	defer func() {
		if err = saveDistDone(dist); err != nil {
			return
		}
		blog.Info("done event dist : %v", dist.DstbID)
	}()
	// if previous done then begin send callback
	if err = SendCallback(sub, dist.Raw); err != nil {
		blog.Errorf("send callback error: %v", err)
		return
	}

	return
}

func popDistInst(subID int64) *types.DistInstCtx {
	eventseletor := common.KvMap{
		"expire": time.Second * 60,
		"key":    []string{types.EventCacheDistQueuePrefix + fmt.Sprint(subID)},
	}
	eventslice := []string{}
	api.GetAPIResource().CacheCli.GetOneByCondition("blpop", nil, eventseletor, &eventslice)

	if len(eventslice) <= 0 {
		return nil
	}

	// Unmarshal event
	eventbytes := []byte(eventslice[1])
	event := types.DistInst{}
	if err := json.Unmarshal(eventbytes, &event); err != nil {
		blog.Errorf("event distribute fail, unmarshal error: %v, date=[%s]", err, eventbytes)
		return nil
	}

	return &types.DistInstCtx{DistInst: event, Raw: eventslice[1]}
}

func saveDistDone(dist *types.DistInstCtx) (err error) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	if err = redisCli.HSet(types.EventCacheDistDonePrefix+fmt.Sprint(dist.SubscriptionID), fmt.Sprint(dist.DstbID), dist.Raw).Err(); err != nil {
		return
	}
	if err = redisCli.Del(types.EventCacheDistRunningPrefix + fmt.Sprintf("%d_%d", dist.SubscriptionID, dist.DstbID)).Err(); err != nil {
		return
	}
	return
}
