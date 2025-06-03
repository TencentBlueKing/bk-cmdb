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

// Package auth TODO
package auth

import (
	"context"
	"errors"
	"net/http"

	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/thirdparty/apigw/iam"
)

// Authorizer TODO
type Authorizer interface {
	// Authorize TODO
	// check if a user's operate resource is already authorized or not.
	Authorize(ctx context.Context, header http.Header, opts *iam.AuthOptions) (*types.Decision, error)

	// AuthorizeBatch TODO
	// check if a user's operate resources is authorized or not batch.
	// Note: being authorized resources must be the same resource.
	AuthorizeBatch(ctx context.Context, header http.Header, opts *iam.AuthBatchOptions) ([]*types.Decision, error)

	// AuthorizeAnyBatch TODO
	// check if a user have any authority of the operate actions batch.
	AuthorizeAnyBatch(ctx context.Context, header http.Header, opts *iam.AuthBatchOptions) ([]*types.Decision, error)

	// ListAuthorizedInstances TODO
	// list a user's all the authorized resource instance list with an action.
	// Note: opts.Resources is not required.
	// the returned list may be huge, we do not do result paging
	ListAuthorizedInstances(ctx context.Context, header http.Header, opts *iam.AuthOptions,
		resourceType iam.IamResourceType) (*iam.AuthorizeList, error)
}

// ResourceFetcher TODO
type ResourceFetcher interface {
	// ListInstancesWithAttributes TODO
	// get "same" resource instances with attributes
	// returned with the resource's instance id list matched with options.
	ListInstancesWithAttributes(ctx context.Context, opts *types.ListWithAttributes) (idList []string, err error)
}

// NewAuth TODO
func NewAuth(fetcher ResourceFetcher) (Authorizer, error) {

	if fetcher == nil {
		return nil, errors.New("fetcher can not be nil")
	}

	return &Authorize{
		iam:     apigwcli.Client().Iam(),
		fetcher: fetcher,
	}, nil
}
