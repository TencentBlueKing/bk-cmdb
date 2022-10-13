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
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

type kubeOperation struct {
}

// New create a new model manager instance
func New() core.KubeOperation {
	container := &kubeOperation{}
	return container
}

func validateNodeData(kit *rest.Kit, bizID int64, node types.OneNodeCreateOption) errors.CCErrorCoder {

	clusterFilter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		types.BKIDField:       node.ClusterID,
	}

	cnt, err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(clusterFilter).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("query database failed, filter: %v, err: %v, rid: %s", clusterFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if cnt != 1 {
		blog.Errorf("query database num error, filter: %v, error: %v, rid: %s", clusterFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	filter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		common.BKHostIDField:  node.HostID,
	}
	_, err = mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query host module config failed, err: %s, rid:%s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if cnt <= 0 {
		blog.Errorf("query host module config count error, err: %s, rid:%s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return nil
}

// GetSysSpecInfoByCond  get the spec redundancy information required by the pod.
func (p *kubeOperation) GetSysSpecInfoByCond(kit *rest.Kit, spec types.SpecInfo, bizID int64,
	hostID int64) (*types.SysSpec, errors.CCErrorCoder) {
	// 通过workload kind 获取表名
	tableName, err := types.GetCollectionWithObject(spec.Ref.Kind)
	if err != nil {
		blog.Errorf("get collection failed, kind: %s, err: %v, rid: %s", spec.Ref.Kind, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	filter := map[string]interface{}{
		common.BKAppIDField:      bizID,
		common.BKOwnerIDField:    kit.SupplierAccount,
		types.BKClusterIDFiled:   *spec.ClusterID,
		types.BKNamespaceIDField: *spec.NamespaceID,
		types.BKIDField:          spec.Ref.ID,
	}
	kubeField := []string{types.ClusterUIDField, types.NamespaceField, types.KubeNameField}

	workload := make([]map[string]interface{}, 0)
	err = mongodb.Client().Table(tableName[0]).Find(filter).Fields(kubeField...).All(kit.Ctx, &workload)
	if err != nil {
		blog.Errorf("query host module config failed, err: %s, rid:%s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(workload) != 1 {
		blog.Errorf("workload gets the wrong amount, filter: %+v, num: %d, rid: %s", filter, len(workload), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	clusterUID := util.GetStrByInterface(workload[0][types.ClusterUIDField])
	if err != nil {
		blog.Errorf("convert clusterUID failed, workload: %+v, rid: %s", workload, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	namespace := util.GetStrByInterface(workload[0][types.NamespaceField])
	if err != nil {
		blog.Errorf("convert namespace failed, workload: %+v, rid: %s", workload, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	workloadName := util.GetStrByInterface(workload[0][types.KubeNameField])
	if err != nil {
		blog.Errorf("convert workloadName failed, workload: %+v, rid: %s", workload, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	nodeFilter := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKOwnerIDField:  kit.SupplierAccount,
		types.BKClusterIDFiled: *spec.ClusterID,
		types.BKIDField:        *spec.NodeID,
	}

	nodes := make([]map[string]interface{}, 0)
	err = mongodb.Client().Table(types.BKTableNameBaseNode).Find(nodeFilter).
		Fields([]string{types.KubeNameField}...).All(kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("query node failed, filter: %+v, err: %s, rid:%s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(nodes) != 1 {
		blog.Errorf("node gets the wrong amount, filter: %+v, num: %d, rid: %s", filter, len(nodes), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	nodeName := util.GetStrByInterface(nodes[0][types.KubeNameField])
	if err != nil {
		blog.Errorf("convert nodeName failed, workload: %+v, rid: %s", nodes, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &types.SysSpec{
		BizID:           bizID,
		SupplierAccount: kit.SupplierAccount,
		ClusterID:       *spec.ClusterID,
		ClusterUID:      clusterUID,
		NameSpaceID:     *spec.NamespaceID,
		NameSpace:       namespace,
		Workload:        types.Ref{Kind: spec.Ref.Kind, Name: workloadName, ID: spec.Ref.ID},
		HostID:          hostID,
		NodeID:          *spec.NodeID,
		Node:            nodeName,
	}, nil
}

//// BatchCreatePod create container node data in batches.
//func (p *kubeOperation) BatchCreatePod(kit *rest.Kit, data []types.PodsInfoArray) ([]types.Pod,
//	errors.CCErrorCoder) {
//	// generate pod ids field
//	ids, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBasePod, len(data))
//	if err != nil {
//		blog.Errorf("create pods failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
//		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
//	}
//	pods := make([]types.Pod, 0)
//	now := time.Now().Unix()
//
//	for _, info := range data {
//		for idx, pod := range info.Pods {
//			// 获取pod表中冗余的信息
//			sysSpec, ccErr := p.GetSysSpecInfoByCond(kit, pod.Spec, info.BizID, pod.HostID)
//			if ccErr != nil {
//				return nil, ccErr
//			}
//			podInfo := types.Pod{
//				ID:            int64(ids[idx]),
//				SysSpec:       *sysSpec,
//				Name:          pod.Name,
//				Priority:      pod.Priority,
//				Labels:        pod.Labels,
//				IP:            pod.IP,
//				IPs:           pod.IPs,
//				Volumes:       pod.Volumes,
//				QOSClass:      pod.QOSClass,
//				NodeSelectors: pod.NodeSelectors,
//				Tolerations:   pod.Tolerations,
//				Revision: table.Revision{
//					CreateTime: now,
//					LastTime:   now,
//					Creator:    kit.User,
//					Modifier:   kit.User,
//				},
//			}
//			pods = append(pods, podInfo)
//			if err := mongodb.Client().Table(types.BKTableNameBasePod).Insert(kit.Ctx, &podInfo); err != nil {
//				blog.Errorf("create pod failed, db insert failed, node: %+v, err: %+v, rid: %s", podInfo, err, kit.Rid)
//				return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
//			}
//
//			// generate pod ids field
//			containerIDs, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBaseContainer,
//				len(pod.Containers))
//			if err != nil {
//				blog.Errorf("create container failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
//				return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
//			}
//
//			for id, info := range pod.Containers {
//				container := &types.Container{
//					ID:                  int64(containerIDs[id]),
//					PodID:               int64(ids[idx]),
//					ContainerBaseFields: info,
//					Revision: table.Revision{
//						CreateTime: now,
//						LastTime:   now,
//						Creator:    kit.User,
//						Modifier:   kit.User,
//					},
//				}
//				err := mongodb.Client().Table(types.BKTableNameBaseContainer).Insert(kit.Ctx, container)
//				if err != nil {
//					blog.Errorf("create container failed, db insert failed, container: %+v, err: %+v, rid: %s",
//						container, err, kit.Rid)
//					return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
//				}
//			}
//		}
//	}
//
//	return pods, nil
//}

func (p *kubeOperation) getClusterInfo(kit *rest.Kit, bizID int64, data []types.OneNodeCreateOption) (
	map[int64]string, errors.CCErrorCoder) {
	clusterIDs := make([]int64, 0)
	clusterIDName := make(map[int64]string)
	for _, info := range data {
		clusterIDs = append(clusterIDs, info.ClusterID)
		clusterIDName[info.ClusterID] = *info.Name
	}

	// 获取cluster信息
	clusterFilter := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKOwnerIDField:  kit.SupplierAccount,
		types.BKClusterIDFiled: map[string]interface{}{common.BKDBIN: clusterIDs},
	}

	clusters := make([]types.Cluster, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(clusterFilter).
		Fields([]string{types.ClusterUIDField, types.BKClusterIDFiled}...).All(kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("query cluster failed, filter: %+v, err: %s, rid:%s", clusterFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	clusterMap := make(map[int64]string)
	for _, cluster := range clusters {
		clusterMap[cluster.ID] = *cluster.Uid
	}
	return clusterMap, nil
}

// BatchCreateNode create container node data in batches.
func (p *kubeOperation) BatchCreateNode(kit *rest.Kit, bizID int64, data []types.OneNodeCreateOption) (
	[]*types.Node, errors.CCErrorCoder) {

	nameClusterID := make(map[string]int64)
	for _, node := range data {
		nameClusterID[*node.Name] = node.ClusterID
	}

	// cluster_id=>cluster_uid
	clusterMap, cErr := p.getClusterInfo(kit, bizID, data)
	if cErr != nil {
		return nil, cErr
	}

	for _, node := range data {
		filter := map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{
					common.BKAppIDField:    bizID,
					types.BKClusterIDFiled: nameClusterID[*node.Name],
					types.KubeNameField:    *node.Name},
				{
					common.BKAppIDField:   bizID,
					types.KubeNameField:   *node.Name,
					types.ClusterUIDField: clusterMap[nameClusterID[*node.Name]],
				},
			},
		}

		count, err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(filter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("query node failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if count > 0 {
			blog.Errorf("create node failed, there are duplicate node names, filter: %+v, rid: %s", filter, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "names")
		}
	}

	// generate ids field
	ids, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBaseNode, len(data))
	if err != nil {
		blog.Errorf("create node failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	result := make([]*types.Node, 0)
	now := time.Now().Unix()
	for idx, node := range data {
		if err := validateNodeData(kit, bizID, node); err != nil {
			return nil, err
		}
		node := &types.Node{
			ID:               int64(ids[idx]),
			BizID:            bizID,
			HostID:           node.HostID,
			ClusterID:        node.ClusterID,
			ClusterUID:       clusterMap[node.ClusterID],
			Name:             node.Name,
			Roles:            node.Roles,
			Labels:           node.Labels,
			Taints:           node.Taints,
			Unschedulable:    node.Unschedulable,
			InternalIP:       node.InternalIP,
			ExternalIP:       node.ExternalIP,
			HostName:         node.HostName,
			RuntimeComponent: node.RuntimeComponent,
			KubeProxyMode:    node.KubeProxyMode,
			PodCidr:          node.PodCidr,
			SupplierAccount:  kit.SupplierAccount,
			Revision: table.Revision{
				CreateTime: now,
				LastTime:   now,
				Creator:    kit.User,
				Modifier:   kit.User,
			},
		}
		if err := mongodb.Client().Table(types.BKTableNameBaseNode).Insert(kit.Ctx, node); err != nil {
			blog.Errorf("create node failed, db insert failed, node: %+v, err: %+v, rid: %s", node, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
		}
		result = append(result, node)
	}

	return result, nil
}

// CreateCluster create cluster instance.
func (p *kubeOperation) CreateCluster(kit *rest.Kit, bizID int64, data *types.Cluster) (*types.Cluster,
	errors.CCErrorCoder) {

	// it is necessary to judge whether there is duplicate data here, to prevent subsequent calls to coreservice
	// directly and lack of verification.
	if err := data.CreateValidate(); err != nil {
		blog.Errorf("create cluster failed, data: %+v, err: %+v, rid: %s", data, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	nameFilter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKAppIDField:     bizID,
				common.BKOwnerIDField:   kit.SupplierAccount,
				types.ContainerUIDField: *data.Uid,
			},
			{
				common.BKAppIDField:     bizID,
				common.BKOwnerIDField:   kit.SupplierAccount,
				types.ContainerUIDField: *data.Name,
			},
		},
	}

	count, err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(nameFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count cluster failed, filter: %+v, err: %+v, rid: %s", nameFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("create cluster failed, name or uid duplicated, name: %s, uid: %s, rid: %s", data.Name,
			data.Uid, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "name or uid")
	}

	// generate id field
	id, err := mongodb.Client().NextSequence(kit.Ctx, types.BKTableNameBaseCluster)
	if err != nil {
		blog.Errorf("create cluster failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	now := time.Now().Unix()
	cluster := &types.Cluster{
		ID:               int64(id),
		BizID:            bizID,
		SupplierAccount:  kit.SupplierAccount,
		Name:             data.Name,
		SchedulingEngine: data.SchedulingEngine,
		Uid:              data.Uid,
		Xid:              data.Xid,
		Version:          data.Version,
		NetworkType:      data.NetworkType,
		Region:           data.Region,
		Vpc:              data.Vpc,
		NetWork:          data.NetWork,
		Type:             data.Type,
		Revision: table.Revision{
			CreateTime: now,
			LastTime:   now,
			Creator:    kit.User,
			Modifier:   kit.User,
		},
	}

	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Insert(kit.Ctx, cluster); err != nil {
		blog.Errorf("create cluster failed, db insert failed, doc: %+v, err: %+v, rid: %s", cluster, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}

	return cluster, nil
}
