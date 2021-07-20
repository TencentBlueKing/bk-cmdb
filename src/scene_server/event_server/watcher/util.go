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
	"errors"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
)

const (
	// which means the the start cursor is not exist error, may be a head cursor.
	startCursorNotExistError = "start cursor not exist error"
)

// GetNodesFromCursor get node start from a cursor, do not contain this cursor's value.
// if start cursor is no event cursor, return events start from the beginning.
func (w *Watcher) GetNodesFromCursor(count int, startCursor string, cursorType watch.CursorType) ([]*watch.ChainNode, error) {
	nodeOpts := &metadata.SearchEventNodesOption{
		Resource:    cursorType,
		StartCursor: startCursor,
		Limit:       count,
	}

	exists, nodes, err := w.cacheCli.Event().SearchFollowingEventChainNodes(w.ctx, w.header, nodeOpts)
	if err != nil {
		blog.Errorf("get latest watch node detail from cache service failed, err: %v", err)
		return nil, err
	}

	if !exists {
		return nil, StartCursorNotExistError
	}

	return nodes, nil
}

var (
	NoEventsError               = errors.New(noEventWarning)
	TailNodeTargetNotExistError = errors.New(tailNodeTargetNotExistError)
	StartCursorNotExistError    = errors.New(startCursorNotExistError)
)

const (
	tailNodeTargetNotExistError = "tail node target detail not exist error"
	noEventWarning              = "no events"
)

// GetLatestEvent get latest event chain node with its detail value.
func (w *Watcher) GetLatestEvent(cursorType watch.CursorType, ) (*watch.ChainNode, error) {
	opts := &metadata.GetLatestEventOption{
		Resource: cursorType,
	}

	node, err := w.cacheCli.Event().GetLatestEvent(w.ctx, w.header, opts)
	if err != nil {
		blog.Errorf("get latest watch node detail from cache service failed, err: %v", err)
		return nil, err
	}

	if !node.ExistsNode {
		return nil, NoEventsError
	}

	return node.Node, nil
}
