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
	"configcenter/src/kube/orm"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

type containerOperation struct {
}

// New create a new model manager instance
func New() core.ContainerOperation {
	container := &containerOperation{}
	return container
}

// BatchDeleteNode delete node instances.
func (p *containerOperation) BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.ArrangeDeleteNodeOption) (
	*metadata.DeletedCount, errors.CCErrorCoder) {

	num := 0
	if option.Flag {
		for uid, names := range option.Option {
			filter := map[string]interface{}{
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: kit.SupplierAccount,
				types.ClusterUIDField: uid,
				types.NodeField: map[string]interface{}{
					common.BKDBIN: names,
				},
			}
			if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Delete(kit.Ctx, filter); err != nil {
				blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
				return nil, kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
			}
			num += len(names)
		}
		return &metadata.DeletedCount{Count: uint64(num)}, nil
	}

	for clusterID, ids := range option.Option {
		filter := map[string]interface{}{
			common.BKAppIDField:    bizID,
			common.BKOwnerIDField:  kit.SupplierAccount,
			types.BKClusterIDFiled: clusterID,
			types.BKNodeIDField: map[string]interface{}{
				common.BKDBIN: ids,
			},
		}
		if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Delete(kit.Ctx, filter); err != nil {
			blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
		}
		num += len(ids)
	}
	return &metadata.DeletedCount{Count: uint64(num)}, nil
}

// DeleteCluster delete cluster instance.
func (p *containerOperation) DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption) (
	*metadata.DeletedCount, errors.CCErrorCoder) {

	filter := make(map[string]interface{}, 0)
	num := 0
	if len(option.IDs) > 0 {
		num = len(option.IDs)
		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: kit.SupplierAccount,
			types.ClusterUIDField: map[string]interface{}{
				common.BKDBIN: option.Uids,
			},
		}
	}

	if len(option.Uids) > 0 {
		num = len(option.Uids)

		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: kit.SupplierAccount,
			types.BKIDField: map[string]interface{}{
				common.BKDBIN: option.IDs,
			},
		}
	}
	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Delete(kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}
	return &metadata.DeletedCount{Count: uint64(num)}, nil
}

func validateNodeData(kit *rest.Kit, bizID int64, node types.NodeBaseFields) errors.CCErrorCoder {

	clusterFilter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		types.BKIDField:       *node.ClusterID,
		types.UidField:        *node.ClusterUID,
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

// BatchCreateNode create container node data in batches.
func (p *containerOperation) BatchCreateNode(kit *rest.Kit, bizID int64, data []types.NodeBaseFields) ([]int64,
	errors.CCErrorCoder) {

	names := make([]string, 0)
	for _, node := range data {
		names = append(names, *node.Name)
	}

	filter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		types.KubeNameField:   map[string]interface{}{common.BKDBIN: names},
	}

	count, err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query node failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("create node failed, there are duplicate node names, names: %+v, rid: %s", names, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "names")
	}

	// generate ids field
	ids, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBaseNode, len(data))
	if nil != err {
		blog.Errorf("create node failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	nodeIDs := make([]int64, 0)

	now := time.Now().Unix()
	for idx, node := range data {

		if err := validateNodeData(kit, bizID, node); err != nil {
			return nil, err
		}
		node := &types.Node{
			ID:              int64(ids[idx]),
			BizID:           bizID,
			NodeBaseFields:  node,
			SupplierAccount: kit.SupplierAccount,
			Revision: table.Revision{
				CreateTime: now,
				LastTime:   now,
				Creator:    kit.User,
				Modifier:   kit.User,
			},
		}
		nodeIDs = append(nodeIDs, int64(ids[idx]))
		if err := mongodb.Client().Table(types.BKTableNameBaseNode).Insert(kit.Ctx, node); err != nil {
			blog.Errorf("create node failed, db insert failed, node: %+v, err: %+v, rid: %s", node, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
		}
	}

	return nodeIDs, nil
}

// SearchNode 查找容器节点
func (p *containerOperation) SearchNode(kit *rest.Kit, input *metadata.QueryCondition) (*types.ResponseNode, error) {

	nodes := make([]types.Node, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("search node error: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	result := &types.ResponseNode{Data: nodes}
	return result, nil
}

// SearchCluster 查找容器集群
func (p *containerOperation) SearchCluster(kit *rest.Kit, input *metadata.QueryCondition) (
	*types.ResponseCluster, error) {

	clusters := make([]types.Cluster, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("search cluster error: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	result := &types.ResponseCluster{Data: clusters}

	return result, nil
}

// CreateContainer create pod instance.
func (p *containerOperation) CreateContainer(kit *rest.Kit, bizID int64, data *types.ContainerCoreInfo) (int64,
	errors.CCErrorCoder) {
	return 0, nil
}

// CreatePod create pod instance.
func (p *containerOperation) CreatePod(kit *rest.Kit, bizID int64, data *types.PodCoreInfo) (int64,
	errors.CCErrorCoder) {
	return 0, nil
}

// UpdateClusterFields create cluster instance.
func (p *containerOperation) UpdateClusterFields(kit *rest.Kit, bizID int64, supplierAccount string,
	data *types.UpdateClusterOption) (*metadata.UpdatedCount, errors.CCErrorCoder) {

	for _, one := range data.Cluster {
		filter := make(map[string]interface{})
		if one.ID != 0 {
			filter = map[string]interface{}{
				types.BKIDField:       one.ID,
				types.BKBizIDField:    bizID,
				common.BKOwnerIDField: supplierAccount,
			}
		}
		if one.UID != "" {
			filter = map[string]interface{}{
				types.UidField:        one.UID,
				types.BKBizIDField:    bizID,
				common.BKOwnerIDField: supplierAccount,
			}
		}
		opts := orm.NewFieldOptions().AddIgnoredFields(common.BKFieldID, types.ClusterUIDField, common.BKFieldName)
		updateData, err := orm.GetUpdateFieldsWithOption(one, opts)
		if err != nil {
			blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", one, err, kit.Rid)

			return &metadata.UpdatedCount{Count: 0}, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
		err = mongodb.Client().Table(types.BKTableNameBaseCluster).Update(kit.Ctx, filter, updateData)
		if err != nil {
			blog.Errorf("update namespace failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
				err, kit.Rid)
			return &metadata.UpdatedCount{Count: 0}, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
	}
	return &metadata.UpdatedCount{Count: uint64(len(data.Cluster))}, nil
}

// CreateCluster create cluster instance.
func (p *containerOperation) CreateCluster(kit *rest.Kit, bizID int64, data *types.ClusterBaseFields) (int64,
	errors.CCErrorCoder) {

	// it is necessary to judge whether there is duplicate data here, to prevent subsequent calls to coreservice
	// directly and lack of verification.
	if err := data.ValidateCreate(); err != nil {
		blog.Errorf("create cluster failed, data: %+v, err: %+v, rid: %s", data, err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	nameFilter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKAppIDField:     bizID,
				common.BKOwnerIDField:   kit.SupplierAccount,
				types.ContainerUIDField: data.Uid,
			},
			{
				common.BKAppIDField:     bizID,
				common.BKOwnerIDField:   kit.SupplierAccount,
				types.ContainerUIDField: data.Name,
			},
		},
	}

	count, err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(nameFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query cluster failed, filter: %+v, err: %+v, rid: %s", nameFilter, err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("create cluster failed, name or uid duplicated, name: %s, uid: %s, rid: %s", data.Name,
			data.Uid, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "name or uid")
	}

	// generate id field
	idTmp, err := mongodb.Client().NextSequence(kit.Ctx, types.BKTableNameBaseCluster)
	if nil != err {
		blog.Errorf("create cluster failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	now := time.Now().Unix()
	id := int64(idTmp)
	cluster := &types.Cluster{
		ID:                &id,
		BizID:             &bizID,
		SupplierAccount:   &kit.SupplierAccount,
		ClusterBaseFields: *data,
		Revision: table.Revision{
			CreateTime: now,
			LastTime:   now,
			Creator:    kit.User,
			Modifier:   kit.User,
		},
	}

	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Insert(kit.Ctx, cluster); err != nil {
		blog.Errorf("create cluster failed, db insert failed, doc: %+v, err: %+v, rid: %s", cluster, err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}

	return *cluster.ID, nil
}
