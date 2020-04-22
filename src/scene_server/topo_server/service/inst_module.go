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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateModule create a new module
func (s *Service) CreateModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("create module failed, failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
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

	module, err := s.Core.ModuleOperation().CreateModule(params, obj, bizID, setID, data)
	if err != nil {
		blog.Errorf("[api-module] create module failed, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	moduleID, err := module.GetInstID()
	if err != nil {
		blog.Errorf("create module failed, unexpected error, create module success, but get id failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	// auth: register module to iam
	if err := s.AuthManager.RegisterModuleByID(params.Context, params.Header, moduleID); err != nil {
		blog.Errorf("create module success, but register module failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return module, nil
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

	// auth: update registered module to iam
	if err := s.AuthManager.UpdateRegisteredModuleByID(params.Context, params.Header, moduleID); err != nil {
		blog.Errorf("update module success, but update registered module failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return nil, nil
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
