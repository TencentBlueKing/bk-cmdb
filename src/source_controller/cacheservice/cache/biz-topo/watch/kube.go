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
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	streamtypes "configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type kubeWatcher struct {
	watcher  *Watcher
	cacheKey key.Key
}

func (w *Watcher) watchKube() error {
	cacheKey, exists := key.TopoKeyMap[types.KubeType]
	if !exists {
		return fmt.Errorf("watch kube but topo type %s is invalid", types.KubeType)
	}

	watcher := &kubeWatcher{
		watcher:  w,
		cacheKey: cacheKey,
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
					bizList = append(bizList, bizID)
				}
			case streamtypes.Delete:
				bizID, exists := delBizIDMap[collOidKey]
				if exists {
					bizList = append(bizList, bizID)
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
	collOidMap map[string][]primitive.ObjectID, rid string) (map[string]int64, error) {

	collOidBizMap := make(map[string]int64)
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
		for _, doc := range docs {
			// parse event to biz id and topo level node
			bizID, node, err := kubeEventDocParserMap[obj](doc)
			if err != nil {
				blog.Errorf("parse %s doc %+v failed, err: %v, rid: %s", coll, doc, err, rid)
				continue
			}
			node.Kind = kind
			bizNodeMap[bizID] = append(bizNodeMap[bizID], node)
			collOidBizMap[genCollOidKey(coll, doc.Oid.Hex())] = bizID
		}

		// delete kube topo level node info in redis
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
}

// handleDeleteTopoLevelEvent handle delete event for kube topo level
func (w *kubeWatcher) handleDeleteTopoLevelEvent(ctx context.Context, obj string, collOidMap map[string][]string,
	rid string) (map[string]int64, error) {

	collOidBizMap := make(map[string]int64)
	for coll, oids := range collOidMap {
		// get del archive data
		cond := mapstr.MapStr{
			"oid":  mapstr.MapStr{common.BKDBIN: oids},
			"coll": coll,
		}

		docs := make([]commonDelArchive, 0)
		err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(cond).All(ctx, &docs)
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
		for _, doc := range docs {
			bizIDMap[doc.Detail.BizID] = append(bizIDMap[doc.Detail.BizID], doc.Detail.ID)
			collOidBizMap[genCollOidKey(coll, doc.Oid)] = doc.Detail.BizID
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
	bizList := make([]int64, 0)
	for idx := range es {
		event := es[idx]

		var bizID int64
		var skip bool
		var err error
		switch event.OperationType {
		case streamtypes.Insert:
			bizID, skip, err = w.handleInsertPodEvent(ctx, event, rid)
		case streamtypes.Delete:
			bizID, skip, err = w.handleDeletePodEvent(ctx, event, rid)
		default:
			// only handle insert and delete event.
			continue
		}

		if err != nil {
			return true
		}

		if skip {
			continue
		}

		blog.V(5).Infof("watch kube topo container count cache, received biz: %d, oid: %s, op-time: %s, %s event, "+
			"rid: %s", bizID, event.Oid, event.ClusterTime.String(), event.OperationType, rid)

		bizList = append(bizList, bizID)
	}

	topolgc.AddRefreshBizTopoTask(types.KubeType, bizList, rid)
	return false
}

type podDelArchive struct {
	Detail *kubetypes.Pod `bson:"detail"`
}

// handleInsertPodEvent handle insert pod event for kube topo container count
func (w *kubeWatcher) handleInsertPodEvent(ctx context.Context, event *streamtypes.Event, rid string) (int64, bool,
	error) {

	pod, ok := event.Document.(*kubetypes.Pod)
	if !ok {
		blog.Errorf("kube pod event %s doc type %T is invalid, rid: %s", event.Oid, event.Document, rid)
		return 0, true, nil
	}

	if err := w.refreshPodWlCountCache(ctx, pod, rid); err != nil {
		return 0, false, err
	}

	return pod.BizID, false, nil
}

// handleDeletePodEvent handle delete pod event for kube topo container count
func (w *kubeWatcher) handleDeletePodEvent(ctx context.Context, event *streamtypes.Event, rid string) (int64, bool,
	error) {

	cond := mapstr.MapStr{
		"oid":  event.Oid,
		"coll": event.Collection,
	}

	archive := new(podDelArchive)
	err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(cond).Fields("detail").One(ctx, archive)
	if err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			blog.Errorf("can not find %s pod del archive, skip, rid: %s", event.Oid, rid)
			return 0, true, nil
		}

		blog.Errorf("get pod del archive by cond: %+v failed, err: %v, rid: %s", cond, err, rid)
		return 0, false, err
	}

	pod := archive.Detail

	if err = w.refreshPodWlCountCache(ctx, pod, rid); err != nil {
		return 0, false, err
	}

	return pod.BizID, false, nil
}

// refreshPodWlCountCache refresh pod related workload topo node container count cache
func (w *kubeWatcher) refreshPodWlCountCache(ctx context.Context, pod *kubetypes.Pod, rid string) error {
	kind := string(pod.Ref.Kind)
	wlID := pod.Ref.ID

	// get pods by pod workloads
	podCond := mapstr.MapStr{
		kubetypes.RefIDField:   wlID,
		kubetypes.RefKindField: kind,
	}

	pods := make([]kubetypes.Pod, 0)
	if err := mongodb.Client().Table(kubetypes.BKTableNameBasePod).Find(podCond).Fields(kubetypes.BKIDField).
		All(ctx, &pods); err != nil {
		blog.Errorf("get pod ids failed, cond: %+v, err: %v, rid: %s", podCond, err, rid)
		return err
	}

	if len(pods) == 0 {
		return nodelgc.AddNodeCountCache(w.cacheKey, pod.BizID, kind, map[int64]int64{wlID: 0}, rid)
	}

	podIDs := make([]int64, len(pods))
	for i, p := range pods {
		podIDs[i] = p.ID
	}

	// get container counts by pods
	containerCond := mapstr.MapStr{
		kubetypes.BKPodIDField: mapstr.MapStr{common.BKDBIN: podIDs},
	}

	cnt, err := mongodb.Client().Table(kubetypes.BKTableNameBaseContainer).Find(containerCond).Count(ctx)
	if err != nil {
		blog.Errorf("count containers failed, cond: %+v, err: %v, rid: %s", containerCond, err, rid)
		return err
	}

	return nodelgc.AddNodeCountCache(w.cacheKey, pod.BizID, kind, map[int64]int64{wlID: int64(cnt)}, rid)
}
