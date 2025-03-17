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

package user

import (
	"context"
	"net/http"

	"configcenter/src/thirdparty/apigw/apigwutil"
)

// ClientI is the bk-login api gateway client
type ClientI interface {
	GetTenants(ctx context.Context, h http.Header) ([]*Tenant, error)
}

type user struct {
	service *apigwutil.ApiGWSrv
}

// NewClient create bk-login api gateway client
func NewClient(options *apigwutil.ApiGWOptions) (ClientI, error) {
	service, err := apigwutil.NewApiGW(options, "apiGW.bkUserApiGatewayUrl")
	if err != nil {
		return nil, err
	}

	return &user{
		service: service,
	}, nil
}
