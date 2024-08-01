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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	nodelgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/node"
	topolgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/topo"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	streamtypes "configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type kubeWatcher struct {
	watcher       *Watcher
	cacheKey      key.Key
	sharedNsCache *cache.SharedNsRelCache
}

func (w *Watcher) watchKube() error {
	cacheKey, exists := key.TopoKeyMap[types.KubeType]
	if !exists {
		return fmt.Errorf("watch kube but topo type %s is invalid", types.KubeType)
	}

	watcher := &kubeWatcher{
		watcher:       w,
		cacheKey:      cacheKey,
		sharedNsCache: w.cacheSet.SharedNsRel,
	}

	kubeObjs := []string{kubetypes.KubeCluster, kubetypes.KubeNamespace, kubetypes.KubeWorkload}
	for _, obj := range kubeObjs {
		if err := watcher.watchTopo(obj, watcher.onTopoLevelChange(obj)); err != nil {
			blog.Errorf("watch kube topo level %s failed, err: %v", obj, err)
			return err
		}
	}

	if err := watcher.watchTopo(kubetypes.KubePod, watcher.onContainerCountChange); err != nil {
		blog.Errorf("watch kube topo container count failed, err: %v", err)
		return err
	}

	return nil
}

// watchTopoLevel watch kube topo event
func (w *kubeWatcher) watchTopo(obj string, doBatch func(es []*streamtypes.Event) bool) error {
	collections, err := kubetypes.GetCollectionWithObject(obj)
	if err != nil {
		blog.Errorf("get collections to watch for kube biz topo obj %s failed, err: %v", obj, err)
		return err
	}

	for _, collection := range collections {
		watchOpts := &streamtypes.WatchOptions{
			Options: streamtypes.Options{
				EventStruct: kubeEventStructMap[obj],
				Collection:  collection,
				Filter:      mapstr.MapStr{},
			},
		}

		tokenHandler := tokenhandler.NewMixTokenHandler(w.cacheKey.Namespace(), collection, mongodb.Client())
		startAtTime, err := tokenHandler.GetStartWatchTime(context.Background())
		if err != nil {
			blog.Errorf("get start watch time for %s failed, err: %v", watchOpts.Collection, err)
			return err
		}
		watchOpts.StartAtTime = startAtTime
		watchOpts.WatchFatalErrorCallback = tokenHandler.ResetWatchToken

		loopOptions := &streamtypes.LoopBatchOptions{
			LoopOptions: streamtypes.LoopOptions{
				Name:         fmt.Sprintf("%s kube biz topo cache", obj),
				WatchOpt:     watchOpts,
				TokenHandler: tokenHandler,
				RetryOptions: &streamtypes.RetryOptions{
					MaxRetryCount: 10,
					RetryDuration: 1 * time.Second,
				},
			},
			EventHandler: &streamtypes.BatchHandler{
				DoBatch: doBatch,
			},
			BatchSize: 200,
		}

		if err = w.watcher.loopW.WithBatch(loopOptions); err != nil {
			blog.Errorf("watch kube biz topo collection %s failed, err: %v", collection, err)
			return err
		}
	}

	return nil
}

var kubeEventStructMap = map[string]interface{}{
	kubetypes.KubeCluster:   new(kubetypes.Cluster),
	kubetypes.KubeNamespace: new(kubetypes.Namespace),
	kubetypes.KubeWorkload:  new(kubetypes.WorkloadBase),
	kubetypes.KubePod:       new(kubetypes.Pod),
}

// onTopoLevelChange handle kube topo level event
func (w *kubeWatcher) onTopoLevelChange(obj string) func(es []*streamtypes.Event) (retry bool) {
	return func(es []*streamtypes.Event) (retry bool) {
		if len(es) == 0 {
			return false
		}

		ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
		rid := es[0].ID()

		upsertCollOidMap := make(map[string][]primitive.ObjectID)
		delCollOidMap := make(map[string][]string)

		for idx := range es {
			one := es[idx]

			switch one.OperationType {
			case streamtypes.Update:
				if _, exists := one.ChangeDesc.UpdatedFields[kubetypes.KubeNameField]; !exists {
					// only handle topo node name changed events
					continue
				}
				fallthrough
			case streamtypes.Insert:
				oid, err := primitive.ObjectIDFromHex(one.Oid)
				if err != nil {
					blog.Errorf("parse %s event oid %s failed, err: %v, rid: %s", one.Collection, one.Oid, err, rid)
					continue
				}
				upsertCollOidMap[one.Collection] = append(upsertCollOidMap[one.Collection], oid)
			case streamtypes.Delete:
				delCollOidMap[one.Collection] = append(delCollOidMap[one.Collection], one.Oid)
			default:
				continue
			}

			blog.V(5).Infof("watch kube topology cache, received coll: %s, oid: %s, op-time: %s, %s event, rid: %s",
				one.Collection, one.Oid, one.ClusterTime.String(), one.OperationType, rid)
		}

		upsertBizIDMap, err := w.handleUpsertTopoLevelEvent(ctx, obj, upsertCollOidMap, rid)
		if err != nil {
			return true
		}

		delBizIDMap, err := w.handleDeleteTopoLevelEvent(ctx, obj, delCollOidMap, rid)
		if err != nil {
			return true
		}

		bizList := make([]int64, 0)
		for _, one := range es {
			collOidKey := genCollOidKey(one.Collection, one.Oid)
			switch one.OperationType {
			case streamtypes.Update:
				if _, exists := one.ChangeDesc.UpdatedFields[kubetypes.KubeNameField]; !exists {
					continue
				}
				fallthrough
			case streamtypes.Insert:
				bizID, exists := upsertBizIDMap[collOidKey]
				if exists {
					bizList = append(bizList, bizID...)
				}
			case streamtypes.Delete:
				bizID, exists := delBizIDMap[collOidKey]
				if exists {
					bizList = append(bizList, bizID...)
				}
			default:
				continue
			}
		}

		topolgc.AddRefreshBizTopoTask(types.KubeType, util.IntArrayUnique(bizList), rid)
		return false
	}
}

func genCollOidKey(coll, oid string) string {
	return coll + "-" + oid
}

type mapStrWithOid struct {
	Oid    primitive.ObjectID     `bson:"_id"`
	BizID  int64                  `bson:"bk_biz_id"`
	ID     int64                  `bson:"id"`
	Name   string                 `bson:"name"`
	MapStr map[string]interface{} `bson:",inline"`
}

var kubeFieldsMap = map[string][]string{
	kubetypes.KubeCluster: {kubetypes.BKBizIDField, kubetypes.BKIDField, kubetypes.KubeNameField},
	kubetypes.KubeNamespace: {kubetypes.BKBizIDField, kubetypes.BKIDField, kubetypes.KubeNameField,
		kubetypes.BKClusterIDFiled},
	kubetypes.KubeWorkload: {kubetypes.BKBizIDField, kubetypes.BKIDField, kubetypes.KubeNameField,
		kubetypes.BKNamespaceIDField},
}

// handleUpsertTopoLevelEvent handle upsert event for kube topo level
func (w *kubeWatcher) handleUpsertTopoLevelEvent(ctx context.Context, obj string,
	collOidMap map[string][]primitive.ObjectID, rid string) (map[string][]int64, error) {

	collOidBizMap := make(map[string][]int64)
	for coll, oids := range collOidMap {
		// get upsert data from db
		cond := mapstr.MapStr{
			"_id": mapstr.MapStr{common.BKDBIN: oids},
		}

		docs := make([]mapStrWithOid, 0)

		findOpt := dbtypes.NewFindOpts().SetWithObjectID(true)
		err := mongodb.Client().Table(coll).Find(cond, findOpt).Fields(kubeFieldsMap[obj]...).All(ctx, &docs)
		if err != nil {
			blog.Errorf("get %s data by cond: %+v failed, err: %v, rid: %s", coll, cond, err, rid)
			return nil, err
		}

		kind, err := getKubeNodeKind(obj, coll)
		if err != nil {
			blog.Errorf("get %s kube node kind by coll %s failed, err: %v, rid: %s", obj, coll, err, rid)
			continue
		}

		bizNodeMap := make(map[int64][]types.Node)
		nsIDs := make([]int64, 0)
		nsNodeMap := make(map[int64][]types.Node)
		nsKeyMap := make(map[int64][]string)
		for _, doc := range docs {
			// parse event to biz id and topo level node
			bizID, node, err := kubeEventDocParserMap[obj](doc)
			if err != nil {
				blog.Errorf("parse %s doc %+v failed, err: %v, rid: %s", coll, doc, err, rid)
				continue
			}

			var nsID int64
			switch obj {
			case kubetypes.KubeNamespace:
				nsID = node.ID
			case kubetypes.KubeWorkload:
				nsID = node.ParentID
			}
			nsIDs = append(nsIDs, nsID)

			node.Kind = kind
			bizNodeMap[bizID] = append(bizNodeMap[bizID], node)
			nsNodeMap[nsID] = append(nsNodeMap[nsID], node)

			collOidKey := genCollOidKey(coll, doc.Oid.Hex())
			collOidBizMap[collOidKey] = []int64{bizID}
			nsKeyMap[nsID] = append(nsKeyMap[nsID], collOidKey)
		}

		// add shared namespace nodes to asst biz's cache
		asstBizInfo, err := w.sharedNsCache.GetAsstBiz(ctx, nsIDs, rid)
		if err != nil {
			return nil, err
		}

		for nsID, nodes := range nsNodeMap {
			asstBizID, exists := asstBizInfo[nsID]
			if !exists {
				continue
			}
			bizNodeMap[asstBizID] = append(bizNodeMap[asstBizID], nodes...)
			for _, collOidKey := range nsKeyMap[nsID] {
				collOidBizMap[collOidKey] = append(collOidBizMap[collOidKey], asstBizID)
			}
		}

		for bizID, nodes := range bizNodeMap {
			// save kube topo level node info to redis
			err = nodelgc.AddNodeInfoCache(w.cacheKey, bizID, kind, nodes, rid)
			if err != nil {
				return nil, err
			}
		}
	}

	return collOidBizMap, nil
}

type commonDelArchive struct {
	Oid    string            `bson:"oid"`
	Detail commonResBaseInfo `bson:"detail"`
}

type commonResBaseInfo struct {
	BizID int64 `bson:"bk_biz_id"`
	ID    int64 `bson:"id"`
	NsID  int64 `bson:"bk_namespace_id"`
}

// handleDeleteTopoLevelEvent handle delete event for kube topo level
func (w *kubeWatcher) handleDeleteTopoLevelEvent(ctx context.Context, obj string, collOidMap map[string][]string,
	rid string) (map[string][]int64, error) {

	collOidBizMap := make(map[string][]int64)
	for coll, oids := range collOidMap {
		// get del archive data
		cond := mapstr.MapStr{
			"oid":  mapstr.MapStr{common.BKDBIN: oids},
			"coll": coll,
		}

		docs := make([]commonDelArchive, 0)
		err := mongodb.Client().Table(common.BKTableNameKubeDelArchive).Find(cond).All(ctx, &docs)
		if err != nil {
			blog.Errorf("get del archive by cond: %+v failed, err: %v, rid: %s", cond, err, rid)
			return nil, err
		}

		kind, err := getKubeNodeKind(obj, coll)
		if err != nil {
			blog.Errorf("get %s kube node kind by coll %s failed, err: %v, rid: %s", obj, coll, err, rid)
			continue
		}

		bizIDMap := make(map[int64][]int64)
		nsIDs := make([]int64, 0)
		nsDocIDsMap := make(map[int64][]int64)
		nsKeyMap := make(map[int64][]string)
		for _, doc := range docs {
			var nsID int64
			switch obj {
			case kubetypes.KubeNamespace:
				nsID = doc.Detail.ID
			case kubetypes.KubeWorkload:
				nsID = doc.Detail.NsID
			}
			nsIDs = append(nsIDs, nsID)

			bizIDMap[doc.Detail.BizID] = append(bizIDMap[doc.Detail.BizID], doc.Detail.ID)
			nsDocIDsMap[nsID] = append(nsDocIDsMap[nsID], doc.Detail.ID)

			collOidKey := genCollOidKey(coll, doc.Oid)
			collOidBizMap[collOidKey] = []int64{doc.Detail.BizID}
			nsKeyMap[nsID] = append(nsKeyMap[nsID], collOidKey)
		}

		// delete shared namespace node info in asst biz's cache
		asstBizInfo, err := w.sharedNsCache.GetAsstBiz(ctx, nsIDs, rid)
		if err != nil {
			return nil, err
		}

		for nsID, docIDs := range nsDocIDsMap {
			asstBizID, exists := asstBizInfo[nsID]
			if !exists {
				continue
			}
			bizIDMap[asstBizID] = append(bizIDMap[asstBizID], docIDs...)
			for _, collOidKey := range nsKeyMap[nsID] {
				collOidBizMap[collOidKey] = append(collOidBizMap[collOidKey], asstBizID)
			}
		}

		// delete kube topo level node info in redis
		for bizID, ids := range bizIDMap {
			err = nodelgc.DeleteNodeInfoCache(w.cacheKey, bizID, kind, ids, rid)
			if err != nil {
				return nil, err
			}

			err = nodelgc.DeleteNodeCountCache(w.cacheKey, bizID, kind, ids, rid)
			if err != nil {
				return nil, err
			}
		}
	}

	return collOidBizMap, nil
}

func getKubeNodeKind(obj string, coll string) (string, error) {
	switch obj {
	case kubetypes.KubeCluster, kubetypes.KubeNamespace:
		return obj, nil
	case kubetypes.KubeWorkload:
		kind, err := kubetypes.GetKindByWorkLoadTableName(coll)
		if err != nil {
			return "", fmt.Errorf("get kube workload kind by coll %s failed, err: %v", coll, err)
		}
		return string(kind), nil
	}
	return "", fmt.Errorf("kube node obj %s is invalid", obj)
}

// kubeEventDocParserMap is the mapping of kube level to kube event document parser
var kubeEventDocParserMap = map[string]func(doc mapStrWithOid) (int64, types.Node, error){
	kubetypes.KubeCluster: func(doc mapStrWithOid) (int64, types.Node, error) {
		node := types.Node{
			ID:   doc.ID,
			Name: doc.Name,
		}
		return doc.BizID, node, nil
	},
	kubetypes.KubeNamespace: func(doc mapStrWithOid) (int64, types.Node, error) {
		parentID, err := util.GetInt64ByInterface(doc.MapStr[kubetypes.BKClusterIDFiled])
		if err != nil {
			return 0, types.Node{}, fmt.Errorf("parse kube workload parent id failed, err: %v", err)
		}

		node := types.Node{
			ID:       doc.ID,
			Name:     doc.Name,
			ParentID: parentID,
		}
		return doc.BizID, node, nil
	},
	kubetypes.KubeWorkload: func(doc mapStrWithOid) (int64, types.Node, error) {
		parentID, err := util.GetInt64ByInterface(doc.MapStr[kubetypes.BKNamespaceIDField])
		if err != nil {
			return 0, types.Node{}, fmt.Errorf("parse kube workload parent id failed, err: %v", err)
		}

		node := types.Node{
			ID:       doc.ID,
			Name:     doc.Name,
			ParentID: parentID,
		}
		return doc.BizID, node, nil
	},
}

// onContainerCountChange handle container count change event
func (w *kubeWatcher) onContainerCountChange(es []*streamtypes.Event) (retry bool) {
	if len(es) == 0 {
		return false
	}

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	rid := es[0].ID()

	insertPodNsIDs := make([]int64, 0)
	delOids := make([]string, 0)

	for idx := range es {
		one := es[idx]

		switch one.OperationType {
		case streamtypes.Insert:
			pod, ok := one.Document.(*kubetypes.Pod)
			if !ok {
				blog.Errorf("kube pod event %s doc type %T is invalid, rid: %s", one.Oid, one.Document, rid)
				continue
			}
			insertPodNsIDs = append(insertPodNsIDs, pod.NamespaceID)
		case streamtypes.Delete:
			delOids = append(delOids, one.Oid)
		default:
			// only handle insert and delete pod event.
			continue
		}

		blog.V(5).Infof("watch kube topo container count cache, received oid: %s, op-time: %s, %s event, rid: %s",
			one.Oid, one.ClusterTime.String(), one.OperationType, rid)
	}

	delOidPodMap, delPodNsIDs, err := w.getDeletePodInfo(ctx, delOids, rid)
	if err != nil {
		return true
	}

	nsIDs := util.IntArrayUnique(append(insertPodNsIDs, delPodNsIDs...))
	asstBizInfo, err := w.sharedNsCache.GetAsstBiz(ctx, nsIDs, rid)
	if err != nil {
		return true
	}

	bizList, err := w.handlePodEvents(ctx, es, delOidPodMap, asstBizInfo, rid)
	if err != nil {
		return true
	}

	topolgc.AddRefreshBizTopoTask(types.KubeType, bizList, rid)
	return false
}

type podDelArchive struct {
	Oid    string         `bson:"oid"`
	Detail *kubetypes.Pod `bson:"detail"`
}

// getDeletePodInfo get delete pod info from del archive
func (w *kubeWatcher) getDeletePodInfo(ctx context.Context, oids []string, rid string) (map[string]*kubetypes.Pod,
	[]int64, error) {

	cond := mapstr.MapStr{
		"oid":  mapstr.MapStr{common.BKDBIN: oids},
		"coll": kubetypes.BKTableNameBasePod,
	}

	archives := make([]podDelArchive, 0)
	err := mongodb.Client().Table(common.BKTableNameKubeDelArchive).Find(cond).Fields("oid", "detail").All(ctx,
		&archives)
	if err != nil {
		blog.Errorf("get pod del archive by cond: %+v failed, err: %v, rid: %s", cond, err, rid)
		return nil, nil, err
	}

	podMap := make(map[string]*kubetypes.Pod)
	nsIDs := make([]int64, 0)
	for _, archive := range archives {
		podMap[archive.Oid] = archive.Detail
		nsIDs = append(nsIDs, archive.Detail.NamespaceID)
	}

	return podMap, nsIDs, nil
}

// handlePodEvents refresh pod events related workload container count, returns biz ids whose topo tree needs refreshing
func (w *kubeWatcher) handlePodEvents(ctx context.Context, es []*streamtypes.Event,
	delOidPodMap map[string]*kubetypes.Pod, asstBizInfo map[int64]int64, rid string) ([]int64, error) {

	// wlKindIDMap is map[workload_kind][]workload_id, stores the pod events related workload info
	wlKindIDMap := make(map[kubetypes.WorkloadType][]int64)
	// wlBizIDMap is map[workload_kind]map[workload_id][]bk_biz_id, stores the workload to its related biz ids
	wlBizIDMap := make(map[kubetypes.WorkloadType]map[int64][]int64)
	// bizList is biz ids whose topo tree needs refreshing, in the order of pod events
	bizList := make([]int64, 0)
	for idx := range es {
		one := es[idx]

		var pod *kubetypes.Pod
		switch one.OperationType {
		case streamtypes.Insert:
			pod = one.Document.(*kubetypes.Pod)
		case streamtypes.Delete:
			var exists bool
			pod, exists = delOidPodMap[one.Oid]
			if !exists {
				continue
			}
		default:
			continue
		}

		// record the workloads that needs to refresh container count
		wlKindIDMap[pod.Ref.Kind] = append(wlKindIDMap[pod.Ref.Kind], pod.Ref.ID)

		// record workload related biz info, including the pod's biz id and shared namespace asst biz id
		_, exists := wlBizIDMap[pod.Ref.Kind]
		if !exists {
			wlBizIDMap[pod.Ref.Kind] = make(map[int64][]int64)
		}
		wlBizIDMap[pod.Ref.Kind][pod.Ref.ID] = []int64{pod.BizID}
		bizList = append(bizList, pod.BizID)

		asstBizID, exists := asstBizInfo[pod.NamespaceID]
		if exists {
			wlBizIDMap[pod.Ref.Kind][pod.Ref.ID] = append(wlBizIDMap[pod.Ref.Kind][pod.Ref.ID], asstBizID)
			bizList = append(bizList, asstBizID)
		}
	}

	// refresh workload topo node container count cache by workload kind
	for wlType, wlIDs := range wlKindIDMap {
		wlIDs = util.IntArrayUnique(wlIDs)
		if err := w.refreshWlCountCache(ctx, string(wlType), wlIDs, wlBizIDMap[wlType], rid); err != nil {
			return nil, err
		}
	}

	return bizList, nil
}

// refreshWlCountCache refresh workload topo node container count cache
func (w *kubeWatcher) refreshWlCountCache(ctx context.Context, kind string, ids []int64, bizMap map[int64][]int64,
	rid string) error {

	// get pods by pod workloads
	podCond := mapstr.MapStr{
		kubetypes.RefIDField:   mapstr.MapStr{common.BKDBIN: ids},
		kubetypes.RefKindField: kind,
	}

	pods := make([]kubetypes.Pod, 0)
	if err := mongodb.Client().Table(kubetypes.BKTableNameBasePod).Find(podCond).Fields(kubetypes.BKIDField,
		kubetypes.RefIDField).All(ctx, &pods); err != nil {
		blog.Errorf("get pod ids failed, cond: %+v, err: %v, rid: %s", podCond, err, rid)
		return err
	}

	wlPodIDMap := make(map[int64][]int64)
	for _, pod := range pods {
		wlPodIDMap[pod.Ref.ID] = append(wlPodIDMap[pod.Ref.ID], pod.ID)
	}

	for wlID, podIDs := range wlPodIDMap {
		containerCnt := int64(0)
		if len(pods) > 0 {
			// get container counts by pods
			containerCond := mapstr.MapStr{
				kubetypes.BKPodIDField: mapstr.MapStr{common.BKDBIN: podIDs},
			}

			cnt, err := mongodb.Client().Table(kubetypes.BKTableNameBaseContainer).Find(containerCond).Count(ctx)
			if err != nil {
				blog.Errorf("count containers failed, cond: %+v, err: %v, rid: %s", containerCond, err, rid)
				return err
			}
			containerCnt = int64(cnt)
		}

		// refresh workload topo node count cache in related bizs
		for _, bizID := range bizMap[wlID] {
			err := nodelgc.AddNodeCountCache(w.cacheKey, bizID, kind, map[int64]int64{wlID: containerCnt}, rid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
