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
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// get resource database table name
func getResourceTableName(resourceType iam.TypeID) string {
	switch resourceType {
	case iam.Host:
		return common.BKTableNameBaseHost
	case iam.SysEventPushing:
		return common.BKTableNameSubscription
	case iam.SysModelGroup:
		return common.BKTableNameObjClassification
	case iam.SysModel, iam.SysInstanceModel:
		return common.BKTableNameObjDes
	case iam.SysInstance:
		return common.BKTableNameBaseInst
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
	case iam.Business:
		return common.BKTableNameBaseApp
	case iam.BizCustomQuery:
		return common.BKTableNameUserAPI
	case iam.BizProcessServiceTemplate:
		return common.BKTableNameServiceTemplate
	case iam.BizProcessServiceCategory:
		return common.BKTableNameServiceCategory
	case iam.BizProcessServiceInstance:
		return common.BKTableNameServiceInstance
	case iam.BizSetTemplate:
		return common.BKTableNameSetTemplate
	//case iam.Set:
	//	return common.BKTableNameBaseSet
	//case iam.Module:
	//	return common.BKTableNameBaseModule
	default:
		return ""
	}
}

// get model instance resource's model id
func GetInstanceResourceObjID(resourceType iam.TypeID) string {
	switch resourceType {
	case iam.Host:
		return common.BKInnerObjIDHost
	case iam.SysCloudArea:
		return common.BKInnerObjIDPlat
	case iam.Business:
		return common.BKInnerObjIDApp
	//case iam.Set:
	//	return common.BKInnerObjIDSet
	//case iam.Module:
	//	return common.BKInnerObjIDModule
	default:
		return ""
	}
}

// generate condition for resource type that have special constraints
func (lgc *Logics) generateSpecialCondition(kit *rest.Kit, resourceType iam.TypeID, condition map[string]interface{}) (map[string]interface{}, error) {
	if condition == nil {
		condition = make(map[string]interface{})
	}

	// not include default business
	if resourceType == iam.Business {
		condition[common.BKDefaultField] = map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		}
		return condition, nil
	}

	// model not include mainline model
	// process and cloud area are temporarily excluded TODO remove this restriction when they are available for user
	if resourceType == iam.SysModel || resourceType == iam.SysInstanceModel {
		excludedObjIDs := []string{common.BKInnerObjIDProc, common.BKInnerObjIDPlat, common.BKInnerObjIDApp,
			common.BKInnerObjIDSet, common.BKInnerObjIDModule}
		// cond := &metadata.QueryCondition{
		// 	Condition: map[string]interface{}{
		// 		common.AssociationKindIDField: common.AssociationKindMainline,
		// 	},
		// }
		// asst, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, cond)
		// if err != nil {
		// 	return nil, err
		// }
		//
		// if !asst.Result {
		// 	return nil, asst.CCError()
		// }
		//
		// for _, mainline := range asst.Data.Info {
		// 	switch mainline.AsstObjID {
		// 	// only exclude these in mainline association tree.
		// 	case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule:
		// 		continue
		// 	default:
		// 	}
		// 	excludedObjIDs = append(excludedObjIDs, mainline.AsstObjID)
		// }

		if resourceType == iam.SysInstanceModel {
			excludedObjIDs = append(excludedObjIDs, common.BKInnerObjIDHost)
		}

		condition[common.BKObjIDField] = map[string]interface{}{
			common.BKDBNIN: excludedObjIDs,
		}
		return condition, nil
	}

	if resourceType != iam.SysResourcePoolDirectory && resourceType != iam.SysHostRscPoolDirectory {
		return condition, nil
	}

	// TODO use cache
	// get resource pool biz id from cache TODO confirm if it need a separate key
	var defaultBizIDVal interface{}
	//businesses, err := cache.GetCacheItemsByKeyRegex(common.BKCacheKeyV3Prefix+common.BKInnerObjIDApp+":id:*", lgc.cache)
	//if err != nil {
	// get biz from db if get it from cache encounters error
	//blog.Errorf("get business from cache failed, try to get from db, error: %s", err.Error())
	input := &metadata.QueryCondition{
		Condition: map[string]interface{}{common.BKDefaultField: common.DefaultAppFlag},
		Page:      metadata.BasePage{Start: 0, Limit: 1, Sort: common.BKAppIDField},
		Fields:    []string{common.BKAppIDField},
	}
	bizResp, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, input)
	if err != nil {
		blog.Errorf("find resource pool biz failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}
	if !bizResp.Result {
		blog.Errorf("find resource pool biz failed, err code: %d, err msg: %s, rid: %s", bizResp.Code, bizResp.ErrMsg, kit.Rid)
		return nil, bizResp.Error()
	}
	if len(bizResp.Data.Info) <= 0 {
		blog.Errorf("find no resource pool biz, rid: %s", kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	defaultBizIDVal = bizResp.Data.Info[0][common.BKAppIDField]
	defaultBizID, err := util.GetInt64ByInterface(defaultBizIDVal)
	if nil != err {
		blog.ErrorJSON("find resource pool biz failed, parse biz id failed, biz: %s, err: %s, rid: %s", defaultBizIDVal, err.Error(), kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}
	condition[common.BKAppIDField] = defaultBizID
	return condition, nil
}
