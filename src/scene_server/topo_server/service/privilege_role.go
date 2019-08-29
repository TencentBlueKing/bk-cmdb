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
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (s *Service) ParseCreateRolePrivilegeOriginData(data []byte) (mapstr.MapStr, error) {
	requestBody := struct {
		Data []string `json:"data" field:"data"`
	}{}
	err := json.Unmarshal(data, &requestBody)
	if nil != err {
		return nil, err
	}
	result := mapstr.MapStr{}
	result.Set("origin", requestBody.Data)
	return result, nil
}

// CreatePrivilege search user group
func (s *Service) CreatePrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	val, exists := data.Get("origin")
	if !exists {
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set anything")
	}

	privilegeData, ok := val.([]string)
	if !ok {
		blog.Infof("CreatePrivilege param invalidate type : %#v, rid: %s", val, params.ReqID)
	}

	err := s.Core.PermissionOperation().Role(params).CreatePermission(params.SupplierAccount, pathParams("bk_obj_id"), pathParams("bk_property_id"), privilegeData)
	return nil, err
}

// GetPrivilege search user group
func (s *Service) GetPrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	return s.Core.PermissionOperation().Role(params).GetPermission(params.SupplierAccount, pathParams("bk_obj_id"), pathParams("bk_property_id"))
}
