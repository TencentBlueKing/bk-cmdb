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

package event

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
)

type Interface interface {
	GetLatestEvent(ctx context.Context, h http.Header, opts *metadata.GetLatestEventOption) (
		*metadata.EventNode, errors.CCErrorCoder)
	SearchFollowingEventChainNodes(ctx context.Context, h http.Header, opts *metadata.SearchEventNodesOption) (
		bool, []*watch.ChainNode, errors.CCErrorCoder)
	SearchEventDetails(ctx context.Context, h http.Header, opts *metadata.SearchEventDetailsOption) ([]string,
		errors.CCErrorCoder)
	WatchEvent(ctx context.Context, h http.Header, opts *watch.WatchEventOptions) (*string, errors.CCErrorCoder)
}

func NewCacheClient(client rest.ClientInterface) Interface {
	return &eventCache{client: client}
}

type eventCache struct {
	client rest.ClientInterface
}
