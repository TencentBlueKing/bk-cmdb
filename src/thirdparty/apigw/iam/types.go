/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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
	"errors"
	"fmt"

	"configcenter/src/ac/iam/types"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/thirdparty/apigw/apigwutil"
)

const (
	codeNotFound     = 1901404
	IamRequestHeader = "X-Request-Id"
)

var (
	// ErrNotFound iam system not found
	ErrNotFound = errors.New("Not Found")
)

// AuthError iam auth server error
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

type apiGWIamPermissionParams struct {
	metadata.IamPermission `json:",inline"`
}

type iamInstanceParams struct {
	metadata.IamInstanceWithCreator `json:",inline"`
}

type iamInstancesParams struct {
	metadata.IamInstancesWithCreator `json:",inline"`
}

type iamPermissionURLResp struct {
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
	apigwutil.ApiGWBaseResponse
}

type iamCreatorActionResp struct {
	apigwutil.ApiGWBaseResponse
	Data []metadata.IamCreatorActionPolicy `json:"data"`
}

type iamBatchOperateInstanceAuthParams struct {
	*metadata.IamBatchOperateInstanceAuthReq `json:",inline"`
}

type iamBatchOperateInstanceAuthResp struct {
	apigwutil.ApiGWBaseResponse
	Data []metadata.IamBatchOperateInstanceAuthRes `json:"data"`
}

// SysConfig TODO
type SysConfig struct {
	Host string `json:"host,omitempty"`
	Auth string `json:"auth,omitempty"`
}

// here is split line

// ResourceType TODO
// describe resource type defined and registered to iam.
type ResourceType struct {
	// unique id
	ID types.TypeID `json:"id"`
	// unique name
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
	// unique description
	Description    string         `json:"description"`
	DescriptionEn  string         `json:"description_en"`
	Parents        []Parent       `json:"parents"`
	ProviderConfig ResourceConfig `json:"provider_config"`
	Version        int64          `json:"version"`
	TenantID       string         `json:"tenant_id,omitempty"`
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
	SystemID   string       `json:"system_id"`
	ResourceID types.TypeID `json:"id"`
}

// ResourceAction TODO
type ResourceAction struct {
	// must be a unique id in the whole system.
	ID types.ActionID `json:"id"`
	// must be a unique name in the whole system.
	Name                 string               `json:"name"`
	NameEn               string               `json:"name_en"`
	Type                 types.ActionType     `json:"type"`
	RelatedResourceTypes []RelateResourceType `json:"related_resource_types"`
	RelatedActions       []types.ActionID     `json:"related_actions"`
	Version              int                  `json:"version"`
	Hidden               bool                 `json:"hidden"`
	TenantID             string               `json:"tenant_id,omitempty"`
}

// RelateResourceType TODO
type RelateResourceType struct {
	SystemID           string                     `json:"system_id"`
	ID                 types.TypeID               `json:"id"`
	NameAlias          string                     `json:"name_alias"`
	NameAliasEn        string                     `json:"name_alias_en"`
	Scope              *Scope                     `json:"scope"`
	SelectionMode      types.SelectionMode        `json:"selection_mode"`
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
	ID       types.InstanceSelectionID `json:"id"`
	SystemID string                    `json:"system_id"`
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
	ID types.ActionID `json:"id"`
}

// InstanceSelection TODO
type InstanceSelection struct {
	// unique
	ID types.InstanceSelectionID `json:"id"`
	// unique
	Name string `json:"name"`
	// unique
	NameEn            string          `json:"name_en"`
	ResourceTypeChain []ResourceChain `json:"resource_type_chain"`
	TenantID          string          `json:"tenant_id,omitempty"`
}

// ResourceChain TODO
type ResourceChain struct {
	SystemID string       `json:"system_id"`
	ID       types.TypeID `json:"id"`
}

// ResourceCreatorActions specifies resource creation actions' related actions that resource creator
// will have permissions to
type ResourceCreatorActions struct {
	Config []ResourceCreatorAction `json:"config"`
}

// ResourceCreatorAction TODO
type ResourceCreatorAction struct {
	ResourceID       types.TypeID            `json:"id"`
	Actions          []CreatorRelatedAction  `json:"actions"`
	SubResourceTypes []ResourceCreatorAction `json:"sub_resource_types,omitempty"`
}

// CreatorRelatedAction TODO
type CreatorRelatedAction struct {
	ID         types.ActionID `json:"id"`
	IsRequired bool           `json:"required"`
}

// CommonAction specifies a common operation's related iam actions
type CommonAction struct {
	Name        string         `json:"name"`
	EnglishName string         `json:"name_en"`
	Actions     []ActionWithID `json:"actions"`
}

// ListPoliciesParams list iam policies parameter
type ListPoliciesParams struct {
	ActionID  types.ActionID
	Page      int64
	PageSize  int64
	Timestamp int64
}

// ListPoliciesResp list iam policies response
type ListPoliciesResp struct {
	apigwutil.ApiGWBaseResponse
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

// SystemResp TODO
type SystemResp struct {
	apigwutil.ApiGWBaseResponse
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

// ----authserver----
// AuthOptions describes a item to be authorized
type AuthOptions struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Action    Action     `json:"action"`
	Resources []Resource `json:"resources"`
}

type GetPolicyOption AuthOptions

// Action define's the use's action, which is must correspond to the registered action ids in iam.
type Action struct {
	ID string `json:"id"`
}

// ActionPolicy TODO
type ActionPolicy struct {
	Action Action           `json:"action"`
	Policy *operator.Policy `json:"condition"`
}

// Resource defines all the information used to authorize a resource.
type Resource struct {
	System    string             `json:"system"`
	Type      IamResourceType    `json:"type"`
	ID        string             `json:"id"`
	Attribute ResourceAttributes `json:"attribute"`
}

// AuthBatch TODO
type AuthBatch struct {
	Action    Action     `json:"action"`
	Resources []Resource `json:"resources"`
}

type ResourceAttributes map[string]interface{}

// GetPolicyResp TODO
type GetPolicyResp struct {
	apigwutil.ApiGWBaseResponse
	Data *operator.Policy `json:"data"`
}

// ListPolicyOptions TODO
type ListPolicyOptions struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Actions   []Action   `json:"actions"`
	Resources []Resource `json:"resources"`
}

// ListPolicyResp TODO
type ListPolicyResp struct {
	apigwutil.ApiGWBaseResponse
	Data []*ActionPolicy `json:"data"`
}

// AuthBatchOptions TODO
type AuthBatchOptions struct {
	System  string       `json:"system"`
	Subject Subject      `json:"subject"`
	Batch   []*AuthBatch `json:"batch"`
}

// AuthorizeList Defines the list structure of authorized instance ids. If the permission type is unlimited, the
// "IsAny" field is true and the "IDS" is empty. Otherwise, the "IsAny" field is false and the "ids" is the specific
// instance ID.
type AuthorizeList struct {
	// ids with permission.
	Ids []string `json:"ids"`
	// is the permission type unrestricted.
	IsAny bool `json:"isAny"`
}

// Validate TODO
func (a AuthBatchOptions) Validate() error {
	if len(a.System) == 0 {
		return errors.New("system is empty")
	}

	if len(a.Subject.Type) == 0 {
		return errors.New("subject.type is empty")
	}

	if len(a.Subject.ID) == 0 {
		return errors.New("subject.id is empty")
	}

	if len(a.Batch) == 0 {
		return nil
	}

	for _, b := range a.Batch {
		if len(b.Action.ID) == 0 {
			return errors.New("empty action id")
		}
	}
	return nil
}

// Subject TODO
type Subject struct {
	Type IamResourceType `json:"type"`
	ID   string          `json:"id"`
}

// IamResourceType TODO
type IamResourceType string

// Validate TODO
func (a AuthOptions) Validate() error {
	if len(a.System) == 0 {
		return errors.New("system is empty")
	}

	if len(a.Subject.Type) == 0 {
		return errors.New("subject.type is empty")
	}

	if len(a.Subject.ID) == 0 {
		return errors.New("subject.id is empty")
	}

	if len(a.Action.ID) == 0 {
		return errors.New("action.id is empty")
	}

	return nil
}
