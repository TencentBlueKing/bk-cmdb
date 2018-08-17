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
	"encoding/json"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage"
)

type HostIdentifier struct {
	HostID          int64              `json:"bk_host_id" bson:"bk_host_id"`
	HostName        string             `json:"bk_host_name" bson:"bk_host_name"`
	SupplierID      int64              `json:"bk_supplier_id"`
	SupplierAccount string             `json:"bk_supplier_account"`
	CloudID         int64              `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName       string             `json:"bk_cloud_name" bson:"bk_cloud_name"`
	InnerIP         string             `json:"bk_host_innerip" bson:"bk_host_innerip"`
	OuterIP         string             `json:"bk_host_outerip" bson:"bk_host_outerip"`
	OSType          string             `json:"bk_os_type" bson:"bk_os_type"`
	OSName          string             `json:"bk_os_name" bson:"bk_os_name"`
	Memory          int64              `json:"bk_mem" bson:"bk_mem"`
	CPU             int64              `json:"bk_cpu" bson:"bk_cpu"`
	Disk            int64              `json:"bk_disk" bson:"bk_disk"`
	Module          map[string]*Module `json:"associations" bson:"associations"`
	Process         []Process          `json:"process" bson:"process"`
}

type Process struct {
	ProcessID       int64   `json:"bk_process_id" bson:"bk_process_id"`               // 进程名称
	ProcessName     string  `json:"bk_process_name" bson:"bk_process_name"`           // 进程名称
	BindIP          string  `json:"bind_ip" bson:"bind_ip"`                           // 绑定IP, 枚举: [{ID: "1", Name: "127.0.0.1"}, {ID: "2", Name: "0.0.0.0"}, {ID: "3", Name: "第一内网IP"}, {ID: "4", Name: "第一外网IP"}]
	PORT            string  `json:"port" bson:"port"`                                 // 端口, 单个端口："8080", 多个连续端口："8080-8089", 多个不连续端口："8080-8089,8199"
	PROTOCOL        string  `json:"protocol" bson:"protocol"`                         // 协议, 枚举: [{ID: "1", Name: "TCP"}, {ID: "2", Name: "UDP"}],
	FuncID          string  `json:"bk_func_id" bson:"bk_func_id"`                     // 功能ID
	FuncName        string  `json:"bk_func_name" bson:"bk_func_name"`                 // 功能名称
	StartParamRegex string  `json:"bk_start_param_regex" bson:"bk_start_param_regex"` // 启动参数匹配规则
	BindModules     []int64 `json:"bind_modules" bson:"bind_modules"`                 // 进程绑定的模块ID，数字数组
}

type Module struct {
	BizID      int64  `json:"bk_biz_id"`
	BizName    string `json:"bk_biz_name"`
	SetID      int64  `json:"bk_set_id"`
	SetName    string `json:"bk_set_name"`
	ModuleID   int64  `json:"bk_module_id"`
	ModuleName string `json:"bk_module_name"`
	SetStatus  string `json:"bk_service_status"`
	SetEnv     string `json:"bk_set_env"`
}

func (iden *HostIdentifier) MarshalBinary() (data []byte, err error) {
	return json.Marshal(iden)
}

func (iden *HostIdentifier) fillIden(cache *redis.Client, db storage.DI) *HostIdentifier {
	// fill cloudName
	cloud, err := getCache(cache, db, common.BKInnerObjIDPlat, iden.CloudID, false)
	if err != nil {
		blog.Errorf("identifier: getCache error %s", err.Error())
		return iden
	}
	iden.CloudName = getString(cloud.data[common.BKCloudNameField])

	// fill module
	for moduleID := range iden.Module {
		biz, err := getCache(cache, db, common.BKInnerObjIDApp, iden.Module[moduleID].BizID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].BizName = getString(biz.data[common.BKAppNameField])
		iden.SupplierAccount = getString(biz.data[common.BKOwnerIDField])
		iden.SupplierID = getInt(biz.data, common.BKSupplierIDField)

		set, err := getCache(cache, db, common.BKInnerObjIDSet, iden.Module[moduleID].SetID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].SetName = getString(set.data[common.BKSetNameField])
		iden.Module[moduleID].SetEnv = getString(set.data[common.BKSetEnvField])
		iden.Module[moduleID].SetStatus = getString(set.data[common.BKSetStatusField])

		module, err := getCache(cache, db, common.BKInnerObjIDModule, iden.Module[moduleID].ModuleID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].ModuleName = getString(module.data[common.BKModuleNameField])
	}

	// fill process
	for procindex := range iden.Process {
		process := &iden.Process[procindex]
		proc, err := getCache(cache, db, common.BKInnerObjIDProc, process.ProcessID, false)
		if err != nil {
			blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDProc, process.ProcessID, err.Error())
			continue
		}
		process.ProcessName = getString(proc.data[common.BKProcessNameField])
		process.FuncID = getString(proc.data[common.BKFuncIDField])
		process.FuncName = getString(proc.data[common.BKFuncName])
		process.BindIP = getString(proc.data[common.BKBindIP])
		process.PROTOCOL = getString(proc.data[common.BKProtocol])
		process.PORT = getString(proc.data[common.BKPort])
		process.StartParamRegex = getString(proc.data["bk_start_param_regex"])
	}

	return iden
}
