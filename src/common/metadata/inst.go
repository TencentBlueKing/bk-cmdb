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

type SetInst struct {
	SetID     int64  `bson:"bk_set_id"`
	SetName   string `bson:"bk_set_name"`
	SetStatus string `bson:"bk_service_status"`
	SetEnv    string `bson:"bk_set_env"`
}
type ModuleInst struct {
	ModuleID   int64  `bson:"bk_module_id"`
	ModuleName string `bson:"bk_module_name"`
}
type BizInst struct {
	BizID           int64  `bson:"bk_biz_id"`
	BizName         string `bson:"bk_biz_name"`
	SupplierID      int64  `bson:"bk_supplier_id"`
	SupplierAccount string `bson:"bk_supplier_account"`
}
type CloudInst struct {
	CloudID   int64  `bson:"bk_cloud_id"`
	CloudName string `bson:"bk_cloud_name"`
}

type HostIdentifier struct {
	// cache     *HostIdenCache
	HostID          int64                       `json:"bk_host_id" bson:"bk_host_id"`           // 主机ID(host_id)								数字
	HostName        string                      `json:"bk_host_name" bson:"bk_host_name"`       // 主机名称
	SupplierID      int64                       `json:"bk_supplier_id"`                         // 开发商ID（bk_supplier_id）				数字
	SupplierAccount string                      `json:"bk_supplier_account"`                    // 开发商帐号（bk_supplier_account）	数字
	CloudID         int64                       `json:"bk_cloud_id" bson:"bk_cloud_id"`         // 所属云区域id(bk_cloud_id)				数字
	CloudName       string                      `json:"bk_cloud_name" bson:"bk_cloud_name"`     // 所属云区域名称(bk_cloud_name)		字符串（最大长度25）
	InnerIP         string                      `json:"bk_host_innerip" bson:"bk_host_innerip"` // 内网IP
	OuterIP         string                      `json:"bk_host_outerip" bson:"bk_host_outerip"` // 外网IP
	OSType          string                      `json:"bk_os_type" bson:"bk_os_type"`           // 操作系统类型
	OSName          string                      `json:"bk_os_name" bson:"bk_os_name"`           // 操作系统名称
	Memory          int64                       `json:"bk_mem" bson:"bk_mem"`                   // 内存容量
	CPU             int64                       `json:"bk_cpu" bson:"bk_cpu"`                   // CPU逻辑核心数
	Disk            int64                       `json:"bk_disk" bson:"bk_disk"`                 // 磁盘容量
	HostIdentModule map[string]*HostIdentModule `json:"associations" bson:"associations"`
}

// Module HostIdentifier module define
type HostIdentModule struct {
	BizID      int64  `json:"bk_biz_id"`         // 业务ID
	BizName    string `json:"bk_biz_name"`       // 业务名称
	SetID      int64  `json:"bk_set_id"`         // 所属集群(bk_set_id)：						数字
	SetName    string `json:"bk_set_name"`       // 所属集群名称(bk_set_name)：			字符串（最大长度25）
	ModuleID   int64  `json:"bk_module_id"`      // 所属模块(bk_module_id)：				数字
	ModuleName string `json:"bk_module_name"`    // 所属模块(bk_module_name)：			字符串（最大长度25）
	SetStatus  string `json:"bk_service_status"` // 集群服务状态（bk_set_status）			数字
	SetEnv     string `json:"bk_set_env"`        // 环境类型（bk_set_type）					数字
}

// SearchIdentifierParam defines the param
type SearchIdentifierParam struct {
	IP   IPParam `json:"ip"`
	Page BasePage
}

type IPParam struct {
	Data    []string `json:"data"`
	CloudID *int64   `json:"bk_cloud_id"`
}

type SearchHostIdentifierResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int              `json:"count"`
		Info  []HostIdentifier `json:"info"`
	} `json:"data"`
}
