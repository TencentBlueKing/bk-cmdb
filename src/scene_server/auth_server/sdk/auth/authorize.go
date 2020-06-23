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

	"configcenter/src/scene_server/auth_server/sdk/client"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

type Authorize struct {
	// iam client
	iam client.Interface
	// fetch resource if needed
	fetcher ResourceGetter
}

func (a Authorize) Authorize(ctx context.Context, opts *types.AuthOptions) (types.Decision, error) {
	panic("implement me")
}

func (a Authorize) AuthorizeBatch(ctx context.Context, opts []*types.AuthOptions) ([]types.Decision, error) {
	panic("implement me")
}
