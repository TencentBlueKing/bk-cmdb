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

package authserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

type AuthServerClientInterface interface {
	AuthorizeBatch(ctx context.Context, h http.Header, input *types.AuthBatchOptions) ([]types.Decision, error)
	AuthorizeAnyBatch(ctx context.Context, h http.Header, input *types.AuthBatchOptions) ([]types.Decision, error)
	ListAuthorizedResources(ctx context.Context, h http.Header, input meta.ListAuthorizedResourcesParam) ([]string, error)
	GetNoAuthSkipUrl(ctx context.Context, h http.Header, input *metadata.IamPermission) (string, error)
	RegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstanceWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)
	BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstancesWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)
}

func NewAuthServerClientInterface(c *util.Capability, version string) AuthServerClientInterface {
	base := fmt.Sprintf("/ac/%s", version)
	return &authServer{
		client: rest.NewRESTClient(c, base),
	}
}

type authServer struct {
	client rest.ClientInterface
}
