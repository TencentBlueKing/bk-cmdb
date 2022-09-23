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
	"fmt"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

// GetPermissionToApply get the permissions to apply
// 用于鉴权没有通过时，根据鉴权的资源信息生成需要申请的权限信息
func (lgc *Logics) GetPermissionToApply(kit *rest.Kit, rs []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	permission := new(metadata.IamPermission)
	permission.SystemID = iam.SystemIDCMDB
	permission.SystemName = iam.SystemNameCMDB
	instTypeIDsMap := make(map[iam.TypeID][]int64)
	actionIDRscTypeMap := make(map[iam.ActionID]iam.TypeID)
	sysInstModelIDNameMap := make(map[int64]string)

	// permissionMap maps ResourceActionID and ResourceTypeID to ResourceInstances
	permissionMap := make(map[iam.ActionID]map[string][][]metadata.IamResourceInstance, 0)
	for _, r := range rs {
		actionID, err := iam.ConvertResourceAction(r.Type, r.Action, r.BusinessID)
		if err != nil {
			return nil, err
		}

		rscType, err := iam.ConvertResourceType(r.Basic.Type, r.BusinessID)
		if err != nil {
			return nil, err
		}

		resource, err := iam.GenIamResource(actionID, *rscType, &r)
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
				iamPath, ok := res.Attribute[types.IamPathKey].([]string)
				if !ok {
					return nil, fmt.Errorf("iam path(%v) is not string array type", res.Attribute[types.IamPathKey])
				}
				ancestors, err := iam.ParseIamPathToAncestors(iamPath)
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
					instTypeIDsMap[iam.TypeID(ancestor.Type)] = append(instTypeIDsMap[iam.TypeID(ancestor.Type)], ancestorID)
				}
			}

			if len(res.ID) > 0 {
				instance = append(instance, metadata.IamResourceInstance{
					Type:     string(res.Type),
					TypeName: getTypeName(iam.TypeID(res.Type), sysInstModelIDNameMap),
					ID:       res.ID,
				})
				instID, err := strconv.ParseInt(res.ID, 10, 64)
				if err != nil {
					blog.Errorf("parse instance id to int failed, instID:%#v, err: %s, rid: %s", res.ID, err, kit.Rid)
					return nil, err
				}
				instTypeIDsMap[iam.TypeID(res.Type)] = append(instTypeIDsMap[iam.TypeID(res.Type)], instID)
			}

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
				SystemID:   iam.SystemIDCMDB,
				SystemName: iam.SystemNameCMDB,
				Type:       rscType,
				TypeName:   getTypeName(iam.TypeID(rscType), sysInstModelIDNameMap),
				Instances:  permissionTypeMap[rscType],
			})
		}
		permission.Actions = append(permission.Actions, action)
	}

	return permission, nil
}

func (lgc *Logics) getInstIDNameMap(kit *rest.Kit, instTypeIDsMap map[iam.TypeID][]int64) (map[int64]string, error) {
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
func getAllSysInstanceModelIDs(rscType iam.TypeID, resource []types.Resource) []int64 {
	modelIDs := make([]int64, 0)
	if iam.IsIAMSysInstance(rscType) {
		modelID, _ := iam.GetModelIDFromIamSysInstance(rscType)
		modelIDs = append(modelIDs, modelID)
	}

	for _, res := range resource {
		if iam.IsIAMSysInstance(iam.TypeID(res.Type)) {
			modelID, _ := iam.GetModelIDFromIamSysInstance(iam.TypeID(res.Type))
			modelIDs = append(modelIDs, modelID)
		}
	}

	return modelIDs
}

func getActionName(actionID iam.ActionID, rscType iam.TypeID, sysInstModelIDNameMap map[int64]string) string {
	name := ""
	if name = iam.ActionIDNameMap[actionID]; name != "" {
		return name
	}

	if iam.IsIAMSysInstance(rscType) {
		modelID, _ := iam.GetModelIDFromIamSysInstance(rscType)
		actionType := iam.GetActionTypeFromIAMSysInstance(actionID)
		return fmt.Sprintf("%s实例%s", sysInstModelIDNameMap[modelID], iam.ActionTypeIDNameMap[actionType])
	}

	return name
}

func getTypeName(rscType iam.TypeID, sysInstModelIDNameMap map[int64]string) string {
	name := ""
	if name = iam.ResourceTypeIDMap[rscType]; name != "" {
		return name
	}

	if iam.IsIAMSysInstance(rscType) {
		modelID, _ := iam.GetModelIDFromIamSysInstance(rscType)
		return sysInstModelIDNameMap[modelID]
	}

	return name
}
