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

package kube

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// updateNodeHasPodField here you need to update the has_pod in the node uniformly
func (s *service) updateNodeHasPodField(kit *rest.Kit, nodeIDs []int64) error {
	if len(nodeIDs) == 0 {
		return nil
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
func (s *service) BatchCreateNode(ctx *rest.Contexts) {
	inputData := make([]types.OneNodeCreateOption, 0)
	if err := ctx.DecodeInto(&inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	for _, data := range inputData {
		if err := data.ValidateCreate(); err.ErrCode != 0 {
			blog.Errorf("node %+v is invalid, err: %v, rid: %s", data, err, ctx.Kit.Rid)
			ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
			return
		}
	}

	nodes, err := s.core.KubeOperation().BatchCreateNode(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nodes)
}

// SearchNodes search nodes
func (s *service) SearchNodes(ctx *rest.Contexts) {
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
func (s *service) BatchUpdateNode(ctx *rest.Contexts) {
	input := new(types.UpdateNodeByIDsOption)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := input.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := map[string]interface{}{
		types.BKIDField: map[string]interface{}{
			common.BKDBIN: input.IDs,
		},
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
func (s *service) BatchDeleteNode(ctx *rest.Contexts) {
	option := new(types.BatchDeleteNodeByIDsOption)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	// obtain the hostID of the deleted node and the corresponding business ID.
	query := map[string]interface{}{
		types.BKIDField: map[string]interface{}{common.BKDBIN: option.IDs},
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
