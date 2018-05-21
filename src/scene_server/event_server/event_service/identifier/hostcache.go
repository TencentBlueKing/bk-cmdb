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
)

func (c *HostIdenCache) getCache(by string, id int) HostIdenSlice {
	switch by {
	case common.BKInnerObjIDApp:
		return c.bizCache[id]
	case common.BKInnerObjIDSet:
		return c.setCache[id]
	case common.BKInnerObjIDModule:
		return c.moduleCache[id]
	case common.BKInnerObjIDPlat:
		return c.platCache[id]
	case common.BKInnerObjIDHost:
		return HostIdenSlice{c.hostCache[id]: true}
	}
	panic("unknow type")
}

// HostIdenCache cache
type HostIdenSlice map[*HostIdentifier]bool

type HostIdenCache struct {
	bizCache    map[int]HostIdenSlice
	setCache    map[int]HostIdenSlice
	moduleCache map[int]HostIdenSlice
	platCache   map[int]HostIdenSlice
	hostCache   map[int]*HostIdentifier
}

func (c *HostIdenCache) AddHost(i *HostIdentifier) {
	c.hostCache[i.HostID] = i
	// c.bizCache[i.BizID][i] = true
	// c.setCache[i.SetID][i] = true
	// c.moduleCache[i.ModuleID][i] = true
	c.platCache[i.CloudID][i] = true
}

// HostIdentifier define
type HostIdentifier struct {
	// cache     *HostIdenCache
	HostID    int            `json:"bk_host_id"`      // 主机ID(host_id)								数字
	HostName  string         `json:"bk_host_name"`    // 主机名称
	CloudID   int            `json:"bk_cloud_id"`     // 所属云区域id(bk_cloud_id)				数字
	CloudName string         `json:"bk_cloud_name"`   // 所属云区域名称(bk_cloud_name)		字符串（最大长度25）
	InnerIP   string         `json:"bk_host_innerip"` // 内网IP
	OuterIP   string         `json:"bk_host_outerip"` // 外网IP
	OSType    string         `json:"bk_os_type"`      // 操作系统类型
	OSName    string         `json:"bk_os_name"`      // 操作系统名称
	Memory    string         `json:"bk_mem"`          // 内存容量
	CPU       string         `json:"bk_cpu"`          // CPU逻辑核心数
	Disk      string         `json:"bk_disk"`         // 磁盘容量
	Module    map[int]Module `json:"belong"`
}

type Module struct {
	SupplierID      int    `json:"bk_supplier_id"`      // 开发商ID（bk_supplier_id）				数字
	SupplierAccount int    `json:"bk_supplier_account"` // 开发商帐号（bk_supplier_account）	数字
	BizID           int    `json:"bk_biz_id"`           // 业务ID
	BizName         string `json:"bk_biz_name"`         // 业务名称
	SetID           int    `json:"bk_set_id"`           // 所属集群(bk_set_id)：						数字
	SetName         string `json:"bk_set_name"`         // 所属集群名称(bk_set_name)：			字符串（最大长度25）
	ModuleID        int    `json:"bk_module_id"`        // 所属模块(bk_module_id)：				数字
	ModuleName      string `json:"bk_module_name"`      // 所属模块(bk_module_name)：			字符串（最大长度25）
	SetStatus       string `json:"bk_service_status"`   // 集群服务状态（bk_set_status）			数字
	SetEnv          string `json:"bk_set_env"`          // 环境类型（bk_set_type）					数字
}

// Set set inst value
// func (i *HostIdentifier) Set(key string, value interface{}) {
// 	switch key {
// 	case "bk_host_name":
// 		i.HostName = value
// 	case "bk_cloud_name":
// 		i.CloudName = value
// 	case "bk_biz_name":
// 		i.BizName = value
// 	case "bk_set_name":
// 		i.SetName = value
// 	case "bk_module_name":
// 		i.ModuleName = value
// 	case "bk_service_status":
// 		i.SetStatus = value
// 	case "bk_set_env":
// 		i.SetEnv = value
// 	case "bk_host_innerip":
// 		i.InnerIP = value
// 	case "bk_host_outerip":
// 		i.OuterIP = value
// 	case "bk_os_type":
// 		i.OSType = value
// 	case "bk_os_name":
// 		i.OSName = value
// 	case "bk_mem":
// 		i.Memory = value
// 	case "bk_cpu":
// 		i.CPU = value
// 	case "bk_disk":
// 		i.Disk = value
// 	}
// }

// func (i *HostIdentifier) SetBizID(id int) {
// 	if id == i.BizID {
// 		return
// 	}
// 	delete(i.cache.bizCache[i.BizID], i)
// 	i.BizID = id
// 	i.cache.bizCache[id][i] = true
// }

// func (i *HostIdentifier) SetSetID(id int) {
// 	if id == i.SetID {
// 		return
// 	}
// 	delete(i.cache.setCache[i.SetID], i)
// 	i.SetID = id
// 	i.cache.setCache[id][i] = true
// }

// func (i *HostIdentifier) SetModuleID(id int) {
// 	if id == i.ModuleID {
// 		return
// 	}
// 	delete(i.cache.moduleCache[i.ModuleID], i)
// 	i.ModuleID = id
// 	i.cache.moduleCache[id][i] = true
// }

// func (i *HostIdentifier) SetPlatID(id int) {
// 	if id == i.CloudID {
// 		return
// 	}
// 	delete(i.cache.platCache[i.CloudID], i)
// 	i.CloudID = id
// 	i.cache.platCache[id][i] = true
// }

// func (i *HostIdentifier) ModuleTransfer(bizID, setID, moduleID int) {
// 	i.SetBizID(bizID)
// 	i.SetSetID(setID)
// 	i.SetModuleID(moduleID)
// }
