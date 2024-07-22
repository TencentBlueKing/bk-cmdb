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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	"configcenter/src/storage/driver/mongodb"
	streamtypes "configcenter/src/storage/stream/types"
)

// watchSharedNsRel watch shared namespace relation event
func (w *Watcher) watchSharedNsRel() error {
	watcher := &sharedNsRelWatcher{
		cache: w.cacheSet.SharedNsRel,
	}

	opt := &watchOptions{
		watchType: SharedNsRelWatchType,
		watchOpts: &streamtypes.WatchOptions{
			Options: streamtypes.Options{
				Filter:      make(mapstr.MapStr),
				EventStruct: new(kubetypes.NsSharedClusterRel),
				Collection:  kubetypes.BKTableNameNsSharedClusterRel,
				Fields:      []string{kubetypes.BKNamespaceIDField, kubetypes.BKAsstBizIDField},
			},
		},
		doBatch: watcher.doBatch,
	}

	tokenExists, err := w.watchCustomResource(opt)
	if err != nil {
		return err
	}

	if !tokenExists {
		ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
		rid := util.GenerateRID()
		blog.Infof("token not exists, start init all shared namespace relation cache task, rid: %s", rid)
		go w.cacheSet.SharedNsRel.RefreshSharedNsRel(ctx, rid)
	}

	return nil
}

type sharedNsRelWatcher struct {
	cache *cache.SharedNsRelCache
}

// doBatch batch handle shared namespace relation event for cache
func (w *sharedNsRelWatcher) doBatch(es []*streamtypes.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	rid := es[0].ID()

	nsAsstBizMap := make(map[int64]int64)
	delOids := make([]string, 0)

	for idx := range es {
		one := es[idx]

		switch one.OperationType {
		case streamtypes.Insert:
			rel := one.Document.(*kubetypes.NsSharedClusterRel)
			nsAsstBizMap[rel.NamespaceID] = rel.AsstBizID
		case streamtypes.Delete:
			delOids = append(delOids, one.Oid)
		default:
			// shared namespace relation can not be updated, so we only need to handle insert and delete event
			continue
		}

		blog.V(5).Infof("watch custom resource cache, received coll: %s, oid: %s, op-time: %s, %s event, rid: %s",
			one.Collection, one.Oid, one.ClusterTime.String(), one.OperationType, rid)
	}

	err := w.cache.UpdateAsstBiz(ctx, nsAsstBizMap, rid)
	if err != nil {
		return true
	}

	delArchives, err := w.getDeletedRelInfo(ctx, delOids, rid)
	if err != nil {
		return true
	}

	delNsIDs := make([]int64, len(delArchives))
	for i, archive := range delArchives {
		delNsIDs[i] = archive.Detail.NamespaceID
	}

	err = w.cache.DeleteAsstBiz(ctx, delNsIDs, rid)
	if err != nil {
		return true
	}

	return false
}

type sharedNsRelDelArchive struct {
	Detail *kubetypes.NsSharedClusterRel `bson:"detail"`
}

// getDeletedRelInfo get deleted shared namespace relation info
func (w *sharedNsRelWatcher) getDeletedRelInfo(ctx context.Context, oids []string, rid string) ([]sharedNsRelDelArchive,
	error) {

	delArchives := make([]sharedNsRelDelArchive, 0)
	if len(oids) == 0 {
		return delArchives, nil
	}

	cond := mapstr.MapStr{
		"oid":  mapstr.MapStr{common.BKDBIN: oids},
		"coll": kubetypes.BKTableNameNsSharedClusterRel,
	}

	err := mongodb.Client().Table(common.BKTableNameKubeDelArchive).Find(cond).Fields("detail.bk_namespace_id").
		All(ctx, &delArchives)
	if err != nil {
		blog.Errorf("get shared ns relation del archive by cond: %+v failed, err: %v, rid: %s", cond, err, rid)
		return nil, err
	}

	return delArchives, nil
}
