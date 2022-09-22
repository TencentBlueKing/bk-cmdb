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
	types2 "configcenter/cmd/scene_server/auth_server/sdk/types"
	iamtype "configcenter/pkg/ac/iam"
	"configcenter/pkg/ac/meta"
	"fmt"
	"strconv"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/http/rest"
	"configcenter/pkg/metadata"
	"configcenter/pkg/util"
)

// GetPermissionToApply get the permissions to apply
// 用于鉴权没有通过时，根据鉴权的资源信息生成需要申请的权限信息
func (lgc *Logics) GetPermissionToApply(kit *rest.Kit, rs []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	permission := new(metadata.IamPermission)
	permission.SystemID = iamtype.SystemIDCMDB
	permission.SystemName = iamtype.SystemNameCMDB
	instTypeIDsMap := make(map[iamtype.TypeID][]int64)
	actionIDRscTypeMap := make(map[iamtype.ActionID]iamtype.TypeID)
	sysInstModelIDNameMap := make(map[int64]string)

	// permissionMap maps ResourceActionID and ResourceTypeID to ResourceInstances
	permissionMap := make(map[iamtype.ActionID]map[string][][]metadata.IamResourceInstance, 0)
	for _, r := range rs {
		actionID, err := iamtype.ConvertResourceAction(r.Type, r.Action, r.BusinessID)
		if err != nil {
			return nil, err
		}

		rscType, err := iamtype.ConvertResourceType(r.Basic.Type, r.BusinessID)
		if err != nil {
			return nil, err
		}

		resource, err := iamtype.GenIamResource(actionID, *rscType, &r)
		if err != nil {
			return nil, err
		}

		actionIDRscTypeMap[actionID] = *rscType

		if _, ok := permissionMap[actionID]; !ok {
			permissionMap[actionID] = make(map[string][][]metadata.IamResourceInstance, 0)
		}

		sysInstModelIDs := getAllSysInstanceModelIDs(*rscType, resource)
		if len(sysInstModelIDs) > 0 {
			var err error
			if sysInstModelIDNameMap, err = lgc.GetModelsIDNameMap(kit, sysInstModelIDs); err != nil {
				return nil, err
			}
		}

		// generate iam resource instances by its paths and itself
		for _, res := range resource {
			if len(res.ID) == 0 && res.Attribute == nil {
				permissionMap[actionID][string(res.Type)] = nil
				continue
			}

			instance := make([]metadata.IamResourceInstance, 0)
			if res.Attribute != nil {
				iamPath, ok := res.Attribute[types2.IamPathKey].([]string)
				if !ok {
					return nil, fmt.Errorf("iam path(%v) is not string array type", res.Attribute[types2.IamPathKey])
				}
				ancestors, err := iamtype.ParseIamPathToAncestors(iamPath)
				if err != nil {
					return nil, err
				}
				instance = append(instance, ancestors...)
				for _, ancestor := range ancestors {
					ancestorID, err := strconv.ParseInt(ancestor.ID, 10, 64)
					if err != nil {
						blog.Errorf("parse instance id to int failed, instID:%#v, err: %s, rid: %s", ancestor.ID, err, kit.Rid)
						return nil, err
					}
					instTypeIDsMap[iamtype.TypeID(ancestor.Type)] = append(instTypeIDsMap[iamtype.TypeID(ancestor.Type)], ancestorID)
				}
			}
			instance = append(instance, metadata.IamResourceInstance{
				Type:     string(res.Type),
				TypeName: getTypeName(iamtype.TypeID(res.Type), sysInstModelIDNameMap),
				ID:       res.ID,
			})
			instID, err := strconv.ParseInt(res.ID, 10, 64)
			if err != nil {
				blog.Errorf("parse instance id to int failed, instID:%#v, err: %s, rid: %s", res.ID, err, kit.Rid)
				return nil, err
			}
			instTypeIDsMap[iamtype.TypeID(res.Type)] = append(instTypeIDsMap[iamtype.TypeID(res.Type)], instID)
			permissionMap[actionID][string(res.Type)] = append(permissionMap[actionID][string(res.Type)], instance)
		}
	}

	instIDNameMap, err := lgc.getInstIDNameMap(kit, instTypeIDsMap)
	if err != nil {
		blog.Errorf("getInstIDNameMap failed, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}

	for actionID, permissionTypeMap := range permissionMap {
		action := metadata.IamAction{
			ID:                   string(actionID),
			Name:                 getActionName(actionID, actionIDRscTypeMap[actionID], sysInstModelIDNameMap),
			RelatedResourceTypes: make([]metadata.IamResourceType, 0),
		}
		for rscType := range permissionTypeMap {
			// set instance name
			for idx := range permissionTypeMap[rscType] {
				for idx2 := range permissionTypeMap[rscType][idx] {
					instID, err := strconv.ParseInt(permissionTypeMap[rscType][idx][idx2].ID, 10, 64)
					if err != nil {
						blog.Errorf("parse instance id to int failed, instID:%#v, err: %s, rid: %s",
							permissionTypeMap[rscType][idx][idx2].ID, err, kit.Rid)
						return nil, err
					}
					permissionTypeMap[rscType][idx][idx2].Name = instIDNameMap[instID]
				}
			}

			action.RelatedResourceTypes = append(action.RelatedResourceTypes, metadata.IamResourceType{
				SystemID:   iamtype.SystemIDCMDB,
				SystemName: iamtype.SystemNameCMDB,
				Type:       rscType,
				TypeName:   getTypeName(iamtype.TypeID(rscType), sysInstModelIDNameMap),
				Instances:  permissionTypeMap[rscType],
			})
		}
		permission.Actions = append(permission.Actions, action)
	}

	return permission, nil
}

func (lgc *Logics) getInstIDNameMap(kit *rest.Kit, instTypeIDsMap map[iamtype.TypeID][]int64) (map[int64]string, error) {
	instIDNameMap := make(map[int64]string)

	for instType, instIDs := range instTypeIDsMap {
		idField := GetResourceIDField(instType)
		nameField := GetResourceNameField(instType)
		if idField == "" || nameField == "" {
			blog.Errorf("request type %s is invalid, rid: %s", instType, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
		}

		param := metadata.PullResourceParam{
			Condition: map[string]interface{}{
				idField: map[string]interface{}{
					common.BKDBIN: instIDs,
				},
			},
			Fields: []string{idField, nameField},
			Limit:  common.BKNoLimit,
		}
		data, err := lgc.searchAuthResource(kit, param, instType)
		if err != nil {
			blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
			return nil, err
		}

		for _, instance := range data.Info {
			instID, err := util.GetInt64ByInterface(instance[idField])
			if err != nil {
				blog.Errorf("parse instance id to int failed, instID:%#v, err: %s, rid: %s", instance[idField], err, kit.Rid)
				return nil, err
			}
			instIDNameMap[instID] = util.GetStrByInterface(instance[nameField])
		}

	}

	return instIDNameMap, nil
}

// getAllSysInstanceModelIDs get all system instance objIDs
func getAllSysInstanceModelIDs(rscType iamtype.TypeID, resource []types2.Resource) []int64 {
	modelIDs := make([]int64, 0)
	if iamtype.IsIAMSysInstance(rscType) {
		modelID, _ := iamtype.GetModelIDFromIamSysInstance(rscType)
		modelIDs = append(modelIDs, modelID)
	}

	for _, res := range resource {
		if iamtype.IsIAMSysInstance(iamtype.TypeID(res.Type)) {
			modelID, _ := iamtype.GetModelIDFromIamSysInstance(iamtype.TypeID(res.Type))
			modelIDs = append(modelIDs, modelID)
		}
	}

	return modelIDs
}

func getActionName(actionID iamtype.ActionID, rscType iamtype.TypeID, sysInstModelIDNameMap map[int64]string) string {
	name := ""
	if name = iamtype.ActionIDNameMap[actionID]; name != "" {
		return name
	}

	if iamtype.IsIAMSysInstance(rscType) {
		modelID, _ := iamtype.GetModelIDFromIamSysInstance(rscType)
		actionType := iamtype.GetActionTypeFromIAMSysInstance(actionID)
		return fmt.Sprintf("%s实例%s", sysInstModelIDNameMap[modelID], iamtype.ActionTypeIDNameMap[actionType])
	}

	return name
}

func getTypeName(rscType iamtype.TypeID, sysInstModelIDNameMap map[int64]string) string {
	name := ""
	if name = iamtype.ResourceTypeIDMap[rscType]; name != "" {
		return name
	}

	if iamtype.IsIAMSysInstance(rscType) {
		modelID, _ := iamtype.GetModelIDFromIamSysInstance(rscType)
		return sysInstModelIDNameMap[modelID]
	}

	return name
}
