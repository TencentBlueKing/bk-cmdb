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

func (s *Service) Subscribe(req *restful.Request, resp *restful.Response) {
	var err error
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID := util.GetOwnerID(header)

	sub := &metadata.Subscription{}
	if err = json.NewDecoder(req.Request.Body).Decode(&sub); err != nil {
		blog.Errorf("add subscription, but decode body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	if len(sub.SubscriptionName) == 0 {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionName")})
		return
	}
	if len(sub.CallbackURL) == 0 {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "CallbackURL")})
		return
	}
	if len(sub.SubscriptionForm) == 0 {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionForm")})
		return
	}
	sub.Operator = util.GetUser(req.Request.Header)
	if sub.TimeOutSeconds <= 0 {
		sub.TimeOutSeconds = 10
	}
	if sub.ConfirmMode != metadata.ConfirmModeHTTPStatus && sub.ConfirmMode != metadata.ConfirmModeRegular {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "ConfirmMode")})
		return
	}
	if sub.ConfirmMode == metadata.ConfirmModeHTTPStatus && sub.ConfirmPattern == "" {
		sub.ConfirmPattern = strconv.FormatInt(http.StatusOK, 10)
	}
	now := metadata.Now()
	sub.LastTime = now
	sub.OwnerID = ownerID

	eventTypesStr := strings.Replace(sub.SubscriptionForm, " ", "", -1)
	eventTypes := strings.Split(eventTypesStr, ",")
	sort.Strings(eventTypes)
	sub.SubscriptionForm = strings.Join(eventTypes, ",")

	// do create or update operation
	existSubscriptions := make([]metadata.Subscription, 0)
	filter := map[string]interface{}{
		common.BKSubscriptionNameField: sub.SubscriptionName,
		common.BKOwnerIDField:          ownerID,
	}
	if err := s.db.Table(common.BKTableNameSubscription).Find(filter).All(s.ctx, &existSubscriptions); err != nil {
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommDuplicateItem, common.BKSubscriptionNameField),
		}
		resp.WriteError(http.StatusOK, result)
		return
	}

	if len(existSubscriptions) > 0 {
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommDuplicateItem, common.BKSubscriptionNameField),
		}
		resp.WriteError(http.StatusOK, result)
		return
	}

	// generate id field
	subscriptionID, err := s.db.NextSequence(s.ctx, common.BKTableNameSubscription)
	sub.SubscriptionID = int64(subscriptionID)
	if nil != err {
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// save to storage
	if err := s.db.Table(common.BKTableNameSubscription).Insert(s.ctx, sub); err != nil {
		blog.Errorf("create subscription failed, error:%s, rid: %s", err.Error(), rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// add new add subscriber to event receivers
	for _, eventType := range eventTypes {
		// TODO: how to clear dirty subscription?
		if err := s.cache.SAdd(types.EventSubscriberCacheKey(ownerID, eventType), sub.SubscriptionID).Err(); err != nil {
			blog.Errorf("create subscription failed, add new add subscriber to event receivers failed, error:%s, rid: %s", err.Error(), rid)
			result := &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed)}
			resp.WriteError(http.StatusInternalServerError, result)
			return
		}
	}
	msg, _ := json.Marshal(&sub)
	s.cache.Publish(types.EventCacheProcessChannel, "create"+string(msg))
	s.cache.Del(types.EventCacheDistCallBackCountPrefix + strconv.FormatInt(sub.SubscriptionID, 10))

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

	result := NewCreateSubscriptionResult(sub.SubscriptionID)
	resp.WriteEntity(result)
}

func NewCreateSubscriptionResult(subscriptionID int64) *metadata.RspSubscriptionCreate {
	result := &metadata.RspSubscriptionCreate{
		BaseResp: metadata.SuccessBaseResp,
		Data: struct {
			SubscriptionID int64 `json:"subscription_id"`
		}{
			SubscriptionID: subscriptionID,
		},
	}
	return result
}

func (s *Service) UnSubscribe(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID := util.GetOwnerID(header)

	id, err := strconv.ParseInt(req.PathParameter("subscribeID"), 10, 64)
	if nil != err {
		result := &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// query old Subscription
	sub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()
	if err := s.db.Table(common.BKTableNameSubscription).Find(condition).One(s.ctx, &sub); err != nil {
		blog.Errorf("fail to get subscription by id %v, error information is %v, rid: %s", id, err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeDeleteFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}
	// execute delete command
	if delErr := s.db.Table(common.BKTableNameSubscription).Delete(s.ctx, condition); nil != delErr {
		blog.Errorf("fail to delete subscription by id %v, error information is %v, rid: %s", id, delErr, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeDeleteFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	subID := fmt.Sprint(id)
	eventTypes := strings.Split(sub.SubscriptionForm, ",")
	for _, eventType := range eventTypes {
		eventType = strings.TrimSpace(eventType)
		if err := s.cache.SRem(types.EventSubscriberCacheKey(ownerID, eventType), subID).Err(); err != nil {
			blog.Errorf("delete subscription failed, error:%s, rid: %s", err.Error(), rid)
			result := &metadata.RespError{
				Msg: defErr.Error(common.CCErrEventSubscribeDeleteFailed),
			}
			resp.WriteError(http.StatusInternalServerError, result)
			return
		}
	}

	s.cache.Del(types.EventCacheDistIDPrefix+subID,
		types.EventCacheDistQueuePrefix+subID,
		types.EventCacheDistDonePrefix+subID)

	msg, _ := json.Marshal(&sub)
	s.cache.Publish(types.EventCacheProcessChannel, "delete"+string(msg))

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

func (s *Service) UpdateSubscription(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID := util.GetOwnerID(header)

	id, err := strconv.ParseInt(req.PathParameter("subscribeID"), 10, 64)
	if nil != err {
		result := &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	blog.Infof("update subscription %v, rid: %s", id, rid)

	sub := &metadata.Subscription{}
	if err = json.NewDecoder(req.Request.Body).Decode(&sub); err != nil {
		blog.Errorf("update subscription, but decode body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	if len(sub.SubscriptionName) == 0 {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionName")})
		return
	}
	if len(sub.CallbackURL) == 0 {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "CallbackURL")})
		return
	}
	if len(sub.SubscriptionForm) == 0 {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "SubscriptionForm")})
		return
	}
	if sub.ConfirmMode != metadata.ConfirmModeHTTPStatus && sub.ConfirmMode != metadata.ConfirmModeRegular {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "ConfirmMode")})
		return
	}
	if sub.ConfirmMode == metadata.ConfirmModeHTTPStatus && sub.ConfirmPattern == "" {
		sub.ConfirmPattern = strconv.FormatInt(http.StatusOK, 10)
	}
	sub.Operator = util.GetUser(req.Request.Header)
	if err = s.updateSubscription(header, id, ownerID, sub); err != nil {
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeUpdateFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
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

func (s *Service) updateSubscription(header http.Header, id int64, ownerID string, sub *metadata.Subscription) error {
	rid := util.GetHTTPCCRequestID(header)
	// query old Subscription
	oldSub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()
	if err := s.db.Table(common.BKTableNameSubscription).Find(condition).One(s.ctx, &oldSub); err != nil {
		blog.Errorf("fail to get subscription by id %v, err: %v, rid: %s", id, err, rid)
		return err
	}
	if oldSub.SubscriptionName != sub.SubscriptionName {
		filter := map[string]interface{}{
			common.BKSubscriptionNameField: sub.SubscriptionName,
			common.BKOwnerIDField:          ownerID,
		}
		count, err := s.db.Table(common.BKTableNameSubscription).Find(filter).Count(s.ctx)
		if err != nil {
			blog.Errorf("get subscription count error: %v, rid: %s", err, rid)
			return err
		}
		if count > 0 {
			blog.Errorf("duplicate subscription name, rid: %s", rid)
			return err
		}
	}

	sub.SubscriptionID = oldSub.SubscriptionID
	if sub.TimeOutSeconds <= 0 {
		sub.TimeOutSeconds = 10
	}
	now := metadata.Now()
	sub.LastTime = now
	sub.OwnerID = ownerID

	sub.SubscriptionForm = strings.Replace(sub.SubscriptionForm, " ", "", -1)
	events := strings.Split(sub.SubscriptionForm, ",")
	sort.Strings(events)
	sub.SubscriptionForm = strings.Join(events, ",")

	filter := map[string]interface{}{
		common.BKSubscriptionIDField: id,
		common.BKOwnerIDField:        ownerID,
	}
	if updateErr := s.db.Table(common.BKTableNameSubscription).Update(s.ctx, filter, sub); nil != updateErr {
		blog.Errorf("fail update subscription by condition, error information is %s, rid: %s", updateErr.Error(), rid)
		return updateErr
	}

	eventTypes := strings.Split(sub.SubscriptionForm, ",")
	oldEventTypes := strings.Split(oldSub.SubscriptionForm, ",")

	subs, plugs := util.CalSliceDiff(oldEventTypes, eventTypes)

	for _, eventType := range subs {
		eventType = strings.TrimSpace(eventType)
		if err := s.cache.SRem(types.EventSubscriberCacheKey(ownerID, eventType), id).Err(); err != nil {
			blog.Errorf("delete subscription failed, error:%s, rid: %s", err.Error(), rid)
			return err
		}
	}
	for _, event := range plugs {
		if err := s.cache.SAdd(types.EventSubscriberCacheKey(ownerID, event), sub.SubscriptionID).Err(); err != nil {
			blog.Errorf("create subscription failed, error:%s, rid: %s", err.Error(), rid)
			return err
		}
	}

	mesg, err := json.Marshal(&sub)
	if err != nil {
		return err
	}
	return s.cache.Publish(types.EventCacheProcessChannel, "update"+string(mesg)).Err()
}

func (s *Service) ListSubscriptions(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID := util.GetOwnerID(header)

	blog.Infof("select subscription, rid: %s", rid)

	var data metadata.ParamSubscriptionSearch
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("search subscription, but decode body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
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
		blog.Errorf("get host count error, input:%+v error:%v, rid: %s", data, err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeSelectFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	results := make([]metadata.Subscription, 0)

	if selErr := s.db.Table(common.BKTableNameSubscription).Find(condition).Fields(fields...).Sort(sortOption).Start(uint64(skip)).Limit(uint64(limit)).All(s.ctx, &results); nil != selErr {
		blog.Errorf("select data failed, error information is %s, input:%v, rid: %s", selErr, data, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeSelectFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
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
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

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
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

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
