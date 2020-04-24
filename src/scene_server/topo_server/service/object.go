/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/operation"
)

// CreateObjectBatch batch to create some objects
func (s *Service) CreateObjectBatch(ctx *rest.Contexts) {
	dataWithMetadata := struct {
		Metadata *metadata.Metadata
		Data     map[string]interface{}
	}{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.Core.ObjectOperation().CreateObjectBatch(ctx.Kit, dataWithMetadata.Data, dataWithMetadata.Metadata))
}

// SearchObjectBatch batch to search some objects
func (s *Service) SearchObjectBatch(ctx *rest.Contexts) {
	data := struct {
		operation.ExportObjectCondition `json:",inline"`
		Metadata                        *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Core.ObjectOperation().FindObjectBatch(ctx.Kit, data.ObjIDS, data.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// CreateObject create a new object
func (s *Service) CreateObject(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Core.ObjectOperation().CreateObject(ctx.Kit, false, dataWithMetadata.Data, dataWithMetadata.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	objAudit := operation.NewObjectAudit(s.Engine.CoreAPI, metadata.ModelRes)
	//get CurData
	err = objAudit.MakeCurrent(resp.ToMapStr())
	if err != nil {
		blog.Errorf("[operation-obj] make Current object failed, id: %+v, err: %s, rid: %s", err.Error())
		ctx.RespAutoError(err)
	}

	//package audit response
	err = objAudit.SaveAuditLog(ctx.Kit, metadata.AuditCreate)
	if err != nil {
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(resp.ToMapStr())
}

// SearchObject search some objects by condition
func (s *Service) SearchObject(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	cond := condition.CreateCondition()
	if err := cond.Parse(dataWithMetadata.Data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Core.ObjectOperation().FindObject(ctx.Kit, cond, dataWithMetadata.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchObjectTopo search the object topo
func (s *Service) SearchObjectTopo(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	cond := condition.CreateCondition()
	err := cond.Parse(dataWithMetadata.Data)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoObjectSelectFailed, err.Error()))
		return
	}

	resp, err := s.Core.ObjectOperation().FindObjectTopo(ctx.Kit, cond, dataWithMetadata.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// UpdateObject update the object
func (s *Service) UpdateObject(ctx *rest.Contexts) {
	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s , rid: %s", idStr, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKFieldID))
		return
	}

	objAudit := operation.NewObjectAudit(s.Engine.CoreAPI, metadata.ModelRes)
	//get PreData
	err = objAudit.WithPrevious(ctx.Kit, id)
	if err != nil {
		blog.Errorf("[operation-obj] find Previous object failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//update model
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err = s.Core.ObjectOperation().UpdateObject(ctx.Kit, data, id)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	//get CurData
	err = objAudit.WithCurrent(ctx.Kit, id)
	if err != nil {
		blog.Errorf("[operation-obj] find Current object failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//package audit response
	err = objAudit.SaveAuditLog(ctx.Kit, metadata.AuditUpdate)
	if err != nil {
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(nil)
}

// DeleteObject delete the object
func (s *Service) DeleteObject(ctx *rest.Contexts) {
	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s , rid: %s", idStr, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}

	objAudit := operation.NewObjectAudit(s.Engine.CoreAPI, metadata.ModelRes)

	//get PreData
	err = objAudit.WithPrevious(ctx.Kit, id)
	if err != nil {
		blog.Errorf("[operation-obj] find Previous object failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//delete model
	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err = s.Core.ObjectOperation().DeleteObject(ctx.Kit, id, true, md.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	//package audit response
	err = objAudit.SaveAuditLog(ctx.Kit, metadata.AuditDelete)
	if err != nil {
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(nil)
}

// GetModelStatistics 用于统计各个模型的实例数(Web页面展示需要)
func (s *Service) GetModelStatistics(ctx *rest.Contexts) {
	result, err := s.Engine.CoreAPI.CoreService().Model().GetModelStatistics(ctx.Kit.Ctx, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.Data)
}
