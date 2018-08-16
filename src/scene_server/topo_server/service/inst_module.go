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
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateModule create a new module
func (s *topoService) CreateModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	return s.core.ModuleOperation().CreateModule(params, obj, bizID, setID, data)

}

// DeleteModule delete the module
func (s *topoService) DeleteModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	moduleID, err := strconv.ParseInt(pathParams("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the module id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "module id")
	}

	return nil, s.core.ModuleOperation().DeleteModule(params, obj, bizID, []int64{setID}, []int64{moduleID})

}

// UpdateModule update the module
func (s *topoService) UpdateModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	moduleID, err := strconv.ParseInt(pathParams("module_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the module id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "module id")
	}

	return nil, s.core.ModuleOperation().UpdateModule(params, data, obj, bizID, setID, moduleID)
}

// SearchModule search the modules
func (s *topoService) SearchModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the module, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-module]failed to parse the set id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	paramsCond := &gparams.SearchParams{}
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

	cnt, instItems, err := s.core.ModuleOperation().FindModule(params, obj, queryCond)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil
}
