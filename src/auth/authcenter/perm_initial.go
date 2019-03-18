/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authcenter

import (
	"context"
	"net/http"
)

func (ac *AuthCenter) Init(ctx context.Context) error {
	header := http.Header{}
	if err := ac.authClient.RegistSystem(ctx, header, expectSystem); err != nil && err != ErrDuplicated {
		return err
	}

	if err := ac.authClient.UpdateSystem(ctx, header, System{SystemID: expectSystem.SystemID, SystemName: expectSystem.SystemName}); err != nil {
		return err
	}

	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDSystem, expectSystemResourceType); err != nil {
		return err
	}
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDBiz, expectBizResourceType); err != nil {
		return err
	}

	if err := ac.authClient.registerResource(ctx, header, &expectModelGroupResourceInst); err != nil && err != ErrDuplicated {
		return err
	}
	if err := ac.authClient.registerResource(ctx, header, &expectModelResourceInst); err != nil && err != ErrDuplicated {
		return err
	}

	return nil
}
