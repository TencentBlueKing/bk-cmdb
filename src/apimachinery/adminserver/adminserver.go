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

package adminserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/apimachinery/util"
)

type AdminServerClientInterface interface {
	ClearDatabase(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	Set(ctx context.Context, ownerID string, h http.Header) (resp *metadata.Response, err error)
	Migrate(ctx context.Context, ownerID string, distribution string, h http.Header) (resp *metadata.Response, err error)
}

func NewAdminServerClientInterface(c *util.Capability, version string) AdminServerClientInterface {
	base := fmt.Sprintf("/migrate/%s", version)
	return &adminServer{
		client: rest.NewRESTClient(c, base),
	}
}

type adminServer struct {
	client rest.ClientInterface
}
