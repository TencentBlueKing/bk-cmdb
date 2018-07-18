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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"configcenter/src/scene_server/event_server/types"
)

func TestSaveDistDone(t *testing.T) {
	initTester()
	dist := &types.DistInstCtx{}
	if err := saveDistDone(dist); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestHandleDist(t *testing.T) {
	initTester()
	dist := &types.DistInst{}
	dist.SubscriptionID = 1
	dist.DstbID = 1
	dist.EventType = "create"
	dist.Action = "create"
	dist.ObjType = "animal"
	dist.CurData = map[string]string{"name": "dog"}
	dist.PreData = map[string]string{"name": "cat"}
	dist.RequestID = "1"

	expect, err := json.Marshal(dist)
	if err != nil {
		t.Fatalf("marshal failed %v", err)
	}

	f := func(resp http.ResponseWriter, req *http.Request) {
		readed, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read failed %v", err)
		}
		if !bytes.Equal(expect, readed) {
			t.Fatalf("expect %s, but receive %s", expect, readed)
		}
	}
	s := httptest.NewServer(http.HandlerFunc(f))
	s.Start()
	defer s.Close()
	sub := &types.Subscription{
		SubscriptionID:   1,
		SubscriptionName: "testsubscription",
		SystemName:       "testsystem",
		CallbackURL:      s.URL,
		ConfirmMode:      types.ConfirmmodeHttpstatus,
		ConfirmPattern:   "200",
		TimeOut:          10,
		SubscriptionForm: "hostadd",
	}
	distCtx := &types.DistInstCtx{}
	distCtx.DistInst = *dist
	distCtx.Raw = string(expect)
	if err := handleDist(sub, distCtx); err != nil {
		t.Fatalf("%v", err)
	}
}
