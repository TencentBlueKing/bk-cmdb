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

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/scene_server/auth_server/sdk/client"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

type Authorizer interface {
	// check if a user's operate resource is already authorized or not.
	Authorize(ctx context.Context, opts *types.AuthOptions) (*types.Decision, error)

	// check if a user's operate resources is authorized or not batch.
	// Note: being authorized resources must be the same resource.
	AuthorizeBatch(ctx context.Context, opts *types.AuthBatchOptions) ([]*types.Decision, error)

	// check if a user have any authority of the operate actions batch.
	AuthorizeAnyBatch(ctx context.Context, opts *types.AuthBatchOptions) ([]*types.Decision, error)

	// list a user's all the authorized resource instance list with an action.
	// Note: opts.Resources is not required.
	// the returned list may be huge, we do not do result paging
	ListAuthorizedInstances(ctx context.Context, opts *types.AuthOptions, resourceType types.ResourceType) ([]string, error)
}

type ResourceFetcher interface {
	// get "same" resource instances with attributes
	// returned with the resource's instance id list matched with options.
	ListInstancesWithAttributes(ctx context.Context, opts *types.ListWithAttributes) (idList []string, err error)
}

func NewAuth(conf types.Config, fetcher ResourceFetcher) (Authorizer, error) {

	if fetcher == nil {
		return nil, errors.New("fetcher can not be nil")
	}

	// initialize iam client.
	iam, err := client.NewClient(conf.Iam, conf.Options)
	if err != nil {
		return nil, fmt.Errorf("new iam client failed, err: %v", err)
	}

	return &Authorize{
		iam:     iam,
		fetcher: fetcher,
	}, nil
}
