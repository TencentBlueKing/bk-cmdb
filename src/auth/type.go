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

package auth

import "configcenter/src/common/metadata"

type Attribute struct {
	Resource Resource
	User     UserInfo
}

type UserInfo struct {
	// the name of this user.
	UserName string
	// the supplier id that this user belongs to.
	SupplierID string
}

type Resource struct {
	// the name of the resource, which could be a model name.
	Name string

	// the instance id of this resource, which could be a model's instance id.
	InstanceID uint64

	// the action that user want to do with this resource.
	Action Action

	// the version of this resource, which is the api version.
	APIVersion string

	// the business id that this resource belongs to, but it's not necessary for
	// a resource that does not belongs to a business.
	BusinessID uint64
}

type Decision string

const (
	// DecisionDeny means that an authorizer decided to deny the action.
	DecisionDeny Decision = "deny"
	// DecisionAllow means that an authorizer decided to allow the action.
	DecisionAllow Decision = "allow"
	// DecisionNoOpinion means that an authorizer has no opinion on whether
	// to allow or deny an action.
	DecisionNoOpinion Decision = "noOpinion"
)

type Action string

const (
	Create     Action = "create"
	CreateMany Action = "createMany"
	Update     Action = "update"
	UpdateMany Action = "updateMany"
	Delete     Action = "delete"
	DeleteMany Action = "deleteMany"
	Find       Action = "find"
	FindMany   Action = "findMany"
)

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
