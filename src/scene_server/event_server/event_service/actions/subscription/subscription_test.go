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
	"configcenter/src/common/core/cc/config"
	paraparse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/event_service"
	"configcenter/src/scene_server/event_server/types"
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
)

func ccapiTester(t *testing.T) *httptest.Server {
	app, err := ccapi.NewCCAPIServer(&config.CCAPIConfig{ExConfig: findConf()})
	require.NoError(t, err)
	app.InitHttpServ(nil)
	return httptest.NewServer(app.HttpServ.GetWebContainer())
}

func findConf() string {
	cur, _ := filepath.Abs(".")
	for len(cur) > 1 {
		if tmp := cur + "/api.conf"; util.FileExists(tmp) {
			return tmp
		} else if tmp = cur + "/conf/api.conf"; util.FileExists(tmp) {
			return tmp
		}
		cur = filepath.Dir(cur)
	}
	return ""
}

func TestSubscribe(t *testing.T) {
	server := ccapiTester(t)
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	defer server.Close()
	id := fmt.Sprint(e.POST("/event/v1/subscribe/0/0").WithJSON(types.Subscription{
		SubscriptionName: "testname",
		SystemName:       "testsystem",
		CallbackURL:      "http://127.0.0.1:58080/callback",
		ConfirmMode:      "httpstatus",
		ConfirmPattern:   "200",
		TimeOut:          10,
		SubscriptionForm: "hostcreate,hostupdate,hostdelete",
	}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("data").Object().Value("ID").Number().Gt(0).Raw())

	e.PUT("/event/v1/subscribe/0/0/"+id).WithJSON(types.Subscription{
		SubscriptionName: "testname",
		SystemName:       "testsystem",
		CallbackURL:      "http://127.0.0.1:58080/callback",
		ConfirmMode:      "httpstatus",
		ConfirmPattern:   "200",
		TimeOut:          16,
		SubscriptionForm: "hostcreate,hostupdate,hostdelete",
	}).
		Expect().Status(http.StatusOK).
		JSON().Object().ValueEqual("data", "success")

	i, _ := strconv.ParseInt(id, 10, 64)
	e.POST("/event/v1/subscribe/search/1/1").WithJSON(
		paraparse.SubscribeCommonSearch{
			Condition: map[string]interface{}{"SubscriptionID": i},
		}).
		Expect().JSON().Object().Path("$.data.info[0].TimeOut").Equal(16)

	e.DELETE("/event/v1/subscribe/0/0/"+id).
		Expect().
		JSON().Object().ValueEqual("data", "success")

	e.POST("/event/v1/subscribe/ping").WithJSON(map[string]interface{}{
		"CallbackURL": "http://127.0.0.1:58080/callback",
		"Data":        types.DistInst{},
	}).
		Expect().
		JSON().Object().ValueEqual("result", false)

	e.POST("/event/v1/subscribe/telnet").WithJSON(map[string]interface{}{
		"CallbackURL": "http://127.0.0.1:58080/callback",
		"Data":        types.DistInst{},
	}).
		Expect().
		JSON().Object().ValueEqual("result", false)
}

func TestTelnet(t *testing.T) {
	uri := "http://127.0.0.1:58080/callback"
	uri, err := getDailAddress(uri)
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:58080", uri)

	uri = "http://www.qq.com/callback"
	uri, err = getDailAddress(uri)
	require.NoError(t, err)
	require.Equal(t, "www.qq.com:80", uri)

	uri = "http://www.qq.com:80/callback"
	uri, err = getDailAddress(uri)
	require.NoError(t, err)
	require.Equal(t, "www.qq.com:80", uri)

}
