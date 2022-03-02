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

package v3v0v8

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

//addBKApp add bk app
func addBKApp(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if count, err := db.Table(common.BKTableNameBaseApp).Find(mapstr.MapStr{common.BKAppNameField: common.BKAppName}).Count(ctx); err != nil {
		return err
	} else if count >= 1 {
		return nil
	}

	// add bk app
	appModelData := map[string]interface{}{}
	appModelData[common.BKAppNameField] = common.BKAppName
	appModelData[common.BKMaintainersField] = admin
	appModelData[common.BKTimeZoneField] = "Asia/Shanghai"
	appModelData[common.BKLanguageField] = "1" // "中文"
	appModelData[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal
	appModelData[common.BKOwnerIDField] = conf.OwnerID
	appModelData[common.BKDefaultField] = common.DefaultFlagDefaultValue
	filled := fillEmptyFields(appModelData, AppRow())
	var preData map[string]interface{}
	bizID, preData, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseApp, appModelData, common.BKAppIDField, []string{common.BKAppNameField, common.BKOwnerIDField}, append(filled, common.BKAppIDField))
	if err != nil {
		blog.Error("add addBKApp error ", err.Error())
		return err
	}

	// add audit log
	id, err := db.NextSequence(ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.Errorf("get next audit log id failed, err: %s", err.Error())
		return err
	}

	action := metadata.AuditCreate
	logDetail := &metadata.BasicContent{
		CurData: appModelData,
	}
	if preData != nil {
		action = metadata.AuditUpdate
		logDetail = &metadata.BasicContent{
			PreData:      preData,
			UpdateFields: appModelData,
		}
	}

	log := metadata.AuditLog{
		ID:              int64(id),
		AuditType:       metadata.BusinessType,
		SupplierAccount: conf.OwnerID,
		User:            conf.User,
		ResourceType:    metadata.BusinessRes,
		Action:          action,
		OperateFrom:     metadata.FromCCSystem,
		BusinessID:      int64(bizID),
		ResourceID:      int64(bizID),
		ResourceName:    common.BKAppName,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: logDetail,
			},
			ModelID: common.BKInnerObjIDApp,
		},
		OperationTime: metadata.Now(),
	}

	if err = db.Table(common.BKTableNameAuditLog).Insert(ctx, log); err != nil {
		blog.ErrorJSON("add audit log %s error %s", log, err.Error())
		return err
	}

	// add bk app default set
	inputSetInfo := make(map[string]interface{})
	inputSetInfo[common.BKAppIDField] = bizID
	inputSetInfo[common.BKInstParentStr] = bizID
	inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
	inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
	inputSetInfo[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(inputSetInfo, SetRow())
	setID, _, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseSet, inputSetInfo, common.BKSetIDField, []string{common.BKOwnerIDField, common.BKAppIDField, common.BKSetNameField}, append(filled, common.BKSetIDField))
	if err != nil {
		blog.Error("add defaultSet error ", err.Error())
		return err
	}

	// add bk app default module
	inputResModuleInfo := make(map[string]interface{})
	inputResModuleInfo[common.BKSetIDField] = setID
	inputResModuleInfo[common.BKInstParentStr] = setID
	inputResModuleInfo[common.BKAppIDField] = bizID
	inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
	inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
	inputResModuleInfo[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(inputResModuleInfo, ModuleRow())
	_, _, err = upgrader.Upsert(ctx, db, common.BKTableNameBaseModule, inputResModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultResModule error ", err.Error())
		return err
	}

	inputFaultModuleInfo := make(map[string]interface{})
	inputFaultModuleInfo[common.BKSetIDField] = setID
	inputFaultModuleInfo[common.BKInstParentStr] = setID
	inputFaultModuleInfo[common.BKAppIDField] = bizID
	inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
	inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
	inputFaultModuleInfo[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(inputFaultModuleInfo, ModuleRow())
	_, _, err = upgrader.Upsert(ctx, db, common.BKTableNameBaseModule, inputFaultModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultFaultModule error ", err.Error())
		return err
	}

	return nil
}
