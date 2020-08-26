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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"

	"github.com/emicklei/go-restful"
)

const (
	// defaultSubTimeoutSeconds is default seconds num for new subscription.
	defaultSubTimeoutSeconds = 10
)

// Subscribe subscribes target resource event in callback mode.
func (s *Service) Subscribe(req *restful.Request, resp *restful.Response) {
	// base request metadatas.
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)

	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	// decode request data.
	sub := &metadata.Subscription{}
	if err := json.NewDecoder(req.Request.Body).Decode(&sub); err != nil {
		blog.Errorf("add new subscription decode request body failed, err: %+v, rid: %s", err, rid)

		// 400, unmarshal failed.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(sub.SubscriptionName) == 0 {
		// 400, empty subscription name.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionName")})
		return
	}
	if len(sub.CallbackURL) == 0 {
		// 400, empty callback url.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "CallbackURL")})
		return
	}
	if len(sub.SubscriptionForm) == 0 {
		// 400, empty subscription form.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionForm")})
		return
	}

	// operator.
	sub.Operator = util.GetUser(req.Request.Header)

	// subscription timeout seconds.
	if sub.TimeOutSeconds <= 0 {
		sub.TimeOutSeconds = defaultSubTimeoutSeconds
	}

	// subscription confirm mode.
	if sub.ConfirmMode != metadata.ConfirmModeHTTPStatus && sub.ConfirmMode != metadata.ConfirmModeRegular {
		// 400, unknown confirm mode.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "ConfirmMode")})
		return
	}
	if sub.ConfirmMode == metadata.ConfirmModeHTTPStatus && len(sub.ConfirmPattern) == 0 {
		sub.ConfirmPattern = strconv.FormatInt(http.StatusOK, 10)
	}

	sub.LastTime = metadata.Now()
	sub.OwnerID = ownerID

	// trim subscription form.
	sub.SubscriptionForm = s.trimSubscriptionForm(sub.SubscriptionForm)

	// create new subscription now.
	existSubscriptions := []metadata.Subscription{}
	filter := map[string]interface{}{
		common.BKSubscriptionNameField: sub.SubscriptionName,
		common.BKOwnerIDField:          ownerID,
	}
	if err := s.db.Table(common.BKTableNameSubscription).Find(filter).All(s.ctx, &existSubscriptions); err != nil {
		// 200, duplicated subscription name of target ownerid.
		// NOTE: maybe just internal system errors.
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommDuplicateItem, common.BKSubscriptionNameField)})
		return
	}

	if len(existSubscriptions) > 0 {
		// 200, duplicated subscription name of target ownerid.
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommDuplicateItem, common.BKSubscriptionNameField)})
		return
	}

	// generate instance id.
	subscriptionID, err := s.db.NextSequence(s.ctx, common.BKTableNameSubscription)
	if err != nil {
		// 500, failed to get sequence to insert a new subscription instance.
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed)})
		return
	}
	sub.SubscriptionID = int64(subscriptionID)

	if err := s.db.Table(common.BKTableNameSubscription).Insert(s.ctx, sub); err != nil {
		// 500, failed to insert a new subscription instance.
		blog.Errorf("create new subscription failed, err:%+v, rid: %s", err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed)})
		return
	}
	s.cache.Del(types.EventCacheDistCallBackCountPrefix + fmt.Sprint(sub.SubscriptionID))

	// register subscription to iam
	iamResource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Name:       sub.SubscriptionName,
			Type:       meta.EventPushing,
			InstanceID: sub.SubscriptionID,
		},
	}
	if err = s.auth.RegisterResource(s.ctx, iamResource); err != nil {
		blog.Errorf("register subscribe to iam failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommRegistResourceToIAMFailed, err)}
		resp.WriteError(http.StatusOK, result)
		return
	}

	result := &metadata.RspSubscriptionCreate{
		BaseResp: metadata.SuccessBaseResp,
		Data: struct {
			SubscriptionID int64 `json:"subscription_id"`
		}{SubscriptionID: int64(subscriptionID)},
	}
	resp.WriteEntity(result)
}

// UnSubscribe unsubscribes target resource event in callback mode.
func (s *Service) UnSubscribe(req *restful.Request, resp *restful.Response) {
	// base request metadatas.
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)

	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	id, err := strconv.ParseInt(req.PathParameter("subscribeID"), 10, 64)
	if err != nil {
		// 400, invalid subscribeID parameter.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// query target subscription info.
	sub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()

	if err := s.db.Table(common.BKTableNameSubscription).Find(condition).One(s.ctx, &sub); err != nil {
		// 500, query target subscription info failed.
		blog.Errorf("query target subscription by id[%d] failed, err: %+v, rid: %s", id, err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeDeleteFailed)})
		return
	}

	// delete subscription.
	if err := s.db.Table(common.BKTableNameSubscription).Delete(s.ctx, condition); err != nil {
		// 500, delete target subscription failed.
		blog.Errorf("delete target subscription by id[%d] failed, err: %+v, rid: %s", id, err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeDeleteFailed)})
		return
	}
	s.cache.Del(types.EventCacheDistIDPrefix+fmt.Sprint(sub.SubscriptionID),
		types.EventCacheSubscriberEventQueueKeyPrefix+fmt.Sprint(sub.SubscriptionID),
		types.EventCacheDistCallBackCountPrefix+fmt.Sprint(sub.SubscriptionID))

	// deregister subscription from iam
	iamResource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Name:       sub.SubscriptionName,
			Type:       meta.EventPushing,
			InstanceID: sub.SubscriptionID,
		},
	}
	if err = s.auth.DeregisterResource(s.ctx, iamResource); err != nil {
		blog.Errorf("deregister subscribe to iam failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommUnRegistResourceToIAMFailed, err)}
		resp.WriteError(http.StatusOK, result)
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// UpdateSubscription updates target subscription in callback mode.
func (s *Service) UpdateSubscription(req *restful.Request, resp *restful.Response) {
	// base request metadatas.
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)

	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	id, err := strconv.ParseInt(req.PathParameter("subscribeID"), 10, 64)
	if err != nil {
		// 400, invalid subscribeID parameter.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// decode request data.
	sub := &metadata.Subscription{}
	if err := json.NewDecoder(req.Request.Body).Decode(&sub); err != nil {
		blog.Errorf("update target subscription decode request body failed, err: %+v, rid: %s", err, rid)

		// 400, unmarshal failed.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(sub.SubscriptionName) == 0 {
		// 400, empty subscription name.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionName")})
		return
	}
	if len(sub.CallbackURL) == 0 {
		// 400, empty callback url.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "CallbackURL")})
		return
	}
	if len(sub.SubscriptionForm) == 0 {
		// 400, empty subscription form.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionForm")})
		return
	}

	// subscription confirm mode.
	if sub.ConfirmMode != metadata.ConfirmModeHTTPStatus && sub.ConfirmMode != metadata.ConfirmModeRegular {
		// 400, unknown confirm mode.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "ConfirmMode")})
		return
	}
	if sub.ConfirmMode == metadata.ConfirmModeHTTPStatus && len(sub.ConfirmPattern) == 0 {
		sub.ConfirmPattern = strconv.FormatInt(http.StatusOK, 10)
	}
	sub.Operator = util.GetUser(req.Request.Header)

	// update subscription.
	if err = s.updateSubscription(header, id, ownerID, sub); err != nil {
		// 400, update target subscription failed.
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeUpdateFailed)})
		return
	}

	// deregister subscription from iam
	iamResource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Name:       sub.SubscriptionName,
			Type:       meta.EventPushing,
			InstanceID: sub.SubscriptionID,
		},
	}
	if err := s.auth.UpdateResource(s.ctx, &iamResource); err != nil {
		blog.Errorf("update subscribe to iam failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommRegistResourceToIAMFailed, err)}
		resp.WriteError(http.StatusOK, result)
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// updateSubscription compares infos and update target subscription.
func (s *Service) updateSubscription(header http.Header, id int64, ownerID string, sub *metadata.Subscription) error {
	rid := util.GetHTTPCCRequestID(header)

	// query target subscription.
	oldSub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()

	if err := s.db.Table(common.BKTableNameSubscription).Find(condition).One(s.ctx, &oldSub); err != nil {
		blog.Errorf("query target subscription by id[%v] failed, err: %+v, rid: %s", id, err, rid)
		return err
	}

	// check duplicated when subscription name changed.
	if oldSub.SubscriptionName != sub.SubscriptionName {
		filter := map[string]interface{}{
			common.BKSubscriptionNameField: sub.SubscriptionName,
			common.BKOwnerIDField:          ownerID,
		}

		count, err := s.db.Table(common.BKTableNameSubscription).Find(filter).Count(s.ctx)
		if err != nil {
			blog.Errorf("query subscription with the name count under target ownerid failed, err: %+v, rid: %s", err, rid)
			return err
		}
		if count > 0 {
			blog.Errorf("can't update target subscription, the name is duplicated, rid: %s", rid)
			return errors.New("duplicate subscriptions with target name")
		}
	}

	// set subscriptionid and other fields.
	sub.SubscriptionID = oldSub.SubscriptionID
	if sub.TimeOutSeconds <= 0 {
		sub.TimeOutSeconds = defaultSubTimeoutSeconds
	}
	sub.LastTime = metadata.Now()
	sub.OwnerID = ownerID

	// trim subscription form.
	sub.SubscriptionForm = s.trimSubscriptionForm(sub.SubscriptionForm)

	filter := map[string]interface{}{
		common.BKSubscriptionIDField: id,
		common.BKOwnerIDField:        ownerID,
	}
	if err := s.db.Table(common.BKTableNameSubscription).Update(s.ctx, filter, sub); err != nil {
		blog.Errorf("update target subscription by condition failed, err: %+v, rid: %s", err, rid)
		return err
	}

	return nil
}

// ListSubscriptions lists all subscriptions in cc.
func (s *Service) ListSubscriptions(req *restful.Request, resp *restful.Response) {
	// base request metadatas.
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)

	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	// decode request data.
	data := metadata.ParamSubscriptionSearch{}
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		// 400, unmarshal failed.
		blog.Errorf("list subscriptions decode request body failed, err: %+v, rid: %s", err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	fields := data.Fields
	condition := data.Condition
	condition = util.SetModOwner(condition, ownerID)

	skip := data.Page.Start
	limit := data.Page.Limit
	if limit <= 0 {
		limit = common.BKNoLimit
	}
	sortOption := data.Page.Sort

	count, err := s.db.Table(common.BKTableNameSubscription).Find(condition).Count(s.ctx)
	if err != nil {
		// 400, query host count failed.
		blog.Errorf("query host count failed, input: %+v err: %+v, rid: %s", data, err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeSelectFailed)})
		return
	}
	results := []metadata.Subscription{}

	if selErr := s.db.Table(common.BKTableNameSubscription).Find(condition).Fields(fields...).Sort(sortOption).Start(uint64(skip)).Limit(uint64(limit)).All(s.ctx, &results); nil != selErr {
		// 400, query source data failed.
		blog.Errorf("query resource data failed, err: %+v, input:%v, rid: %s", selErr, data, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeSelectFailed)})
		return
	}

	for index := range results {
		val := s.cache.HGetAll(types.EventCacheDistCallBackCountPrefix + fmt.Sprint(results[index].SubscriptionID)).Val()

		failure, err := strconv.ParseInt(val["failue"], 10, 64)
		if nil != err {
			blog.Warnf("get failure value error %s, rid: %s", err.Error(), rid)
		}

		total, err := strconv.ParseInt(val["total"], 10, 64)
		if nil != err {
			blog.Warnf("get total value error %s, rid: %s", err.Error(), rid)
		}

		results[index].Statistics = &metadata.Statistics{
			Total:   total,
			Failure: failure,
		}
	}

	info := make(map[string]interface{})
	info["count"] = count
	info["info"] = results

	result := metadata.RspSubscriptionSearch{
		Count: count,
		Info:  results,
	}

	resp.WriteEntity(metadata.NewSuccessResp(result))
}

func (s *Service) Ping(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	var data metadata.ParamSubscriptionTestCallback
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("ping subscription failed, decode request body failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	callbackUrl := data.CallbackUrl
	callbackBody := data.Data

	blog.Infof("requesting callback url: %s, data: %s, rid: %s", callbackUrl, callbackBody, rid)
	callbackReq, _ := http.NewRequest(http.MethodPost, callbackUrl, bytes.NewBufferString(callbackBody))
	callbackResp, err := http.DefaultClient.Do(callbackReq)
	if err != nil {
		blog.Errorf("test distribute failed, do http request failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribePingFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	defer callbackResp.Body.Close()

	callbackRespBody, err := ioutil.ReadAll(callbackResp.Body)
	if err != nil {
		blog.Errorf("test distribute failed, read response body failed, err:%v, rid: %s", err, rid)
	}
	result := metadata.RspSubscriptionTestCallback{}
	result.HttpStatus = callbackResp.StatusCode
	result.ResponseBody = string(callbackRespBody)

	resp.WriteEntity(metadata.NewSuccessResp(result))
}

func (s *Service) Telnet(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	var data metadata.ParamSubscriptionTelnet
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("telnet subscription failed, decode request body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	callbackUrl := data.CallbackUrl
	uri, err := util.GetDailAddress(callbackUrl)
	if err != nil {
		blog.Errorf("telnet callback failed, err:%+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "bk_callback_url"),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	blog.Infof("telnet url: %, rid: %s", uri, rid)

	conn, err := net.Dial("tcp", uri)
	if err != nil {
		blog.Errorf("telnet callback failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeTelnetFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	conn.Close()

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// trimSubscriptionForm trims space on subscription form.
func (s *Service) trimSubscriptionForm(subscriptionForm string) string {
	subscriptionFormStr := strings.Replace(subscriptionForm, " ", "", -1)
	subscriptionForms := strings.Split(subscriptionFormStr, ",")

	sort.Strings(subscriptionForms)
	return strings.Join(subscriptionForms, ",")
}
