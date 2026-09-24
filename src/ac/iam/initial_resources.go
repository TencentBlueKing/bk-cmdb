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
	"configcenter/pkg/tenant/tools"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/iam"
)

var (
	businessParent = []iamtypes.TypeID{iamtypes.Business}
)

// ResourceTypeIDMap TODO
var ResourceTypeIDMap = map[iamtypes.TypeID]string{
	iamtypes.Business:                 "业务",
	iamtypes.BizSet:                   "业务集",
	iamtypes.Project:                  "项目",
	iamtypes.SysCloudArea:             "管控区域",
	iamtypes.SysResourcePoolDirectory: "主机池目录",
	iamtypes.SysHost:                  "主机池主机",
	iamtypes.SysModelGroup:            "模型分组",
	iamtypes.SysInstanceModel:         "实例模型",
	iamtypes.SysModel:                 "模型",
	iamtypes.SysModelEvent:            "模型列表",
	iamtypes.MainlineModelEvent:       "资源事件",
	iamtypes.InstAsstEvent:            "实例关联事件",
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
		resourceParentMap[resourceType.ID] = resourceType.Ancestors
	}
	return resourceParentMap
}

func genBusinessResources() []iam.ResourceType {
	return []iam.ResourceType{
		{
			ID:        iamtypes.Host,
			Name:      ResourceTypeIDMap[iamtypes.Host],
			NameEn:    "Host",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizHostApply,
			Name:      ResourceTypeIDMap[iamtypes.BizHostApply],
			NameEn:    "Host Apply",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizCustomQuery,
			Name:      ResourceTypeIDMap[iamtypes.BizCustomQuery],
			NameEn:    "Dynamic Grouping",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizCustomField,
			Name:      ResourceTypeIDMap[iamtypes.BizCustomField],
			NameEn:    "Custom Field",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizProcessServiceInstance,
			Name:      ResourceTypeIDMap[iamtypes.BizProcessServiceInstance],
			NameEn:    "Service Instance",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizProcessServiceCategory,
			Name:      ResourceTypeIDMap[iamtypes.BizProcessServiceCategory],
			NameEn:    "Service Category",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizSetTemplate,
			Name:      ResourceTypeIDMap[iamtypes.BizSetTemplate],
			NameEn:    "Set Template",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizTopology,
			Name:      ResourceTypeIDMap[iamtypes.BizTopology],
			NameEn:    "Business Topology",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.BizProcessServiceTemplate,
			Name:      ResourceTypeIDMap[iamtypes.BizProcessServiceTemplate],
			NameEn:    "Service Template",
			Ancestors: businessParent,
		},
		// only for biz topology usage, not related to actions
		{
			ID:        iamtypes.Set,
			Name:      ResourceTypeIDMap[iamtypes.Set],
			NameEn:    "Set",
			Ancestors: businessParent,
		},
		{
			ID:        iamtypes.Module,
			Name:      ResourceTypeIDMap[iamtypes.Module],
			NameEn:    "Module",
			Ancestors: []iamtypes.TypeID{iamtypes.Set},
		},
	}
}

func genPublicResources() []iam.ResourceType {
	return []iam.ResourceType{
		{
			ID:     iamtypes.BizSet,
			Name:   ResourceTypeIDMap[iamtypes.BizSet],
			NameEn: "Business Set",
		},
		{
			ID:     iamtypes.Business,
			Name:   ResourceTypeIDMap[iamtypes.Business],
			NameEn: "Business",
		},
		{
			ID:     iamtypes.Project,
			Name:   ResourceTypeIDMap[iamtypes.Project],
			NameEn: "Project",
		},
		{
			ID:     iamtypes.SysCloudArea,
			Name:   ResourceTypeIDMap[iamtypes.SysCloudArea],
			NameEn: "Cloud Area",
		},
		{
			ID:     iamtypes.SysResourcePoolDirectory,
			Name:   ResourceTypeIDMap[iamtypes.SysResourcePoolDirectory],
			NameEn: "Host Pool Directory",
		},
		{
			ID:        iamtypes.SysHost,
			Name:      ResourceTypeIDMap[iamtypes.SysHost],
			NameEn:    "Host In Host Pool Directory",
			Ancestors: []iamtypes.TypeID{iamtypes.SysResourcePoolDirectory},
		},
		{
			ID:     iamtypes.SysModelGroup,
			Name:   ResourceTypeIDMap[iamtypes.SysModelGroup],
			NameEn: "Model Group",
		},
		{
			ID:     iamtypes.SysInstanceModel,
			Name:   ResourceTypeIDMap[iamtypes.SysInstanceModel],
			NameEn: "InstanceModel",
		},
		{
			ID:     iamtypes.SysModel,
			Name:   ResourceTypeIDMap[iamtypes.SysModel],
			NameEn: "Model",
		},
		{
			ID:     iamtypes.SysAssociationType,
			Name:   ResourceTypeIDMap[iamtypes.SysAssociationType],
			NameEn: "Association Type",
		},
		{
			ID:     iamtypes.SysAuditLog,
			Name:   ResourceTypeIDMap[iamtypes.SysAuditLog],
			NameEn: "Operation Audit",
		},
		{
			ID:     iamtypes.SysEventWatch,
			Name:   ResourceTypeIDMap[iamtypes.SysEventWatch],
			NameEn: "Event Listen",
		},
		{
			ID:     iamtypes.SysModelEvent,
			Name:   ResourceTypeIDMap[iamtypes.SysModelEvent],
			NameEn: "Model List",
		},
		{
			ID:     iamtypes.MainlineModelEvent,
			Name:   ResourceTypeIDMap[iamtypes.MainlineModelEvent],
			NameEn: "Resource Event",
		},
		{
			ID:     iamtypes.InstAsstEvent,
			Name:   ResourceTypeIDMap[iamtypes.InstAsstEvent],
			NameEn: "Instance Association Event",
		},
		{
			ID:     iamtypes.FieldGroupingTemplate,
			Name:   ResourceTypeIDMap[iamtypes.FieldGroupingTemplate],
			NameEn: "Field Grouping Template",
		},
		{
			ID:     iamtypes.GeneralCache,
			Name:   ResourceTypeIDMap[iamtypes.GeneralCache],
			NameEn: "General Resource Cache",
		},
	}
}

func genTenantSetResources() []iam.ResourceType {
	if tools.GetDefaultTenant() != common.BKDefaultTenantID {
		return make([]iam.ResourceType, 0)
	}

	return []iam.ResourceType{
		{
			ID:       iamtypes.TenantSet,
			Name:     ResourceTypeIDMap[iamtypes.TenantSet],
			NameEn:   "Tenant Set",
			TenantID: common.BKDefaultTenantID,
		},
	}
}
