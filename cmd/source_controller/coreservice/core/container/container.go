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

package container

import (
	"configcenter/cmd/source_controller/coreservice/core"
	"configcenter/pkg/kube/orm"
	types2 "configcenter/pkg/kube/types"
	"configcenter/pkg/storage/dal/table"
	"configcenter/pkg/storage/driver/mongodb"
	"time"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/errors"
	"configcenter/pkg/http/rest"
	"configcenter/pkg/metadata"
	"configcenter/pkg/util"
)

type containerOperation struct {
}

// New create a new model manager instance
func New() core.ContainerOperation {
	container := &containerOperation{}
	return container
}

// BatchDeleteNode delete node instances.
func (p *containerOperation) BatchDeleteNode(kit *rest.Kit, bizID int64, option *types2.BatchDeleteNodeOption) (
	*metadata.DeletedCount, errors.CCErrorCoder) {

	num := 0
	for _, node := range option.Nodes {
		filter := map[string]interface{}{
			common.BKAppIDField:     bizID,
			common.BKOwnerIDField:   kit.SupplierAccount,
			types2.BKClusterIDFiled: node.ClusterID,
			types2.BKIDField: map[string]interface{}{
				common.BKDBIN: node.IDs,
			},
		}
		if err := mongodb.Client().Table(types2.BKTableNameBaseNode).Delete(kit.Ctx, filter); err != nil {
			blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
		}
		num += len(node.IDs)
	}

	return &metadata.DeletedCount{Count: uint64(num)}, nil
}

// DeleteCluster delete cluster instance.
func (p *containerOperation) DeleteCluster(kit *rest.Kit, bizID int64, option *types2.DeleteClusterOption) (
	*metadata.DeletedCount, errors.CCErrorCoder) {

	filter := make(map[string]interface{}, 0)
	num := 0
	if len(option.IDs) > 0 {
		num = len(option.IDs)
		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: kit.SupplierAccount,
			types2.BKIDField: map[string]interface{}{
				common.BKDBIN: option.IDs,
			},
		}
	}
	if err := mongodb.Client().Table(types2.BKTableNameBaseCluster).Delete(kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}
	return &metadata.DeletedCount{Count: uint64(num)}, nil
}

func validateNodeData(kit *rest.Kit, bizID int64, node types2.OneNodeCreateOption) errors.CCErrorCoder {

	clusterFilter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		types2.BKIDField:      node.ClusterID,
	}

	cnt, err := mongodb.Client().Table(types2.BKTableNameBaseCluster).Find(clusterFilter).Count(kit.Ctx)
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

func (p *containerOperation) getSysSpecInfoByCond(kit *rest.Kit, cmdbSpec *types2.SpecInfo, bizID int64,
	hostID int64) (*types2.SysSpec, errors.CCErrorCoder) {
	// 通过workload kind 获取表名
	tableName, err := types2.GetCollectionWithObject(*cmdbSpec.WorkloadKind)
	if err != nil {
		blog.Errorf("get collection failed, kind: %s, err: %v, rid: %s", *cmdbSpec.WorkloadKind, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	filter := map[string]interface{}{
		common.BKAppIDField:       bizID,
		common.BKOwnerIDField:     kit.SupplierAccount,
		types2.BKClusterIDFiled:   *cmdbSpec.ClusterID,
		types2.BKNamespaceIDField: *cmdbSpec.NamespaceID,
		types2.BKIDField:          *cmdbSpec.WorkloadID,
	}
	kubeField := []string{types2.ClusterUIDField, types2.NamespaceField, types2.KubeNameField}
	workload := make([]map[string]interface{}, 0)

	err = mongodb.Client().Table(tableName[0]).Find(filter).Fields(kubeField...).All(kit.Ctx, &workload)
	if err != nil {
		blog.Errorf("query host module config failed, err: %s, rid:%s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(workload) == 0 {
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(workload) > 1 {
		return nil, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	clusterUID := util.GetStrByInterface(workload[0][types2.ClusterUIDField])
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	namespace := util.GetStrByInterface(workload[0][types2.NamespaceField])
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	workloadname := util.GetStrByInterface(workload[0][types2.KubeNameField])
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	nodeFilter := map[string]interface{}{
		common.BKAppIDField:     bizID,
		common.BKOwnerIDField:   kit.SupplierAccount,
		types2.BKClusterIDFiled: *cmdbSpec.ClusterID,
		types2.BKIDField:        *cmdbSpec.NodeID,
	}

	nodes := make([]map[string]interface{}, 0)
	err = mongodb.Client().Table(types2.BKTableNameBaseNode).Find(nodeFilter).
		Fields([]string{types2.KubeNameField}...).All(kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("query node failed, filter: %+v, err: %s, rid:%s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(nodes) == 0 {
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(nodes) > 1 {
		return nil, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	nodename := util.GetStrByInterface(nodes[0][types2.KubeNameField])
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &types2.SysSpec{
		BizID:           bizID,
		SupplierAccount: kit.SupplierAccount,
		ClusterID:       *cmdbSpec.ClusterID,
		ClusterUID:      clusterUID,
		NameSpaceID:     *cmdbSpec.NamespaceID,
		NameSpace:       namespace,
		Workload: types2.Ref{
			Kind: *cmdbSpec.WorkloadKind,
			Name: workloadname,
			ID:   *cmdbSpec.WorkloadID,
		},
		HostID: hostID,
		NodeID: *cmdbSpec.NodeID,
		Node:   nodename,
	}, nil
}

// BatchCreatePod create container node data in batches.
func (p *containerOperation) BatchCreatePod(kit *rest.Kit, bizID int64, data []types2.PodsInfo) ([]types2.Pod,
	errors.CCErrorCoder) {
	// generate pod ids field
	ids, err := mongodb.Client().NextSequences(kit.Ctx, types2.BKTableNameBasePod, len(data))
	if nil != err {
		blog.Errorf("create pod failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	pods := make([]types2.Pod, 0)
	now := time.Now().Unix()

	// 这里有一个问题，就是如果用户传的是ccID  那么应该是不知道他上一级的kube信息。所以这里要获一下上一级信息。同理如果传的是kube信息，那么
	// 也需要获取一下上一级的ccID信息。
	for idx, pod := range data {
		sysSpec, ccErr := p.getSysSpecInfoByCond(kit, pod.Spec, bizID, pod.HostID)
		if ccErr != nil {
			return nil, ccErr
		}
		podInfo := types2.Pod{
			ID:            int64(ids[idx]),
			SysSpec:       *sysSpec,
			PodBaseFields: pod.PodBaseFields,
			Revision: table.Revision{
				CreateTime: now,
				LastTime:   now,
				Creator:    kit.User,
				Modifier:   kit.User,
			},
		}
		pods = append(pods, podInfo)
		if err := mongodb.Client().Table(types2.BKTableNameBasePod).Insert(kit.Ctx, &podInfo); err != nil {
			blog.Errorf("create pod failed, db insert failed, node: %+v, err: %+v, rid: %s", podInfo, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
		}

		// generate pod ids field
		containerIDs, err := mongodb.Client().NextSequences(kit.Ctx, types2.BKTableNameBaseContainer, len(pod.Containers))
		if nil != err {
			blog.Errorf("create container failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
		}

		for id, info := range pod.Containers {
			container := &types2.Container{
				ID:                  int64(containerIDs[id]),
				PodID:               int64(ids[idx]),
				ContainerBaseFields: info,
				Revision: table.Revision{
					CreateTime: now,
					LastTime:   now,
					Creator:    kit.User,
					Modifier:   kit.User,
				},
			}
			if err := mongodb.Client().Table(types2.BKTableNameBaseContainer).Insert(kit.Ctx, container); err != nil {
				blog.Errorf("create container failed, db insert failed, container: %+v, err: %+v, rid: %s",
					container, err, kit.Rid)
				return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
			}
		}
	}
	return pods, nil
}

// BatchCreateNode create container node data in batches.
func (p *containerOperation) BatchCreateNode(kit *rest.Kit, bizID int64, data []types2.OneNodeCreateOption) (
	[]*types2.Node, errors.CCErrorCoder) {

	names := make([]string, 0)
	for _, node := range data {
		names = append(names, *node.Name)
	}

	filter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		types2.KubeNameField:  map[string]interface{}{common.BKDBIN: names},
	}

	count, err := mongodb.Client().Table(types2.BKTableNameBaseNode).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query node failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("create node failed, there are duplicate node names, names: %+v, rid: %s", names, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "names")
	}

	// generate ids field
	ids, err := mongodb.Client().NextSequences(kit.Ctx, types2.BKTableNameBaseNode, len(data))
	if nil != err {
		blog.Errorf("create node failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	result := make([]*types2.Node, 0)
	hasPod := false
	now := time.Now().Unix()
	for idx, node := range data {
		if err := validateNodeData(kit, bizID, node); err != nil {
			return nil, err
		}
		node.NodeBaseFields.HasPod = &hasPod
		node := &types2.Node{
			ID:              int64(ids[idx]),
			BizID:           bizID,
			HostID:          node.HostID,
			ClusterID:       node.ClusterID,
			NodeBaseFields:  node.NodeBaseFields,
			SupplierAccount: kit.SupplierAccount,
			Revision: table.Revision{
				CreateTime: now,
				LastTime:   now,
				Creator:    kit.User,
				Modifier:   kit.User,
			},
		}
		if err := mongodb.Client().Table(types2.BKTableNameBaseNode).Insert(kit.Ctx, node); err != nil {
			blog.Errorf("create node failed, db insert failed, node: %+v, err: %+v, rid: %s", node, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
		}
		result = append(result, node)
	}

	return result, nil
}

// SearchNode find container nodes
func (p *containerOperation) SearchNode(kit *rest.Kit, input *metadata.QueryCondition) (*types2.SearchNodeRsp, error) {

	nodes := make([]types2.Node, 0)
	err := mongodb.Client().Table(types2.BKTableNameBaseNode).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("search node error: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	result := &types2.SearchNodeRsp{Data: nodes}
	return result, nil
}

// SearchCluster find container clusters
func (p *containerOperation) SearchCluster(kit *rest.Kit, input *metadata.QueryCondition) (
	*types2.ResponseCluster, error) {

	clusters := make([]types2.Cluster, 0)
	err := mongodb.Client().Table(types2.BKTableNameBaseCluster).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("search cluster error: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	result := &types2.ResponseCluster{Data: clusters}

	return result, nil
}

// UpdateNodeFields create cluster instance.
func (p *containerOperation) UpdateNodeFields(kit *rest.Kit, bizID int64, supplierAccount string,
	data *types2.UpdateNodeOption) (*metadata.UpdatedCount, errors.CCErrorCoder) {

	filter := map[string]interface{}{
		types2.BKBizIDField:   bizID,
		common.BKOwnerIDField: supplierAccount,
	}
	for _, node := range data.Nodes {

		filter[types2.BKIDField] = map[string]interface{}{
			common.BKDBIN: node.NodeIDs,
		}

		opts := orm.NewFieldOptions().AddIgnoredFields(common.BKFieldID, types2.ClusterUIDField, common.BKFieldName)
		updateData, err := orm.GetUpdateFieldsWithOption(node.Data, opts)
		if err != nil {
			blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", node.Data, err, kit.Rid)
			return &metadata.UpdatedCount{Count: 0}, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}

		err = mongodb.Client().Table(types2.BKTableNameBaseNode).Update(kit.Ctx, filter, updateData)
		if err != nil {
			blog.Errorf("update node failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
				err, kit.Rid)
			return &metadata.UpdatedCount{Count: 0}, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}

	}
	return &metadata.UpdatedCount{Count: uint64(len(data.Nodes))}, nil
}

// UpdateClusterFields create cluster instance.
func (p *containerOperation) UpdateClusterFields(kit *rest.Kit, bizID int64, supplierAccount string,
	data *types2.UpdateClusterOption) (*metadata.UpdatedCount, errors.CCErrorCoder) {

	for _, one := range data.Clusters {
		filter := map[string]interface{}{
			types2.BKBizIDField:   bizID,
			common.BKOwnerIDField: supplierAccount,
		}

		if one.ID != 0 {
			filter[types2.BKIDField] = one.ID
		}

		opts := orm.NewFieldOptions().AddIgnoredFields(common.BKFieldID, types2.ClusterUIDField, common.BKFieldName)
		updateData, err := orm.GetUpdateFieldsWithOption(one.Data, opts)
		if err != nil {
			blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", one, err, kit.Rid)
			return &metadata.UpdatedCount{Count: 0}, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}

		err = mongodb.Client().Table(types2.BKTableNameBaseCluster).Update(kit.Ctx, filter, updateData)
		if err != nil {
			blog.Errorf("update cluster failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
				err, kit.Rid)
			return &metadata.UpdatedCount{Count: 0}, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
	}
	return &metadata.UpdatedCount{Count: uint64(len(data.Clusters))}, nil
}

// CreateCluster create cluster instance.
func (p *containerOperation) CreateCluster(kit *rest.Kit, bizID int64, data *types2.ClusterBaseFields) (*types2.Cluster,
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
				common.BKAppIDField:      bizID,
				common.BKOwnerIDField:    kit.SupplierAccount,
				types2.ContainerUIDField: data.Uid,
			},
			{
				common.BKAppIDField:      bizID,
				common.BKOwnerIDField:    kit.SupplierAccount,
				types2.ContainerUIDField: data.Name,
			},
		},
	}

	count, err := mongodb.Client().Table(types2.BKTableNameBaseCluster).Find(nameFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query cluster failed, filter: %+v, err: %+v, rid: %s", nameFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("create cluster failed, name or uid duplicated, name: %s, uid: %s, rid: %s", data.Name,
			data.Uid, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "name or uid")
	}

	// generate id field
	idTmp, err := mongodb.Client().NextSequence(kit.Ctx, types2.BKTableNameBaseCluster)
	if nil != err {
		blog.Errorf("create cluster failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	now := time.Now().Unix()
	id := int64(idTmp)
	cluster := &types2.Cluster{
		ID:                id,
		BizID:             bizID,
		SupplierAccount:   kit.SupplierAccount,
		ClusterBaseFields: *data,
		Revision: table.Revision{
			CreateTime: now,
			LastTime:   now,
			Creator:    kit.User,
			Modifier:   kit.User,
		},
	}

	if err := mongodb.Client().Table(types2.BKTableNameBaseCluster).Insert(kit.Ctx, cluster); err != nil {
		blog.Errorf("create cluster failed, db insert failed, doc: %+v, err: %+v, rid: %s", cluster, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}

	return cluster, nil
}
