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

package user

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type UserInterface interface {
	AddUserConfig(ctx context.Context, h http.Header, dat *metadata.UserConfig) (resp *metadata.IDResult, err error)
	UpdateUserConfig(ctx context.Context, businessID string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error)
	DeleteUserConfig(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.BaseResp, err error)
	GetUserConfig(ctx context.Context, h http.Header, opt *metadata.QueryInput) (resp *metadata.GetUserConfigResult, err error)
	GetUserConfigDetail(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.GetUserConfigDetailResult, err error)

	AddUserCustom(ctx context.Context, user string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error)
	UpdateUserCustomByID(ctx context.Context, user string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error)
	GetUserCustomByUser(ctx context.Context, user string, h http.Header) (resp *metadata.GetUserCustomResult, err error)
	GetDefaultUserCustom(ctx context.Context, user string, h http.Header) (resp *metadata.GetUserCustomResult, err error)
}

func NewUserInterface(client rest.ClientInterface) UserInterface {
	return &user{client: client}
}

type user struct {
	client rest.ClientInterface
}
