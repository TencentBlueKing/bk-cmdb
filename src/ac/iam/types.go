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

import "fmt"

const (
	iamRequestHeader   = "X-Request-Id"
	iamAppCodeHeader   = "X-Bk-App-Code"
	iamAppSecretHeader = "X-Bk-App-Secret"

	systemID = "bk_iam"
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
	SysSystemBase         ResourceTypeID = "sys_system_base"
	SysBusinessInstance   ResourceTypeID = "sys_business_instance"
	SysHostInstance       ResourceTypeID = "sys_host_instance"
	SysEventPushing       ResourceTypeID = "sys_event_pushing"
	SysModelGroup         ResourceTypeID = "sys_model_group"
	SysModel              ResourceTypeID = "sys_model"
	SysInstance           ResourceTypeID = "sys_instance"
	SysAssociationType    ResourceTypeID = "sys_association_type"
	SysAuditLog           ResourceTypeID = "sys_audit_log"
	SysOperationStatistic ResourceTypeID = "sys_operation_statistic"
)

const (
	Business                  ResourceTypeID = "business"
	BizHostInstance           ResourceTypeID = "biz_host_instance"
	BizCustomQuery            ResourceTypeID = "biz_custom_query"
	BizTopology               ResourceTypeID = "biz_topology"
	BizModelGroup             ResourceTypeID = "biz_model_group"
	BizModel                  ResourceTypeID = "biz_model"
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
	CreateBusinessHost ResourceActionID = "create_biz_host"
	EditBusinessHost   ResourceActionID = "edit_biz_host"
	RemoveBusinessHost ResourceActionID = "remove_biz_host"

	CreateBusinessCustomQuery ResourceActionID = "create_biz_dynamic_query"
	EditBusinessCustomQuery   ResourceActionID = "edit_biz_dynamic_query"
	DeleteBusinessCustomQuery ResourceActionID = "delete_biz_dynamic_query"
	FindBusinessCustomQuery   ResourceActionID = "find_biz_dynamic_query"

	CreateBusinessModel ResourceActionID = "create_biz_model"
	EditBusinessModel   ResourceActionID = "edit_biz_model"
	DeleteBusinessModel ResourceActionID = "delete_biz_model"

	CreateBusinessModelGroup ResourceActionID = "create_biz_model_group"
	EditBusinessModelGroup   ResourceActionID = "edit_biz_model_group"
	DeleteBusinessModelGroup ResourceActionID = "delete_biz_model_group"

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
