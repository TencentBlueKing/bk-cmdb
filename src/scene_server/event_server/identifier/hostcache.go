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
	"fmt"

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

	for moduleID := range iden.Module {

		biz, err := getCache(cache, db, common.BKInnerObjIDApp, iden.Module[moduleID].BizID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].BizName = fmt.Sprint(biz.data[common.BKAppNameField])
		iden.SupplierAccount = fmt.Sprint(biz.data[common.BKOwnerIDField])
		iden.SupplierID = getInt(biz.data, common.BKSupplierIDField)

		set, err := getCache(cache, db, common.BKInnerObjIDSet, iden.Module[moduleID].SetID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].SetName = fmt.Sprint(set.data[common.BKSetNameField])
		iden.Module[moduleID].SetEnv = fmt.Sprint(set.data[common.BKSetEnvField])
		iden.Module[moduleID].SetStatus = fmt.Sprint(set.data[common.BKSetStatusField])

		module, err := getCache(cache, db, common.BKInnerObjIDModule, iden.Module[moduleID].ModuleID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		iden.Module[moduleID].ModuleName = fmt.Sprint(module.data[common.BKModuleNameField])

	}
	cloud, err := getCache(cache, db, common.BKInnerObjIDPlat, iden.CloudID, false)
	if err != nil {
		blog.Errorf("identifier: getCache error %s", err.Error())
		return iden
	}
	iden.CloudName = fmt.Sprint(cloud.data[common.BKCloudNameField])

	return iden
}
