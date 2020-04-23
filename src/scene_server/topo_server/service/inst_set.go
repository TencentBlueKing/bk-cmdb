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
	"configcenter/src/scene_server/topo_server/core/operation"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateSet create a new set
func (s *Service) CreateSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	set, err := s.Core.SetOperation().CreateSet(params, obj, bizID, data)
	if err != nil {
		return nil, err
	}

	setID, err := set.GetInstID()
	if err != nil {
		blog.Errorf("unexpected error, create set success, but get id field failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	// auth: register set
	if err := s.AuthManager.RegisterSetByID(params.Context, params.Header, setID); err != nil {
		blog.Errorf("create set success,but register to iam failed, err:  %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}
	return set, nil
}

func (s *Service) DeleteSets(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDSet)

	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	cond := &operation.OpCondition{}
	if err = data.MarshalJSONInto(cond); nil != err {
		blog.Errorf("[api-set] failed to parse to the operation condition, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	// auth: deregister set
	if err := s.AuthManager.DeregisterSetByID(params.Context, params.Header, cond.Delete.InstID...); err != nil {
		blog.Errorf("delete sets failed, deregister sets from iam failed, %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	err = s.Core.SetOperation().DeleteSet(params, obj, bizID, cond.Delete.InstID)
	if err != nil {
		blog.Errorf("delete set %v failed, error info is %s, rid: %s", cond.Delete.InstID, err.Error(), params.ReqID)
		return nil, err
	}
	return nil, nil
}

// DeleteSet delete the set
func (s *Service) DeleteSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	if "batch" == pathParams("set_id") {
		return s.DeleteSets(params, pathParams, queryParams, data)
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDSet)

	if nil != err {
		blog.Errorf("delete set failed, failed to search the set, %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	// auth: deregister set
	if err := s.AuthManager.DeregisterSetByID(params.Context, params.Header, setID); err != nil {
		blog.Errorf("delete set failed, deregister set from iam failed, %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	err = s.Core.SetOperation().DeleteSet(params, obj, bizID, []int64{setID})
	if err != nil {
		blog.Errorf("delete sets failed, %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return nil, nil
}

// UpdateSet update the set
func (s *Service) UpdateSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	setID, err := strconv.ParseInt(pathParams("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the set id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "set id")
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("update set failed,failed to search the set, %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	err = s.Core.SetOperation().UpdateSet(params, data, obj, bizID, setID)
	if err != nil {
		blog.Errorf("update set failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	// auth: update register set
	if err := s.AuthManager.UpdateRegisteredSetByID(params.Context, params.Header, setID); err != nil {
		blog.Errorf("update set success, but update registered set failed, %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}
	return nil, nil
}

// SearchSet search the set
func (s *Service) SearchSet(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("[api-set]failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	paramsCond := &parser.SearchParams{
		Condition: mapstr.New(),
	}
	if err = data.MarshalJSONInto(paramsCond); nil != err {
		blog.Errorf("search set failed, decode parameter condition failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommParamsInvalid)
	}

	paramsCond.Condition[common.BKAppIDField] = bizID

	queryCond := &metadata.QueryInput{}
	queryCond.Condition = paramsCond.Condition
	queryCond.Fields = strings.Join(paramsCond.Fields, ",")
	page := metadata.ParsePage(paramsCond.Page)
	queryCond.Start = page.Start
	queryCond.Sort = page.Sort
	queryCond.Limit = page.Limit

	cnt, instItems, err := s.Core.SetOperation().FindSet(params, obj, queryCond)
	if nil != err {
		blog.Errorf("[api-set] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil

}
