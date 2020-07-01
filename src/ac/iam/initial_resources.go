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
	businessParent = Parent{
		SystemID:   SystemIDCMDB,
		ResourceID: Business,
	}
)

// GenerateResourceTypes generate all the resource types registered to IAM.
func GenerateResourceTypes() []ResourceType {
	resourceTypeList := make([]ResourceType, 0)

	// add public resources
	resourceTypeList = append(resourceTypeList, genPublicResources()...)

	// add business resources
	resourceTypeList = append(resourceTypeList, genBusinessResources()...)

	return resourceTypeList
}

func genBusinessResources() []ResourceType {
	return []ResourceType{
		{
			ID:            Host,
			Name:          "主机",
			NameEn:        "Host",
			Description:   "主机",
			DescriptionEn: "hosts under a business or in resource pool",
			Parents: []Parent{{
				SystemID: SystemIDCMDB,
				//ResourceID: Module,
				ResourceID: Business,
			}, {
				SystemID:   SystemIDCMDB,
				ResourceID: SysResourcePoolDirectory,
			}},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/instance/resource",
			},
			Version: 1,
		},
		{
			ID:            BizHostApply,
			Name:          "主机自动应用",
			NameEn:        "Host Apply",
			Description:   "自动应用业务主机的属性信息",
			DescriptionEn: "apply business host attribute automatically",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            BizCustomQuery,
			Name:          "动态分组",
			NameEn:        "Dynamic Grouping",
			Description:   "根据条件查询主机信息",
			DescriptionEn: "custom query the host instances",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/business/resource",
			},
			Version: 1,
		},
		{
			ID:            BizCustomField,
			Name:          "自定义字段",
			NameEn:        "Custom Field",
			Description:   "模型在业务下的自定义字段",
			DescriptionEn: "model's custom field under a business",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            BizProcessServiceInstance,
			Name:          "服务实例",
			NameEn:        "Service Instance",
			Description:   "服务实例",
			DescriptionEn: "service instance",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            BizProcessServiceCategory,
			Name:          "服务分类",
			NameEn:        "Service Category",
			Description:   "服务分类用于分类服务实例",
			DescriptionEn: "service category is to classify service instances",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/business/resource",
			},
			Version: 1,
		},
		{
			ID:            BizSetTemplate,
			Name:          "集群模板",
			NameEn:        "Set Template",
			Description:   "集群模板用于实例化集群",
			DescriptionEn: "set template is used to instantiate a set",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/business/resource",
			},
			Version: 1,
		},
		{
			ID:            BizTopology,
			Name:          "业务拓扑",
			NameEn:        "Business Topology",
			Description:   "业务拓扑包含了业务拓扑树中所有的相关元素",
			DescriptionEn: "business topology contains all elements related to the business topology tree",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            BizProcessServiceTemplate,
			Name:          "服务模板",
			NameEn:        "Service Template",
			Description:   "服务模板用于实例化服务实例",
			DescriptionEn: "service template is used to instantiate a service instance ",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/business/resource",
			},
			Version: 1,
		},
		// only for host topology usage, not related to actions
		//{
		//	ID:            Set,
		//	Name:          "集群",
		//	NameEn:        "Set",
		//	Description:   "集群列表",
		//	DescriptionEn: "all the sets in blueking cmdb.",
		//	Parents:       []Parent{businessParent},
		//	ProviderConfig: ResourceConfig{
		//		Path: "/auth/v3/find/business/resource",
		//	},
		//	Version: 1,
		//},
		//{
		//	ID:            Module,
		//	Name:          "模块",
		//	NameEn:        "Module",
		//	Description:   "模块列表",
		//	DescriptionEn: "all the modules in blueking cmdb.",
		//	Parents: []Parent{{
		//		SystemID:   SystemIDCMDB,
		//		ResourceID: Set,
		//	}},
		//	ProviderConfig: ResourceConfig{
		//		Path: "/auth/v3/find/business/resource",
		//	},
		//	Version: 1,
		//},
	}
}

func genPublicResources() []ResourceType {
	return []ResourceType{
		{
			ID:            Business,
			Name:          "业务",
			NameEn:        "Business",
			Description:   "业务列表",
			DescriptionEn: "all the business in blueking cmdb.",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/instance/resource",
			},
			Version: 1,
		},
		{
			ID:            SysCloudArea,
			Name:          "云区域",
			NameEn:        "Cloud Area",
			Description:   "云区域",
			DescriptionEn: "cloud area",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/instance/resource",
			},
			Version: 1,
		},
		{
			ID:            SysResourcePoolDirectory,
			Name:          "主机池目录",
			NameEn:        "Host Pool Directory",
			Description:   "主机池目录",
			DescriptionEn: "host pool directory",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysModelGroup,
			Name:          "模型分组",
			NameEn:        "Model Group",
			Description:   "模型分组用于对模型进行归类",
			DescriptionEn: "group models by model group",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysInstanceModel,
			Name:          "实例模型",
			NameEn:        "InstanceModel",
			Description:   "实例模型",
			DescriptionEn: "instance model",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysModel,
			Name:          "模型",
			NameEn:        "Model",
			Description:   "模型",
			DescriptionEn: "model",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysInstance,
			Name:          "实例",
			NameEn:        "Instance",
			Description:   "模型实例",
			DescriptionEn: "model instance",
			Parents: []Parent{{
				SystemID:   SystemIDCMDB,
				ResourceID: SysModel,
			}},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/instance/resource",
			},
			Version: 1,
		},
		{
			ID:            SysAssociationType,
			Name:          "关联类型",
			NameEn:        "Association Type",
			Description:   "关联类型是模型关联关系的分类",
			DescriptionEn: "association type is the classification of model association",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysEventPushing,
			Name:          "事件订阅",
			NameEn:        "Event Subscription",
			Description:   "当配置发生变化时推送事件",
			DescriptionEn: "push event when configuration changes",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysSystemBase,
			Name:          "系统基础",
			NameEn:        "System",
			Description:   "基础系统资源",
			DescriptionEn: "basic system resource",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            SysOperationStatistic,
			Name:          "运营统计",
			NameEn:        "Operational Statistics",
			Description:   "运营统计",
			DescriptionEn: "operational statistics",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            SysAuditLog,
			Name:          "操作审计",
			NameEn:        "Operation Audit",
			Description:   "操作审计",
			DescriptionEn: "audit log",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
		{
			ID:            SysCloudAccount,
			Name:          "云账户",
			NameEn:        "Cloud Account",
			Description:   "云账户",
			DescriptionEn: "cloud account",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysCloudResourceTask,
			Name:          "云资源任务",
			NameEn:        "Cloud Resource Task",
			Description:   "云资源任务",
			DescriptionEn: "cloud resource task",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/system/resource",
			},
			Version: 1,
		},
		{
			ID:            SysEventWatch,
			Name:          "事件监听",
			NameEn:        "Event Listen",
			Description:   "事件监听",
			DescriptionEn: "event watch",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/empty/resource",
			},
			Version: 1,
		},
	}
}
