/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package event

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
)

const (
	// defaultSubTimeoutSeconds is default seconds num for new subscription.
	defaultSubTimeoutSeconds = 10
)

func (e *eventOperation) Subscribe(kit *rest.Kit, subscription *metadata.Subscription) (*metadata.Subscription, errors.CCErrorCoder) {
	// create new subscription now.
	existSubscriptions := []metadata.Subscription{}
	filter := map[string]interface{}{
		common.BKSubscriptionNameField: subscription.SubscriptionName,
		common.BKOwnerIDField:          kit.SupplierAccount,
	}
	if err := e.dbProxy.Table(common.BKTableNameSubscription).Find(filter).All(kit.Ctx, &existSubscriptions); err != nil {
		// 200, duplicated subscription name of target ownerid.
		// NOTE: maybe just internal system errors.
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKSubscriptionNameField)
	}

	if len(existSubscriptions) > 0 {
		// 200, duplicated subscription name of target ownerid.
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKSubscriptionNameField)
	}

	// generate instance id.
	subscriptionID, err := e.dbProxy.NextSequence(kit.Ctx, common.BKTableNameSubscription)
	if err != nil {
		// 500, failed to get sequence to insert a new subscription instance.
		blog.Errorf("Subscribe failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrEventSubscribeInsertFailed)

	}
	subscription.SubscriptionID = int64(subscriptionID)

	if err := e.dbProxy.Table(common.BKTableNameSubscription).Insert(kit.Ctx, subscription); err != nil {
		// 500, failed to insert a new subscription instance.
		blog.Errorf("create new subscription failed, err:%+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrEventSubscribeInsertFailed)
	}

	e.cache.Del(context.Background(), types.EventCacheDistCallBackCountPrefix+fmt.Sprint(subscription.SubscriptionID))

	return subscription, nil
}

func (e *eventOperation) UnSubscribe(kit *rest.Kit, subscribeID int64) errors.CCErrorCoder {
	// query target subscription info.
	sub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, subscribeID, common.BKOwnerIDField, kit.SupplierAccount).Build()

	if err := e.dbProxy.Table(common.BKTableNameSubscription).Find(condition).One(kit.Ctx, &sub); err != nil {
		// 500, query target subscription info failed.
		blog.Errorf("query target subscription by id[%d] failed, err: %+v, rid: %s", subscribeID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrEventSubscribeDeleteFailed)
	}

	// delete subscription.
	if err := e.dbProxy.Table(common.BKTableNameSubscription).Delete(kit.Ctx, condition); err != nil {
		// 500, delete target subscription failed.
		blog.Errorf("delete target subscription by id[%d] failed, err: %+v, rid: %s", subscribeID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrEventSubscribeDeleteFailed)
	}

	e.cache.Del(context.Background(), types.EventCacheDistIDPrefix+fmt.Sprint(sub.SubscriptionID),
		types.EventCacheSubscriberEventQueueKeyPrefix+fmt.Sprint(sub.SubscriptionID),
		types.EventCacheDistCallBackCountPrefix+fmt.Sprint(sub.SubscriptionID))

	return nil
}

func (e *eventOperation) UpdateSubscription(kit *rest.Kit, subscribeID int64, sub *metadata.Subscription) errors.CCErrorCoder {
	// query target subscription.
	oldSub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, subscribeID, common.BKOwnerIDField, kit.SupplierAccount).Build()

	if err := e.dbProxy.Table(common.BKTableNameSubscription).Find(condition).One(kit.Ctx, &oldSub); err != nil {
		blog.Errorf("query target subscription by id[%v] failed, err: %+v, rid: %s", subscribeID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrEventSubscribeUpdateFailed)
	}

	// check duplicated when subscription name changed.
	if oldSub.SubscriptionName != sub.SubscriptionName {
		filter := map[string]interface{}{
			common.BKSubscriptionNameField: sub.SubscriptionName,
			common.BKOwnerIDField:          kit.SupplierAccount,
		}

		count, err := e.dbProxy.Table(common.BKTableNameSubscription).Find(filter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("query subscription with the name count under target ownerid failed, err: %+v, rid: %s", err, kit.Rid)
			return kit.CCError.CCError(common.CCErrEventSubscribeUpdateFailed)
		}
		if count > 0 {
			blog.Errorf("can't update target subscription, the name is duplicated, rid: %s", kit.Rid)
			return kit.CCError.CCError(common.CCErrEventSubscribeUpdateFailed)
		}
	}

	// set subscriptionid and other fields.
	sub.SubscriptionID = oldSub.SubscriptionID
	if sub.TimeOutSeconds <= 0 {
		sub.TimeOutSeconds = defaultSubTimeoutSeconds
	}
	sub.LastTime = metadata.Now()
	sub.OwnerID = kit.SupplierAccount

	filter := map[string]interface{}{
		common.BKSubscriptionIDField: subscribeID,
		common.BKOwnerIDField:        kit.SupplierAccount,
	}
	if err := e.dbProxy.Table(common.BKTableNameSubscription).Update(kit.Ctx, filter, sub); err != nil {
		blog.Errorf("update target subscription by condition failed, err: %+v, rid: %s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrEventSubscribeUpdateFailed)
	}

	return nil
}

func (e *eventOperation) ListSubscriptions(kit *rest.Kit, data *metadata.ParamSubscriptionSearch) (*metadata.RspSubscriptionSearch, errors.CCErrorCoder) {
	if data.Page.Limit <= 0 {
		data.Page.Limit = common.BKNoLimit
	}
	util.SetQueryOwner(data.Condition, kit.SupplierAccount)
	count, err := e.dbProxy.Table(common.BKTableNameSubscription).Find(data.Condition).Count(kit.Ctx)
	if err != nil {
		// 400, query host count failed.
		blog.Errorf("query host count failed, input: %+v err: %+v, rid: %s", data, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrEventSubscribeSelectFailed)

	}
	results := []metadata.Subscription{}

	if selErr := e.dbProxy.Table(common.BKTableNameSubscription).Find(data.Condition).Fields(data.Fields...).Sort(data.Page.Sort).Start(uint64(data.Page.Start)).Limit(uint64(data.Page.Limit)).All(kit.Ctx, &results); nil != selErr {
		// 400, query source data failed.
		blog.Errorf("query resource data failed, err: %+v, input:%v, rid: %s", selErr, data, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrEventSubscribeSelectFailed)
	}

	for index := range results {
		val := e.cache.HGetAll(context.Background(), types.EventCacheDistCallBackCountPrefix+fmt.Sprint(results[index].SubscriptionID)).Val()
		failure, err := strconv.ParseInt(val["failue"], 10, 64)
		if nil != err {
			blog.Warnf("get failure value error %s, rid: %s", err.Error(), kit.Rid)
		}

		total, err := strconv.ParseInt(val["total"], 10, 64)
		if nil != err {
			blog.Warnf("get total value error %s, rid: %s", err.Error(), kit.Rid)
		}

		results[index].Statistics = &metadata.Statistics{
			Total:   total,
			Failure: failure,
		}
	}

	info := make(map[string]interface{})
	info["count"] = count
	info["info"] = results

	return &metadata.RspSubscriptionSearch{Count: count, Info: results}, nil
}
