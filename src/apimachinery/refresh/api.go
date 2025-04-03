/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package refresh

import (
	"context"
	"net/http"

	"configcenter/pkg/tenant/types"
)

// RefreshTenant refresh tenant info
func (a *refresh) RefreshTenant(ctx context.Context, h http.Header) ([]types.Tenant, error) {

	resp := new(types.AllTenantsResult)
	subPath := "/refresh/tenant"
	err := a.client.Post().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return resp.Data, nil
}
