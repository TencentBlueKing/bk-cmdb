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
	businessResource = iam.RelateResourceType{
		SystemID:    iamtypes.SystemIDCMDB,
		ID:          iamtypes.Business,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.BusinessSelection,
		}},
	}

	resourcePoolDirResource = iam.RelateResourceType{
		SystemID:    iamtypes.SystemIDCMDB,
		ID:          iamtypes.SysResourcePoolDirectory,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.SysResourcePoolDirectorySelection,
		}},
	}
)

// ActionIDNameMap TODO
var ActionIDNameMap = map[iamtypes.ActionID]string{
	iamtypes.EditBusinessHost:                    "业务主机编辑",
	iamtypes.BusinessHostTransferToResourcePool:  "主机归还主机池",
	iamtypes.HostTransferAcrossBusiness:          "主机转移到其他业务",
	iamtypes.CreateBusinessCustomQuery:           "动态分组新建",
	iamtypes.EditBusinessCustomQuery:             "动态分组编辑",
	iamtypes.DeleteBusinessCustomQuery:           "动态分组删除",
	iamtypes.EditBusinessCustomField:             "业务自定义字段编辑",
	iamtypes.CreateBusinessServiceCategory:       "服务分类新建",
	iamtypes.EditBusinessServiceCategory:         "服务分类编辑",
	iamtypes.DeleteBusinessServiceCategory:       "服务分类删除",
	iamtypes.CreateBusinessServiceInstance:       "服务实例新建",
	iamtypes.EditBusinessServiceInstance:         "服务实例编辑",
	iamtypes.DeleteBusinessServiceInstance:       "服务实例删除",
	iamtypes.CreateBusinessServiceTemplate:       "服务模板新建",
	iamtypes.EditBusinessServiceTemplate:         "服务模板编辑",
	iamtypes.DeleteBusinessServiceTemplate:       "服务模板删除",
	iamtypes.CreateBusinessSetTemplate:           "集群模板新建",
	iamtypes.EditBusinessSetTemplate:             "集群模板编辑",
	iamtypes.DeleteBusinessSetTemplate:           "集群模板删除",
	iamtypes.CreateBusinessTopology:              "业务拓扑新建",
	iamtypes.EditBusinessTopology:                "业务拓扑编辑",
	iamtypes.DeleteBusinessTopology:              "业务拓扑删除",
	iamtypes.EditBusinessHostApply:               "主机自动应用编辑",
	iamtypes.ViewResourcePoolHost:                "主机池主机查看",
	iamtypes.CreateResourcePoolHost:              "主机池主机创建",
	iamtypes.EditResourcePoolHost:                "主机池主机编辑",
	iamtypes.DeleteResourcePoolHost:              "主机池主机删除",
	iamtypes.ResourcePoolHostTransferToBusiness:  "主机池主机分配到业务",
	iamtypes.ResourcePoolHostTransferToDirectory: "主机池主机分配到目录",
	iamtypes.CreateResourcePoolDirectory:         "主机池目录创建",
	iamtypes.EditResourcePoolDirectory:           "主机池目录编辑",
	iamtypes.DeleteResourcePoolDirectory:         "主机池目录删除",
	iamtypes.CreateBusiness:                      "业务创建",
	iamtypes.EditBusiness:                        "业务编辑",
	iamtypes.ArchiveBusiness:                     "业务归档",
	iamtypes.FindBusiness:                        "业务查询",
	iamtypes.ViewBusinessResource:                "业务访问",
	iamtypes.CreateBizSet:                        "业务集新增",
	iamtypes.EditBizSet:                          "业务集编辑",
	iamtypes.DeleteBizSet:                        "业务集删除",
	iamtypes.ViewBizSet:                          "业务集查看",
	iamtypes.AccessBizSet:                        "业务集访问",
	iamtypes.CreateProject:                       "项目新建",
	iamtypes.EditProject:                         "项目编辑",
	iamtypes.DeleteProject:                       "项目删除",
	iamtypes.ViewProject:                         "项目查看",
	iamtypes.ViewCloudArea:                       "管控区域查看",
	iamtypes.CreateCloudArea:                     "管控区域创建",
	iamtypes.EditCloudArea:                       "管控区域编辑",
	iamtypes.DeleteCloudArea:                     "管控区域删除",
	iamtypes.ViewSysModel:                        "模型查看",
	iamtypes.CreateSysModel:                      "模型新建",
	iamtypes.EditSysModel:                        "模型编辑",
	iamtypes.DeleteSysModel:                      "模型删除",
	iamtypes.CreateAssociationType:               "关联类型新建",
	iamtypes.EditAssociationType:                 "关联类型编辑",
	iamtypes.DeleteAssociationType:               "关联类型删除",
	iamtypes.CreateModelGroup:                    "模型分组新建",
	iamtypes.EditModelGroup:                      "模型分组编辑",
	iamtypes.DeleteModelGroup:                    "模型分组删除",
	iamtypes.ViewModelTopo:                       "模型拓扑查看",
	iamtypes.EditBusinessLayer:                   "业务层级编辑",
	iamtypes.EditModelTopologyView:               "模型拓扑视图编辑",
	iamtypes.FindAuditLog:                        "操作审计查询",
	iamtypes.WatchHostEvent:                      "主机事件监听",
	iamtypes.WatchHostRelationEvent:              "主机关系事件监听",
	iamtypes.WatchBizEvent:                       "业务事件监听",
	iamtypes.WatchSetEvent:                       "集群事件监听",
	iamtypes.WatchModuleEvent:                    "模块数据监听",
	iamtypes.WatchProcessEvent:                   "进程数据监听",
	iamtypes.WatchCommonInstanceEvent:            "模型实例事件监听",
	iamtypes.WatchMainlineInstanceEvent:          "自定义拓扑层级事件监听",
	iamtypes.WatchInstAsstEvent:                  "实例关联事件监听",
	iamtypes.WatchBizSetEvent:                    "业务集事件监听",
	iamtypes.WatchPlatEvent:                      "管控区域事件监听",
	iamtypes.WatchKubeClusterEvent:               "容器集群事件监听",
	iamtypes.WatchKubeNodeEvent:                  "容器节点事件监听",
	iamtypes.WatchKubeNamespaceEvent:             "容器命名空间事件监听",
	iamtypes.WatchKubeWorkloadEvent:              "容器工作负载事件监听",
	iamtypes.WatchKubePodEvent:                   "容器Pod事件监听",
	iamtypes.WatchProjectEvent:                   "项目事件监听",
	iamtypes.GlobalSettings:                      "全局设置",
	iamtypes.ManageHostAgentID:                   "主机AgentID管理",
	iamtypes.CreateContainerCluster:              "容器集群新建",
	iamtypes.EditContainerCluster:                "容器集群编辑",
	iamtypes.DeleteContainerCluster:              "容器集群删除",
	iamtypes.CreateContainerNode:                 "容器集群节点新建",
	iamtypes.EditContainerNode:                   "容器集群节点编辑",
	iamtypes.DeleteContainerNode:                 "容器集群节点删除",
	iamtypes.CreateContainerNamespace:            "容器命名空间新建",
	iamtypes.EditContainerNamespace:              "容器命名空间编辑",
	iamtypes.DeleteContainerNamespace:            "容器命名空间删除",
	iamtypes.CreateContainerWorkload:             "容器工作负载新建",
	iamtypes.EditContainerWorkload:               "容器工作负载编辑",
	iamtypes.DeleteContainerWorkload:             "容器工作负载删除",
	iamtypes.CreateContainerPod:                  "容器Pod新建",
	iamtypes.DeleteContainerPod:                  "容器Pod删除",
	iamtypes.UseFulltextSearch:                   "全文检索",
	iamtypes.CreateFieldGroupingTemplate:         "字段组合模板新建",
	iamtypes.ViewFieldGroupingTemplate:           "字段组合模板查看",
	iamtypes.EditFieldGroupingTemplate:           "字段组合模板编辑",
	iamtypes.DeleteFieldGroupingTemplate:         "字段组合模板删除",
	iamtypes.EditIDRuleIncrID:                    "ID规则自增ID编辑",
	iamtypes.CreateFullSyncCond:                  "全量同步缓存条件新建",
	iamtypes.ViewFullSyncCond:                    "全量同步缓存条件查看",
	iamtypes.EditFullSyncCond:                    "全量同步缓存条件编辑",
	iamtypes.DeleteFullSyncCond:                  "全量同步缓存条件删除",
	iamtypes.ViewGeneralCache:                    "通用缓存查询",
	iamtypes.ViewTenantSet:                       "租户集查看",
	iamtypes.AccessTenantSet:                     "租户集访问",
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
	resourceActionList = append(resourceActionList, genKubeEventWatchActions()...)
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
	hostSelection := []iam.RelatedInstanceSelection{{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.BizHostInstanceSelection,
	}}

	relatedResource := []iam.RelateResourceType{{
		SystemID:    iamtypes.SystemIDCMDB,
		ID:          iamtypes.Host,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
		SelectionMode:      iamtypes.ModeAll,
		InstanceSelections: hostSelection,
	}}

	actions := make([]iam.ResourceAction, 0)

	// edit business's host actions
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessHost,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessHost],
		NameEn:               "Edit Business Hosts",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	// business host transfer to resource pool actions
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.BusinessHostTransferToResourcePool,
		Name:                 ActionIDNameMap[iamtypes.BusinessHostTransferToResourcePool],
		NameEn:               "Return Hosts To Pool",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource, resourcePoolDirResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	businessHostResource := iam.RelateResourceType{
		SystemID:    iamtypes.SystemIDCMDB,
		ID:          iamtypes.BusinessForHostTrans,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.BusinessHostTransferSelection,
		}},
	}

	// business host transfer to another business actions
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.HostTransferAcrossBusiness,
		Name:                 ActionIDNameMap[iamtypes.HostTransferAcrossBusiness],
		NameEn:               "Assigned Host To Other Business",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessHostResource, businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessCustomQueryActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.BizCustomQuerySelection,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.BizCustomQuery,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusinessCustomQuery,
		Name:                 ActionIDNameMap[iamtypes.CreateBusinessCustomQuery],
		NameEn:               "Create Dynamic Grouping",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessCustomQuery,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessCustomQuery],
		NameEn:               "Edit Dynamic Grouping",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBusinessCustomQuery,
		Name:                 ActionIDNameMap[iamtypes.DeleteBusinessCustomQuery],
		NameEn:               "Delete Dynamic Grouping",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessCustomFieldActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessCustomField,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessCustomField],
		NameEn:               "Edit Custom Field",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceCategoryActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusinessServiceCategory,
		Name:                 ActionIDNameMap[iamtypes.CreateBusinessServiceCategory],
		NameEn:               "Create Service Category",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessServiceCategory,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessServiceCategory],
		NameEn:               "Edit Service Category",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBusinessServiceCategory,
		Name:                 ActionIDNameMap[iamtypes.DeleteBusinessServiceCategory],
		NameEn:               "Delete Service Category",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceInstanceActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusinessServiceInstance,
		Name:                 ActionIDNameMap[iamtypes.CreateBusinessServiceInstance],
		NameEn:               "Create Service Instance",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessServiceInstance,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessServiceInstance],
		NameEn:               "Edit Service Instance",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBusinessServiceInstance,
		Name:                 ActionIDNameMap[iamtypes.DeleteBusinessServiceInstance],
		NameEn:               "Delete Service Instance",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceTemplateActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID:       iamtypes.SystemIDCMDB,
		ID:             iamtypes.BizProcessServiceTemplateSelection,
		IgnoreAuthPath: true,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.BizProcessServiceTemplate,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusinessServiceTemplate,
		Name:                 ActionIDNameMap[iamtypes.CreateBusinessServiceTemplate],
		NameEn:               "Create Service Template",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessServiceTemplate,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessServiceTemplate],
		NameEn:               "Edit Service Template",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBusinessServiceTemplate,
		Name:                 ActionIDNameMap[iamtypes.DeleteBusinessServiceTemplate],
		NameEn:               "Delete Service Template",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessSetTemplateActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID:       iamtypes.SystemIDCMDB,
		ID:             iamtypes.BizSetTemplateSelection,
		IgnoreAuthPath: true,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.BizSetTemplate,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusinessSetTemplate,
		Name:                 ActionIDNameMap[iamtypes.CreateBusinessSetTemplate],
		NameEn:               "Create Set Template",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessSetTemplate,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessSetTemplate],
		NameEn:               "Edit Set Template",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBusinessSetTemplate,
		Name:                 ActionIDNameMap[iamtypes.DeleteBusinessSetTemplate],
		NameEn:               "Delete Set Template",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessTopologyActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusinessTopology,
		Name:                 ActionIDNameMap[iamtypes.CreateBusinessTopology],
		NameEn:               "Create Business Topo",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessTopology,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessTopology],
		NameEn:               "Edit Business Topo",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBusinessTopology,
		Name:                 ActionIDNameMap[iamtypes.DeleteBusinessTopology],
		NameEn:               "Delete Business Topo",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessHostApplyActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessHostApply,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessHostApply],
		NameEn:               "Edit Host Apply",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genResourcePoolHostActions() []iam.ResourceAction {
	hostSelection := []iam.RelatedInstanceSelection{{SystemID: iamtypes.SystemIDCMDB,
		ID: iamtypes.SysHostInstanceSelection}}

	relatedResource := []iam.RelateResourceType{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.Host, NameAlias: "",
		NameAliasEn: "", Scope: nil,
		// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
		SelectionMode: iamtypes.ModeAll, InstanceSelections: hostSelection,
	}}

	actions := make([]iam.ResourceAction, 0)

	actions = append(actions, iam.ResourceAction{
		ID:   iamtypes.ViewResourcePoolHost,
		Name: ActionIDNameMap[iamtypes.ViewResourcePoolHost], NameEn: "View Resource Pool Hosts",
		Type: iamtypes.View, RelatedResourceTypes: nil,
		RelatedActions: nil, Version: 1,
	})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.CreateResourcePoolHost, Name: ActionIDNameMap[iamtypes.CreateResourcePoolHost],
		NameEn: "Create Pool Hosts", Type: iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil, Version: 1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:   iamtypes.EditResourcePoolHost,
		Name: ActionIDNameMap[iamtypes.EditResourcePoolHost], NameEn: "Edit Pool Hosts",
		Type: iamtypes.Edit, RelatedResourceTypes: relatedResource,
		RelatedActions: []iamtypes.ActionID{iamtypes.ViewResourcePoolHost}, Version: 1,
	})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.DeleteResourcePoolHost, Name: ActionIDNameMap[iamtypes.DeleteResourcePoolHost],
		NameEn: "Delete Pool Hosts",
		Type:   iamtypes.Delete, RelatedResourceTypes: relatedResource,
		RelatedActions: []iamtypes.ActionID{iamtypes.ViewResourcePoolHost}, Version: 1,
	})

	relatedHostResource := []iam.RelateResourceType{{
		SystemID:  iamtypes.SystemIDCMDB,
		ID:        iamtypes.SysHostRscPoolDirectory,
		NameAlias: "", NameAliasEn: "", Scope: nil,
		InstanceSelections: []iam.RelatedInstanceSelection{{SystemID: iamtypes.SystemIDCMDB,
			ID: iamtypes.SysHostRscPoolDirectorySelection}},
	}}

	transferToBusinessRelatedResource := append(relatedHostResource, businessResource)
	actions = append(actions, iam.ResourceAction{
		ID:     iamtypes.ResourcePoolHostTransferToBusiness,
		Name:   ActionIDNameMap[iamtypes.ResourcePoolHostTransferToBusiness],
		NameEn: "Assigned Pool Hosts To Business", Type: iamtypes.Edit,
		RelatedResourceTypes: transferToBusinessRelatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewResourcePoolHost}, Version: 1,
	})

	transferToDirectoryRelatedResource := append(relatedHostResource, resourcePoolDirResource)
	actions = append(actions, iam.ResourceAction{
		ID:     iamtypes.ResourcePoolHostTransferToDirectory,
		Name:   ActionIDNameMap[iamtypes.ResourcePoolHostTransferToDirectory],
		NameEn: "Assigned Pool Hosts To Directory", Type: iamtypes.Edit,
		RelatedResourceTypes: transferToDirectoryRelatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewResourcePoolHost}, Version: 1,
	})

	actions = append(actions,
		iam.ResourceAction{ID: iamtypes.ManageHostAgentID, Name: ActionIDNameMap[iamtypes.ManageHostAgentID],
			NameEn: "Manage Host AgentID", Type: iamtypes.Edit, Version: 1})

	return actions
}

func genResourcePoolDirectoryActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateResourcePoolDirectory,
		Name:                 ActionIDNameMap[iamtypes.CreateResourcePoolDirectory],
		NameEn:               "Create Pool Directory",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditResourcePoolDirectory,
		Name:                 ActionIDNameMap[iamtypes.EditResourcePoolDirectory],
		NameEn:               "Edit Pool Directory",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteResourcePoolDirectory,
		Name:                 ActionIDNameMap[iamtypes.DeleteResourcePoolDirectory],
		NameEn:               "Delete Pool Directory",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: []iam.RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBusiness,
		Name:                 ActionIDNameMap[iamtypes.CreateBusiness],
		NameEn:               "Create Business",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusiness,
		Name:                 ActionIDNameMap[iamtypes.EditBusiness],
		NameEn:               "Edit Business",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.FindBusiness},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.ArchiveBusiness,
		Name:                 ActionIDNameMap[iamtypes.ArchiveBusiness],
		NameEn:               "Archive Business",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.FindBusiness},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.FindBusiness,
		Name:                 ActionIDNameMap[iamtypes.FindBusiness],
		NameEn:               "View Business",
		Type:                 iamtypes.View,
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:     iamtypes.ViewBusinessResource,
		Name:   ActionIDNameMap[iamtypes.ViewBusinessResource],
		NameEn: "View Business Resource",
		Type:   iamtypes.View,
		// TODO add business collection resource
		RelatedResourceTypes: []iam.RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBizSetActions() []iam.ResourceAction {
	bizSetResource := iam.RelateResourceType{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.BizSet,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.BizSetSelection,
		}},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateBizSet,
		Name:                 ActionIDNameMap[iamtypes.CreateBizSet],
		NameEn:               "Create Business Set",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBizSet},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBizSet,
		Name:                 ActionIDNameMap[iamtypes.EditBizSet],
		NameEn:               "Edit Business Set",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{bizSetResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBizSet},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteBizSet,
		Name:                 ActionIDNameMap[iamtypes.DeleteBizSet],
		NameEn:               "Delete Business Set",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: []iam.RelateResourceType{bizSetResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewBizSet},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.ViewBizSet,
		Name:                 ActionIDNameMap[iamtypes.ViewBizSet],
		NameEn:               "View Business Set",
		Type:                 iamtypes.View,
		RelatedResourceTypes: []iam.RelateResourceType{bizSetResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.AccessBizSet,
		Name:                 ActionIDNameMap[iamtypes.AccessBizSet],
		NameEn:               "Access Business Set",
		Type:                 iamtypes.View,
		RelatedResourceTypes: []iam.RelateResourceType{bizSetResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genProjectActions() []iam.ResourceAction {
	projectResource := iam.RelateResourceType{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.Project,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.ProjectSelection,
		}},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateProject,
		Name:                 ActionIDNameMap[iamtypes.CreateProject],
		NameEn:               "Create Project",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditProject,
		Name:                 ActionIDNameMap[iamtypes.EditProject],
		NameEn:               "Edit Project",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: []iam.RelateResourceType{projectResource},
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewProject},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteProject,
		Name:                 ActionIDNameMap[iamtypes.DeleteProject],
		NameEn:               "Delete Project",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: []iam.RelateResourceType{projectResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.ViewProject,
		Name:                 ActionIDNameMap[iamtypes.ViewProject],
		NameEn:               "View Project",
		Type:                 iamtypes.View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudAreaActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.SysCloudAreaSelection,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.SysCloudArea,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.ViewCloudArea,
		Name:                 ActionIDNameMap[iamtypes.ViewCloudArea],
		NameEn:               "View Cloud Area",
		Type:                 iamtypes.View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateCloudArea,
		Name:                 ActionIDNameMap[iamtypes.CreateCloudArea],
		NameEn:               "Create Cloud Area",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditCloudArea,
		Name:                 ActionIDNameMap[iamtypes.EditCloudArea],
		NameEn:               "Edit Cloud Area",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewCloudArea},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteCloudArea,
		Name:                 ActionIDNameMap[iamtypes.DeleteCloudArea],
		NameEn:               "Delete Cloud Area",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewCloudArea},
		Version:              1,
	})

	return actions
}

func genModelActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.SysModelSelection,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.SysModel,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.ViewSysModel,
		Name:                 ActionIDNameMap[iamtypes.ViewSysModel],
		NameEn:               "View Model",
		Type:                 iamtypes.View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:     iamtypes.CreateSysModel,
		Name:   ActionIDNameMap[iamtypes.CreateSysModel],
		NameEn: "Create Model",
		Type:   iamtypes.Create,
		RelatedResourceTypes: []iam.RelateResourceType{
			{
				SystemID:    iamtypes.SystemIDCMDB,
				ID:          iamtypes.SysModelGroup,
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				InstanceSelections: []iam.RelatedInstanceSelection{{
					SystemID: iamtypes.SystemIDCMDB,
					ID:       iamtypes.SysModelGroupSelection,
				}},
			},
		},
		RelatedActions: nil,
		Version:        1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditSysModel,
		Name:                 ActionIDNameMap[iamtypes.EditSysModel],
		NameEn:               "Edit Model",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewSysModel},
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteSysModel,
		Name:                 ActionIDNameMap[iamtypes.DeleteSysModel],
		NameEn:               "Delete Model",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewSysModel},
		Version:              1,
	})

	return actions
}

func genAssociationTypeActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.SysAssociationTypeSelection,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.SysAssociationType,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateAssociationType,
		Name:                 ActionIDNameMap[iamtypes.CreateAssociationType],
		NameEn:               "Create Association Type",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditAssociationType,
		Name:                 ActionIDNameMap[iamtypes.EditAssociationType],
		NameEn:               "Edit Association Type",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteAssociationType,
		Name:                 ActionIDNameMap[iamtypes.DeleteAssociationType],
		NameEn:               "Delete Association Type",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelGroupActions() []iam.ResourceAction {
	selection := []iam.RelatedInstanceSelection{{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.SysModelGroupSelection,
	}}

	relatedResource := []iam.RelateResourceType{
		{
			SystemID:           iamtypes.SystemIDCMDB,
			ID:                 iamtypes.SysModelGroup,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.CreateModelGroup,
		Name:                 ActionIDNameMap[iamtypes.CreateModelGroup],
		NameEn:               "Create Model Group",
		Type:                 iamtypes.Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditModelGroup,
		Name:                 ActionIDNameMap[iamtypes.EditModelGroup],
		NameEn:               "Edit Model Group",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.DeleteModelGroup,
		Name:                 ActionIDNameMap[iamtypes.DeleteModelGroup],
		NameEn:               "Delete Model Group",
		Type:                 iamtypes.Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessLayerActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditBusinessLayer,
		Name:                 ActionIDNameMap[iamtypes.EditBusinessLayer],
		NameEn:               "Edit Business Level",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewModelTopo},
		Version:              1,
	})
	return actions
}

func genModelTopologyViewActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:      iamtypes.ViewModelTopo,
		Name:    ActionIDNameMap[iamtypes.ViewModelTopo],
		NameEn:  "View Model Topo",
		Type:    iamtypes.View,
		Version: 1,
	})

	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.EditModelTopologyView,
		Name:                 ActionIDNameMap[iamtypes.EditModelTopologyView],
		NameEn:               "Edit Model Topo View",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       []iamtypes.ActionID{iamtypes.ViewModelTopo},
		Version:              1,
	})
	return actions
}

func genAuditLogActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.FindAuditLog,
		Name:                 ActionIDNameMap[iamtypes.FindAuditLog],
		NameEn:               "View Operation Audit",
		Type:                 iamtypes.View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genEventWatchActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{ID: iamtypes.WatchHostEvent,
		Name: ActionIDNameMap[iamtypes.WatchHostEvent], NameEn: "Host Event Listen",
		Type: iamtypes.View, RelatedResourceTypes: nil, RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{ID: iamtypes.WatchHostRelationEvent,
		Name:   ActionIDNameMap[iamtypes.WatchHostRelationEvent],
		NameEn: "Host Relation Event Listen", Type: iamtypes.View,
		RelatedResourceTypes: nil, RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{ID: iamtypes.WatchBizEvent,
		Name: ActionIDNameMap[iamtypes.WatchBizEvent], NameEn: "Business Event Listen",
		Type: iamtypes.View, RelatedResourceTypes: nil, RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchSetEvent, Name: ActionIDNameMap[iamtypes.WatchSetEvent],
		NameEn: "Set Event Listen", Type: iamtypes.View, RelatedResourceTypes: nil,
		RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchModuleEvent, Name: ActionIDNameMap[iamtypes.WatchModuleEvent],
		NameEn: "Module Event Listen", Type: iamtypes.View,
		RelatedResourceTypes: nil, RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchProcessEvent, Name: ActionIDNameMap[iamtypes.WatchProcessEvent],
		NameEn: "Process Event Listen", Type: iamtypes.View,
		RelatedResourceTypes: nil, RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchBizSetEvent, Name: ActionIDNameMap[iamtypes.WatchBizSetEvent],
		NameEn: "Business Set Event Listen", Type: iamtypes.View, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchPlatEvent, Name: ActionIDNameMap[iamtypes.WatchPlatEvent],
		NameEn: "Cloud Area Event Listen", Type: iamtypes.View, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchProjectEvent, Name: ActionIDNameMap[iamtypes.WatchProjectEvent],
		NameEn: "Project Event Listen", Type: iamtypes.View, Version: 1})

	modelSelection := []iam.RelatedInstanceSelection{{SystemID: iamtypes.SystemIDCMDB,
		ID: iamtypes.SysModelEventSelection}}

	modelResource := []iam.RelateResourceType{
		{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysModelEvent, NameAlias: "",
			NameAliasEn: "", Scope: nil, InstanceSelections: modelSelection}}

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchCommonInstanceEvent, Name: ActionIDNameMap[iamtypes.WatchCommonInstanceEvent],
		NameEn: "Common Model Instance Event Listen", Type: iamtypes.View,
		RelatedResourceTypes: modelResource, RelatedActions: nil, Version: 1})

	mainlineModelSelection := []iam.RelatedInstanceSelection{{SystemID: iamtypes.SystemIDCMDB,
		ID: iamtypes.MainlineModelEventSelection}}

	mainlineModelResource := []iam.RelateResourceType{
		{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.MainlineModelEvent, NameAlias: "",
			NameAliasEn: "", Scope: nil, InstanceSelections: mainlineModelSelection},
	}

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchMainlineInstanceEvent, Name: ActionIDNameMap[iamtypes.WatchMainlineInstanceEvent],
		NameEn: "Custom Topo Layer Event Listen", Type: iamtypes.View,
		RelatedResourceTypes: mainlineModelResource, RelatedActions: nil, Version: 1})

	actions = append(actions, iam.ResourceAction{
		ID: iamtypes.WatchInstAsstEvent, Name: ActionIDNameMap[iamtypes.WatchInstAsstEvent],
		NameEn: "Instance Association Event Listen", Type: iamtypes.View,
		RelatedResourceTypes: []iam.RelateResourceType{
			{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.InstAsstEvent, NameAlias: "", NameAliasEn: "", Scope: nil,
				InstanceSelections: []iam.RelatedInstanceSelection{{SystemID: iamtypes.SystemIDCMDB,
					ID: iamtypes.InstAsstEventSelection}}}}, RelatedActions: nil, Version: 1})
	return actions
}

func genKubeEventWatchActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.WatchKubeClusterEvent,
			Name:    ActionIDNameMap[iamtypes.WatchKubeClusterEvent],
			NameEn:  "Kube Cluster Event Listen",
			Type:    iamtypes.View,
			Version: 1,
		},
		{
			ID:      iamtypes.WatchKubeNodeEvent,
			Name:    ActionIDNameMap[iamtypes.WatchKubeNodeEvent],
			NameEn:  "Kube Node Event Listen",
			Type:    iamtypes.View,
			Version: 1,
		},
		{
			ID:      iamtypes.WatchKubeNamespaceEvent,
			Name:    ActionIDNameMap[iamtypes.WatchKubeNamespaceEvent],
			NameEn:  "Kube Namespace Event Listen",
			Type:    iamtypes.View,
			Version: 1,
		},
		{
			ID:     iamtypes.WatchKubeWorkloadEvent,
			Name:   ActionIDNameMap[iamtypes.WatchKubeWorkloadEvent],
			NameEn: "Kube Workload Event Listen",
			Type:   iamtypes.View,
			RelatedResourceTypes: []iam.RelateResourceType{
				{
					SystemID: iamtypes.SystemIDCMDB,
					ID:       iamtypes.KubeWorkloadEvent,
					InstanceSelections: []iam.RelatedInstanceSelection{{
						SystemID: iamtypes.SystemIDCMDB,
						ID:       iamtypes.KubeWorkloadEventSelection,
					}},
				},
			},
			Version: 1,
		},
		{
			ID:      iamtypes.WatchKubePodEvent,
			Name:    ActionIDNameMap[iamtypes.WatchKubePodEvent],
			NameEn:  "Kube Pod Event Listen",
			Type:    iamtypes.View,
			Version: 1,
		},
	}
}

func genConfigAdminActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:                   iamtypes.GlobalSettings,
		Name:                 ActionIDNameMap[iamtypes.GlobalSettings],
		NameEn:               "Global Settings",
		Type:                 iamtypes.Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
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
			ID:      iamtypes.CreateContainerCluster,
			Name:    ActionIDNameMap[iamtypes.CreateContainerCluster],
			NameEn:  "Create Container Cluster",
			Type:    iamtypes.Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.EditContainerCluster,
			Name:    ActionIDNameMap[iamtypes.EditContainerCluster],
			NameEn:  "Edit Container Cluster",
			Type:    iamtypes.Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.DeleteContainerCluster,
			Name:    ActionIDNameMap[iamtypes.DeleteContainerCluster],
			NameEn:  "Delete Container Cluster",
			Type:    iamtypes.Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerNodeActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.CreateContainerNode,
			Name:    ActionIDNameMap[iamtypes.CreateContainerNode],
			NameEn:  "Create Container Node",
			Type:    iamtypes.Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.EditContainerNode,
			Name:    ActionIDNameMap[iamtypes.EditContainerNode],
			NameEn:  "Edit Container Node",
			Type:    iamtypes.Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.DeleteContainerNode,
			Name:    ActionIDNameMap[iamtypes.DeleteContainerNode],
			NameEn:  "Delete Container Node",
			Type:    iamtypes.Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerNamespaceActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.CreateContainerNamespace,
			Name:    ActionIDNameMap[iamtypes.CreateContainerNamespace],
			NameEn:  "Create Container Namespace",
			Type:    iamtypes.Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.EditContainerNamespace,
			Name:    ActionIDNameMap[iamtypes.EditContainerNamespace],
			NameEn:  "Edit Container Namespace",
			Type:    iamtypes.Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.DeleteContainerNamespace,
			Name:    ActionIDNameMap[iamtypes.DeleteContainerNamespace],
			NameEn:  "Delete Container Namespace",
			Type:    iamtypes.Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerWorkloadActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.CreateContainerWorkload,
			Name:    ActionIDNameMap[iamtypes.CreateContainerWorkload],
			NameEn:  "Create Container Workload",
			Type:    iamtypes.Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.EditContainerWorkload,
			Name:    ActionIDNameMap[iamtypes.EditContainerWorkload],
			NameEn:  "Edit Container Workload",
			Type:    iamtypes.Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.DeleteContainerWorkload,
			Name:    ActionIDNameMap[iamtypes.DeleteContainerWorkload],
			NameEn:  "Delete Container Workload",
			Type:    iamtypes.Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genContainerPodActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.CreateContainerPod,
			Name:    ActionIDNameMap[iamtypes.CreateContainerPod],
			NameEn:  "Create Container Pod",
			Type:    iamtypes.Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.DeleteContainerPod,
			Name:    ActionIDNameMap[iamtypes.DeleteContainerPod],
			NameEn:  "Delete Container Pod",
			Type:    iamtypes.Delete,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genFulltextSearchActions() []iam.ResourceAction {
	actions := make([]iam.ResourceAction, 0)
	actions = append(actions, iam.ResourceAction{
		ID:      iamtypes.UseFulltextSearch,
		Name:    ActionIDNameMap[iamtypes.UseFulltextSearch],
		NameEn:  "Fulltext Search",
		Type:    iamtypes.View,
		Version: 1,
	})
	return actions
}

func genFieldGroupingTemplateActions() []iam.ResourceAction {
	templateResource := iam.RelateResourceType{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.FieldGroupingTemplate,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.FieldGroupingTemplateSelection,
		}},
	}

	return []iam.ResourceAction{
		{
			ID:      iamtypes.CreateFieldGroupingTemplate,
			Name:    ActionIDNameMap[iamtypes.CreateFieldGroupingTemplate],
			NameEn:  "Create Field Grouping Template",
			Type:    iamtypes.Create,
			Version: 1,
		},
		{
			ID:                   iamtypes.ViewFieldGroupingTemplate,
			Name:                 ActionIDNameMap[iamtypes.ViewFieldGroupingTemplate],
			NameEn:               "View Field Grouping Template",
			Type:                 iamtypes.View,
			RelatedResourceTypes: []iam.RelateResourceType{templateResource},
			Version:              1,
		},
		{
			ID:                   iamtypes.EditFieldGroupingTemplate,
			Name:                 ActionIDNameMap[iamtypes.EditFieldGroupingTemplate],
			NameEn:               "Edit Field Grouping Template",
			Type:                 iamtypes.Edit,
			RelatedResourceTypes: []iam.RelateResourceType{templateResource},
			RelatedActions:       []iamtypes.ActionID{iamtypes.ViewFieldGroupingTemplate},
			Version:              1,
		},
		{
			ID:                   iamtypes.DeleteFieldGroupingTemplate,
			Name:                 ActionIDNameMap[iamtypes.DeleteFieldGroupingTemplate],
			NameEn:               "Delete Field Grouping Template",
			Type:                 iamtypes.Delete,
			RelatedResourceTypes: []iam.RelateResourceType{templateResource},
			RelatedActions:       []iamtypes.ActionID{iamtypes.ViewFieldGroupingTemplate},
			Version:              1,
		},
	}
}

func genIDRuleActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.EditIDRuleIncrID,
			Name:    ActionIDNameMap[iamtypes.EditIDRuleIncrID],
			NameEn:  "Edit ID Rule",
			Type:    iamtypes.Edit,
			Version: 1,
		},
	}
}

func genFullSyncCondActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:      iamtypes.CreateFullSyncCond,
			Name:    ActionIDNameMap[iamtypes.CreateFullSyncCond],
			NameEn:  "Create Full Sync Cond",
			Type:    iamtypes.Create,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.EditFullSyncCond,
			Name:    ActionIDNameMap[iamtypes.EditFullSyncCond],
			NameEn:  "Edit Full Sync Cond",
			Type:    iamtypes.Edit,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.DeleteFullSyncCond,
			Name:    ActionIDNameMap[iamtypes.DeleteFullSyncCond],
			NameEn:  "Delete Full Sync Cond",
			Type:    iamtypes.Delete,
			Version: 1,
			Hidden:  true,
		},
		{
			ID:      iamtypes.ViewFullSyncCond,
			Name:    ActionIDNameMap[iamtypes.ViewFullSyncCond],
			NameEn:  "View Full Sync Cond",
			Type:    iamtypes.View,
			Version: 1,
			Hidden:  true,
		},
	}
}

func genCacheActions() []iam.ResourceAction {
	return []iam.ResourceAction{
		{
			ID:     iamtypes.ViewGeneralCache,
			Name:   ActionIDNameMap[iamtypes.ViewGeneralCache],
			NameEn: "View General Resource Cache",
			Type:   iamtypes.View,
			RelatedResourceTypes: []iam.RelateResourceType{{
				SystemID: iamtypes.SystemIDCMDB,
				ID:       iamtypes.GeneralCache,
				InstanceSelections: []iam.RelatedInstanceSelection{{
					SystemID: iamtypes.SystemIDCMDB,
					ID:       iamtypes.GeneralCacheSelection,
				}},
			}},
			Version: 1,
			Hidden:  true,
		},
	}
}

func genTenantSetActions() []iam.ResourceAction {
	if tools.GetDefaultTenant() != common.BKDefaultTenantID {
		return make([]iam.ResourceAction, 0)
	}

	tenantSetResource := iam.RelateResourceType{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.TenantSet,
		InstanceSelections: []iam.RelatedInstanceSelection{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       iamtypes.TenantSetSelection,
		}},
	}

	return []iam.ResourceAction{
		{
			ID:                   iamtypes.ViewTenantSet,
			Name:                 ActionIDNameMap[iamtypes.ViewTenantSet],
			NameEn:               "View Tenant Set",
			Type:                 iamtypes.View,
			RelatedResourceTypes: []iam.RelateResourceType{tenantSetResource},
			Version:              1,
			Hidden:               true,
			TenantID:             common.BKDefaultTenantID,
		}, {
			ID:                   iamtypes.AccessTenantSet,
			Name:                 ActionIDNameMap[iamtypes.AccessTenantSet],
			NameEn:               "Access Tenant Set",
			Type:                 iamtypes.View,
			RelatedResourceTypes: []iam.RelateResourceType{tenantSetResource},
			Version:              1,
			Hidden:               true,
			TenantID:             common.BKDefaultTenantID,
		},
	}
}
