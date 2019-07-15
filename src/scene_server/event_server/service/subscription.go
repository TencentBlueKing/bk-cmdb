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

// Subscribe  Subscribe events
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
	now := metadata.Now()
	sub.Operator = util.GetUser(req.Request.Header)
	if sub.TimeOut <= 0 {
		sub.TimeOut = 10
	}
	if sub.ConfirmMode == metadata.ConfirmmodeHttpstatus && sub.ConfirmPattern == "" {
		sub.ConfirmPattern = "200"
	}
	sub.LastTime = now
	sub.OwnerID = ownerID

	sub.SubscriptionForm = strings.Replace(sub.SubscriptionForm, " ", "", -1)

	events := strings.Split(sub.SubscriptionForm, ",")
	sort.Strings(events)
	sub.SubscriptionForm = strings.Join(events, ",")

	exists := make([]metadata.Subscription, 0)
	err = s.db.Table(common.BKTableNameSubscription).Find(map[string]interface{}{common.BKSubscriptionNameField: sub.SubscriptionName, common.BKOwnerIDField: ownerID}).All(s.ctx, &exists)
	if err != nil {
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommDuplicateItem, "subscription_name"),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if len(exists) > 0 {
		if err = s.rebook(header, exists[0].SubscriptionID, ownerID, sub); err != nil {
			result := &metadata.RespError{
				Msg: defErr.Error(common.CCErrEventSubscribeUpdateFailed),
			}
			resp.WriteError(http.StatusBadRequest, result)
			return
		}
	} else {
		nid, err := s.db.NextSequence(s.ctx, common.BKTableNameSubscription)
		sub.SubscriptionID = int64(nid)
		if nil != err {
			result := &metadata.RespError{
				Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed),
			}
			resp.WriteError(http.StatusInternalServerError, result)
			return
		}
		// save to the storage
		if err := s.db.Table(common.BKTableNameSubscription).Insert(s.ctx, sub); err != nil {
			blog.Errorf("create subscription failed, error:%s, rid: %s", err.Error(), rid)
			result := &metadata.RespError{
				Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed),
			}
			resp.WriteError(http.StatusInternalServerError, result)
			return
		}

		// save to subscribeForm in cache
		for _, event := range events {
			if err := s.cache.SAdd(types.EventSubscriberCacheKey(ownerID, event), sub.SubscriptionID).Err(); err != nil {
				blog.Errorf("create subscription failed, error:%s, rid: %s", err.Error(), rid)
				result := &metadata.RespError{Msg: defErr.Error(common.CCErrEventSubscribeInsertFailed)}
				resp.WriteError(http.StatusInternalServerError, result)
				return
			}
		}

		mesg, _ := json.Marshal(&sub)
		s.cache.Publish(types.EventCacheProcessChannel, "create"+string(mesg))
		s.cache.Del(types.EventCacheDistCallBackCountPrefix + fmt.Sprint(sub.SubscriptionID))
	}

	if err = s.auth.RegisterResource(s.ctx, meta.ResourceAttribute{
		Basic: meta.Basic{
			Name:       sub.SubscriptionName,
			Type:       meta.EventPushing,
			InstanceID: sub.SubscriptionID,
		},
	}); err != nil {
		blog.Errorf("permission Deny for create subscribe, %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommRegistResourceToIAMFailed, err)}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	result := metadata.RspSubscriptionCreate{
		BaseResp: metadata.SuccessBaseResp,
		Data: struct {
			SubscriptionID int64 `json:"subscription_id"`
		}{
			SubscriptionID: sub.SubscriptionID,
		},
	}
	resp.WriteEntity(result)

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
		}
	}

	s.cache.Del(types.EventCacheDistIDPrefix+subID,
		types.EventCacheDistQueuePrefix+subID,
		types.EventCacheDistDonePrefix+subID)

	msg, _ := json.Marshal(&sub)
	s.cache.Publish(types.EventCacheProcessChannel, "delete"+string(msg))

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) Rebook(req *restful.Request, resp *restful.Response) {
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
	sub.Operator = util.GetUser(req.Request.Header)
	if err = s.rebook(header, id, ownerID, sub); err != nil {
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeUpdateFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) rebook(header http.Header, id int64, ownerID string, sub *metadata.Subscription) error {
	rid := util.GetHTTPCCRequestID(header)
	// query old Subscription
	oldSub := metadata.Subscription{}
	condition := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()
	if err := s.db.Table(common.BKTableNameSubscription).Find(condition).One(s.ctx, &oldSub); err != nil {
		blog.Errorf("fail to get subscription by id %v, error information is %v, rid: %s", id, err, rid)
		return err
	}
	if oldSub.SubscriptionName != sub.SubscriptionName {
		count, err := s.db.Table(common.BKTableNameSubscription).Find(map[string]interface{}{common.BKSubscriptionNameField: sub.SubscriptionName, common.BKOwnerIDField: ownerID}).Count(s.ctx)
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
	if sub.TimeOut <= 0 {
		sub.TimeOut = 10
	}
	now := metadata.Now()
	sub.LastTime = now
	sub.OwnerID = ownerID

	sub.SubscriptionForm = strings.Replace(sub.SubscriptionForm, " ", "", -1)
	events := strings.Split(sub.SubscriptionForm, ",")
	sort.Strings(events)
	sub.SubscriptionForm = strings.Join(events, ",")

	if updateErr := s.db.Table(common.BKTableNameSubscription).Update(s.ctx, util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build(), sub); nil != updateErr {
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

func (s *Service) Query(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID := util.GetOwnerID(header)

	blog.Infof("select subscription, rid: %s", rid)

	var dat metadata.ParamSubscriptionSearch
	if err := json.NewDecoder(req.Request.Body).Decode(&dat); err != nil {
		blog.Errorf("search subscription, but decode body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	fields := dat.Fields
	condition := dat.Condition
	condition = util.SetModOwner(condition, ownerID)

	skip := dat.Page.Start
	limit := dat.Page.Limit
	if limit <= 0 {
		limit = common.BKNoLimit
	}
	sortOption := dat.Page.Sort

	count, err := s.db.Table(common.BKTableNameSubscription).Find(condition).Count(s.ctx)
	if err != nil {
		blog.Errorf("get host count error, input:%+v error:%v, rid: %s", dat, err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeSelectFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	results := make([]metadata.Subscription, 0)
	blog.Debug("selector:%+v, rid: %s", condition, rid)

	if selErr := s.db.Table(common.BKTableNameSubscription).Find(condition).Fields(fields...).Sort(sortOption).Start(uint64(skip)).Limit(uint64(limit)).All(s.ctx, &results); nil != selErr {
		blog.Errorf("select data failed, error information is %s, input:%v, rid: %s", selErr, dat, rid)
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

	var dat metadata.ParamSubscriptionTestCallback
	if err := json.NewDecoder(req.Request.Body).Decode(&dat); err != nil {
		blog.Errorf("ping subscription, but decode body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	callbackUrl := dat.CallbackUrl
	callbackBody := dat.Data

	blog.Infof("requesting callback: %v,%s, rid: %s", callbackUrl, callbackBody, rid)
	callbackReq, _ := http.NewRequest("POST", callbackUrl, bytes.NewBufferString(callbackBody))
	callbackResp, err := http.DefaultClient.Do(callbackReq)
	if err != nil {
		blog.Errorf("test distribute error:%v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribePingFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	defer callbackResp.Body.Close()

	callbackRespBody, err := ioutil.ReadAll(callbackResp.Body)
	if err != nil {
		blog.Errorf("test distribute error:%v, rid: %s", err, rid)
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
	var dat metadata.ParamSubscriptionTelnet
	if err := json.NewDecoder(req.Request.Body).Decode(&dat); err != nil {
		blog.Errorf("telnet subscription, but decode body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	callbackUrl := dat.CallbackUrl
	uri, err := util.GetDailAddress(callbackUrl)
	if err != nil {
		blog.Errorf("telnet callback failed, err:%v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "bk_callback_url"),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	blog.Infof("telnet %, rid: %s", uri, rid)

	conn, err := net.Dial("tcp", uri)
	if err != nil {
		blog.Errorf("telnet callback error:%v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrEventSubscribeTelnetFailed),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	conn.Close()

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}
