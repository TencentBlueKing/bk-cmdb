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

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// BatchCreatePod batch create nodes
func (s *coreService) BatchCreatePod(ctx *rest.Contexts) {

	inputData := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ContainerOperation().BatchCreatePod(ctx.Kit, inputData.Data))

}

// BatchCreateNode batch create nodes
func (s *coreService) BatchCreateNode(ctx *rest.Contexts) {

	inputData := new(types.CreateNodesOption)
	if err := ctx.DecodeInto(inputData); nil != err {
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

	ctx.RespEntityWithError(s.core.ContainerOperation().BatchCreateNode(ctx.Kit, bizID, inputData.Nodes))
}

// SearchClusters search container clusters
func (s *coreService) SearchClusters(ctx *rest.Contexts) {
	inputData := new(metadata.QueryCondition)

	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ContainerOperation().SearchCluster(ctx.Kit, inputData))
}

// SearchNodes search container nodes
func (s *coreService) SearchNodes(ctx *rest.Contexts) {
	inputData := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := inputData.Page.ValidateLimit(common.BKMaxUpdateOrCreatePageSize); err != nil {
		blog.Errorf("page params is invalid, input: %+v, err: %v, rid: %s", inputData, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	nodes := make([]types.Node, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(inputData.Condition).
		Start(uint64(inputData.Page.Start)).
		Limit(uint64(inputData.Page.Limit)).
		Sort(inputData.Page.Sort).
		Fields(inputData.Fields...).All(ctx.Kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("search nodes failed, input %+v, err: %v, rid: %s", inputData, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.SearchNodeRsp{Data: nodes}
	ctx.RespEntityWithError(result, nil)
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

	supplierAccount := ctx.Request.PathParameter("supplierAccount")
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKOwnerIDField))
		return
	}

	ctx.RespEntityWithError(s.core.ContainerOperation().BatchUpdateNodeFields(ctx.Kit, bizID, supplierAccount, input))
}

// BatchUpdateCluster update cluster.
func (s *coreService) BatchUpdateCluster(ctx *rest.Contexts) {

	input := new(types.UpdateClusterOption)

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

	supplierAccount := ctx.Request.PathParameter("supplierAccount")
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKOwnerIDField))
		return
	}
	result, err := s.core.ContainerOperation().BatchUpdateClusterFields(ctx.Kit, bizID, supplierAccount, input)
	ctx.RespEntityWithError(result, err)
}

// CreateCluster create cluster instance.
func (s *coreService) CreateCluster(ctx *rest.Contexts) {

	inputData := new(types.Cluster)

	if err := ctx.DecodeInto(inputData); nil != err {
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

	ctx.RespEntityWithError(s.core.ContainerOperation().CreateCluster(ctx.Kit, bizID, inputData))
}

// BatchDeleteCluster delete clusters.
func (s *coreService) BatchDeleteCluster(ctx *rest.Contexts) {
	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); nil != err {
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

	ctx.RespEntityWithError(s.core.ContainerOperation().DeleteCluster(ctx.Kit, bizID, option))
}

// BatchDeleteNode delete cluster instance.
func (s *coreService) BatchDeleteNode(ctx *rest.Contexts) {
	option := new(types.BatchDeleteNodeOption)
	if err := ctx.DecodeInto(option); nil != err {
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

	ctx.RespEntityWithError(s.core.ContainerOperation().BatchDeleteNode(ctx.Kit, bizID, option))
}

// ListContainer list container
func (s *coreService) ListContainer(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	containers := make([]types.Container, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseContainer).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &containers)
	if err != nil {
		blog.Errorf("search container failed, cond: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.ContainerDataResp{Info: containers}
	ctx.RespEntity(result)
}
