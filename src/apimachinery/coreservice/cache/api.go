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

package cache

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/cache/topo_tree"
)

type Interface interface {
	SearchTopologyTree(ctx context.Context, h http.Header, opt *topo_tree.SearchOption) ([]topo_tree.Topology, error)
	SearchHostWithInnerIP(ctx context.Context, h http.Header, opt *metadata.SearchHostWithInnerIPOption) (jsonString string, err error)
	SearchHostWithHostID(ctx context.Context, h http.Header, opt *metadata.SearchHostWithIDOption) (jsonString string, err error)
	ListHostWithHostID(ctx context.Context, h http.Header, opt *metadata.ListHostWithIDOption) (jsonString string, err error)
	SearchBusiness(ctx context.Context, h http.Header, bizID int64) (jsonString string, err error)
	SearchSet(ctx context.Context, h http.Header, setID int64) (jsonString string, err error)
	SearchModule(ctx context.Context, h http.Header, moduleID int64) (jsonString string, err error)
	SearchCustomLayer(ctx context.Context, h http.Header, objID string, instID int64) (jsonString string, err error)
}

func NewCacheClient(client rest.ClientInterface) Interface {
	return &baseCache{client: client}
}

type baseCache struct {
	client rest.ClientInterface
}
