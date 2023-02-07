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

// Package meta TODO
package meta

// AuthAttribute TODO
type AuthAttribute struct {
	User      UserInfo
	Resources []ResourceAttribute
}

// UserInfo TODO
type UserInfo struct {
	// the name of this user.
	UserName string
	// the supplier id that this user belongs to.
	SupplierAccount string
}

// Item TODO
type Item Basic

// Layers TODO
type Layers []Item

// ResourceAttribute represent one iam resource
type ResourceAttribute struct {
	Basic

	SupplierAccount string `json:"supplier_account"`
	BusinessID      int64  `json:"business_id"`
	// if this object belongs to a topology, like mainline topology,
	// layers means each object's item before this object.
	Layers Layers `json:"layers"`
}

// Basic defines the basic info for a resource.
type Basic struct {
	// the name of the affiliated resource, which could be a model name.
	Type ResourceType `json:"type"`

	// the action that user want to do with this resource.
	// this field should be empty when it's used in resource handle operation.
	Action Action `json:"action"`

	// the name of the resource, which could be a bk-route, etc.
	// this filed is not necessary for all the resources.
	Name string `json:"name"`

	// the instance id of this resource, which could be a model's instance id.
	InstanceID int64

	// InstanceIDEx is a extend for instanceID which can only be integer, but some resources only have string format id.
	InstanceIDEx string
}

// BackendResourceLayer represent one resource layer
type BackendResourceLayer struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id,omitempty"`
	ResourceName string `json:"resource_name"`
}

// BackendResource represent a resource in auth backend, like iam.
type BackendResource []BackendResourceLayer

// CommonInfo contains common field which can be extracted from restful.Request
type CommonInfo struct {
	User UserInfo
}

// Decision TODO
type Decision struct {
	// the authorize decision, whether a user has been authorized or not.
	Authorized bool

	// the detailed reason for this authorize.
	Reason string
}

// ListAuthorizedResourcesParam TODO
type ListAuthorizedResourcesParam struct {
	UserName     string       `json:"user_name"`
	BizID        int64        `json:"bk_biz_id"`
	ResourceType ResourceType `json:"resource_type"`
	Action       Action       `json:"action"`
}

// Action TODO
type Action string

// String 用于打印
func (a Action) String() string {
	return string(a)
}

const (
	// Create TODO
	Create Action = "create"
	// CreateMany TODO
	CreateMany Action = "createMany"
	// Update TODO
	Update Action = "update"
	// UpdateMany TODO
	UpdateMany Action = "updateMany"
	// Delete TODO
	Delete Action = "delete"
	// DeleteMany TODO
	DeleteMany Action = "deleteMany"
	// Archive TODO
	Archive Action = "archive"
	// Find TODO
	Find Action = "find"
	// FindMany TODO
	FindMany Action = "findMany"
	// Unknown action, which is also unsupported actions.
	Unknown Action = "unknown"
	// EmptyAction TODO
	EmptyAction Action = "" // used for register resources
	// SkipAction TODO
	SkipAction Action = "skip"

	// Execute TODO
	Execute Action = "execute"
	// DefaultHostApply TODO
	DefaultHostApply Action = "default"

	// MoveResPoolHostToBizIdleModule TODO
	// move resource pool hosts to a business idle module
	MoveResPoolHostToBizIdleModule Action = "moveResPoolHostToBizIdleModule"
	// MoveResPoolHostToDirectory TODO
	MoveResPoolHostToDirectory Action = "moveResPoolHostToDirectory"
	// AddHostToResourcePool TODO
	AddHostToResourcePool Action = "addHostToResourcePool"
	// MoveBizHostFromModuleToResPool TODO
	MoveBizHostFromModuleToResPool Action = "moveHostFromModuleToResPool"
	// MoveHostToAnotherBizModule TODO
	MoveHostToAnotherBizModule Action = "moveHostToAnotherBizModule"

	// ModelTopologyView TODO
	// system base
	ModelTopologyView Action = "modelTopologyView"
	// ModelTopologyOperation TODO
	ModelTopologyOperation Action = "modelTopologyOperation"

	// WatchHost TODO
	// event watch
	WatchHost Action = "host"
	// WatchHostRelation TODO
	WatchHostRelation Action = "host_relation"
	// WatchBiz TODO
	WatchBiz Action = "biz"
	// WatchSet TODO
	WatchSet Action = "set"
	// WatchModule TODO
	WatchModule Action = "module"
	// WatchProcess TODO
	WatchProcess Action = "process"
	// WatchCommonInstance TODO
	WatchCommonInstance Action = "object_instance"
	// WatchMainlineInstance TODO
	WatchMainlineInstance Action = "mainline_instance"
	// WatchInstAsst TODO
	WatchInstAsst Action = "inst_asst"
	// WatchBizSet TODO
	WatchBizSet Action = "biz_set"
	// WatchPlat watch cloud area event cc action
	WatchPlat Action = "plat"

	// kube related event watch cc actions

	// WatchKubeCluster watch kube cluster event cc action
	WatchKubeCluster Action = "kube_cluster"
	// WatchKubeNode watch kube node event cc action
	WatchKubeNode Action = "kube_node"
	// WatchKubeNamespace watch kube namespace event cc action
	WatchKubeNamespace Action = "kube_namespace"
	// WatchKubeWorkload watch kube workload event cc action
	WatchKubeWorkload Action = "kube_workload"
	// WatchKubePod watch kube pod event cc action
	WatchKubePod Action = "kube_pod"

	// ViewBusinessResource view business related resources action, including business and business collection resources
	ViewBusinessResource Action = "viewBusinessResource"

	// AccessBizSet access business set related resources, including business and business related resources
	AccessBizSet Action = "accessBizSet"
)
