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
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

// addBKApp add bk app
func addBKApp(ctx context.Context, db dal.RDB, conf *history.Config) error {
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
	appModelData["bk_supplier_account"] = conf.TenantID
	appModelData[common.BKDefaultField] = common.DefaultFlagDefaultValue
	filled := fillEmptyFields(appModelData, AppRow())
	var preData map[string]interface{}
	bizID, preData, err := history.Upsert(ctx, db, common.BKTableNameBaseApp, appModelData, common.BKAppIDField,
		[]string{common.BKAppNameField, "bk_supplier_account"}, append(filled, common.BKAppIDField))
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

	log := AuditLog{
		ID:              int64(id),
		AuditType:       metadata.BusinessType,
		SupplierAccount: conf.TenantID,
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
	inputSetInfo["bk_supplier_account"] = conf.TenantID
	filled = fillEmptyFields(inputSetInfo, SetRow())
	setID, _, err := history.Upsert(ctx, db, common.BKTableNameBaseSet, inputSetInfo, common.BKSetIDField,
		[]string{"bk_supplier_account", common.BKAppIDField, common.BKSetNameField},
		append(filled, common.BKSetIDField))
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
	inputResModuleInfo["bk_supplier_account"] = conf.TenantID
	filled = fillEmptyFields(inputResModuleInfo, ModuleRow())
	_, _, err = history.Upsert(ctx, db, common.BKTableNameBaseModule, inputResModuleInfo, common.BKModuleIDField,
		[]string{"bk_supplier_account", common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField},
		append(filled, common.BKModuleIDField))
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
	inputFaultModuleInfo["bk_supplier_account"] = conf.TenantID
	filled = fillEmptyFields(inputFaultModuleInfo, ModuleRow())
	_, _, err = history.Upsert(ctx, db, common.BKTableNameBaseModule, inputFaultModuleInfo, common.BKModuleIDField,
		[]string{"bk_supplier_account", common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField},
		append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultFaultModule error ", err.Error())
		return err
	}

	return nil
}

// AuditLog audit log struct
type AuditLog struct {
	ID int64 `json:"id" bson:"id"`
	// AuditType is a high level abstract of the resource managed by this cmdb.
	// Each kind of concept, resource must belongs to one of the resource type.
	AuditType metadata.AuditType `json:"audit_type" bson:"audit_type"`
	// the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	// name of the one who triggered this operation.
	User string `json:"user" bson:"user"`
	// the operated resource by the user
	ResourceType metadata.ResourceType `json:"resource_type" bson:"resource_type"`
	// ActionType represent the user's operation type, like CUD etc.
	Action metadata.ActionType `json:"action" bson:"action"`
	// OperateFrom describe which form does this audit come from.
	OperateFrom metadata.OperateFromType `json:"operate_from" bson:"operate_from"`
	// OperationDetail describe the details information by a user.
	// Note: when the ResourceType relevant to Business, then the business id field must
	// be bk_biz_id, otherwise the user can not search this operation log with business id.
	OperationDetail metadata.DetailFactory `json:"operation_detail" bson:"operation_detail"`
	// OperationTime is the time that user do the operation.
	OperationTime metadata.Time `json:"operation_time" bson:"operation_time"`
	// the business id of the resource if it belongs to a business.
	BusinessID int64 `json:"bk_biz_id,omitempty" bson:"bk_biz_id,omitempty"`
	// ResourceID is the id of the resource instance. which is a unique id, dynamic grouping id is string type.
	// for service instance audit log,
	ResourceID interface{} `json:"resource_id" bson:"resource_id"`
	// ResourceName is the name of the resource, such as a switch model has a name "switch"
	ResourceName string `json:"resource_name" bson:"resource_name"`
	// AppCode is the app code of the system where the request comes from
	AppCode string `json:"code,omitempty" bson:"code,omitempty"`
	// RequestID is the request id of the request
	RequestID string `json:"rid,omitempty" bson:"rid,omitempty"`
	// todo ExtendResourceName for the temporary solution of ipv6
	ExtendResourceName string `json:"extend_resource_name" bson:"extend_resource_name"`
}
