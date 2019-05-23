/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instances

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

// EventHandle event data tmp handle
type EventHandle struct {
	data            map[int64]metadata.EventData
	instanceManager *instanceManager
}

// NewEventHandle new event push handle
func (m *instanceManager) NewEventHandle(objID string) *EventHandle {
	return &EventHandle{
		data:            make(map[int64]metadata.EventData, 0),
		instanceManager: m,
	}
}

// SetPreData set inst befor  data
func (eh *EventHandle) SetPreData(id int64, data interface{}) {

	item, ok := eh.data[id]
	if !ok {
		item = metadata.EventData{
			PreData: data,
		}
	} else {
		item.PreData = data
	}
	eh.data[id] = item
}

// SetCurData set inst current data
func (eh *EventHandle) SetCurData(id int64, data interface{}) {
	item, ok := eh.data[id]
	if !ok {
		item = metadata.EventData{
			CurData: data,
		}
	} else {
		item.CurData = data
	}
	eh.data[id] = item
}

// SetCurDataAndPush get current instance info and push data
func (eh *EventHandle) SetCurDataAndPush(ctx core.ContextParams, objID, eventAction string, cond mapstr.MapStr) error {
	instIDFieldName := common.GetInstIDField(objID)

	insts, _, err := eh.instanceManager.getInsts(ctx, objID, cond)
	if nil != err {
		blog.ErrorJSON("EventHandle SetCurDataAndPush find objID(%s) error. error: %s, condition:%s, rid: %s", objID, cond, ctx.ReqID)
		return err
	}
	for _, inst := range insts {
		id, err := inst.Int64(instIDFieldName)
		if err != nil {
			blog.ErrorJSON("EventHandle SetCurDataAndPush objID(%s) field(%s) convert to int64 error. err:%s, inst:%s, rid:%s", objID, instIDFieldName, err.Error(), inst, ctx.ReqID)
			//  convert %s  field %s to %s error %s
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, objID, instIDFieldName, "integer", err.Error())
		}
		eh.SetCurData(id, inst)
	}
	err = eh.Push(ctx, objID, eventAction)
	if err != nil {
		blog.ErrorJSON("EventHandle SetCurDataAndPush objID(%s) Push event error. event action:%s, condition:%s, rid:%s", objID, eventAction, cond, ctx.ReqID)
		return err
	}
	return nil
}

// Push push event to event server
func (eh *EventHandle) Push(ctx core.ContextParams, objType, eventAction string) error {
	if eh.data == nil {
		return nil
	}
	var eventInstArr []*metadata.EventInst
	for _, item := range eh.data {

		srcevent := eventclient.NewEventWithHeader(ctx.Header)
		srcevent.EventType = metadata.EventTypeInstData
		srcevent.ObjType = objType
		srcevent.Action = eventAction
		srcevent.Data = []metadata.EventData{item}
		eventInstArr = append(eventInstArr, srcevent)
	}
	err := eh.instanceManager.EventC.Push(ctx, eventInstArr...)
	if err != nil {
		blog.ErrorJSON("Push objType(%s) change to event server error. data:%s, rid:%s", objType, eventInstArr, ctx.ReqID)
		ctx.Error.Errorf(common.CCErrEventPushEventFailed)
	}
	return nil
}
