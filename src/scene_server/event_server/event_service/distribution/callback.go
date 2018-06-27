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
	"bytes"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/event_server/types"
	"fmt"
	redis "gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func SendCallback(receiver *types.Subscription, event string) (err error) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	redisCli.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "total", 1)

	body := bytes.NewBufferString(event)
	req, err := http.NewRequest("POST", receiver.CallbackURL, body)
	if err != nil {
		redisCli.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
		return fmt.Errorf("event distribute fail, build request error: %v, date=[%s]", err, event)
	}
	var duration time.Duration
	if receiver.TimeOut == 0 {
		duration = timeout
	} else {
		duration = receiver.GetTimeout()
	}
	resp, err := httpCli.DoWithTimeout(duration, req)
	if err != nil {
		redisCli.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
		return fmt.Errorf("event distribute fail, send request error: %v, date=[%s]", err, event)
	}
	defer resp.Body.Close()
	respdata, _ := ioutil.ReadAll(resp.Body)
	if receiver.ConfirmMode == types.ConfirmmodeHttpstatus {
		if strconv.Itoa(resp.StatusCode) != receiver.ConfirmPattern {
			redisCli.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
			return fmt.Errorf("event distribute fail, received response %s, date=[%s]", respdata, event)
		}
	} else if receiver.ConfirmMode == types.ConfirmmodeRegular {
		pattern, err := regexp.Compile(receiver.ConfirmPattern)
		if err != nil {
			return fmt.Errorf("event distribute fail, build regexp error: %v", err)
		}
		if !pattern.Match(respdata) {
			redisCli.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
			return fmt.Errorf("event distribute fail, received response %s, date=[%s]", respdata, event)
		}
		return nil
	}

	return
}

var httpCli = httpclient.NewHttpClient()
