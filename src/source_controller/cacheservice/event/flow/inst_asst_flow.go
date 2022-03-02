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

package flow

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"github.com/tidwall/gjson"
)

func newInstAsstFlow(ctx context.Context, opts flowOptions, getDeleteEventDetails getDeleteEventDetailsFunc,
	parseEvent parseEventFunc) error {

	flow, err := NewFlow(opts, getDeleteEventDetails, parseEvent)
	if err != nil {
		return err
	}
	instAsstFlow := InstAsstFlow{
		Flow: flow,
	}

	return instAsstFlow.RunFlow(ctx)
}

// InstAsstFlow instance association event watch flow
type InstAsstFlow struct {
	Flow
}

// RunFlow run instance association event watch flow
func (f *InstAsstFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.watchDB, f.metrics)

	startAtTime, err := f.tokenHandler.getStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get start watch time for %s failed, err: %v", f.key.Collection(), err)
		return err
	}

	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: f.EventStruct,
			// watch all tables with the prefix of instance association table
			CollectionFilter: map[string]interface{}{
				common.BKDBLIKE: event.InstAsstTablePrefixRegex,
			},
			StartAtTime:             startAtTime,
			WatchFatalErrorCallback: f.tokenHandler.resetWatchToken,
		},
	}

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         f.key.Namespace(),
			WatchOpt:     watchOpts,
			TokenHandler: f.tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: batchSize,
	}

	if err := f.watch.WithBatch(opts); err != nil {
		blog.Errorf("run flow, but watch batch failed, err: %v", err)
		return err
	}

	return nil
}

// parseInstAsstEvent parse instance association event into db chain nodes to store in db and details to store in redis
func parseInstAsstEvent(key event.Key, e *types.Event, oidDetailMap map[oidCollKey][]byte, id uint64, rid string) (
	*watch.ChainNode, []byte, bool, error) {

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
			blog.Errorf("run flow, received %s event, but delete doc[oid: %s] detail not exists, rid: %s",
				key.Collection(), e.Oid, rid)
			return nil, nil, false, nil
		}
		// update delete event detail doc bytes from del archive
		e.DocBytes = doc

		if err := key.Validate(doc); err != nil {
			blog.Errorf("run flow, received %s event, but got invalid event, doc: %s, oid: %s, err: %v, rid: %s",
				key.Collection(), e.DocBytes, e.Oid, err, rid)
			return nil, nil, false, nil
		}
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
		ID:          id,
		Oid:         oid,
		ClusterTime: e.ClusterTime,
		EventType:   watch.ConvertOperateType(e.OperationType),
		Token:       e.Token.Data,
		Cursor:      currentCursor,
		InstanceID:  instAsstID,
	}

	objID := gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()
	asstObjID := gjson.GetBytes(e.DocBytes, common.BKAsstObjIDField).String()
	chainNode.SubResource = []string{objID, asstObjID}

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
