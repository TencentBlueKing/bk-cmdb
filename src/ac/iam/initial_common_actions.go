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

package iam

import (
	"configcenter/src/ac/iam/types"
	"configcenter/src/thirdparty/apigw/iam"
)

// GenerateCommonActions generate all the common actions registered to IAM.
func GenerateCommonActions() []iam.CommonAction {
	return []iam.CommonAction{
		{
			Name:        "业务运维",
			EnglishName: "Business Maintainer",
			Actions: []iam.ActionWithID{{ID: types.ViewBusinessResource}, {ID: types.EditBusinessHost},
				{ID: types.BusinessHostTransferToResourcePool}, {ID: types.CreateBusinessTopology},
				{ID: types.EditBusinessTopology}, {ID: types.DeleteBusinessTopology},
				{ID: types.CreateBusinessServiceInstance}, {ID: types.EditBusinessServiceInstance},
				{ID: types.DeleteBusinessServiceInstance}, {ID: types.CreateBusinessServiceTemplate},
				{ID: types.EditBusinessServiceTemplate}, {ID: types.DeleteBusinessServiceTemplate},
				{ID: types.CreateBusinessSetTemplate}, {ID: types.EditBusinessSetTemplate},
				{ID: types.DeleteBusinessSetTemplate}, {ID: types.CreateBusinessServiceCategory},
				{ID: types.EditBusinessServiceCategory}, {ID: types.DeleteBusinessServiceCategory},
				{ID: types.CreateBusinessCustomQuery}, {ID: types.EditBusinessCustomQuery},
				{ID: types.DeleteBusinessCustomQuery}, {ID: types.EditBusinessCustomField},
				{ID: types.EditBusinessHostApply}, {ID: types.FindBusiness}},
		},
		{
			Name:        "业务只读",
			EnglishName: "Business Visitor",
			Actions:     []iam.ActionWithID{{ID: types.ViewBusinessResource}, {ID: types.FindBusiness}},
		},
		{
			Name:        "业务集运维",
			EnglishName: "Biz-set Maintainer",
			Actions: []iam.ActionWithID{{ID: types.AccessBizSet}, {ID: types.DeleteBizSet},
				{ID: types.ViewBizSet}},
		},
		{
			Name:        "业务集只读",
			EnglishName: "Biz-set Visitor",
			Actions:     []iam.ActionWithID{{ID: types.AccessBizSet}, {ID: types.ViewBizSet}},
		},
		{
			Name:        "主机资源管理员",
			EnglishName: "Host Maintainer",
			Actions: []iam.ActionWithID{{ID: types.ViewResourcePoolHost}, {ID: types.CreateResourcePoolHost},
				{ID: types.EditResourcePoolHost}, {ID: types.DeleteResourcePoolHost},
				{ID: types.ResourcePoolHostTransferToBusiness}, {ID: types.ResourcePoolHostTransferToDirectory},
				{ID: types.ManageHostAgentID}, {ID: types.CreateResourcePoolDirectory},
				{ID: types.EditResourcePoolDirectory}, {ID: types.DeleteResourcePoolDirectory}},
		},
		{
			Name:        "开发者",
			EnglishName: "Developer",
			Actions: []iam.ActionWithID{{ID: types.WatchHostEvent}, {ID: types.WatchHostRelationEvent},
				{ID: types.WatchBizEvent}, {ID: types.WatchSetEvent}, {ID: types.WatchModuleEvent},
				{ID: types.WatchProcessEvent}, {ID: types.WatchCommonInstanceEvent}},
		},
		{
			Name:        "模型关系维护人",
			EnglishName: "Model Maintainer",
			Actions: []iam.ActionWithID{{ID: types.CreateModelGroup}, {ID: types.EditModelGroup},
				{ID: types.DeleteModelGroup}, {ID: types.EditBusinessLayer}, {ID: types.EditModelTopologyView},
				{ID: types.CreateSysModel}, {ID: types.EditSysModel}, {ID: types.DeleteSysModel},
				{ID: types.CreateAssociationType}, {ID: types.EditAssociationType}, {ID: types.DeleteAssociationType}},
		},
		{
			Name:        "审计员",
			EnglishName: "Auditor",
			Actions:     []iam.ActionWithID{{ID: types.FindAuditLog}},
		},
	}
}
