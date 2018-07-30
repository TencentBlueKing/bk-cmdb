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

package models

//common group
const (
	BaseInfo     = "default"
	BaseInfoName = "基础信息"

	GroupNone = "none"
)

//app group info
const (
	AppRole     = "role"
	AppRoleName = "角色"
)

//host group info

const (
	HostAutoFields     = "auto"
	HostAutoFieldsName = "自动发现信息（需要安装agent）"
)

//process group info
const (
	ProcPort     = "proc_port"
	ProcPortName = "监听端口"

	ProcGsekitBaseInfo     = "gsekit_baseinfo"
	ProcGsekitBaseInfoName = "GSEkit 基本信息"

	ProcGsekitManageInfo     = "gsekit_manage"
	ProcGsekitManageInfoName = "GSEkit 进程管理信息"
)
