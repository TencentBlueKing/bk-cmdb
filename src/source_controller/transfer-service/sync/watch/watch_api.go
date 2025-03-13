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

package watch

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/transfer-service/sync/util"
)

// resTypeCursorMap is cmdb data sync resource type to event cursor type map
var resTypeCursorMap = map[types.ResType][]watch.CursorType{
	types.Biz:             {watch.Biz},
	types.Set:             {watch.Set},
	types.Module:          {watch.Module},
	types.Host:            {watch.Host},
	types.HostRelation:    {watch.ModuleHostRelation},
	types.ObjectInstance:  {watch.ObjectBase, watch.MainlineInstance},
	types.InstAsst:        {watch.InstAsst},
	types.Process:         {watch.Process},
	types.ProcessRelation: {watch.ProcessInstanceRelation},
}

// watchAPI watch events by api
func (w *Watcher) watchAPI(kit *rest.Kit, resType types.ResType, cursorType watch.CursorType) {
	handler := w.tokenHandlers[resType]

	prevStatus := false

	opt := &watch.WatchEventOptions{
		Resource: cursorType,
	}

	for {
		isMaster := w.isMaster.IsMaster()
		if !isMaster {
			prevStatus = false
			blog.V(4).Infof("watch %s event by api, but not master, skip", resType)
			time.Sleep(time.Minute)
			continue
		}

		// is master status changed, re-watch event from new cursor
		if !prevStatus {
			prevStatus = true
			token, err := handler.getWatchCursorInfo(kit)
			if err != nil {
				blog.Errorf("get %s watch cursor info failed, err: %v, rid: %s", resType, err, kit.Rid)
				time.Sleep(time.Second)
				continue
			}

			opt.Cursor = token.Cursor[cursorType]
			if opt.Cursor == "" && token.StartAtTime != nil {
				opt.StartFrom = token.StartAtTime.Unix()
			}
		}

		// watch events by api and set next watch cursor
		var lastCursor string
		var err error
		util.RetryWrapper(5, func() (bool, error) {
			lastCursor, err = w.doWatchAPI(kit, resType, opt)
			if err != nil {
				blog.Errorf("watch %s events by api failed, err: %v, opt: %+v, rid: %s", resType, err, *opt,
					kit.Rid)
				return true, err
			}
			return false, nil
		})

		opt.Cursor = lastCursor
		opt.StartFrom = 0

		watchTokenInfo := mapstr.MapStr{
			fmt.Sprintf("%s.%s", common.BKCursorField, cursorType): lastCursor,
			common.BKStartAtTimeField:                              &metadata.Time{Time: time.Now()},
		}
		err = handler.setWatchCursorInfo(kit, watchTokenInfo)
		if err != nil {
			blog.Errorf("set %s watch cursor %s failed, err: %v, rid: %s", resType, lastCursor, err, kit.Rid)
			time.Sleep(time.Second)
			continue
		}
	}
}

func (w *Watcher) doWatchAPI(kit *rest.Kit, resType types.ResType, opt *watch.WatchEventOptions) (string, error) {
	// watch events by api
	watchRes, ccErr := w.cacheCli.Cache().Event().InnerWatchEvent(kit.Ctx, kit.Header, opt)
	if ccErr != nil {
		if ccErr.GetCode() == common.CCErrEventChainNodeNotExist {
			blog.Errorf("watch opt(%+v) is invalid, re-watch from now, err: %v, rid: %s", *opt, ccErr, kit.Rid)
			return "", nil
		}
		blog.Errorf("watch event by opt(%+v) failed, err: %v, rid: %s", *opt, ccErr, kit.Rid)
		return "", ccErr
	}

	if len(watchRes.Events) == 0 {
		blog.Errorf("watch response(%+v) has no event, rid: %s", *watchRes, kit.Rid)
		return "", errors.New("watch response has no event")
	}

	if !watchRes.Watched {
		return watchRes.Events[0].Cursor, nil
	}

	// convert event to incremental sync event array
	eventInfos := make([]*types.EventInfo, 0)
	for _, event := range watchRes.Events {
		if event.Detail == nil {
			blog.Errorf("%s watch event(%+v) detail is nil, rid: %s", resType, *event, kit.Rid)
			continue
		}

		// get oid from event cursor
		cursor := new(watch.Cursor)
		if err := cursor.Decode(event.Cursor); err != nil {
			blog.Errorf("%s watch event(%+v) cursor %s is invalid, rid: %s", resType, *event, event.Cursor, kit.Rid)
			continue
		}
		oid := cursor.Oid

		// parse event detail
		var detail string
		switch t := event.Detail.(type) {
		case watch.JsonString:
			detail = string(t)
		default:
			blog.Errorf("%s watch event detail(%+v) is invalid, rid: %s", resType, event.Detail, kit.Rid)
			continue
		}

		eventInfo, needSync := w.metadata.ParseEventDetail(event.EventType, resType, oid, json.RawMessage(detail))
		if !needSync {
			continue
		}
		eventInfos = append(eventInfos, eventInfo)
	}

	// push incremental sync data to transfer medium
	err := w.pushSyncData(kit, eventInfos)
	if err != nil {
		return "", err
	}

	return watchRes.Events[len(watchRes.Events)-1].Cursor, nil
}

type eventDetailMap struct {
	update, create, delete map[string]json.RawMessage
}

func (w *Watcher) pushSyncData(kit *rest.Kit, events []*types.EventInfo) error {
	eventInfoMap := w.classifyEvents(kit, events)

	// push upsert and delete event info to transfer medium
	for resType, subResMap := range eventInfoMap {
		if len(subResMap) == 0 {
			continue
		}

		upsertInfo, deleteInfo := make(map[string][]json.RawMessage), make(map[string][]json.RawMessage)
		for subRes, detailMap := range subResMap {
			if len(detailMap.create) != 0 || len(detailMap.update) != 0 {
				upsertInfo[subRes] = append(convEventDetail(detailMap.create), convEventDetail(detailMap.update)...)
			}
			if len(detailMap.delete) != 0 {
				deleteInfo[subRes] = convEventDetail(detailMap.delete)
			}
		}

		if len(upsertInfo) == 0 && len(deleteInfo) == 0 {
			continue
		}

		pushOpt := &types.PushSyncDataOpt{
			ResType:     resType,
			IsIncrement: true,
			Data: &types.IncrSyncTransData{
				Name:       w.name,
				UpsertInfo: upsertInfo,
				DeleteInfo: deleteInfo,
			},
		}
		err := w.transMedium.PushSyncData(kit.Ctx, kit.Header, pushOpt)
		if err != nil {
			blog.Errorf("push %s incr sync data failed, err: %v, opt: %+v, rid: %s", resType, err, *pushOpt, kit.Rid)
			return err
		}
	}

	return nil
}

// classify events by resource type and sub resources and event type
func (w *Watcher) classifyEvents(kit *rest.Kit, events []*types.EventInfo) map[types.ResType]map[string]eventDetailMap {
	eventInfoMap := make(map[types.ResType]map[string]eventDetailMap)

	for _, event := range events {
		if len(event.SubRes) == 0 {
			event.SubRes = []string{""}
		}

		for _, subRes := range event.SubRes {
			if _, exists := eventInfoMap[event.ResType]; !exists {
				eventInfoMap[event.ResType] = make(map[string]eventDetailMap)
			}
			if _, exists := eventInfoMap[event.ResType][subRes]; !exists {
				eventInfoMap[event.ResType][subRes] = eventDetailMap{
					update: make(map[string]json.RawMessage),
					create: make(map[string]json.RawMessage),
					delete: make(map[string]json.RawMessage),
				}
			}

			switch event.EventType {
			case watch.Create:
				eventInfoMap[event.ResType][subRes].create[event.Oid] = event.Detail
			case watch.Update:
				_, exists := eventInfoMap[event.ResType][subRes].create[event.Oid]
				if exists {
					eventInfoMap[event.ResType][subRes].create[event.Oid] = event.Detail
					continue
				}
				eventInfoMap[event.ResType][subRes].update[event.Oid] = event.Detail
			case watch.Delete:
				// if the data is created then deleted, remove the create and delete event
				_, exists := eventInfoMap[event.ResType][subRes].create[event.Oid]
				if exists {
					delete(eventInfoMap[event.ResType][subRes].create, event.Oid)
					continue
				}

				// data is deleted, remove the update event
				delete(eventInfoMap[event.ResType][subRes].update, event.Oid)

				eventInfoMap[event.ResType][subRes].delete[event.Oid] = event.Detail
			default:
				blog.Errorf("event(%+v) type(%s) is invalid, rid: %s", *event, event.EventType, kit.Rid)
				continue
			}
		}
	}

	return eventInfoMap
}

func convEventDetail(oidDetailMap map[string]json.RawMessage) []json.RawMessage {
	details := make([]json.RawMessage, 0)
	for _, detail := range oidDetailMap {
		details = append(details, detail)
	}
	return details
}
