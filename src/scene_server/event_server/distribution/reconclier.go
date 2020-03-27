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
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal"

	"github.com/tidwall/gjson"
	"gopkg.in/redis.v5"
)

type reconciler struct {
	db                   dal.RDB
	cache                *redis.Client
	cached               map[string][]string
	persisted            map[string][]string
	cachedSubscribers    []string
	persistedSubscribers []string
	processID            string
	ctx                  context.Context
}

func newReconciler(ctx context.Context, cache *redis.Client, db dal.RDB) *reconciler {
	return &reconciler{
		ctx:                  ctx,
		db:                   db,
		cache:                cache,
		cached:               map[string][]string{},
		persisted:            map[string][]string{},
		persistedSubscribers: []string{},
	}
}

var MsgChan = make(chan string, 3)

func (r *reconciler) loadAll() {
	r.cached = map[string][]string{}
	r.persisted = map[string][]string{}
	r.persistedSubscribers = []string{}
	r.loadAllCached()
	r.loadAllPersisted()
}

func (r *reconciler) loadAllCached() {
	r.cached = map[string][]string{}
	for _, formKey := range r.cache.Keys(types.EventCacheSubscribeFormKey + "*").Val() {
		if formKey != "" && formKey != nilStr && formKey != "redis" {
			r.cached[strings.TrimPrefix(formKey, types.EventCacheSubscribeFormKey)] = r.cache.SMembers(formKey).Val()
		}
	}
}

func (r *reconciler) loadAllPersisted() {
	r.persisted = map[string][]string{}
	r.persistedSubscribers = []string{}
	subscriptions := make([]metadata.Subscription, 0)
	if err := r.db.Table(common.BKTableNameSubscription).Find(nil).All(r.ctx, &subscriptions); err != nil {
		blog.Errorf("reconcile err: %v", err)
	}
	blog.Infof("loaded %v subscriptions from persistent", len(subscriptions))
	for _, sub := range subscriptions {
		eventNames := strings.Split(sub.SubscriptionForm, ",")
		r.persistedSubscribers = append(r.persistedSubscribers, sub.GetCacheKey())
		for _, eventName := range eventNames {
			eventName = sub.OwnerID + ":" + eventName
			r.persisted[eventName] = append(r.persisted[eventName], fmt.Sprint(sub.SubscriptionID))
		}
	}
}

func (r *reconciler) reconcile() {

	for k, v := range r.persisted {
		subs, plugs := util.CalSliceDiff(r.cached[k], v)
		if len(subs) > 0 {
			subss, _ := util.GetMapInterfaceByInerface(subs)
			if err := r.cache.SRem(types.EventCacheSubscribeFormKey+k, subss...).Err(); err != nil {
				blog.Errorf("reconcile err: %v", err)
			}
		}
		if len(plugs) > 0 {
			plugss, _ := util.GetMapInterfaceByInerface(plugs)
			if err := r.cache.SAdd(types.EventCacheSubscribeFormKey+k, plugss...).Err(); err != nil {
				blog.Errorf("reconcile err: %v", err)
			}
		}
		delete(r.cached, k)
	}

	for k := range r.cached {
		r.cache.Del(types.EventCacheSubscribeFormKey + k)
	}

}

func SubscribeChannel(redisCli *redis.Client) (err error) {
	subChan, err := redisCli.PSubscribe(types.EventCacheProcessChannel)
	if err != nil {
		return err
	}
	blog.Info("start receiving massages from redis")
	for {
		msgIf, err := subChan.Receive()
		if err == redis.Nil || err == io.EOF {
			continue
		}
		if nil != err {
			blog.Warnf("SubscribeChannel err %s, continue", err.Error())
			if err := subChan.Unsubscribe(types.EventCacheProcessChannel); err != nil {
				blog.Errorf("Unsubscribe channel %s failed, err: %+v", types.EventCacheProcessChannel, err)
			}
			time.Sleep(time.Second)
			if err := subChan.Subscribe(types.EventCacheProcessChannel); err != nil {
				blog.Errorf("Subscribe channel %s failed, err: %+v", types.EventCacheProcessChannel, err)
			}
			continue
		}
		msg, ok := msgIf.(*redis.Message)
		if !ok {
			blog.Warnf("SubscribeChannel receive a message of unexpect type: %v, msg: %+v, continue", reflect.TypeOf(msgIf).String(), msgIf)
			continue
		}
		if "" == msg.Payload {
			blog.Warnf("SubscribeChannel ignore empty Payload empty")
			continue
		}
		MsgChan <- msg.Payload
	}
}

func cleanExpiredEvents(redisCli *redis.Client) {
	var err error
	timeout := time.Hour * 1
	tick := util.NewTicker(timeout)
	tick.Tick()
	for range tick.C {
		blog.Infof("starting clean expired events")
		var keys = make([]string, 0)
		if err = redisCli.Keys(types.EventCacheDistDonePrefix + "*").ScanSlice(&keys); err != nil {
			blog.Errorf("fetch expired event keys failed: %v", err)
		}
		keys = append(keys, types.EventCacheEventDoneKey)

		for _, key := range keys {
			iter := redisCli.HScan(key, 0, "*", 10).Iterator()
			for iter.Next() {
				if strings.HasPrefix(iter.Val(), "{") {
					if time.Now().Sub(gjson.Get(iter.Val(), "action_time").Time()) > timeout {
						if err = redisCli.HDel(key, gjson.Get(iter.Val(), "event_id").String()).Err(); err != nil {
							blog.Errorf("remove expired event %s failed: %v", iter.Val(), err)
						}
					}
				}
			}
			if err := iter.Err(); err != nil {
				blog.Errorf("scan expired events failed: %v", err)
			}
		}
	}
}
