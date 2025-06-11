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
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache/object"
	streamtypes "configcenter/src/storage/stream/types"
)

// watchObject watch object event
func (w *Watcher) watchObject() error {
	watcher := &objectWatcher{
		cache: w.cacheSet.Object,
	}

	opt := &watchOptions{
		watchType: ObjectWatchType,
		watchOpts: &streamtypes.WatchCollOptions{
			OperationType: []streamtypes.OperType{streamtypes.Insert, streamtypes.Delete},
			CollectionOptions: streamtypes.CollectionOptions{
				EventStruct: new(metadata.Object),
				CollectionFilter: &streamtypes.CollectionFilter{
					Regex: fmt.Sprintf("_%s$", common.BKTableNameObjDes),
				},
				Fields: []string{common.BKObjIDField, metadata.ModelFieldObjUUID},
			},
		},
		doBatch: watcher.doBatch,
	}

	tokenExists, err := w.watchCustomResource(opt)
	if err != nil {
		return err
	}

	if !tokenExists {
		rid := util.GenerateRID()
		blog.Infof("token not exists, start init all object cache task, rid: %s", rid)
		go w.cacheSet.Object.RefreshCache(rid)
	}

	return nil
}

type objectWatcher struct {
	cache *object.ObjectCache
}

// doBatch batch handle object event for cache
func (w *objectWatcher) doBatch(dbInfo *streamtypes.DBInfo, es []*streamtypes.Event) bool {
	if len(es) == 0 {
		return false
	}

	kit := rest.NewKit().WithCtx(util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)).
		WithRid(es[0].ID())

	objUUIDMap := make(map[string]map[string]string)
	delObjIDsMap := make(map[string][]string)

	for idx := range es {
		one := es[idx]

		tenantID := one.TenantID
		obj := one.Document.(*metadata.Object)

		switch one.OperationType {
		case streamtypes.Insert:
			_, exists := objUUIDMap[tenantID]
			if !exists {
				objUUIDMap[tenantID] = make(map[string]string)
			}
			objUUIDMap[tenantID][obj.ObjectID] = obj.UUID

		case streamtypes.Delete:
			objMap, exists := objUUIDMap[tenantID]
			if exists {
				_, ok := objMap[obj.ObjectID]
				if ok {
					delete(objMap, obj.ObjectID)
				}
				continue
			}

			delObjIDsMap[tenantID] = append(delObjIDsMap[tenantID], obj.ObjectID)

		default:
			// object uuid can not be updated, so we only need to handle insert and delete event
			continue
		}

		blog.V(5).Infof("watch custom resource cache, received coll: %s, oid: %s, op-time: %s, %s event, rid: %s",
			one.Collection, one.Oid, one.ClusterTime.String(), one.OperationType, kit.Rid)
	}

	for tenantID, objMap := range objUUIDMap {
		kit = kit.WithTenant(tenantID)
		err := w.cache.UpdateObjUUIDCache(kit, objMap)
		if err != nil {
			return true
		}
	}

	for tenantID, delObjIDs := range delObjIDsMap {
		kit = kit.WithTenant(tenantID)
		err := w.cache.DeleteObjUUIDCache(kit, delObjIDs)
		if err != nil {
			return true
		}
	}

	return false
}
