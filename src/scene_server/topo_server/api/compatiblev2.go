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

package api

import (
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initCompatiblev2)
}

func (cli *topoAPI) initCompatiblev2() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/searchAll", HandlerFunc: cli.SearchAllApp})

	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/openapi/set/multi/{appid}", HandlerFunc: cli.UpdateMultiSet})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/openapi/set/multi/{appid}", HandlerFunc: cli.DeleteMultiSet})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/openapi/set/setHost/{appid}", HandlerFunc: cli.DeleteSetHost})

	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/openapi/module/multi/{" + common.BKAppIDField + "}", HandlerFunc: cli.UpdateMultiModule})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/openapi/module/searchByApp/{" + common.BKAppIDField + "}", HandlerFunc: cli.SearchModuleByApp})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/openapi/module/searchByProperty/{" + common.BKAppIDField + "}", HandlerFunc: cli.SearchModuleBySetProperty})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/openapi/module/multi", HandlerFunc: cli.AddMultiModule})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/openapi/module/multi/{" + common.BKAppIDField + "}", HandlerFunc: cli.DeleteMultiModule})

}

func (cli *topoAPI) SearchAllApp(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

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

	return cli.core.CompatibleV2Operation().Business(params).SearchAllApp(gfields, cond)
}

func (cli *topoAPI) UpdateMultiSet(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

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

	err = cli.core.CompatibleV2Operation().Set(params).UpdateMultiSet(bizID, innerData, cond)
	return nil, err
}

func (cli *topoAPI) DeleteMultiSet(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("bizID", pathParams("appid"))
	bizID, err := paramPath.Int64("bizID")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the path params bizid(%s), error info is %s ", pathParams("appid"), err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
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

	err = cli.core.CompatibleV2Operation().Set(params).DeleteMultiSet(bizID, setIDS)
	return nil, err
}

func (cli *topoAPI) DeleteSetHost(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

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

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).In(setIDS)

	err = cli.core.CompatibleV2Operation().Set(params).DeleteSetHost(bizID, cond)
	return nil, err
}

func (cli *topoAPI) UpdateMultiModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("bizID", pathParams(common.BKAppIDField))
	bizID, err := paramPath.Int64("bizID")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the path params bizid(%s), error info is %s ", pathParams(common.BKAppIDField), err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
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

	err = cli.core.CompatibleV2Operation().Module(params).UpdateMultiModule(bizID, moduleIDS, innerData)
	return nil, err
}

func (cli *topoAPI) SearchModuleByApp(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("bizID", pathParams(common.BKAppIDField))
	bizID, err := paramPath.Int64("bizID")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the path params bizid(%s), error info is %s ", pathParams(common.BKAppIDField), err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	cond, err := data.MapStr("condition")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the condition, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	page, err := data.MapStr("page")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the condition, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	intoPage := &metadata.BasePage{}
	if err = page.MarshalJSONInto(intoPage); nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse page , error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	fields, err := data.MapStr("fields")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the condition, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	intoFIelds := make([]string, 0)
	if err = fields.MarshalJSONInto(intoFIelds); nil != err {
		blog.Errorf("[api-compatiblev2] faied to parse the fields, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	cond.Set(common.BKAppIDField, bizID)

	query := &metadata.QueryInput{}
	query.Condition = cond
	query.Fields = strings.Join(intoFIelds, ",")
	query.Start = intoPage.Start
	query.Limit = intoPage.Limit
	query.Sort = intoPage.Sort

	return cli.core.CompatibleV2Operation().Module(params).SearchModuleByApp(query)
}

func (cli *topoAPI) SearchModuleBySetProperty(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("bizID", pathParams(common.BKAppIDField))
	bizID, err := paramPath.Int64("bizID")
	if nil != err {
		blog.Errorf("[api-compatiblev2] failed to parse the path params bizid(%s), error info is %s ", pathParams(common.BKAppIDField), err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	cond := condition.CreateCondition()

	data.ForEach(func(key string, val interface{}) {
		cond.Field(key).In(val)
	})

	return cli.core.CompatibleV2Operation().Module(params).SearchModuleBySetProperty(bizID, cond)
}

func (cli *topoAPI) AddMultiModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

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

	return cli.core.CompatibleV2Operation().Module(params).AddMultiModule(bizID, setID, strings.Split(moduleNameStr, ","), cond)
}

func (cli *topoAPI) DeleteMultiModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	return nil, nil
}
