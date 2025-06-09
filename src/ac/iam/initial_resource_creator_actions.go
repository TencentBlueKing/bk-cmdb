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

// GenerateResourceCreatorActions generate all the resource creator actions registered to IAM.
func GenerateResourceCreatorActions() iam.ResourceCreatorActions {
	return iam.ResourceCreatorActions{
		Config: []iam.ResourceCreatorAction{
			{
				ResourceID: types.SysResourcePoolDirectory,
				Actions: []iam.CreatorRelatedAction{{ID: types.EditResourcePoolDirectory, IsRequired: false},
					{ID: types.DeleteResourcePoolDirectory, IsRequired: false}}, SubResourceTypes: nil},
			{
				ResourceID: types.Business,
				Actions: []iam.CreatorRelatedAction{{ID: types.EditBusiness, IsRequired: false},
					{ID: types.ArchiveBusiness, IsRequired: false}, {ID: types.FindBusiness, IsRequired: false}},
				SubResourceTypes: nil,
			},
			{
				ResourceID: types.SysCloudArea,
				Actions: []iam.CreatorRelatedAction{{ID: types.EditCloudArea, IsRequired: false},
					{ID: types.DeleteCloudArea, IsRequired: false}}, SubResourceTypes: nil},
			{
				ResourceID: types.SysModelGroup,
				Actions: []iam.CreatorRelatedAction{{ID: types.EditModelGroup, IsRequired: false},
					{ID: types.DeleteModelGroup, IsRequired: false}},
				SubResourceTypes: nil,
			},
			{
				ResourceID: types.SysModel,
				Actions: []iam.CreatorRelatedAction{{ID: types.ViewSysModel, IsRequired: false},
					{ID: types.EditSysModel, IsRequired: false},
					{ID: types.DeleteSysModel, IsRequired: false}}, SubResourceTypes: nil},
			{
				ResourceID: types.SysAssociationType,
				Actions: []iam.CreatorRelatedAction{{ID: types.EditAssociationType, IsRequired: false},
					{ID: types.DeleteAssociationType, IsRequired: false}},
				SubResourceTypes: nil,
			},
			{
				ResourceID: types.BizProcessServiceTemplate,
				Actions: []iam.CreatorRelatedAction{
					{ID: types.EditBusinessServiceTemplate, IsRequired: false},
					{ID: types.DeleteBusinessServiceTemplate, IsRequired: false},
				}, SubResourceTypes: nil,
			},
			{
				ResourceID: types.BizSetTemplate,
				Actions: []iam.CreatorRelatedAction{
					{ID: types.EditBusinessSetTemplate, IsRequired: false},
					{ID: types.DeleteBusinessSetTemplate, IsRequired: false},
				}, SubResourceTypes: nil,
			},
			{
				ResourceID: types.BizCustomQuery,
				Actions: []iam.CreatorRelatedAction{
					{ID: types.EditBusinessCustomQuery, IsRequired: false},
					{ID: types.DeleteBusinessCustomQuery, IsRequired: false},
				}, SubResourceTypes: nil,
			},
			{
				ResourceID: types.FieldGroupingTemplate,
				Actions: []iam.CreatorRelatedAction{
					{ID: types.EditFieldGroupingTemplate, IsRequired: false},
					{ID: types.DeleteFieldGroupingTemplate, IsRequired: false},
					{ID: types.ViewFieldGroupingTemplate, IsRequired: false},
				}, SubResourceTypes: nil,
			},
			{
				ResourceID: types.BizSet, Actions: []iam.CreatorRelatedAction{{ID: types.EditBizSet, IsRequired: false},
					{ID: types.DeleteBizSet, IsRequired: false}, {ID: types.ViewBizSet, IsRequired: false}},
				SubResourceTypes: nil,
			},
			{
				ResourceID: types.Project,
				Actions: []iam.CreatorRelatedAction{{ID: types.EditProject, IsRequired: false},
					{ID: types.DeleteProject, IsRequired: false}, {ID: types.ViewProject, IsRequired: false}},
				SubResourceTypes: nil,
			},
		},
	}
}
