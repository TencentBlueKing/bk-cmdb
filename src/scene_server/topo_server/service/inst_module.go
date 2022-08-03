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
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

const (
	// bKMaxPageSize maximum page size
	bKMaxPageSize = 500
	// bkSetIdSMaxSize maximum number of set's id
	bkSetIdSMaxSize = 200
)

// isSetInitializedByTemplate check if set initialized by template
func (s *Service) isSetInitializedByTemplate(kit *rest.Kit, setID int64) (bool, errors.CCErrorCoder) {
	qc := &metadata.QueryCondition{
		Fields: []string{common.BKSetTemplateIDField},
		Condition: map[string]interface{}{
			common.BKSetIDField: setID,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet,
		qc)
	if err != nil {
		blog.Errorf("failed to search set instance, setID: %d, err: %v, rid: %s", setID, err, kit.Rid)
		return false, errors.NewFromStdError(err, common.CCErrCommHTTPDoRequestFailed)
	}

	if len(result.Info) == 0 {
		blog.Errorf("set instance not exist, setID: %d, rid: %s", setID, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if len(result.Info) > 1 {
		blog.Errorf("got multiple set instance, setID: %d, rid: %s", setID, kit.Rid)

		return false, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	setData := result.Info[0]

	setTemplateID, err := util.GetInt64ByInterface(setData[common.BKSetTemplateIDField])
	if err != nil {
		blog.Errorf("decode set failed, data: %#v, err: %v, rid: %s", setData, err, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	return setTemplateID > 0, nil
}

// CreateModule create a new module
func (s *Service) CreateModule(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id: %s, err: %v, rid: %s", ctx.Request.PathParameter("app_id"),
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the set id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKSetIDField))
		return
	}

	// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
	initializedByTemplate, err := s.isSetInitializedByTemplate(ctx.Kit, setID)
	if err != nil {
		blog.Errorf("set initialized template failed, setID: %d, err: %v, rid: %s", setID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if initializedByTemplate {
		blog.Errorf("forbidden add module to set initialized by template, setID: %d, rid: %s", setID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate))
		return
	}

	module := make(mapstr.MapStr)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		module, err = s.Logics.ModuleOperation().CreateModule(ctx.Kit, bizID, setID, data)
		if err != nil {
			blog.Errorf("create module failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(module)
}

// checkIsBuiltInModule check if object is built-in object
func (s *Service) checkIsBuiltInModule(kit *rest.Kit, moduleIDs ...int64) error {
	// 检查是否时内置集群
	input := &metadata.Condition{
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
			// 当default值不等于0或4时为内置模块，其中4表示主机池中用户自定义创建的模块
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNIN: []int{common.DefaultFlagDefaultValue, common.DefaultResSelfDefinedModuleFlag},
			},
		},
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		input)
	if err != nil {
		blog.Errorf("failed read module instance, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if rsp.Count > 0 {
		return kit.CCError.CCError(common.CCErrorTopoForbiddenDeleteOrUpdateBuiltInSetModule)
	}

	return nil
}

// DeleteModule delete the module
func (s *Service) DeleteModule(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("parse the biz id from path failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if err != nil {
		blog.Errorf("parse the set id from path failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
	initializedByTemplate, err := s.isSetInitializedByTemplate(ctx.Kit, setID)
	if err != nil {
		blog.Errorf("set initialized template failed, setID: %d, err: %v, rid: %s", setID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if initializedByTemplate {
		blog.V(3).Infof("forbidden add module to set initialized by template, setID: %d, rid: %s", setID,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate))
		return
	}

	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter("module_id"), 10, 64)
	if err != nil {
		blog.Errorf("parse the module id from path, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "module id"))
		return
	}

	// 不允许直接删除内置模块
	if err := s.checkIsBuiltInModule(ctx.Kit, moduleID); err != nil {
		blog.Errorf("check is built in module failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.ModuleOperation().DeleteModule(ctx.Kit, bizID, []int64{setID}, []int64{moduleID})
		if err != nil {
			blog.Errorf("delete module failed, delete operation failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// UpdateModule update the module
func (s *Service) UpdateModule(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("parse the biz id from path, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("parse the set id from the path, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("parse the module id from the path, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "module id"))
		return
	}

	// 不允许修改内置模块
	if err := s.checkIsBuiltInModule(ctx.Kit, moduleID); err != nil {
		blog.Errorf("check is builtIn module failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.ModuleOperation().UpdateModule(ctx.Kit, data, bizID, setID, moduleID)
		if err != nil {
			blog.Errorf("update module failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// ListModulesByServiceTemplateID search object by service template ID
func (s *Service) ListModulesByServiceTemplateID(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse bk_biz_id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	serviceTemplateID, e := strconv.ParseInt(ctx.Request.PathParameter(common.BKServiceTemplateIDField), 10, 64)
	if e != nil {
		blog.Errorf("parse service_template_id field failed, err: %v, rid: %s", e, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKServiceTemplateIDField))
		return
	}

	requestBody := struct {
		Page    *metadata.BasePage `field:"page" json:"page" mapstructure:"page"`
		Keyword string             `field:"keyword" json:"keyword" mapstructure:"keyword"`
		Modules []int64            `field:"bk_module_ids" json:"bk_module_ids" mapstructure:"bk_module_ids"`
		Fields  []string           `field:"fields" json:"fields" mapstructure:"fields"`
	}{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// check and set page's limit value
	if requestBody.Page == nil {
		requestBody.Page = &metadata.BasePage{
			Limit: common.BKNoLimit,
		}
	} else {
		if requestBody.Page.Limit == 0 {
			requestBody.Page.Limit = common.BKDefaultLimit
		}
		if requestBody.Page.IsIllegal() {
			blog.Errorf("page is illegal, page: %+v, rid: %s", requestBody.Page, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
			return
		}
	}

	filter := map[string]interface{}{
		common.BKServiceTemplateIDField: serviceTemplateID,
		common.BKAppIDField:             bizID,
	}

	if len(requestBody.Modules) > 0 {
		filter[common.BKModuleIDField] = mapstr.MapStr{
			common.BKDBIN: requestBody.Modules,
		}
	}

	if len(requestBody.Keyword) != 0 {
		filter[common.BKModuleNameField] = map[string]interface{}{
			common.BKDBLIKE: requestBody.Keyword,
		}
	}
	qc := &metadata.QueryCondition{
		Page:      *requestBody.Page,
		Condition: filter,
	}

	if len(requestBody.Fields) > 0 {
		qc.Fields = requestBody.Fields
	}

	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("list modules by service templateID failed, err: %v, cond: %#v, rid: %s", err, qc,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(instanceResult)
}

// SearchModule search module in one set
func (s *Service) SearchModule(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse the biz id from the path failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKSetIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse the set id from the path failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKSetIDField))
		return
	}

	s.searchModule(ctx, bizID, setID)
}

// SearchModuleByCondition search module in one biz
func (s *Service) SearchModuleByCondition(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse the biz id from the path failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	s.searchModule(ctx, bizID, 0)
}

func (s *Service) searchModule(ctx *rest.Contexts, bizID, setID int64) {
	searchCond := new(metadata.SearchModuleCondition)
	if err := ctx.DecodeInto(searchCond); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if searchCond.Condition == nil {
		searchCond.Condition = mapstr.New()
	}
	searchCond.Condition[common.BKAppIDField] = bizID

	// compatible for api /module/search/{owner_id}/{app_id}/{set_id}
	if searchCond.SetID > 0 {
		searchCond.Condition[common.BKSetIDField] = searchCond.SetID
	}
	if setID > 0 {
		searchCond.Condition[common.BKSetIDField] = setID
	}

	queryCond := &metadata.QueryCondition{
		Fields:    searchCond.Fields,
		Condition: searchCond.Condition,
		Page:      searchCond.Page,
	}

	instItems, err := s.Logics.InstOperation().FindInst(ctx.Kit, common.BKInnerObjIDModule, queryCond)
	if err != nil {
		blog.Errorf("search module inst failed, err: %v, cond: %#v, rid: %s", err, queryCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(instItems)
	return
}

// SearchModuleBatch search the modules by module IDs in one biz
func (s *Service) SearchModuleBatch(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.SearchInstBatchOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	moduleIDs := util.IntArrayUnique(option.IDs)
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKModuleIDField: mapstr.MapStr{
			common.BKDBIN: moduleIDs,
		},
	}

	qc := &metadata.QueryCondition{
		Fields: option.Fields,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: cond,
	}
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("search module batch failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(instanceResult.Info)
}

// SearchModuleWithRelation search the modules by set's ids and service template's ids under application
func (s *Service) SearchModuleWithRelation(ctx *rest.Contexts) {
	// parsing input params
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	data := struct {
		BkSetIdS             []int64           `json:"bk_set_ids"`
		BkServiceTemplateIds []int64           `json:"bk_service_template_ids"`
		Fields               []string          `json:"fields"`
		Page                 metadata.BasePage `json:"page"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	// check input params
	if len(data.Fields) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "fields"))
		return
	}
	if len(data.BkSetIdS) > bkSetIdSMaxSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommOverLimit, "the number of bk_set_ids"))
		return
	}
	if data.Page.Limit > bKMaxPageSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommOverLimit, "page"))
		return
	}

	// set query condition
	cond := mapstr.MapStr{common.BKAppIDField: bizID}
	if len(data.BkSetIdS) != 0 {
		bkSetIdS := util.IntArrayUnique(data.BkSetIdS)
		cond[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: bkSetIdS}
	}
	if len(data.BkServiceTemplateIds) != 0 {
		bkServiceTemplateIds := util.IntArrayUnique(data.BkServiceTemplateIds)
		cond[common.BKServiceTemplateIDField] = mapstr.MapStr{common.BKDBIN: bkServiceTemplateIds}
	}
	qc := &metadata.QueryCondition{
		Fields:    data.Fields,
		Page:      data.Page,
		Condition: cond,
	}

	// query and check result
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("search module with relation failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithCount(int64(instanceResult.Count), instanceResult.Info)

	return
}

// SearchRuleRelatedTopoNodes search rule related topo nodes
func (s *Service) SearchRuleRelatedTopoNodes(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse bk_biz_id from the path failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	requestBody := metadata.SearchRuleRelatedModulesOption{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if requestBody.QueryFilter == nil {
		blog.Errorf("search query_filter should not be empty, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter"))
		return
	}
	if key, err := requestBody.QueryFilter.Validate(&querybuilder.RuleOption{NeedSameSliceElementType: true}); err != nil {
		blog.Errorf("search query_filter.%s validate failed, err: %v, rid: %s", key, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter."+key))
		return
	}

	modules, err := s.Engine.CoreAPI.CoreService().HostApplyRule().SearchRuleRelatedModules(ctx.Kit.Ctx,
		ctx.Kit.Header, bizID, requestBody)
	if err != nil {
		blog.Errorf("search rule related modules failed, err: %v, rid: %s", err, ctx.Kit.Rid)
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

	ctx.RespEntity(matchNodes)
}

// UpdateModuleHostApplyEnableStatus update object host if apply's status is enabled
func (s *Service) UpdateModuleHostApplyEnableStatus(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse bk_biz_id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}
	request := metadata.UpdateHostApplyEnableStatusOption{}
	if err := ctx.DecodeInto(&request); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(request.IDs) == 0 {
		blog.Errorf("module ids must be set, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "bk_module_ids"))
		return
	}

	if err := request.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	updateOption := &metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKAppIDField:    bizID,
			common.BKModuleIDField: mapstr.MapStr{common.BKDBIN: request.IDs},
		},
		Data: map[string]interface{}{
			common.HostApplyEnabledField: request.Enable,
		},
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		_, err = s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDModule, updateOption)
		if err != nil {
			blog.Errorf("search rule related modules failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		// If this request is to enable the host to automatically apply, then the cleanup rules are not involved, and
		// return directly here。
		if request.Enable {
			return nil
		}
		if request.ClearRules {
			listRuleOption := metadata.ListHostApplyRuleOption{
				ModuleIDs: request.IDs,
				Page: metadata.BasePage{
					Limit: common.BKNoLimit,
				},
			}
			listRuleResult, ccErr := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx,
				ctx.Kit.Header, bizID, listRuleOption)
			if ccErr != nil {
				blog.Errorf("get list host apply rule failed, bizID: %d,listRuleOption: %#v, rid: %s", bizID,
					listRuleOption, ctx.Kit.Rid)
				return ccErr
			}
			ruleIDs := make([]int64, 0)
			for _, item := range listRuleResult.Info {
				ruleIDs = append(ruleIDs, item.ID)
			}
			if len(ruleIDs) > 0 {
				deleteRuleOption := metadata.DeleteHostApplyRuleOption{
					RuleIDs:   ruleIDs,
					ModuleIDs: request.IDs,
				}
				if ccErr := s.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx,
					ctx.Kit.Header, bizID, deleteRuleOption); ccErr != nil {
					blog.Errorf("delete list host apply rule failed, bizID: %d, listRuleOption: %#v, rid: %s",
						bizID, listRuleOption, ctx.Kit.Rid)
					return ccErr
				}
			}
		}
		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// GetInternalModuleWithStatistics get internal object by statistics
func (s *Service) GetInternalModuleWithStatistics(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, innerAppTopo, err := s.Logics.ModuleOperation().GetInternalModule(ctx.Kit, bizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if innerAppTopo == nil {
		blog.Errorf("get internal module with statistics failed, type: %#v, rid: %s", innerAppTopo, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	moduleIDArr := make([]int64, 0)
	for _, item := range innerAppTopo.Module {
		moduleIDArr = append(moduleIDArr, item.ModuleID)
	}

	// count host apply rules
	listApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDArr,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostApplyRules, err := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx,
		ctx.Kit.Header, bizID, listApplyRuleOption)
	if err != nil {
		blog.Errorf("get list host apply rule failed, bizID: %d, err: %v, rid: %s", bizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	moduleRuleCount := make(map[int64]int64)
	for _, item := range hostApplyRules.Info {
		if _, exist := moduleRuleCount[item.ModuleID]; !exist {
			moduleRuleCount[item.ModuleID] = 0
		}
		moduleRuleCount[item.ModuleID] += 1
	}

	set := mapstr.NewFromStruct(innerAppTopo, "field")
	modules := make([]mapstr.MapStr, 0)
	for _, module := range innerAppTopo.Module {
		moduleItem := mapstr.NewFromStruct(module, "field")
		moduleItem["host_apply_rule_count"] = 0
		if ruleCount, ok := moduleRuleCount[module.ModuleID]; ok {
			moduleItem["host_apply_rule_count"] = ruleCount
		}
		modules = append(modules, moduleItem)
	}
	set["module"] = modules
	ctx.RespEntity(set)
}

// GetInternalModule get internal module
func (s *Service) GetInternalModule(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error()))
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	_, result, err := s.Logics.ModuleOperation().GetInternalModule(ctx.Kit, bizID)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// UpdateGlobalSetOrModuleConfig update platform_setting，注意： 此接口只给前端的管理员使用不能上ESB
func (s *Service) UpdateGlobalSetOrModuleConfig(ctx *rest.Contexts) {

	option := new(metadata.ConfigUpdateSettingOption)
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		blog.Errorf("update global config fail, param is invalid, input: %v, error: %v, rid: %s", option, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoAppSearchFailed))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.BusinessOperation().UpdateBusinessIdleSetOrModule(ctx.Kit, option)
		if err != nil {
			blog.Errorf("update business set or module fail, option: %v, err: %v, rid: %s", option, err,
				ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteUserModulesSettingConfig delete user config module，注意： 此接口只给前端的管理员使用不能上ESB
func (s *Service) DeleteUserModulesSettingConfig(ctx *rest.Contexts) {

	option := new(metadata.BuiltInModuleDeleteOption)
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		blog.Errorf("option is illegal option: %+v, rid: %s", option, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoAppSearchFailed, "module key and name must be set"))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.BusinessOperation().DeleteBusinessGlobalUserModule(ctx.Kit, option)
		if err != nil {
			blog.Errorf("create business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		return nil

	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}
