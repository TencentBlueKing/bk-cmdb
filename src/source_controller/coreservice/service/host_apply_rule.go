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
	"configcenter/src/storage/driver/mongodb"
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

// updateModuleHostApplyStatus after judging the deletion of the module rule, whether there is a corresponding host
// automatic application rule in the cc_HostApplyRule table, if not, the host automatic application state corresponding
// to this module needs to be turned off.
func (s *coreService) updateModuleHostApplyStatus(kit *rest.Kit, bizID int64, moduleIDs []int64) error {

	modIDs := make([]int64, 0)

	// 判断每个模块的主机规则是否存在，如果没有的话 需要关闭主机属性开关
	for _, moduleID := range moduleIDs {
		filter := map[string]interface{}{
			common.BKAppIDField:    bizID,
			common.BKModuleIDField: moduleID,
		}
		count, err := mongodb.Client().Table(common.BKTableNameHostApplyRule).Find(filter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("update host apply enable status failed, table: %s, filter: %+v,  err: %+v, rid: %s",
				common.BKTableNameBaseModule, filter, err, kit.Rid)
			return err
		}
		if count == 0 {
			modIDs = append(modIDs, moduleID)
		}
	}
	if len(modIDs) == 0 {
		return nil
	}
	enabled := map[string]interface{}{
		common.HostApplyEnabledField: false,
	}

	filter := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: modIDs},
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseModule).Update(kit.Ctx, filter, enabled); nil != err {
		blog.Errorf("update host apply enable status failed, table: %s, filter: %+v,  err: %+v, rid: %s",
			common.BKTableNameBaseModule, filter, err, kit.Rid)
		return err
	}
	return nil
}

// updateServiceTemplateHostApplyStatus after judging the host automatic application rule corresponding to the deleted
// template, whether the template has other corresponding host automatic application rules in the cc_HostApplyRule
// table, if not, the host automatic application state corresponding to this template needs to be turned off.
func (s *coreService) updateTemplateHostApplyStatus(kit *rest.Kit, bizID int64, serviceTemplateIDs []int64) error {

	templateIDs := make([]int64, 0)
	for _, templateID := range serviceTemplateIDs {
		filter := map[string]interface{}{
			common.BKAppIDField:             bizID,
			common.BKServiceTemplateIDField: templateID,
		}
		count, err := mongodb.Client().Table(common.BKTableNameHostApplyRule).Find(filter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("update host apply enable status failed, table: %s, filter: %+v,  err: %+v, rid: %s",
				common.BKTableNameBaseModule, filter, err, kit.Rid)
			return err
		}
		if count == 0 {
			templateIDs = append(templateIDs, templateID)
		}
	}
	if len(templateIDs) == 0 {
		return nil
	}

	enabled := map[string]interface{}{
		common.HostApplyEnabledField: false,
	}

	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldID:    map[string]interface{}{common.BKDBIN: templateIDs},
	}

	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Update(kit.Ctx, filter,
		enabled); nil != err {
		blog.Errorf("update service template host apply enable status failed, filter: %+v,  err: %+v, rid: %s",
			filter, err, kit.Rid)
		return err
	}
	return nil
}

func (s *coreService) updateHostaApplyEnableStatus(kit *rest.Kit, bizID int64,
	option metadata.DeleteHostApplyRuleOption) error {

	if len(option.ModuleIDs) > 0 {
		// update module host apply enabled status
		if err := s.updateModuleHostApplyStatus(kit, bizID, option.ModuleIDs); err != nil {
			return err
		}
	} else {
		if err := s.updateTemplateHostApplyStatus(kit, bizID, option.ServiceTemplateIDs); err != nil {
			return err
		}
	}
	return nil
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

	// Check whether there are other rules in the cc_HostApplyRule table after deleting the rules according to ModuleIDs
	// or templateIDs, if not, then you need to turn off the corresponding host automatic application status.
	if e := s.updateHostaApplyEnableStatus(ctx.Kit, bizID, option); e != nil {
		ctx.RespAutoError(e)
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
