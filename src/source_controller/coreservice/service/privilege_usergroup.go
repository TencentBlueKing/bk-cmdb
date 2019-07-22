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

	"github.com/rs/xid"
)

// CreateUserGroup create group
func (s *coreService) CreateUserGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	guid := xid.New()
	data[common.BKUserGroupIDField] = guid.String()
	data = util.SetModOwner(data, params.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserGroup).Insert(params.Context, data)
	if nil != err {
		blog.Errorf("create user group error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}

	return nil, nil
}

// UpdateUserGroup create group
func (s *coreService) UpdateUserGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	groupID := pathParams("group_id")
	cond := make(map[string]interface{})
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, params.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserGroup).Update(params.Context, cond, data)
	if nil != err {
		blog.Errorf("update user group error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}
	return nil, nil
}

// DeleteUserGroup create group
func (s *coreService) DeleteUserGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	groupID := pathParams("group_id")
	cond := make(map[string]interface{})
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, params.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserGroup).Delete(params.Context, cond)
	if nil != err {
		blog.Errorf("delete user group error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}
	return nil, nil
}

// SearchUserGroup create group
func (s *coreService) SearchUserGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cond := util.SetModOwner(data, params.SupplierAccount)
	var result []interface{}
	err := s.db.Table(common.BKTableNameUserGroup).Find(cond).All(params.Context, &result)
	if nil != err {
		blog.Errorf("get user group error :%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrObjectDBOpErrno)
	}

	return result, nil
}
