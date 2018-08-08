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
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateUserGroup create user goup
func (s *topoService) CreateUserGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	userGroup := &metadata.UserGroup{}
	_, err := userGroup.Parse(data)
	if nil != err {
		blog.Errorf("[api-privilege] failed to parse the input data, error info is %s ", err.Error())
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	err = s.core.PermissionOperation().UserGroup(params).CreateUserGroup(params.SupplierAccount, userGroup)
	return nil, err
}

// DeleteUserGroup delete user goup
func (s *topoService) DeleteUserGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	err := s.core.PermissionOperation().UserGroup(params).DeleteUserGroup(pathParams("bk_supplier_account"), pathParams("group_id"))
	return nil, err
}

// UpdateUserGroup update user goup
func (s *topoService) UpdateUserGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	err := s.core.PermissionOperation().UserGroup(params).UpdateUserGroup(pathParams("bk_supplier_account"), pathParams("group_id"), data)
	return nil, err
}

// SearchUserGroup search user goup
func (s *topoService) SearchUserGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()

	data.ForEach(func(key string, val interface{}) {
		cond.Field(key).Like(val)
	})

	return s.core.PermissionOperation().UserGroup(params).SearchUserGroup(pathParams("bk_supplier_account"), cond)
}
