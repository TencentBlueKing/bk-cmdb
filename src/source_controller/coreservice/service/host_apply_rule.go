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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateHostApplyRule TODO
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

// UpdateHostApplyRule TODO
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
func (s *coreService) updateModuleHostApplyStatus(kit *rest.Kit, bizID int64, moduleIDs []int64, enabled bool) error {

	filter := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: moduleIDs},
	}
	fields := []string{common.BKModuleIDField}

	rules := make([]metadata.HostApplyRule, 0)
	err := mongodb.Client().Table(common.BKTableNameHostApplyRule).Find(filter).Fields(fields...).All(kit.Ctx, &rules)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return err
	}
	dbModIDs := make([]int64, 0)
	for _, rule := range rules {
		dbModIDs = append(dbModIDs, rule.ModuleID)
	}

	modIDs := util.IntArrDeleteElements(moduleIDs, dbModIDs)

	enabledField := map[string]interface{}{
		common.HostApplyEnabledField: enabled,
	}

	option := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: modIDs},
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseModule).Update(kit.Ctx, option, enabledField); nil != err {
		blog.Errorf("update host apply enable status failed, table: %s, filter: %+v,  err: %+v, rid: %s",
			common.BKTableNameBaseModule, filter, err, kit.Rid)
		return err
	}
	return nil
}

// updateTemplateHostApplyStatus after judging the host automatic application rule corresponding to the deleted
// template, whether the template has other corresponding host automatic application rules in the cc_HostApplyRule
// table, if not, the host automatic application state corresponding to this template needs to be turned off.
func (s *coreService) updateTemplateHostApplyStatus(kit *rest.Kit, bizID int64, serviceTemplateIDs []int64,
	enabled bool) error {

	filter := map[string]interface{}{
		common.BKAppIDField:             bizID,
		common.BKServiceTemplateIDField: map[string]interface{}{common.BKDBIN: serviceTemplateIDs},
	}
	fields := []string{common.BKServiceTemplateIDField}
	rules := make([]metadata.HostApplyRule, 0)
	err := mongodb.Client().Table(common.BKTableNameHostApplyRule).Find(filter).Fields(fields...).All(kit.Ctx, &rules)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return err
	}

	dbServiceTemplateIDs := make([]int64, 0)
	for _, rule := range rules {
		dbServiceTemplateIDs = append(dbServiceTemplateIDs, rule.ServiceTemplateID)
	}

	templateIDs := util.IntArrDeleteElements(serviceTemplateIDs, dbServiceTemplateIDs)

	enabledField := map[string]interface{}{
		common.HostApplyEnabledField: enabled,
	}

	updateFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldID:    map[string]interface{}{common.BKDBIN: templateIDs},
	}

	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Update(kit.Ctx, updateFilter,
		enabledField); nil != err {
		blog.Errorf("update service template host apply enable status failed, filter: %+v,  err: %+v, rid: %s",
			filter, err, kit.Rid)
		return err
	}
	return nil
}

func (s *coreService) updateHostApplyEnableStatus(kit *rest.Kit, bizID int64,
	option metadata.DeleteHostApplyRuleOption) error {

	if len(option.ModuleIDs) > 0 {
		// update module host apply enabled status
		if err := s.updateModuleHostApplyStatus(kit, bizID, option.ModuleIDs, false); err != nil {
			return err
		}
	} else {
		if err := s.updateTemplateHostApplyStatus(kit, bizID, option.ServiceTemplateIDs, false); err != nil {
			return err
		}
	}
	return nil
}

// DeleteHostApplyRule TODO
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

	if err := s.core.HostApplyRuleOperation().DeleteHostApplyRule(ctx.Kit, bizID, option); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, bizID: %d, ruleID: %d, err: %+v, rid: %s", bizID, option.RuleIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// Check whether there are other rules in the cc_HostApplyRule table after deleting the rules according to ModuleIDs
	// or templateIDs, if not, then you need to turn off the corresponding host automatic application status.
	if e := s.updateHostApplyEnableStatus(ctx.Kit, bizID, option); e != nil {
		ctx.RespAutoError(e)
		return
	}
	ctx.RespEntity(nil)
}

// GetHostApplyRule TODO
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

// ListHostApplyRule TODO
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

// GenerateApplyPlan TODO
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

// SearchRuleRelatedModules TODO
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

// BatchUpdateHostApplyRule TODO
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

// UpdateHostByHostApplyRule TODO
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

	relationFilter := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: option.HostIDs}}
	relations := make([]metadata.ModuleHost, 0)
	err = mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(relationFilter).All(ctx.Kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("find %s failed, filter: %s, err: %v, rid: %s", common.BKTableNameModuleHostConfig, relationFilter,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.HostApplyRuleOperation().RunHostApplyOnHosts(ctx.Kit, bizID, relations)
	if err != nil {
		blog.Errorf("UpdateHostByHostApplyRule failed, RunHostApplyOnHosts failed, option: %+v, err: %+v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SearchRuleRelatedServiceTemplates search rule related service templates
func (s *coreService) SearchRuleRelatedServiceTemplates(ctx *rest.Contexts) {
	option := metadata.RuleRelatedServiceTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	serviceTemplates, err := s.core.HostApplyRuleOperation().SearchRuleRelatedServiceTemplates(ctx.Kit, option)
	if err != nil {
		blog.Errorf("search templates failed, option: %v, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(serviceTemplates)
}
