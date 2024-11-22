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

package kube

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/util/errors"
	"configcenter/src/kube/types"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

type kubeOperation struct {
}

// New create a new model manager instance
func New() core.KubeOperation {
	kube := new(kubeOperation)
	return kube
}

// GetSysSpecInfoByCond get the spec redundancy information required by the pod.
func (p *kubeOperation) GetSysSpecInfoByCond(kit *rest.Kit, podsInfos []types.PodsInfoArray) ([]types.SysSpec, []int64,
	ccErr.CCErrorCoder) {

	if len(podsInfos) == 0 {
		return make([]types.SysSpec, 0), make([]int64, 0), nil
	}

	wlIDMap := make(map[types.WorkloadType][]int64)
	nodeIDs := make([]int64, 0)
	for _, info := range podsInfos {
		for _, pod := range info.Pods {
			spec := pod.Spec
			wlIDMap[spec.Ref.Kind] = append(wlIDMap[spec.Ref.Kind], spec.Ref.ID)
			nodeIDs = append(nodeIDs, spec.NodeID)
		}
	}

	wlMap, err := p.getWlInfo(kit, wlIDMap)
	if err != nil {
		return nil, nil, err
	}

	clusterPlatBizMap, err := p.getSharedClusterInfo(kit, podsInfos, wlMap)
	if err != nil {
		return nil, nil, err
	}

	nodeMap, noPodNodeIDs, err := p.getNodeInfo(kit, nodeIDs)
	if err != nil {
		return nil, nil, err
	}

	// combine all pod related info to system spec array
	sysSpecArr := make([]types.SysSpec, 0)
	for _, info := range podsInfos {
		for _, pod := range info.Pods {
			spec := pod.Spec
			workload := wlMap[spec.Ref.Kind][spec.Ref.ID]

			node, exists := nodeMap[spec.NodeID]
			if !exists {
				blog.Errorf("pod related node(id: %d) not exists, rid: %s", spec.NodeID, kit.Rid)
				return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeNode)
			}

			if spec.ClusterID != node.ClusterID || clusterPlatBizMap[node.ClusterID] != node.BizID {
				blog.Errorf("pod related node(%+v) is invalid, rid: %s", node, kit.Rid)
				return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeNode)
			}

			if node.Name == nil {
				blog.Errorf("node(id: %d) name is nil, rid: %s", spec.NodeID, kit.Rid)
				return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeNode)
			}

			sysSpecArr = append(sysSpecArr, types.SysSpec{
				WorkloadSpec: types.WorkloadSpec{
					NamespaceSpec: types.NamespaceSpec{
						ClusterSpec: types.ClusterSpec{
							BizID:      workload.BizID,
							ClusterUID: workload.ClusterUID,
							ClusterID:  spec.ClusterID,
						},
						Namespace:   workload.Namespace,
						NamespaceID: spec.NamespaceID,
					},
					Ref: &types.Reference{Kind: spec.Ref.Kind, Name: workload.Name, ID: spec.Ref.ID},
				},
				TenantID: kit.TenantID,
				HostID:   pod.HostID,
				NodeID:   spec.NodeID,
				Node:     *node.Name,
			})
		}
	}

	return sysSpecArr, noPodNodeIDs, nil
}

// getWlInfo get workload id to workload info map
func (p *kubeOperation) getWlInfo(kit *rest.Kit, wlIDMap map[types.WorkloadType][]int64) (
	map[types.WorkloadType]map[int64]types.WorkloadBase, ccErr.CCErrorCoder) {

	wlMap := make(map[types.WorkloadType]map[int64]types.WorkloadBase)

	for kind, ids := range wlIDMap {
		// 通过workload kind 获取表名
		tableName, err := kind.Table()
		if err != nil {
			blog.Errorf("get collection failed, kind: %s, err: %v, rid: %s", kind, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeWorkload)
		}

		filter := map[string]interface{}{
			types.BKIDField: mapstr.MapStr{common.BKDBIN: ids},
		}

		kubeField := []string{types.BKIDField, types.KubeNameField, types.BKBizIDField, types.BKClusterIDFiled,
			types.ClusterUIDField, types.BKNamespaceIDField, types.NamespaceField}

		workloads := make([]types.WorkloadBase, 0)
		err = mongodb.Client().Table(tableName).Find(filter).Fields(kubeField...).All(kit.Ctx, &workloads)
		if err != nil {
			blog.Errorf("get workload failed, err: %v, filter: %v, rid:%s", err, filter, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}

		wlMap[kind] = make(map[int64]types.WorkloadBase)

		for _, workload := range workloads {
			wlMap[kind][workload.ID] = workload
		}
	}

	return wlMap, nil
}

func (p *kubeOperation) getSharedClusterInfo(kit *rest.Kit, podsInfos []types.PodsInfoArray,
	wlMap map[types.WorkloadType]map[int64]types.WorkloadBase) (map[int64]int64, ccErr.CCErrorCoder) {

	if len(podsInfos) == 0 {
		return make(map[int64]int64), nil
	}

	// check if workload info matches pod info, and check if namespace biz id is the same with input biz id
	matchingNsIDMap, mismatchNsIDMap := make(map[int64][]int64), make(map[int64][]int64)
	matchingNsClusterIDMap, mismatchNsBizIDMap := make(map[int64]int64), make(map[int64]int64)

	for _, info := range podsInfos {
		for _, pod := range info.Pods {
			spec := pod.Spec

			workload, exists := wlMap[spec.Ref.Kind][spec.Ref.ID]
			if !exists {
				blog.Errorf("pod related workload(%+v) not exists, rid: %s", spec.Ref, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeWorkload)
			}

			if spec.NamespaceID != workload.NamespaceID || spec.ClusterID != workload.ClusterID {
				blog.Errorf("pod(spec: %+v) related workload(%+v) is invalid, rid: %s", spec, workload, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeWorkload)
			}

			if info.BizID == workload.BizID {
				matchingNsClusterIDMap[spec.NamespaceID] = spec.ClusterID
				matchingNsIDMap[info.BizID] = append(matchingNsIDMap[info.BizID], spec.NamespaceID)
				continue
			}

			mismatchNsBizIDMap[spec.NamespaceID] = info.BizID
			mismatchNsIDMap[info.BizID] = append(mismatchNsIDMap[info.BizID], spec.NamespaceID)
		}
	}

	// get namespace related shared cluster relations
	sharedNsCond := make([]mapstr.MapStr, 0)
	for bizID, ids := range matchingNsIDMap {
		sharedNsCond = append(sharedNsCond, mapstr.MapStr{
			types.BKBizIDField:       bizID,
			types.BKNamespaceIDField: mapstr.MapStr{common.BKDBIN: ids},
		})
	}

	for bizID, ids := range mismatchNsIDMap {
		sharedNsCond = append(sharedNsCond, mapstr.MapStr{
			types.BKAsstBizIDField:   bizID,
			types.BKNamespaceIDField: mapstr.MapStr{common.BKDBIN: ids},
		})
	}

	clusterPlatBizMap, err := p.getSharedClusterInfoByCond(kit, sharedNsCond, matchingNsClusterIDMap,
		mismatchNsBizIDMap, matchingNsIDMap)
	if err != nil {
		return nil, err
	}

	return clusterPlatBizMap, nil
}

func (p *kubeOperation) getSharedClusterInfoByCond(kit *rest.Kit, sharedNsCond []mapstr.MapStr,
	matchingNsClusterIDMap, mismatchNsBizIDMap map[int64]int64, matchingNsIDMap map[int64][]int64) (map[int64]int64,
	ccErr.CCErrorCoder) {

	sharedRelCond := mapstr.MapStr{common.BKDBOR: sharedNsCond}
	sharedRel := make([]types.NsSharedClusterRel, 0)
	err := mongodb.Client().Table(types.BKTableNameNsSharedClusterRel).Find(sharedRelCond).All(kit.Ctx, &sharedRel)
	if err != nil {
		blog.Errorf("get shared cluster relations failed, err: %v, cond: %v, rid:%s", err, sharedRelCond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// checks if namespace biz ids are valid, then generate cluster and its plat biz map for node check
	clusterPlatBizMap := make(map[int64]int64)
	for _, rel := range sharedRel {
		clusterPlatBizMap[rel.ClusterID] = rel.AsstBizID

		inputBiz, exists := mismatchNsBizIDMap[rel.NamespaceID]
		if !exists {
			continue
		}

		if inputBiz != rel.AsstBizID {
			blog.Errorf("ns %d input biz %d != asst biz %d, rid: %s", rel.NamespaceID, inputBiz, rel.AsstBizID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeNamespace)
		}

		delete(mismatchNsBizIDMap, rel.NamespaceID)
	}

	if len(mismatchNsBizIDMap) > 0 {
		blog.Errorf("some shared ns info(%+v) is invalid, rid: %s", mismatchNsBizIDMap, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KubeNamespace)
	}

	// fill matching ns related normal cluster and its id to clusterPlatBizMap
	for bizID, nsIDs := range matchingNsIDMap {
		for _, nsID := range nsIDs {
			clusterID := matchingNsClusterIDMap[nsID]
			_, exists := clusterPlatBizMap[clusterID]
			if !exists {
				clusterPlatBizMap[clusterID] = bizID
			}
		}
	}
	return clusterPlatBizMap, nil
}

// getNodeInfo get node id to node info map
func (p *kubeOperation) getNodeInfo(kit *rest.Kit, nodeIDs []int64) (map[int64]types.Node, []int64,
	ccErr.CCErrorCoder) {

	filter := map[string]interface{}{
		types.BKIDField: mapstr.MapStr{common.BKDBIN: nodeIDs},
	}

	fields := []string{types.BKIDField, types.KubeNameField, types.BKBizIDField, types.BKClusterIDFiled,
		types.HasPodField}

	nodes := make([]types.Node, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(filter).Fields(fields...).All(kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("get node failed, err: %v, filter: %v, rid:%s", err, filter, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	nodeMap := make(map[int64]types.Node)
	noPodNodeIDs := make([]int64, 0, len(nodes))
	for _, node := range nodes {
		nodeMap[node.ID] = node

		// if node does not have pods previously, returns these ids to set their hasPod flag to true
		if node.HasPod == nil || !*node.HasPod {
			noPodNodeIDs = append(noPodNodeIDs, node.ID)
		}
	}

	return nodeMap, noPodNodeIDs, nil
}

func (p *kubeOperation) getNodeRelatedInfo(kit *rest.Kit, data []types.OneNodeCreateOption) (map[int64]types.Cluster,
	map[int64]int64, ccErr.CCErrorCoder) {

	clusterIDs := make([]int64, len(data))
	hostIDs := make([]int64, len(data))
	for i, info := range data {
		clusterIDs[i] = info.ClusterID
		hostIDs[i] = info.HostID
	}

	// 获取cluster信息
	clusterFilter := map[string]interface{}{
		types.BKIDField: map[string]interface{}{common.BKDBIN: clusterIDs},
	}

	clusters := make([]types.Cluster, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(clusterFilter).
		Fields(types.BKBizIDField, types.UidField, types.BKIDField).All(kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("query cluster failed, filter: %+v, err: %s, rid:%s", clusterFilter, err, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	clusterMap := make(map[int64]types.Cluster)
	for _, cluster := range clusters {
		clusterMap[cluster.ID] = cluster
	}

	// get host id to biz id map
	hostFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
	}

	relations := make([]metadata.ModuleHost, 0)
	err = mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(hostFilter).
		Fields(common.BKAppIDField, common.BKHostIDField).All(kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("query host relation failed, filter: %+v, err: %s, rid:%s", hostFilter, err, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	hostMap := make(map[int64]int64)
	for _, rel := range relations {
		hostMap[rel.HostID] = rel.AppID
	}
	return clusterMap, hostMap, nil
}

// BatchCreateNode create container node data in batches.
func (p *kubeOperation) BatchCreateNode(kit *rest.Kit, data []types.OneNodeCreateOption) (
	[]*types.Node, ccErr.CCErrorCoder) {

	nameClusterID := make(map[string]int64)
	for _, node := range data {
		if node.Name == nil {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "name")
		}
		nameClusterID[*node.Name] = node.ClusterID
	}

	// get map[clusterID]cluster & map[hostID]bizID for validation
	clusterMap, hostMap, cErr := p.getNodeRelatedInfo(kit, data)
	if cErr != nil {
		return nil, cErr
	}

	// generate ids field
	ids, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBaseNode, len(data))
	if err != nil {
		blog.Errorf("create node failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	result := make([]*types.Node, 0)
	now := time.Now().Unix()
	hasPod := false

	for idx, node := range data {
		cluster, exists := clusterMap[node.ClusterID]
		if !exists || cluster.BizID != node.BizID || cluster.Uid == nil {
			blog.Errorf("node(biz %d) related cluster %d is invalid, rid: %s", node.BizID, node.ClusterID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKClusterIDFiled)
		}

		if hostMap[node.HostID] != node.BizID {
			blog.Errorf("node(biz %d) related host %d is invalid, rid: %s", node.BizID, node.HostID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}

		node := &types.Node{
			ID:               int64(ids[idx]),
			BizID:            node.BizID,
			HostID:           node.HostID,
			ClusterID:        node.ClusterID,
			ClusterUID:       *cluster.Uid,
			Name:             node.Name,
			Roles:            node.Roles,
			Labels:           node.Labels,
			Taints:           node.Taints,
			Unschedulable:    node.Unschedulable,
			InternalIP:       node.InternalIP,
			ExternalIP:       node.ExternalIP,
			HasPod:           &hasPod,
			HostName:         node.HostName,
			RuntimeComponent: node.RuntimeComponent,
			KubeProxyMode:    node.KubeProxyMode,
			PodCidr:          node.PodCidr,
			TenantID:         kit.TenantID,
			Revision: table.Revision{
				CreateTime: now,
				LastTime:   now,
				Creator:    kit.User,
				Modifier:   kit.User,
			},
		}
		if err := mongodb.Client().Table(types.BKTableNameBaseNode).Insert(kit.Ctx, node); err != nil {
			blog.Errorf("create node failed, db insert failed, node: %+v, err: %+v, rid: %s", node, err, kit.Rid)
			return nil, errors.ConvDBInsertError(kit, mongodb.Client(), err)
		}
		result = append(result, node)
	}

	return result, nil
}

// CheckPlatBizSharedNs check if platform biz's ns is in shared cluster and if its related biz matches the plat biz
func (p *kubeOperation) CheckPlatBizSharedNs(kit *rest.Kit, bizNsMap map[int64][]int64) ccErr.CCErrorCoder {
	if len(bizNsMap) == 0 {
		return nil
	}

	nsCnt := 0
	conds := make([]mapstr.MapStr, 0)

	for bizID, nsIDs := range bizNsMap {
		nsIDs = util.IntArrayUnique(nsIDs)
		nsCnt += len(nsIDs)
		conds = append(conds, mapstr.MapStr{
			types.BKAsstBizIDField:   bizID,
			types.BKNamespaceIDField: map[string]interface{}{common.BKDBIN: nsIDs},
		})
	}

	cond := mapstr.MapStr{common.BKDBOR: conds}
	cnt, err := mongodb.Client().Table(types.BKTableNameNsSharedClusterRel).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count shared ns failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if int(cnt) != nsCnt {
		blog.Errorf("shared ns count %d is invalid, cond: %+v, rid: %s", cnt, cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	return nil
}
