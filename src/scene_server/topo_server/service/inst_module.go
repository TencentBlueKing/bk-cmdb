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
	"configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
)

const (
	// bKMaxPageSize maximum page size
	bKMaxPageSize = 500
	// bkSetIdSMaxSize maximum number of set's id
	bkSetIdSMaxSize = 200
)

func (s *Service) IsSetInitializedByTemplate(kit *rest.Kit, setID int64) (bool, errors.CCErrorCoder) {
	qc := &metadata.QueryCondition{
		Fields: []string{common.BKSetTemplateIDField},
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
	setTemplateID, err := util.GetInt64ByInterface(setData[common.BKSetTemplateIDField])
	if err != nil {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, decode set failed, data: %s, err: %s, rid: %s", setData)
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

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule)
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

	var module inst.Inst
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		module, err = s.Core.ModuleOperation().CreateModule(ctx.Kit, obj, bizID, setID, data)
		if err != nil {
			blog.Errorf("[api-module] create module failed, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.ModuleOperation().DeleteModule(ctx.Kit, bizID, []int64{setID}, []int64{moduleID})
		if err != nil {
			blog.Errorf("delete module failed, delete operation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
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

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule)
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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.ModuleOperation().UpdateModule(ctx.Kit, data, obj, bizID, setID, moduleID)
		if err != nil {
			blog.Errorf("update module failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
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
		Modules []int64            `field:"bk_module_ids" json:"bk_module_ids" mapstructure:"bk_module_ids"`
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
			blog.Errorf("ListModulesByServiceTemplateID failed, Page is IsIllegal, rid:%s, page:%+v", ctx.Kit.Rid, requestBody.Page)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
			return
		}
	}

	filter := map[string]interface{}{
		common.BKServiceTemplateIDField: serviceTemplateID,
		common.BKAppIDField:             bizID,
	}

	if requestBody.Modules != nil {
		filter[common.BKModuleIDField] = mapstr.MapStr{common.BKDBIN: requestBody.Modules}
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

// SearchModuleInOneSet search module in one set
func (s *Service) SearchModule(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKSetIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKSetIDField))
		return
	}

	s.searchModule(ctx, bizID, setID)
}

// SearchModuleByCondition search module in one biz
func (s *Service) SearchModuleByCondition(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	s.searchModule(ctx, bizID, 0)
}

func (s *Service) searchModule(ctx *rest.Contexts, bizID, setID int64) {
	data := struct {
		parser.SearchParams `json:",inline"`
		// compatible for api /module/search/{owner_id}/{app_id}/{set_id}
		SetID int64 `json:"bk_set_id"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	paramsCond := data.SearchParams
	if paramsCond.Condition == nil {
		paramsCond.Condition = mapstr.New()
	}

	paramsCond.Condition[common.BKAppIDField] = bizID

	// compatible for api /module/search/{owner_id}/{app_id}/{set_id}
	if data.SetID > 0 {
		paramsCond.Condition[common.BKSetIDField] = data.SetID
	}
	if setID > 0 {
		paramsCond.Condition[common.BKSetIDField] = setID
	}

	queryCond := &metadata.QueryInput{}
	queryCond.Condition = paramsCond.Condition
	queryCond.Fields = strings.Join(paramsCond.Fields, ",")
	page := metadata.ParsePage(paramsCond.Page)
	queryCond.Limit = page.Limit
	queryCond.Sort = page.Sort
	queryCond.Start = page.Start

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

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
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("SearchModuleBatch failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !instanceResult.Result {
		blog.ErrorJSON("SearchModuleBatch failed, ReadInstance failed, filter: %s, response: %s, rid: %s", qc, instanceResult, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(instanceResult.Code, instanceResult.ErrMsg))
		return
	}
	ctx.RespEntity(instanceResult.Data.Info)
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
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("SearchModuleWithRelation failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !instanceResult.Result {
		blog.ErrorJSON("SearchModuleWithRelation failed, ReadInstance failed, filter: %s, response: %s, rid: %s", qc, instanceResult, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(instanceResult.Code, instanceResult.ErrMsg))
		return
	}

	ctx.RespEntityWithCount(int64(instanceResult.Data.Count), instanceResult.Data.Info)
	return
}

func (s *Service) SearchRuleRelatedTopoNodes(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("SearchRuleRelatedModules failed, parse bk_biz_id failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	requestBody := metadata.SearchRuleRelatedModulesOption{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
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

	var result *metadata.UpdatedOptionResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		result, err = s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, updateOption)
		if err != nil {
			blog.Errorf("SearchRuleRelatedModules failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if ccErr := result.CCError(); ccErr != nil {
			blog.ErrorJSON("SearchRuleRelatedModules failed, update module instance failed, updateOption: %s, response: %s, rid: %s", updateOption, result, ctx.Kit.Rid)
			return ccErr
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
				return ccErr
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
	ctx.RespEntity(result.Data)
}
