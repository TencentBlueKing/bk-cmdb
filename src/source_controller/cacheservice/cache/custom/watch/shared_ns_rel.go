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
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	streamtypes "configcenter/src/storage/stream/types"
)

// watchSharedNsRel watch shared namespace relation event
func (w *Watcher) watchSharedNsRel() error {
	watcher := &sharedNsRelWatcher{
		cache: w.cacheSet.SharedNsRel,
	}

	opt := &watchOptions{
		watchType: SharedNsRelWatchType,
		watchOpts: &streamtypes.WatchCollOptions{
			CollectionOptions: streamtypes.CollectionOptions{
				EventStruct: new(kubetypes.NsSharedClusterRel),
				CollectionFilter: &streamtypes.CollectionFilter{
					Regex: fmt.Sprintf("_%s$", kubetypes.BKTableNameNsSharedClusterRel),
				},
				Fields: []string{kubetypes.BKNamespaceIDField, kubetypes.BKAsstBizIDField},
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
		blog.Infof("token not exists, start init all shared namespace relation cache task, rid: %s", rid)
		go w.cacheSet.SharedNsRel.RefreshSharedNsRel(rid)
	}

	return nil
}

type sharedNsRelWatcher struct {
	cache *cache.SharedNsRelCache
}

// doBatch batch handle shared namespace relation event for cache
func (w *sharedNsRelWatcher) doBatch(dbInfo *streamtypes.DBInfo, es []*streamtypes.Event) bool {
	if len(es) == 0 {
		return false
	}

	kit := rest.NewKit().WithCtx(util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)).
		WithRid(es[0].ID())

	nsAsstBizMap := make(map[string]map[int64]int64)
	delNsIDsMap := make(map[string][]int64)

	for idx := range es {
		one := es[idx]

		tenantID := one.TenantID
		rel := one.Document.(*kubetypes.NsSharedClusterRel)

		switch one.OperationType {
		case streamtypes.Insert:
			_, exists := nsAsstBizMap[tenantID]
			if !exists {
				nsAsstBizMap[tenantID] = make(map[int64]int64)
			}
			nsAsstBizMap[tenantID][rel.NamespaceID] = rel.AsstBizID

		case streamtypes.Delete:
			delNsIDsMap[tenantID] = append(delNsIDsMap[tenantID], rel.NamespaceID)

		default:
			// shared namespace relation can not be updated, so we only need to handle insert and delete event
			continue
		}

		blog.V(5).Infof("watch custom resource cache, received coll: %s, oid: %s, op-time: %s, %s event, rid: %s",
			one.Collection, one.Oid, one.ClusterTime.String(), one.OperationType, kit.Rid)
	}

	for tenantID, nsAsstBizInfo := range nsAsstBizMap {
		kit = kit.WithTenant(tenantID)
		err := w.cache.UpdateAsstBiz(kit, nsAsstBizInfo)
		if err != nil {
			return true
		}
	}

	for tenantID, delNsIDs := range delNsIDsMap {
		kit = kit.WithTenant(tenantID)
		err := w.cache.DeleteAsstBiz(kit, delNsIDs)
		if err != nil {
			return true
		}
	}

	return false
}
