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

// GenerateActionGroups generate all the resource action groups registered to IAM.
func GenerateActionGroups(tenantObjects map[string][]metadata.Object) []iam.ActionGroup {
	ActionGroups := GenerateStaticActionGroups()

	// generate model instance manage action groups, contains model instance related actions which are dynamic
	ActionGroups = append(ActionGroups, GenModelInstanceManageActionGroups(tenantObjects)...)

	return ActionGroups
}

// GenerateStaticActionGroups generate all the static resource action groups.
func GenerateStaticActionGroups() []iam.ActionGroup {
	ActionGroups := make([]iam.ActionGroup, 0)

	// generate business set manage action groups, contains fulltext search related actions
	ActionGroups = append(ActionGroups, genFulltextSearchServiceActionGroups()...)

	// generate business set manage action groups, contains business set related actions
	ActionGroups = append(ActionGroups, genBizSetManageActionGroups()...)

	// generate business manage action groups, contains business related actions
	ActionGroups = append(ActionGroups, genBusinessManageActionGroups()...)

	// generate resource manage action groups, contains resource related actions
	ActionGroups = append(ActionGroups, genResourceManageActionGroups()...)

	// generate model manage action groups, contains model related actions
	ActionGroups = append(ActionGroups, genModelManageActionGroups()...)

	// generate operation statistic action groups, contains operation statistic and audit related actions
	ActionGroups = append(ActionGroups, genOperationStatisticActionGroups()...)

	// generate global settings action groups, contains global settings related actions
	ActionGroups = append(ActionGroups, genGlobalSettingsActionGroups()...)

	// generate container management action groups, contains container related actions
	ActionGroups = append(ActionGroups, genContainerManagementActionGroups()...)

	// generate tenant set action groups
	ActionGroups = append(ActionGroups, genTenantSetActionGroups()...)

	return ActionGroups
}

func genFulltextSearchServiceActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "检索服务",
			NameEn: "Fulltext Search Service",
			Actions: []iam.ActionWithID{
				{
					ID: iamtypes.UseFulltextSearch,
				},
			},
		},
	}
}

func genBizSetManageActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "业务集管理",
			NameEn: "Business Set Manage",
			Actions: []iam.ActionWithID{
				{
					ID: iamtypes.AccessBizSet,
				},
			},
		},
	}
}

func genBusinessManageActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{Name: "业务管理",
			NameEn:  "Business Manage",
			Actions: []iam.ActionWithID{{ID: iamtypes.ViewBusinessResource}},
			SubGroups: []iam.ActionGroup{
				{
					Name:   "业务主机",
					NameEn: "Business Host",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.EditBusinessHost},
						{ID: iamtypes.BusinessHostTransferToResourcePool},
						{ID: iamtypes.HostTransferAcrossBusiness},
					},
				},
				{
					Name:   "业务拓扑",
					NameEn: "Business Topology",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateBusinessTopology}, {ID: iamtypes.EditBusinessTopology},
						{ID: iamtypes.DeleteBusinessTopology},
					},
				},
				{
					Name:   "服务实例",
					NameEn: "Service Instance",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateBusinessServiceInstance}, {ID: iamtypes.EditBusinessServiceInstance},
						{ID: iamtypes.DeleteBusinessServiceInstance},
					},
				},
				{
					Name:   "服务模版",
					NameEn: "Service Template",
					Actions: []iam.ActionWithID{{ID: iamtypes.CreateBusinessServiceTemplate},
						{ID: iamtypes.EditBusinessServiceTemplate}, {ID: iamtypes.DeleteBusinessServiceTemplate},
					},
				},
				{
					Name:   "集群模版",
					NameEn: "Set Template",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateBusinessSetTemplate}, {ID: iamtypes.EditBusinessSetTemplate},
						{ID: iamtypes.DeleteBusinessSetTemplate},
					},
				},
				{
					Name:   "服务分类",
					NameEn: "Service Category",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateBusinessServiceCategory}, {ID: iamtypes.EditBusinessServiceCategory},
						{ID: iamtypes.DeleteBusinessServiceCategory},
					},
				},
				{
					Name:   "动态分组",
					NameEn: "Dynamic Grouping",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateBusinessCustomQuery}, {ID: iamtypes.EditBusinessCustomQuery},
						{ID: iamtypes.DeleteBusinessCustomQuery},
					},
				},
				{
					Name:    "业务自定义字段",
					NameEn:  "Business Custom Field",
					Actions: []iam.ActionWithID{{ID: iamtypes.EditBusinessCustomField}},
				},
				{
					Name:   "主机自动应用",
					NameEn: "Business Host Apply",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.EditBusinessHostApply},
					},
				},
			},
		},
	}
}

func genResourceManageActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "资源管理",
			NameEn: "Resource Manage",
			SubGroups: []iam.ActionGroup{
				{
					Name:   "主机池",
					NameEn: "Host Pool",
					Actions: []iam.ActionWithID{{ID: iamtypes.ViewResourcePoolHost},
						{ID: iamtypes.CreateResourcePoolHost}, {ID: iamtypes.EditResourcePoolHost},
						{ID: iamtypes.DeleteResourcePoolHost}, {ID: iamtypes.ResourcePoolHostTransferToBusiness},
						{ID: iamtypes.ResourcePoolHostTransferToDirectory}, {ID: iamtypes.CreateResourcePoolDirectory},
						{ID: iamtypes.EditResourcePoolDirectory}, {ID: iamtypes.DeleteResourcePoolDirectory},
						{ID: iamtypes.ManageHostAgentID}},
				},
				{
					Name: "业务", NameEn: "Business",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateBusiness}, {ID: iamtypes.EditBusiness}, {ID: iamtypes.ArchiveBusiness},
						{ID: iamtypes.FindBusiness}}},
				{
					Name: "项目", NameEn: "Project",
					Actions: []iam.ActionWithID{{ID: iamtypes.CreateProject}, {ID: iamtypes.EditProject},
						{ID: iamtypes.DeleteProject}, {ID: iamtypes.ViewProject}},
				},
				{
					Name: "业务集", NameEn: "BizSet",
					Actions: []iam.ActionWithID{{ID: iamtypes.CreateBizSet}, {ID: iamtypes.EditBizSet},
						{ID: iamtypes.DeleteBizSet}, {ID: iamtypes.ViewBizSet}},
				},
				{
					Name:   "管控区域",
					NameEn: "Cloud Area",
					Actions: []iam.ActionWithID{{ID: iamtypes.ViewCloudArea}, {ID: iamtypes.CreateCloudArea},
						{ID: iamtypes.EditCloudArea}, {ID: iamtypes.DeleteCloudArea}},
				},
				{
					Name:   "事件监听",
					NameEn: "Event Watch",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.WatchHostEvent}, {ID: iamtypes.WatchHostRelationEvent},
						{ID: iamtypes.WatchBizEvent}, {ID: iamtypes.WatchSetEvent},
						{ID: iamtypes.WatchModuleEvent}, {ID: iamtypes.WatchProcessEvent},
						{ID: iamtypes.WatchCommonInstanceEvent}, {ID: iamtypes.WatchMainlineInstanceEvent},
						{ID: iamtypes.WatchInstAsstEvent}, {ID: iamtypes.WatchBizSetEvent},
						{ID: iamtypes.WatchPlatEvent}, {ID: iamtypes.WatchKubeClusterEvent},
						{ID: iamtypes.WatchKubeNodeEvent}, {ID: iamtypes.WatchKubeNamespaceEvent},
						{ID: iamtypes.WatchKubeWorkloadEvent}, {ID: iamtypes.WatchKubePodEvent},
						{ID: iamtypes.WatchProjectEvent}},
				},
				{
					Name: "全量同步缓存条件", NameEn: "Full Sync Condition",
					Actions: []iam.ActionWithID{{ID: iamtypes.CreateFullSyncCond},
						{ID: iamtypes.ViewFullSyncCond}, {ID: iamtypes.EditFullSyncCond},
						{ID: iamtypes.DeleteFullSyncCond}}},
				{
					Name: "缓存", NameEn: "Cache", Actions: []iam.ActionWithID{{ID: iamtypes.ViewGeneralCache}}},
			},
		},
	}
}

func genModelManageActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "模型管理",
			NameEn: "Model Manage",
			SubGroups: []iam.ActionGroup{
				{
					Name:   "模型分组",
					NameEn: "Model Group",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateModelGroup},
						{ID: iamtypes.EditModelGroup},
						{ID: iamtypes.DeleteModelGroup},
					},
				},
				{
					Name:   "模型关系",
					NameEn: "Model Relation",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.ViewModelTopo,
						},
						{
							ID: iamtypes.EditBusinessLayer,
						},
						{
							ID: iamtypes.EditModelTopologyView,
						},
					},
				},
				{
					Name:   "模型",
					NameEn: "Model",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.ViewSysModel,
						},
						{
							ID: iamtypes.CreateSysModel,
						},
						{
							ID: iamtypes.EditSysModel,
						},
						{
							ID: iamtypes.DeleteSysModel,
						},
					},
				},
				{
					Name:   "关联类型",
					NameEn: "Association Type",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateAssociationType},
						{ID: iamtypes.EditAssociationType},
						{ID: iamtypes.DeleteAssociationType},
					},
				},
				{
					Name:   "字段组合模板",
					NameEn: "Field Grouping Template",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.CreateFieldGroupingTemplate},
						{ID: iamtypes.ViewFieldGroupingTemplate},
						{ID: iamtypes.EditFieldGroupingTemplate},
						{ID: iamtypes.DeleteFieldGroupingTemplate},
					},
				},
				{
					Name:   "ID规则自增ID",
					NameEn: "ID Rule Self-increasing ID",
					Actions: []iam.ActionWithID{
						{ID: iamtypes.EditIDRuleIncrID},
					},
				},
			},
		},
	}
}

// modelInstManageActionGroupName is the name of the model instance management action group
const modelInstManageActionGroupName = "模型实例管理"

// GenModelInstanceManageActionGroups TODO
func GenModelInstanceManageActionGroups(tenantObjects map[string][]metadata.Object) []iam.ActionGroup {
	subGroups := make([]iam.ActionGroup, 0)
	for _, objects := range tenantObjects {
		for _, obj := range objects {
			subGroups = append(subGroups, genDynamicActionSubGroup(obj))
		}
	}

	if len(subGroups) == 0 {
		return make([]iam.ActionGroup, 0)
	}

	return []iam.ActionGroup{
		{
			Name:      modelInstManageActionGroupName,
			NameEn:    "Model instance Manage",
			SubGroups: subGroups,
		},
	}
}

func genContainerManagementActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "容器资源管理",
			NameEn: "Container Management",
			SubGroups: []iam.ActionGroup{
				{
					Name:   "容器 Cluster",
					NameEn: "Container Cluster",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.CreateContainerCluster,
						},
						{
							ID: iamtypes.EditContainerCluster,
						},
						{
							ID: iamtypes.DeleteContainerCluster,
						},
					},
				}, {
					Name:   "容器 Node",
					NameEn: "Container Node",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.CreateContainerNode,
						},
						{
							ID: iamtypes.EditContainerNode,
						},
						{
							ID: iamtypes.DeleteContainerNode,
						},
					},
				}, {
					Name:   "容器命名空间",
					NameEn: "Container Namespace",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.CreateContainerNamespace,
						},
						{
							ID: iamtypes.EditContainerNamespace,
						},
						{
							ID: iamtypes.DeleteContainerNamespace,
						},
					},
				}, {
					Name:   "容器工作负载",
					NameEn: "Container Workload",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.CreateContainerWorkload,
						},
						{
							ID: iamtypes.EditContainerWorkload,
						},
						{
							ID: iamtypes.DeleteContainerWorkload,
						},
					},
				}, {
					Name:   "容器 Pod",
					NameEn: "Container Pod",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.CreateContainerPod,
						},
						{
							ID: iamtypes.DeleteContainerPod,
						},
					},
				},
			},
		},
	}
}

func genOperationStatisticActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "运营统计",
			NameEn: "Operation Statistic",
			SubGroups: []iam.ActionGroup{
				{
					Name:   "操作审计",
					NameEn: "Operation Audit",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.FindAuditLog,
						},
					},
				},
			},
		},
	}
}

func genGlobalSettingsActionGroups() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name:   "全局设置",
			NameEn: "Global Settings",
			SubGroups: []iam.ActionGroup{
				{
					Name:   "全局设置",
					NameEn: "Global Settings",
					Actions: []iam.ActionWithID{
						{
							ID: iamtypes.GlobalSettings,
						},
					},
				},
			},
		},
	}
}

func genTenantSetActionGroups() []iam.ActionGroup {
	if tools.GetDefaultTenant() != common.BKDefaultTenantID {
		return make([]iam.ActionGroup, 0)
	}

	return []iam.ActionGroup{
		{
			Name:   "租户集",
			NameEn: "Tenant Set",
			Actions: []iam.ActionWithID{
				{ID: iamtypes.ViewTenantSet},
				{ID: iamtypes.AccessTenantSet},
			},
		},
	}
}
