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
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateBusiness create a new business
func (s *topoService) CreateBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	return s.core.BusinessOperation().CreateBusiness(params, obj, data)
}

// DeleteBusiness delete the business
func (s *topoService) DeleteBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	return nil, s.core.BusinessOperation().DeleteBusiness(params, obj, bizID)
}

// UpdateBusiness update the business
func (s *topoService) UpdateBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	return nil, s.core.BusinessOperation().UpdateBusiness(params, data, obj, bizID)

}

// UpdateBusinessStatus update the business status
func (s *topoService) UpdateBusinessStatus(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	switch common.DataStatusFlag(pathParams("flag")) {
	case common.DataStatusDisabled, common.DataStatusEnable:
		data.Set(common.BKDataStatusField, pathParams("flag"))
	default:
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, pathParams("flag"))
	}

	return nil, s.core.BusinessOperation().UpdateBusiness(params, data, obj, bizID)
}

// SearchBusiness search the business by condition
func (s *topoService) SearchBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	innerCond := condition.CreateCondition()
	if err = innerCond.Parse(data); nil != err {
		blog.Errorf("[api-biz] failed to parse the input data, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cnt, instItems, err := s.core.BusinessOperation().FindBusiness(params, obj, innerCond)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil
}

// SearchDefaultBusiness search the business by condition
func (s *topoService) SearchDefaultBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	innerCond := condition.CreateCondition()
	if err = innerCond.Parse(data); nil != err {
		blog.Errorf("[api-biz] failed to parse the input data, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	cnt, instItems, err := s.core.BusinessOperation().FindBusiness(params, obj, innerCond)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}
	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// CreateDefaultBusiness create the default business
func (s *topoService) CreateDefaultBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	return s.core.BusinessOperation().CreateBusiness(params, obj, data)
}

func (s *topoService) GetInternalModule(params types.ContextParams, pathParams, queryparams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	obj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	_, result, err := s.core.BusinessOperation().GetInternalModule(params, obj, bizID)
	if nil != err {
		return nil, err
	}

	return result, nil
}
