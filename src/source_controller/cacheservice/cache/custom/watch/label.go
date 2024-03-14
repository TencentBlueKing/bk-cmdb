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

	"configcenter/pkg/conv"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	"configcenter/src/storage/driver/mongodb"
	streamtypes "configcenter/src/storage/stream/types"
)

// watchPodLabel watch pod event for label key and value cache
func (w *Watcher) watchPodLabel() error {
	labelWatcher := &podLabelWatcher{
		labelCache:    w.cacheSet.Label,
		sharedNsCache: w.cacheSet.SharedNsRel,
	}

	opt := &watchOptions{
		watchType: PodLabelWatchType,
		watchOpts: &streamtypes.WatchOptions{
			Options: streamtypes.Options{
				Filter:      make(mapstr.MapStr),
				EventStruct: new(kubetypes.Pod),
				Collection:  kubetypes.BKTableNameBasePod,
				Fields: []string{kubetypes.BKIDField, kubetypes.BKBizIDField, kubetypes.LabelsField,
					kubetypes.BKNamespaceIDField},
			},
		},
		doBatch: labelWatcher.doBatch,
	}

	tokenExists, err := w.watchCustomResource(opt)
	if err != nil {
		return err
	}

	if !tokenExists {
		rid := util.GenerateRID()
		blog.Infof("token not exists, start init all pod label cache task, rid: %s", rid)
		go w.cacheSet.Label.RefreshCache(rid)
	}

	return nil
}

type podLabelWatcher struct {
	labelCache    *cache.PodLabelCache
	sharedNsCache *cache.SharedNsRelCache
}

// doBatch batch handle pod event for label key and value cache
func (w *podLabelWatcher) doBatch(es []*streamtypes.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	rid := es[0].ID()

	// group inserted and deleted pod events
	insertPodMap := make(map[string]*kubetypes.Pod)
	delOids := make([]string, 0)
	nsIDs := make([]int64, 0)

	for idx := range es {
		one := es[idx]

		switch one.OperationType {
		case streamtypes.Insert:
			pod := one.Document.(*kubetypes.Pod)
			insertPodMap[one.Oid] = pod
			nsIDs = append(nsIDs, pod.NamespaceID)
		case streamtypes.Delete:
			_, exists := insertPodMap[one.Oid]
			if exists {
				delete(insertPodMap, one.Oid)
				continue
			}
			delOids = append(delOids, one.Oid)
		default:
			// right now, pod can not be updated, so we only need to handle insert and delete event
			continue
		}

		blog.V(5).Infof("watch custom resource cache, received coll: %s, oid: %s, op-time: %s, %s event, rid: %s",
			one.Collection, one.Oid, one.ClusterTime.String(), one.OperationType, rid)
	}

	delArchives, err := w.getDeletedPodInfo(ctx, delOids, rid)
	if err != nil {
		return true
	}

	for _, archive := range delArchives {
		nsIDs = append(nsIDs, archive.Detail.NamespaceID)
	}

	nsIDs = util.IntArrayUnique(nsIDs)
	asstBizInfo, err := w.sharedNsCache.GetAsstBiz(ctx, nsIDs, rid)
	if err != nil {
		return false
	}

	// get biz to pod label key and value count map
	keyCnt := make(map[int64]map[string]int64)
	valueCnt := make(map[int64]map[string]map[string]int64)

	for _, pod := range insertPodMap {
		w.countPodLabel(pod, asstBizInfo, keyCnt, valueCnt, 1)
	}

	for _, archive := range delArchives {
		w.countPodLabel(archive.Detail, asstBizInfo, keyCnt, valueCnt, -1)
	}

	// update changed pod label key and value cache
	if err = w.labelCache.UpdateKeyCount(ctx, keyCnt, rid); err != nil {
		return true
	}

	if err = w.labelCache.UpdateValueCount(ctx, valueCnt, rid); err != nil {
		return true
	}

	return false
}

// countPodLabel count pod label key and value by biz id
func (w *podLabelWatcher) countPodLabel(pod *kubetypes.Pod, asstBiz map[int64]int64, keyCnt map[int64]map[string]int64,
	valueCnt map[int64]map[string]map[string]int64, cnt int64) {

	if pod == nil || pod.Labels == nil || len(*pod.Labels) == 0 {
		return
	}

	bizs := []int64{pod.BizID}
	asstBizID, exists := asstBiz[pod.NamespaceID]
	if exists {
		bizs = append(bizs, asstBizID)
	}

	for _, bizID := range bizs {
		_, exists = keyCnt[bizID]
		if !exists {
			keyCnt[bizID] = make(map[string]int64)
			valueCnt[bizID] = make(map[string]map[string]int64)
		}

		for key, value := range *pod.Labels {
			key = conv.DecodeDot(key)
			keyCnt[bizID][key] += cnt
			_, exists = valueCnt[bizID][key]
			if !exists {
				valueCnt[bizID][key] = make(map[string]int64)
			}
			valueCnt[bizID][key][value] += cnt
		}
	}
}

type podDelArchive struct {
	Detail *kubetypes.Pod `bson:"detail"`
}

// getDeletedPodInfo get deleted pod info
func (w *podLabelWatcher) getDeletedPodInfo(ctx context.Context, oids []string, rid string) ([]podDelArchive, error) {
	delArchives := make([]podDelArchive, 0)
	if len(oids) == 0 {
		return delArchives, nil
	}

	cond := mapstr.MapStr{
		"oid":  mapstr.MapStr{common.BKDBIN: oids},
		"coll": kubetypes.BKTableNameBasePod,
	}

	err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(cond).Fields("detail.labels", "detail.bk_biz_id",
		"detail.bk_namespace_id").All(ctx, &delArchives)
	if err != nil {
		blog.Errorf("get pod del archive by cond: %+v failed, err: %v, rid: %s", cond, err, rid)
		return nil, err
	}

	return delArchives, nil
}
