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

package service

import (
	"strconv"
	
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// Subscribe subscribes target resource event in callback mode.
func (s *coreService) Subscribe(ctx *rest.Contexts) {
	subscription := metadata.Subscription{}
	if err := ctx.DecodeInto(&subscription); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.EventOperation().Subscribe(ctx.Kit, &subscription)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// UnSubscribe unsubscribes target resource event in callback mode.
func (s *coreService) UnSubscribe(ctx *rest.Contexts) {
	//get subscribeID
	subscribeIDStr := ctx.Request.PathParameter(common.BKSubscribeID)
	subscribeID, err := strconv.ParseInt(subscribeIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSubscribeID))
		return
	}

	err = s.core.EventOperation().UnSubscribe(ctx.Kit, subscribeID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateSubscription updates target subscription in callback mode.
func (s *coreService) UpdateSubscription(ctx *rest.Contexts) {
	//get subscribeID
	subscribeIDStr := ctx.Request.PathParameter(common.BKSubscribeID)
	subscribeID, err := strconv.ParseInt(subscribeIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSubscribeID))
		return
	}

	sub := &metadata.Subscription{}
	if err := ctx.DecodeInto(&sub); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err = s.core.EventOperation().UpdateSubscription(ctx.Kit, subscribeID, sub)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// ListSubscriptions lists all subscriptions in cc.
func (s *coreService) ListSubscriptions(ctx *rest.Contexts) {
	option := metadata.ParamSubscriptionSearch{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.EventOperation().ListSubscriptions(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}
