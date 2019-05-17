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

package inst

// FieldName the field name
type FieldName string

const (

	// HostID host id
	HostID = "bk_host_id"

	// PlatID the plat id
	PlatID = "bk_cloud_id"
	// Plat the object id
	Plat = "plat"
	// Business the business object
	Business = "biz"

	// Module the module object
	Module = "module"

	// Set the set object
	Set = "set"

	// BusinessID the business id
	BusinessID = "bk_biz_id"
	// BusinessNameField the business name
	BusinessNameField = "bk_biz_name"
	// InstID the common inst id
	InstID = "bk_inst_id"

	ParentID = "bk_parent_id"

	// InstName the common inst name
	InstName = "bk_inst_name"

	// SetID the set id
	SetID = "bk_set_id"

	// SetName the set name
	SetName = "bk_set_name"

	// ModuleID the module id
	ModuleID = "bk_module_id"

	// ModuleName the module name
	ModuleName = "bk_module_name"

	// PlatName the plat name
	PlatName = "bk_cloud_name"

	// DefaultLimit the limit num
	DefaultLimit = 1000
	// HostIDField the host id field
	HostIDField = "bk_host_id"
	// HostNameField the host name field
	HostNameField = "bk_host_name"

	// HostInnerIP the host innerip
	HostInnerIP = "bk_host_innerip"
)

// Maintaince the operation method
type Maintaince interface {
	IsExists() (bool, error)
	Create() error
	Update() error
	Save() error
}
