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

package api

// plat
const (
	fieldObjectID = "bk_obj_id"
	plat          = "plat"
)

// set fields
const (
	fieldParentID        = "bk_parent_id"
	fieldSetID           = "bk_set_id"
	fieldSetName         = "bk_set_name"
	fieldPlatID          = "bk_cloud_id"
	fieldPlatName        = "bk_cloud_name"
	fieldSupplierAccount = "bk_supplier_account"
	fieldSupplierID      = "bk_supplier_id"
	fieldBusinessID      = "bk_biz_id"
	fieldCapacity        = "bk_capacity"
	fieldServiceStatus   = "bk_service_status"
	fieldSetDesc         = "bk_set_desc"
	fieldSetEnv          = "bk_set_env"
	fieldObjID           = "bk_obj_id"
	fieldDescription     = "description"
)

// module fields
const (
	fieldModuleID    = "bk_module_id"
	fieldModuleName  = "bk_module_name"
	fieldBakOperator = "bk_bak_operator"
	fieldModuleTYpe  = "bk_module_type"
	fieldOperator    = "bk_operator"
)

// business fields
const (
	fieldBizDeveloper  = "bk_biz_developer"
	fieldBizID         = "bk_biz_id"
	fieldBizMaintainer = "bk_biz_maintainer"
	fieldBizName       = "bk_biz_name"
	fieldBizProductor  = "bk_biz_productor"
	fieldBizTester     = "bk_biz_tester"
	fieldLifeCycle     = "life_cycle"
	fieldBizOperator   = "operator"
)

// host fields
const (
	fieldOsBit        = "bk_os_bit"
	fieldSLA          = "bk_sla"
	fieldCloudID      = "bk_cloud_id"
	fieldHostInnerIP  = "bk_host_innerip"
	fieldCPU          = "bk_cpu"
	fieldCPUMhz       = "bk_cpu_mhz"
	fieldOsType       = "bk_os_type"
	fieldDisk         = "bk_disk"
	fieldHostID       = "bk_host_id"
	fieldHostOuterIP  = "bk_host_outerip"
	fieldAssetID      = "bk_asset_id"
	fieldMac          = "bk_mac"
	fieldProvinceName = "bk_provinceName"
	fieldSN           = "bk_sn"
	fieldCPUModule    = "bk_cpu_module"
	fieldHostName     = "bk_host_name"
	fieldISPName      = "bk_isp_name"
	fieldOuterMac     = "bk_outer_mac"
	fieldServiceTerm  = "bk_service_term"
	fieldComment      = "bk_comment"
	fieldMem          = "bk_mem"
	fieldOsName       = "bk_os_name"
	fieldOsVersion    = "bk_os_version"
	fieldImportFrom   = "import_from"
	fieldHostOperator = "operator"
)

// Enum definition
const (
	HostSLALevel1            = "1"
	HostSLALevel2            = "2"
	HostSLALevel3            = "3"
	HostOSTypeLinux          = "1"
	HostOSTypeWindows        = "2"
	HostImportFromExcel      = "1"
	HostImportFromAgent      = "2"
	HostImportFromAPI        = "3"
	BusinessLifeCycleTesting = "1"
	BusinessLifeCycleOnLine  = "2"
	BusinessLifeCycleStopped = "3"
	SetEnvTesting            = "1"
	SetEnvGuest              = "2"
	SetEnvNormal             = "3"
	SetServiceOpen           = "1"
	SetServiceClose          = "2"
)

type HostModuleActionType string

const (
	HostAppendModule  HostModuleActionType = "append"
	HostReplaceModule HostModuleActionType = "replace"
)
