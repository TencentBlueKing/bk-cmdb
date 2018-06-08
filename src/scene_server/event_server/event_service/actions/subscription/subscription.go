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

package subscription

import (
	"bytes"
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	paraparse "configcenter/src/common/paraparse"
	commontypes "configcenter/src/common/types"
	"configcenter/src/common/util"
	sencecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/common/instdata"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	redis "gopkg.in/redis.v5"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
)

var eventSubscription = &subscriptionAction{}

type subscriptionAction struct {
	base.BaseAction
}

// Collect collect events
func (cli *subscriptionAction) Subscribe(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		blog.Info("add subscription")
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		sub := &types.Subscription{}
		if err = json.Unmarshal([]byte(value), sub); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		now := commontypes.Now()
		sub.Operator = sencecommon.GetUserFromHeader(req)
		if sub.TimeOut <= 0 {
			sub.TimeOut = 10
		}
		if sub.ConfirmMode == types.ConfirmmodeHttpstatus &&
			sub.ConfirmPattern == "" {
			sub.ConfirmPattern = "200"
		}
		sub.LastTime = &now
		sub.SubscriptionForm = strings.Replace(sub.SubscriptionForm, " ", "", 0)
		sub.OwnerID = ownerID

		count, err := instdata.GetSubscriptionCntByCondition(map[string]interface{}{"subscription_name": sub.SubscriptionName})
		if err != nil || count > 0 {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDuplicateItem)
		}
		// save to the storage
		if _, err := instdata.CreateSubscription(sub); err != nil {
			blog.Error("create subscription failed, error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeInsertFailed)
		}

		// save to subscribeform in cache
		events := strings.Split(sub.SubscriptionForm, ",")
		for _, event := range events {
			cacheValue := common.KvMap{
				"key":    types.EventCacheSubscribeformKey + event,
				"values": []string{fmt.Sprint(sub.SubscriptionID)},
			}
			if _, err := cli.CC.CacheCli.Insert("sadd", cacheValue); err != nil {
				blog.Error("create subscription failed, error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeInsertFailed)
			}
		}

		mesg, _ := json.Marshal(&sub)
		redisCli := cli.CC.CacheCli.GetSession().(*redis.Client)
		redisCli.Publish(types.EventCacheProcessChannel, "create"+string(mesg))
		redisCli.Del(types.EventCacheDistCallBackCountPrefix + fmt.Sprint(sub.SubscriptionID))

		info := make(map[string]int64)
		info[common.BKSubscriptionIDField] = sub.SubscriptionID
		return http.StatusOK, info, nil
	}, resp)

}

func (cli *subscriptionAction) UnSubscribe(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		var id int64
		pathParameters := req.PathParameters()
		if nil != cli.GetParams(cli.CC, &pathParameters, "subscribeID", &id, resp) {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		blog.Info("delete Subscribe %v", id)

		// query old Subscription
		sub := types.Subscription{}
		condiction := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()
		if err := instdata.GetOneSubscriptionByCondition(condiction, &sub); err != nil {
			blog.Error("fail to get subscription by id %v, error information is %v", id, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrEventSubscribeDeleteFailed)
		}
		// execute delete command
		if delerr := instdata.DelSubscriptionByCondition(condiction); nil != delerr {
			blog.Error("fail to delete subscription by id %v, error information is %v", id, delerr)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrEventSubscribeDeleteFailed)
		}

		subID := fmt.Sprint(id)
		redisCli := cli.CC.CacheCli.GetSession().(*redis.Client)
		eventTypes := strings.Split(sub.SubscriptionForm, ",")
		for _, eventType := range eventTypes {
			eventType = strings.TrimSpace(eventType)
			if err := redisCli.SRem(types.EventCacheSubscribeformKey+eventType, subID).Err(); err != nil {
				blog.Error("delete subscription failed, error:%s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrEventSubscribeDeleteFailed)
			}
		}

		redisCli.Del(types.EventCacheDistIDPrefix+subID,
			types.EventCacheDistQueuePrefix+subID,
			types.EventCacheDistDonePrefix+subID)

		mesg, _ := json.Marshal(&sub)
		redisCli.Publish(types.EventCacheProcessChannel, "delete"+string(mesg))

		return http.StatusOK, nil, nil
	}, resp)
}

func (cli *subscriptionAction) Rebook(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		var id int64
		pathParameters := req.PathParameters()
		if nil != cli.GetParams(cli.CC, &pathParameters, "subscribeID", &id, resp) {
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "subscription_id")
		}
		blog.Info("update subscription %v", id)

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		sub := &types.Subscription{}
		if jserr := json.Unmarshal([]byte(value), sub); nil != jserr {
			blog.Error("unmarshal json failed, error information is %v", jserr)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// query old Subscription
		oldsub := types.Subscription{}
		condiction := util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()
		if err := instdata.GetOneSubscriptionByCondition(condiction, &oldsub); err != nil {
			blog.Error("fail to get subscription by id %v, error information is %v", id, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrEventSubscribeUpdateFailed)
		}
		if oldsub.SubscriptionName != sub.SubscriptionName {
			count, err := instdata.GetSubscriptionCntByCondition(map[string]interface{}{"subscription_name": sub.SubscriptionName})
			if err != nil {
				blog.Errorf("get subscription count error: %v", err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeUpdateFailed)
			}
			if count > 0 {
				blog.Error("duplicate subscription name")
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDuplicateItem)
			}
		}

		sub.SubscriptionID = oldsub.SubscriptionID
		if sub.TimeOut <= 0 {
			sub.TimeOut = 10
		}
		now := commontypes.Now()
		sub.LastTime = &now
		sub.SubscriptionForm = strings.Replace(sub.SubscriptionForm, " ", "", 0)
		sub.Operator = sencecommon.GetUserFromHeader(req)
		if updateerr := instdata.UpdateSubscriptionByCondition(sub, util.NewMapBuilder(common.BKSubscriptionIDField, id, common.BKOwnerIDField, ownerID).Build()); nil != updateerr {
			blog.Error("fail update subscription by condition, error information is %s", updateerr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeUpdateFailed)
		}

		eventTypes := strings.Split(sub.SubscriptionForm, ",")
		oldeventTypes := strings.Split(oldsub.SubscriptionForm, ",")

		subs, plugs := util.CalSliceDiff(oldeventTypes, eventTypes)

		for _, eventType := range subs {
			eventType = strings.TrimSpace(eventType)
			cacheValue := common.KvMap{
				"key":    types.EventCacheSubscribeformKey + eventType,
				"values": []interface{}{fmt.Sprint(id)},
			}
			if err := cli.CC.CacheCli.DelByCondition("srem", cacheValue); err != nil {
				blog.Error("delete subscription failed, error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeUpdateFailed)
			}
		}
		for _, event := range plugs {
			cacheValue := common.KvMap{
				"key":    types.EventCacheSubscribeformKey + event,
				"values": []string{fmt.Sprint(sub.SubscriptionID)},
			}
			if _, err := cli.CC.CacheCli.Insert("sadd", cacheValue); err != nil {
				blog.Error("create subscription failed, error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeUpdateFailed)
			}
		}

		mesg, _ := json.Marshal(&sub)
		redisCli := cli.CC.CacheCli.GetSession().(*redis.Client)
		redisCli.Publish(types.EventCacheProcessChannel, "update"+string(mesg))

		return http.StatusOK, nil, nil
	}, resp)
}

func (cli *subscriptionAction) Query(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		blog.Info("select subscription")

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var dat paraparse.SubscribeCommonSearch
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("get subscription: input:%v error:%v", value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		fields := dat.Fields
		condition := dat.Condition
		condition = util.SetModOwner(condition, ownerID)

		skip := dat.Page.Start
		limit := dat.Page.Limit
		if limit <= 0 {
			limit = common.BKNoLimit
		}
		sort := dat.Page.Sort

		count, err := instdata.GetSubscriptionCntByCondition(condition)
		if err != nil {
			blog.Error("get host count error, input:%v error:%v", value, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeSelectFailed)
		}
		results := []*types.Subscription{}
		blog.Debug("selector:%+v", condition)
		if selerr := instdata.GetSubscriptionByCondition(fields, condition, &results, sort, skip, limit); nil != selerr {
			blog.Error("select data failed, error information is %s, input:%v", selerr, dat)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribeSelectFailed)
		}

		redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
		for _, sub := range results {
			val := redisCli.HGetAll(types.EventCacheDistCallBackCountPrefix + fmt.Sprint(sub.SubscriptionID)).Val()
			failue, _ := strconv.ParseInt(val["failue"], 10, 64)
			total, _ := strconv.ParseInt(val["total"], 10, 64)
			sub.Statistics = &types.Statistics{
				Total:   total,
				Failure: failue,
			}
		}

		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = results
		return http.StatusOK, info, nil
	}, resp)

}

func (cli *subscriptionAction) Ping(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		pjson := gjson.ParseBytes(value)
		callbackurl := pjson.Get("callback_url").String()
		callbackBody := pjson.Get("data").String()

		blog.Infof("requesting callback: %v,%s", callbackurl, callbackBody)
		callbackreq, _ := http.NewRequest("POST", callbackurl, bytes.NewBufferString(callbackBody))
		callbackResp, err := http.DefaultClient.Do(callbackreq)
		if err != nil {
			blog.Error("test distribute error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrEventSubscribePingFailed)
		}
		defer callbackResp.Body.Close()

		callbackRespBody, err := ioutil.ReadAll(callbackResp.Body)
		if err != nil {
			blog.Error("test distribute error:%v", err)
		}
		result := map[string]interface{}{}
		result["http_status"] = callbackResp.StatusCode
		result["response_body"] = string(callbackRespBody)

		return http.StatusOK, result, nil
	}, resp)
}

func (cli *subscriptionAction) Telnet(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		pjson := gjson.ParseBytes(value)
		callbackurl := pjson.Get("callback_url").String()
		uri, err := util.GetDailAddress(callbackurl)
		if err != nil {
			blog.Error("telent callback error:%v", err)
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "bk_callback_url")
		}
		blog.Infof("telnet %", uri)

		conn, err := net.Dial("tcp", uri)
		if err != nil {
			blog.Error("telent callback error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrEventSubscribeTelnetFailed)
		}
		conn.Close()

		return http.StatusOK, nil, nil
	}, resp)
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/subscribe/search/{ownerID}/{appID}", Params: nil, Handler: eventSubscription.Query})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/subscribe/ping", Params: nil, Handler: eventSubscription.Ping})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/subscribe/telnet", Params: nil, Handler: eventSubscription.Telnet})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/subscribe/{ownerID}/{appID}", Params: nil, Handler: eventSubscription.Subscribe})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/subscribe/{ownerID}/{appID}/{subscribeID}", Params: nil, Handler: eventSubscription.UnSubscribe})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/subscribe/{ownerID}/{appID}/{subscribeID}", Params: nil, Handler: eventSubscription.Rebook})

	// create cc subscription
	eventSubscription.CreateAction()
}
