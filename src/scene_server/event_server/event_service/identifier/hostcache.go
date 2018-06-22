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

package identifier

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"encoding/json"
	"fmt"
)

// HostIdentifier define
type HostIdentifier struct {
	// cache     *HostIdenCache
	HostID          int                `json:"bk_host_id" bson:"bk_host_id"`           // 主机ID(host_id)								数字
	HostName        string             `json:"bk_host_name" bson:"bk_host_name"`       // 主机名称
	SupplierID      int                `json:"bk_supplier_id"`                         // 开发商ID（bk_supplier_id）				数字
	SupplierAccount string             `json:"bk_supplier_account"`                    // 开发商帐号（bk_supplier_account）	数字
	CloudID         int                `json:"bk_cloud_id" bson:"bk_cloud_id"`         // 所属云区域id(bk_cloud_id)				数字
	CloudName       string             `json:"bk_cloud_name" bson:"bk_cloud_name"`     // 所属云区域名称(bk_cloud_name)		字符串（最大长度25）
	InnerIP         string             `json:"bk_host_innerip" bson:"bk_host_innerip"` // 内网IP
	OuterIP         string             `json:"bk_host_outerip" bson:"bk_host_outerip"` // 外网IP
	OSType          string             `json:"bk_os_type" bson:"bk_os_type"`           // 操作系统类型
	OSName          string             `json:"bk_os_name" bson:"bk_os_name"`           // 操作系统名称
	Memory          int64              `json:"bk_mem" bson:"bk_mem"`                   // 内存容量
	CPU             int64              `json:"bk_cpu" bson:"bk_cpu"`                   // CPU逻辑核心数
	Disk            int64              `json:"bk_disk" bson:"bk_disk"`                 // 磁盘容量
	Module          map[string]*Module `json:"associations" bson:"associations"`
}

// Module HostIdentifier module define
type Module struct {
	BizID      int    `json:"bk_biz_id"`         // 业务ID
	BizName    string `json:"bk_biz_name"`       // 业务名称
	SetID      int    `json:"bk_set_id"`         // 所属集群(bk_set_id)：						数字
	SetName    string `json:"bk_set_name"`       // 所属集群名称(bk_set_name)：			字符串（最大长度25）
	ModuleID   int    `json:"bk_module_id"`      // 所属模块(bk_module_id)：				数字
	ModuleName string `json:"bk_module_name"`    // 所属模块(bk_module_name)：			字符串（最大长度25）
	SetStatus  string `json:"bk_service_status"` // 集群服务状态（bk_set_status）			数字
	SetEnv     string `json:"bk_set_env"`        // 环境类型（bk_set_type）					数字
}

// MarshalBinary implement MarshalBinary interface
func (iden *HostIdentifier) MarshalBinary() (data []byte, err error) {
	return json.Marshal(iden)
}

func (iden *HostIdentifier) fillIden() *HostIdentifier {

	for moduleID := range iden.Module {

		biz, err := getCache(common.BKInnerObjIDApp, iden.Module[moduleID].BizID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].BizName = fmt.Sprint(biz.data[common.BKAppNameField])
		iden.SupplierAccount = fmt.Sprint(biz.data[common.BKOwnerIDField])
		iden.SupplierID = getInt(biz.data, common.BKSupplierIDField)

		set, err := getCache(common.BKInnerObjIDSet, iden.Module[moduleID].SetID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].SetName = fmt.Sprint(set.data[common.BKSetNameField])
		iden.Module[moduleID].SetEnv = fmt.Sprint(set.data[common.BKSetEnvField])
		iden.Module[moduleID].SetStatus = fmt.Sprint(set.data[common.BKSetStatusField])

		module, err := getCache(common.BKInnerObjIDModule, iden.Module[moduleID].ModuleID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].ModuleName = fmt.Sprint(module.data[common.BKModuleNameField])

	}
	cloud, err := getCache(common.BKInnerObjIDPlat, iden.CloudID, false)
	if err != nil {
		blog.Errorf("identifier: getCache error %s", err.Error())
		return iden
	}
	iden.CloudName = fmt.Sprint(cloud.data[common.BKCloudNameField])

	return iden
}
