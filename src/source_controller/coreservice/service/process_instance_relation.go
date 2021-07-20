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
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (s *coreService) CreateProcessInstanceRelation(ctx *rest.Contexts) {
	relation := &metadata.ProcessInstanceRelation{}
	if err := ctx.DecodeInto(relation); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateProcessInstanceRelation(ctx.Kit, relation)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) CreateProcessInstanceRelations(ctx *rest.Contexts) {
	relations := make([]*metadata.ProcessInstanceRelation, 0)
	if err := ctx.DecodeInto(&relations); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateProcessInstanceRelations(ctx.Kit, relations)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelations failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) GetProcessInstanceRelation(ctx *rest.Contexts) {
	processInstanceIDStr := ctx.Request.PathParameter(common.BKProcIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("GetProcessInstanceRelation failed, path parameter `%s` empty, rid: %s", common.BKProcIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcIDField))
		return
	}

	serviceTemplateID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcIDField, processInstanceIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcIDField))
		return
	}

	result, err := s.core.ProcessOperation().GetProcessInstanceRelation(ctx.Kit, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListProcessInstanceRelation(ctx *rest.Contexts) {
	// filter parameter
	fp := metadata.ListProcessInstanceRelationOption{}

	if err := ctx.DecodeInto(&fp); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListProcessInstanceRelation failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	result, err := s.core.ProcessOperation().ListProcessInstanceRelation(ctx.Kit, fp)
	if err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateProcessInstanceRelation(ctx *rest.Contexts) {
	processInstanceIDStr := ctx.Request.PathParameter(common.BKProcIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("UpdateProcessInstanceRelation failed, path parameter `%s` empty, rid: %s", common.BKProcIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcIDField))
		return
	}

	processInstanceID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcIDField, processInstanceIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcIDField))
		return
	}

	relation := metadata.ProcessInstanceRelation{}
	if err := ctx.DecodeInto(&relation); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().UpdateProcessInstanceRelation(ctx.Kit, processInstanceID, relation)
	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) DeleteProcessInstanceRelation(ctx *rest.Contexts) {
	option := metadata.DeleteProcessInstanceRelationOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.ProcessOperation().DeleteProcessInstanceRelation(ctx.Kit, option); err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) CreateProcessInstance(kit *rest.Kit, process *metadata.Process) (*metadata.Process, errors.CCErrorCoder) {
	processBytes, err := json.Marshal(process)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommJsonEncode)
	}
	mData := mapstr.MapStr{}
	if err := json.Unmarshal(processBytes, &mData); nil != err && 0 != len(processBytes) {
		return nil, kit.CCError.CCError(common.CCErrCommJsonDecode)
	}
	inputParam := metadata.CreateModelInstance{
		Data: mData,
	}
	result, err := s.core.InstanceOperation().CreateModelInstance(kit, common.BKProcessObjectName, inputParam)
	if err != nil {
		blog.Errorf("CreateProcessInstance failed, CreateModelInstance failed, inputParam: %+v, err: %+v, rid: %s", inputParam, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrProcCreateProcessFailed)
	}
	process.ProcessID = int64(result.Created.ID)
	return process, nil
}

func (s *coreService) CreateProcessInstances(kit *rest.Kit, processes []*metadata.Process) ([]*metadata.Process, errors.CCErrorCoder) {
	processesBytes, err := json.Marshal(processes)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommJsonEncode)
	}
	mData := []mapstr.MapStr{}
	if err := json.Unmarshal(processesBytes, &mData); nil != err && 0 != len(processesBytes) {
		return nil, kit.CCError.CCError(common.CCErrCommJsonDecode)
	}
	inputParam := metadata.CreateManyModelInstance{
		Datas: mData,
	}
	result, err := s.core.InstanceOperation().CreateManyModelInstance(kit, common.BKProcessObjectName, inputParam)
	if err != nil {
		blog.Errorf("CreateProcessInstances failed, CreateManyModelInstance failed, inputParam: %#v, err: %v, rid: %s", inputParam, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrProcCreateProcessFailed)
	}
	if len(processes) != len(result.Created) {
		blog.Errorf("CreateProcessInstances failed, len(processes) != len(result.Created), inputParam: %#v, rid: %s", inputParam, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrProcCreateProcessFailed)
	}

	for idx, created := range result.Created {
		processes[idx].ProcessID = int64(created.ID)
	}

	return processes, nil
}
