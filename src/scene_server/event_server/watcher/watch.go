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

package watcher

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/storage/dal/redis"
)

/* eventserver watcher defines, just created base on old service/watch.go */

// Watcher is resource events watcher in eventserver.
type Watcher struct {
	ctx    context.Context
	header http.Header

	// cache is cc redis client.
	cache redis.Client

	// cacheCli is cc cache service client set
	cacheCli cacheservice.Cache
}

// NewWatcher creates a new Watcher object.
func NewWatcher(ctx context.Context, header http.Header, cache redis.Client, cacheCli cacheservice.Cache) *Watcher {
	return &Watcher{ctx: ctx, header: header, cache: cache, cacheCli: cacheCli}
}

// GetEventDetailsWithCursorNodes gets event detail strings base on target hit chain nodes.
func (w *Watcher) GetEventDetailsWithCursorNodes(cursorType watch.CursorType, hitNodes []*watch.ChainNode) (
	[]string, error) {

	if len(hitNodes) == 0 {
		return make([]string, 0), nil
	}

	cursors := make([]string, len(hitNodes))
	for index, node := range hitNodes {
		cursors[index] = node.Cursor
	}

	detailOpts := &metadata.SearchEventDetailsOption{
		Resource: cursorType,
		Cursors:  cursors,
	}

	details, err := w.cacheCli.Event().SearchEventDetails(w.ctx, w.header, detailOpts)
	if err != nil {
		blog.Errorf("search event details failed, err: %v, cursors: %+v", err, cursors)
		return nil, err
	}

	return details, nil
}

func (w *Watcher) GetHitNodeWithEventType(nodes []*watch.ChainNode, typs []watch.EventType) []*watch.ChainNode {
	if len(typs) == 0 {
		return nodes
	}

	if len(nodes) == 0 {
		return nodes
	}

	m := make(map[watch.EventType]bool)
	for _, t := range typs {
		m[t] = true
	}

	hitNodes := make([]*watch.ChainNode, 0)
	for _, node := range nodes {
		_, hit := m[node.EventType]
		if hit {
			hitNodes = append(hitNodes, node)
			continue
		}
	}
	return hitNodes
}

// ResetRequestID reset a new request id for watcher's header
func (w *Watcher) ResetRequestID() {
	w.header.Set(common.BKHTTPCCRequestID, util.GenerateRID())
}

// GetRid get request id from watcher's header
func (w *Watcher) GetRid() string {
	return w.header.Get(common.BKHTTPCCRequestID)
}
