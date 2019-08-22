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
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// UpdateUserGroupPrivilege search user group
func (s *Service) UpdateUserGroupPrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	privilegeData := &metadata.PrivilegeUserGroup{}
	if err := data.MarshalJSONInto(privilegeData); nil != err {
		blog.Errorf("[api-privilege] failed to parse the input data, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	err := s.Core.PermissionOperation().Permission(params).SetUserGroupPermission(params.SupplierAccount, pathParams("group_id"), privilegeData)
	return nil, err
}

// GetUserGroupPrivilege search user group
func (s *Service) GetUserGroupPrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	return s.Core.PermissionOperation().Permission(params).GetUserGroupPermission(params.SupplierAccount, pathParams("group_id"))
}

// GetUserPrivilege search user group
func (s *Service) GetUserPrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	return s.Core.PermissionOperation().Permission(params).GetUserPermission(params.SupplierAccount, pathParams("user_name"))
}
