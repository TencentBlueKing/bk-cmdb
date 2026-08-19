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

// Package types provides types for iam
package types

import "configcenter/src/ac/meta"

// TypeID TODO
type TypeID string

const (
	// SysModelGroup TODO
	SysModelGroup TypeID = "sys_model_group"

	// FieldGroupingTemplate defines object field template
	FieldGroupingTemplate TypeID = "field_grouping_template"

	// SysInstanceModel TODO
	// special model resource for selection of instance, not including models whose instances are managed separately
	SysInstanceModel TypeID = "sys_instance_model"
	// SysModel TODO
	SysModel TypeID = "sys_model"
	// SysModelEvent special model resource for resource watch, not including inner and mainline models
	SysModelEvent TypeID = "sys_model_event"
	// MainlineModelEvent special mainline model resource for resource watch
	MainlineModelEvent TypeID = "mainline_model_event"
	// SysInstance TODO
	SysInstance TypeID = "sys_instance"
	// SysAssociationType TODO
	SysAssociationType TypeID = "sys_association_type"
	// SysAuditLog TODO
	SysAuditLog TypeID = "sys_audit_log"
	// SysResourcePoolDirectory TODO
	SysResourcePoolDirectory TypeID = "sys_resource_pool_directory"
	// SysCloudArea TODO
	SysCloudArea TypeID = "sys_cloud_area"
	// SysEventWatch TODO
	SysEventWatch TypeID = "event_watch"
	// Host TODO
	Host TypeID = "host"
	// SysHost is the host resource in host pool directory.
	SysHost TypeID = "sys_host"
	// UserCustom TODO
	UserCustom TypeID = "usercustom"
	// InstAsstEvent instance association resource for resource watch
	InstAsstEvent TypeID = "inst_asst_event"

	// GeneralCache defines general resource cache auth type
	GeneralCache TypeID = "general_cache"

	// Set is set auth type
	Set TypeID = "set"
	// Module is module auth type
	Module TypeID = "module"

	// SkipType TODO
	// for resource type, which is not need to be authorized
	SkipType TypeID = "skip_type"
)

const (
	// BizSet TODO
	BizSet TypeID = "business_set"
	// Business TODO
	Business TypeID = "biz"
	// BizCustomQuery TODO
	BizCustomQuery TypeID = "biz_custom_query"
	// BizTopology TODO
	BizTopology TypeID = "biz_topology"
	// BizCustomField TODO
	BizCustomField TypeID = "biz_custom_field"
	// BizProcessServiceTemplate TODO
	BizProcessServiceTemplate TypeID = "biz_process_service_template"
	// BizProcessServiceCategory TODO
	BizProcessServiceCategory TypeID = "biz_process_service_category"
	// BizProcessServiceInstance TODO
	BizProcessServiceInstance TypeID = "biz_process_service_instance"
	// BizSetTemplate TODO
	BizSetTemplate TypeID = "biz_set_template"
	// BizHostApply TODO
	BizHostApply TypeID = "biz_host_apply"
	// Project project type id
	Project TypeID = "project"
	// TenantSet is the tenant set type id
	TenantSet TypeID = "tenant_set"
)

// ActionID TODO
type ActionID string

// container related iam action id
const (
	// CreateContainerCluster iam action id
	CreateContainerCluster ActionID = "create_container_cluster"

	// EditContainerCluster iam action id
	EditContainerCluster ActionID = "edit_container_cluster"

	// DeleteContainerCluster iam action id
	DeleteContainerCluster ActionID = "delete_container_cluster"

	// CreateContainerNode iam action id
	CreateContainerNode ActionID = "create_container_node"

	// EditContainerNode iam action id
	EditContainerNode ActionID = "edit_container_node"

	// DeleteContainerNode iam action id
	DeleteContainerNode ActionID = "delete_container_node"

	// CreateContainerNamespace iam action id
	CreateContainerNamespace ActionID = "create_container_namespace"

	// EditContainerNamespace iam action id
	EditContainerNamespace ActionID = "edit_container_namespace"

	// DeleteContainerNamespace iam action id
	DeleteContainerNamespace ActionID = "delete_container_namespace"

	// CreateContainerWorkload iam action id, including create action of deployment, statefulSet, daemonSet ...
	CreateContainerWorkload ActionID = "create_container_workload"

	// EditContainerWorkload iam action id, including edit action of deployment, statefulSet, daemonSet ...
	EditContainerWorkload ActionID = "edit_container_workload"

	// DeleteContainerWorkload iam action id, including delete action of deployment, statefulSet, daemonSet ...
	DeleteContainerWorkload ActionID = "delete_container_workload"

	// CreateContainerPod iam action id
	CreateContainerPod ActionID = "create_container_pod"

	// DeleteContainerPod iam action id
	DeleteContainerPod ActionID = "delete_container_pod"
)

const (
	// EditBusinessHost TODO
	EditBusinessHost ActionID = "edit_biz_host"
	// TransferHostOutOfBiz transfer host out of biz (source side)
	TransferHostOutOfBiz ActionID = "transfer_host_out_of_biz"
	// TransferHostIntoBiz transfer host into biz (destination side)
	TransferHostIntoBiz ActionID = "transfer_host_into_biz"
	// TransferHostOutOfResPoolDir transfer host out of resource pool directory (source side)
	TransferHostOutOfResPoolDir ActionID = "trans_host_out_of_res_pool_dir"
	// TransferHostToResPoolDir transfer host to resource pool directory (destination side)
	TransferHostToResPoolDir ActionID = "trans_host_to_res_pool_dir"

	// CreateBusinessCustomQuery TODO
	CreateBusinessCustomQuery ActionID = "create_biz_dynamic_query"
	// EditBusinessCustomQuery TODO
	EditBusinessCustomQuery ActionID = "edit_biz_dynamic_query"
	// DeleteBusinessCustomQuery TODO
	DeleteBusinessCustomQuery ActionID = "delete_biz_dynamic_query"

	// EditBusinessCustomField TODO
	EditBusinessCustomField ActionID = "edit_biz_custom_field"

	// CreateBusinessServiceCategory TODO
	CreateBusinessServiceCategory ActionID = "create_biz_service_category"
	// EditBusinessServiceCategory TODO
	EditBusinessServiceCategory ActionID = "edit_biz_service_category"
	// DeleteBusinessServiceCategory TODO
	DeleteBusinessServiceCategory ActionID = "delete_biz_service_category"

	// CreateBusinessServiceInstance TODO
	CreateBusinessServiceInstance ActionID = "create_biz_service_instance"
	// EditBusinessServiceInstance TODO
	EditBusinessServiceInstance ActionID = "edit_biz_service_instance"
	// DeleteBusinessServiceInstance TODO
	DeleteBusinessServiceInstance ActionID = "delete_biz_service_instance"

	// CreateBusinessServiceTemplate TODO
	CreateBusinessServiceTemplate ActionID = "create_biz_service_template"
	// EditBusinessServiceTemplate TODO
	EditBusinessServiceTemplate ActionID = "edit_biz_service_template"
	// DeleteBusinessServiceTemplate TODO
	DeleteBusinessServiceTemplate ActionID = "delete_biz_service_template"

	// CreateBusinessSetTemplate TODO
	CreateBusinessSetTemplate ActionID = "create_biz_set_template"
	// EditBusinessSetTemplate TODO
	EditBusinessSetTemplate ActionID = "edit_biz_set_template"
	// DeleteBusinessSetTemplate TODO
	DeleteBusinessSetTemplate ActionID = "delete_biz_set_template"

	// CreateBusinessTopology TODO
	CreateBusinessTopology ActionID = "create_biz_topology"
	// EditBusinessTopology TODO
	EditBusinessTopology ActionID = "edit_biz_topology"
	// DeleteBusinessTopology TODO
	DeleteBusinessTopology ActionID = "delete_biz_topology"

	// EditBusinessHostApply TODO
	EditBusinessHostApply ActionID = "edit_biz_host_apply"

	// ViewResourcePoolHost view resource pool host
	ViewResourcePoolHost ActionID = "view_resource_pool_host"
	// CreateResourcePoolHost TODO
	CreateResourcePoolHost ActionID = "create_resource_pool_host"
	// EditResourcePoolHost TODO
	EditResourcePoolHost ActionID = "edit_resource_pool_host"
	// DeleteResourcePoolHost TODO
	DeleteResourcePoolHost ActionID = "delete_resource_pool_host"
	ManageHostAgentID      ActionID = "manage_host_agent_id"

	// CreateResourcePoolDirectory TODO
	CreateResourcePoolDirectory ActionID = "create_resource_pool_directory"
	// EditResourcePoolDirectory TODO
	EditResourcePoolDirectory ActionID = "edit_resource_pool_directory"
	// DeleteResourcePoolDirectory TODO
	DeleteResourcePoolDirectory ActionID = "delete_resource_pool_directory"

	// CreateBusiness TODO
	CreateBusiness ActionID = "create_business"
	// EditBusiness TODO
	EditBusiness ActionID = "edit_business"
	// ArchiveBusiness TODO
	ArchiveBusiness ActionID = "archive_business"
	// FindBusiness TODO
	FindBusiness ActionID = "find_business"
	// ViewBusinessResource TODO
	ViewBusinessResource ActionID = "find_business_resource"

	// CreateBizSet TODO
	CreateBizSet ActionID = "create_business_set"
	// EditBizSet TODO
	EditBizSet ActionID = "edit_business_set"
	// DeleteBizSet TODO
	DeleteBizSet ActionID = "delete_business_set"
	// ViewBizSet TODO
	ViewBizSet ActionID = "view_business_set"
	// AccessBizSet TODO
	AccessBizSet ActionID = "access_business_set"

	// CreateProject create project action id
	CreateProject ActionID = "create_project"
	// EditProject edit project action id
	EditProject ActionID = "edit_project"
	// DeleteProject delete project action id
	DeleteProject ActionID = "delete_project"
	// ViewProject view project action id
	ViewProject ActionID = "view_project"

	// ViewCloudArea view cloud area
	ViewCloudArea ActionID = "view_cloud_area"
	// CreateCloudArea TODO
	CreateCloudArea ActionID = "create_cloud_area"
	// EditCloudArea TODO
	EditCloudArea ActionID = "edit_cloud_area"
	// DeleteCloudArea TODO
	DeleteCloudArea ActionID = "delete_cloud_area"

	// ViewSysModel view system model
	ViewSysModel ActionID = "view_sys_model"
	// CreateSysModel TODO
	CreateSysModel ActionID = "create_sys_model"
	// EditSysModel TODO
	EditSysModel ActionID = "edit_sys_model"
	// DeleteSysModel TODO
	DeleteSysModel ActionID = "delete_sys_model"

	// CreateAssociationType TODO
	CreateAssociationType ActionID = "create_association_type"
	// EditAssociationType TODO
	EditAssociationType ActionID = "edit_association_type"
	// DeleteAssociationType TODO
	DeleteAssociationType ActionID = "delete_association_type"

	// CreateModelGroup TODO
	CreateModelGroup ActionID = "create_model_group"
	// EditModelGroup TODO
	EditModelGroup ActionID = "edit_model_group"
	// DeleteModelGroup TODO
	DeleteModelGroup ActionID = "delete_model_group"

	// ViewModelTopo view model topo
	ViewModelTopo ActionID = "view_model_topo"
	// EditBusinessLayer TODO
	EditBusinessLayer ActionID = "edit_business_layer"
	// EditModelTopologyView TODO
	EditModelTopologyView ActionID = "edit_model_topology_view"

	// FindAuditLog TODO
	FindAuditLog ActionID = "find_audit_log"

	// WatchHostEvent TODO
	WatchHostEvent ActionID = "watch_host_event"
	// WatchHostRelationEvent TODO
	WatchHostRelationEvent ActionID = "watch_host_relation_event"
	// WatchBizEvent TODO
	WatchBizEvent ActionID = "watch_biz_event"
	// WatchSetEvent TODO
	WatchSetEvent ActionID = "watch_set_event"
	// WatchModuleEvent TODO
	WatchModuleEvent ActionID = "watch_module_event"
	// WatchProcessEvent TODO
	WatchProcessEvent ActionID = "watch_process_event"
	// WatchCommonInstanceEvent TODO
	WatchCommonInstanceEvent ActionID = "watch_comm_model_inst_event"
	// WatchMainlineInstanceEvent TODO
	WatchMainlineInstanceEvent ActionID = "watch_custom_topo_layer_event"
	// WatchInstAsstEvent TODO
	WatchInstAsstEvent ActionID = "watch_inst_asst_event"
	// WatchBizSetEvent TODO
	WatchBizSetEvent ActionID = "watch_biz_set_event"
	// WatchPlatEvent watch cloud area event action id
	WatchPlatEvent ActionID = "watch_plat_event"
	// WatchProjectEvent watch project event action id
	WatchProjectEvent ActionID = "watch_project_event"

	// CreateFieldGroupingTemplate create field grouping template action id
	CreateFieldGroupingTemplate = "create_field_grouping_template"
	// ViewFieldGroupingTemplate view field grouping template action id
	ViewFieldGroupingTemplate = "view_field_grouping_template"
	// EditFieldGroupingTemplate edit field grouping template action id
	EditFieldGroupingTemplate = "edit_field_grouping_template"
	// DeleteFieldGroupingTemplate delete field grouping template action id
	DeleteFieldGroupingTemplate = "delete_field_grouping_template"

	// EditIDRuleIncrID edit id rule self-increasing id action id
	EditIDRuleIncrID ActionID = "edit_id_rule_incr_id"

	// GlobalSettings TODO
	GlobalSettings ActionID = "global_settings"

	// UseFulltextSearch use fulltext search
	UseFulltextSearch ActionID = "use_fulltext_search"

	// CreateFullSyncCond create full sync cond action id
	CreateFullSyncCond = "create_full_sync_cond"
	// ViewFullSyncCond view full sync cond action id
	ViewFullSyncCond = "view_full_sync_cond"
	// EditFullSyncCond edit full sync cond action id
	EditFullSyncCond = "edit_full_sync_cond"
	// DeleteFullSyncCond delete full sync cond action id
	DeleteFullSyncCond = "delete_full_sync_cond"

	// ViewGeneralCache view general resource cache
	ViewGeneralCache = "view_general_cache"

	// ViewTenantSet is the view tenant set action id
	ViewTenantSet ActionID = "view_tenant_set"
	// AccessTenantSet is the access tenant set action id
	AccessTenantSet ActionID = "access_tenant_set"

	// Unsupported TODO
	// Unknown is an action that can not be recognized
	Unsupported ActionID = "unsupported"
	// Skip is an action that no need to auth
	Skip ActionID = "skip"
)

// ActionType TODO
type ActionType string

const (
	// Create TODO
	Create ActionType = "create"
	// Delete TODO
	Delete ActionType = "delete"
	// View TODO
	View ActionType = "view"
	// Edit TODO
	Edit ActionType = "edit"
)

// ActionTypeIDNameMap TODO
var ActionTypeIDNameMap = map[ActionType]string{
	Create: "新建",
	Edit:   "编辑",
	Delete: "删除",
	View:   "查询",
}

const (
	// IAMSysInstTypePrefix is the prefix of iam system instance resource type
	IAMSysInstTypePrefix = meta.CMDBSysInstTypePrefix
)

const (
	// IamRequestHeader TODO
	IamRequestHeader = "X-Request-Id"

	// SystemIDCMDB TODO
	SystemIDCMDB = "bk_cmdb"
	// SystemNameCMDBEn TODO
	SystemNameCMDBEn = "cmdb"
	// SystemNameCMDB TODO
	SystemNameCMDB = "配置平台"

	// SystemIDIAM TODO
	SystemIDIAM = "bk_iam"

	// RegisterIamLock defines the lock key for register iam operation
	RegisterIamLock = "register_iam_lock"
)

// RoleID is IAM V4 role id
type RoleID string

// DeleteCMDBResourceParam is the parameter to delete unused CMDB resources in IAM
type DeleteCMDBResourceParam struct {
	ActionIDs []ActionID
	TypeIDs   []TypeID
	RoleIDs   []RoleID
}

// DynamicAction is dynamic model action
type DynamicAction struct {
	ActionID     ActionID
	ActionType   ActionType
	ActionNameCN string
	ActionNameEN string
}
