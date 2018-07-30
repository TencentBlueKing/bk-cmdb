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
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

// SearchAllApp search all business
func (s *topoService) SearchAllApp(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond, err := data.MapStr("condition")
	if nil != err {
		blog.Errorf("[api-compatiblev2] not set the condition in the search conditons")
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set the search condition")
	}

	gfields := ""
	if data.Exists("fields") {
		fields, err := data.String("fields")
		if nil != err {
			blog.Errorf("[api-compatiblev2] failed to parse the fields, error  info is %s", err.Error())
			return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}
		gfields = fields
	}

	return s.core.CompatibleV2Operation().Business(params).SearchAllApp(gfields, cond)
}

// UpdateMultiSet update multi sets
func (s *topoService) UpdateMultiSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("bizID", pathParams("appid"))
	bizID, err := paramPath.Int64("bizID")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the path params bizid(%s), error info is %s ", pathParams("appid"), err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	setIDS, exists := data.Get(common.BKSetIDField)
	if !exists {
		blog.Errorf("[api-compatiblev2] failed to get the set ids, the input data is %#v", data)
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, common.BKSetIDField)
	}

	innerData, err := data.MapStr("data")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to get the new data, the input data is %#v, error info is %s", data, err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).In(setIDS)

	err = s.core.CompatibleV2Operation().Set(params).UpdateMultiSet(bizID, innerData, cond)
	return nil, err
}

// DeleteMultiSet delete multi sets
func (s *topoService) DeleteMultiSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams("appid"), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setIDSStr, err := data.String(common.BKSetIDField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to get the set ids, the input data is %#v, error info is %s", data, err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	setIDSArr := strings.Split(setIDSStr, ",")
	setIDS, err := util.SliceStrToInt64(setIDSArr)
	if nil != err {
		blog.Errorf("[api-compatiblev2] the set id is invalid, the input set ids is %s, error info is %s", setIDSStr, err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	err = s.core.CompatibleV2Operation().Set(params).DeleteMultiSet(bizID, setIDS)
	return nil, err
}

// DeleteSetHost delete hosts in some sets
func (s *topoService) DeleteSetHost(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams("appid"), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setIDS, exists := data.Get(common.BKSetIDField)
	if !exists {
		blog.Errorf("[api-compatiblev2] failed to get the set ids, the input data is %#v", data)
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, common.BKSetIDField)
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).In(setIDS)
	err = s.core.CompatibleV2Operation().Set(params).DeleteSetHost(bizID, cond)
	return nil, err
}

// UpdateMultiModule update multi modules
func (s *topoService) UpdateMultiModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	innerData, err := data.MapStr("data")
	if nil != err {
		blog.Error("[api-compatiblev2] failed to parse the data, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	moduleIDS, exists := data.Get(common.BKModuleIDField)
	if !exists {
		blog.Error("[api-compatiblev2] failed to parse the module ids")
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, common.BKModuleIDField)
	}

	err = s.core.CompatibleV2Operation().Module(params).UpdateMultiModule(bizID, moduleIDS, innerData)
	return nil, err
}

// SearchModuleByApp search module by business
func (s *topoService) SearchModuleByApp(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	cond := &compatiblev2Condition{}
	if err := data.MarshalJSONInto(cond); nil != err {
		return nil, err
	}

	cond.Condition.Set(common.BKAppIDField, bizID)
	query := &metadata.QueryInput{}
	query.Condition = cond.Condition
	query.Fields = strings.Join(cond.Fields, ",")
	query.Start = cond.Page.Start
	query.Limit = cond.Page.Limit
	query.Sort = cond.Page.Sort

	return s.core.CompatibleV2Operation().Module(params).SearchModuleByApp(query)
}

// SearchModuleBySetProperty search module by set property
func (s *topoService) SearchModuleBySetProperty(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}
	cond := condition.CreateCondition()

	data.ForEach(func(key string, val interface{}) {
		cond.Field(key).In([]interface{}{val})
	})
	cond.Field(common.BKAppIDField).Eq(bizID)
	return s.core.CompatibleV2Operation().Module(params).SearchModuleBySetProperty(bizID, cond)
}

// AddMultiModule add multi modules
func (s *topoService) AddMultiModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := data.Int64(common.BKAppIDField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse business id, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	setID, err := data.Int64(common.BKSetIDField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse set id, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	moduleNameStr, err := data.String(common.BKModuleNameField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse module name, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	// prepare the data
	data.ForEach(func(key string, val interface{}) {
		switch key {
		case common.BKSetIDField, common.BKOperatorField, common.BKBakOperatorField, common.BKModuleTypeField:
			return
		}
		// clear the unused key
		data.Remove(key)
	})

	err = s.core.CompatibleV2Operation().Module(params).AddMultiModule(bizID, setID, strings.Split(moduleNameStr, ","), data)
	return nil, err
}

// DeleteMultiModule delete multi modules
func (s *topoService) DeleteMultiModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse business id, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	inputParams := &struct {
		BizID     int64   `json:"bk_biz_id"`
		ModuleIDS []int64 `json:"bk_module_id"`
	}{BizID: bizID}

	if err := data.MarshalJSONInto(inputParams); nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the data (%#v), error info is %s", data, err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	return nil, s.core.CompatibleV2Operation().Module(params).DeleteMultiModule(inputParams.BizID, inputParams.ModuleIDS)

}
