/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	iamtypes "configcenter/src/ac/iam/types"
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

	plan, err := a.iam.HybridPlan(ctx, header, &iam.PlanReq{
		Subject: iam.AuthSubject{
			Type: iam.SubjectType(opts.Subject.Type),
			ID:   opts.Subject.ID,
		},
		ActionID: iamtypes.ActionID(opts.Action.ID),
	})
	if err != nil {
		return nil, err
	}

	authorized, err := a.calculatePolicy(ctx, opts.Resources, plan)
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

	plans, err := a.listUserPlanBatchWithCompress(ctx, header, opts)
	if err != nil {
		return nil, fmt.Errorf("list user plan failed, err: %v", err)
	}

	var hitError error
	decisions := make([]*types.Decision, len(opts.Batch))

	pipe := make(chan struct{}, 50)
	wg := sync.WaitGroup{}
	for idx, b := range opts.Batch {
		wg.Add(1)

		pipe <- struct{}{}
		go func(idx int, resources []iam.Resource, plan *operator.Plan) {
			defer func() {
				wg.Done()
				<-pipe
			}()

			var authorized bool
			var err error
			if exact {
				authorized, err = a.calculatePolicy(ctx, resources, plan)
			} else {
				authorized, err = a.calculateAnyPolicy(ctx, resources, plan)
			}
			if err != nil {
				hitError = err
				return
			}

			// save the result with index
			decisions[idx] = &types.Decision{Authorized: authorized}
		}(idx, b.Resources, plans[idx])
	}
	// wait all the policy are calculated
	wg.Wait()

	if hitError != nil {
		return nil, fmt.Errorf("batch calculate policy failed, err: %v", hitError)
	}

	return decisions, nil
}

func (a *Authorize) listUserPlanBatchWithCompress(ctx context.Context, header http.Header,
	opts *iam.AuthBatchOptions) ([]*operator.Plan, error) {

	// because these resource are the same, so we can unique the action id,
	// so that we can cut off the request to iam, and improve the performance.
	actionIDMap := make(map[string]struct{})
	for _, b := range opts.Batch {
		actionIDMap[b.Action.ID] = struct{}{}
	}

	actionIDs := make([]iamtypes.ActionID, 0, len(actionIDMap))
	for actionID := range actionIDMap {
		actionIDs = append(actionIDs, iamtypes.ActionID(actionID))
	}

	plans, err := a.iam.HybridPlanByActions(ctx, header, &iam.PlanByActionsReq{
		Subject: iam.AuthSubject{
			Type: iam.SubjectType(opts.Subject.Type),
			ID:   opts.Subject.ID,
		},
		ActionIDs: actionIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("list user's plan failed, err: %s", err)
	}

	planMap := make(map[string]*operator.Plan)
	for i := range plans {
		plan := plans[i].Plan
		planMap[string(plans[i].ActionID)] = &plan
	}

	allPlans := make([]*operator.Plan, len(opts.Batch))
	for idx, b := range opts.Batch {
		plan, exist := planMap[b.Action.ID]
		if !exist {
			return nil, fmt.Errorf("list user's auth plan, but can not find action id %s in response", b.Action.ID)
		}
		allPlans[idx] = plan
	}

	return allPlans, nil
}

// ListAuthorizedInstances list a user's all the authorized resource instance list with an action.
func (a *Authorize) ListAuthorizedInstances(ctx context.Context, header http.Header, opts *iam.AuthOptions,
	resourceType iam.IamResourceType) (*iam.AuthorizeList, error) {

	plan, err := a.iam.HybridPlan(ctx, header, &iam.PlanReq{
		Subject: iam.AuthSubject{
			Type: iam.SubjectType(opts.Subject.Type),
			ID:   opts.Subject.ID,
		},
		ActionID: iamtypes.ActionID(opts.Action.ID),
	})
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.Kind == operator.AlwaysDeniedKind || plan.Kind == "" {
		return &iam.AuthorizeList{}, nil
	}
	return a.countPolicy(ctx, plan, resourceType)
}
