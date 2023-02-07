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
	"errors"
	"fmt"
	"strings"
	"sync"

	"configcenter/src/ac/meta"
	"configcenter/src/common/auth"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/scene_server/auth_server/sdk/operator"
)

const (
	// IamRequestHeader TODO
	IamRequestHeader   = "X-Request-Id"
	iamAppCodeHeader   = "X-Bk-App-Code"
	iamAppSecretHeader = "X-Bk-App-Secret"

	// SystemIDCMDB TODO
	SystemIDCMDB = "bk_cmdb"
	// SystemNameCMDBEn TODO
	SystemNameCMDBEn = "cmdb"
	// SystemNameCMDB TODO
	SystemNameCMDB = "配置平台"

	// SystemIDIAM TODO
	SystemIDIAM = "bk_iam"
)

// AuthConfig TODO
type AuthConfig struct {
	// blueking's auth center addresses
	Address []string
	// app code is used for authorize used.
	AppCode string
	// app secret is used for authorized
	AppSecret string
	// the system id that cmdb used in auth center.
	// default value: bk_cmdb
	SystemID string
}

// ParseConfigFromKV TODO
func ParseConfigFromKV(prefix string, configMap map[string]string) (AuthConfig, error) {
	var cfg AuthConfig
	if !auth.EnableAuthorize() {
		return AuthConfig{}, nil
	}
	address, err := cc.String(prefix + ".address")
	if err != nil {
		return cfg, errors.New(`missing "address" configuration for auth center`)
	}

	cfg.Address = strings.Split(strings.Replace(address, " ", "", -1), ",")
	if len(cfg.Address) == 0 {
		return cfg, errors.New(`invalid "address" configuration for auth center`)
	}
	for i := range cfg.Address {
		if !strings.HasSuffix(cfg.Address[i], "/") {
			cfg.Address[i] = cfg.Address[i] + "/"
		}
	}

	appSecret, err := cc.String(prefix + ".appSecret")
	if err != nil {
		return cfg, errors.New(`invalid "appSecret" configuration for auth center`)
	}
	cfg.AppSecret = appSecret
	if len(cfg.AppSecret) == 0 {
		return cfg, errors.New(`invalid "appSecret" configuration for auth center`)
	}

	appCode, err := cc.String(prefix + ".appCode")
	if err != nil {
		return cfg, errors.New(`missing "appCode" configuration for auth center`)
	}
	cfg.AppCode = appCode
	if len(cfg.AppCode) == 0 {
		return cfg, errors.New(`invalid "appCode" configuration for auth center`)
	}

	cfg.SystemID = SystemIDCMDB

	return cfg, nil
}

// System TODO
type System struct {
	ID                 string     `json:"id,omitempty"`
	Name               string     `json:"name,omitempty"`
	EnglishName        string     `json:"name_en,omitempty"`
	Description        string     `json:"description,omitempty"`
	EnglishDescription string     `json:"description_en,omitempty"`
	Clients            string     `json:"clients,omitempty"`
	ProviderConfig     *SysConfig `json:"provider_config"`
}

// SystemQueryField is system query field for searching system info
type SystemQueryField string

const (
	// FieldBaseInfo TODO
	FieldBaseInfo SystemQueryField = "base_info"
	// FieldResourceTypes TODO
	FieldResourceTypes SystemQueryField = "resource_types"
	// FieldActions TODO
	FieldActions SystemQueryField = "actions"
	// FieldActionGroups TODO
	FieldActionGroups SystemQueryField = "action_groups"
	// FieldInstanceSelections TODO
	FieldInstanceSelections SystemQueryField = "instance_selections"
	// FieldResourceCreatorActions TODO
	FieldResourceCreatorActions SystemQueryField = "resource_creator_actions"
	// FieldCommonActions TODO
	FieldCommonActions SystemQueryField = "common_actions"
)

// SysConfig TODO
type SysConfig struct {
	Host string `json:"host,omitempty"`
	Auth string `json:"auth,omitempty"`
}

// SystemResp TODO
type SystemResp struct {
	BaseResponse
	Data RegisteredSystemInfo `json:"data"`
}

// RegisteredSystemInfo TODO
type RegisteredSystemInfo struct {
	BaseInfo               System                 `json:"base_info"`
	ResourceTypes          []ResourceType         `json:"resource_types"`
	Actions                []ResourceAction       `json:"actions"`
	ActionGroups           []ActionGroup          `json:"action_groups"`
	InstanceSelections     []InstanceSelection    `json:"instance_selections"`
	ResourceCreatorActions ResourceCreatorActions `json:"resource_creator_actions"`
	CommonActions          []CommonAction         `json:"common_actions"`
}

// BaseResponse TODO
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AuthError TODO
type AuthError struct {
	RequestID string
	Reason    error
}

// Error 用于错误处理
func (a *AuthError) Error() string {
	if len(a.RequestID) == 0 {
		return a.Reason.Error()
	}
	return fmt.Sprintf("iam request id: %s, err: %s", a.RequestID, a.Reason.Error())
}

// TypeID TODO
type TypeID string

const (
	// SysModelGroup TODO
	SysModelGroup TypeID = "sys_model_group"

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
	// SysOperationStatistic TODO
	SysOperationStatistic TypeID = "sys_operation_statistic"
	// SysResourcePoolDirectory TODO
	SysResourcePoolDirectory TypeID = "sys_resource_pool_directory"
	// SysHostRscPoolDirectory TODO
	SysHostRscPoolDirectory TypeID = "sys_host_rsc_pool_directory"
	// SysCloudArea TODO
	SysCloudArea TypeID = "sys_cloud_area"
	// SysCloudAccount TODO
	SysCloudAccount TypeID = "sys_cloud_account"
	// SysCloudResourceTask TODO
	SysCloudResourceTask TypeID = "sys_cloud_resource_task"
	// SysEventWatch TODO
	SysEventWatch TypeID = "event_watch"
	// Host TODO
	Host TypeID = "host"
	// UserCustom TODO
	UserCustom TypeID = "usercustom"
	// InstAsstEvent instance association resource for resource watch
	InstAsstEvent TypeID = "inst_asst_event"
	// KubeWorkloadEvent kube workload resource for resource watch
	KubeWorkloadEvent TypeID = "kube_workload_event"

	// SkipType TODO
	// for resource type, which is not need to be authorized
	SkipType TypeID = "skip_type"
)

const (
	// BizSet TODO
	BizSet TypeID = "business_set"
	// Business TODO
	Business TypeID = "biz"
	// BusinessForHostTrans TODO
	BusinessForHostTrans TypeID = "biz_for_host_trans"
	// BizCustomQuery TODO
	// Set                       ResourceTypeID = "set"
	// Module                    ResourceTypeID = "module"
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
)

// ResourceType TODO
// describe resource type defined and registered to iam.
type ResourceType struct {
	// unique id
	ID TypeID `json:"id"`
	// unique name
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
	// unique description
	Description    string         `json:"description"`
	DescriptionEn  string         `json:"description_en"`
	Parents        []Parent       `json:"parents"`
	ProviderConfig ResourceConfig `json:"provider_config"`
	Version        int64          `json:"version"`
}

// ResourceConfig TODO
type ResourceConfig struct {
	// the url to get this resource.
	Path string `json:"path"`
}

// Parent TODO
type Parent struct {
	// only one value for cmdb.
	// default value: bk_cmdb
	SystemID   string `json:"system_id"`
	ResourceID TypeID `json:"id"`
}

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
	// List TODO
	List ActionType = "list"
)

// ActionTypeIDNameMap TODO
var ActionTypeIDNameMap = map[ActionType]string{
	Create: "新建",
	Edit:   "编辑",
	Delete: "删除",
	View:   "查询",
}

// ActionID TODO
type ActionID string

const (
	// EditBusinessHost TODO
	EditBusinessHost ActionID = "edit_biz_host"
	// BusinessHostTransferToResourcePool TODO
	BusinessHostTransferToResourcePool ActionID = "unassign_biz_host"
	// HostTransferAcrossBusiness TODO
	HostTransferAcrossBusiness ActionID = "host_transfer_across_business"

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

	// CreateResourcePoolHost TODO
	CreateResourcePoolHost ActionID = "create_resource_pool_host"
	// EditResourcePoolHost TODO
	EditResourcePoolHost ActionID = "edit_resource_pool_host"
	// DeleteResourcePoolHost TODO
	DeleteResourcePoolHost ActionID = "delete_resource_pool_host"
	// ResourcePoolHostTransferToBusiness TODO
	ResourcePoolHostTransferToBusiness ActionID = "assign_host_to_biz"
	// ResourcePoolHostTransferToDirectory TODO
	ResourcePoolHostTransferToDirectory ActionID = "host_transfer_in_resource_pool"

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

	// CreateCloudArea TODO
	CreateCloudArea ActionID = "create_cloud_area"
	// EditCloudArea TODO
	EditCloudArea ActionID = "edit_cloud_area"
	// DeleteCloudArea TODO
	DeleteCloudArea ActionID = "delete_cloud_area"

	// CreateCloudAccount TODO
	CreateCloudAccount ActionID = "create_cloud_account"
	// EditCloudAccount TODO
	EditCloudAccount ActionID = "edit_cloud_account"
	// DeleteCloudAccount TODO
	DeleteCloudAccount ActionID = "delete_cloud_account"
	// FindCloudAccount TODO
	FindCloudAccount ActionID = "find_cloud_account"

	// CreateCloudResourceTask TODO
	CreateCloudResourceTask ActionID = "create_cloud_resource_task"
	// EditCloudResourceTask TODO
	EditCloudResourceTask ActionID = "edit_cloud_resource_task"
	// DeleteCloudResourceTask TODO
	DeleteCloudResourceTask ActionID = "delete_cloud_resource_task"
	// FindCloudResourceTask TODO
	FindCloudResourceTask ActionID = "find_cloud_resource_task"

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

	// EditBusinessLayer TODO
	EditBusinessLayer ActionID = "edit_business_layer"

	// EditModelTopologyView TODO
	EditModelTopologyView ActionID = "edit_model_topology_view"

	// FindOperationStatistic TODO
	FindOperationStatistic ActionID = "find_operation_statistic"
	// EditOperationStatistic TODO
	EditOperationStatistic ActionID = "edit_operation_statistic"

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

	// watch kube related event actions

	// WatchKubeClusterEvent watch kube cluster event action id
	WatchKubeClusterEvent ActionID = "watch_kube_cluster"
	// WatchKubeNodeEvent watch kube node event action id
	WatchKubeNodeEvent ActionID = "watch_kube_node"
	// WatchKubeNamespaceEvent watch kube namespace event action id
	WatchKubeNamespaceEvent ActionID = "watch_kube_namespace"
	// WatchKubeWorkloadEvent watch kube workload event action id, authorized by workload type as sub-resource
	WatchKubeWorkloadEvent ActionID = "watch_kube_workload"
	// WatchKubePodEvent watch kube pod event action id, its event detail includes containers in it
	WatchKubePodEvent ActionID = "watch_kube_pod"

	// GlobalSettings TODO
	GlobalSettings ActionID = "global_settings"

	// Unsupported TODO
	// Unknown is an action that can not be recognized
	Unsupported ActionID = "unsupported"
	// Skip is an action that no need to auth
	Skip ActionID = "skip"
)

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
	// IAMSysInstTypePrefix TODO
	// IAM侧资源的通用模型实例前缀标识
	IAMSysInstTypePrefix = meta.CMDBSysInstTypePrefix
)

// ResourceAction TODO
type ResourceAction struct {
	// must be a unique id in the whole system.
	ID ActionID `json:"id"`
	// must be a unique name in the whole system.
	Name                 string               `json:"name"`
	NameEn               string               `json:"name_en"`
	Type                 ActionType           `json:"type"`
	RelatedResourceTypes []RelateResourceType `json:"related_resource_types"`
	RelatedActions       []ActionID           `json:"related_actions"`
	Version              int                  `json:"version"`
}

// SelectionMode 选择类型, 资源在权限中心产品上配置权限时的作用范围
type SelectionMode string

const (
	// 仅可选择实例, 默认值
	modeInstance SelectionMode = "instance"
	// 仅可配置属性, 此时instance_selections配置不生效
	modeAttribute SelectionMode = "attribute"
	// 可以同时选择实例和配置属性
	modeAll SelectionMode = "all"
)

// RelateResourceType TODO
type RelateResourceType struct {
	SystemID           string                     `json:"system_id"`
	ID                 TypeID                     `json:"id"`
	NameAlias          string                     `json:"name_alias"`
	NameAliasEn        string                     `json:"name_alias_en"`
	Scope              *Scope                     `json:"scope"`
	SelectionMode      SelectionMode              `json:"selection_mode"`
	InstanceSelections []RelatedInstanceSelection `json:"related_instance_selections"`
}

// Scope TODO
type Scope struct {
	Op      string         `json:"op"`
	Content []ScopeContent `json:"content"`
}

// ScopeContent TODO
type ScopeContent struct {
	Op    string `json:"op"`
	Field string `json:"field"`
	Value string `json:"value"`
}

// RelatedInstanceSelection TODO
type RelatedInstanceSelection struct {
	ID       InstanceSelectionID `json:"id"`
	SystemID string              `json:"system_id"`
	// if true, then this selected instance with not be calculated to calculate the auth.
	// as is will be ignored, the only usage for this selection is to support a convenient
	// way for user to find it's resource instances.
	IgnoreAuthPath bool `json:"ignore_iam_path"`
}

// ActionGroup TODO
// groups related resource actions to make action selection more organized
type ActionGroup struct {
	// must be a unique name in the whole system.
	Name      string         `json:"name"`
	NameEn    string         `json:"name_en"`
	SubGroups []ActionGroup  `json:"sub_groups,omitempty"`
	Actions   []ActionWithID `json:"actions,omitempty"`
}

// ActionWithID TODO
type ActionWithID struct {
	ID ActionID `json:"id"`
}

// InstanceSelectionID TODO
type InstanceSelectionID string

const (
	// BusinessSelection 业务的两种视图，管理的资源也相同，仅名称做区分
	BusinessSelection InstanceSelectionID = "business"
	// BusinessHostTransferSelection TODO
	BusinessHostTransferSelection InstanceSelectionID = "business_host_transfer"
	// BizSetSelection TODO
	BizSetSelection InstanceSelectionID = "business_set_list"
	// BizHostInstanceSelection TODO
	BizHostInstanceSelection InstanceSelectionID = "biz_host_instance"
	// BizCustomQuerySelection TODO
	BizCustomQuerySelection InstanceSelectionID = "biz_custom_query"
	// BizProcessServiceTemplateSelection TODO
	BizProcessServiceTemplateSelection InstanceSelectionID = "biz_process_service_template"
	// BizSetTemplateSelection TODO
	BizSetTemplateSelection InstanceSelectionID = "biz_set_template"
	// SysHostInstanceSelection TODO
	SysHostInstanceSelection InstanceSelectionID = "sys_host_instance"
	// SysModelGroupSelection TODO
	SysModelGroupSelection InstanceSelectionID = "sys_model_group"
	// SysModelSelection TODO
	SysModelSelection InstanceSelectionID = "sys_model"
	// SysModelEventSelection TODO
	SysModelEventSelection InstanceSelectionID = "sys_model_event"
	// MainlineModelEventSelection TODO
	MainlineModelEventSelection InstanceSelectionID = "mainline_model_event"
	// KubeWorkloadEventSelection k8s workload event selection id
	KubeWorkloadEventSelection InstanceSelectionID = "kube_workload_event"
	// SysInstanceModelSelection TODO
	SysInstanceModelSelection InstanceSelectionID = "sys_instance_model"
	// SysAssociationTypeSelection TODO
	SysAssociationTypeSelection InstanceSelectionID = "sys_association_type"
	// SysCloudAreaSelection TODO
	SysCloudAreaSelection InstanceSelectionID = "sys_cloud_area"
	// SysCloudAccountSelection TODO
	SysCloudAccountSelection InstanceSelectionID = "sys_cloud_account"
	// SysCloudResourceTaskSelection TODO
	SysCloudResourceTaskSelection InstanceSelectionID = "sys_cloud_resource_task"
	// InstAsstEventSelection TODO
	InstAsstEventSelection InstanceSelectionID = "inst_asst_event"
	// SysResourcePoolDirectorySelection 主机池目录的两种视图，管理的资源也相同，仅名称做区分
	SysResourcePoolDirectorySelection InstanceSelectionID = "sys_resource_pool_directory"
	// SysHostRscPoolDirectorySelection TODO
	SysHostRscPoolDirectorySelection InstanceSelectionID = "sys_host_rsc_pool_directory"
)

// InstanceSelection TODO
type InstanceSelection struct {
	// unique
	ID InstanceSelectionID `json:"id"`
	// unique
	Name string `json:"name"`
	// unique
	NameEn            string          `json:"name_en"`
	ResourceTypeChain []ResourceChain `json:"resource_type_chain"`
}

// ResourceChain TODO
type ResourceChain struct {
	SystemID string `json:"system_id"`
	ID       TypeID `json:"id"`
}

type iamDiscovery struct {
	servers []string
	index   int
	sync.Mutex
}

// GetServers TODO
func (s *iamDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
}

// GetServersChan TODO
func (s *iamDiscovery) GetServersChan() chan []string {
	return nil
}

// RscTypeAndID TODO
// resource type with id, used to represent resource layer from root to leaf
type RscTypeAndID struct {
	ResourceType TypeID `json:"resource_type"`
	ResourceID   string `json:"resource_id,omitempty"`
}

// Resource TODO
// iam resource, system is resource's iam system id, type is resource type, resource id and attribute are used for filtering
type Resource struct {
	System    string                 `json:"system"`
	Type      TypeID                 `json:"type"`
	ID        string                 `json:"id,omitempty"`
	Attribute map[string]interface{} `json:"attribute,omitempty"`
}

// ResourceCreatorActions TODO
// specifies resource creation actions' related actions that resource creator will have permissions to
type ResourceCreatorActions struct {
	Config []ResourceCreatorAction `json:"config"`
}

// ResourceCreatorAction TODO
type ResourceCreatorAction struct {
	ResourceID       TypeID                  `json:"id"`
	Actions          []CreatorRelatedAction  `json:"actions"`
	SubResourceTypes []ResourceCreatorAction `json:"sub_resource_types,omitempty"`
}

// CreatorRelatedAction TODO
type CreatorRelatedAction struct {
	ID         ActionID `json:"id"`
	IsRequired bool     `json:"required"`
}

// CommonAction specifies a common operation's related iam actions
type CommonAction struct {
	Name        string         `json:"name"`
	EnglishName string         `json:"name_en"`
	Actions     []ActionWithID `json:"actions"`
}

// DynamicAction is dynamic model action
type DynamicAction struct {
	ActionID     ActionID
	ActionType   ActionType
	ActionNameCN string
	ActionNameEN string
}

// DeleteCMDBResourceParam TODO
type DeleteCMDBResourceParam struct {
	ActionIDs            []ActionID
	InstanceSelectionIDs []InstanceSelectionID
	TypeIDs              []TypeID
}

// ListPoliciesParams list iam policies parameter
type ListPoliciesParams struct {
	ActionID  ActionID
	Page      int64
	PageSize  int64
	Timestamp int64
}

// ListPoliciesResp list iam policies response
type ListPoliciesResp struct {
	BaseResponse
	Data *ListPoliciesData `json:"data"`
}

// ListPoliciesData list policy data, which represents iam policies
type ListPoliciesData struct {
	Metadata PolicyMetadata `json:"metadata"`
	Count    int64          `json:"count"`
	Results  []PolicyResult `json:"results"`
}

// PolicyMetadata iam policy metadata
type PolicyMetadata struct {
	System    string       `json:"system"`
	Action    ActionWithID `json:"action"`
	Timestamp int64        `json:"timestamp"`
}

// PolicyResult iam policy result
type PolicyResult struct {
	Version    string           `json:"version"`
	ID         int64            `json:"id"`
	Subject    PolicySubject    `json:"subject"`
	Expression *operator.Policy `json:"expression"`
	ExpiredAt  int64            `json:"expired_at"`
}

// PolicySubject policy subject, which represents user or user group for now
type PolicySubject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SimplifiedInstance simplified instance with only id and name
type SimplifiedInstance struct {
	InstanceID   int64  `json:"bk_inst_id" bson:"bk_inst_id"`
	InstanceName string `json:"bk_inst_name" bson:"bk_inst_name"`
}
