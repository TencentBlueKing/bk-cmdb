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
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

const (
	// defaultSubTimeoutSeconds is default seconds num for new subscription.
	defaultSubTimeoutSeconds = 10
)

// Subscribe subscribes target resource event in callback mode.
func (s *Service) Subscribe(ctx *rest.Contexts) {
	// decode request data.
	sub := &metadata.Subscription{}
	if err := ctx.DecodeInto(&sub); err != nil {
		blog.Errorf("add new subscription decode request body failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		// 400, unmarshal failed.
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	if len(sub.SubscriptionName) == 0 {
		// 400, empty subscription name.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionName"))
		return
	}
	if len(sub.CallbackURL) == 0 {
		// 400, empty callback url.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "CallbackURL"))
		return
	}
	if len(sub.SubscriptionForm) == 0 {
		// 400, empty subscription form.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionForm"))
		return
	}

	// operator.
	sub.Operator = ctx.Kit.User

	// subscription timeout seconds.
	if sub.TimeOutSeconds <= 0 {
		sub.TimeOutSeconds = defaultSubTimeoutSeconds
	}

	// subscription confirm mode.
	if sub.ConfirmMode != metadata.ConfirmModeHTTPStatus && sub.ConfirmMode != metadata.ConfirmModeRegular {
		// 400, unknown confirm mode.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "ConfirmMode"))
		return
	}
	if sub.ConfirmMode == metadata.ConfirmModeHTTPStatus && len(sub.ConfirmPattern) == 0 {
		sub.ConfirmPattern = strconv.FormatInt(http.StatusOK, 10)
	}

	sub.LastTime = metadata.Now()
	sub.OwnerID = ctx.Kit.SupplierAccount

	// trim subscription form.
	sub.SubscriptionForm = s.trimSubscriptionForm(sub.SubscriptionForm)

	res, err := s.engine.CoreAPI.CoreService().Event().Subscribe(ctx.Kit.Ctx, ctx.Kit.Header, sub)
	if err != nil {
		blog.Errorf("Subscribe failed, Subscribe err:%s, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// register cloud sync task resource creator action to iam
	if auth.EnableAuthorize() {
		iamInstance := metadata.IamInstanceWithCreator{
			Type:    string(iam.SysEventPushing),
			ID:      strconv.FormatInt(res.SubscriptionID, 10),
			Name:    res.SubscriptionName,
			Creator: res.Operator,
		}
		_, err := s.authorizer.RegisterResourceCreatorAction(s.ctx, ctx.Kit.Header, iamInstance)
		if err != nil {
			blog.Errorf("register created event subscription to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	data := struct {
		SubscriptionID int64 `json:"subscription_id"`
	}{SubscriptionID: res.SubscriptionID}

	ctx.RespEntity(data)
}

// UnSubscribe unsubscribes target resource event in callback mode.
func (s *Service) UnSubscribe(ctx *rest.Contexts) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("subscribeID"), 10, 64)
	if err != nil {
		// 400, invalid subscribeID parameter.
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	if err = s.engine.CoreAPI.CoreService().Event().UnSubscribe(ctx.Kit.Ctx, ctx.Kit.Header, id); err != nil {
		blog.Errorf("delete target subscription by id[%d] failed, err: %+v, rid: %s", id, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// UpdateSubscription updates target subscription in callback mode.
func (s *Service) UpdateSubscription(ctx *rest.Contexts) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("subscribeID"), 10, 64)
	if err != nil {
		// 400, invalid subscribeID parameter.
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	// decode request data.
	sub := &metadata.Subscription{}
	if err := ctx.DecodeInto(&sub); err != nil {
		blog.Errorf("update target subscription decode request body failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		// 400, unmarshal failed.
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	if len(sub.SubscriptionName) == 0 {
		// 400, empty subscription name.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionName"))
		return
	}
	if len(sub.CallbackURL) == 0 {
		// 400, empty callback url.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "CallbackURL"))
		return
	}
	if len(sub.SubscriptionForm) == 0 {
		// 400, empty subscription form.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionForm"))
		return
	}

	// subscription confirm mode.
	if sub.ConfirmMode != metadata.ConfirmModeHTTPStatus && sub.ConfirmMode != metadata.ConfirmModeRegular {
		// 400, unknown confirm mode.
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "ConfirmMode"))
		return
	}
	if sub.ConfirmMode == metadata.ConfirmModeHTTPStatus && len(sub.ConfirmPattern) == 0 {
		sub.ConfirmPattern = strconv.FormatInt(http.StatusOK, 10)
	}
	sub.Operator = ctx.Kit.User

	// trim subscription form.
	sub.SubscriptionForm = s.trimSubscriptionForm(sub.SubscriptionForm)

	// update subscription.
	if err = s.engine.CoreAPI.CoreService().Event().UpdateSubscription(ctx.Kit.Ctx, ctx.Kit.Header, id, sub); err != nil {
		// 400, update target subscription failed.
		blog.Errorf("update target subscription by condition failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrEventSubscribeUpdateFailed))
		return
	}

	ctx.RespEntity(nil)
}

// ListSubscriptions lists all subscriptions in cc.
func (s *Service) ListSubscriptions(ctx *rest.Contexts) {
	// decode request data.
	data := metadata.ParamSubscriptionSearch{}
	if err := ctx.DecodeInto(&data); err != nil {
		// 400, unmarshal failed.
		blog.Errorf("list subscriptions decode request body failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	data.Condition = util.SetModOwner(data.Condition, ctx.Kit.SupplierAccount)

	// get authorized event subscription ids if auth is enabled
	if auth.EnableAuthorize() {
		authInput := meta.ListAuthorizedResourcesParam{
			UserName:     ctx.Kit.User,
			ResourceType: meta.EventPushing,
			Action:       meta.Find,
		}

		authorizedResources, err := s.authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
		if err != nil {
			blog.ErrorJSON("list authorized subscribe resources failed, err: %v, cond: %s, rid: %s", err, authInput, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		subscriptions := make([]int64, 0)
		for _, resourceID := range authorizedResources {
			subscriptionID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse resourceID(%s) failed, err: %v, rid: %s", resourceID, err, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}
			subscriptions = append(subscriptions, subscriptionID)
		}

		data.Condition = map[string]interface{}{
			common.BKDBAND: []map[string]interface{}{
				data.Condition,
				{
					common.BKSubscriptionIDField: map[string]interface{}{
						common.BKDBIN: subscriptions,
					},
				},
			},
		}
	}

	res, err := s.engine.CoreAPI.CoreService().Event().ListSubscriptions(ctx.Kit.Ctx, ctx.Kit.Header, &data)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

func (s *Service) Ping(ctx *rest.Contexts) {
	var data metadata.ParamSubscriptionTestCallback
	if err := ctx.DecodeInto(&data); err != nil {
		blog.Errorf("ping subscription failed, decode request body failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	callbackUrl := data.CallbackUrl
	callbackBody := data.Data

	blog.Infof("requesting callback url: %s, data: %s, rid: %s", callbackUrl, callbackBody, ctx.Kit.Rid)
	callbackReq, _ := http.NewRequest(http.MethodPost, callbackUrl, bytes.NewBufferString(callbackBody))
	callbackResp, err := http.DefaultClient.Do(callbackReq)
	if err != nil {
		blog.Errorf("test distribute failed, do http request failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrEventSubscribePingFailed))
		return
	}
	defer callbackResp.Body.Close()

	callbackRespBody, err := ioutil.ReadAll(callbackResp.Body)
	if err != nil {
		blog.Errorf("test distribute failed, read response body failed, err:%v, rid: %s", err, ctx.Kit.Rid)
	}
	result := metadata.RspSubscriptionTestCallback{}
	result.HttpStatus = callbackResp.StatusCode
	result.ResponseBody = string(callbackRespBody)

	ctx.RespEntity(result)
}

func (s *Service) Telnet(ctx *rest.Contexts) {
	var data metadata.ParamSubscriptionTelnet
	if err := ctx.DecodeInto(&data); nil != err {
		blog.Errorf("telnet subscription failed, decode request body failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	callbackUrl := data.CallbackUrl
	uri, err := util.GetDailAddress(callbackUrl)
	if err != nil {
		blog.Errorf("telnet callback failed, err:%+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_callback_url"))
		return
	}
	blog.Infof("telnet url: %, rid: %s", uri, ctx.Kit.Rid)

	conn, err := net.Dial("tcp", uri)
	if err != nil {
		blog.Errorf("telnet callback failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrEventSubscribeTelnetFailed))
		return
	}
	conn.Close()

	ctx.RespEntity(nil)
}

// trimSubscriptionForm trims space on subscription form.
func (s *Service) trimSubscriptionForm(subscriptionForm string) string {
	subscriptionFormStr := strings.Replace(subscriptionForm, " ", "", -1)
	subscriptionForms := strings.Split(subscriptionFormStr, ",")

	sort.Strings(subscriptionForms)
	return strings.Join(subscriptionForms, ",")
}
