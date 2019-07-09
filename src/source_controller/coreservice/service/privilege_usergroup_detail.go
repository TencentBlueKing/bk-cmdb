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

// CreateUserGroupPrivi create group privi
func (s *coreService) CreateUserGroupPrivi(params core.ContextParams, pathParams, queryParams ParamsGetter, info mapstr.MapStr) (interface{}, error) {
	groupID := pathParams("group_id")
	data := make(map[string]interface{})
	data[common.BKUserGroupIDField] = groupID
	data[common.BKPrivilegeField] = info
	data = util.SetModOwner(data, params.SupplierAccount)

	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = params.SupplierAccount
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, params.SupplierAccount)
	cnt, err := s.db.Table(common.BKTableNameUserGroupPrivilege).Find(cond).Count(params.Context)
	if nil != err && !s.db.IsNotFoundError(err) {
		blog.Errorf("get user group privi error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}
	if cnt > 0 {
		blog.V(3).Infof("update user group privi: %+v, by condition %+v, rid: %s", data, cond, params.ReqID)
		err = s.db.Table(common.BKTableNameUserGroupPrivilege).Update(params.Context, cond, data)
		if nil != err {
			blog.Errorf("update user group privi error :%v, rid: %s", err, params.ReqID)
			return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
		}
		return nil, nil
	}

	blog.V(3).Infof("create user group privi: %+v, rid: %s", data, params.ReqID)
	err = s.db.Table(common.BKTableNameUserGroupPrivilege).Insert(params.Context, data)
	if nil != err {
		blog.Errorf("insert user group privi error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}

	return nil, nil
}

// UpdateUserGroupPrivi update group privi
func (s *coreService) UpdateUserGroupPrivi(params core.ContextParams, pathParams, queryParams ParamsGetter, info mapstr.MapStr) (interface{}, error) {
	groupID := pathParams("group_id")
	cond := make(map[string]interface{})
	data := make(map[string]interface{})
	cond[common.BKUserGroupIDField] = groupID
	data[common.BKPrivilegeField] = info
	cond = util.SetModOwner(cond, params.SupplierAccount)
	blog.V(3).Infof("update user group privi: %+v, by condition %+v, rid: %s", data, cond, params.ReqID)
	err := s.db.Table(common.BKTableNameUserGroupPrivilege).Update(params.Context, cond, data)
	if nil != err {
		blog.Errorf("update user group privi error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}
	return nil, nil
}

// GetUserGroupPrivi get group privi
func (s *coreService) GetUserGroupPrivi(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	groupID := pathParams("group_id")

	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = params.SupplierAccount
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, params.SupplierAccount)

	blog.V(3).Infof("get user group privi by condition %+v, rid: %s", cond, params.ReqID)
	cnt, err := s.db.Table(common.BKTableNameUserGroupPrivilege).Find(cond).Count(params.Context)
	if nil != err && !s.db.IsNotFoundError(err) {
		blog.Errorf("get user group privi error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}
	if 0 == cnt {
		data := make(map[string]interface{})
		data[common.BKOwnerIDField] = params.SupplierAccount
		data[common.BKUserGroupIDField] = groupID
		data[common.BKPrivilegeField] = common.KvMap{}
		blog.V(3).Infof("get user group privi by condition %+v, returns %+v, rid: %s", cond, data, params.ReqID)
		return data, nil
	}

	var result interface{}
	err = s.db.Table(common.BKTableNameUserGroupPrivilege).Find(cond).One(params.Context, &result)
	if nil != err {
		blog.Errorf("get user group privi error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}

	blog.V(3).Infof("get user group privi by condition %+v, returns %+v, rid: %s", cond, result, params.ReqID)
	return result, nil
}
