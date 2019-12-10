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

	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

func fillIdentifier(identifier *metadata.HostIdentifier, ctx context.Context, cache *redis.Client, db dal.RDB) (*metadata.HostIdentifier, error) {
	// fill cloudName
	cloud, err := getCache(ctx, cache, db, common.BKInnerObjIDPlat, identifier.CloudID, false)
	if err != nil {
		blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDPlat, identifier.CloudID, err.Error())
		return nil, err
	}
	identifier.CloudName = getString(cloud.data[common.BKCloudNameField])

	// fill module
	for _, hostIdentModule := range identifier.HostIdentModule {
		biz, err := getCache(ctx, cache, db, common.BKInnerObjIDApp, hostIdentModule.BizID, false)
		if err != nil {
			blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDApp, hostIdentModule.BizID, err.Error())
			return nil, err
		}
		hostIdentModule.BizName = getString(biz.data[common.BKAppNameField])
		identifier.SupplierAccount = getString(biz.data[common.BKOwnerIDField])
		identifier.SupplierID, err = getInt(biz.data, common.BKSupplierIDField)
		if err != nil {
			blog.Errorf("identifier: convert instID failed the raw is %+v", biz.data[common.BKSupplierIDField])
			return nil, err
		}

		set, err := getCache(ctx, cache, db, common.BKInnerObjIDSet, hostIdentModule.SetID, false)
		if err != nil {
			blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDSet, hostIdentModule.SetID, err.Error())
			return nil, err
		}
		hostIdentModule.SetName = getString(set.data[common.BKSetNameField])
		hostIdentModule.SetEnv = getString(set.data[common.BKSetEnvField])
		hostIdentModule.SetStatus = getString(set.data[common.BKSetStatusField])

		module, err := getCache(ctx, cache, db, common.BKInnerObjIDModule, hostIdentModule.ModuleID, false)
		if err != nil {
			blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDModule, hostIdentModule.ModuleID, err.Error())
			return nil, err
		}
		hostIdentModule.ModuleName = getString(module.data[common.BKModuleNameField])

		// fill host layer info
		asstMap := make(map[string]string)
		asstArr := make([]metadata.Association, 0)
		cond := condition.CreateCondition().Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)
		condMap := util.SetQueryOwner(cond.ToMapStr().ToMapInterface(), identifier.SupplierAccount)
		err = db.Table(common.BKTableNameObjAsst).Find(condMap).All(ctx, &asstArr)
		if err != nil {
			blog.ErrorJSON("findHostLayerInfo query mainline association info error. condition:%s", condMap)
			return nil, err
		}
		for _, asst := range asstArr {
			asstMap[asst.ObjectID] = asst.AsstObjID
		}

		parentID, err := getInt(set.data, common.BKParentIDField)
		if err != nil {
			blog.Errorf("identifier: convert set bk_parent_id failed, the raw is %+v", set.data[common.BKParentIDField])
			return nil, err
		}

		curObj, ok := asstMap[common.BKInnerObjIDSet]
		if !ok {
			continue
		}

		var layer *metadata.Layer
		for curObj != "" && curObj != common.BKInnerObjIDApp {
			// TODO: add event for layer to use cache, right now data are obtained from db
			objLayer, err := getCache(ctx, cache, db, curObj, parentID, true)
			if err != nil {
				blog.Errorf("identifier: getCache for %s %d error %s", curObj, parentID, err.Error())
				return nil, err
			}

			instID, err := getInt(objLayer.data, common.BKInstIDField)
			if err != nil {
				blog.Errorf("identifier: convert %s bk_inst_id failed, the raw is %+v", curObj, objLayer.data[common.BKInstIDField])
				return nil, err
			}

			layer = &metadata.Layer{
				InstID:   instID,
				InstName: getString(objLayer.data[common.BKInstNameField]),
				ObjID:    curObj,
				Child:    layer,
			}

			curObj = asstMap[curObj]
			parentID, err = getInt(objLayer.data, common.BKParentIDField)
			if err != nil {
				blog.Errorf("identifier: convert set bk_parent_id failed, the raw is %+v", set.data[common.BKParentIDField])
				return nil, err
			}
		}
		hostIdentModule.Layer = layer
	}

	// fill process
	for procIndex := range identifier.Process {
		process := &identifier.Process[procIndex]
		proc, err := getCache(ctx, cache, db, common.BKInnerObjIDProc, process.ProcessID, false)
		if err != nil {
			blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDProc, process.ProcessID, err.Error())
			return nil, err
		}
		process.ProcessName = getString(proc.data[common.BKProcessNameField])
		process.FuncID = getString(proc.data[common.BKFuncIDField])
		process.FuncName = getString(proc.data[common.BKFuncName])
		process.BindIP = getString(proc.data[common.BKBindIP])
		process.Protocol = getString(proc.data[common.BKProtocol])
		process.Port = getString(proc.data[common.BKPort])
		process.StartParamRegex = getString(proc.data[common.BKStartParamRegex])
	}

	return identifier, nil
}
