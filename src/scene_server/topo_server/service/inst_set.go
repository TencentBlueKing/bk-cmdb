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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateSet create a new set
func (s *Service) CreateSet(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("CreateSet")
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID).Field(common.BKObjIDField).Eq(common.BKInnerObjIDSet)

	objItems, err := s.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, pathParams("app_id"))

	for _, item := range objItems {
		return s.core.InstOperation().CreateInst(params, item, data) // should only item
	}

	return nil, nil
}

// DeleteSet delete the set
func (s *Service) DeleteSet(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDSet)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))
	cond.Field(common.BKSetIDField).Eq(pathParams("set_id"))

	objItems, err := s.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	for _, item := range objItems {
		if err = s.core.InstOperation().DeleteInst(params, item, cond); nil != err {
			return nil, err
		}
	}

	return nil, err
}

// UpdateSet update the set
func (s *Service) UpdateSet(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDSet)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))
	cond.Field(common.BKSetIDField).Eq(pathParams("set_id"))

	objItems, err := s.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, pathParams("app_id"))
	data.Set(common.BKSetIDField, pathParams("set_id"))

	for _, item := range objItems {
		if err = s.core.InstOperation().UpdateInst(params, data, item, cond); nil != err {
			return nil, err
		}
	}

	return nil, err
}

// SearchSet search the set
func (s *Service) SearchSet(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDSet)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))

	objItems, err := s.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}
	for _, objItem := range objItems {

		cnt, instItems, err := s.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-set] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
			return nil, err
		}
		count = count + cnt
		instRst = append(instRst, instItems...)
	}

	result := frtypes.MapStr{}
	result.Set("count", count)
	result.Set("info", instRst)

	return result, nil

}
