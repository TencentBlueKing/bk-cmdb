/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// getResourceTableName TODO
// get resource database table name
func getResourceTableName(resourceType iam.TypeID) string {
	switch resourceType {
	case iam.Host:
		return common.BKTableNameBaseHost
	case iam.SysModelGroup:
		return common.BKTableNameObjClassification
	case iam.SysModel, iam.SysInstanceModel, iam.SysModelEvent, iam.MainlineModelEvent, iam.InstAsstEvent:
		return common.BKTableNameObjDes
	case iam.SysAssociationType:
		return common.BKTableNameAsstDes
	case iam.SysResourcePoolDirectory, iam.SysHostRscPoolDirectory:
		return common.BKTableNameBaseModule
	case iam.SysCloudArea:
		return common.BKTableNameBasePlat
	case iam.SysCloudAccount:
		return common.BKTableNameCloudAccount
	case iam.SysCloudResourceTask:
		return common.BKTableNameCloudSyncTask
	case iam.Business, iam.BusinessForHostTrans:
		return common.BKTableNameBaseApp
	case iam.BizSet:
		return common.BKTableNameBaseBizSet
	case iam.BizCustomQuery:
		return common.BKTableNameDynamicGroup
	case iam.BizProcessServiceTemplate:
		return common.BKTableNameServiceTemplate
	case iam.BizProcessServiceCategory:
		return common.BKTableNameServiceCategory
	case iam.BizProcessServiceInstance:
		return common.BKTableNameServiceInstance
	case iam.BizSetTemplate:
		return common.BKTableNameSetTemplate
	// case iam.Set:
	//	return common.BKTableNameBaseSet
	// case iam.Module:
	//	return common.BKTableNameBaseModule
	default:
		return ""
	}
}

func isResourceIDStringType(resourceType iam.TypeID) bool {
	switch resourceType {
	case iam.BizCustomQuery:
		return true
	}
	return false
}

// getInstanceResourceObjID TODO
// get model instance resource's model id
func getInstanceResourceObjID(resourceType iam.TypeID) string {
	switch resourceType {
	case iam.Host:
		return common.BKInnerObjIDHost
	case iam.SysCloudArea:
		return common.BKInnerObjIDPlat
	case iam.Business, iam.BusinessForHostTrans:
		return common.BKInnerObjIDApp
	case iam.BizSet:
		return common.BKInnerObjIDBizSet
	// case iam.Set:
	//	return common.BKInnerObjIDSet
	// case iam.Module:
	//	return common.BKInnerObjIDModule
	default:
		return ""
	}
}

var resourcePoolBizID int64

// GetResourcePoolBizID search bizID of resource pool
func (lgc *Logics) GetResourcePoolBizID(kit *rest.Kit) (int64, error) {
	if resourcePoolBizID != 0 {
		return resourcePoolBizID, nil
	}

	input := &metadata.QueryCondition{
		Condition: map[string]interface{}{common.BKDefaultField: common.DefaultAppFlag},
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKAppIDField},
		Fields:    []string{common.BKAppIDField, common.BkSupplierAccount},
	}

	bizResp, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp,
		input)
	if err != nil {
		blog.Errorf("find resource pool biz failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return 0, err
	}

	if len(bizResp.Info) <= 0 {
		blog.Errorf("find no resource pool biz, rid: %s", kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	for _, biz := range bizResp.Info {
		if util.GetStrByInterface(biz[common.BkSupplierAccount]) == kit.SupplierAccount {
			resourcePoolBizID, err = util.GetInt64ByInterface(biz[common.BKAppIDField])
			if nil != err {
				blog.ErrorJSON("find resource pool biz failed, parse biz id failed, biz: %s, err: %s, rid: %s",
					bizResp.Info[0][common.BKAppIDField], err.Error(), kit.Rid)
				return 0, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp,
					common.BKAppIDField, "int", err.Error())
			}
			return resourcePoolBizID, nil
		}
	}

	return 0, kit.CCError.CCError(common.CCErrCommNotFound)
}

// getCloudNameMapByIDs get cloud area ID to name map by ID to generate host display name
func (lgc *Logics) getCloudNameMapByIDs(kit *rest.Kit, cloudIDs []int64) (map[int64]string, error) {
	cloudParam := metadata.QueryCondition{
		Fields:    []string{common.BKCloudIDField, common.BKCloudNameField},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{common.BKCloudIDField: map[string]interface{}{common.BKDBIN: cloudIDs}},
	}
	cloudRsp, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDPlat,
		&cloudParam)
	if err != nil {
		blog.Errorf("get cloud areas failed, err: %v,cloudIDs: %+v", err, cloudIDs)
		return nil, err
	}

	cloudMap := make(map[int64]string)
	for _, cloud := range cloudRsp.Info {
		cloudID, err := util.GetInt64ByInterface(cloud[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("parse cloud area id failed, err: %v,cloud: %+v", err, cloud)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
		}

		cloudMap[cloudID] = util.GetStrByInterface(cloud[common.BKCloudNameField])
	}

	return cloudMap, nil
}

func getHostDisplayName(innerIP string, cloudName string) string {
	return innerIP + "(" + cloudName + ")"
}

// GetModelsIDNameMap get a map, key is id, value is bk_obj_name
func (lgc *Logics) GetModelsIDNameMap(kit *rest.Kit, modelIDs []int64) (map[int64]string, error) {
	cond := &metadata.QueryCondition{
		Fields: []string{common.BKObjNameField, common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: modelIDs,
			},
		},
	}

	resp, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.ErrorJSON("get models failed, err:%s, cond:%s, rid:%s", err, cond, kit.Rid)
		return nil, fmt.Errorf("get models failed, err: %+v", err)
	}
	if len(resp.Info) == 0 {
		blog.ErrorJSON("get models failed, no model was found, cond:%s, rid:%s", cond, kit.Rid)
		return nil, fmt.Errorf("get models failed, no model was found")
	}

	objIDNameMap := make(map[int64]string)
	for _, item := range resp.Info {
		objIDNameMap[item.ID] = item.ObjectName
	}

	return objIDNameMap, nil
}

// modelIDObjIDMap is a concurrent safe type
type modelIDObjIDMap struct {
	sync.RWMutex
	data map[int64]string
}

// modelObjIDMap is map whose key is modelID , value is objID
// eg : {7:"bk_switch"}
var modelObjIDMap = &modelIDObjIDMap{
	data: make(map[int64]string),
}

func (m *modelIDObjIDMap) get(objID int64) string {
	m.RLock()
	defer m.RUnlock()
	return m.data[objID]
}

func (m *modelIDObjIDMap) set(modelID int64, objID string) {
	m.Lock()
	defer m.Unlock()
	m.data[modelID] = objID
}

// GetObjIDFromResourceType get objID from resourceType
func (lgc *Logics) GetObjIDFromResourceType(ctx context.Context, header http.Header, resourceType iam.TypeID) (string,
	error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	modelID, err := iam.GetModelIDFromIamSysInstance(resourceType)
	if err != nil {
		return "", err
	}

	if objID := modelObjIDMap.get(modelID); objID != "" {
		return objID, nil
	}

	cond := &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: modelID,
		},
	}

	resp, err := lgc.CoreAPI.CoreService().Model().ReadModel(ctx, header, cond)
	if err != nil {
		blog.ErrorJSON("get model failed, err:%s, cond:%s, rid:%s", err, cond, rid)
		return "", fmt.Errorf("get model failed, err: %+v", err)
	}
	if len(resp.Info) == 0 {
		blog.ErrorJSON("get model failed, no model was found, cond:%s, rid:%s", cond, rid)
		return "", fmt.Errorf("get model failed, no model was found")
	}

	for _, item := range resp.Info {
		modelObjIDMap.set(item.ID, item.ObjectID)
		return item.ObjectID, nil
	}

	return modelObjIDMap.get(modelID), nil
}
