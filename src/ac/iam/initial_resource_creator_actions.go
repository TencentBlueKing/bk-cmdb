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

package iam

import (
	iamtypes "configcenter/src/ac/iam/types"
)

// ResourceCreatorRoleMap maps resource type to the role granted after resource creation.
// Reuse manager role when its actions equal the old creator related actions; otherwise grant a dedicated *_owner role.
var ResourceCreatorRoleMap = map[iamtypes.TypeID]iamtypes.RoleID{
	iamtypes.Business:                  iamtypes.BizOwner,
	iamtypes.BizSet:                    iamtypes.BizSetOwner,
	iamtypes.SysCloudArea:              iamtypes.CloudAreaOwner,
	iamtypes.Project:                   iamtypes.ProjectManager,
	iamtypes.SysModelGroup:             iamtypes.ModelGroupManager,
	iamtypes.SysModel:                  iamtypes.SysModelManager,
	iamtypes.SysAssociationType:        iamtypes.AsstTypeManager,
	iamtypes.SysResourcePoolDirectory:  iamtypes.HostPoolDirManager,
	iamtypes.FieldGroupingTemplate:     iamtypes.FieldTplManager,
	iamtypes.BizProcessServiceTemplate: iamtypes.BizSvcTplManager,
	iamtypes.BizSetTemplate:            iamtypes.BizSetTplManager,
	iamtypes.BizCustomQuery:            iamtypes.BizDynQueryManager,
}

// GetResourceCreatorRole get the role granted to the creator of the specified static resource type
func GetResourceCreatorRole(typeID iamtypes.TypeID) (iamtypes.RoleID, bool) {
	roleID, exists := ResourceCreatorRoleMap[typeID]
	return roleID, exists
}
