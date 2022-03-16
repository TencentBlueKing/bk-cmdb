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
)

func (s *coreService) CreateHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.CreateHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.HostApplyRuleOperation().CreateHostApplyRule(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	ruleIDStr := ctx.Request.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField))
		return
	}

	option := metadata.UpdateHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.HostApplyRuleOperation().UpdateHostApplyRule(ctx.Kit, bizID, ruleID, option)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, ruleID: %d, option: %+v, err: %+v, rid: %s", ruleID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) DeleteHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.DeleteHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.HostApplyRuleOperation().DeleteHostApplyRule(ctx.Kit, bizID, option.RuleIDs...); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, bizID: %d, ruleID: %d, err: %+v, rid: %s", bizID, option.RuleIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) GetHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	hostApplyRuleIDStr := ctx.Request.PathParameter(common.HostApplyRuleIDField)
	hostApplyRuleID, err := strconv.ParseInt(hostApplyRuleIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField))
		return
	}

	rule, err := s.core.HostApplyRuleOperation().GetHostApplyRule(ctx.Kit, bizID, hostApplyRuleID)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, bizID: %d, ruleID: %d, err: %+v, rid: %s", bizID, hostApplyRuleID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(rule)
}

func (s *coreService) ListHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	hostApplyRuleResult, err := s.core.HostApplyRuleOperation().ListHostApplyRule(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(hostApplyRuleResult)
}

func (s *coreService) GenerateApplyPlan(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.HostApplyPlanOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	applyPlans, err := s.core.HostApplyRuleOperation().GenerateApplyPlan(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(applyPlans)
}

func (s *coreService) SearchRuleRelatedModules(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.SearchRuleRelatedModulesOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	modules, err := s.core.HostApplyRuleOperation().SearchRuleRelatedModules(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(modules)
}

func (s *coreService) BatchUpdateHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.BatchCreateOrUpdateApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.HostApplyRuleOperation().BatchUpdateHostApplyRule(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("BatchUpdateHostApplyRule failed, option: %+v, err: %+v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateHostByHostApplyRule(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.UpdateHostByHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.HostApplyRuleOperation().RunHostApplyOnHosts(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("UpdateHostByHostApplyRule failed, RunHostApplyOnHosts failed, option: %+v, err: %+v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}
