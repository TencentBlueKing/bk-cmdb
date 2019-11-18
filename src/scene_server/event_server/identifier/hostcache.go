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
	"context"
	"encoding/json"
	"sort"

	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

func MarshalBinary(identifier *metadata.HostIdentifier) (data []byte, err error) {
	sort.Sort(metadata.HostIdentProcessSorter(identifier.Process))
	return json.Marshal(identifier)
}

func fillIdentifier(identifier *metadata.HostIdentifier, ctx context.Context, cache *redis.Client, db dal.RDB) *metadata.HostIdentifier {
	// fill cloudName
	cloud, err := getCache(ctx, cache, db, common.BKInnerObjIDPlat, identifier.CloudID, false)
	if err != nil {
		blog.Errorf("identifier: getCache error %s", err.Error())
		return identifier
	}
	identifier.CloudName = getString(cloud.data[common.BKCloudNameField])

	// fill module
	for moduleID := range identifier.HostIdentModule {
		biz, err := getCache(ctx, cache, db, common.BKInnerObjIDApp, identifier.HostIdentModule[moduleID].BizID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		identifier.HostIdentModule[moduleID].BizName = getString(biz.data[common.BKAppNameField])
		identifier.SupplierAccount = getString(biz.data[common.BKOwnerIDField])
		identifier.SupplierID = getInt(biz.data, common.BKSupplierIDField)

		set, err := getCache(ctx, cache, db, common.BKInnerObjIDSet, identifier.HostIdentModule[moduleID].SetID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		identifier.HostIdentModule[moduleID].SetName = getString(set.data[common.BKSetNameField])
		identifier.HostIdentModule[moduleID].SetEnv = getString(set.data[common.BKSetEnvField])
		identifier.HostIdentModule[moduleID].SetStatus = getString(set.data[common.BKSetStatusField])

		module, err := getCache(ctx, cache, db, common.BKInnerObjIDModule, identifier.HostIdentModule[moduleID].ModuleID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %s", err.Error())
			continue
		}
		identifier.HostIdentModule[moduleID].ModuleName = getString(module.data[common.BKModuleNameField])
	}

	// fill process
	for procIndex := range identifier.Process {
		process := &identifier.Process[procIndex]
		proc, err := getCache(ctx, cache, db, common.BKInnerObjIDProc, process.ProcessID, false)
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

	return identifier
}
