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

// GenerateResourceCreatorActions generate all the resource creator actions registered to IAM.
func GenerateResourceCreatorActions() ResourceCreatorActions {
	return ResourceCreatorActions{
		Config: []ResourceCreatorAction{
			{
				ResourceID: SysResourcePoolDirectory,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditResourcePoolDirectory,
						IsRequired: false,
					},
					{
						ID:         DeleteResourcePoolDirectory,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: Business,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditBusiness,
						IsRequired: false,
					},
					{
						ID:         ArchiveBusiness,
						IsRequired: false,
					},
					{
						ID:         FindBusiness,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysCloudAccount,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditCloudAccount,
						IsRequired: false,
					},
					{
						ID:         DeleteCloudAccount,
						IsRequired: false,
					},
					{
						ID:         FindCloudAccount,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysCloudResourceTask,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditCloudResourceTask,
						IsRequired: false,
					},
					{
						ID:         DeleteCloudResourceTask,
						IsRequired: false,
					},
					{
						ID:         FindCloudResourceTask,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysCloudArea,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditCloudArea,
						IsRequired: false,
					},
					{
						ID:         DeleteCloudArea,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysEventPushing,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditEventPushing,
						IsRequired: false,
					},
					{
						ID:         DeleteEventPushing,
						IsRequired: false,
					},
					{
						ID:         FindEventPushing,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysModelGroup,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditModelGroup,
						IsRequired: false,
					},
					{
						ID:         DeleteModelGroup,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysModel,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditSysModel,
						IsRequired: false,
					},
					{
						ID:         DeleteSysModel,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: SysAssociationType,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditAssociationType,
						IsRequired: false,
					},
					{
						ID:         DeleteAssociationType,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: BizProcessServiceTemplate,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditBusinessServiceTemplate,
						IsRequired: false,
					},
					{
						ID:         DeleteBusinessServiceTemplate,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: BizSetTemplate,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditBusinessSetTemplate,
						IsRequired: false,
					},
					{
						ID:         DeleteBusinessSetTemplate,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
			{
				ResourceID: BizCustomQuery,
				Actions: []CreatorRelatedAction{
					{
						ID:         EditBusinessCustomQuery,
						IsRequired: false,
					},
					{
						ID:         DeleteBusinessCustomQuery,
						IsRequired: false,
					},
				},
				SubResourceTypes: nil,
			},
		},
	}
}
