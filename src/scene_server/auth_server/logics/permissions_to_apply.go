/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package logics

import (
	"fmt"
	"strconv"

	tenantset "configcenter/pkg/types/tenant-set"
	"configcenter/src/ac/iam"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	apigwiam "configcenter/src/thirdparty/apigw/iam"
)

// GetPermissionToApply get the permissions to apply
// 用于鉴权没有通过时，根据鉴权的资源信息生成需要申请的权限信息
func (lgc *Logics) GetPermissionToApply(kit *rest.Kit, rs []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	permission := new(metadata.IamPermission)
	permission.SystemID = iamtypes.SystemIDCMDB
	permission.SystemName = iamtypes.SystemNameCMDB

	data, err := lgc.getPermissionData(kit, rs)
	if err != nil {
		return nil, err
	}

	instIDNameMap, err := lgc.getInstIDNameMap(kit, data.instTypeIDsMap)
	if err != nil {
		blog.Errorf("getInstIDNameMap failed, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}

	for actionID, permissionTypeMap := range data.permissionMap {
		action := metadata.IamAction{
			ID: string(actionID),
			Name: getActionName(actionID, data.actionIDRscTypeMap[actionID],
				data.sysInstModelIDNameMap),
			RelatedResourceTypes: make([]metadata.IamResourceType, 0),
		}
		for rscType := range permissionTypeMap {
			// set instance name
			for idx := range permissionTypeMap[rscType] {
				for idx2 := range permissionTypeMap[rscType][idx] {
					inst := permissionTypeMap[rscType][idx][idx2]
					instID, err := parseInstID(iamtypes.TypeID(inst.Type), inst.ID)
					if err != nil {
						blog.Errorf("parse instance id failed, instID:%#v, err: %s, rid: %s", inst.ID, err, kit.Rid)
						return nil, err
					}
					permissionTypeMap[rscType][idx][idx2].Name = instIDNameMap[instID]
				}
			}

			action.RelatedResourceTypes = append(action.RelatedResourceTypes, metadata.IamResourceType{
				SystemID:   iamtypes.SystemIDCMDB,
				SystemName: iamtypes.SystemNameCMDB,
				Type:       rscType,
				TypeName:   getTypeName(iamtypes.TypeID(rscType), data.sysInstModelIDNameMap),
				Instances:  permissionTypeMap[rscType],
			})
		}
		permission.Actions = append(permission.Actions, action)
	}

	return permission, nil
}

type permissionData struct {
	permissionMap         map[iamtypes.ActionID]map[string][][]metadata.IamResourceInstance
	actionIDRscTypeMap    map[iamtypes.ActionID]iamtypes.TypeID
	sysInstModelIDNameMap map[int64]string
	instTypeIDsMap        map[iamtypes.TypeID][]interface{}
}

func (lgc *Logics) getPermissionData(kit *rest.Kit, rs []meta.ResourceAttribute) (*permissionData, error) {
	sysInstModelIDNameMap := make(map[int64]string)
	actionIDRscTypeMap := make(map[iamtypes.ActionID]iamtypes.TypeID)
	instTypeIDsMap := make(map[iamtypes.TypeID][]interface{})

	// permissionMap maps ResourceActionID and ResourceTypeID to ResourceInstances
	permissionMap := make(map[iamtypes.ActionID]map[string][][]metadata.IamResourceInstance, 0)
	for _, r := range rs {
		opts, err := iam.AdaptAuthOptions(&r)
		if err != nil {
			return nil, err
		}

		for _, opt := range opts {
			if opt.Action == iamtypes.Skip {
				continue
			}

			actionIDRscTypeMap[opt.Action] = opt.Type
			if _, ok := permissionMap[opt.Action]; !ok {
				permissionMap[opt.Action] = make(map[string][][]metadata.IamResourceInstance, 0)
			}
			sysInstModelIDs := getAllSysInstanceModelIDs(opt.Type, opt.Resources)
			if len(sysInstModelIDs) > 0 {
				sysInstModelIDNameMap, err = lgc.GetModelsIDNameMap(kit, sysInstModelIDs)
				if err != nil {
					return nil, err
				}
			}
			// generate iam resource instances by its ancestors and itself
			for _, res := range opt.Resources {
				if len(res.ID) == 0 && res.Attribute == nil {
					permissionMap[opt.Action][string(res.Type)] = nil
					continue
				}
				instance := make([]metadata.IamResourceInstance, 0)
				if res.Attribute != nil {
					ancestors, err := iam.ParseAncestorAttributes(res.Attribute)
					if err != nil {
						return nil, err
					}
					instance = append(instance, ancestors...)
					for _, ancestor := range ancestors {
						ancestorID, err := parseInstID(iamtypes.TypeID(ancestor.Type), ancestor.ID)
						if err != nil {
							blog.Errorf("parse instance id failed, instID: %#v, err: %v, rid: %s", ancestor.ID, err,
								kit.Rid)
							return nil, err
						}
						instTypeIDsMap[iamtypes.TypeID(ancestor.Type)] = append(
							instTypeIDsMap[iamtypes.TypeID(ancestor.Type)], ancestorID)
					}
				}
				if len(res.ID) > 0 {
					instance = append(instance, metadata.IamResourceInstance{
						Type:     string(res.Type),
						TypeName: getTypeName(iamtypes.TypeID(res.Type), sysInstModelIDNameMap),
						ID:       res.ID,
					})
					instID, err := parseInstID(iamtypes.TypeID(res.Type), res.ID)
					if err != nil {
						blog.Errorf("parse instance id failed, instID: %#v, err: %v, rid: %s", res.ID, err, kit.Rid)
						return nil, err
					}
					instTypeIDsMap[iamtypes.TypeID(res.Type)] = append(instTypeIDsMap[iamtypes.TypeID(res.Type)],
						instID)
				}

				permissionMap[opt.Action][string(res.Type)] = append(permissionMap[opt.Action][string(res.Type)],
					instance)
			}
		}
	}
	return &permissionData{permissionMap: permissionMap, actionIDRscTypeMap: actionIDRscTypeMap,
		sysInstModelIDNameMap: sysInstModelIDNameMap, instTypeIDsMap: instTypeIDsMap,
	}, nil
}

// parseInstID keeps numeric ids as int64 and string ids, such as dynamic group ids, as string.
func parseInstID(instType iamtypes.TypeID, id string) (interface{}, error) {
	if isResourceIDStringType(instType) {
		return id, nil
	}
	return strconv.ParseInt(id, 10, 64)
}

func (lgc *Logics) getInstIDNameMap(kit *rest.Kit, instTypeIDsMap map[iamtypes.TypeID][]interface{}) (
	map[interface{}]string, error) {
	instIDNameMap := make(map[interface{}]string)

	for instType, instIDs := range instTypeIDsMap {
		if instType == iamtypes.TenantSet {
			for _, id := range instIDs {
				if id == int64(tenantset.DefaultTenantSetID) {
					instIDNameMap[id] = tenantset.DefaultTenantSetName
					break
				}
			}
			continue
		}
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
			var instID interface{}
			if isResourceIDStringType(instType) {
				instID = util.GetStrByInterface(instance[idField])
			} else {
				instID, err = util.GetInt64ByInterface(instance[idField])
				if err != nil {
					blog.Errorf("parse instance id to int failed, instID: %#v, err: %v, rid: %s", instance[idField],
						err, kit.Rid)
					return nil, err
				}
			}
			instIDNameMap[instID] = util.GetStrByInterface(instance[nameField])
		}

	}

	return instIDNameMap, nil
}

// getAllSysInstanceModelIDs get all system instance objIDs
func getAllSysInstanceModelIDs(rscType iamtypes.TypeID, resource []apigwiam.Resource) []int64 {
	modelIDs := make([]int64, 0)
	if iam.IsIAMSysInstance(rscType) {
		modelID, _ := iam.GetModelIDFromIamSysInstance(rscType)
		modelIDs = append(modelIDs, modelID)
	}

	for _, res := range resource {
		if iam.IsIAMSysInstance(iamtypes.TypeID(res.Type)) {
			modelID, _ := iam.GetModelIDFromIamSysInstance(iamtypes.TypeID(res.Type))
			modelIDs = append(modelIDs, modelID)
		}
	}

	return modelIDs
}

func getActionName(actionID iamtypes.ActionID, rscType iamtypes.TypeID, sysInstModelIDNameMap map[int64]string) string {
	name := ""
	if name = iam.ActionIDNameMap[actionID]; name != "" {
		return name
	}

	if iam.IsIAMSysInstance(rscType) {
		modelID, _ := iam.GetModelIDFromIamSysInstance(rscType)
		actionType := iam.GetActionTypeFromIAMSysInstance(actionID)
		return fmt.Sprintf("%s实例%s", sysInstModelIDNameMap[modelID], iamtypes.ActionTypeIDNameMap[actionType])
	}

	return name
}

func getTypeName(rscType iamtypes.TypeID, sysInstModelIDNameMap map[int64]string) string {
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
