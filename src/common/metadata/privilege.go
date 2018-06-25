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

package metadata

import (
	types "configcenter/src/common/mapstr"
)

// UserGroup the privilege user group definition
type UserGroup struct {
	GroupName       string `field:"group_name" json:"group_name"`
	UserList        string `field:"user_list" json:"user_list"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account"`
	GroupID         string `field:"group_id" json:"group_id"`
}

// Parse load the data from mapstr object into object instance
func (u *UserGroup) Parse(data types.MapStr) (*UserGroup, error) {

	err := SetValueToStructByTags(u, data)
	if nil != err {
		return nil, err
	}

	return u, err
}

// ToMapStr to mapstr
func (u *UserGroup) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(u)
}

// PrivilegeUserGroup the user group permission configure
type PrivilegeUserGroup struct {
	SupplierAccount string       `field:"bk_supplier_account" json:"bk_supplier_account"`
	GroupID         string       `field:"group_id" json:"group_id"`
	ModelConfig     types.MapStr `field:"model_config" json:"model_config"`
	SystemConfig    types.MapStr `field:"sys_config" json:"sys_config"`
}

// Parse load the data from mapstr object into object instance
func (p *PrivilegeUserGroup) Parse(data types.MapStr) (*PrivilegeUserGroup, error) {

	err := SetValueToStructByTags(p, data)
	if nil != err {
		return nil, err
	}

	return p, err
}

// ToMapStr to mapstr
func (p *PrivilegeUserGroup) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(p)
}
