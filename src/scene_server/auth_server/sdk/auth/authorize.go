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
	"net/http"
	"sync"

	"configcenter/src/scene_server/auth_server/sdk/client"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/thirdparty/apigw/iam"
)

// Authorize TODO
type Authorize struct {
	// iam client
	iam client.Interface
	// fetch resource if needed
	fetcher ResourceFetcher
}

// Authorize TODO
func (a *Authorize) Authorize(ctx context.Context, header http.Header, opts *iam.AuthOptions) (*types.Decision, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// find user's policy with action
	getOpt := iam.GetPolicyOption{
		System:  opts.System,
		Subject: opts.Subject,
		Action:  opts.Action,
		// do not use user's policy, so that we can get all the user's policy.
		Resources: make([]iam.Resource, 0),
	}

	policy, err := a.iam.GetUserPolicy(ctx, header, &getOpt)
	if err != nil {
		return nil, err
	}

	authorized, err := a.calculatePolicy(ctx, opts.Resources, policy)
	if err != nil {
		return nil, fmt.Errorf("calculate user's auth policy failed, err: %v", err)
	}

	return &types.Decision{Authorized: authorized}, nil
}

// AuthorizeBatch TODO
func (a *Authorize) AuthorizeBatch(ctx context.Context, header http.Header, opts *iam.AuthBatchOptions) (
	[]*types.Decision, error) {

	return a.authorizeBatch(ctx, header, opts, true)
}

// AuthorizeAnyBatch TODO
func (a *Authorize) AuthorizeAnyBatch(ctx context.Context, header http.Header,
	opts *iam.AuthBatchOptions) ([]*types.Decision, error) {

	return a.authorizeBatch(ctx, header, opts, false)
}

func (a *Authorize) authorizeBatch(ctx context.Context, header http.Header, opts *iam.AuthBatchOptions,
	exact bool) ([]*types.Decision, error) {

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if len(opts.Batch) == 0 {
		return nil, errors.New("no resource instance need to authorize")
	}

	policies, err := a.listUserPolicyBatchWithCompress(ctx, header, opts)
	if err != nil {
		return nil, fmt.Errorf("list user policy failed, err: %v", err)
	}

	var hitError error
	decisions := make([]*types.Decision, len(opts.Batch))

	pipe := make(chan struct{}, 50)
	wg := sync.WaitGroup{}
	for idx, b := range opts.Batch {
		wg.Add(1)

		pipe <- struct{}{}
		go func(idx int, resources []iam.Resource, policy *operator.Policy) {
			defer func() {
				wg.Done()
				<-pipe
			}()

			var authorized bool
			var err error
			if exact {
				authorized, err = a.calculatePolicy(ctx, resources, policy)
			} else {
				authorized, err = a.calculateAnyPolicy(ctx, resources, policy)
			}
			if err != nil {
				hitError = err
				return
			}

			// save the result with index
			decisions[idx] = &types.Decision{Authorized: authorized}
		}(idx, b.Resources, policies[idx])
	}
	// wait all the policy are calculated
	wg.Wait()

	if hitError != nil {
		return nil, fmt.Errorf("batch calculate policy failed, err: %v", hitError)
	}

	return decisions, nil
}

func (a *Authorize) listUserPolicyBatchWithCompress(ctx context.Context, header http.Header,
	opts *iam.AuthBatchOptions) ([]*operator.Policy,
	error) {

	// because these resource are the same, so we can unique the action id,
	// so that we can cut off the request to iam, and improve the performance.
	actionIDMap := make(map[string]iam.Action)
	for _, b := range opts.Batch {
		actionIDMap[b.Action.ID] = b.Action
	}

	actions := make([]iam.Action, 0)
	for _, action := range actionIDMap {
		actions = append(actions, action)
	}

	listOpts := &iam.ListPolicyOptions{
		System:  opts.System,
		Subject: opts.Subject,
		Actions: actions,
		// get all policies with these actions
		Resources: nil,
	}

	policies, err := a.iam.ListUserPolicies(ctx, header, listOpts)
	if err != nil {
		return nil, fmt.Errorf("list user's policy failed, err: %s", err)
	}

	policyMap := make(map[string]*operator.Policy)
	for _, p := range policies {
		policyMap[p.Action.ID] = p.Policy
	}

	allPolicies := make([]*operator.Policy, len(opts.Batch))
	for idx, b := range opts.Batch {
		policy, exist := policyMap[b.Action.ID]
		if !exist {
			return nil, fmt.Errorf("list user's auth policy, but can not find action id %s in response", b.Action.ID)
		}
		allPolicies[idx] = policy
	}

	return allPolicies, nil
}

// ListAuthorizedInstances list a user's all the authorized resource instance list with an action.
func (a *Authorize) ListAuthorizedInstances(ctx context.Context, header http.Header, opts *iam.AuthOptions,
	resourceType iam.IamResourceType) (*iam.AuthorizeList, error) {

	// find user's policy with action
	getOpt := iam.GetPolicyOption{
		System:  opts.System,
		Subject: opts.Subject,
		Action:  opts.Action,
		// do not use user's policy, so that we can get all the user's policy.
		Resources: opts.Resources,
	}
	policy, err := a.iam.GetUserPolicy(ctx, header, &getOpt)
	if err != nil {
		return nil, err
	}
	if policy == nil || policy.Operator == "" {
		return &iam.AuthorizeList{}, nil
	}
	return a.countPolicy(ctx, policy, resourceType)
}
