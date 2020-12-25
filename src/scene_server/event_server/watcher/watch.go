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
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/stream/types"

	rawRedis "github.com/go-redis/redis/v7"
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
	ctx context.Context

	// cache is cc redis client.
	cache redis.Client
}

// NewWatcher creates a new Watcher object.
func NewWatcher(ctx context.Context, cache redis.Client) *Watcher {
	return &Watcher{ctx: ctx, cache: cache}
}

// WatchWithStartFrom watches target resource base on timestamp.
func (w *Watcher) WatchWithStartFrom(key event.Key, opts *watch.WatchEventOptions, rid string) ([]*watch.WatchEventDetail, error) {
	// validate start from value is in the range or not
	headTarget, tailTarget, err := w.GetHeadTailNodeTargetNode(key)
	if err != nil {
		blog.Errorf("get head and tail targeted node detail failed, err: %v, rid: %s", err, rid)

		// tail node is not initialized, which means no events.
		if err == TailNodeNotExistError {
			return []*watch.WatchEventDetail{{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}}, nil
		}

		return nil, err
	}

	// not one event occurs.
	if headTarget.NextCursor == key.TailKey() || tailTarget.NextCursor == key.HeadKey() {
		// validate start from time with key's ttl
		diff := time.Now().Unix() - opts.StartFrom
		if diff < 0 || diff > key.TTLSeconds() {
			// this is invalid.
			return nil, errors.New("bk_start_from value is out of range")
		}
	}

	// start from is too old, not allowed.
	if int64(headTarget.ClusterTime.Sec) > opts.StartFrom {
		return nil, errors.New("bk_start_from value is too small")
	}

	// start from is ahead of the latest's event time, watch from now.
	if int64(tailTarget.ClusterTime.Sec) < opts.StartFrom {

		latestEvent, err := w.WatchFromNow(key, opts, rid)
		if err != nil {
			blog.Errorf("watch with start from: %d, result in watch from now, get latest event failed, err: %v, rid: %s",
				opts.StartFrom, err, rid)
			return nil, err
		}

		return []*watch.WatchEventDetail{latestEvent}, nil
	}

	// keep scan the cursor chain until to the tail cursor.
	// start from the head key.
	nextCursor := key.HeadKey()
	timeout := time.After(timeoutWatchLoopSeconds * time.Second)
	for {
		select {
		case <-timeout:
			// scan the event's too long time, need to exist immediately.
			blog.Errorf("watch with start from: %d, scan the cursor chain, but scan too long time, rid: %s", opts.StartFrom, rid)
			return nil, errors.New("scan the event cost too long time")
		default:

		}

		// scan event node from head, returned nodes does not contain tail node.
		nodes, err := w.GetNodesFromCursor(eventStep, nextCursor, key)
		if err != nil {
			blog.Errorf("get event from head failed, err: %v, rid: %s", err, rid)
			if err == HeadNodeNotExistError {
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
			resp := &watch.WatchEventDetail{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}

			// at least the tail node should can be scan, so something goes wrong.
			blog.V(5).Infof("watch with start from %s, but no event found in the chain, rid: %s", opts.StartFrom, rid)
			return []*watch.WatchEventDetail{resp}, nil
		}

		hitNodes := w.GetHitNodeWithEventType(nodes, opts.EventTypes)
		matchedNodes := make([]*watch.ChainNode, 0)
		for _, node := range hitNodes {
			// find node that cluster time is larger than the start from seconds.
			if int64(node.ClusterTime.Sec) >= opts.StartFrom {
				matchedNodes = append(matchedNodes, node)
			}
		}

		if len(matchedNodes) != 0 {
			// matched event has been found, get them all.
			return w.GetEventsWithCursorNodes(opts, matchedNodes, key, rid)
		}

		// not even one is hit.
		// check if nodes has already scan to the end
		lastNode := nodes[len(nodes)-1]
		if lastNode.NextCursor == key.TailKey() {
			resp := &watch.WatchEventDetail{
				Cursor:   lastNode.Cursor,
				Resource: opts.Resource,
				Detail:   nil,
			}
			return []*watch.WatchEventDetail{resp}, nil
		}

		// update nextCursor and do next scan round.
		nextCursor = lastNode.Cursor
	}
}

func (w *Watcher) GetEventsWithCursorNodes(opts *watch.WatchEventOptions, hitNodes []*watch.ChainNode,
	key event.Key, rid string) ([]*watch.WatchEventDetail, error) {

	results := make([]*rawRedis.StringCmd, 0)
	pipe := w.cache.Pipeline()
	for _, node := range hitNodes {
		if node.Cursor == key.TailKey() {
			continue
		}
		results = append(results, pipe.Get(key.DetailKey(node.Cursor)))
	}

	// cursor is end to tail node.
	if len(results) == 0 {
		return make([]*watch.WatchEventDetail, 0), nil
	}

	_, err := pipe.Exec()
	if err != nil {
		blog.ErrorJSON("watch with start from: %s, resource: %s, hit events, but get event detail failed, hit nodes: %s, err: %v, rid: %s",
			opts.StartFrom, opts.Resource, hitNodes, err, rid)
		return nil, err
	}
	resp := make([]*watch.WatchEventDetail, 0)
	for idx, result := range results {
		jsonStr := types.GetEventDetail(result.Val())

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
func (w *Watcher) GetEventDetailsWithCursorNodes(hitNodes []*watch.ChainNode, key event.Key) ([]string, error) {
	results := make([]*rawRedis.StringCmd, 0)

	pipe := w.cache.Pipeline()
	for _, node := range hitNodes {
		if node.Cursor == key.TailKey() {
			continue
		}
		results = append(results, pipe.Get(key.DetailKey(node.Cursor)))
	}

	if len(results) == 0 {
		return []string{}, nil
	}

	if _, err := pipe.Exec(); err != nil {
		blog.ErrorJSON("get event detail strings failed, hit nodes: %s, err: %+v", hitNodes, err)
		return nil, err
	}

	resp := []string{}
	for _, result := range results {
		resp = append(resp, result.Val())
	}

	return resp, nil
}

// WatchFromNow watches target resource events from now.
func (w *Watcher) WatchFromNow(key event.Key, opts *watch.WatchEventOptions, rid string) (*watch.WatchEventDetail, error) {
	node, tailTarget, err := w.GetLatestEventDetail(key)
	if err != nil {
		blog.Errorf("watch from now, but get latest event failed, key, err: %v, rid: %s", err, rid)

		if err == TailNodeNotExistError || err == NoEventsError {
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

	jsonStr := types.GetEventDetail(tailTarget)
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
	if startCursor == watch.NoEventCursor {
		// user got no events because of no event occurs in the system in the previous watch around,
		// we should watch from the head cursor in this round, so that user can not miss any events.
		startCursor = key.HeadKey()
	}

	start := time.Now().Unix()
	for {
		nodes, err := w.GetNodesFromCursor(eventStep, startCursor, key)
		if err != nil {
			blog.Errorf("watch event from cursor: %s, but get cursors failed, err: %v, rid: %s", opts.Cursor, err, rid)

			if err == HeadNodeNotExistError {

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
			if hitNodes[0].Cursor == key.TailKey() {
				// to the end
				resp := &watch.WatchEventDetail{
					Cursor:    watch.NoEventCursor,
					Resource:  opts.Resource,
					EventType: "",
					Detail:    nil,
				}

				// at least the tail node should can be scan, so something goes wrong.
				blog.V(5).Infof("watch with cursor %s, but no events found in the chain, rid: %s", opts.Cursor, rid)
				return []*watch.WatchEventDetail{resp}, nil
			}

			// matched event has been found, get them all.
			blog.V(5).Infof("watch key: %s with resource: %s, hit events, return immediately. rid: %s", key.Namespace(), opts.Resource, rid)
			return w.GetEventsWithCursorNodes(opts, hitNodes, key, rid)
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
