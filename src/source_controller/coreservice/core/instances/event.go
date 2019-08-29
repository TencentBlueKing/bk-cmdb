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

// EventClient save event data to cache temporarily and push to event server at calling push
type EventClient struct {
	// eventID --> eventData
	cache           map[int64]metadata.EventData
	instanceManager *instanceManager
	eventCli        eventclient.Client
}

// NewEventHandle new event client
func (m *instanceManager) NewEventClient(objID string) *EventClient {
	return &EventClient{
		cache:           make(map[int64]metadata.EventData, 0),
		instanceManager: m,
		eventCli:        m.EventCli,
	}
}

// SetPreData set inst before data
func (eh *EventClient) SetPreData(eventID int64, data interface{}) {

	item, ok := eh.cache[eventID]
	if !ok {
		item = metadata.EventData{
			PreData: data,
		}
	} else {
		item.PreData = data
	}
	eh.cache[eventID] = item
}

// SetCurData set inst current data
func (eh *EventClient) SetCurData(eventID int64, data interface{}) {
	item, ok := eh.cache[eventID]
	if !ok {
		item = metadata.EventData{
			CurData: data,
		}
	} else {
		item.CurData = data
	}
	eh.cache[eventID] = item
}

// SetCurDataAndPush get current instance info and push data
func (eh *EventClient) SetCurDataAndPush(ctx core.ContextParams, objID, eventAction string, cond mapstr.MapStr) error {
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
			// convert %s  field %s to %s error %s
			return ctx.Error.Errorf(common.CCErrCommInstFieldConvertFail, objID, instIDFieldName, "integer", err.Error())
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
func (eh *EventClient) Push(ctx core.ContextParams, objType, eventAction string) error {
	if eh.cache == nil {
		return nil
	}
	var eventInstArr []*metadata.EventInst
	for _, item := range eh.cache {
		srcEvent := eventclient.NewEventWithHeader(ctx.Header)
		srcEvent.EventType = metadata.EventTypeInstData
		srcEvent.ObjType = objType
		srcEvent.Action = eventAction
		srcEvent.Data = []metadata.EventData{item}
		eventInstArr = append(eventInstArr, srcEvent)
	}
	err := eh.eventCli.Push(ctx, eventInstArr...)
	if err != nil {
		blog.ErrorJSON("Push objType(%s) change to event server error. data:%s, rid:%s", objType, eventInstArr, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrEventPushEventFailed)
	}
	return nil
}
