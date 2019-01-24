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

type Authorizer interface {
	// Authorize works to check if a user has the authority to operate resources.
	Authorize(a *Attribute) (authorized Decision, reason string, err error)
}

// ResourceManager is used to handle the resources register to authorize center.
// request id is a identifier for a request, returned by IAM.
type ResourceHandler interface {
	// register a resource
	Register(r *ResourceAttribute) (requestID string, err error)
	// deregister a resource
	Deregister(r *ResourceAttribute) (requestID string, err error)
	// update a resource's info
	Update(r *ResourceAttribute) (requestID string, err error)
	// get a resource's info
	Get() error
}

func NewAuthorizer() (Authorizer, err error) {
	panic("implement me")
}

func NewResourceHandler() (ResourceHandler, error) {
	panic("implement me")
}
