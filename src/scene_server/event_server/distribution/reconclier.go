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
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage"
)

type reconciler struct {
	db                   storage.DI
	cache                *redis.Client
	cached               map[string][]string
	persisted            map[string][]string
	cachedSubscribers    []string
	persistedSubscribers []string
	processID            string
}

func newReconciler(cache *redis.Client, db storage.DI) *reconciler {
	return &reconciler{
		db:                   db,
		cache:                cache,
		cached:               map[string][]string{},
		persisted:            map[string][]string{},
		persistedSubscribers: []string{},
	}
}

var MsgChan = make(chan string, 3)

func (r *reconciler) loadAll() {
	r.loadAllCached()
	r.loadAllPersisted()
}

func (r *reconciler) loadAllCached() {
	r.cached = map[string][]string{}
	for _, formkey := range r.cache.Keys(types.EventCacheSubscribeformKey + "*").Val() {
		if formkey != "" && formkey != "nil" && formkey != "redis" {
			r.cached[strings.TrimPrefix(formkey, types.EventCacheSubscribeformKey)] = r.cache.SMembers(formkey).Val()
		}
	}
}

func (r *reconciler) loadAllPersisted() {
	r.persisted = map[string][]string{}
	r.persistedSubscribers = []string{}
	subscriptions := []metadata.Subscription{}
	if err := r.db.GetMutilByCondition(common.BKTableNameSubscription, nil, nil, &subscriptions, "", 0, 0); err != nil {
		blog.Errorf("reconcile err: %v", err)
	}
	blog.Infof("loaded %v subscriptions from persistent", len(subscriptions))
	for _, sub := range subscriptions {
		eventnames := strings.Split(sub.SubscriptionForm, ",")
		r.persistedSubscribers = append(r.persistedSubscribers, sub.GetCacheKey())
		for _, eventname := range eventnames {
			eventname = sub.OwnerID + ":" + eventname
			r.persisted[eventname] = append(r.persisted[eventname], fmt.Sprint(sub.SubscriptionID))
		}
	}
}

func (r *reconciler) reconcile() {

	for k, v := range r.persisted {
		subs, plugs := util.CalSliceDiff(r.cached[k], v)
		if len(subs) > 0 {
			subss, _ := util.GetMapInterfaceByInerface(subs)
			if err := r.cache.SRem(types.EventCacheSubscribeformKey+k, subss...).Err(); err != nil {
				blog.Errorf("reconcile err: %v", err)
			}
		}
		if len(plugs) > 0 {
			plugss, _ := util.GetMapInterfaceByInerface(plugs)
			if err := r.cache.SAdd(types.EventCacheSubscribeformKey+k, plugss...).Err(); err != nil {
				blog.Errorf("reconcile err: %v", err)
			}
		}
		delete(r.cached, k)
	}

	for k := range r.cached {
		r.cache.Del(types.EventCacheSubscribeformKey + k)
	}

}

func SubscribeChannel(redisCli *redis.Client) (err error) {
	subChan, err := redisCli.PSubscribe(types.EventCacheProcessChannel)
	if err != nil {
		return err
	}
	blog.Info("receiving massages")
	for {
		mesg, err := subChan.Receive()
		if err == redis.Nil || err == io.EOF {
			continue
		}
		if nil != err {
			blog.Warnf("SubscribeChannel err %s,, continue", err.Error())
			subChan.Unsubscribe(types.EventCacheProcessChannel)
			time.Sleep(time.Second)
			subChan.Subscribe(types.EventCacheProcessChannel)
			continue
		}
		msg, ok := mesg.(*redis.Message)
		if !ok {
			blog.Warnf("SubscribeChannel msg not message type but %v, continue", reflect.TypeOf(mesg).String())
			continue
		}
		if "" == msg.Payload {
			blog.Warnf("SubscribeChannel Payload empty, continue")
			continue
		}
		MsgChan <- msg.Payload
	}
}
