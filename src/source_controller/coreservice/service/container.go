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
)

// BatchCreatePod batch create nodes
func (s *coreService) BatchCreatePod(ctx *rest.Contexts) {

	inputData := new(types.CreatePodsOption)
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

	ctx.RespEntityWithError(s.core.ContainerOperation().BatchCreatePod(ctx.Kit, bizID, inputData.Pods))

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

// SearchClusterInstances 查找集群实例
func (s *coreService) SearchClusterInstances(ctx *rest.Contexts) {
	inputData := new(metadata.QueryCondition)

	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ContainerOperation().SearchCluster(ctx.Kit, inputData))
}

// SearchNodeInstances 查找集群实例
func (s *coreService) SearchNodeInstances(ctx *rest.Contexts) {
	inputData := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ContainerOperation().SearchNode(ctx.Kit, inputData))
}

// UpdateNodeInstance update node instance.
func (s *coreService) UpdateNodeInstance(ctx *rest.Contexts) {

	inputData := new(types.UpdateNodeOption)
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

	supplierAccount := ctx.Request.PathParameter("supplierAccount")
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKOwnerIDField))
		return
	}

	ctx.RespEntityWithError(s.core.ContainerOperation().UpdateNodeFields(ctx.Kit, bizID, supplierAccount, inputData))
}

// UpdateClusterInstance update cluster instance.
func (s *coreService) UpdateClusterInstance(ctx *rest.Contexts) {

	inputData := new(types.UpdateClusterOption)

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

	supplierAccount := ctx.Request.PathParameter("supplierAccount")
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKOwnerIDField))
		return
	}

	ctx.RespEntityWithError(s.core.ContainerOperation().UpdateClusterFields(ctx.Kit, bizID, supplierAccount, inputData))
}

// CreateClusterInstance create cluster instance.
func (s *coreService) CreateClusterInstance(ctx *rest.Contexts) {

	inputData := new(types.ClusterBaseFields)

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

// DeleteClusterInstance delete cluster instance.
func (s *coreService) DeleteClusterInstance(ctx *rest.Contexts) {
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

// BatchDeleteNodeInstance delete cluster instance.
func (s *coreService) BatchDeleteNodeInstance(ctx *rest.Contexts) {
	option := new(types.ArrangeDeleteNodeOption)
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
