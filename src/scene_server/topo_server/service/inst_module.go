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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

func (s *Service) IsSetInitializedByTemplate(kit *rest.Kit, setID int64) (bool, errors.CCErrorCoder) {
	qc := &metadata.QueryCondition{
		Fields: []string{common.BKSetTemplateIDField, common.BKSetIDField},
		Condition: map[string]interface{}{
			common.BKSetIDField: setID,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, qc)
	if err != nil {
		blog.Errorf("IsSetInitializedByTemplate failed, failed to search set instance, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
		return false, errors.NewFromStdError(err, common.CCErrCommHTTPDoRequestFailed)
	}
	if result.Code != 0 {
		return false, errors.NewCCError(result.Code, result.ErrMsg)
	}
	if len(result.Data.Info) == 0 {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, set:%d not found, rid: %s", setID, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(result.Data.Info) > 1 {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, set:%d got multiple, rid: %s", setID, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	setData := result.Data.Info[0]
	set := metadata.SetInst{}
	if err := mapstruct.Decode2Struct(setData, &set); err != nil {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, decode set failed, data: %s, err: %s, rid: %s", setData)
		return false, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	return set.SetTemplateID > 0, nil
}

// CreateModule create a new module
func (s *Service) CreateModule(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule, dataWithMetadata.Metadata)
	if nil != err {
		blog.Errorf("create module failed, failed to search set model, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module] create module failed, failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module] create module failed, failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKSetIDField))
		return
	}

	// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
	initializedByTemplate, err := s.IsSetInitializedByTemplate(ctx.Kit, setID)
	if err != nil {
		blog.Errorf("CreateModule failed, IsSetInitializedByTemplate failed, setID: %d, err: %s, rid: %s", setID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if initializedByTemplate == true {
		blog.V(3).Infof("CreateModule failed, forbidden add module to set initialized by template, setID: %d, rid: %s", setID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate))
		return
	}

	module, err := s.Core.ModuleOperation().CreateModule(ctx.Kit, obj, bizID, setID, dataWithMetadata.Data)
	if err != nil {
		blog.Errorf("[api-module] create module failed, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(module)
}

func (s *Service) CheckIsBuiltInModule(kit *rest.Kit, moduleIDs ...int64) errors.CCErrorCoder {
	// 检查是否时内置集群
	qc := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: 0,
		},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: common.DefaultFlagDefaultValue,
			},
		},
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, qc)
	if nil != err {
		blog.Errorf("[operation-module] failed read module instance, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if rsp.Result == false || rsp.Code != 0 {
		blog.ErrorJSON("[operation-set] failed read module instance, option: %s, response: %s, rid: %s", qc, rsp, kit.Rid)
		return errors.New(rsp.Code, rsp.ErrMsg)
	}
	if rsp.Data.Count > 0 {
		return kit.CCError.CCError(common.CCErrorTopoForbiddenDeleteBuiltInSetModule)
	}
	return nil
}

// DeleteModule delete the module
func (s *Service) DeleteModule(ctx *rest.Contexts) {
	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule, md.Metadata)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
	initializedByTemplate, err := s.IsSetInitializedByTemplate(ctx.Kit, setID)
	if err != nil {
		blog.Errorf("DeleteModule failed, IsSetInitializedByTemplate failed, setID: %d, err: %s, rid: %s", setID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if initializedByTemplate == true {
		blog.V(3).Infof("DeleteModule failed, forbidden add module to set initialized by template, setID: %d, rid: %s", setID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate))
		return
	}

	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the module id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "module id"))
		return
	}

	// 不允许直接删除内置模块
	if err := s.CheckIsBuiltInModule(ctx.Kit, moduleID); err != nil {
		blog.Errorf("[api-module]DeleteModule failed, CheckIsBuiltInModule failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	err = s.Core.ModuleOperation().DeleteModule(ctx.Kit, obj, bizID, []int64{setID}, []int64{moduleID})
	if err != nil {
		blog.Errorf("delete module failed, delete operation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// auth: deregister module to iam
	if err := s.AuthManager.DeregisterModuleByID(ctx.Kit.Ctx, ctx.Kit.Header, moduleID); err != nil {
		blog.Errorf("delete module failed, deregister module failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
        ctx.RespAutoError(err)
        return
	}

	ctx.RespEntity(nil)
}

// UpdateModule update the module
func (s *Service) UpdateModule(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule, dataWithMetadata.Metadata)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	/*
		// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
		initializedByTemplate, err := s.IsSetInitializedByTemplate(ctx.Kit, setID)
		if err != nil {
			blog.Errorf("UpdateModule failed, IsSetInitializedByTemplate failed, setID: %d, err: %s, rid: %s", setID, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		if initializedByTemplate == true {
			blog.V(3).Infof("UpdateModule failed, forbidden add module to set initialized by template, setID: %d, rid: %s", setID, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate))
			return
		}
	*/

	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the module id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "module id"))
		return
	}

	err = s.Core.ModuleOperation().UpdateModule(ctx.Kit, dataWithMetadata.Data, obj, bizID, setID, moduleID)
	if err != nil {
		blog.Errorf("update module failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) ListModulesByServiceTemplateID(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("ListModulesByServiceTemplateID failed, parse bk_biz_id failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	serviceTemplateID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKServiceTemplateIDField), 10, 64)
	if nil != err {
		blog.Errorf("ListModulesByServiceTemplateID failed, parse service_template_id field failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKServiceTemplateIDField))
		return
	}

	requestBody := struct {
		Page    *metadata.BasePage `field:"page" json:"page" mapstructure:"page"`
		Keyword string             `field:"keyword" json:"keyword" mapstructure:"keyword"`
	}{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	start := 0
	limit := common.BKDefaultLimit
	if requestBody.Page != nil {
		limit = requestBody.Page.Limit
		start = requestBody.Page.Start
	}
	filter := map[string]interface{}{
		common.BKServiceTemplateIDField: serviceTemplateID,
		common.BKAppIDField:             bizID,
	}
	if len(requestBody.Keyword) != 0 {
		filter[common.BKModuleNameField] = map[string]interface{}{
			common.BKDBLIKE: requestBody.Keyword,
		}
	}
	qc := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Start: start,
			Limit: limit,
			Sort:  requestBody.Page.Sort,
		},
		Condition: filter,
	}
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("ListModulesByServiceTemplateID failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if instanceResult.Code != 0 {
		blog.ErrorJSON("ListModulesByServiceTemplateID failed, ReadInstance failed, filter: %s, response: %s, rid: %s", qc, instanceResult, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(instanceResult.Code, instanceResult.ErrMsg))
		return
	}
	ctx.RespEntity(instanceResult.Data)
}

// SearchModule search the modules
func (s *Service) SearchModule(ctx *rest.Contexts) {
	data := struct {
		parser.SearchParams `json:",inline"`
		Metadata            *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	paramsCond := data.SearchParams
	if paramsCond.Condition == nil {
		paramsCond.Condition = mapstr.New()
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule, data.Metadata)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	paramsCond.Condition[common.BKAppIDField] = bizID
	paramsCond.Condition[common.BKSetIDField] = setID

	queryCond := &metadata.QueryInput{}
	queryCond.Condition = paramsCond.Condition
	queryCond.Fields = strings.Join(paramsCond.Fields, ",")
	page := metadata.ParsePage(paramsCond.Page)
	queryCond.Limit = page.Limit
	queryCond.Sort = page.Sort
	queryCond.Start = page.Start

	cnt, instItems, err := s.Core.ModuleOperation().FindModule(ctx.Kit, obj, queryCond)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	ctx.RespEntity(result)
	return
}

func (s *Service) SearchRuleRelatedTopoNodes(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("SearchRuleRelatedModules failed, parse bk_biz_id failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	requestBody := metadata.SearchRuleRelatedModulesOption{}
	if err := mapstruct.Decode2Struct(data, &requestBody); err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, parse request body failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if requestBody.QueryFilter == nil {
		blog.V(3).Info("SearchRuleRelatedModules failed, search query_filter should'nt be empty, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter"))
		return
	}
	if key, err := requestBody.QueryFilter.Validate(); err != nil {
		blog.V(3).Info("SearchRuleRelatedModules failed, search query_filter.%s validate failed, err: %+v, rid: %s", key, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter."+key))
		return
	}

	modules, err := s.Engine.CoreAPI.CoreService().HostApplyRule().SearchRuleRelatedModules(ctx.Kit.Ctx, ctx.Kit.Header, bizID, requestBody)
	if err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	topoRoot, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx, ctx.Kit.Header, bizID, false)
	if err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, SearchMainlineInstanceTopo failed, bizID: %d, err: %s, rid: %s", bizID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	matchNodes := make([]metadata.TopoNode, 0)
	for _, module := range modules {
		matchNodes = append(matchNodes, metadata.TopoNode{
			ObjectID:   common.BKInnerObjIDModule,
			InstanceID: module.ModuleID,
		})
	}
	topoRoot.DeepFirstTraversal(func(node *metadata.TopoInstanceNode) {
		matched := requestBody.QueryFilter.Match(func(r querybuilder.AtomRule) bool {
			if r.Field != metadata.TopoNodeKeyword {
				return false
			}
			valueStr, ok := r.Value.(string)
			if ok == false {
				return false
			}
			// case-insensitive contains
			if r.Operator == querybuilder.OperatorContains {
				if util.CaseInsensitiveContains(node.InstanceName, valueStr) {
					return true
				}
			}
			return false
		})
		if matched {
			matchNodes = append(matchNodes, metadata.TopoNode{
				ObjectID:   node.ObjectID,
				InstanceID: node.InstanceID,
			})
		}
	})

	// unique result
	finalNodes := make([]metadata.TopoNode, 0)
	existMap := make(map[string]bool)
	for _, item := range matchNodes {
		if _, exist := existMap[item.Key()]; exist == true {
			continue
		}
		existMap[item.Key()] = true
		finalNodes = append(finalNodes, item)
	}
	ctx.RespEntity(finalNodes)
}

func (s *Service) UpdateModuleHostApplyEnableStatus(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("UpdateModuleHostApplyEnableStatus failed, parse bk_biz_id failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKModuleIDField), 10, 64)
	if nil != err {
		blog.Errorf("UpdateModuleHostApplyEnableStatus failed, parse bk_module_id failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKModuleIDField))
		return
	}

	requestBody := metadata.UpdateModuleHostApplyEnableStatusOption{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}
	updateOption := &metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKAppIDField:    bizID,
			common.BKModuleIDField: moduleID,
		},
		Data: map[string]interface{}{
			common.HostApplyEnabledField: requestBody.Enable,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, updateOption)
	if err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if ccErr := result.CCError(); ccErr != nil {
		blog.ErrorJSON("SearchRuleRelatedModules failed, update module instance failed, updateOption: %s, response: %s, rid: %s", updateOption, result, ctx.Kit.Rid)
		ctx.RespAutoError(ccErr)
		return
	}
	if requestBody.ClearRules {
		listRuleOption := metadata.ListHostApplyRuleOption{
			ModuleIDs: []int64{moduleID},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		listRuleResult, ccErr := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, listRuleOption)
		if ccErr != nil {
			blog.ErrorJSON("SearchRuleRelatedModules failed, ListHostApplyRule failed, bizID: %s, listRuleOption: %s, rid: %s", bizID, listRuleOption, ctx.Kit.Rid)
			ctx.RespAutoError(ccErr)
			return
		}
		ruleIDs := make([]int64, 0)
		for _, item := range listRuleResult.Info {
			ruleIDs = append(ruleIDs, item.ID)
		}
		if len(ruleIDs) > 0 {
			deleteRuleOption := metadata.DeleteHostApplyRuleOption{
				RuleIDs: ruleIDs,
			}
			if ccErr := s.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, deleteRuleOption); ccErr != nil {
				blog.ErrorJSON("SearchRuleRelatedModules failed, ListHostApplyRule failed, bizID: %s, listRuleOption: %s, rid: %s", bizID, listRuleOption, ctx.Kit.Rid)
				ctx.RespAutoError(ccErr)
				return
			}
		}
	}
	ctx.RespEntity(result.Data)
	return
}
