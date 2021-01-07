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

type PermissionSystemResponse struct {
	BaseResp `json:",inline"`
	Data     types.MapStr `json:"data"`
}

type PermissionGroupListResult struct {
	BaseResp `json:",inline"`
	Data     []UserGroup `json:"data"`
}

type Gprivilege struct {
	ModelConfig    map[string]map[string][]string `json:"model_config" bson:"model_config"`
	SysConfig      SysConfigStruct                `json:"sys_config,omitempty" bson:"sys_config"`
	IsHostCrossBiz bool                           `json:"is_host_cross_biz" bson:"is_host_cross_biz"`
}

type Privilege struct {
	ModelConfig map[string]map[string][]string `json:"model_config,omitempty" bson:"model_config"`
	SysConfig   *SysConfigStruct               `json:"sys_config,omitempty" bson:"sys_config"`
}

type SysConfigStruct struct {
	Globalbusi []string `json:"global_busi"`
	BackConfig []string `json:"back_config"`
}

type UserPrivilege struct {
	GroupID     string                         `json:"bk_group_id" bson:"bk_group_id"`
	ModelConfig map[string]map[string][]string `json:"model_config" bson:"model_config"`
	SysConfig   SysConfigStruct                `json:"sys_config" bson:"sys_config"`
}

type UserPriviResult struct {
	Result  bool          `json:"result"`
	Code    int           `json:"code"`
	Message interface{}   `json:"message"`
	Data    UserPrivilege `json:"data"`
}

type GroupPrivilege struct {
	GroupID   string     `json:"group_id" bson:"group_id"`
	OwnerID   string     `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Privilege *Privilege `json:"privilege"`
}

type GroupPriviResult struct {
	BaseResp `json:",inline"`
	Data     GroupPrivilege `json:"data"`
}

type SearchGroup struct {
	Code    int         `json:"code"`
	Result  bool        `json:"result"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type SearchMainLine struct {
	Code    int                      `json:"code"`
	Result  bool                     `json:"result"`
	Message interface{}              `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

// UserGroup the privilege user group definition
type UserGroup struct {
	GroupName       string `field:"group_name" json:"group_name" bson:"group_name"`
	UserList        string `field:"user_list" json:"user_list" bson:"user_list"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	GroupID         string `field:"group_id" json:"group_id" bson:"group_id"`
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
	SupplierAccount string       `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	GroupID         string       `field:"group_id" json:"group_id" bson:"bk_supplier_account"`
	ModelConfig     types.MapStr `field:"model_config" json:"model_config" bson:"model_config"`
	SystemConfig    types.MapStr `field:"sys_config" json:"sys_config" bson:"sys_config"`
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
