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
	"net/http"

	"configcenter/src/ac/iam/types"
	"configcenter/src/scene_server/auth_server/sdk/operator"
)

const (
	IamRequestHeader = "X-Request-Id"
	// iamOperatorHeader is the required operator header of IAM authorization APIs.
	iamOperatorHeader = "X-Bkiam-Operator"

	// MaxListPageSize is the max page size of IAM model resource list apis.
	MaxListPageSize = 100
	// MaxAddAuthorizationSize is the max item count of one add_authorization request.
	MaxAddAuthorizationSize = 20
	// MaxAuthorizationExpireDays is the max expire days of IAM authorization.
	MaxAuthorizationExpireDays = 365
	// pageParam is the query parameter name of page number.
	pageParam = "page"
	// pageSizeParam is the query parameter name of page size.
	pageSizeParam = "page_size"
	// idsParam is the query parameter name of the comma separated id list.
	idsParam = "ids"
)

// AuthError iam auth server error
type AuthError struct {
	RequestID  string
	StatusCode int
	Reason     error
}

// Error returns the iam auth error message
func (a *AuthError) Error() string {
	msg := fmt.Sprintf("status: %d, err: %v", a.StatusCode, a.Reason)
	if len(a.RequestID) == 0 {
		return msg
	}
	return fmt.Sprintf("iam request id: %s, %s", a.RequestID, msg)
}

// IsSystemNotExistErr judge whether the error means that the cmdb system is not registered in IAM.
func IsSystemNotExistErr(err error) bool {
	var authErr *AuthError
	if !errors.As(err, &authErr) {
		return false
	}
	return authErr.StatusCode == http.StatusNotFound
}

// PermApplyURLReq is the request of IAM generate_perm_apply_url API.
type PermApplyURLReq struct {
	SystemID    string          `json:"system_id"`
	Permissions []PermApplyItem `json:"permissions"`
}

// PermApplyItem is one action's permission to apply.
type PermApplyItem struct {
	ActionID  string              `json:"action_id"`
	Resources []PermApplyResource `json:"resources,omitempty"`
}

// PermApplyResource is a resource instance in permission apply request.
type PermApplyResource struct {
	ID        string              `json:"id"`
	Type      string              `json:"type"`
	Ancestors []PermApplyAncestor `json:"ancestors,omitempty"`
}

// PermApplyAncestor is an ancestor of a resource instance.
type PermApplyAncestor struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// permApplyURLData is the response data of IAM generate_perm_apply_url API.
type permApplyURLData struct {
	URL string `json:"url"`
}

// System is IAM V4 system info, used by create system request and retrieve system response.
type System struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn      string `json:"name_en,omitempty"`
	Description string `json:"description,omitempty"`
	// TODO: IAM currently does not support DescriptionEn, need to confirm with IAM how to handle it.
	DescriptionEn string   `json:"description_en,omitempty"`
	Managers      []string `json:"managers,omitempty"`
	Clients       []string `json:"clients,omitempty"`
	CallbackURL   string   `json:"callback_url,omitempty"`
}

// RegisterSystemData is the response data of registering a system
type RegisterSystemData struct {
	ID string `json:"id"`
}

// ResourceType is IAM V4 resource type
type ResourceType struct {
	ID   types.TypeID `json:"id"`
	Name string       `json:"name"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn    string         `json:"name_en,omitempty"`
	Ancestors []types.TypeID `json:"ancestors"`
	TenantID  string         `json:"tenant_id,omitempty"`
}

// ResourceAction is IAM V4 action
type ResourceAction struct {
	ID   types.ActionID `json:"id"`
	Name string         `json:"name"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn         string         `json:"name_en,omitempty"`
	ResourceTypeID types.TypeID   `json:"resource_type_id"`
	AuthMode       types.AuthMode `json:"auth_mode,omitempty"`
	// TODO: IAM currently does not support Hidden, need to confirm with IAM how to handle it.
	Hidden   bool   `json:"hidden,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
}

// UpdateResourceTypeReq is the request body of updating a resource type
type UpdateResourceTypeReq struct {
	Name string `json:"name,omitempty"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn    string         `json:"name_en,omitempty"`
	Ancestors []types.TypeID `json:"ancestors,omitempty"`
}

// UpdateActionReq is the request body of updating an action
type UpdateActionReq struct {
	Name string `json:"name"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn   string         `json:"name_en,omitempty"`
	AuthMode types.AuthMode `json:"auth_mode,omitempty"`
}

// Role is IAM V4 role
type Role struct {
	ID   types.RoleID `json:"id"`
	Name string       `json:"name"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn      string `json:"name_en,omitempty"`
	Description string `json:"description,omitempty"`
	// TODO: IAM currently does not support DescriptionEn, need to confirm with IAM how to handle it.
	DescriptionEn string       `json:"description_en,omitempty"`
	Actions       []RoleAction `json:"actions"`
	TenantID      string       `json:"tenant_id,omitempty"`
}

// UpdateRoleReq is the request body of updating a role
type UpdateRoleReq struct {
	Name string `json:"name,omitempty"`
	// TODO: IAM currently does not support NameEn, need to confirm with IAM how to handle it.
	NameEn      string `json:"name_en,omitempty"`
	Description string `json:"description,omitempty"`
	// TODO: IAM currently does not support DescriptionEn, need to confirm with IAM how to handle it.
	DescriptionEn string `json:"description_en,omitempty"`
}

// RoleAction is the action bound to a role
type RoleAction struct {
	ID             types.ActionID `json:"id"`
	ResourceTypeID types.TypeID   `json:"resource_type_id"`
}

// IamErrorResp is IAM V4 error response
type IamErrorResp struct {
	Error *IamErrorData `json:"error"`
}

// IamErrorData is IAM V4 error detail
type IamErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ListResourceTypesData is resource type page data
type ListResourceTypesData struct {
	Count   int64          `json:"count"`
	Results []ResourceType `json:"results"`
}

// ListActionsData is action page data
type ListActionsData struct {
	Count   int64            `json:"count"`
	Results []ResourceAction `json:"results"`
}

// ListRolesData is role page data
type ListRolesData struct {
	Count   int64  `json:"count"`
	Results []Role `json:"results"`
}

// systemAuthToken is the auth token of the cmdb system, which is used to validate if the request is from IAM
type systemAuthToken struct {
	AuthToken string `json:"auth_token"`
}

// SubjectType is the authorized subject type.
type SubjectType string

const (
	// UserSubjectType is the authorized subject type for user.
	UserSubjectType SubjectType = "user"
)

// AuthSubject is the subject to be authorized.
type AuthSubject struct {
	Type SubjectType `json:"type"`
	ID   string      `json:"id"`
}

type AuthAction struct {
	ID types.ActionID `json:"id"`
}

// PlanReq is the request of the hybrid plan api.
type PlanReq struct {
	Subject AuthSubject `json:"subject"`
	Action  AuthAction  `json:"action"`
}

// PlanByActionsReq is the request of the hybrid plan by actions api.
type PlanByActionsReq struct {
	Subject AuthSubject  `json:"subject"`
	Actions []AuthAction `json:"actions"`
}

// ActionPlanRes is one action's authorization plan of the hybrid plan by actions api.
type ActionPlanRes struct {
	ActionID      types.ActionID `json:"action_id"`
	operator.Plan `json:",inline"`
}

// AuthResource is a resource instance used by IAM authorization APIs.
type AuthResource struct {
	Type types.TypeID `json:"type"`
	ID   string       `json:"id"`
}

// AddAuthorizationReq is one item of IAM add_authorization API.
type AddAuthorizationReq struct {
	Subject               AuthSubject    `json:"subject"`
	RoleID                types.RoleID   `json:"role_id"`
	RelatedResourceTypeID types.TypeID   `json:"related_resource_type_id,omitempty"`
	Resources             []AuthResource `json:"resources,omitempty"`
	ExpiredAt             int64          `json:"expired_at"`
}

// ----authserver----
// AuthOptions describes a item to be authorized
type AuthOptions struct {
	System    string     `json:"system"`
	Subject   Subject    `json:"subject"`
	Action    Action     `json:"action"`
	Resources []Resource `json:"resources"`
}

// Action define's the use's action, which is must correspond to the registered action ids in iam.
type Action struct {
	ID string `json:"id"`
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
