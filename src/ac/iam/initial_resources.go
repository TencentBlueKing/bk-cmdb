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
	"configcenter/pkg/tenant/tools"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/iam"
)

var (
	businessParent = iam.Parent{
		SystemID:   iamtypes.SystemIDCMDB,
		ResourceID: iamtypes.Business,
	}
)

// ResourceTypeIDMap TODO
var ResourceTypeIDMap = map[iamtypes.TypeID]string{
	iamtypes.Business:                 "业务",
	iamtypes.BizSet:                   "业务集",
	iamtypes.Project:                  "项目",
	iamtypes.BusinessForHostTrans:     "业务主机",
	iamtypes.SysCloudArea:             "管控区域",
	iamtypes.SysResourcePoolDirectory: "主机池目录",
	iamtypes.SysHostRscPoolDirectory:  "主机池主机",
	iamtypes.SysModelGroup:            "模型分组",
	iamtypes.SysInstanceModel:         "实例模型",
	iamtypes.SysModel:                 "模型",
	iamtypes.SysModelEvent:            "模型列表",
	iamtypes.MainlineModelEvent:       "资源事件",
	iamtypes.InstAsstEvent:            "实例关联事件",
	iamtypes.KubeWorkloadEvent:        "容器工作负载事件",
	// SysInstance:               "实例",
	iamtypes.SysAssociationType:        "关联类型",
	iamtypes.SysAuditLog:               "操作审计",
	iamtypes.SysEventWatch:             "事件监听",
	iamtypes.Host:                      "主机",
	iamtypes.BizHostApply:              "主机自动应用",
	iamtypes.BizCustomQuery:            "动态分组",
	iamtypes.BizCustomField:            "自定义字段",
	iamtypes.BizProcessServiceInstance: "服务实例",
	iamtypes.BizProcessServiceCategory: "服务分类",
	iamtypes.BizSetTemplate:            "集群模板",
	iamtypes.BizTopology:               "业务拓扑",
	iamtypes.BizProcessServiceTemplate: "服务模板",
	iamtypes.FieldGroupingTemplate:     "字段组合模板",
	iamtypes.GeneralCache:              "通用缓存",
	iamtypes.Set:                       "集群",
	iamtypes.Module:                    "模块",
	iamtypes.TenantSet:                 "租户集",
}

// GenerateResourceTypes generate all the resource types registered to IAM.
func GenerateResourceTypes(tenantObjects map[string][]metadata.Object) []iam.ResourceType {
	resourceTypeList := make([]iam.ResourceType, 0)

	// add public and business resources
	resourceTypeList = append(resourceTypeList, GenerateStaticResourceTypes()...)

	// add dynamic resources
	resourceTypeList = append(resourceTypeList, genDynamicResourceTypes(tenantObjects)...)

	return resourceTypeList
}

// GenerateStaticResourceTypes TODO
func GenerateStaticResourceTypes() []iam.ResourceType {
	resourceTypeList := make([]iam.ResourceType, 0)

	// add public resources
	resourceTypeList = append(resourceTypeList, genPublicResources()...)
	resourceTypeList = append(resourceTypeList, genTenantSetResources()...)

	// add business resources
	resourceTypeList = append(resourceTypeList, genBusinessResources()...)
	return resourceTypeList
}

// GetResourceParentMap generate resource types' mapping to parents.
func GetResourceParentMap() map[iamtypes.TypeID][]iamtypes.TypeID {
	resourceParentMap := make(map[iamtypes.TypeID][]iamtypes.TypeID, 0)
	for _, resourceType := range GenerateStaticResourceTypes() {
		for _, parent := range resourceType.Parents {
			resourceParentMap[resourceType.ID] = append(resourceParentMap[resourceType.ID], parent.ResourceID)
		}
	}
	return resourceParentMap
}

func genBusinessResources() []iam.ResourceType {
	return []iam.ResourceType{
		{
			ID:            iamtypes.Host,
			Name:          ResourceTypeIDMap[iamtypes.Host],
			NameEn:        "Host",
			Description:   "主机",
			DescriptionEn: "hosts under a business or in resource pool",
			Parents: []iam.Parent{{
				SystemID: iamtypes.SystemIDCMDB,
				// ResourceID: Module,
				ResourceID: iamtypes.Business,
			}, {
				SystemID:   iamtypes.SystemIDCMDB,
				ResourceID: iamtypes.SysResourcePoolDirectory,
			}},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizHostApply,
			Name:          ResourceTypeIDMap[iamtypes.BizHostApply],
			NameEn:        "Host Apply",
			Description:   "自动应用业务主机的属性信息",
			DescriptionEn: "apply business host attribute automatically",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizCustomQuery,
			Name:          ResourceTypeIDMap[iamtypes.BizCustomQuery],
			NameEn:        "Dynamic Grouping",
			Description:   "根据条件查询主机信息",
			DescriptionEn: "custom query the host instances",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizCustomField,
			Name:          ResourceTypeIDMap[iamtypes.BizCustomField],
			NameEn:        "Custom Field",
			Description:   "模型在业务下的自定义字段",
			DescriptionEn: "model's custom field under a business",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizProcessServiceInstance,
			Name:          ResourceTypeIDMap[iamtypes.BizProcessServiceInstance],
			NameEn:        "Service Instance",
			Description:   "服务实例",
			DescriptionEn: "service instance",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizProcessServiceCategory,
			Name:          ResourceTypeIDMap[iamtypes.BizProcessServiceCategory],
			NameEn:        "Service Category",
			Description:   "服务分类用于分类服务实例",
			DescriptionEn: "service category is to classify service instances",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizSetTemplate,
			Name:          ResourceTypeIDMap[iamtypes.BizSetTemplate],
			NameEn:        "Set Template",
			Description:   "集群模板用于实例化集群",
			DescriptionEn: "set template is used to instantiate a set",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizTopology,
			Name:          ResourceTypeIDMap[iamtypes.BizTopology],
			NameEn:        "Business Topology",
			Description:   "业务拓扑包含了业务拓扑树中所有的相关元素",
			DescriptionEn: "business topology contains all elements related to the business topology tree",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BizProcessServiceTemplate,
			Name:          ResourceTypeIDMap[iamtypes.BizProcessServiceTemplate],
			NameEn:        "Service Template",
			Description:   "服务模板用于实例化服务实例",
			DescriptionEn: "service template is used to instantiate a service instance ",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		// only for biz topology usage, not related to actions
		{
			ID:            iamtypes.Set,
			Name:          ResourceTypeIDMap[iamtypes.Set],
			NameEn:        "Set",
			Description:   "业务拓扑集群",
			DescriptionEn: "business topology set",
			Parents:       []iam.Parent{businessParent},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.Module,
			Name:          ResourceTypeIDMap[iamtypes.Module],
			NameEn:        "Module",
			Description:   "业务拓扑模块",
			DescriptionEn: "business topology module",
			Parents: []iam.Parent{{
				SystemID:   iamtypes.SystemIDCMDB,
				ResourceID: iamtypes.Set,
			}},
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
	}
}

func genPublicResources() []iam.ResourceType {
	return []iam.ResourceType{
		{
			ID:            iamtypes.BizSet,
			Name:          ResourceTypeIDMap[iamtypes.BizSet],
			NameEn:        "Business Set",
			Description:   "业务集",
			DescriptionEn: "business set",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.Business,
			Name:          ResourceTypeIDMap[iamtypes.Business],
			NameEn:        "Business",
			Description:   "业务列表",
			DescriptionEn: "all the business in blueking cmdb.",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.Project,
			Name:          ResourceTypeIDMap[iamtypes.Project],
			NameEn:        "Project",
			Description:   "项目列表",
			DescriptionEn: "all the project in blueking cmdb.",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.BusinessForHostTrans,
			Name:          ResourceTypeIDMap[iamtypes.BusinessForHostTrans],
			NameEn:        "Host In Business",
			Description:   "业务主机",
			DescriptionEn: "host in business",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysCloudArea,
			Name:          ResourceTypeIDMap[iamtypes.SysCloudArea],
			NameEn:        "Cloud Area",
			Description:   "管控区域",
			DescriptionEn: "cloud area",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysResourcePoolDirectory,
			Name:          ResourceTypeIDMap[iamtypes.SysResourcePoolDirectory],
			NameEn:        "Host Pool Directory",
			Description:   "主机池目录",
			DescriptionEn: "host pool directory",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysHostRscPoolDirectory,
			Name:          ResourceTypeIDMap[iamtypes.SysHostRscPoolDirectory],
			NameEn:        "Host In Host Pool Directory",
			Description:   "主机池主机",
			DescriptionEn: "host in host pool directory",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysModelGroup,
			Name:          ResourceTypeIDMap[iamtypes.SysModelGroup],
			NameEn:        "Model Group",
			Description:   "模型分组用于对模型进行归类",
			DescriptionEn: "group models by model group",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysInstanceModel,
			Name:          ResourceTypeIDMap[iamtypes.SysInstanceModel],
			NameEn:        "InstanceModel",
			Description:   "实例模型",
			DescriptionEn: "instance model",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysModel,
			Name:          ResourceTypeIDMap[iamtypes.SysModel],
			NameEn:        "Model",
			Description:   "模型",
			DescriptionEn: "model",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysAssociationType,
			Name:          ResourceTypeIDMap[iamtypes.SysAssociationType],
			NameEn:        "Association Type",
			Description:   "关联类型是模型关联关系的分类",
			DescriptionEn: "association type is the classification of model association",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysAuditLog,
			Name:          ResourceTypeIDMap[iamtypes.SysAuditLog],
			NameEn:        "Operation Audit",
			Description:   "操作审计",
			DescriptionEn: "audit log",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysEventWatch,
			Name:          ResourceTypeIDMap[iamtypes.SysEventWatch],
			NameEn:        "Event Listen",
			Description:   "事件监听",
			DescriptionEn: "event watch",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.SysModelEvent,
			Name:          ResourceTypeIDMap[iamtypes.SysModelEvent],
			NameEn:        "Model List",
			Description:   "模型列表",
			DescriptionEn: "model list",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.MainlineModelEvent,
			Name:          ResourceTypeIDMap[iamtypes.MainlineModelEvent],
			NameEn:        "Resource Event",
			Description:   "资源事件",
			DescriptionEn: "resource event",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.InstAsstEvent,
			Name:          ResourceTypeIDMap[iamtypes.InstAsstEvent],
			NameEn:        "Instance Association Event",
			Description:   "实例关联事件",
			DescriptionEn: "instance association event",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.KubeWorkloadEvent,
			Name:          ResourceTypeIDMap[iamtypes.KubeWorkloadEvent],
			NameEn:        "Kube Workload Event",
			Description:   "容器工作负载事件",
			DescriptionEn: "kube workload event",
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.FieldGroupingTemplate,
			Name:          ResourceTypeIDMap[iamtypes.FieldGroupingTemplate],
			NameEn:        "Field Grouping Template",
			Description:   "字段组合模板",
			DescriptionEn: "Field Grouping Template",
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            iamtypes.GeneralCache,
			Name:          ResourceTypeIDMap[iamtypes.GeneralCache],
			NameEn:        "General Resource Cache",
			Description:   "通用缓存",
			DescriptionEn: "general resource cache",
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
	}
}

func genTenantSetResources() []iam.ResourceType {
	if tools.GetDefaultTenant() != common.BKDefaultTenantID {
		return make([]iam.ResourceType, 0)
	}

	return []iam.ResourceType{
		{
			ID:            iamtypes.TenantSet,
			Name:          ResourceTypeIDMap[iamtypes.TenantSet],
			NameEn:        "Tenant Set",
			Description:   "租户集",
			DescriptionEn: "tenant set",
			Parents:       nil,
			ProviderConfig: iam.ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version:  1,
			TenantID: common.BKDefaultTenantID,
		},
	}
}
