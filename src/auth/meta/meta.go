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

package meta

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

type Affiliated Basic

type Resource struct {
	Basic

	// the action that user want to do with this resource.
	Action Action

	// the business id that this resource belongs to, but it's not necessary for
	// a resource that does not belongs to a business.
	BusinessID int64

	// affiliated resource info
	Affiliated Affiliated
}

// Basic defines the basic info for a resource.
type Basic struct {
	// the name of the affiliated resource, which could be a model name.
	Type ResourceType

	// the name of the resource, which could be a bk-route, etc.
	// this filed is not necessary for all the resources.
	Name string

	// the instance id of this resource, which could be a model's instance id.
	InstanceID int64
}

type Decision struct {
	// the authorize decision, whether a user has been authorized or not.
	Authorized bool

	// the detailed reason for this authorize.
	Reason string
}

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

type Item Basic

type ResourceAttribute struct {
	Basic

	BusinessID int64
	// if this object belongs to a topology, like mainline topology,
	// layers means each object's item before this object.
	Layers []Item
}
