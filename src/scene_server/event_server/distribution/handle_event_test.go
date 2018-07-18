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
	"encoding/json"
	"testing"

	"configcenter/src/common"
	"configcenter/src/scene_server/event_server/types"
)

func TestPrepareDistInst(t *testing.T) {
	initTester()
	subscriber := "1"
	event := types.EventInstCtx{}
	event.EventType = types.EventTypeInstData
	event.Action = types.EventActionUpdate
	event.ObjType = "object"
	event.CurData = map[string]string{"name": "dog", common.BKObjIDField: "animal"}
	event.PreData = map[string]string{"name": "cat", common.BKObjIDField: "animal"}
	event.RequestID = "1"

	dist := prepareDistInst(subscriber, event)

	if dist.ObjType != "animal" {
		t.Fatalf("expected animal but got %s", dist.ObjType)
	}

	if dist.DstbID <= 0 {
		t.Fatalf("expected dstbid > 0 but got %d", dist.DstbID)
	}
}

func TestPushToQueue(t *testing.T) {
	initTester()
	key := common.BKCacheKeyV3Prefix + "event:testqueue"
	value := `{"dist":"testdata"}`
	if err := pushToQueue(key, value); err != nil {
		t.Fatalf("pushToQueue failed %v", err)
	}
}

func TestNextDistID(t *testing.T) {
	initTester()
	nextID, err := nextDistID("1")
	if err != nil {
		t.Fatalf("nextDistID failed %v", err)
	}
	if nextID <= 0 {
		t.Fatalf("expected nextDistID > 0 but got %d", nextID)
	}
}

func TestSaveEventDone(t *testing.T) {
	initTester()
	eventCtx := testEventCtx()
	err := SaveEventDone(eventCtx)
	if err != nil {
		t.Fatalf("SaveEventDone failed %v", err)
	}
}

func TestCheckFromDone(t *testing.T) {
	initTester()
	_, err := checkFromDone(types.EventCacheEventDoneKey, "0")
	if err != nil {
		t.Fatalf("checkFromDone failed %v", err)
	}
}

func TestFindEventTypeSubscribers(t *testing.T) {
	initTester()
	finded := findEventTypeSubscribers("animalcreate")
	if len(finded) < 1 {
		t.Fatalf("findEventTypeSubscribers failed")
	}
}

func TestHandleInst(t *testing.T) {
	initTester()
	event := testEventCtx()
	err := handleInst(event)
	if err != nil {
		t.Fatalf("handleInst failed %v", err)
	}
}

func testEventCtx() *types.EventInstCtx {
	initTester()
	event := &types.EventInst{}
	event.EventType = types.EventTypeInstData
	event.Action = types.EventActionUpdate
	event.ObjType = "object"
	event.CurData = map[string]string{"name": "dog", common.BKObjIDField: "animal"}
	event.PreData = map[string]string{"name": "cat", common.BKObjIDField: "animal"}
	event.RequestID = "1"

	out, _ := json.Marshal(event)

	eventCtx := &types.EventInstCtx{}
	eventCtx.Raw = string(out)
	eventCtx.EventInst = *event
	return eventCtx
}
