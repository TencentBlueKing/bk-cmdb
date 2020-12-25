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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
)

func fillIdentifier(identifier *metadata.HostIdentifier, ctx context.Context, cache redis.Client, clientSet apimachinery.ClientSetInterface, db dal.RDB) (*metadata.HostIdentifier, error) {
	// fill cloudName
	cloud, err := getCache(ctx, cache, clientSet, db, common.BKInnerObjIDPlat, identifier.CloudID)
	if err != nil {
		blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDPlat, identifier.CloudID, err.Error())
		return nil, err
	}
	identifier.CloudName = getString(cloud.data[common.BKCloudNameField])

	customLayers, err := getCustomLayers(ctx, db, identifier.SupplierAccount)
	if err != nil {
		blog.ErrorJSON("identifier: getCustomLayers error %s", err)
		return nil, err
	}
	// fill module
	for _, hostIdentModule := range identifier.HostIdentModule {
		err = fillModule(identifier, hostIdentModule, customLayers, ctx, cache, clientSet, db)
		if err != nil {
			blog.ErrorJSON("identifier: fillModule error %s, hostIdentModule: %s", err, hostIdentModule)
			return nil, err
		}
	}

	// fill process
	for index := range identifier.Process {
		err = fillProcess(&identifier.Process[index], ctx, cache, clientSet, db)
		if err != nil {
			blog.ErrorJSON("identifier: fillProcess error %s, process: %s", err, identifier.Process[index])
			return nil, err
		}
	}

	return identifier, nil
}

func fillProcess(process *metadata.HostIdentProcess, ctx context.Context, cache redis.Client, clientSet apimachinery.ClientSetInterface, db dal.RDB) error {
	proc, err := getCache(ctx, cache, clientSet, db, common.BKInnerObjIDProc, process.ProcessID)
	if err != nil {
		blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDProc, process.ProcessID, err.Error())
		return err
	}

	ip, port, protocol, enable, bindInfoArr := getBindInfo(proc.data[common.BKProcBindInfo])
	process.ProcessName = getString(proc.data[common.BKProcessNameField])
	process.FuncName = getString(proc.data[common.BKFuncName])
	process.BindIP = ip         //getString(proc.data[common.BKBindIP])
	process.Protocol = protocol // getString(proc.data[common.BKProtocol])
	process.Port = port         // getString(proc.data[common.BKPort])
	process.PortEnable = enable
	process.BindInfo = bindInfoArr
	process.StartParamRegex = getString(proc.data[common.BKStartParamRegex])
	return nil
}

func fillModule(identifier *metadata.HostIdentifier, hostIdentModule *metadata.HostIdentModule, customLayers []string,
	ctx context.Context, cache redis.Client, clientSet apimachinery.ClientSetInterface, db dal.RDB) error {

	biz, err := getCache(ctx, cache, clientSet, db, common.BKInnerObjIDApp, hostIdentModule.BizID)
	if err != nil {
		blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDApp, hostIdentModule.BizID, err.Error())
		return err
	}
	hostIdentModule.BizName = getString(biz.data[common.BKAppNameField])
	identifier.SupplierAccount = getString(biz.data[common.BKOwnerIDField])

	set, err := getCache(ctx, cache, clientSet, db, common.BKInnerObjIDSet, hostIdentModule.SetID)
	if err != nil {
		blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDSet, hostIdentModule.SetID, err.Error())
		return err
	}
	hostIdentModule.SetName = getString(set.data[common.BKSetNameField])
	hostIdentModule.SetEnv = getString(set.data[common.BKSetEnvField])
	hostIdentModule.SetStatus = getString(set.data[common.BKSetStatusField])

	module, err := getCache(ctx, cache, clientSet, db, common.BKInnerObjIDModule, hostIdentModule.ModuleID)
	if err != nil {
		blog.Errorf("identifier: getCache for %s %d error %s", common.BKInnerObjIDModule, hostIdentModule.ModuleID, err.Error())
		return err
	}
	hostIdentModule.ModuleName = getString(module.data[common.BKModuleNameField])

	// fill host layer info
	parentID, err := getInt(set.data, common.BKParentIDField)
	if err != nil {
		blog.Errorf("identifier: convert set bk_parent_id failed, the raw is %+v", set.data[common.BKParentIDField])
		return err
	}
	if len(customLayers) == 0 {
		customLayers, err = getCustomLayers(ctx, db, identifier.SupplierAccount)
		if err != nil {
			blog.ErrorJSON("identifier: getCustomLayers error %s", err)
			return err
		}
	}
	var layer *metadata.Layer
	for _, curObj := range customLayers {
		objLayer, err := getCache(ctx, cache, clientSet, db, curObj, parentID)
		if err != nil {
			blog.Errorf("identifier: getCache for %s %d error %s", curObj, parentID, err.Error())
			return err
		}

		instID, err := getInt(objLayer.data, common.BKInstIDField)
		if err != nil {
			blog.Errorf("identifier: convert %s bk_inst_id failed, the raw is %+v", curObj, objLayer.data[common.BKInstIDField])
			return err
		}

		layer = &metadata.Layer{
			InstID:   instID,
			InstName: getString(objLayer.data[common.BKInstNameField]),
			ObjID:    curObj,
			Child:    layer,
		}
		parentID, err = getInt(objLayer.data, common.BKParentIDField)
		if err != nil {
			blog.Errorf("identifier: convert set bk_parent_id failed, the raw is %+v", set.data[common.BKParentIDField])
			return err
		}
	}
	hostIdentModule.Layer = layer
	return nil
}

// get custom layer objects TODO use cache when it supports refreshing
func getCustomLayers(ctx context.Context, db dal.RDB, supplierAccount string) ([]string, error) {
	asstMap := make(map[string]string)
	asstArr := make([]metadata.Association, 0)
	cond := condition.CreateCondition().Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)
	condMap := util.SetQueryOwner(cond.ToMapStr().ToMapInterface(), supplierAccount)
	err := db.Table(common.BKTableNameObjAsst).Find(condMap).All(ctx, &asstArr)
	if err != nil {
		blog.ErrorJSON("findHostLayerInfo query mainline association info error. condition:%s", condMap)
		return nil, err
	}
	for _, asst := range asstArr {
		asstMap[asst.ObjectID] = asst.AsstObjID
	}
	customLayers := make([]string, 0)
	for obj := asstMap[common.BKInnerObjIDSet]; obj != "" && obj != common.BKInnerObjIDApp; obj = asstMap[obj] {
		customLayers = append(customLayers, obj)
	}
	return customLayers, err
}
