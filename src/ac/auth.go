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

// auth control
package ac

import (
	"context"
	"errors"
	"net/http"

	"configcenter/src/ac/meta"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

var NoAuthorizeError = errors.New("no authorize")

type AuthInterface interface {
	RegisterSystem(ctx context.Context, host string) error
}

type AuthorizeInterface interface {
	AuthorizeBatch(ctx context.Context, h http.Header, user meta.UserInfo, resources ...meta.ResourceAttribute) (
		[]types.Decision, error)

	AuthorizeAnyBatch(ctx context.Context, h http.Header, user meta.UserInfo, resources ...meta.ResourceAttribute) (
		[]types.Decision, error)

	ListAuthorizedResources(ctx context.Context, h http.Header, input meta.ListAuthorizedResourcesParam) ([]string, error)
	GetNoAuthSkipUrl(ctx context.Context, h http.Header, input *metadata.IamPermission) (string, error)
	RegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstanceWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)

	BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstancesWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)
}
