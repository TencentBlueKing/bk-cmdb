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
 * an "AS IS" BASIS, WITHOUT WARRAcommon.BKOwnerIDField: supplierAccountNTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"errors"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// updateNodeField here you need to update the has_pod in the node uniformly
func (s *coreService) updateNodeField(kit *rest.Kit, nodeIDMap map[int64]struct{}) error {

	if len(nodeIDMap) == 0 {
		return nil
	}

	nodeIDs := make([]int64, 0)
	for id := range nodeIDMap {
		nodeIDs = append(nodeIDs, id)
	}

	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: nodeIDs,
		},
	}

	updateData := map[string]interface{}{
		types.HasPodField: true,
	}
	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Update(kit.Ctx, filter, updateData); err != nil {
		blog.Errorf("update node has_pod field failed, filter: %v, err: %+v, rid: %s", filter, err, kit.Rid)
		return err
	}
	return nil
}

func getClusterSpecInfo(kit *rest.Kit, bizID int64, data []types.OneNodeCreateOption) (
	map[int64]types.ClusterSpec, ccErr.CCErrorCoder) {

	clusterIDs := make([]int64, 0)
	for _, info := range data {
		clusterIDs = append(clusterIDs, info.ClusterID)
	}

	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		types.BKIDField:     map[string]interface{}{common.BKDBIN: clusterIDs},
	}
	util.SetModOwner(filter, kit.SupplierAccount)
	clusters := make([]types.Cluster, 0)
	fields := []string{types.UidField, types.BKIDField, types.TypeField, types.BKAsstBizIDField, common.BKAppIDField}
	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(filter).
		Fields(fields...).All(kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("query cluster failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(clusters) == 0 {
		blog.Errorf("no cluster founded, filter: %+v,  rid:%s", filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}

	clusterMap := make(map[int64]types.ClusterSpec)
	for _, cluster := range clusters {
		if cluster.Uid == nil || cluster.Type == nil {
			blog.Errorf("query cluster uid or type failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed, "cluster uid or type")
		}
		// 如果集群中的bizID与请求中的bizID不一致，那么此集群必定为共享集群
		if cluster.BizID != bizID && *cluster.Type != types.ClusterShareTypeField {
			blog.Errorf("bizID(%d) in the request is inconsistent with the bizID(%d) in the cluster, "+
				"and the cluster type must be a shared cluster, type is %s, filter: %+v, err: %+v, rid: %s", bizID,
				cluster.BizID, *cluster.Type, filter, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed, errors.New("cluster must be share type"))
		}

		clusterMap[cluster.ID] = types.ClusterSpec{
			BizID:       bizID,
			ClusterUID:  *cluster.Uid,
			ClusterID:   cluster.ID,
			ClusterType: *cluster.Type,
			BizAsstID:   cluster.BizID,
		}
	}
	return clusterMap, nil
}

// validateNodeData 目前的逻辑是node所在对应的host是一定需要在cc中的。
func validateNodeData(kit *rest.Kit, bizID int64, hostIDs []int64) ccErr.CCErrorCoder {

	cond := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	util.SetModOwner(cond, kit.SupplierAccount)
	cnt, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Distinct(kit.Ctx, common.BKHostIDField, cond)
	if err != nil {
		blog.Errorf("query host module config failed, err: %s, rid:%s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	hostIDNum, cntNum := len(hostIDs), len(cnt)
	if cntNum != hostIDNum {
		blog.Errorf("hostID num not as expected, filter: %+v, cnt: %d, hostIDs num: %+v, rid:%s", cond,
			cntNum, hostIDNum, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsIsInvalid)
	}

	return nil
}

// batchCreateNode create container node data in batches.
func batchCreateNode(kit *rest.Kit, bizID int64, data []types.OneNodeCreateOption) (
	[]types.Node, ccErr.CCErrorCoder) {

	hostIDMap := make(map[int64]struct{})
	for _, node := range data {
		hostIDMap[node.HostID] = struct{}{}
	}

	hostIDs := make([]int64, 0)
	for id := range hostIDMap {
		hostIDs = append(hostIDs, id)
	}

	if err := validateNodeData(kit, bizID, hostIDs); err != nil {
		return nil, err
	}

	clusterMap, cErr := getClusterSpecInfo(kit, bizID, data)
	if cErr != nil {
		return nil, cErr
	}

	// generate ids field
	ids, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBaseNode, len(data))
	if err != nil {
		blog.Errorf("create node failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	result := make([]types.Node, 0)
	now := time.Now().Unix()
	noPod := false

	for idx, node := range data {
		node := types.Node{
			ID:               int64(ids[idx]),
			ClusterSpec:      clusterMap[node.ClusterID],
			HostID:           node.HostID,
			Name:             node.Name,
			Roles:            node.Roles,
			Labels:           node.Labels,
			Taints:           node.Taints,
			Unschedulable:    node.Unschedulable,
			InternalIP:       node.InternalIP,
			ExternalIP:       node.ExternalIP,
			HasPod:           &noPod,
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
		result = append(result, node)
	}

	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Insert(kit.Ctx, result); err != nil {
		blog.Errorf("create nodes failed, db insert failed, node: %+v, err: %+v, rid: %s", result, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}
	blog.Errorf("000000000000000000000 result: %+v", result)
	return result, nil
}

// BatchCreateNode batch create nodes
func (s *coreService) BatchCreateNode(ctx *rest.Contexts) {

	inputData := new(types.CreateNodesOption)
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizStr := ctx.Request.PathParameter("bk_biz_id")
	bizID, err := strconv.ParseInt(bizStr, 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", bizStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}
	nodes, err := batchCreateNode(ctx.Kit, bizID, inputData.Nodes)

	ctx.RespEntityWithError(nodes, err)
}

// SearchNodes search nodes
func (s *coreService) SearchNodes(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)

	nodes := make([]types.Node, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(input.Condition).
		Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("search nodes failed, input %+v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.SearchNodeRsp{Data: nodes}
	ctx.RespEntity(result)
}

// BatchUpdateNode batch update node.
func (s *coreService) BatchUpdateNode(ctx *rest.Contexts) {

	input := new(types.UpdateNodeOption)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	filter := map[string]interface{}{
		types.BKBizIDField: bizID,
	}
	filter[types.BKIDField] = map[string]interface{}{
		common.BKDBIN: input.IDs,
	}
	util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	opts := orm.NewFieldOptions().AddIgnoredFields(types.IgnoredUpdateNodeFields...)
	updateData, err := orm.GetUpdateFieldsWithOption(input.Data, opts)
	if err != nil {
		blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", input.Data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	err = mongodb.Client().Table(types.BKTableNameBaseNode).Update(ctx.Kit.Ctx, filter, updateData)
	if err != nil {
		blog.Errorf("update node failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// BatchDeleteNode batch delete nodes.
func (s *coreService) BatchDeleteNode(ctx *rest.Contexts) {
	option := new(types.BatchDeleteNodeOption)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	bizStr := ctx.Request.PathParameter("bk_biz_id")
	bizID, err := strconv.ParseInt(bizStr, 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", bizStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	// obtain the hostID of the deleted node and the corresponding business ID.
	query := map[string]interface{}{
		types.BKIDField:     map[string]interface{}{common.BKDBIN: option.IDs},
		common.BKAppIDField: bizID,
	}
	util.SetQueryOwner(query, ctx.Kit.SupplierAccount)
	nodes := make([]types.Node, 0)
	fields := []string{common.BKFieldID, types.BKClusterIDField, types.BKAsstBizIDField, types.BKBizIDField}

	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(query).
		Fields(fields...).All(ctx.Kit.Ctx, &nodes); err != nil {
		blog.Errorf("query node failed, filter: %+v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	hostBizMap := make(map[int64]int64)
	for _, node := range nodes {
		hostBizMap[node.HostID] = node.BizID
	}

	// delete nodes.
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		types.BKIDField: map[string]interface{}{
			common.BKDBIN: option.IDs,
		},
	}
	util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	delConds := make([]map[string]interface{}, 0)
	for _, node := range nodes {
		delConds = append(delConds, map[string]interface{}{
			types.BKNamespaceIDField: node.ID,
			types.BKAsstBizIDField:   node.BizAsstID,
			types.BKBizIDField:       node.BizID,
			types.BKClusterIDField:   node.ClusterID,
		})
	}

	cond := map[string]interface{}{common.BKDBOR: delConds}
	cond = util.SetModOwner(cond, ctx.Kit.SupplierAccount)
	if err := mongodb.Client().Table(types.BKTableNodeClusterRelation).Delete(ctx.Kit.Ctx, cond); err != nil {
		blog.Errorf("delete node relation failed, cond: %+v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	ctx.RespEntity(nil)
}
