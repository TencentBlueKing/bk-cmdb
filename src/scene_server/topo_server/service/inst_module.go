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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (s *Service) IsSetInitializedByTemplate(params types.ContextParams, setID int64) (bool, errors.CCErrorCoder) {
	qc := &metadata.QueryCondition{
		Fields: []string{common.BKSetTemplateIDField, common.BKSetIDField},
		Condition: map[string]interface{}{
			common.BKSetIDField: setID,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDSet, qc)
	if err != nil {
		blog.Errorf("IsSetInitializedByTemplate failed, failed to search set instance, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return false, errors.NewFromStdError(err, common.CCErrCommHTTPDoRequestFailed)
	}
	if result.Code != 0 {
		return false, errors.NewCCError(result.Code, result.ErrMsg)
	}
	if len(result.Data.Info) == 0 {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, set:%d not found, rid: %s", setID, params.ReqID)
		return false, params.Err.CCError(common.CCErrCommNotFound)
	}
	if len(result.Data.Info) > 1 {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, set:%d got multiple, rid: %s", setID, params.ReqID)
		return false, params.Err.CCError(common.CCErrCommGetMultipleObject)
	}
	setData := result.Data.Info[0]
	set := metadata.SetInst{}
	if err := mapstruct.Decode2Struct(setData, &set); err != nil {
		blog.ErrorJSON("IsSetInitializedByTemplate failed, decode set failed, data: %s, err: %s, rid: %s", setData)
		return false, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	return set.SetTemplateID > 0, nil
}

// CreateModule create a new module
func (s *Service) CreateModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("create module failed, failed to search set model, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module] create module failed, failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module] create module failed, failed to parse the set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, common.BKSetIDField)
	}

	// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
	initializedByTemplate, err := s.IsSetInitializedByTemplate(params, setID)
	if err != nil {
		blog.Errorf("CreateModule failed, IsSetInitializedByTemplate failed, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return nil, err
	}
	if initializedByTemplate == true {
		blog.V(3).Infof("CreateModule failed, forbidden add module to set initialized by template, setID: %d, rid: %s", setID, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate)
	}

	module, err := s.Core.ModuleOperation().CreateModule(params, obj, bizID, setID, data)
	if err != nil {
		blog.Errorf("[api-module] create module failed, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	return module, nil
}

func (s *Service) CheckIsBuiltInModule(params types.ContextParams, moduleIDs ...int64) errors.CCErrorCoder {
	// 检查是否时内置集群
	qc := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
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
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDModule, qc)
	if nil != err {
		blog.Errorf("[operation-module] failed read module instance, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if rsp.Result == false || rsp.Code != 0 {
		blog.ErrorJSON("[operation-set] failed read module instance, option: %s, response: %s, rid: %s", qc, rsp, params.ReqID)
		return errors.New(rsp.Code, rsp.ErrMsg)
	}
	if rsp.Data.Count > 0 {
		return params.Err.CCError(common.CCErrorTopoForbiddenDeleteBuiltInSetModule)
	}
	return nil
}

// DeleteModule delete the module
func (s *Service) DeleteModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
	initializedByTemplate, err := s.IsSetInitializedByTemplate(params, setID)
	if err != nil {
		blog.Errorf("DeleteModule failed, IsSetInitializedByTemplate failed, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return nil, err
	}
	if initializedByTemplate == true {
		blog.V(3).Infof("DeleteModule failed, forbidden add module to set initialized by template, setID: %d, rid: %s", setID, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate)
	}

	moduleID, err := strconv.ParseInt(pathParams("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the module id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "module id")
	}

	// auth: deregister module to iam
	if err := s.AuthManager.DeregisterModuleByID(params.Context, params.Header, moduleID); err != nil {
		blog.Errorf("delete module failed, deregister module failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	err = s.Core.ModuleOperation().DeleteModule(params, obj, bizID, []int64{setID}, []int64{moduleID})
	if err != nil {
		blog.Errorf("delete module failed, delete operation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return nil, nil
}

// UpdateModule update the module
func (s *Service) UpdateModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	/*
		// 通过集群模板创建的模板禁止直接操作(只能通过集群模板同步)
		initializedByTemplate, err := s.IsSetInitializedByTemplate(params, setID)
		if err != nil {
			blog.Errorf("UpdateModule failed, IsSetInitializedByTemplate failed, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
			return nil, err
		}
		if initializedByTemplate == true {
			blog.V(3).Infof("UpdateModule failed, forbidden add module to set initialized by template, setID: %d, rid: %s", setID, params.ReqID)
			return nil, params.Err.Error(common.CCErrorTopoForbiddenOperateModuleOnSetInitializedByTemplate)
		}
	*/

	moduleID, err := strconv.ParseInt(pathParams("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the module id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "module id")
	}

	err = s.Core.ModuleOperation().UpdateModule(params, data, obj, bizID, setID, moduleID)
	if err != nil {
		blog.Errorf("update module failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return nil, nil
}

func (s *Service) ListModulesByServiceTemplateID(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("ListModulesByServiceTemplateID failed, parse bk_biz_id failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
	}

	serviceTemplateID, err := strconv.ParseInt(pathParams(common.BKServiceTemplateIDField), 10, 64)
	if nil != err {
		blog.Errorf("ListModulesByServiceTemplateID failed, parse service_template_id field failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, common.BKServiceTemplateIDField)
	}

	requestBody := struct {
		Page    *metadata.BasePage `field:"page" json:"page" mapstructure:"page"`
		Keyword string             `field:"keyword" json:"keyword" mapstructure:"keyword"`
	}{}
	if err := mapstruct.Decode2Struct(data, &requestBody); err != nil {
		blog.Errorf("ListModulesByServiceTemplateID failed, parse request body failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	start := int64(0)
	limit := int64(common.BKDefaultLimit)
	sortArr := make([]metadata.SearchSort, 0)
	if requestBody.Page != nil {
		limit = int64(requestBody.Page.Limit)
		start = int64(requestBody.Page.Start)
		sortArr = requestBody.Page.ToSearchSort()
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
		Limit: metadata.SearchLimit{
			Offset: start,
			Limit:  limit,
		},
		SortArr:   sortArr,
		Condition: filter,
	}
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("ListModulesByServiceTemplateID failed, http request failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if instanceResult.Code != 0 {
		blog.ErrorJSON("ListModulesByServiceTemplateID failed, ReadInstance failed, filter: %s, response: %s, rid: %s", qc, instanceResult, params.ReqID)
		return nil, errors.New(instanceResult.Code, instanceResult.ErrMsg)
	}
	return instanceResult.Data, nil
}

// SearchModule search the modules
func (s *Service) SearchModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	paramsCond := &parser.SearchParams{
		Condition: mapstr.New(),
	}
	if err = data.MarshalJSONInto(paramsCond); nil != err {
		return nil, err
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

	cnt, instItems, err := s.Core.ModuleOperation().FindModule(params, obj, queryCond)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil
}
