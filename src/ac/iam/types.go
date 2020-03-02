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

	"configcenter/src/common/auth"
)

const (
	iamRequestHeader   = "X-Request-Id"
	iamAppCodeHeader   = "X-Bk-App-Code"
	iamAppSecretHeader = "X-Bk-App-Secret"

	SystemIDCMDB   = "bk_cmdb"
	SystemNameCMDB = "配置平台"

	ScopeTypeIDSystem     = "system"
	ScopeTypeIDSystemName = "全局"

	ScopeTypeIDBiz     = "biz"
	ScopeTypeIDBizName = "业务"
)

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

func ParseConfigFromKV(prefix string, configMap map[string]string) (AuthConfig, error) {
	var cfg AuthConfig

	if !auth.IsAuthed() {
		return AuthConfig{}, nil
	}

	address, exist := configMap[prefix+".address"]
	if !exist {
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

	cfg.AppSecret, exist = configMap[prefix+".appSecret"]
	if !exist {
		return cfg, errors.New(`missing "appSecret" configuration for auth center`)
	}
	if len(cfg.AppSecret) == 0 {
		return cfg, errors.New(`invalid "appSecret" configuration for auth center`)
	}

	cfg.AppCode, exist = configMap[prefix+".appCode"]
	if !exist {
		return cfg, errors.New(`missing "appCode" configuration for auth center`)
	}
	if len(cfg.AppCode) == 0 {
		return cfg, errors.New(`invalid "appCode" configuration for auth center`)
	}

	cfg.SystemID = SystemIDCMDB

	return cfg, nil
}

type System struct {
	ID                 string     `json:"id,omitempty"`
	Name               string     `json:"name,omitempty"`
	EnglishName        string     `json:"name_en,omitempty"`
	Description        string     `json:"description,omitempty"`
	EnglishDescription string     `json:"description_en,omitempty"`
	Clients            string     `json:"clients,omitempty"`
	ProviderConfig     *SysConfig `json:"provider_config"`
}

type SysConfig struct {
	Host string `json:"host,omitempty"`
	Auth string `json:"auth,omitempty"`
}

type SystemResp struct {
	BaseResponse
	Data RegisteredSystemInfo `json:"data"`
}

type RegisteredSystemInfo struct {
	BaseInfo      System           `json:"base_info"`
	ResourceTypes []ResourceType   `json:"resource_types"`
	Actions       []ResourceAction `json:"actions"`
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AuthError struct {
	RequestID string
	Reason    error
}

func (a *AuthError) Error() string {
	if len(a.RequestID) == 0 {
		return a.Reason.Error()
	}
	return fmt.Sprintf("iam request id: %s, err: %s", a.RequestID, a.Reason.Error())
}

type ResourceTypeID string

const (
	SysSystemBase            ResourceTypeID = "sys_system_base"
	SysHostInstance          ResourceTypeID = "sys_host_instance"
	SysEventPushing          ResourceTypeID = "sys_event_pushing"
	SysModelGroup            ResourceTypeID = "sys_model_group"
	SysModel                 ResourceTypeID = "sys_model"
	SysInstance              ResourceTypeID = "sys_instance"
	SysAssociationType       ResourceTypeID = "sys_association_type"
	SysAuditLog              ResourceTypeID = "sys_audit_log"
	SysOperationStatistic    ResourceTypeID = "sys_operation_statistic"
	SysResourcePoolDirectory ResourceTypeID = "sys_resource_pool_directory"
	SysCloudArea             ResourceTypeID = "sys_cloud_area"
	SysCloudAccount          ResourceTypeID = "sys_cloud_account"
	SysCloudResourceTask     ResourceTypeID = "sys_cloud_resource_task"
)

const (
	Business                  ResourceTypeID = "business"
	BizHostInstance           ResourceTypeID = "biz_host_instance"
	BizCustomQuery            ResourceTypeID = "biz_custom_query"
	BizTopology               ResourceTypeID = "biz_topology"
	BizCustomField            ResourceTypeID = "biz_custom_field"
	BizProcessServiceTemplate ResourceTypeID = "biz_process_service_template"
	BizProcessServiceCategory ResourceTypeID = "biz_process_service_category"
	BizProcessServiceInstance ResourceTypeID = "biz_process_service_instance"
	BizSetTemplate            ResourceTypeID = "biz_set_template"
	BizHostApply              ResourceTypeID = "biz_host_apply"
)

// describe resource type defined and registered to iam.
type ResourceType struct {
	ID             ResourceTypeID `json:"id"`
	Name           string         `json:"name"`
	NameEn         string         `json:"name_en"`
	Description    string         `json:"description"`
	DescriptionEn  string         `json:"description_en"`
	Parents        []Parent       `json:"parents"`
	ProviderConfig ResourceConfig `json:"provider_config"`
	Version        int64          `json:"version"`
}

type ResourceConfig struct {
	// the url to get this resource.
	Path string `json:"path"`
}

type Parent struct {
	// only one value for cmdb.
	// default value: bk_cmdb
	SystemID   string         `json:"system_id"`
	ResourceID ResourceTypeID `json:"id"`
}

type ActionType string

const (
	Create ActionType = "create"
	Delete ActionType = "delete"
	View   ActionType = "view"
	Edit   ActionType = "edit"
	List   ActionType = "list"
)

type ResourceActionID string

const (
	EditBusinessHost                   ResourceActionID = "edit_biz_host"
	BusinessHostTransferToResourcePool ResourceActionID = "biz_host_transfer_to_resource_pool"

	CreateBusinessCustomQuery ResourceActionID = "create_biz_dynamic_query"
	EditBusinessCustomQuery   ResourceActionID = "edit_biz_dynamic_query"
	DeleteBusinessCustomQuery ResourceActionID = "delete_biz_dynamic_query"
	FindBusinessCustomQuery   ResourceActionID = "find_biz_dynamic_query"

	EditBusinessCustomField ResourceActionID = "edit_biz_custom_field"

	CreateBusinessServiceCategory ResourceActionID = "create_biz_service_category"
	EditBusinessServiceCategory   ResourceActionID = "edit_biz_service_category"
	DeleteBusinessServiceCategory ResourceActionID = "delete_biz_service_category"

	CreateBusinessServiceInstance ResourceActionID = "create_biz_service_instance"
	EditBusinessServiceInstance   ResourceActionID = "edit_biz_service_instance"
	DeleteBusinessServiceInstance ResourceActionID = "delete_biz_service_instance"

	CreateBusinessServiceTemplate ResourceActionID = "create_biz_service_template"
	EditBusinessServiceTemplate   ResourceActionID = "edit_biz_service_template"
	DeleteBusinessServiceTemplate ResourceActionID = "delete_biz_service_template"

	CreateBusinessSetTemplate ResourceActionID = "create_biz_set_template"
	EditBusinessSetTemplate   ResourceActionID = "edit_biz_set_template"
	DeleteBusinessSetTemplate ResourceActionID = "delete_biz_set_template"

	CreateBusinessTopology ResourceActionID = "create_biz_topology"
	EditBusinessTopology   ResourceActionID = "edit_biz_topology"
	DeleteBusinessTopology ResourceActionID = "delete_biz_topology"

	EditBusinessHostApply ResourceActionID = "edit_biz_host_apply"

	CreateResourcePoolHost              ResourceActionID = "create_resource_pool_host"
	EditResourcePoolHost                ResourceActionID = "edit_resource_pool_host"
	DeleteResourcePoolHost              ResourceActionID = "delete_resource_pool_host"
	ResourcePoolHostTransferToBusiness  ResourceActionID = "resource_pool_host_transfer_to_biz"
	ResourcePoolHostTransferToDirectory ResourceActionID = "resource_pool_host_transfer_to_directory"

	CreateResourcePoolDirectory ResourceActionID = "create_resource_pool_directory"
	EditResourcePoolDirectory   ResourceActionID = "edit_resource_pool_directory"
	DeleteResourcePoolDirectory ResourceActionID = "delete_resource_pool_directory"

	CreateBusiness  ResourceActionID = "create_business"
	EditBusiness    ResourceActionID = "edit_business"
	ArchiveBusiness ResourceActionID = "archive_business"
	FindBusiness    ResourceActionID = "find_business"

	CreateCloudArea ResourceActionID = "create_cloud_area"
	EditCloudArea   ResourceActionID = "edit_cloud_area"
	DeleteCloudArea ResourceActionID = "delete_cloud_area"

	CreateInstance ResourceActionID = "create_instance"
	EditInstance   ResourceActionID = "edit_instance"
	DeleteInstance ResourceActionID = "delete_instance"
	FindInstance   ResourceActionID = "find_instance"

	CreateEventPushing ResourceActionID = "create_event_subscription"
	EditEventPushing   ResourceActionID = "edit_event_subscription"
	DeleteEventPushing ResourceActionID = "delete_event_subscription"
	FindEventPushing   ResourceActionID = "find_event_subscription"

	CreateCloudAccount ResourceActionID = "create_cloud_account"
	EditCloudAccount   ResourceActionID = "edit_cloud_account"
	DeleteCloudAccount ResourceActionID = "delete_cloud_account"
	FindCloudAccount   ResourceActionID = "find_cloud_account"

	CreateCloudResourceTask ResourceActionID = "create_cloud_resource_task"
	EditCloudResourceTask   ResourceActionID = "edit_cloud_resource_task"
	DeleteCloudResourceTask ResourceActionID = "delete_cloud_resource_task"
	FindCloudResourceTask   ResourceActionID = "find_cloud_resource_task"

	CreateModel ResourceActionID = "create_model"
	EditModel   ResourceActionID = "edit_model"
	DeleteModel ResourceActionID = "delete_model"
	FindModel   ResourceActionID = "find_model"

	CreateAssociationType ResourceActionID = "create_association_type"
	EditAssociationType   ResourceActionID = "edit_association_type"
	DeleteAssociationType ResourceActionID = "delete_association_type"

	CreateModelGroup ResourceActionID = "create_model_group"
	EditModelGroup   ResourceActionID = "edit_model_group"
	DeleteModelGroup ResourceActionID = "delete_model_group"

	EditBusinessLayer ResourceActionID = "edit_business_layer"

	EditModelTopologyView ResourceActionID = "edit_model_topology_view"

	FindOperationStatistic ResourceActionID = "find_operation_statistic"
	EditOperationStatistic ResourceActionID = "edit_operation_statistic"

	FindAuditLog ResourceActionID = "find_audit_log"

	// Unknown is an action that can not be recognized
	Unknown ResourceActionID = "unknown"
)

type ResourceAction struct {
	// must be a unique id in the whole system.
	ID ResourceActionID `json:"id"`
	// must be a unique name in the whole system.
	Name                 string               `json:"name"`
	NameEn               string               `json:"name_en"`
	Type                 ActionType           `json:"type"`
	RelatedResourceTypes []RelateResourceType `json:"related_resource_types"`
	RelatedActions       []string             `json:"related_actions"`
	Version              int                  `json:"version"`
}

type RelateResourceType struct {
	SystemID           string              `json:"system_id"`
	ID                 ResourceTypeID      `json:"id"`
	NameAlias          string              `json:"name_alias"`
	NameAliasEn        string              `json:"name_alias_en"`
	Scope              *Scope              `json:"scope"`
	InstanceSelections []InstanceSelection `json:"instance_selections"`
}

type Scope struct {
	Op      string         `json:"op"`
	Content []ScopeContent `json:"content"`
}

type ScopeContent struct {
	Op    string `json:"op"`
	Field string `json:"field"`
	Value string `json:"value"`
}

type InstanceSelection struct {
	Name              string          `json:"name"`
	NameEn            string          `json:"name_en"`
	ResourceTypeChain []ResourceChain `json:"resource_type_chain"`
}

type ResourceChain struct {
	SystemID string         `json:"system_id"`
	ID       ResourceTypeID `json:"id"`
}

type iamDiscovery struct {
	servers []string
	index   int
	sync.Mutex
}

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
