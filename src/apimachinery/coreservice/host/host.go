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

package host

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type HostClientInterface interface {
	TransferToInnerModule(ctx context.Context, h http.Header, input *metadata.TransferHostToInnerModule) (resp *metadata.OperaterException, err error)
	TransferToNormalModule(ctx context.Context, header http.Header, input *metadata.HostsModuleRelation) (resp *metadata.OperaterException, err error)
	TransferToAnotherBusiness(ctx context.Context, header http.Header, input *metadata.TransferHostsCrossBusinessRequest) (resp *metadata.OperaterException, err error)

	RemoveFromModule(ctx context.Context, header http.Header, input *metadata.RemoveHostsFromModuleOption) (resp *metadata.OperaterException, err error)
	DeleteHost(ctx context.Context, header http.Header, input *metadata.DeleteHostRequest) (resp *metadata.OperaterException, err error)

	GetHostModuleRelation(ctx context.Context, header http.Header, input *metadata.HostModuleRelationRequest) (resp *metadata.HostConfig, err error)
	FindIdentifier(ctx context.Context, header http.Header, input *metadata.SearchHostIdentifierParam) (resp *metadata.SearchHostIdentifierResult, err error)

	GetHostByID(ctx context.Context, header http.Header, hostID string) (resp *metadata.HostInstanceResult, err error)
	GetHosts(ctx context.Context, header http.Header, opt *metadata.QueryInput) (resp *metadata.GetHostsResult, err error)
}

func NewHostClientInterface(client rest.ClientInterface) HostClientInterface {
	return &host{client: client}
}

type host struct {
	client rest.ClientInterface
}
