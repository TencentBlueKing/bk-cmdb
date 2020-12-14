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
	"errors"
	"net/http"
	"time"

	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/stream/types"
)

/* eventserver watcher defines, just created base on old service/watch.go */

const (
	// timeoutWatchLoopSeconds 25s timeout
	timeoutWatchLoopSeconds = 25

	// loopInternal watch loop internal duration
	loopInternal = 250 * time.Millisecond
)

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

// WatchWithStartFrom watches target resource base on timestamp.
func (w *Watcher) WatchWithStartFrom(key event.Key, opts *watch.WatchEventOptions, rid string) ([]*watch.WatchEventDetail, error) {
	// validate start from time with key's ttl
	diff := time.Now().Unix() - opts.StartFrom
	if diff < 0 || diff > key.TTLSeconds() {
		return nil, errors.New("bk_start_from value is out of range")
	}

	// validate if start from value is in the range or not
	tailTarget, err := w.GetLatestEvent(opts.Resource, true)
	if err != nil {
		blog.Errorf("get head and tail targeted node detail failed, err: %v, rid: %s", err, rid)

		// no events.
		if err == NoEventsError {
			return []*watch.WatchEventDetail{{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}}, nil
		}

		return nil, err
	}

	// start from is ahead of the latest's event time, watch from now.
	tailNode := tailTarget.Node
	if int64(tailNode.ClusterTime.Sec) <= opts.StartFrom {
		hit := w.GetHitNodeWithEventType([]*watch.ChainNode{tailNode}, opts.EventTypes)
		if len(hit) == 0 {
			// not matched, set to no event cursor with empty detail
			return []*watch.WatchEventDetail{{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}}, nil
		}

		jsonStr := types.GetEventDetail(tailTarget.Detail)
		cut := json.CutJsonDataWithFields(&jsonStr, opts.Fields)
		// matched the event type.
		return []*watch.WatchEventDetail{{
			Cursor:    tailNode.Cursor,
			Resource:  opts.Resource,
			EventType: tailNode.EventType,
			Detail:    watch.JsonString(*cut),
		}}, nil
	}

	// find nodes whose cluster time is larger than the start from seconds with length of eventStep.
	filter := map[string]interface{}{
		common.BKClusterTimeField: map[string]interface{}{common.BKDBGT: time.Unix(opts.StartFrom, 0)},
	}

	timeout := time.After(timeoutWatchLoopSeconds * time.Second)
	for {
		select {
		case <-timeout:
			// scan the event's too long time, need to exist immediately.
			blog.Errorf("watch with start from: %d, scan the cursor chain, but scan too long time, rid: %s", opts.StartFrom, rid)
			return nil, errors.New("scan the event cost too long time")
		default:

		}

		nodes, err := w.GetNodesFromFilter(eventStep, filter, opts.Resource)
		if err != nil {
			blog.ErrorJSON("get event failed, err: %s, rid: %s, filter: %s", err, rid, filter)
			return nil, err
		}

		// nodes has already scan to the end
		if len(nodes) == 0 {
			resp := &watch.WatchEventDetail{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}
			// at least the tail node should can be scan, so something goes wrong.
			blog.V(5).Infof("watch with start from %s, but no event found in the chain, rid: %s",
				opts.StartFrom, rid)
			return []*watch.WatchEventDetail{resp}, nil
		}

		hitNodes := w.GetHitNodeWithEventType(nodes, opts.EventTypes)

		if len(hitNodes) != 0 {
			// matched event has been found, get them all.
			return w.GetEventsWithCursorNodes(opts, hitNodes, rid)
		}

		// update filter and do next scan round.
		filter = map[string]interface{}{
			common.BKCursorField: nodes[len(nodes)-1].Cursor,
		}
	}
}

func (w *Watcher) GetEventsWithCursorNodes(opts *watch.WatchEventOptions, hitNodes []*watch.ChainNode, rid string) (
	[]*watch.WatchEventDetail, error) {

	if len(hitNodes) == 0 {
		return make([]*watch.WatchEventDetail, 0), nil
	}

	details, err := w.GetEventDetailsWithCursorNodes(opts.Resource, hitNodes)
	if err != nil {
		blog.Errorf("search event details failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	resp := make([]*watch.WatchEventDetail, 0)
	for idx, jsonStr := range details {
		cut := json.CutJsonDataWithFields(&jsonStr, opts.Fields)
		resp = append(resp, &watch.WatchEventDetail{
			Cursor:    hitNodes[idx].Cursor,
			Resource:  opts.Resource,
			EventType: hitNodes[idx].EventType,
			Detail:    watch.JsonString(*cut),
		})
	}
	return resp, nil
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

// WatchFromNow watches target resource events from now.
func (w *Watcher) WatchFromNow(key event.Key, opts *watch.WatchEventOptions, rid string) (*watch.WatchEventDetail, error) {
	nodeWithDetail, err := w.GetLatestEvent(opts.Resource, true)
	if err != nil {
		blog.Errorf("watch from now, but get latest event failed, key, err: %v, rid: %s", err, rid)

		if err == NoEventsError {
			// event chain list is empty, which means no event and not be initialized.
			return &watch.WatchEventDetail{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}, nil
		}

		return nil, err
	}
	node := nodeWithDetail.Node

	hit := w.GetHitNodeWithEventType([]*watch.ChainNode{node}, opts.EventTypes)
	if len(hit) == 0 {
		// not matched, set to no event cursor with empty detail
		return &watch.WatchEventDetail{
			Cursor:    watch.NoEventCursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}, nil
	}

	jsonStr := types.GetEventDetail(nodeWithDetail.Detail)
	cut := json.CutJsonDataWithFields(&jsonStr, opts.Fields)
	// matched the event type.
	return &watch.WatchEventDetail{
		Cursor:    node.Cursor,
		Resource:  opts.Resource,
		EventType: node.EventType,
		Detail:    watch.JsonString(*cut),
	}, nil
}

// watchWithCursor get events with the start cursor which is offered by user.
// it will hold the request for timeout seconds if no matched event is hit.
// if event has been hit in a round, then events will be returned immediately.
// if no events hit, then will loop the event every 200ms until timeout and return
// with a special cursor named "NoEventCursor", then we will help the user watch
// event from the head cursor.
func (w *Watcher) WatchWithCursor(key event.Key, opts *watch.WatchEventOptions, rid string) ([]*watch.WatchEventDetail, error) {
	startCursor := opts.Cursor

	start := time.Now().Unix()
	for {
		nodes, err := w.GetNodesFromCursor(eventStep, startCursor, opts.Resource)
		if err != nil {
			blog.Errorf("watch event from cursor: %s, but get cursors failed, err: %v, rid: %s", opts.Cursor, err, rid)

			if err == NoEventsError {
				resp := &watch.WatchEventDetail{
					Cursor:    watch.NoEventCursor,
					Resource:  opts.Resource,
					EventType: "",
					Detail:    nil,
				}

				return []*watch.WatchEventDetail{resp}, nil
			}

			return nil, err
		}

		if len(nodes) == 0 {

			if time.Now().Unix()-start > timeoutWatchLoopSeconds {
				// has already looped for timeout seconds, and we still got one event.
				// return with NoEventCursor and empty detail
				resp := &watch.WatchEventDetail{
					Cursor:    watch.NoEventCursor,
					Resource:  opts.Resource,
					EventType: "",
					Detail:    nil,
				}

				// at least the tail node should can be scan, so something goes wrong.
				blog.V(5).Infof("watch with cursor %s, timeout and no event found in the chain, rid: %s", opts.Cursor, rid)
				return []*watch.WatchEventDetail{resp}, nil
			}

			// we got not event one event, sleep a little, and then try to continue the loop watch
			time.Sleep(loopInternal)
			blog.V(5).Infof("watch key: %s with resource: %s, got nothing, try next round. rid: %s", key.Namespace(), opts.Resource, rid)
			continue
		}

		hitNodes := w.GetHitNodeWithEventType(nodes, opts.EventTypes)
		if len(hitNodes) != 0 {
			// matched event has been found, get them all.
			blog.V(5).Infof("watch key: %s with resource: %s, hit events, return immediately. rid: %s", key.Namespace(), opts.Resource, rid)
			return w.GetEventsWithCursorNodes(opts, hitNodes, rid)
		}

		if time.Now().Unix()-start > timeoutWatchLoopSeconds {
			// no event is hit, but timeout, we return the last event cursor with nil detail
			// because it's not what the use want, return the last cursor to help user can
			// watch from here later for next watch round.
			lastNode := nodes[len(nodes)-1]
			resp := &watch.WatchEventDetail{
				Cursor:   lastNode.Cursor,
				Resource: opts.Resource,
				Detail:   nil,
			}

			// at least the tail node should can be scan, so something goes wrong.
			blog.V(5).Infof("watch with cursor %s, but no event matched in the chain, rid: %s", opts.Cursor, rid)
			return []*watch.WatchEventDetail{resp}, nil
		}
		// not event one event is hit, sleep a little, and then try to continue the loop watch
		time.Sleep(loopInternal)
		blog.V(5).Infof("watch key: %s with resource: %s, hit nothing, try next round. rid: %s", key.Namespace(), opts.Resource, rid)
		continue
	}
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
