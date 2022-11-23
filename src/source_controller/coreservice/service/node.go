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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
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

// BatchCreateNode batch create nodes
func (s *coreService) BatchCreateNode(ctx *rest.Contexts) {

	inputData := new(types.CreateNodesOption)
	if err := ctx.DecodeInto(inputData); err != nil {
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
	nodes, err := s.core.KubeOperation().BatchCreateNode(ctx.Kit, bizID, inputData.Nodes)
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
	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(query).
		Fields(common.BKHostIDField, common.BKAppIDField).All(ctx.Kit.Ctx, &nodes); err != nil {
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

	ctx.RespEntity(nil)
}
