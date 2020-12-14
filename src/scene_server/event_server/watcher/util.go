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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
)

const (
	eventStep = 200
	// which means the the start cursor is not exist error, may be a head cursor.
	startCursorNotExistError = "start cursor not exist error"
)

// GetNodesFromCursor get node start from a cursor, do not contain this cursor's value.
// if cursor is no event cursor, return events start from the beginning.
func (w *Watcher) GetNodesFromCursor(count int, startCursor string, cursorType watch.CursorType) ([]*watch.ChainNode, error) {
	startCursorFilter := make(map[string]interface{})
	if startCursor != watch.NoEventCursor {
		startCursorFilter[common.BKCursorField] = startCursor
	}

	nodes, err := w.GetNodesFromFilter(count, startCursorFilter, cursorType)
	if err != nil {
		blog.Errorf("get latest watch node detail from cache service failed, err: %v", err)
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, StartCursorNotExistError
	}

	if startCursor == watch.NoEventCursor {
		return nodes, nil
	}

	return nodes[1:], nil
}

// GetNodesFromFilter get node start from the node specified by filter, containing this node's value.
func (w *Watcher) GetNodesFromFilter(count int, filter map[string]interface{}, cursorType watch.CursorType) (
	[]*watch.ChainNode, error) {

	nodeOpts := &metadata.SearchFollowingEventNodesOption{
		Resource: cursorType,
		Filter:   filter,
		Sort:     common.BKFieldID,
		Limit:    count,
	}

	nodes, err := w.cacheCli.Event().SearchFollowingEventChainNodes(w.ctx, w.header, nodeOpts)
	if err != nil {
		blog.Errorf("get latest watch node detail from cache service failed, err: %v", err)

		if err.GetCode() == common.CCErrEventChainNodeNotExist {
			return nil, StartCursorNotExistError
		}
		return nil, err
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
func (w *Watcher) GetLatestEvent(cursorType watch.CursorType, withDetail bool) (*metadata.EventNodeWithDetail, error) {
	opts := &metadata.SearchEventNodeOption{
		Resource:   cursorType,
		Filter:     make(map[string]interface{}),
		Sort:       common.BKFieldID + ":-1",
		WithDetail: withDetail,
	}

	node, err := w.cacheCli.Event().SearchEventChainNode(w.ctx, w.header, opts)
	if err != nil {
		blog.Errorf("get latest watch node detail from cache service failed, err: %v", err)

		if err.GetCode() == common.CCErrEventChainNodeNotExist {
			return nil, NoEventsError
		}
		if err.GetCode() == common.CCErrEventDetailNotExist {
			return nil, TailNodeTargetNotExistError
		}
		return nil, err
	}

	return node, nil
}
