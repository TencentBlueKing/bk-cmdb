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
func getResourceTableName(resourceType iam.ResourceTypeID) string {
	switch resourceType {
	case iam.SysHostInstance, iam.BizHostInstance:
		return common.BKTableNameBaseHost
	case iam.SysEventPushing:
		return common.BKTableNameSubscription
	case iam.SysModelGroup:
		return common.BKTableNameObjClassification
	case iam.SysModel:
		return common.BKTableNameObjDes
	case iam.SysInstance:
		return common.BKTableNameBaseInst
	case iam.SysAssociationType:
		return common.BKTableNameAsstDes
	case iam.SysResourcePoolDirectory:
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
	default:
		return ""
	}
}

func isResourceIDStringType(resourceType iam.ResourceTypeID) bool {
	switch resourceType {
	case iam.SysModelGroup, iam.SysModel, iam.SysAssociationType:
		return true
	}
	return false
}

// get model instance resource's model id
func GetInstanceResourceObjID(resourceType iam.ResourceTypeID) string {
	switch resourceType {
	case iam.SysHostInstance, iam.BizHostInstance:
		return common.BKInnerObjIDHost
	case iam.SysCloudArea:
		return common.BKInnerObjIDPlat
	case iam.Business:
		return common.BKInnerObjIDApp
	default:
		return ""
	}
}

// generate condition for resource type that have special constraints
func (lgc *Logics) generateSpecialCondition(kit *rest.Kit, resourceType iam.ResourceTypeID, condition map[string]interface{}) (map[string]interface{}, error) {
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
	if resourceType == iam.SysModel {
		excludedObjIDs := []string{common.BKInnerObjIDProc, common.BKInnerObjIDPlat}
		cond := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				common.AssociationKindIDField: common.AssociationKindMainline,
			},
		}
		asst, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, cond)
		if err != nil {
			return nil, err
		}
		if !asst.Result {
			return nil, asst.CCError()
		}
		for _, mainline := range asst.Data.Info {
			excludedObjIDs = append(excludedObjIDs, mainline.AsstObjID)
		}
		condition[common.BKObjIDField] = map[string]interface{}{
			common.BKDBNIN: excludedObjIDs,
		}
		return condition, nil
	}

	if resourceType != iam.BizHostInstance && resourceType != iam.SysHostInstance && resourceType != iam.SysResourcePoolDirectory {
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
	//} else {
	//	for _, biz := range businesses {
	//		if biz[common.BKDefaultField] == strconv.Itoa(common.DefaultAppFlag) {
	//			defaultBizIDVal = biz[common.BKAppIDField]
	//			break
	//		}
	//	}
	//}
	defaultBizID, err := util.GetInt64ByInterface(defaultBizIDVal)
	if nil != err {
		blog.ErrorJSON("find resource pool biz failed, parse biz id failed, biz: %s, err: %s, rid: %s", defaultBizIDVal, err.Error(), kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}

	if resourceType == iam.SysResourcePoolDirectory {
		condition[common.BKAppIDField] = defaultBizID
		return condition, nil
	}

	// TODO use cache
	// get resource pool host IDs from cache
	hostIDs := make([]int64, 0)
	//relations, err := cache.GetCacheItemsByKeyRegex(common.BKCacheKeyV3Prefix+"host_module:*", lgc.cache)
	//if err == nil {
	//	for _, relation := range relations {
	//		if relation[common.BKAppIDField] == strconv.FormatInt(defaultBizID, 10) {
	//			hostID, err := strconv.ParseInt(relation[common.BKHostIDField], 10, 64)
	//			if err != nil {
	//				blog.ErrorJSON("parse cached relation host id %s failed, err: %s, rid: %s", relation[common.BKHostIDField], err.Error(), kit.Rid)
	//				break
	//			}
	//			hostIDs = append(hostIDs, hostID)
	//		}
	//	}
	//} else {
	// get host module relation from db if get it from cache encounters error
	//blog.Errorf("get host module relation from cache failed, try to get from db, error: %s", err.Error())
	hostRsp, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, &metadata.HostModuleRelationRequest{ApplicationID: defaultBizID})
	if err != nil {
		blog.Errorf("get resource pool host ids failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}
	if !hostRsp.Result {
		blog.Errorf("get resource pool host ids failed, err code: %d, err msg: %s, rid: %s", hostRsp.Code, hostRsp.ErrMsg, kit.Rid)
		return nil, hostRsp.Error()
	}
	for _, relation := range hostRsp.Data.Info {
		hostIDs = append(hostIDs, relation.HostID)
	}
	//}

	var hostCond map[string]interface{}
	switch resourceType {
	case iam.BizHostInstance:
		hostCond = map[string]interface{}{
			common.BKDBNIN: hostIDs,
		}
	case iam.SysHostInstance:
		hostCond = map[string]interface{}{
			common.BKDBIN: hostIDs,
		}
	}
	if condition[common.BKHostIDField] == nil {
		condition[common.BKHostIDField] = hostCond
		return condition, nil
	}
	if andCond := condition[common.BKDBAND]; andCond != nil {
		switch andCondArr := andCond.(type) {
		case []interface{}:
			andCondArr = append(andCondArr, map[string]interface{}{common.BKHostIDField: hostCond})
			condition[common.BKDBAND] = andCondArr
		case []map[string]interface{}:
			andCondArr = append(andCondArr, map[string]interface{}{common.BKHostIDField: hostCond})
			condition[common.BKDBAND] = andCondArr
		default:
			blog.Errorf("and condition is not slice type, rid: %s", kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "condition.$and")
		}
		return condition, nil
	}
	condition[common.BKDBAND] = []map[string]interface{}{{common.BKHostIDField: hostCond}}
	return condition, nil
}

// get model instance resource's model id
//func GetResourceTopoPath(resourceType iam.ResourceTypeID) []types. {
//	return
//}
