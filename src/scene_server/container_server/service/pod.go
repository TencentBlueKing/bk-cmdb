/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreatePod create a pod
func (s *ContainerService) CreatePod(ctx *rest.Contexts) {

	bkBizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bkBizID, err := util.GetInt64ByInterface(bkBizIDStr)
	if err != nil {
		blog.Warnf("rid: %s, get bk_biz_id failed, invalid value %s, err %s",
			ctx.Kit.Rid, bkBizIDStr, err.Error())
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDataFailed))
		return
	}
	inputData := metadata.CreatePod{}
	if err := ctx.DecodeInto(&inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	inputData.BizID = bkBizID
	resp, err := s.core.PodOperation().CreatePod(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if !resp.Result {
		blog.Warnf("rid: %s CreatePod failed, result %+v", ctx.Kit.Rid, resp)
		ctx.RespErrorCodeOnly(resp.Code, resp.ErrMsg)
		return
	}
	ctx.RespEntity(resp.Data)
}

// CreateManyPod create or update multi pods
func (s *ContainerService) CreateManyPod(ctx *rest.Contexts) {
	bkBizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bkBizID, err := util.GetInt64ByInterface(bkBizIDStr)
	if err != nil {
		blog.Warnf("rid: %s, get bk_biz_id failed, invalid value %s, err %s",
			ctx.Kit.Rid, bkBizIDStr, err.Error())
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDataFailed))
		return
	}
	inputData := metadata.CreateManyPod{}
	if err := ctx.DecodeInto(&inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	inputData.BizID = bkBizID
	if len(inputData.PodList) == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsLostField, "PodList lost")
		return
	}

	resp, err := s.core.PodOperation().CreateManyPod(ctx.Kit, inputData)
	if err != nil {
		blog.Warnf("rid: %s CreateManyPod failed, err %s", ctx.Kit.Rid, err.Error())
		ctx.RespAutoError(err)
		return
	}
	if !resp.Result {
		blog.Warnf("rid: %s CreateManyPod failed, result %+v", ctx.Kit.Rid, resp)
		ctx.RespErrorCodeOnly(resp.Code, resp.ErrMsg)
		return
	}
	ctx.RespEntity(resp.Data)
}

// UpdatePod update a pod
func (s *ContainerService) UpdatePod(ctx *rest.Contexts) {
	bkBizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bkBizID, err := util.GetInt64ByInterface(bkBizIDStr)
	if err != nil {
		blog.Warnf("rid: %s, get bk_biz_id failed, invalid value %s, err %s",
			ctx.Kit.Rid, bkBizIDStr, err.Error())
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDataFailed))
		return
	}
	inputData := metadata.UpdatePod{}
	if err := ctx.DecodeInto(&inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	inputData.BizID = bkBizID
	resp, err := s.core.PodOperation().UpdatePod(ctx.Kit, inputData)
	if err != nil {
		blog.Warnf("rid: %s UpdatePod failed, err %s", ctx.Kit.Rid, err.Error())
		ctx.RespAutoError(err)
		return
	}
	if !resp.Result {
		blog.Warnf("rid: %s UpdatePod failed, result %+v", ctx.Kit.Rid, resp)
		ctx.RespErrorCodeOnly(resp.Code, resp.ErrMsg)
		return
	}
	ctx.RespEntity(resp.Data)
}

// DeletePod delete a pod
func (s *ContainerService) DeletePod(ctx *rest.Contexts) {
	bkBizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bkBizID, err := util.GetInt64ByInterface(bkBizIDStr)
	if err != nil {
		blog.Warnf("rid: %s, get bk_biz_id failed, invalid value %s, err %s",
			ctx.Kit.Rid, bkBizIDStr, err.Error())
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDataFailed))
		return
	}
	inputData := metadata.DeletePod{}
	if err := ctx.DecodeInto(&inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	inputData.BizID = bkBizID

	resp, err := s.core.PodOperation().DeletePod(ctx.Kit, inputData)
	if err != nil {
		blog.Warnf("rid: %s DeletePod failed, err %s", ctx.Kit.Rid, err.Error())
		ctx.RespAutoError(err)
		return
	}
	if !resp.Result {
		blog.Warnf("rid: %s DeletePod failed, result %+v", ctx.Kit.Rid, resp)
		ctx.RespErrorCodeOnly(resp.Code, resp.ErrMsg)
		return
	}
	ctx.RespEntity(resp.Data)
}

// ListPods list pods
func (s *ContainerService) ListPods(ctx *rest.Contexts) {
	bkBizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bkBizID, err := util.GetInt64ByInterface(bkBizIDStr)
	if err != nil {
		blog.Warnf("rid: %s, get bk_biz_id failed, invalid value %s, err %s",
			ctx.Kit.Rid, bkBizIDStr, err.Error())
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDataFailed))
		return
	}
	inputData := metadata.ListPods{}
	if err := ctx.DecodeInto(&inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	inputData.BizID = bkBizID
	resp, err := s.core.PodOperation().ListPods(ctx.Kit, inputData)
	if err != nil {
		blog.Warnf("rid: %s DeletePod failed, err %s", ctx.Kit.Rid, err.Error())
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}
