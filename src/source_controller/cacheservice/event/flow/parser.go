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
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/watch"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"github.com/tidwall/gjson"
)

// parseEventFunc function type for parsing db event into chain node and detail
type parseEventFunc func(db dal.DB, key event.Key, e *types.Event, oidDetailMap map[oidCollKey][]byte, id uint64,
	rid string) (*watch.ChainNode, []byte, bool, error)

// parseEvent parse event into db chain nodes to store in db and details to store in redis
func parseEvent(db dal.DB, key event.Key, e *types.Event, oidDetailMap map[oidCollKey][]byte, id uint64, rid string) (
	*watch.ChainNode, []byte, bool, error) {

	switch e.OperationType {
	case types.Insert, types.Update, types.Replace:
		// validate the event is valid or not.
		// the invalid event will be dropped.
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}
	case types.Delete:
		doc, exist := oidDetailMap[oidCollKey{oid: e.Oid, coll: e.Collection}]
		if !exist {
			blog.Errorf("run flow, received %s event, but delete doc[oid: %s] detail not exists, rid: %s",
				key.Collection(), e.Oid, rid)
			return nil, nil, false, nil
		}
		// update delete event detail doc bytes.
		e.DocBytes = doc

		// validate the event is valid or not.
		// the invalid event will be dropped.
		if err := key.Validate(doc); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}

	// since following event cannot be parsed, skip them and do not retry
	case types.Invalidate:
		blog.Errorf("loop flow, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
		return nil, nil, false, nil
	case types.Drop:
		blog.Errorf("loop flow, received drop collection event operation type, **delete object will send a drop "+
			"instance collection event, ignore it**. doc: %s, rid: %s", e.DocBytes, rid)
		return nil, nil, false, nil
	default:
		blog.Errorf("loop flow, received unsupported event operation type: %s, doc: %s, rid: %s",
			e.OperationType, e.DocBytes, rid)
		return nil, nil, false, nil
	}

	return parseEventToNodeAndDetail(key, e, id, rid)
}

// parseInstAsstEvent parse instance association event into db chain nodes to store in db and details to store in redis
func parseInstAsstEvent(db dal.DB, key event.Key, e *types.Event, oidDetailMap map[oidCollKey][]byte, id uint64,
	rid string) (*watch.ChainNode, []byte, bool, error) {

	switch e.OperationType {
	case types.Insert:
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}
	case types.Delete:
		doc, exist := oidDetailMap[oidCollKey{oid: e.Oid, coll: e.Collection}]
		if !exist {
			blog.Errorf("%s event delete doc[oid: %s] detail not exists, rid: %s", key.Collection(), e.Oid, rid)
			return nil, nil, false, nil
		}
		// update delete event detail doc bytes from del archive
		e.DocBytes = doc

		if err := key.Validate(doc); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}

	// since following event cannot be parsed, skip them and do not retry
	case types.Invalidate:
		blog.Errorf("loop flow, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
		return nil, nil, false, nil
	case types.Drop:
		blog.Errorf("loop flow, received drop collection event operation type, **delete object will send a drop "+
			"instance association collection event, ignore it**. doc: %s, rid: %s", e.DocBytes, rid)
		return nil, nil, false, nil
	default:
		blog.Errorf("loop flow, received invalid event op type: %s, doc: %s, rid: %s", e.OperationType, e.DocBytes, rid)
		return nil, nil, false, nil
	}

	instAsstID := key.InstanceID(e.DocBytes)
	if instAsstID == 0 {
		blog.Errorf("loop flow, received invalid event id, doc: %s, rid: %s", e.DocBytes, rid)
		return nil, nil, false, nil
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

		return nil, nil, false, err
	}

	chainNode := &watch.ChainNode{
		ID:              id,
		Oid:             oid,
		ClusterTime:     e.ClusterTime,
		EventType:       watch.ConvertOperateType(e.OperationType),
		Token:           e.Token.Data,
		Cursor:          currentCursor,
		InstanceID:      instAsstID,
		SupplierAccount: key.SupplierAccount(e.DocBytes),
	}

	chainNode.SubResource = []string{gjson.GetBytes(e.DocBytes, common.BKObjIDField).String(),
		gjson.GetBytes(e.DocBytes, common.BKAsstObjIDField).String()}

	detail := types.EventDetail{
		Detail:        types.JsonString(e.DocBytes),
		UpdatedFields: e.ChangeDesc.UpdatedFields,
		RemovedFields: e.ChangeDesc.RemovedFields,
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run flow, %s, marshal detail failed, detail: %+v, err: %v, oid: %s, rid: %s", key.Collection(),
			detail, err, e.ID(), rid)
		return nil, nil, false, err
	}

	return chainNode, detailBytes, false, nil
}

// parsePodEvent parse pod events into db chain nodes and details, pod detail includes its containers
func parsePodEvent(db dal.DB, key event.Key, e *types.Event, oidDetailMap map[oidCollKey][]byte, id uint64,
	rid string) (*watch.ChainNode, []byte, bool, error) {

	switch e.OperationType {
	case types.Insert:
		if err := key.Validate(e.DocBytes); err != nil {
			blog.Errorf("pod event is invalid, doc: %s, oid: %s, err: %v, rid: %s", e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}

		podID := key.InstanceID(e.DocBytes)

		// get pod containers from container table and del archive for the scenario of delete right after create
		filter := map[string]interface{}{
			kubetypes.BKPodIDField: podID,
		}

		containers := make([]interface{}, 0)
		err := db.Table(kubetypes.BKTableNameBaseContainer).Find(filter).All(context.Background(), &containers)
		if err != nil {
			blog.Errorf("get pod containers failed, pod id: %d, err: %v, rid: %s", podID, err, rid)
			return nil, nil, true, err
		}

		delContainers, retry, err := getDeletedContainerDetail(db, podID, rid)
		if err != nil {
			return nil, nil, retry, err
		}

		// set container details to deleted pod detail, and update the event detail doc bytes
		podDetail := *e.Document.(*map[string]interface{})
		podDetail["containers"] = append(containers, delContainers...)

		byt, err := json.Marshal(podDetail)
		if err != nil {
			blog.Errorf("marshal pod with container detail(%+v) failed, err: %v, rid: %s", podDetail, err, rid)
			return nil, nil, false, err
		}
		e.DocBytes = byt

	case types.Delete:
		doc, exist := oidDetailMap[oidCollKey{oid: e.Oid, coll: e.Collection}]
		if !exist {
			blog.Errorf("pod event delete doc not exists, oid: %s, rid: %s", e.Oid, rid)
			return nil, nil, false, nil
		}

		if err := key.Validate(doc); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}

		// get deleted container details, put it in deleted pod detail, and update the event detail doc bytes
		podID := key.InstanceID(doc)
		containers, retry, err := getDeletedContainerDetail(db, podID, rid)
		if err != nil {
			return nil, nil, retry, err
		}

		podDetail := make(map[string]interface{})
		if err = json.Unmarshal(doc, &podDetail); err != nil {
			blog.Errorf("unmarshal pod detail(%s) failed, err: %v, rid: %s", string(doc), err, rid)
			return nil, nil, false, err
		}
		podDetail["containers"] = containers

		byt, err := json.Marshal(podDetail)
		if err != nil {
			blog.Errorf("marshal pod with container detail(%+v) failed, err: %v, rid: %s", podDetail, err, rid)
			return nil, nil, false, err
		}
		e.DocBytes = byt

	// since invalid event cannot be parsed, skip them and do not retry
	default:
		blog.Errorf("pod event %s op type %s is invalid, doc: %s, rid: %s", e.Oid, e.OperationType, e.DocBytes, rid)
		return nil, nil, false, nil
	}

	return parseEventToNodeAndDetail(key, e, id, rid)
}

// getDeletedContainerDetail get deleted pod's containers details
func getDeletedContainerDetail(db dal.DB, podID int64, rid string) ([]interface{}, bool, error) {
	filter := map[string]interface{}{
		"detail.bk_pod_id": podID,
		"coll":             kubetypes.BKTableNameBaseContainer,
	}

	docs := make([]map[string]interface{}, 0)
	err := db.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
	if err != nil {
		blog.Errorf("get archive deleted doc failed, filter: %+v, err: %v, rid: %s", filter, err, rid)
		return nil, true, err
	}

	containers := make([]interface{}, len(docs))
	for idx, doc := range docs {
		containers[idx] = doc["detail"]
	}

	return containers, false, nil
}

// parseKubeWorkloadEvent parse kube workload event, its sub resource is its workload kind
func parseKubeWorkloadEvent(db dal.DB, key event.Key, e *types.Event, oidDetailMap map[oidCollKey][]byte, id uint64,
	rid string) (*watch.ChainNode, []byte, bool, error) {

	chainNode, details, retry, err := parseEvent(db, key, e, oidDetailMap, id, rid)
	if err != nil {
		return nil, nil, retry, err
	}

	// ignore invalid event
	if chainNode == nil {
		return nil, nil, false, nil
	}

	// get workload sub resource by event collection
	switch e.Collection {
	case kubetypes.BKTableNameBaseDeployment:
		chainNode.SubResource = []string{string(kubetypes.KubeDeployment)}
	case kubetypes.BKTableNameBaseStatefulSet:
		chainNode.SubResource = []string{string(kubetypes.KubeStatefulSet)}
	case kubetypes.BKTableNameBaseDaemonSet:
		chainNode.SubResource = []string{string(kubetypes.KubeDaemonSet)}
	case kubetypes.BKTableNameGameStatefulSet:
		chainNode.SubResource = []string{string(kubetypes.KubeGameStatefulSet)}
	case kubetypes.BKTableNameGameDeployment:
		chainNode.SubResource = []string{string(kubetypes.KubeGameDeployment)}
	case kubetypes.BKTableNameBaseCronJob:
		chainNode.SubResource = []string{string(kubetypes.KubeCronJob)}
	case kubetypes.BKTableNameBaseJob:
		chainNode.SubResource = []string{string(kubetypes.KubeJob)}
	case kubetypes.BKTableNameBasePodWorkload:
		chainNode.SubResource = []string{string(kubetypes.KubePodWorkload)}

	// since invalid event cannot be parsed, skip them and do not retry
	default:
		blog.Errorf("kube workload event coll %s is invalid, doc: %s, rid: %s", e.Collection, e.DocBytes, rid)
		return nil, nil, false, nil
	}

	return chainNode, details, false, nil
}

// parseEventToNodeAndDetail parse validated event into db chain nodes to store in db and details to store in redis
func parseEventToNodeAndDetail(key event.Key, e *types.Event, id uint64, rid string) (*watch.ChainNode, []byte, bool,
	error) {

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

		return nil, nil, false, err
	}

	chainNode := &watch.ChainNode{
		ID:              id,
		ClusterTime:     e.ClusterTime,
		Oid:             e.Oid,
		EventType:       watch.ConvertOperateType(e.OperationType),
		Token:           e.Token.Data,
		Cursor:          currentCursor,
		SupplierAccount: key.SupplierAccount(e.DocBytes),
	}

	if instID > 0 {
		chainNode.InstanceID = instID
	}

	detail := types.EventDetail{
		Detail:        types.JsonString(e.DocBytes),
		UpdatedFields: e.ChangeDesc.UpdatedFields,
		RemovedFields: e.ChangeDesc.RemovedFields,
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run flow, %s, marshal detail failed, name: %s, detail: %+v, err: %v, oid: %s, rid: %s",
			key.Collection(), name, detail, err, e.ID(), rid)
		return nil, nil, false, err
	}

	return chainNode, detailBytes, false, nil
}
