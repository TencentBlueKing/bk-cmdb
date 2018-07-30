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
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (s *topoService) ParseCreateRolePrivilegeOriginData(data []byte) (frtypes.MapStr, error) {
	rst := []string{}
	err := json.Unmarshal(data, &rst)
	if nil != err {
		return nil, err
	}
	result := frtypes.MapStr{}
	result.Set("origin", rst)
	return result, nil
}

// CreatePrivilege search user goup
func (s *topoService) CreatePrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	datas := make([]string, 0)
	val, exists := data.Get("origin")
	if !exists {
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set anything")
	}

	datas, _ = val.([]string)
	err := s.core.PermissionOperation().Role(params).CreatePermission(params.SupplierAccount, pathParams("bk_obj_id"), pathParams("bk_property_id"), datas)
	return nil, err
}

// GetPrivilege search user goup
func (s *topoService) GetPrivilege(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	return s.core.PermissionOperation().Role(params).GetPermission(params.SupplierAccount, pathParams("bk_obj_id"), pathParams("bk_property_id"))
}
