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

var (
	businessChain = ResourceChain{
		SystemID: SystemIDCMDB,
		ID:       Business,
	}
)

// GenerateInstanceSelections generate all the instance selections registered to IAM.
func GenerateInstanceSelections() []InstanceSelection {
	return []InstanceSelection{
		{
			ID:                BusinessSelection,
			Name:              "业务",
			NameEn:            "business",
			ResourceTypeChain: []ResourceChain{businessChain},
		},
		{
			ID:     SysResourcePoolDirectorySelection,
			Name:   "资源池目录",
			NameEn: "Resource Pool Directory",
			ResourceTypeChain: []ResourceChain{{
				SystemID: SystemIDCMDB,
				ID:       SysResourcePoolDirectory,
			}},
		},
		{
			ID:     BizHostInstanceSelection,
			Name:   "业务主机",
			NameEn: "Business's Hosts",
			ResourceTypeChain: []ResourceChain{
				// select the business at first.
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       Set,
				},
				{
					SystemID: SystemIDCMDB,
					ID:       Module,
				},
				// then select the host instances.
				{
					SystemID: SystemIDCMDB,
					ID:       BizHostInstance,
				},
			},
		},
		{
			ID:     BizCustomQuerySelection,
			Name:   "动态分组",
			NameEn: "Custom Query",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizCustomQuery,
				},
			},
		},
		{
			ID:     BizProcessServiceCategorySelection,
			Name:   "服务分类",
			NameEn: "Service Category",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizProcessServiceCategory,
				},
			},
		},
		{
			ID:     BizProcessServiceInstanceSelection,
			Name:   "服务实例",
			NameEn: "Service Instance",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizProcessServiceInstance,
				},
			},
		},
		{
			ID:     BizProcessServiceTemplateSelection,
			Name:   "服务模板",
			NameEn: "Service Template",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizProcessServiceTemplate,
				},
			},
		},
		{
			ID:     BizSetTemplateSelection,
			Name:   "集群模板",
			NameEn: "Set Template",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizSetTemplate,
				},
			},
		},
		{
			ID:     SysHostInstanceSelection,
			Name:   "资源池主机",
			NameEn: "Resource Pool Host",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysResourcePoolDirectory,
				},
				{
					SystemID: SystemIDCMDB,
					ID:       SysHostInstance,
				},
			},
		},
		{
			ID:     SysCloudAreaSelection,
			Name:   "云区域",
			NameEn: "Cloud Area",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysCloudArea,
				},
			},
		},
		{
			ID:     SysInstanceSelection,
			Name:   "实例",
			NameEn: "Instance",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysModel,
				},
				{
					SystemID: SystemIDCMDB,
					ID:       SysInstance,
				},
			},
		},
		{
			ID:     SysEventPushingSelection,
			Name:   "事件推送",
			NameEn: "Event Pushing",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysEventPushing,
				},
			},
		},
		{
			ID:     SysCloudAccountSelection,
			Name:   "云账户",
			NameEn: "Cloud Account",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysCloudAccount,
				},
			},
		},
		{
			ID:     SysCloudResourceTaskSelection,
			Name:   "云资源任务",
			NameEn: "Cloud Resource Task",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysCloudResourceTask,
				},
			},
		},
		{
			ID:     SysModelSelection,
			Name:   "模型",
			NameEn: "Model",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysModel,
				},
			},
		},
		{
			ID:     SysAssociationTypeSelection,
			Name:   "关联类型",
			NameEn: "Association Type",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysAssociationType,
				},
			},
		},
		{
			ID:     SysModelGroupSelection,
			Name:   "模型分组",
			NameEn: "Model Group",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysModelGroup,
				},
			},
		},
	}
}
