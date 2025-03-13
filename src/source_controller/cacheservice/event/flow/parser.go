/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package flow

import (
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"github.com/tidwall/gjson"
)

// parseEventFunc function type for parsing db event into chain node and detail
type parseEventFunc func(db dal.DB, key event.Key, e *types.Event, id uint64, rid string) (string, *watch.ChainNode,
	*eventDetail, bool, error)

// eventDetail is the parsed event detail
type eventDetail struct {
	// eventInfo is the mongodb event info, it contains the update fields & deleted fields etc.
	eventInfo []byte
	// resDetail is the resource detail of the event, the latest data is stored in general resource detail cache
	resDetail []byte
}

// parseEvent parse event into db chain nodes to store in db and details to store in redis
func parseEvent(db dal.DB, key event.Key, e *types.Event, id uint64, rid string) (string, *watch.ChainNode,
	*eventDetail, bool, error) {

	switch e.OperationType {
	case types.Insert, types.Update, types.Replace:
		// validate the event is valid or not.
		// the invalid event will be dropped.
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return "", nil, nil, false, nil
		}
	case types.Delete:
		// validate the event is valid or not.
		// the invalid event will be dropped.
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return "", nil, nil, false, nil
		}

	// since following event cannot be parsed, skip them and do not retry
	case types.Invalidate:
		blog.Errorf("loop flow, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
		return "", nil, nil, false, nil
	case types.Drop:
		blog.Errorf("loop flow, received drop collection event operation type, **delete object will send a drop "+
			"instance collection event, ignore it**. doc: %s, rid: %s", e.DocBytes, rid)
		return "", nil, nil, false, nil
	default:
		blog.Errorf("loop flow, received unsupported event operation type: %s, doc: %s, rid: %s",
			e.OperationType, e.DocBytes, rid)
		return "", nil, nil, false, nil
	}

	return parseEventToNodeAndDetail(key, e, id, rid)
}

// parseEventToNodeAndDetail parse validated event into db chain nodes to store in db and details to store in redis
func parseEventToNodeAndDetail(key event.Key, e *types.Event, id uint64, rid string) (string, *watch.ChainNode,
	*eventDetail, bool, error) {

	name := key.Name(e.DocBytes)
	instID := key.InstanceID(e.DocBytes)
	currentCursor, err := watch.GetEventCursor(key.Collection(), e, instID)
	if err != nil {
		blog.Errorf("get %s event cursor failed, name: %s, err: %v, oid: %s, rid: %s", key.Collection(), name,
			err, e.ID(), rid)

		monitor.Collect(&meta.Alarm{
			RequestID: rid,
			Type:      meta.FlowFatalError,
			Detail: fmt.Sprintf("run event flow, but get invalid %s cursor, inst id: %d, name: %s",
				key.Collection(), instID, name),
			Module:    types2.CC_MODULE_CACHESERVICE,
			Dimension: map[string]string{"hit_invalid_cursor": "yes"},
		})

		return "", nil, nil, false, err
	}

	chainNode := &watch.ChainNode{
		ID:          id,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		EventType:   watch.ConvertOperateType(e.OperationType),
		Token:       e.Token.Data,
		Cursor:      currentCursor,
	}

	if instID > 0 {
		chainNode.InstanceID = instID
	}

	detail := types.EventInfo{
		UpdatedFields: e.ChangeDesc.UpdatedFields,
		RemovedFields: e.ChangeDesc.RemovedFields,
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run flow, %s, marshal detail failed, name: %s, detail: %+v, err: %v, oid: %s, rid: %s",
			key.Collection(), name, detail, err, e.ID(), rid)
		return "", nil, nil, false, err
	}

	return e.TenantID, chainNode, &eventDetail{eventInfo: detailBytes, resDetail: e.DocBytes}, false, nil
}

// parseInstAsstEvent parse instance association event into db chain nodes to store in db and details to store in redis
func parseInstAsstEvent(db dal.DB, key event.Key, e *types.Event, id uint64, rid string) (string, *watch.ChainNode,
	*eventDetail, bool, error) {

	switch e.OperationType {
	case types.Insert:
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return "", nil, nil, false, nil
		}
	case types.Delete:
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return "", nil, nil, false, nil
		}

	// since following event cannot be parsed, skip them and do not retry
	case types.Invalidate:
		blog.Errorf("loop flow, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
		return "", nil, nil, false, nil
	case types.Drop:
		blog.Errorf("loop flow, received drop collection event operation type, **delete object will send a drop "+
			"instance association collection event, ignore it**. doc: %s, rid: %s", e.DocBytes, rid)
		return "", nil, nil, false, nil
	default:
		blog.Errorf("loop flow, received invalid event op type: %s, doc: %s, rid: %s", e.OperationType, e.DocBytes, rid)
		return "", nil, nil, false, nil
	}

	return parseInstAsstEventToNodeAndDetail(key, e, id, rid)
}

// parseInstAsstEventToNodeAndDetail parse inst asst event into db chain nodes and details
func parseInstAsstEventToNodeAndDetail(key event.Key, e *types.Event, id uint64, rid string) (string, *watch.ChainNode,
	*eventDetail, bool, error) {

	instAsstID := key.InstanceID(e.DocBytes)
	if instAsstID == 0 {
		blog.Errorf("loop flow, received invalid event id, doc: %s, rid: %s", e.DocBytes, rid)
		return "", nil, nil, false, nil
	}

	// since instance association is saved in both source and target object inst asst table, one change will generate 2
	// events, so we change the oid to id so that the cursor of them will be the same for deduplicate
	oid := e.Oid
	e.Oid = strconv.FormatInt(instAsstID, 10)
	currentCursor, err := watch.GetEventCursor(key.Collection(), e, instAsstID)
	if err != nil {
		blog.Errorf("get %s event cursor failed, err: %v, oid: %s, rid: %s", key.Collection(), err, e.ID(), rid)

		monitor.Collect(&meta.Alarm{
			RequestID: rid,
			Type:      meta.FlowFatalError,
			Detail:    fmt.Sprintf("run event flow, but get invalid %s cursor, id: %d", key.Collection(), instAsstID),
			Module:    types2.CC_MODULE_CACHESERVICE,
			Dimension: map[string]string{"hit_invalid_cursor": "yes"},
		})

		return "", nil, nil, false, err
	}

	chainNode := &watch.ChainNode{
		ID:          id,
		Oid:         oid,
		ClusterTime: e.ClusterTime,
		EventType:   watch.ConvertOperateType(e.OperationType),
		Token:       e.Token.Data,
		Cursor:      currentCursor,
		InstanceID:  instAsstID,
	}

	chainNode.SubResource = []string{gjson.GetBytes(e.DocBytes, common.BKObjIDField).String(),
		gjson.GetBytes(e.DocBytes, common.BKAsstObjIDField).String()}

	detail := types.EventInfo{
		UpdatedFields: e.ChangeDesc.UpdatedFields,
		RemovedFields: e.ChangeDesc.RemovedFields,
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run flow, %s, marshal detail failed, detail: %+v, err: %v, oid: %s, rid: %s", key.Collection(),
			detail, err, e.ID(), rid)
		return "", nil, nil, false, err
	}

	return e.TenantID, chainNode, &eventDetail{eventInfo: detailBytes, resDetail: e.DocBytes}, false, nil
}
