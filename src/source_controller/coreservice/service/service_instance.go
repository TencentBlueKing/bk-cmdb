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
	"configcenter/src/common/metadata"
)

func (s *coreService) CreateServiceInstance(ctx *rest.Contexts) {
	instance := new(metadata.ServiceInstance)
	if err := ctx.DecodeInto(instance); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateServiceInstance(ctx.Kit, instance)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) CreateServiceInstances(ctx *rest.Contexts) {
	instances := []*metadata.ServiceInstance{}
	if err := ctx.DecodeInto(&instances); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateServiceInstances(ctx.Kit, instances)
	if err != nil {
		blog.Errorf("CreateServiceInstances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ConstructServiceInstanceName(ctx *rest.Contexts) {
	params := new(metadata.SrvInstNameParams)
	if err := ctx.DecodeInto(params); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.ProcessOperation().ConstructServiceInstanceName(ctx.Kit, params.ServiceInstanceID, params.Host, params.Process); err != nil {
		blog.Errorf("ConstructServiceInstanceName failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) ReconstructServiceInstanceName(ctx *rest.Contexts) {
	serviceInstanceIDStr := ctx.Request.PathParameter(common.BKServiceInstanceIDField)
	if len(serviceInstanceIDStr) == 0 {
		blog.Errorf("GetServiceInstance failed, path parameter `%s` empty, rid: %s", common.BKServiceInstanceIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField))
		return
	}

	serviceInstanceID, err := strconv.ParseInt(serviceInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceInstance failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceInstanceIDField, serviceInstanceIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField))
		return
	}

	if err := s.core.ProcessOperation().ReconstructServiceInstanceName(ctx.Kit, serviceInstanceID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) GetServiceInstance(ctx *rest.Contexts) {
	serviceInstanceIDStr := ctx.Request.PathParameter(common.BKServiceInstanceIDField)
	if len(serviceInstanceIDStr) == 0 {
		blog.Errorf("GetServiceInstance failed, path parameter `%s` empty, rid: %s", common.BKServiceInstanceIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField))
		return
	}

	serviceInstanceID, err := strconv.ParseInt(serviceInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceInstance failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceInstanceIDField, serviceInstanceIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField))
		return
	}

	result, err := s.core.ProcessOperation().GetServiceInstance(ctx.Kit, serviceInstanceID)
	if err != nil {
		blog.Errorf("GetServiceInstance failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListServiceInstances(ctx *rest.Contexts) {
	// filter parameter
	fp := metadata.ListServiceInstanceOption{}

	if err := ctx.DecodeInto(&fp); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceTemplates failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	result, err := s.core.ProcessOperation().ListServiceInstance(ctx.Kit, fp)
	if err != nil {
		blog.Errorf("ListServiceInstance failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListServiceInstanceDetail(ctx *rest.Contexts) {
	// filter parameter
	fp := metadata.ListServiceInstanceDetailOption{}

	if err := ctx.DecodeInto(&fp); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceInstanceDetail failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	result, err := s.core.ProcessOperation().ListServiceInstanceDetail(ctx.Kit, fp)
	if err != nil {
		blog.Errorf("ListServiceInstanceDetail failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateServiceInstances(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	if len(bizIDStr) == 0 {
		blog.Errorf("UpdateServiceInstances failed, path parameter `%s` empty, rid: %s", common.BKServiceInstanceIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField))
		return
	}

	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceInstances failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s",
			common.BKServiceInstanceIDField, bizIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField))
		return
	}

	option := new(metadata.UpdateServiceInstanceOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.ProcessOperation().UpdateServiceInstances(ctx.Kit, bizID, option); err != nil {
		blog.Errorf("UpdateServiceInstances failed, err: %+v, option:%#v, rid: %s", err, *option, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) DeleteServiceInstance(ctx *rest.Contexts) {
	option := metadata.CoreDeleteServiceInstanceOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.ProcessOperation().DeleteServiceInstance(ctx.Kit, option.ServiceInstanceIDs); err != nil {
		blog.Errorf("DeleteServiceInstance failed, err: %+v, rid: %s", err, common.BKServiceInstanceIDField)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) GetBusinessDefaultSetModuleInfo(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	if len(bizIDStr) == 0 {
		blog.Errorf("GetBusinessDefaultSetModuleInfo failed, path parameter `%s` empty, rid: %s", common.BKAppIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetBusinessDefaultSetModuleInfo failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKAppIDField, bizIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	defaultSetModuleInfo, err := s.core.ProcessOperation().GetBusinessDefaultSetModuleInfo(ctx.Kit, bizID)
	if err != nil {
		blog.Errorf("GetBusinessDefaultSetModuleInfo failed, bizID: %d, err: %+v, rid: %s", bizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(defaultSetModuleInfo)
}

// AutoCreateServiceInstanceModuleHost is dependence for host
func (s *coreService) AutoCreateServiceInstanceModuleHost(kit *rest.Kit, hostIDs []int64, moduleIDs []int64) errors.CCErrorCoder {
	err := s.core.ProcessOperation().AutoCreateServiceInstanceModuleHost(kit, hostIDs, moduleIDs)
	if err != nil {
		blog.Errorf("AutoCreateServiceInstanceModuleHost failed, hostID: %+v, moduleID: %+v, err: %+v, rid: %s", hostIDs, moduleIDs, err, kit.Rid)
		return err
	}
	return nil
}

func (s *coreService) RemoveTemplateBindingOnModule(ctx *rest.Contexts) {
	moduleIDStr := ctx.Request.PathParameter(common.BKModuleIDField)
	moduleID, err := strconv.ParseInt(moduleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("RemoveTemplateBindingOnModule failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKAppIDField, moduleIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField))
		return
	}

	if err := s.core.ProcessOperation().RemoveTemplateBindingOnModule(ctx.Kit, moduleID); err != nil {
		blog.Errorf("RemoveTemplateBindingOnModule failed, moduleID: %d, err: %+v, rid: %s", moduleID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
