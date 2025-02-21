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

// Package tenant defines tenant related logics
package tenant

import (
	"context"
	"net/http"

	"configcenter/pkg/tenant/types"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
)

// TenantClientInterface tenant client interface
type TenantClientInterface interface {
	GetAllTenants(ctx context.Context, header http.Header) ([]types.Tenant, errors.CCErrorCoder)
}

// New new tenant client interface
func New(client rest.ClientInterface) TenantClientInterface {
	return &tenant{client: client}
}

type tenant struct {
	client rest.ClientInterface
}
