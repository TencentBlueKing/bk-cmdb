package authcenter

import (
	"fmt"

	"configcenter/src/auth/meta"
)

// system constant
const (
	SystemIDCMDB   = "bk_cmdb"
	SystemNameCMDB = "配置平台"
)

// ScopeTypeID constant
const (
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
	SystemID string
	// enable sync auth data to iam
	EnableSync          bool
	SyncWorkerCount     int
	SyncIntervalMinutes int
}

type RegisterInfo struct {
	CreatorType string           `json:"creator_type"`
	CreatorID   string           `json:"creator_id"`
	Resources   []ResourceEntity `json:"resources,omitempty"`
}

type ResourceEntity struct {
	ResourceType ResourceTypeID `json:"resource_type"`
	ScopeInfo
	ResourceName string         `json:"resource_name,omitempty"`
	ResourceID   []RscTypeAndID `json:"resource_id,omitempty"`
}

type RscTypeAndID struct {
	ResourceType ResourceTypeID `json:"resource_type"`
	ResourceID   string         `json:"resource_id,omitempty"`
}

type ResourceInfo struct {
	ResourceType ResourceTypeID `json:"resource_type"`
	// this filed is not always used, it's decided by the api
	// that is used.
	ResourceEntity
}

type ScopeInfo struct {
	ScopeType string `json:"scope_type,omitempty"`
	ScopeID   string `json:"scope_id,omitempty"`
}

type ResourceResult struct {
	BaseResponse
	RequestID string       `json:"request_id"`
	Data      ResultStatus `json:"data"`
}

type ResultStatus struct {
	// for create resource result confirm use,
	// which true means register a resource success.
	IsCreated bool `json:"is_created"`
	// for deregister resource result confirm use,
	// which true means deregister success.
	IsDeleted bool `json:"is_deleted"`
	// for update resource result confirm use,
	// which true means update a resource success.
	IsUpdated bool `json:"is_updated"`
}

type DeregisterInfo struct {
	Resources []ResourceEntity `json:"resources"`
}

type UpdateInfo struct {
	ScopeInfo
	ResourceInfo
}

type Principal struct {
	Type string `json:"principal_type"`
	ID   string `json:"principal_id"`
}

type AuthBatch struct {
	Principal
	ScopeInfo
	ResourceActions []ResourceAction `json:"resources_actions"`
}

type BatchResult struct {
	BaseResponse
	RequestID string        `json:"request_id"`
	Data      []BatchStatus `json:"data"`
}

type ResourceAction struct {
	ResourceType ResourceTypeID `json:"resource_type"`
	ResourceID   []RscTypeAndID `json:"resource_id,omitempty"`
	ActionID     ActionID       `json:"action_id"`
}

type BatchStatus struct {
	ActionID     string         `json:"action_id"`
	ResourceType ResourceTypeID `json:"resource_type"`
	// for authorize confirm use, define if a user have
	// the permission to this request.
	IsPass bool `json:"is_pass"`
}

type AuthError struct {
	RequestID string
	Reason    error
}

func (a *AuthError) Error() string {
	if len(a.RequestID) == 0 {
		return a.Reason.Error()
	}
	return fmt.Sprintf("request id: %s, err: %s", a.RequestID, a.Reason.Error())
}

type System struct {
	SystemID   string `json:"system_id,omitempty"`
	SystemName string `json:"system_name"`
	Desc       string `json:"desc"`
	// 可为空，在使用注册资源的方式时
	QueryInterface string `json:"query_interface"`
	//  关联的资源所属，有业务、全局、项目等
	ReleatedScopeTypes string `json:"releated_scope_types"`
	// 管理者，可通过权限中心产品页面修改模型相关信息
	Managers string `json:"managers"`
	// 更新者，可为system
	Updater string `json:"updater,omitempty"`
	// 创建者，可为system
	Creator string `json:"creator,omitempty"`
}

type ResourceType struct {
	ResourceTypeID       ResourceTypeID `json:"resource_type"`
	ResourceTypeName     string         `json:"resource_type_name"`
	ParentResourceTypeID ResourceTypeID `json:"parent_resource_type"`
	Share                bool           `json:"is_share"`
	Actions              []Action       `json:"actions"`
}

type Action struct {
	ActionID          ActionID `json:"action_id"`
	ActionName        string   `json:"action_name"`
	IsRelatedResource bool     `json:"is_related_resource"`
}

type SystemDetail struct {
	System
	Scopes []struct {
		ScopeTypeID   string         `json:"scope_type_id"`
		ResourceTypes []ResourceType `json:"resource_types"`
	} `json:"scopes"`
}

type BaseResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Result    bool   `json:"result"`
	RequestID string `json:"request_id"`
}

type SearchCondition struct {
	ScopeInfo
	ResourceType    ResourceTypeID `json:"resource_type"`
	ParentResources []RscTypeAndID `json:"parent_resources"`
}

type SearchResult struct {
	BaseResponse
	RequestID string                 `json:"request_id"`
	Data      []meta.BackendResource `json:"data"`
}

type PageBackendResource struct {
	Count   int64                  `json:"count"`
	Results []meta.BackendResource `json:"results"`
}

type SearchPageResult struct {
	BaseResponse
	RequestID string              `json:"request_id"`
	Data      PageBackendResource `json:"data"`
}

func (br BaseResponse) ErrorString() string {
	return fmt.Sprintf("request id: %s, error code: %d, message: %s", br.RequestID, br.Code, br.Message)
}

type ListAuthorizedResources struct {
	Principal   `json:",inline"`
	ScopeInfo   `json:",inline"`
	TypeActions []TypeAction `json:"resource_types_actions"`
	// array or string
	DataType string `json:"resource_data_type"`
	Exact    bool   `json:"is_exact_resource"`
}

type TypeAction struct {
	ActionID     ActionID       `json:"action_id"`
	ResourceType ResourceTypeID `json:"resource_type"`
}

type ListAuthorizedResourcesResult struct {
	BaseResponse
	Data []AuthorizedResource `json:"data"`
}

type ListAuthorizedScopeResult struct {
	BaseResponse
	Data []string `json:"data"`
}

type AuthorizedResource struct {
	ActionID     ActionID       `json:"action_id"`
	ResourceType ResourceTypeID `json:"resource_type"`
	ResourceIDs  []IamResource  `json:"resource_ids"`
}

type IamResource []RscTypeAndID

type RoleWithAuthResources struct {
	RoleTemplateName string       `json:"perm_template_name"`
	TemplateID       string       `json:"template_id"`
	Desc             string       `json:"desc"`
	ResourceActions  []RoleAction `json:"resource_types_actions"`
}

type RoleAction struct {
	ScopeTypeID    string         `json:"scope_type_id"`
	ResourceTypeID ResourceTypeID `json:"resource_type_id"`
	ActionID       ActionID       `json:"action_id"`
}

type RegisterRoleResult struct {
	BaseResponse
	Data struct {
		TemplateID int64 `json:"perm_template_id"`
	} `json:"data"`
}

type GetSkipUrlResult struct {
	BaseResponse
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
}

type UserGroupMembersResult struct {
	BaseResponse
	Data []UserGroupMembers `json:"data"`
}

type UserGroupMembers struct {
	ID int64 `json:"group_id"`
	// user's group name, should be one of follows:
	// bk_biz_maintainer, bk_biz_productor, bk_biz_test, bk_biz_developer, operator
	Name  string   `json:"group_code"`
	Users []string `json:"users"`
}
