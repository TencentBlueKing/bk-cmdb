/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
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

// BatchCreateNode batch create nodes
func (s *coreService) BatchCreateNode(ctx *rest.Contexts) {

	inputData := new(types.CreateNodesReq)
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
	blog.Errorf("0000000000000000000 inputData: %+v", inputData)

	ctx.RespEntityWithError(s.core.ContainerOperation().SearchCluster(ctx.Kit, inputData))
}

// SearchNodeInstances 查找集群实例
func (s *coreService) SearchNodeInstances(ctx *rest.Contexts) {
	inputData := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	blog.Errorf("0000000000000000000 inputData: %+v", inputData)

	ctx.RespEntityWithError(s.core.ContainerOperation().SearchNode(ctx.Kit, inputData))
}

// CreateClusterInstance create cluster instance.
func (s *coreService) CreateClusterInstance(ctx *rest.Contexts) {
	inputData := new(types.ClusterBaseFields)

	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	blog.Errorf("0000000000000000000 Name: %+v, Uid: %v, NetworkType: %v", *inputData.Name, *inputData.Uid, *inputData.NetworkType)

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
