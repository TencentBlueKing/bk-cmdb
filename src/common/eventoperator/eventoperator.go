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

package eventoperator

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

// EventOperator save event data to cache temporarily and push to event server at calling push
type EventOperator struct {
	// eventID --> eventData
	cache    map[int64]metadata.EventData
	eventCli eventclient.Client
	db       dal.DB
}

// NewEventOperator new a event operator
func NewEventOperator(ec eventclient.Client, db dal.DB) *EventOperator {
	return &EventOperator{
		cache:    make(map[int64]metadata.EventData, 0),
		eventCli: ec,
		db:       db,
	}
}

// SetPreData set inst previous data
func (e *EventOperator) SetPreData(eventID int64, data interface{}) {

	item, ok := e.cache[eventID]
	if !ok {
		item = metadata.EventData{
			PreData: data,
		}
	} else {
		item.PreData = data
	}
	e.cache[eventID] = item
}

// SetPreDataByCond set inst by condition before data change
func (e *EventOperator) SetPreDataByCond(kit *rest.Kit, objID string, cond mapstr.MapStr) error {
	insts, err := e.getInsts(kit, objID, cond)
	if nil != err {
		blog.ErrorJSON("SetPreDataByCond failed,  getInsts error: %s, objID:%s, condition:%s, rid: %s", err, objID, cond, kit.Rid)
		return err
	}

	instIDFieldName := common.GetInstIDField(objID)
	for _, inst := range insts {
		id, err := inst.Int64(instIDFieldName)
		if err != nil {
			blog.ErrorJSON("SetPreDataByCond failed, objID(%s) field(%s) convert to int64 error. err:%s, inst:%s, rid:%s", objID, instIDFieldName, err.Error(), inst, kit.Rid)
			// convert %s  field %s to %s error %s
			return kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, objID, instIDFieldName, "integer", err.Error())
		}
		e.SetPreData(id, inst)
	}

	return nil
}

// SetCurData set inst current data
func (e *EventOperator) SetCurData(eventID int64, data interface{}) {
	item, ok := e.cache[eventID]
	if !ok {
		item = metadata.EventData{
			CurData: data,
		}
	} else {
		item.CurData = data
	}
	e.cache[eventID] = item
}

// SetCurDataByCond set inst by condition after data change
func (e *EventOperator) SetCurDataByCond(kit *rest.Kit, objID string, cond mapstr.MapStr) error {
	insts, err := e.getInsts(kit, objID, cond)
	if nil != err {
		blog.ErrorJSON("SetCurDataByCond failed,  getInsts error: %s, objID:%s, condition:%s, rid: %s", err, objID, cond, kit.Rid)
		return err
	}

	instIDFieldName := common.GetInstIDField(objID)
	for _, inst := range insts {
		id, err := inst.Int64(instIDFieldName)
		if err != nil {
			blog.ErrorJSON("SetCurDataByCond failed, objID(%s) field(%s) convert to int64 error. err:%s, inst:%s, rid:%s", objID, instIDFieldName, err.Error(), inst, kit.Rid)
			// convert %s  field %s to %s error %s
			return kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, objID, instIDFieldName, "integer", err.Error())
		}
		e.SetCurData(id, inst)
	}

	return nil
}

// SetCurDataAndPush get current instance info and push data
func (e *EventOperator) SetCurDataAndPush(kit *rest.Kit, objID, eventAction string, cond mapstr.MapStr) error {
	err := e.SetCurDataByCond(kit, objID, cond)
	if err != nil {
		blog.ErrorJSON("SetCurDataAndPush failed,  SetCurDataByCond error: %s, objID:%s, condition:%s, rid: %s", err, objID, cond, kit.Rid)
		return err
	}

	err = e.Push(kit, objID, eventAction)
	if err != nil {
		blog.ErrorJSON("SetCurDataAndPush failed, Push error: %s, objID:%s, eventAction:%s, condition:%s, rid:%s", err, objID, eventAction, cond, kit.Rid)
		return err
	}
	return nil
}

// Push push model instance event to event server
func (e *EventOperator) Push(kit *rest.Kit, objType, eventAction string) error {
	if e.cache == nil || len(e.cache) == 0 {
		return nil
	}
	var eventInstArr []*metadata.EventInst
	for _, item := range e.cache {
		srcEvent := eventclient.NewEventWithHeader(kit.Header)
		srcEvent.EventType = metadata.EventTypeInstData
		srcEvent.ObjType = objType
		srcEvent.Action = eventAction
		srcEvent.Data = []metadata.EventData{item}
		eventInstArr = append(eventInstArr, srcEvent)
	}
	err := e.eventCli.Push(kit.Ctx, eventInstArr...)
	if err != nil {
		blog.ErrorJSON("Push objType(%s) change to event server error. data:%s, rid:%s", objType, eventInstArr, kit.Rid)
		return kit.CCError.Errorf(common.CCErrEventPushEventFailed)
	}
	return nil
}

// getInsts get model instances
func (e *EventOperator) getInsts(kit *rest.Kit, objID string, cond mapstr.MapStr) (result []mapstr.MapStr, err error) {
	result = make([]mapstr.MapStr, 0)
	tableName := common.GetInstTableName(objID)
	if !util.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	if objID == common.BKInnerObjIDHost {
		hosts := make([]metadata.HostMapStr, 0)
		err = e.db.Table(tableName).Find(cond).All(kit.Ctx, &hosts)
		for _, host := range hosts {
			result = append(result, mapstr.MapStr(host))
		}
	} else {
		err = e.db.Table(tableName).Find(cond).All(kit.Ctx, &result)
	}
	return result, err
}
