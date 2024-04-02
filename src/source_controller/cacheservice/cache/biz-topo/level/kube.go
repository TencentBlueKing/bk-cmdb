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

package level

import (
	"context"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	nlgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/node"
	topotypes "configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	"configcenter/src/source_controller/cacheservice/cache/tools"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

var kubeTopLevel = newKubeClusterLevel()

// GetKubeTopLevel get the top level of kube biz topology
func GetKubeTopLevel() LevelI {
	return kubeTopLevel
}

type kubeClusterLevel struct {
	*commonCacheLevel
}

func newKubeClusterLevel() *kubeClusterLevel {
	return &kubeClusterLevel{
		commonCacheLevel: newCommonCacheLevel(topotypes.KubeType, newKubeNsLevel(), types.KubeCluster),
	}
}

// GetNodesByDB get all nodes that belongs to the topology level
func (l *kubeClusterLevel) GetNodesByDB(ctx context.Context, bizID int64, _ []mapstr.MapStr, rid string) (
	[]topotypes.Node, error) {

	sharedCond := mapstr.MapStr{types.BKBizIDField: bizID}

	relations := make([]types.NsSharedClusterRel, 0)
	err := mongodb.Client().Table(types.BKTableNameNsSharedClusterRel).Find(sharedCond).
		Fields(types.BKClusterIDFiled).All(ctx, &relations)
	if err != nil {
		blog.Errorf("list kube shared cluster rel failed, err: %v, cond: %+v, rid: %v", err, sharedCond, rid)
		return nil, err
	}

	sharedClusterIDs := make([]int64, 0)
	for _, relation := range relations {
		sharedClusterIDs = append(sharedClusterIDs, relation.ClusterID)
	}

	cond := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{types.BKBizIDField: bizID},
			{types.BKIDField: mapstr.MapStr{common.BKDBIN: sharedClusterIDs}},
		},
	}

	fields := []string{types.BKIDField, types.KubeNameField}
	nodes, err := nlgc.PagedGetNodes(ctx, types.BKTableNameBaseCluster, cond, fields, l.clusterNodeParser(), rid)
	if err != nil {
		return nil, err
	}

	err = nlgc.CrossCompareNodeInfoCache(l.topoKey, bizID, types.KubeCluster, nodes, rid)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nodes, nil
	}

	childNodes, err := l.nextLevel.GetNodesByDB(ctx, bizID, nil, rid)
	if err != nil {
		return nil, err
	}

	nodes = nlgc.CombineChildNodes(nodes, childNodes)
	return nodes, nil
}

func (l *kubeClusterLevel) clusterNodeParser() nlgc.NodeParser {
	return func(ctx context.Context, data []mapstr.MapStr, rid string) ([]topotypes.Node, error) {
		nodes := make([]topotypes.Node, len(data))
		for i, item := range data {
			id, err := util.GetInt64ByInterface(item[types.BKIDField])
			if err != nil {
				blog.Errorf("parse kube cluster id failed, err: %v, item: %+v, rid: %s", err, item, rid)
				return nil, err
			}

			nodes[i] = topotypes.Node{
				Kind: types.KubeCluster,
				ID:   id,
				Name: util.GetStrByInterface(item[types.KubeNameField]),
			}
		}

		return nodes, nil
	}
}

type kubeNsLevel struct {
	*commonCacheLevel
}

func newKubeNsLevel() *kubeNsLevel {
	return &kubeNsLevel{
		commonCacheLevel: newCommonCacheLevel(topotypes.KubeType, newKubeWlLevel(), types.KubeNamespace),
	}
}

// GetNodesByDB get all nodes that belongs to the topology level
func (l *kubeNsLevel) GetNodesByDB(ctx context.Context, bizID int64, _ []mapstr.MapStr, rid string) ([]topotypes.Node,
	error) {

	cond, err := tools.GenKubeSharedNsCond(ctx, bizID, types.BKIDField, rid)
	if err != nil {
		return nil, err
	}

	fields := []string{types.BKIDField, types.KubeNameField, types.BKClusterIDFiled}
	nodes, err := nlgc.PagedGetNodes(ctx, types.BKTableNameBaseNamespace, cond, fields, l.nsNodeParser, rid)
	if err != nil {
		return nil, err
	}

	err = nlgc.CrossCompareNodeInfoCache(l.topoKey, bizID, types.KubeNamespace, nodes, rid)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nodes, nil
	}

	// get next level nodes by namespace ids
	nsIDs := make([]int64, 0)
	nextConds := make([]mapstr.MapStr, 0)
	for _, node := range nodes {
		nsIDs = append(nsIDs, node.ID)
		if len(nsIDs) == 100 {
			nextConds = append(nextConds, mapstr.MapStr{types.BKNamespaceIDField: mapstr.MapStr{common.BKDBIN: nsIDs}})
			nsIDs = make([]int64, 0)
		}
	}
	if len(nsIDs) > 0 {
		nextConds = append(nextConds, mapstr.MapStr{types.BKNamespaceIDField: mapstr.MapStr{common.BKDBIN: nsIDs}})
	}

	childNodes, err := l.nextLevel.GetNodesByDB(ctx, bizID, nextConds, rid)
	if err != nil {
		return nil, err
	}

	nodes = nlgc.CombineChildNodes(nodes, childNodes)

	return nodes, nil
}

func (l *kubeNsLevel) nsNodeParser(ctx context.Context, data []mapstr.MapStr, rid string) ([]topotypes.Node, error) {
	// parse nodes
	nodes := make([]topotypes.Node, len(data))
	for i, item := range data {
		id, err := util.GetInt64ByInterface(item[types.BKIDField])
		if err != nil {
			blog.Errorf("parse kube ns id failed, err: %v, item: %+v, rid: %s", err, item, rid)
			return nil, err
		}

		parentID, err := util.GetInt64ByInterface(item[types.BKClusterIDFiled])
		if err != nil {
			blog.Errorf("parse kube ns cluster id failed, err: %v, item: %+v, rid: %s", err, item, rid)
			return nil, err
		}

		nodes[i] = topotypes.Node{
			Kind:     types.KubeNamespace,
			ID:       id,
			Name:     util.GetStrByInterface(item[types.KubeNameField]),
			ParentID: parentID,
		}
	}
	return nodes, nil
}

type kubeWlLevel struct {
	topoKey key.Key
}

func newKubeWlLevel() *kubeWlLevel {
	return &kubeWlLevel{
		topoKey: key.TopoKeyMap[topotypes.KubeType],
	}
}

// GetNodesByDB get all nodes that belongs to the topology level
func (l *kubeWlLevel) GetNodesByDB(ctx context.Context, bizID int64, conds []mapstr.MapStr, rid string) (
	[]topotypes.Node, error) {

	if len(conds) == 0 {
		cond, err := tools.GenKubeSharedNsCond(ctx, bizID, types.BKNamespaceIDField, rid)
		if err != nil {
			return nil, err
		}
		conds = append(conds, cond)
	}

	fields := []string{types.BKIDField, types.KubeNameField, types.BKNamespaceIDField}
	allNodes := make([]topotypes.Node, 0)
	for _, wlType := range types.GetWorkLoadObjects() {
		wlNodes := make([]topotypes.Node, 0)
		for _, cond := range conds {
			table, err := wlType.Table()
			if err != nil {
				return nil, err
			}

			nodes, err := nlgc.PagedGetNodes(ctx, table, cond, fields, wlNodeParser(wlType, bizID), rid)
			if err != nil {
				return nil, err
			}

			wlNodes = append(wlNodes, nodes...)
			allNodes = append(allNodes, nodes...)
		}

		err := nlgc.CrossCompareNodeInfoCache(l.topoKey, bizID, string(wlType), wlNodes, rid)
		if err != nil {
			return nil, err
		}
	}

	return allNodes, nil
}

func wlNodeParser(wlType types.WorkloadType, bizID int64) nlgc.NodeParser {
	return func(ctx context.Context, data []mapstr.MapStr, rid string) ([]topotypes.Node, error) {
		zeroCount := int64(0)

		nodes := make([]topotypes.Node, len(data))
		wlIDs := make([]int64, len(data))
		for i, item := range data {
			id, err := util.GetInt64ByInterface(item[types.BKIDField])
			if err != nil {
				blog.Errorf("parse kube wl id failed, err: %v, item: %+v, rid: %s", err, item, rid)
				return nil, err
			}

			parentID, err := util.GetInt64ByInterface(item[types.BKNamespaceIDField])
			if err != nil {
				blog.Errorf("parse kube wl cluster id failed, err: %v, item: %+v, rid: %s", err, item, rid)
				return nil, err
			}

			nodes[i] = topotypes.Node{
				Kind:     string(wlType),
				ID:       id,
				Name:     util.GetStrByInterface(item[types.KubeNameField]),
				ParentID: parentID,
				Count:    &zeroCount,
			}
			wlIDs[i] = id
		}

		return fillKubeWlNodeCount(ctx, bizID, string(wlType), wlIDs, nodes, rid)
	}
}

// fillKubeWlNodeCount fill workload nodes with container count
func fillKubeWlNodeCount(ctx context.Context, bizID int64, kind string, wlIDs []int64, nodes []topotypes.Node,
	rid string) ([]topotypes.Node, error) {

	if len(wlIDs) == 0 {
		return nodes, nil
	}

	// get pods by workloads
	podCond := mapstr.MapStr{
		types.RefIDField:   mapstr.MapStr{common.BKDBIN: wlIDs},
		types.RefKindField: kind,
	}

	pods := make([]types.Pod, 0)
	if err := mongodb.Client().Table(types.BKTableNameBasePod).Find(podCond).Fields(types.BKIDField,
		types.RefIDField).All(ctx, &pods); err != nil {
		blog.Errorf("get pod ids failed, cond: %+v, err: %v, rid: %s", podCond, err, rid)
		return nil, err
	}

	if len(pods) == 0 {
		if err := nlgc.DeleteNodeCountCache(key.TopoKeyMap[topotypes.KubeType], bizID, kind, wlIDs, rid); err != nil {
			return nil, err
		}
		return nodes, nil
	}

	podIDs := make([]int64, len(pods))
	podWlMap := make(map[int64]int64)
	for i, pod := range pods {
		podIDs[i] = pod.ID
		podWlMap[pod.ID] = pod.Ref.ID
	}

	// get container counts by pods
	containerCond := mapstr.MapStr{
		types.BKPodIDField: mapstr.MapStr{common.BKDBIN: podIDs},
	}

	containers := make([]types.Container, 0)
	if err := mongodb.Client().Table(types.BKTableNameBaseContainer).Find(containerCond).Fields(types.BKPodIDField).
		All(ctx, &containers); err != nil {
		blog.Errorf("get containers failed, cond: %+v, err: %v, rid: %s", containerCond, err, rid)
		return nil, err
	}

	if len(containers) == 0 {
		if err := nlgc.DeleteNodeCountCache(key.TopoKeyMap[topotypes.KubeType], bizID, kind, wlIDs, rid); err != nil {
			return nil, err
		}
		return nodes, nil
	}

	wlCountMap := make(map[int64]int64)
	for _, container := range containers {
		wlID := podWlMap[container.PodID]
		wlCountMap[wlID]++
	}

	for i, node := range nodes {
		cnt := wlCountMap[node.ID]
		nodes[i].Count = &cnt
	}

	// add workload topo nodes count cache
	if err := nlgc.AddNodeCountCache(key.TopoKeyMap[topotypes.KubeType], bizID, kind, wlCountMap, rid); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (l *kubeWlLevel) GetNodesByCache(ctx context.Context, bizID int64, rid string) ([]topotypes.Node, error) {
	allNodes := make([]topotypes.Node, 0)
	for _, wlType := range types.GetWorkLoadObjects() {
		kind := string(wlType)

		// get nodes from cache
		nodes, err := nlgc.GetNodeInfoCache(l.topoKey, bizID, kind, rid)
		if err != nil {
			return nil, err
		}

		// get workload topo nodes count from cache
		pagedNodes := make([]topotypes.Node, 0)
		for _, node := range nodes {
			pagedNodes = append(pagedNodes, node)
			if len(pagedNodes) == common.BKMaxPageSize {
				if err = l.fillWlNodesCountFromCache(ctx, bizID, kind, pagedNodes, rid); err != nil {
					return nil, err
				}
				allNodes = append(allNodes, pagedNodes...)
				pagedNodes = make([]topotypes.Node, 0)
			}
		}

		if len(pagedNodes) > 0 {
			if err = l.fillWlNodesCountFromCache(ctx, bizID, kind, pagedNodes, rid); err != nil {
				return nil, err
			}
			allNodes = append(allNodes, pagedNodes...)
		}
	}

	return allNodes, nil
}

func (l *kubeWlLevel) fillWlNodesCountFromCache(ctx context.Context, bizID int64, kind string,
	nodes []topotypes.Node, rid string) error {

	cntKeys := make([]string, len(nodes))
	for i, node := range nodes {
		cntKeys[i] = nlgc.GenNodeInfoCntKey(l.topoKey, bizID, kind, node.ID)
	}

	cntRes, err := redis.Client().MGet(ctx, cntKeys...).Result()
	if err != nil {
		blog.Errorf("get node count cache by keys: %+v failed, err: %v, rid: %s", cntKeys, err, rid)
		return err
	}

	zero := int64(0)
	for i, cnt := range cntRes {
		if cnt == nil {
			nodes[i].Count = &zero
			continue
		}

		cntStr, ok := cnt.(string)
		if !ok {
			blog.Errorf("%s workload %d count cache %v type %T is invalid, rid: %s", kind, nodes[i].ID, cnt, cnt, rid)
			nodes[i].Count = &zero
			continue
		}

		count, err := strconv.ParseInt(cntStr, 10, 64)
		if err != nil {
			blog.Errorf("parse %s wl %d count cache %s failed, err: %v, rid: %s", kind, nodes[i].ID, cntStr, err, rid)
			return err
		}

		nodes[i].Count = &count
	}

	return nil
}
