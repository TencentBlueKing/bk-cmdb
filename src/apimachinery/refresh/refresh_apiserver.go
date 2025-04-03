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
	"fmt"
	"net/http"

	"configcenter/pkg/tenant/types"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

// RefreshClientInterface refresh tenant info, skip tenant verify
type RefreshClientInterface interface {
	RefreshTenant(ctx context.Context, h http.Header) ([]types.Tenant, error)
}

// NewRefreshClientInterface new refresh tenant info client
func NewRefreshClientInterface(c *util.Capability, version string) RefreshClientInterface {
	base := fmt.Sprintf("/refresh/%s", version)
	return &refresh{
		client: rest.NewRESTClient(c, base),
	}
}

type refresh struct {
	client rest.ClientInterface
}
