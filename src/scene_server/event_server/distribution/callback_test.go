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
	"net/http"
	"net/http/httptest"
	"testing"

	"configcenter/src/common/core/cc/api"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal"
)

func initTester() {
	config := map[string]string{
		dal.RDB_REDIS + ".host":     "127.0.0.1",
		dal.RDB_REDIS + ".port":     "6379",
		dal.RDB_REDIS + ".usr":      "cc",
		dal.RDB_REDIS + ".pwd":      "cc",
		dal.RDB_REDIS + ".database": "0",

		dal.RDB_MONGO + ".host":     "127.0.0.1",
		dal.RDB_MONGO + ".port":     "27017",
		dal.RDB_MONGO + ".usr":      "cc",
		dal.RDB_MONGO + ".pwd":      "cc",
		dal.RDB_MONGO + ".database": "0",
	}
	a := api.NewAPIResource()
	a.GetDataCli(config, dal.RDB_REDIS)
	a.GetDataCli(config, dal.RDB_MONGO)
}

func TestSendCallback(t *testing.T) {
	initTester()
	f := func(http.ResponseWriter, *http.Request) {}
	s := httptest.NewServer(http.HandlerFunc(f))
	defer s.Close()
	var receiver = &types.Subscription{
		CallbackURL:    s.URL,
		ConfirmMode:    types.ConfirmmodeHttpstatus,
		ConfirmPattern: "200",
		TimeOut:        10,
	}
	if err := SendCallback(receiver, "test message"); err != nil {
		t.Fail()
	}
}
