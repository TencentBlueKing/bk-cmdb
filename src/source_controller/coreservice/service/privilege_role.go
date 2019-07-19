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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

// GetRolePri get role privilege
func (s *coreService) GetRolePri(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")
	propertyID := pathParams("bk_property_id")
	cond := make(map[string]interface{})
	cond[common.BKObjIDField] = objID
	cond[common.BKPropertyIDField] = propertyID
	var result map[string]interface{}
	cond = util.SetModOwner(cond, params.SupplierAccount)

	cnt, err := s.db.Table(common.BKTableNamePrivilege).Find(cond).Count(params.Context)
	if nil != err {
		blog.Errorf("get user group privi error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrObjectDBOpErrno, err.Error())
	}
	if 0 == cnt {
		blog.V(3).Infof("failed to find the cnt, rid: %s", params.ReqID)
		info := make([]string, 0)
		return info, nil
	}

	err = s.db.Table(common.BKTableNamePrivilege).Find(cond).One(params.Context, &result)
	if nil != err {
		blog.Errorf("get role pri field error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommDBSelectFailed, err.Error())
	}
	privilege, ok := result["privilege"]
	if !ok {
		blog.Errorf("not privilege, the origin data is %#v, rid: %s", result, params.ReqID)
		info := make([]string, 0)
		return info, nil
	}
	return privilege, nil
}

// CreateRolePri create role privilege
func (s *coreService) CreateRolePri(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")
	propertyID := pathParams("bk_property_id")
	requestBody := struct {
		Privileges []string `json:"privilege" field:"privilege" bson:"privilege"`
	}{}
	err := data.MarshalJSONInto(&requestBody)
	if err != nil {
		blog.Errorf("read json data error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	input := make(map[string]interface{})
	input[common.BKOwnerIDField] = params.SupplierAccount
	input[common.BKObjIDField] = objID
	input[common.BKPropertyIDField] = propertyID
	input[common.BKPrivilegeField] = requestBody.Privileges
	input = util.SetModOwner(input, params.SupplierAccount)

	err = s.db.Table(common.BKTableNamePrivilege).Insert(params.Context, input)
	if nil != err {
		blog.Errorf("create role privilege error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}

	return nil, nil
}

// UpdateRolePri update role privilege
func (s *coreService) UpdateRolePri(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")
	propertyID := pathParams("bk_property_id")
	requestBody := struct {
		Privileges []string `json:"privilege" field:"privilege" bson:"privilege"`
	}{}
	err := data.MarshalJSONInto(&requestBody)
	if err != nil {
		blog.Errorf("read json data error: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	input := make(map[string]interface{})
	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = params.SupplierAccount
	cond[common.BKObjIDField] = objID
	cond[common.BKPropertyIDField] = propertyID
	input[common.BKPrivilegeField] = requestBody.Privileges
	cond = util.SetModOwner(cond, params.SupplierAccount)

	// do update or create operation
	count, err := s.db.Table(common.BKTableNamePrivilege).Find(cond).Count(params.Context)
	if count == 0 {
		input[common.BKOwnerIDField] = params.SupplierAccount
		input[common.BKObjIDField] = objID
		input[common.BKPropertyIDField] = propertyID
		input[common.BKPrivilegeField] = requestBody.Privileges
		input = util.SetModOwner(input, params.SupplierAccount)
		err = s.db.Table(common.BKTableNamePrivilege).Insert(params.Context, input)
		if nil != err {
			blog.Errorf("create role privilege failed, err: %+v, rid: %s", err, params.ReqID)
			return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
		}
	} else {
		err = s.db.Table(common.BKTableNamePrivilege).Update(params.Context, cond, input)
		if nil != err {
			blog.Errorf("update role privilege failed, err: %v, rid: %s", err, params.ReqID)
			return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
		}
	}

	return nil, nil
}
