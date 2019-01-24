package authcenter

import "configcenter/src/common/metadata"

type RegisterInfo struct {
	CreatorType  string `json:"creator_type"`
	CreatorID    string `json:"creator_id"`
	ScopeInfo    `json:",inline"`
	ResourceInfo `json:",inline"`
}

type ResourceInfo struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name,omitempty"`
	ResourceID   string `json:"resource_id"`
}

type ScopeInfo struct {
	ScopeType string `json:"scope_type"`
	ScopeID   string `json:"scope_id"`
}

type ResourceResult struct {
	metadata.BaseResp `json:",inline"`
	RequestID         string       `json:"request_id"`
	Data              ResultStatus `json:"data"`
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
	ScopeInfo    `json:",inline"`
	ResourceInfo `json:",inline"`
}
