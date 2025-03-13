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

	"configcenter/pkg/conv"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
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
		watchOpts: &streamtypes.WatchCollOptions{
			CollectionOptions: streamtypes.CollectionOptions{
				EventStruct: new(kubetypes.Pod),
				CollectionFilter: &streamtypes.CollectionFilter{
					Regex: fmt.Sprintf("_%s$", kubetypes.BKTableNameBasePod),
				},
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
func (w *podLabelWatcher) doBatch(dbInfo *streamtypes.DBInfo, es []*streamtypes.Event) bool {
	if len(es) == 0 {
		return false
	}

	kit := rest.NewKit().WithCtx(util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)).
		WithRid(es[0].ID())

	// group inserted and deleted pod events
	insertPodMap := make(map[string]map[string]*kubetypes.Pod)
	delPodMap := make(map[string]map[string]*kubetypes.Pod)
	nsIDMap := make(map[string][]int64)

	for idx := range es {
		one := es[idx]

		tenantID := one.TenantID
		pod := one.Document.(*kubetypes.Pod)

		switch one.OperationType {
		case streamtypes.Insert:
			_, exists := insertPodMap[tenantID]
			if !exists {
				insertPodMap[tenantID] = make(map[string]*kubetypes.Pod)
			}
			insertPodMap[tenantID][one.Oid] = pod
			nsIDMap[tenantID] = append(nsIDMap[tenantID], pod.NamespaceID)
		case streamtypes.Delete:
			_, exists := insertPodMap[one.Oid]
			if exists {
				delete(insertPodMap, one.Oid)
				continue
			}
			_, exists = delPodMap[tenantID]
			if !exists {
				delPodMap[tenantID] = make(map[string]*kubetypes.Pod)
			}
			delPodMap[tenantID][one.Oid] = pod
			nsIDMap[tenantID] = append(nsIDMap[tenantID], pod.NamespaceID)
		default:
			// right now, pod can not be updated, so we only need to handle insert and delete event
			continue
		}

		blog.V(5).Infof("watch custom resource cache, received coll: %s, oid: %s, op-time: %s, %s event, rid: %s",
			one.Collection, one.Oid, one.ClusterTime.String(), one.OperationType, kit.Rid)
	}

	for tenantID, nsIDs := range nsIDMap {
		nsIDs = util.IntArrayUnique(nsIDs)
		asstBiz, err := w.sharedNsCache.GetAsstBiz(kit, nsIDs)
		if err != nil {
			return false
		}

		// get biz to pod label key and value count map
		keyCnt := make(map[int64]map[string]int64)
		valueCnt := make(map[int64]map[string]map[string]int64)

		for _, pod := range insertPodMap[tenantID] {
			w.countPodLabel(pod, asstBiz, keyCnt, valueCnt, 1)
		}
		for _, pod := range delPodMap[tenantID] {
			w.countPodLabel(pod, asstBiz, keyCnt, valueCnt, -1)
		}

		kit = kit.WithTenant(tenantID)

		// update changed pod label key and value cache
		if err := w.labelCache.UpdateKeyCount(kit, keyCnt); err != nil {
			return true
		}

		if err := w.labelCache.UpdateValueCount(kit, valueCnt); err != nil {
			return true
		}
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
