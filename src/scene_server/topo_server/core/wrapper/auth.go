/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wrapper

import (
	"context"

	"configcenter/src/auth"
)

var _ auth.Authorizer = (*AuthAPI)(nil)
var _ auth.ResourceHandler = (*AuthAPI)(nil)

// AuthAPI wrapper API for auth
type AuthAPI struct {
	authorizer      auth.Authorizer
	resourceHandler auth.ResourceHandler
}

// NewAuthAPI return a new auth wrapper
func NewAuthAPI() (AuthAPI, error) {

	authorizer, err := auth.NewAuthorizer()
	if nil != err {
		return AuthAPI{}, err
	}

	resourceHandler, err := auth.NewResourceHandler()
	if nil != err {
		return AuthAPI{}, err
	}

	return AuthAPI{
		authorizer:      authorizer,
		resourceHandler: resourceHandler,
	}, nil
}

// Authorize works to check if a user has the authority to operate resources.
func (w AuthAPI) Authorize(a *auth.Attribute) (authorized auth.Decision, reason string, err error) {
	return w.authorizer.Authorize(a)
}

// Register register a resource
func (w AuthAPI) Register(ctx context.Context, r *auth.ResourceAttribute) (requestID string, err error) {
	return w.resourceHandler.Register(ctx, r)
}

// Deregister deregister a resource
func (w AuthAPI) Deregister(ctx context.Context, r *auth.ResourceAttribute) (requestID string, err error) {
	return w.resourceHandler.Deregister(ctx, r)
}

// Update update a resource's info
func (w AuthAPI) Update(ctx context.Context, r *auth.ResourceAttribute) (requestID string, err error) {
	return w.resourceHandler.Update(ctx, r)
}

// Get get a resource's info
func (w AuthAPI) Get(ctx context.Context) error {
	return w.resourceHandler.Get(ctx)
}
