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

type Attribute struct {
	// the version of this resource, which is the api version.
	APIVersion string
	Resources  []Resource
	User       UserInfo
}

type UserInfo struct {
	// the name of this user.
	UserName string
	// the supplier id that this user belongs to.
	SupplierID string
}

type Resource struct {
	// the name of the resource, which could be a model name.
	Name ResourceType

	// the instance id of this resource, which could be a model's instance id.
	InstanceID uint64

	// the action that user want to do with this resource.
	Action Action

	// the business id that this resource belongs to, but it's not necessary for
	// a resource that does not belongs to a business.
	BusinessID uint64

	// affiliated resource info
	Affiliated Affiliated
}

type Affiliated struct {
	// the name of the affiliated resource, which could be a model name.
	Name ResourceType
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
	// unknown action, which is also unsupported actions.
	Unknown Action = "unknown"
)

type ResourceAttribute struct {
	// object's id
	Object string
	// object's name
	// it's not be needed when it's used to deregister a resource.
	ObjectName string
	// if this object belongs to a topology, like mainline topology,
	// layers means each object's item before this object.
	Layers []Item
}

type Item struct {
	// object's id
	Object string
	// this object's instance id
	InstanceID int64
}
