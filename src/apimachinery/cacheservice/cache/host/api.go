/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package host

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

// Interface TODO
type Interface interface {
	SearchHostWithInnerIPForStatic(ctx context.Context, h http.Header, opt *metadata.SearchHostWithInnerIPOption) (
		jsonString string, err error)
	SearchHostWithAgentID(ctx context.Context, h http.Header, opt *metadata.SearchHostWithAgentID) (
		jsonString string, err error)
	SearchHostWithHostID(ctx context.Context, h http.Header, opt *metadata.SearchHostWithIDOption) (
		jsonString string, err error)
	ListHostWithHostID(ctx context.Context, h http.Header, opt *metadata.ListWithIDOption) (
		jsonString string, err error)
	ListHostWithPage(ctx context.Context, h http.Header, opt *metadata.ListHostWithPage) (
		cnt int64, jsonString string, err error)
}

// NewCacheClient TODO
func NewCacheClient(client rest.ClientInterface) Interface {
	return &baseCache{client: client}
}

type baseCache struct {
	client rest.ClientInterface
}
