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
	"fmt"

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/iam"
)

// GenerateRoles generate all roles registered to IAM, including dynamic model roles
func GenerateRoles(tenantObjects map[string][]metadata.Object) []iam.Role {
	roles := make([]iam.Role, 0)
	roles = append(roles, genCompatibleRoles()...)
	roles = append(roles, genManagerRoles()...)
	roles = append(roles, genCreatorRoles()...)
	roles = append(roles, genViewerRoles()...)
	roles = append(roles, genOwnerRoles()...)
	roles = append(roles, genPlatformRoles()...)
	roles = append(roles, genDynamicInstRoles(tenantObjects)...)
	return roles
}

// genCompatibleRoles generate CommonAction compatible roles
// NOCC:golint/fnsize(角色定义需要放在一起)
func genCompatibleRoles() []iam.Role {
	return []iam.Role{
		{
			ID:          "business_maintainer",
			Name:        "业务运维",
			Description: "业务运维角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewBusinessResource, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessHost, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.CreateBusinessTopology, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessTopology, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.DeleteBusinessTopology, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.CreateBusinessServiceInstance, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessServiceInstance, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.DeleteBusinessServiceInstance, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.CreateBusinessServiceTemplate, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessServiceTemplate, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.DeleteBusinessServiceTemplate, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.CreateBusinessSetTemplate, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessSetTemplate, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.DeleteBusinessSetTemplate, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.CreateBusinessServiceCategory, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessServiceCategory, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.DeleteBusinessServiceCategory, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.CreateBusinessCustomQuery, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessCustomQuery, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.DeleteBusinessCustomQuery, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessCustomField, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.EditBusinessHostApply, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.FindBusiness, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.TransferHostIntoBiz, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.TransferHostOutOfBiz, ResourceTypeID: iamtypes.Business},
			},
		},
		{
			ID:          "business_visitor",
			Name:        "业务只读",
			Description: "业务只读角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewBusinessResource, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.FindBusiness, ResourceTypeID: iamtypes.Business},
			},
		},
		{
			ID:          "biz_set_maintainer",
			Name:        "业务集运维",
			Description: "业务集运维角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.AccessBizSet, ResourceTypeID: iamtypes.BizSet},
				{ID: iamtypes.DeleteBizSet, ResourceTypeID: iamtypes.BizSet},
				{ID: iamtypes.ViewBizSet, ResourceTypeID: iamtypes.BizSet},
			},
		},
		{
			ID:          "biz_set_visitor",
			Name:        "业务集只读",
			Description: "业务集只读角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.AccessBizSet, ResourceTypeID: iamtypes.BizSet},
				{ID: iamtypes.ViewBizSet, ResourceTypeID: iamtypes.BizSet},
			},
		},
		{
			ID:          "host_maintainer",
			Name:        "主机资源管理员",
			Description: "主机资源管理员角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewResourcePoolHost},
				{ID: iamtypes.CreateResourcePoolHost, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
				{ID: iamtypes.CreateResourcePoolDirectory},
				{ID: iamtypes.EditResourcePoolDirectory, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
				{ID: iamtypes.DeleteResourcePoolDirectory, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
				{ID: iamtypes.EditResourcePoolHost, ResourceTypeID: iamtypes.SysHost},
				{ID: iamtypes.DeleteResourcePoolHost, ResourceTypeID: iamtypes.SysHost},
				{ID: iamtypes.TransferHostOutOfResPoolDir, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
				{ID: iamtypes.TransferHostToResPoolDir, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
				{ID: iamtypes.ManageHostAgentID},
			},
		},
		{
			ID:          "model_maintainer",
			Name:        "模型关系维护人",
			Description: "模型关系维护人角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateModelGroup},
				{ID: iamtypes.EditModelGroup, ResourceTypeID: iamtypes.SysModelGroup},
				{ID: iamtypes.DeleteModelGroup, ResourceTypeID: iamtypes.SysModelGroup},
				{ID: iamtypes.EditBusinessLayer},
				{ID: iamtypes.EditModelTopologyView},
				{ID: iamtypes.CreateSysModel, ResourceTypeID: iamtypes.SysModelGroup},
				{ID: iamtypes.EditSysModel, ResourceTypeID: iamtypes.SysModel},
				{ID: iamtypes.DeleteSysModel, ResourceTypeID: iamtypes.SysModel},
				{ID: iamtypes.CreateAssociationType},
				{ID: iamtypes.EditAssociationType, ResourceTypeID: iamtypes.SysAssociationType},
				{ID: iamtypes.DeleteAssociationType, ResourceTypeID: iamtypes.SysAssociationType},
			},
		},
		{
			ID:          "developer",
			Name:        "开发者",
			Description: "开发者角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.WatchHostEvent},
				{ID: iamtypes.WatchHostRelationEvent},
				{ID: iamtypes.WatchBizEvent},
				{ID: iamtypes.WatchSetEvent},
				{ID: iamtypes.WatchModuleEvent},
				{ID: iamtypes.WatchProcessEvent},
				{ID: iamtypes.WatchBizSetEvent},
				{ID: iamtypes.WatchPlatEvent},
				{ID: iamtypes.WatchProjectEvent},
				{ID: iamtypes.WatchCommonInstanceEvent, ResourceTypeID: iamtypes.SysModelEvent},
				{ID: iamtypes.WatchMainlineInstanceEvent, ResourceTypeID: iamtypes.MainlineModelEvent},
				{ID: iamtypes.WatchInstAsstEvent, ResourceTypeID: iamtypes.InstAsstEvent},
			},
		},
		{
			ID:          "auditor",
			Name:        "审计员",
			Description: "审计员角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.FindAuditLog},
			},
		},
	}
}

// genManagerRoles generate resource manager roles
// NOCC:golint/fnsize(角色定义需要放在一起)
func genManagerRoles() []iam.Role {
	return []iam.Role{
		{
			ID:          "project_manager",
			Name:        "项目管理员",
			Description: "项目管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewProject},
				{ID: iamtypes.EditProject, ResourceTypeID: iamtypes.Project},
				{ID: iamtypes.DeleteProject, ResourceTypeID: iamtypes.Project},
			},
		},
		{
			ID:          "cloud_area_manager",
			Name:        "管控区域管理员",
			Description: "管控区域管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewCloudArea},
				{ID: iamtypes.EditCloudArea, ResourceTypeID: iamtypes.SysCloudArea},
				{ID: iamtypes.DeleteCloudArea, ResourceTypeID: iamtypes.SysCloudArea},
			},
		},
		{
			ID:          "model_group_manager",
			Name:        "模型分组管理员",
			Description: "模型分组管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditModelGroup, ResourceTypeID: iamtypes.SysModelGroup},
				{ID: iamtypes.DeleteModelGroup, ResourceTypeID: iamtypes.SysModelGroup},
			},
		},
		{
			ID:          "sys_model_manager",
			Name:        "模型管理员",
			Description: "模型管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewSysModel, ResourceTypeID: iamtypes.SysModel},
				{ID: iamtypes.EditSysModel, ResourceTypeID: iamtypes.SysModel},
				{ID: iamtypes.DeleteSysModel, ResourceTypeID: iamtypes.SysModel},
			},
		},
		{
			ID:          "asst_type_manager",
			Name:        "关联类型管理员",
			Description: "关联类型管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditAssociationType, ResourceTypeID: iamtypes.SysAssociationType},
				{ID: iamtypes.DeleteAssociationType, ResourceTypeID: iamtypes.SysAssociationType},
			},
		},
		{
			ID:          "host_pool_dir_manager",
			Name:        "主机池目录管理员",
			Description: "主机池目录管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditResourcePoolDirectory, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
				{ID: iamtypes.DeleteResourcePoolDirectory, ResourceTypeID: iamtypes.SysResourcePoolDirectory},
			},
		},
		{
			ID:          "field_tpl_manager",
			Name:        "字段组合模板管理员",
			Description: "字段组合模板管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewFieldGroupingTemplate, ResourceTypeID: iamtypes.FieldGroupingTemplate},
				{ID: iamtypes.EditFieldGroupingTemplate, ResourceTypeID: iamtypes.FieldGroupingTemplate},
				{ID: iamtypes.DeleteFieldGroupingTemplate, ResourceTypeID: iamtypes.FieldGroupingTemplate},
			},
		},
		{
			ID:          "host_manager",
			Name:        "业务主机管理员",
			Description: "业务主机管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditBusinessHost, ResourceTypeID: iamtypes.Host},
			},
		},
		{
			ID:          "biz_svc_tpl_manager",
			Name:        "服务模板管理员",
			Description: "服务模板管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditBusinessServiceTemplate, ResourceTypeID: iamtypes.BizProcessServiceTemplate},
				{ID: iamtypes.DeleteBusinessServiceTemplate, ResourceTypeID: iamtypes.BizProcessServiceTemplate},
			},
		},
		{
			ID:          "biz_set_tpl_manager",
			Name:        "集群模板管理员",
			Description: "集群模板管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditBusinessSetTemplate, ResourceTypeID: iamtypes.BizSetTemplate},
				{ID: iamtypes.DeleteBusinessSetTemplate, ResourceTypeID: iamtypes.BizSetTemplate},
			},
		},
		{
			ID:          "biz_dyn_query_manager",
			Name:        "动态分组管理员",
			Description: "动态分组管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditBusinessCustomQuery, ResourceTypeID: iamtypes.BizCustomQuery},
				{ID: iamtypes.DeleteBusinessCustomQuery, ResourceTypeID: iamtypes.BizCustomQuery},
			},
		},
	}
}

// genCreatorRoles generate create-dimension roles
func genCreatorRoles() []iam.Role {
	return []iam.Role{
		{
			ID:          "business_creator",
			Name:        "业务创建者",
			Description: "业务创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateBusiness},
			},
		},
		{
			ID:          "biz_set_creator",
			Name:        "业务集创建者",
			Description: "业务集创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateBizSet},
			},
		},
		{
			ID:          "project_creator",
			Name:        "项目创建者",
			Description: "项目创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateProject},
			},
		},
		{
			ID:          "cloud_area_creator",
			Name:        "管控区域创建者",
			Description: "管控区域创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateCloudArea},
			},
		},
		{
			ID:          "model_group_creator",
			Name:        "模型分组创建者",
			Description: "模型分组创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateModelGroup},
			},
		},
		{
			ID:          "sys_model_creator",
			Name:        "模型创建者",
			Description: "模型创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateSysModel, ResourceTypeID: iamtypes.SysModelGroup},
			},
		},
		{
			ID:          "asst_type_creator",
			Name:        "关联类型创建者",
			Description: "关联类型创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateAssociationType},
			},
		},
		{
			ID:          "host_pool_dir_creator",
			Name:        "主机池目录创建者",
			Description: "主机池目录创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateResourcePoolDirectory},
			},
		},
		{
			ID:          "field_tpl_creator",
			Name:        "字段组合模板创建者",
			Description: "字段组合模板创建者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateFieldGroupingTemplate},
			},
		},
	}
}

// genViewerRoles generate extra view-dimension roles
func genViewerRoles() []iam.Role {
	return []iam.Role{
		{
			ID:          "fulltext_search_user",
			Name:        "检索服务使用者",
			Description: "检索服务使用者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.UseFulltextSearch},
			},
		},
		{
			ID:          "project_viewer",
			Name:        "项目查看者",
			Description: "项目查看者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewProject},
			},
		},
		{
			ID:          "cloud_area_viewer",
			Name:        "管控区域查看者",
			Description: "管控区域查看者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewCloudArea},
			},
		},
		{
			ID:          "sys_model_viewer",
			Name:        "模型查看者",
			Description: "模型查看者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewSysModel, ResourceTypeID: iamtypes.SysModel},
				{ID: iamtypes.ViewModelTopo},
			},
		},
		{
			ID:          "field_tpl_viewer",
			Name:        "字段组合模板查看者",
			Description: "字段组合模板查看者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewFieldGroupingTemplate, ResourceTypeID: iamtypes.FieldGroupingTemplate},
			},
		},
	}
}

// genOwnerRoles generate creator-related owner roles when manager action set is not equal
func genOwnerRoles() []iam.Role {
	return []iam.Role{
		{
			ID:          "biz_owner",
			Name:        "业务属主",
			Description: "业务创建者关联授权角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditBusiness, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.ArchiveBusiness, ResourceTypeID: iamtypes.Business},
				{ID: iamtypes.FindBusiness, ResourceTypeID: iamtypes.Business},
			},
		},
		{
			ID:          "biz_set_owner",
			Name:        "业务集属主",
			Description: "业务集创建者关联授权角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditBizSet, ResourceTypeID: iamtypes.BizSet},
				{ID: iamtypes.DeleteBizSet, ResourceTypeID: iamtypes.BizSet},
				{ID: iamtypes.ViewBizSet, ResourceTypeID: iamtypes.BizSet},
			},
		},
		{
			ID:          "cloud_area_owner",
			Name:        "管控区域属主",
			Description: "管控区域创建者关联授权角色",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditCloudArea, ResourceTypeID: iamtypes.SysCloudArea},
				{ID: iamtypes.DeleteCloudArea, ResourceTypeID: iamtypes.SysCloudArea},
			},
		},
	}
}

// genPlatformRoles generate roles for platform and hidden capabilities, one role per type
// TODO support hidden role when iam support this feature
func genPlatformRoles() []iam.Role {
	roles := []iam.Role{
		{
			ID:          "global_settings_manager",
			Name:        "全局设置管理员",
			Description: "全局设置管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.GlobalSettings},
			},
		},
		{
			ID:          "id_rule_manager",
			Name:        "ID规则管理员",
			Description: "ID规则管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.EditIDRuleIncrID},
			},
		},
		{
			ID:          "full_sync_cond_manager",
			Name:        "全量同步缓存条件管理员",
			Description: "全量同步缓存条件管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateFullSyncCond},
				{ID: iamtypes.ViewFullSyncCond},
				{ID: iamtypes.EditFullSyncCond},
				{ID: iamtypes.DeleteFullSyncCond},
			},
		},
		{
			ID:          "full_sync_cond_viewer",
			Name:        "全量同步缓存条件查看者",
			Description: "全量同步缓存条件查看者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewFullSyncCond},
			},
		},
		{
			ID:          "general_cache_viewer",
			Name:        "通用缓存查看者",
			Description: "通用缓存查看者",
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewGeneralCache, ResourceTypeID: iamtypes.GeneralCache},
			},
		},
		{
			ID:          "kube_manager",
			Name:        "容器管理员",
			Description: "容器管理员",
			Actions: []iam.RoleAction{
				{ID: iamtypes.CreateContainerCluster},
				{ID: iamtypes.EditContainerCluster},
				{ID: iamtypes.DeleteContainerCluster},
				{ID: iamtypes.CreateContainerNode},
				{ID: iamtypes.EditContainerNode},
				{ID: iamtypes.DeleteContainerNode},
				{ID: iamtypes.CreateContainerNamespace},
				{ID: iamtypes.EditContainerNamespace},
				{ID: iamtypes.DeleteContainerNamespace},
				{ID: iamtypes.CreateContainerWorkload},
				{ID: iamtypes.EditContainerWorkload},
				{ID: iamtypes.DeleteContainerWorkload},
				{ID: iamtypes.CreateContainerPod},
				{ID: iamtypes.DeleteContainerPod},
			},
		},
		{
			ID:          "tenant_set_visitor",
			Name:        "租户集访问者",
			Description: "租户集访问者",
			// TODO set tenant after iam support it
			Actions: []iam.RoleAction{
				{ID: iamtypes.ViewTenantSet, ResourceTypeID: iamtypes.TenantSet},
				{ID: iamtypes.AccessTenantSet, ResourceTypeID: iamtypes.TenantSet},
			},
		},
	}

	return roles
}

// genDynamicInstRoles generate manager and viewer roles for each dynamic model
func genDynamicInstRoles(tenantObjects map[string][]metadata.Object) []iam.Role {
	roles := make([]iam.Role, 0)
	for _, objects := range tenantObjects {
		for _, obj := range objects {
			managerActions, viewerActions := make([]iam.RoleAction, 0), make([]iam.RoleAction, 0)

			for _, a := range genDynamicAction(obj) {
				item := iam.RoleAction{ID: a.ActionID}
				if a.ActionType == iamtypes.Edit || a.ActionType == iamtypes.Delete {
					item.ResourceTypeID = GenIAMDynamicResTypeID(obj.ID)
				}
				managerActions = append(managerActions, item)
				if a.ActionType == iamtypes.View {
					viewerActions = []iam.RoleAction{item}
				}
			}

			// TODO set tenant after iam support it
			roles = append(roles, iam.Role{
				ID:          GenIAMDynamicRoleID(obj.ID, "manager"),
				Name:        fmt.Sprintf("%s实例管理员", obj.ObjectName),
				Description: fmt.Sprintf("%s实例管理员", obj.ObjectName),
				Actions:     managerActions,
			}, iam.Role{
				ID:          GenIAMDynamicRoleID(obj.ID, "viewer"),
				Name:        fmt.Sprintf("%s实例查看者", obj.ObjectName),
				Description: fmt.Sprintf("%s实例查看者", obj.ObjectName),
				Actions:     viewerActions,
			})
		}
	}
	return roles
}
