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

// ActionIDNameMap TODO
var ActionIDNameMap = map[iamtypes.ActionID]string{
	iamtypes.EditBusinessHost:              "业务主机编辑",
	iamtypes.TransferHostOutOfBiz:          "主机转出业务",
	iamtypes.TransferHostIntoBiz:           "主机转入业务",
	iamtypes.CreateBusinessCustomQuery:     "动态分组新建",
	iamtypes.EditBusinessCustomQuery:       "动态分组编辑",
	iamtypes.DeleteBusinessCustomQuery:     "动态分组删除",
	iamtypes.EditBusinessCustomField:       "业务自定义字段编辑",
	iamtypes.CreateBusinessServiceCategory: "服务分类新建",
	iamtypes.EditBusinessServiceCategory:   "服务分类编辑",
	iamtypes.DeleteBusinessServiceCategory: "服务分类删除",
	iamtypes.CreateBusinessServiceInstance: "服务实例新建",
	iamtypes.EditBusinessServiceInstance:   "服务实例编辑",
	iamtypes.DeleteBusinessServiceInstance: "服务实例删除",
	iamtypes.CreateBusinessServiceTemplate: "服务模板新建",
	iamtypes.EditBusinessServiceTemplate:   "服务模板编辑",
	iamtypes.DeleteBusinessServiceTemplate: "服务模板删除",
	iamtypes.CreateBusinessSetTemplate:     "集群模板新建",
	iamtypes.EditBusinessSetTemplate:       "集群模板编辑",
	iamtypes.DeleteBusinessSetTemplate:     "集群模板删除",
	iamtypes.CreateBusinessTopology:        "业务拓扑新建",
	iamtypes.EditBusinessTopology:          "业务拓扑编辑",
	iamtypes.DeleteBusinessTopology:        "业务拓扑删除",
	iamtypes.EditBusinessHostApply:         "主机自动应用编辑",
	iamtypes.ViewResourcePoolHost:          "主机池主机查看",
	iamtypes.CreateResourcePoolHost:        "主机池主机创建",
	iamtypes.EditResourcePoolHost:          "主机池主机编辑",
	iamtypes.DeleteResourcePoolHost:        "主机池主机删除",
	iamtypes.TransferHostOutOfResPoolDir:   "主机转出主机池目录",
	iamtypes.TransferHostToResPoolDir:      "主机转入主机池目录",
	iamtypes.CreateResourcePoolDirectory:   "主机池目录创建",
	iamtypes.EditResourcePoolDirectory:     "主机池目录编辑",
	iamtypes.DeleteResourcePoolDirectory:   "主机池目录删除",
	iamtypes.CreateBusiness:                "业务创建",
	iamtypes.EditBusiness:                  "业务编辑",
	iamtypes.ArchiveBusiness:               "业务归档",
	iamtypes.FindBusiness:                  "业务查询",
	iamtypes.ViewBusinessResource:          "业务访问",
	iamtypes.CreateBizSet:                  "业务集新增",
	iamtypes.EditBizSet:                    "业务集编辑",
	iamtypes.DeleteBizSet:                  "业务集删除",
	iamtypes.ViewBizSet:                    "业务集查看",
	iamtypes.AccessBizSet:                  "业务集访问",
	iamtypes.CreateProject:                 "项目新建",
	iamtypes.EditProject:                   "项目编辑",
	iamtypes.DeleteProject:                 "项目删除",
	iamtypes.ViewProject:                   "项目查看",
	iamtypes.ViewCloudArea:                 "管控区域查看",
	iamtypes.CreateCloudArea:               "管控区域创建",
	iamtypes.EditCloudArea:                 "管控区域编辑",
	iamtypes.DeleteCloudArea:               "管控区域删除",
	iamtypes.ViewSysModel:                  "模型查看",
	iamtypes.CreateSysModel:                "模型新建",
	iamtypes.EditSysModel:                  "模型编辑",
	iamtypes.DeleteSysModel:                "模型删除",
	iamtypes.CreateAssociationType:         "关联类型新建",
	iamtypes.EditAssociationType:           "关联类型编辑",
	iamtypes.DeleteAssociationType:         "关联类型删除",
	iamtypes.CreateModelGroup:              "模型分组新建",
	iamtypes.EditModelGroup:                "模型分组编辑",
	iamtypes.DeleteModelGroup:              "模型分组删除",
	iamtypes.ViewModelTopo:                 "模型拓扑查看",
	iamtypes.EditBusinessLayer:             "业务层级编辑",
	iamtypes.EditModelTopologyView:         "模型拓扑视图编辑",
	iamtypes.FindAuditLog:                  "操作审计查询",
	iamtypes.WatchHostEvent:                "主机事件监听",
	iamtypes.WatchHostRelationEvent:        "主机关系事件监听",
	iamtypes.WatchBizEvent:                 "业务事件监听",
	iamtypes.WatchSetEvent:                 "集群事件监听",
	iamtypes.WatchModuleEvent:              "模块数据监听",
	iamtypes.WatchProcessEvent:             "进程数据监听",
	iamtypes.WatchCommonInstanceEvent:      "模型实例事件监听",
	iamtypes.WatchMainlineInstanceEvent:    "自定义拓扑层级事件监听",
	iamtypes.WatchInstAsstEvent:            "实例关联事件监听",
	iamtypes.WatchBizSetEvent:              "业务集事件监听",
	iamtypes.WatchPlatEvent:                "管控区域事件监听",
	iamtypes.WatchProjectEvent:             "项目事件监听",
	iamtypes.GlobalSettings:                "全局设置",
	iamtypes.ManageHostAgentID:             "主机AgentID管理",
	iamtypes.CreateContainerCluster:        "容器集群新建",
	iamtypes.EditContainerCluster:          "容器集群编辑",
	iamtypes.DeleteContainerCluster:        "容器集群删除",
	iamtypes.CreateContainerNode:           "容器集群节点新建",
	iamtypes.EditContainerNode:             "容器集群节点编辑",
	iamtypes.DeleteContainerNode:           "容器集群节点删除",
	iamtypes.CreateContainerNamespace:      "容器命名空间新建",
	iamtypes.EditContainerNamespace:        "容器命名空间编辑",
	iamtypes.DeleteContainerNamespace:      "容器命名空间删除",
	iamtypes.CreateContainerWorkload:       "容器工作负载新建",
	iamtypes.EditContainerWorkload:         "容器工作负载编辑",
	iamtypes.DeleteContainerWorkload:       "容器工作负载删除",
	iamtypes.CreateContainerPod:            "容器Pod新建",
	iamtypes.DeleteContainerPod:            "容器Pod删除",
	iamtypes.UseFulltextSearch:             "全文检索",
	iamtypes.CreateFieldGroupingTemplate:   "字段组合模板新建",
	iamtypes.ViewFieldGroupingTemplate:     "字段组合模板查看",
	iamtypes.EditFieldGroupingTemplate:     "字段组合模板编辑",
	iamtypes.DeleteFieldGroupingTemplate:   "字段组合模板删除",
	iamtypes.EditIDRuleIncrID:              "ID规则自增ID编辑",
	iamtypes.CreateFullSyncCond:            "全量同步缓存条件新建",
	iamtypes.ViewFullSyncCond:              "全量同步缓存条件查看",
	iamtypes.EditFullSyncCond:              "全量同步缓存条件编辑",
	iamtypes.DeleteFullSyncCond:            "全量同步缓存条件删除",
	iamtypes.ViewGeneralCache:              "通用缓存查询",
	iamtypes.ViewTenantSet:                 "租户集查看",
	iamtypes.AccessTenantSet:               "租户集访问",
}

// GenerateActions generate all the actions registered to IAM.
func GenerateActions(tenantObjects map[string][]metadata.Object) []iam.ResourceAction {
	resourceActionList := GenerateStaticActions()
	resourceActionList = append(resourceActionList, genDynamicActions(tenantObjects)...)
	return resourceActionList
}

// GenerateStaticActions TODO
func GenerateStaticActions() []iam.ResourceAction {
	resourceActionList := make([]iam.ResourceAction, 0)
	// add business resource actions
	resourceActionList = append(resourceActionList, genBusinessHostActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomQueryActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomFieldActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceCategoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceInstanceActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessSetTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessTopologyActions()...)
	resourceActionList = append(resourceActionList, genBusinessHostApplyActions()...)

	// add public resource actions
	resourceActionList = append(resourceActionList, genResourcePoolHostActions()...)
	resourceActionList = append(resourceActionList, genResourcePoolDirectoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessActions()...)
	resourceActionList = append(resourceActionList, genBizSetActions()...)
	resourceActionList = append(resourceActionList, genProjectActions()...)
	resourceActionList = append(resourceActionList, genCloudAreaActions()...)
	resourceActionList = append(resourceActionList, genModelActions()...)
	resourceActionList = append(resourceActionList, genAssociationTypeActions()...)
	resourceActionList = append(resourceActionList, genModelGroupActions()...)
	resourceActionList = append(resourceActionList, genBusinessLayerActions()...)
	resourceActionList = append(resourceActionList, genModelTopologyViewActions()...)
	resourceActionList = append(resourceActionList, genAuditLogActions()...)
	resourceActionList = append(resourceActionList, genEventWatchActions()...)
	resourceActionList = append(resourceActionList, genConfigAdminActions()...)
	resourceActionList = append(resourceActionList, genContainerManagementActions()...)
	resourceActionList = append(resourceActionList, genFulltextSearchActions()...)
	resourceActionList = append(resourceActionList, genFieldGroupingTemplateActions()...)
	resourceActionList = append(resourceActionList, genIDRuleActions()...)
	resourceActionList = append(resourceActionList, genFullSyncCondActions()...)
	resourceActionList = append(resourceActionList, genCacheActions()...)
	resourceActionList = append(resourceActionList, genTenantSetActions()...)

	return resourceActionList
}

func genBusinessHostActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.EditBusinessHost,
			Name:           ActionIDNameMap[iamtypes.EditBusinessHost],
			NameEn:         "Edit Business Hosts",
			ResourceTypeID: iamtypes.Host,
		},
		{
			ID:             iamtypes.TransferHostOutOfBiz,
			Name:           ActionIDNameMap[iamtypes.TransferHostOutOfBiz],
			NameEn:         "Transfer Host Out Of Business",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.TransferHostIntoBiz,
			Name:           ActionIDNameMap[iamtypes.TransferHostIntoBiz],
			NameEn:         "Transfer Host Into Business",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genBusinessCustomQueryActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.CreateBusinessCustomQuery,
			Name:           ActionIDNameMap[iamtypes.CreateBusinessCustomQuery],
			NameEn:         "Create Dynamic Grouping",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.EditBusinessCustomQuery,
			Name:           ActionIDNameMap[iamtypes.EditBusinessCustomQuery],
			NameEn:         "Edit Dynamic Grouping",
			ResourceTypeID: iamtypes.BizCustomQuery,
		},
		{
			ID:             iamtypes.DeleteBusinessCustomQuery,
			Name:           ActionIDNameMap[iamtypes.DeleteBusinessCustomQuery],
			NameEn:         "Delete Dynamic Grouping",
			ResourceTypeID: iamtypes.BizCustomQuery,
		},
	}
}

func genBusinessCustomFieldActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.EditBusinessCustomField,
			Name:           ActionIDNameMap[iamtypes.EditBusinessCustomField],
			NameEn:         "Edit Custom Field",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genBusinessServiceCategoryActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.CreateBusinessServiceCategory,
			Name:           ActionIDNameMap[iamtypes.CreateBusinessServiceCategory],
			NameEn:         "Create Service Category",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.EditBusinessServiceCategory,
			Name:           ActionIDNameMap[iamtypes.EditBusinessServiceCategory],
			NameEn:         "Edit Service Category",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.DeleteBusinessServiceCategory,
			Name:           ActionIDNameMap[iamtypes.DeleteBusinessServiceCategory],
			NameEn:         "Delete Service Category",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genBusinessServiceInstanceActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.CreateBusinessServiceInstance,
			Name:           ActionIDNameMap[iamtypes.CreateBusinessServiceInstance],
			NameEn:         "Create Service Instance",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.EditBusinessServiceInstance,
			Name:           ActionIDNameMap[iamtypes.EditBusinessServiceInstance],
			NameEn:         "Edit Service Instance",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.DeleteBusinessServiceInstance,
			Name:           ActionIDNameMap[iamtypes.DeleteBusinessServiceInstance],
			NameEn:         "Delete Service Instance",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genBusinessServiceTemplateActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.CreateBusinessServiceTemplate,
			Name:           ActionIDNameMap[iamtypes.CreateBusinessServiceTemplate],
			NameEn:         "Create Service Template",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.EditBusinessServiceTemplate,
			Name:           ActionIDNameMap[iamtypes.EditBusinessServiceTemplate],
			NameEn:         "Edit Service Template",
			ResourceTypeID: iamtypes.BizProcessServiceTemplate,
		},
		{
			ID:             iamtypes.DeleteBusinessServiceTemplate,
			Name:           ActionIDNameMap[iamtypes.DeleteBusinessServiceTemplate],
			NameEn:         "Delete Service Template",
			ResourceTypeID: iamtypes.BizProcessServiceTemplate,
		},
	}
}

func genBusinessSetTemplateActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.CreateBusinessSetTemplate,
			Name:           ActionIDNameMap[iamtypes.CreateBusinessSetTemplate],
			NameEn:         "Create Set Template",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.EditBusinessSetTemplate,
			Name:           ActionIDNameMap[iamtypes.EditBusinessSetTemplate],
			NameEn:         "Edit Set Template",
			ResourceTypeID: iamtypes.BizSetTemplate,
		},
		{
			ID:             iamtypes.DeleteBusinessSetTemplate,
			Name:           ActionIDNameMap[iamtypes.DeleteBusinessSetTemplate],
			NameEn:         "Delete Set Template",
			ResourceTypeID: iamtypes.BizSetTemplate,
		},
	}
}

func genBusinessTopologyActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.CreateBusinessTopology,
			Name:           ActionIDNameMap[iamtypes.CreateBusinessTopology],
			NameEn:         "Create Business Topo",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.EditBusinessTopology,
			Name:           ActionIDNameMap[iamtypes.EditBusinessTopology],
			NameEn:         "Edit Business Topo",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.DeleteBusinessTopology,
			Name:           ActionIDNameMap[iamtypes.DeleteBusinessTopology],
			NameEn:         "Delete Business Topo",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genBusinessHostApplyActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.EditBusinessHostApply,
			Name:           ActionIDNameMap[iamtypes.EditBusinessHostApply],
			NameEn:         "Edit Host Apply",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genResourcePoolHostActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.ViewResourcePoolHost,
			Name:   ActionIDNameMap[iamtypes.ViewResourcePoolHost],
			NameEn: "View Resource Pool Hosts",
		},
		{
			ID:             iamtypes.CreateResourcePoolHost,
			Name:           ActionIDNameMap[iamtypes.CreateResourcePoolHost],
			NameEn:         "Create Pool Hosts",
			ResourceTypeID: iamtypes.SysResourcePoolDirectory,
		},
		{
			ID:             iamtypes.EditResourcePoolHost,
			Name:           ActionIDNameMap[iamtypes.EditResourcePoolHost],
			NameEn:         "Edit Pool Hosts",
			ResourceTypeID: iamtypes.SysHost,
		},
		{
			ID:             iamtypes.DeleteResourcePoolHost,
			Name:           ActionIDNameMap[iamtypes.DeleteResourcePoolHost],
			NameEn:         "Delete Pool Hosts",
			ResourceTypeID: iamtypes.SysHost,
		},
		{
			ID:             iamtypes.TransferHostOutOfResPoolDir,
			Name:           ActionIDNameMap[iamtypes.TransferHostOutOfResPoolDir],
			NameEn:         "Transfer Host Out Of Pool Directory",
			ResourceTypeID: iamtypes.SysResourcePoolDirectory,
		},
		{
			ID:             iamtypes.TransferHostToResPoolDir,
			Name:           ActionIDNameMap[iamtypes.TransferHostToResPoolDir],
			NameEn:         "Transfer Host To Pool Directory",
			ResourceTypeID: iamtypes.SysResourcePoolDirectory,
		},
		{
			ID:     iamtypes.ManageHostAgentID,
			Name:   ActionIDNameMap[iamtypes.ManageHostAgentID],
			NameEn: "Manage Host AgentID",
		},
	}
}

func genResourcePoolDirectoryActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateResourcePoolDirectory,
			Name:   ActionIDNameMap[iamtypes.CreateResourcePoolDirectory],
			NameEn: "Create Pool Directory",
		},
		{
			ID:             iamtypes.EditResourcePoolDirectory,
			Name:           ActionIDNameMap[iamtypes.EditResourcePoolDirectory],
			NameEn:         "Edit Pool Directory",
			ResourceTypeID: iamtypes.SysResourcePoolDirectory,
		},
		{
			ID:             iamtypes.DeleteResourcePoolDirectory,
			Name:           ActionIDNameMap[iamtypes.DeleteResourcePoolDirectory],
			NameEn:         "Delete Pool Directory",
			ResourceTypeID: iamtypes.SysResourcePoolDirectory,
		},
	}
}

func genBusinessActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateBusiness,
			Name:   ActionIDNameMap[iamtypes.CreateBusiness],
			NameEn: "Create Business",
		},
		{
			ID:             iamtypes.EditBusiness,
			Name:           ActionIDNameMap[iamtypes.EditBusiness],
			NameEn:         "Edit Business",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.ArchiveBusiness,
			Name:           ActionIDNameMap[iamtypes.ArchiveBusiness],
			NameEn:         "Archive Business",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.FindBusiness,
			Name:           ActionIDNameMap[iamtypes.FindBusiness],
			NameEn:         "View Business",
			ResourceTypeID: iamtypes.Business,
		},
		{
			ID:             iamtypes.ViewBusinessResource,
			Name:           ActionIDNameMap[iamtypes.ViewBusinessResource],
			NameEn:         "View Business Resource",
			ResourceTypeID: iamtypes.Business,
		},
	}
}

func genBizSetActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateBizSet,
			Name:   ActionIDNameMap[iamtypes.CreateBizSet],
			NameEn: "Create Business Set",
		},
		{
			ID:             iamtypes.EditBizSet,
			Name:           ActionIDNameMap[iamtypes.EditBizSet],
			NameEn:         "Edit Business Set",
			ResourceTypeID: iamtypes.BizSet,
		},
		{
			ID:             iamtypes.DeleteBizSet,
			Name:           ActionIDNameMap[iamtypes.DeleteBizSet],
			NameEn:         "Delete Business Set",
			ResourceTypeID: iamtypes.BizSet,
		},
		{
			ID:             iamtypes.ViewBizSet,
			Name:           ActionIDNameMap[iamtypes.ViewBizSet],
			NameEn:         "View Business Set",
			ResourceTypeID: iamtypes.BizSet,
		},
		{
			ID:             iamtypes.AccessBizSet,
			Name:           ActionIDNameMap[iamtypes.AccessBizSet],
			NameEn:         "Access Business Set",
			ResourceTypeID: iamtypes.BizSet,
		},
	}
}

func genProjectActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateProject,
			Name:   ActionIDNameMap[iamtypes.CreateProject],
			NameEn: "Create Project",
		},
		{
			ID:             iamtypes.EditProject,
			Name:           ActionIDNameMap[iamtypes.EditProject],
			NameEn:         "Edit Project",
			ResourceTypeID: iamtypes.Project,
		},
		{
			ID:             iamtypes.DeleteProject,
			Name:           ActionIDNameMap[iamtypes.DeleteProject],
			NameEn:         "Delete Project",
			ResourceTypeID: iamtypes.Project,
		},
		{
			ID:     iamtypes.ViewProject,
			Name:   ActionIDNameMap[iamtypes.ViewProject],
			NameEn: "View Project",
		},
	}
}

func genCloudAreaActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.ViewCloudArea,
			Name:   ActionIDNameMap[iamtypes.ViewCloudArea],
			NameEn: "View Cloud Area",
		},
		{
			ID:     iamtypes.CreateCloudArea,
			Name:   ActionIDNameMap[iamtypes.CreateCloudArea],
			NameEn: "Create Cloud Area",
		},
		{
			ID:             iamtypes.EditCloudArea,
			Name:           ActionIDNameMap[iamtypes.EditCloudArea],
			NameEn:         "Edit Cloud Area",
			ResourceTypeID: iamtypes.SysCloudArea,
		},
		{
			ID:             iamtypes.DeleteCloudArea,
			Name:           ActionIDNameMap[iamtypes.DeleteCloudArea],
			NameEn:         "Delete Cloud Area",
			ResourceTypeID: iamtypes.SysCloudArea,
		},
	}
}

func genModelActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.ViewSysModel,
			Name:           ActionIDNameMap[iamtypes.ViewSysModel],
			NameEn:         "View Model",
			ResourceTypeID: iamtypes.SysModel,
		},
		{
			ID:             iamtypes.CreateSysModel,
			Name:           ActionIDNameMap[iamtypes.CreateSysModel],
			NameEn:         "Create Model",
			ResourceTypeID: iamtypes.SysModelGroup,
		},
		{
			ID:             iamtypes.EditSysModel,
			Name:           ActionIDNameMap[iamtypes.EditSysModel],
			NameEn:         "Edit Model",
			ResourceTypeID: iamtypes.SysModel,
		},
		{
			ID:             iamtypes.DeleteSysModel,
			Name:           ActionIDNameMap[iamtypes.DeleteSysModel],
			NameEn:         "Delete Model",
			ResourceTypeID: iamtypes.SysModel,
		},
	}
}

func genAssociationTypeActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateAssociationType,
			Name:   ActionIDNameMap[iamtypes.CreateAssociationType],
			NameEn: "Create Association Type",
		},
		{
			ID:             iamtypes.EditAssociationType,
			Name:           ActionIDNameMap[iamtypes.EditAssociationType],
			NameEn:         "Edit Association Type",
			ResourceTypeID: iamtypes.SysAssociationType,
		},
		{
			ID:             iamtypes.DeleteAssociationType,
			Name:           ActionIDNameMap[iamtypes.DeleteAssociationType],
			NameEn:         "Delete Association Type",
			ResourceTypeID: iamtypes.SysAssociationType,
		},
	}
}

func genModelGroupActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateModelGroup,
			Name:   ActionIDNameMap[iamtypes.CreateModelGroup],
			NameEn: "Create Model Group",
		},
		{
			ID:             iamtypes.EditModelGroup,
			Name:           ActionIDNameMap[iamtypes.EditModelGroup],
			NameEn:         "Edit Model Group",
			ResourceTypeID: iamtypes.SysModelGroup,
		},
		{
			ID:             iamtypes.DeleteModelGroup,
			Name:           ActionIDNameMap[iamtypes.DeleteModelGroup],
			NameEn:         "Delete Model Group",
			ResourceTypeID: iamtypes.SysModelGroup,
		},
	}
}

func genBusinessLayerActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.EditBusinessLayer,
			Name:   ActionIDNameMap[iamtypes.EditBusinessLayer],
			NameEn: "Edit Business Level",
		},
	}
}

func genModelTopologyViewActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.ViewModelTopo,
			Name:   ActionIDNameMap[iamtypes.ViewModelTopo],
			NameEn: "View Model Topo",
		},
		{
			ID:     iamtypes.EditModelTopologyView,
			Name:   ActionIDNameMap[iamtypes.EditModelTopologyView],
			NameEn: "Edit Model Topo View",
		},
	}
}

func genAuditLogActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.FindAuditLog,
			Name:   ActionIDNameMap[iamtypes.FindAuditLog],
			NameEn: "View Operation Audit",
		},
	}
}

func genEventWatchActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.WatchHostEvent,
			Name:   ActionIDNameMap[iamtypes.WatchHostEvent],
			NameEn: "Host Event Listen",
		},
		{
			ID:     iamtypes.WatchHostRelationEvent,
			Name:   ActionIDNameMap[iamtypes.WatchHostRelationEvent],
			NameEn: "Host Relation Event Listen",
		},
		{
			ID:     iamtypes.WatchBizEvent,
			Name:   ActionIDNameMap[iamtypes.WatchBizEvent],
			NameEn: "Business Event Listen",
		},
		{
			ID:     iamtypes.WatchSetEvent,
			Name:   ActionIDNameMap[iamtypes.WatchSetEvent],
			NameEn: "Set Event Listen",
		},
		{
			ID:     iamtypes.WatchModuleEvent,
			Name:   ActionIDNameMap[iamtypes.WatchModuleEvent],
			NameEn: "Module Event Listen",
		},
		{
			ID:     iamtypes.WatchProcessEvent,
			Name:   ActionIDNameMap[iamtypes.WatchProcessEvent],
			NameEn: "Process Event Listen",
		},
		{
			ID:     iamtypes.WatchBizSetEvent,
			Name:   ActionIDNameMap[iamtypes.WatchBizSetEvent],
			NameEn: "Business Set Event Listen",
		},
		{
			ID:     iamtypes.WatchPlatEvent,
			Name:   ActionIDNameMap[iamtypes.WatchPlatEvent],
			NameEn: "Cloud Area Event Listen",
		},
		{
			ID:     iamtypes.WatchProjectEvent,
			Name:   ActionIDNameMap[iamtypes.WatchProjectEvent],
			NameEn: "Project Event Listen",
		},
		{
			ID:             iamtypes.WatchCommonInstanceEvent,
			Name:           ActionIDNameMap[iamtypes.WatchCommonInstanceEvent],
			NameEn:         "Common Model Instance Event Listen",
			ResourceTypeID: iamtypes.SysModelEvent,
		},
		{
			ID:             iamtypes.WatchMainlineInstanceEvent,
			Name:           ActionIDNameMap[iamtypes.WatchMainlineInstanceEvent],
			NameEn:         "Custom Topo Layer Event Listen",
			ResourceTypeID: iamtypes.MainlineModelEvent,
		},
		{
			ID:             iamtypes.WatchInstAsstEvent,
			Name:           ActionIDNameMap[iamtypes.WatchInstAsstEvent],
			NameEn:         "Instance Association Event Listen",
			ResourceTypeID: iamtypes.InstAsstEvent,
		},
	}
}

func genConfigAdminActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.GlobalSettings,
			Name:   ActionIDNameMap[iamtypes.GlobalSettings],
			NameEn: "Global Settings",
		},
	}
}

func genContainerManagementActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, genContainerClusterActions()...)
	actions = append(actions, genContainerNodeActions()...)
	actions = append(actions, genContainerNamespaceActions()...)
	actions = append(actions, genContainerWorkloadActions()...)
	actions = append(actions, genContainerPodActions()...)

	return actions
}

func genContainerClusterActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateContainerCluster,
			Name:   ActionIDNameMap[iamtypes.CreateContainerCluster],
			NameEn: "Create Container Cluster",
			Hidden: true,
		},
		{
			ID:     iamtypes.EditContainerCluster,
			Name:   ActionIDNameMap[iamtypes.EditContainerCluster],
			NameEn: "Edit Container Cluster",
			Hidden: true,
		},
		{
			ID:     iamtypes.DeleteContainerCluster,
			Name:   ActionIDNameMap[iamtypes.DeleteContainerCluster],
			NameEn: "Delete Container Cluster",
			Hidden: true,
		},
	}
}

func genContainerNodeActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateContainerNode,
			Name:   ActionIDNameMap[iamtypes.CreateContainerNode],
			NameEn: "Create Container Node",
			Hidden: true,
		},
		{
			ID:     iamtypes.EditContainerNode,
			Name:   ActionIDNameMap[iamtypes.EditContainerNode],
			NameEn: "Edit Container Node",
			Hidden: true,
		},
		{
			ID:     iamtypes.DeleteContainerNode,
			Name:   ActionIDNameMap[iamtypes.DeleteContainerNode],
			NameEn: "Delete Container Node",
			Hidden: true,
		},
	}
}

func genContainerNamespaceActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateContainerNamespace,
			Name:   ActionIDNameMap[iamtypes.CreateContainerNamespace],
			NameEn: "Create Container Namespace",
			Hidden: true,
		},
		{
			ID:     iamtypes.EditContainerNamespace,
			Name:   ActionIDNameMap[iamtypes.EditContainerNamespace],
			NameEn: "Edit Container Namespace",
			Hidden: true,
		},
		{
			ID:     iamtypes.DeleteContainerNamespace,
			Name:   ActionIDNameMap[iamtypes.DeleteContainerNamespace],
			NameEn: "Delete Container Namespace",
			Hidden: true,
		},
	}
}

func genContainerWorkloadActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateContainerWorkload,
			Name:   ActionIDNameMap[iamtypes.CreateContainerWorkload],
			NameEn: "Create Container Workload",
			Hidden: true,
		},
		{
			ID:     iamtypes.EditContainerWorkload,
			Name:   ActionIDNameMap[iamtypes.EditContainerWorkload],
			NameEn: "Edit Container Workload",
			Hidden: true,
		},
		{
			ID:     iamtypes.DeleteContainerWorkload,
			Name:   ActionIDNameMap[iamtypes.DeleteContainerWorkload],
			NameEn: "Delete Container Workload",
			Hidden: true,
		},
	}
}

func genContainerPodActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateContainerPod,
			Name:   ActionIDNameMap[iamtypes.CreateContainerPod],
			NameEn: "Create Container Pod",
			Hidden: true,
		},
		{
			ID:     iamtypes.DeleteContainerPod,
			Name:   ActionIDNameMap[iamtypes.DeleteContainerPod],
			NameEn: "Delete Container Pod",
			Hidden: true,
		},
	}
}

func genFulltextSearchActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.UseFulltextSearch,
			Name:   ActionIDNameMap[iamtypes.UseFulltextSearch],
			NameEn: "Fulltext Search",
		},
	}
}

func genFieldGroupingTemplateActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateFieldGroupingTemplate,
			Name:   ActionIDNameMap[iamtypes.CreateFieldGroupingTemplate],
			NameEn: "Create Field Grouping Template",
		},
		{
			ID:             iamtypes.ViewFieldGroupingTemplate,
			Name:           ActionIDNameMap[iamtypes.ViewFieldGroupingTemplate],
			NameEn:         "View Field Grouping Template",
			ResourceTypeID: iamtypes.FieldGroupingTemplate,
		},
		{
			ID:             iamtypes.EditFieldGroupingTemplate,
			Name:           ActionIDNameMap[iamtypes.EditFieldGroupingTemplate],
			NameEn:         "Edit Field Grouping Template",
			ResourceTypeID: iamtypes.FieldGroupingTemplate,
		},
		{
			ID:             iamtypes.DeleteFieldGroupingTemplate,
			Name:           ActionIDNameMap[iamtypes.DeleteFieldGroupingTemplate],
			NameEn:         "Delete Field Grouping Template",
			ResourceTypeID: iamtypes.FieldGroupingTemplate,
		},
	}
}

func genIDRuleActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.EditIDRuleIncrID,
			Name:   ActionIDNameMap[iamtypes.EditIDRuleIncrID],
			NameEn: "Edit ID Rule",
		},
	}
}

func genFullSyncCondActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.CreateFullSyncCond,
			Name:   ActionIDNameMap[iamtypes.CreateFullSyncCond],
			NameEn: "Create Full Sync Cond",
			Hidden: true,
		},
		{
			ID:     iamtypes.EditFullSyncCond,
			Name:   ActionIDNameMap[iamtypes.EditFullSyncCond],
			NameEn: "Edit Full Sync Cond",
			Hidden: true,
		},
		{
			ID:     iamtypes.DeleteFullSyncCond,
			Name:   ActionIDNameMap[iamtypes.DeleteFullSyncCond],
			NameEn: "Delete Full Sync Cond",
			Hidden: true,
		},
		{
			ID:     iamtypes.ViewFullSyncCond,
			Name:   ActionIDNameMap[iamtypes.ViewFullSyncCond],
			NameEn: "View Full Sync Cond",
			Hidden: true,
		},
	}
}

func genCacheActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:             iamtypes.ViewGeneralCache,
			Name:           ActionIDNameMap[iamtypes.ViewGeneralCache],
			NameEn:         "View General Resource Cache",
			ResourceTypeID: iamtypes.GeneralCache,
			Hidden:         true,
		},
	}
}

func genTenantSetActions() []iam.ResourceAction {
	if tools.GetDefaultTenant() != common.BKDefaultTenantID {
		return make([]iam.ResourceAction, 0)
	}

	return []iam.ResourceAction{
		{
			ID:             iamtypes.ViewTenantSet,
			Name:           ActionIDNameMap[iamtypes.ViewTenantSet],
			NameEn:         "View Tenant Set",
			ResourceTypeID: iamtypes.TenantSet,
			Hidden:         true,
			TenantID:       common.BKDefaultTenantID,
		},
		{
			ID:             iamtypes.AccessTenantSet,
			Name:           ActionIDNameMap[iamtypes.AccessTenantSet],
			NameEn:         "Access Tenant Set",
			ResourceTypeID: iamtypes.TenantSet,
			Hidden:         true,
			TenantID:       common.BKDefaultTenantID,
		},
	}
}
