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
func (s *Service) SearchAllApp(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond, err := data.MapStr("condition")
	if nil != err {
		blog.Errorf("[api-compatiblev2] not set the condition in the search conditons, rid: %s", params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set the search condition")
	}

	gfields := ""
	if data.Exists("fields") {
		fields, err := data.String("fields")
		if nil != err {
			blog.Errorf("[api-compatiblev2] failed to parse the fields, error  info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}
		gfields = fields
	}

	return s.Core.CompatibleV2Operation().Business(params).SearchAllApp(gfields, cond)
}

// UpdateMultiSet update multi sets
func (s *Service) UpdateMultiSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("bizID", pathParams("appid"))
	bizID, err := paramPath.Int64("bizID")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the path params bizid(%s), error info is %s, rid: %s", pathParams("appid"), err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	setIDS, exists := data.Get(common.BKSetIDField)
	if !exists {
		blog.Errorf("[api-compatiblev2] failed to get the set ids, the input data is %#v, rid: %s", data, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, common.BKSetIDField)
	}

	innerData, err := data.MapStr("data")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to get the new data, the input data is %#v, error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).In(setIDS)

	err = s.Core.CompatibleV2Operation().Set(params).UpdateMultiSet(bizID, innerData, cond)
	return nil, err
}

// DeleteMultiSet delete multi sets
func (s *Service) DeleteMultiSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams("appid"), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setIDSStr, err := data.String(common.BKSetIDField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to get the set ids, the input data is %#v, error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	setIDSArr := strings.Split(setIDSStr, ",")
	setIDS, err := util.SliceStrToInt64(setIDSArr)
	if nil != err {
		blog.Errorf("[api-compatiblev2] the set id is invalid, the input set ids is %s, error info is %s, rid: %s", setIDSStr, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	err = s.Core.CompatibleV2Operation().Set(params).DeleteMultiSet(bizID, setIDS)
	return nil, err
}

// DeleteSetHost delete hosts in some sets
func (s *Service) DeleteSetHost(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams("appid"), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	requestBody := struct {
		SetIDs []int64 `json:"bk_set_id" field:"bk_set_id" bson:"bk_set_id"`
	}{}
	if err := data.MarshalJSONInto(&requestBody); err != nil {
		blog.Errorf("[api-compatiblev2] parse request body failed, the input data is %#v, rid: %s", data, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}

	err = s.Core.CompatibleV2Operation().Set(params).DeleteSetHost(bizID, requestBody.SetIDs)
	return nil, err
}

// UpdateMultiModule update multi modules
func (s *Service) UpdateMultiModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	innerData, err := data.MapStr("data")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the data, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	moduleIDS, exists := data.Get(common.BKModuleIDField)
	if !exists {
		blog.Errorf("[api-compatiblev2] failed to parse the module ids, rid: %s", params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, common.BKModuleIDField)
	}

	err = s.Core.CompatibleV2Operation().Module(params).UpdateMultiModule(bizID, moduleIDS, innerData)
	return nil, err
}

// SearchModuleByApp search module by business
func (s *Service) SearchModuleByApp(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	cond := &compatibleV2Condition{}
	if err := data.MarshalJSONInto(cond); nil != err {
		return nil, err
	}

	cond.Condition.Set(common.BKAppIDField, bizID)
	sortParser := metadata.NewSearchSortParse()
	sortArr := sortParser.String(cond.Page.Sort).ToSearchSortArr()
	inputParam := &metadata.QueryCondition{
		Fields: cond.Fields,
		Limit: metadata.SearchLimit{
			Offset: int64(cond.Page.Start),
			Limit:  int64(cond.Page.Limit),
		},
		SortArr:   sortArr,
		Condition: cond.Condition,
	}
	return s.Core.CompatibleV2Operation().Module(params).SearchModuleByApp(inputParam)
}

// SearchModuleBySetProperty search module by set property
func (s *Service) SearchModuleBySetProperty(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatible-v2]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}
	cond := condition.CreateCondition()

	data.ForEach(func(key string, val interface{}) error {
		cond.Field(key).In(util.ConverToInterfaceSlice(val))
		return nil
	})
	cond.Field(common.BKAppIDField).Eq(bizID)
	return s.Core.CompatibleV2Operation().Module(params).SearchModuleBySetProperty(bizID, cond)
}

// AddMultiModule add multi modules
func (s *Service) AddMultiModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := data.Int64(common.BKAppIDField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse business id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	setID, err := data.Int64(common.BKSetIDField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	moduleNameStr, err := data.String(common.BKModuleNameField)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse module name, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	// prepare the data
	_ = data.ForEach(func(key string, val interface{}) error {
		switch key {
		case common.BKSetIDField, common.BKOperatorField, common.BKBakOperatorField, common.BKModuleTypeField:
			return nil
		}
		// clear the unused key
		data.Remove(key)
		return nil
	})

	err = s.Core.CompatibleV2Operation().Module(params).AddMultiModule(bizID, setID, strings.Split(moduleNameStr, ","), data)
	return nil, err
}

// DeleteMultiModule delete multi modules
func (s *Service) DeleteMultiModule(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse business id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	inputParams := &struct {
		BizID     int64   `json:"bk_biz_id"`
		ModuleIDS []int64 `json:"bk_module_id"`
	}{BizID: bizID}

	if err := data.MarshalJSONInto(inputParams); nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the data (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	return nil, s.Core.CompatibleV2Operation().Module(params).DeleteMultiModule(inputParams.BizID, inputParams.ModuleIDS)

}
