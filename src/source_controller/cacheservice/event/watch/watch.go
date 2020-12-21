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

package watch

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/stream/types"
)

/* eventserver watcher defines, just created base on old service/watch.go */

const (
	// timeoutWatchLoopSeconds 25s timeout
	timeoutWatchLoopSeconds = 25

	// loopInternal watch loop internal duration
	loopInternal = 250 * time.Millisecond

	// the number of events to read in one step TODO: increase this later
	eventStep = 200
)

// WatchWithStartFrom watches target resource base on timestamp.
func (c *Client) WatchWithStartFrom(kit *rest.Kit, key event.Key, opts *watch.WatchEventOptions) (
	[]*watch.WatchEventDetail, error) {

	rid := kit.Rid

	// validate start from time with key's ttl
	diff := time.Now().Unix() - opts.StartFrom
	if diff < 0 || diff > key.TTLSeconds() {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_start_from")
	}

	// validate if start from value is in the range or not
	tailNode, exists, err := c.getLatestEvent(kit, key)
	if err != nil {
		blog.Errorf("get head and tail targeted node detail failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	// no events
	if !exists {
		return []*watch.WatchEventDetail{{
			Cursor:    watch.NoEventCursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}}, nil
	}

	// start from is ahead of the latest's event time, watch from now.
	if int64(tailNode.ClusterTime.Sec) <= opts.StartFrom {
		hit := c.GetHitNodeWithEventType([]*watch.ChainNode{tailNode}, opts.EventTypes)
		if len(hit) == 0 {
			// not matched, set to no event cursor with empty detail
			return []*watch.WatchEventDetail{{
				Cursor:    watch.NoEventCursor,
				Resource:  opts.Resource,
				EventType: "",
				Detail:    nil,
			}}, nil
		}

		detail, exists, err := c.getEventDetail(kit, tailNode, opts.Fields, key)
		if err != nil {
			blog.Errorf("get latest event detail failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		if !exists {
			return nil, kit.CCError.CCError(common.CCErrEventDetailNotExist)
		}

		// matched the event type.
		return []*watch.WatchEventDetail{{
			Cursor:    tailNode.Cursor,
			Resource:  opts.Resource,
			EventType: tailNode.EventType,
			Detail:    watch.JsonString(detail),
		}}, nil
	}

	// find the first node with a larger cluster time than the start from parameter
	filter := map[string]interface{}{
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGT: metadata.Time{Time: time.Unix(opts.StartFrom, 0).Local()},
		},
	}

	node := new(watch.ChainNode)
	err = c.watchDB.Table(key.ChainCollection()).Find(filter).Sort(common.BKFieldID).One(kit.Ctx, node)
	if err != nil {
		blog.ErrorJSON("get chain node from mongo failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
		if !c.watchDB.IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}

		return []*watch.WatchEventDetail{{
			Cursor:    watch.NoEventCursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}}, nil
	}

	timeout := time.After(timeoutWatchLoopSeconds * time.Second)
	isFirstLoop := true
	for {
		select {
		case <-timeout:
			// scan the event's too long time, need to exist immediately.
			blog.Errorf("watch with start from: %d, scan the cursor chain, but scan too long time, rid: %s", opts.StartFrom, rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_start_from")
		default:

		}

		nodes, err := c.searchFollowingEventChainNodesByID(kit, node.ID, eventStep, key)
		if err != nil {
			blog.ErrorJSON("get event failed, err: %s, rid: %s, filter: %s", err, rid, filter)
			return nil, err
		}

		// since the first node is after the start time, we need to include it in the nodes after the start time
		if isFirstLoop {
			nodes = append([]*watch.ChainNode{node}, nodes...)
			isFirstLoop = false
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

		hitNodes := c.GetHitNodeWithEventType(nodes, opts.EventTypes)

		if len(hitNodes) != 0 {
			// matched event has been found, get them all.
			return c.getEventDetailsWithNodes(kit, opts, hitNodes, key)
		}

		// update filter and do next scan round.
		node = nodes[len(nodes)-1]
	}
}

// getEventDetailsWithNodes get event details with nodes, first get from redis, then get failed ones from mongo
func (c *Client) getEventDetailsWithNodes(kit *rest.Kit, opts *watch.WatchEventOptions, hitNodes []*watch.ChainNode, key event.Key) (
	[]*watch.WatchEventDetail, error) {

	if len(hitNodes) == 0 {
		return make([]*watch.WatchEventDetail, 0), nil
	}

	cursors := make([]string, len(hitNodes))
	for index, node := range hitNodes {
		cursors[index] = node.Cursor
	}

	details, errCursors, errCursorIndexMap, err := c.searchEventDetailsFromRedis(kit, cursors, key)
	if err != nil {
		return nil, err
	}

	if len(errCursors) == 0 {
		resp := make([]*watch.WatchEventDetail, len(details))
		for idx, detail := range details {
			resp[idx] = &watch.WatchEventDetail{
				Cursor:    hitNodes[idx].Cursor,
				Resource:  opts.Resource,
				EventType: hitNodes[idx].EventType,
				Detail:    watch.JsonString(detail),
			}
		}
		return resp, nil
	}

	// get event chain nodes from db for cursors that failed when reading redis
	errCursorsExistMap := make(map[string]struct{})
	for _, errCursor := range errCursors {
		errCursorsExistMap[errCursor] = struct{}{}
	}

	errNodes := make([]*watch.ChainNode, 0)
	for _, node := range hitNodes {
		if _, exists := errCursorsExistMap[node.Cursor]; exists {
			errNodes = append(errNodes, node)
		}
	}

	indexDetailMap, err := c.searchEventDetailsFromMongo(kit, errNodes, opts.Fields, errCursorIndexMap, key)
	if err != nil {
		blog.Errorf("get details from mongo failed, err: %v, cursors: %+v, rid: %s", err, errCursors, kit.Rid)
		return nil, err
	}

	resp := make([]*watch.WatchEventDetail, len(details))
	for idx, detail := range details {
		errDetail, exists := indexDetailMap[idx]
		if exists {
			detail = errDetail
		} else {
			jsonStr := types.GetEventDetail(&detail)
			detail = *json.CutJsonDataWithFields(jsonStr, opts.Fields)
		}

		resp[idx] = &watch.WatchEventDetail{
			Cursor:    hitNodes[idx].Cursor,
			Resource:  opts.Resource,
			EventType: hitNodes[idx].EventType,
			Detail:    watch.JsonString(detail),
		}
	}
	return resp, nil
}

// WatchFromNow watches target resource events from noc.
func (c *Client) WatchFromNow(kit *rest.Kit, key event.Key, opts *watch.WatchEventOptions) (
	*watch.WatchEventDetail, error) {

	rid := kit.Rid

	node, exists, err := c.getLatestEvent(kit, key)
	if err != nil {
		blog.Errorf("watch from now, but get latest event failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	if !exists {
		// event chain list is empty, which means no event and not be initialized.
		return &watch.WatchEventDetail{
			Cursor:    watch.NoEventCursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}, nil
	}

	hit := c.GetHitNodeWithEventType([]*watch.ChainNode{node}, opts.EventTypes)
	if len(hit) == 0 {
		// not matched, set to no event cursor with empty detail
		return &watch.WatchEventDetail{
			Cursor:    watch.NoEventCursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}, nil
	}

	detail, exists, err := c.getEventDetail(kit, node, opts.Fields, key)
	if err != nil {
		blog.Errorf("watch from now, but get latest event detail failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	if !exists {
		return nil, kit.CCError.CCError(common.CCErrEventDetailNotExist)
	}

	// matched the event type.
	return &watch.WatchEventDetail{
		Cursor:    node.Cursor,
		Resource:  opts.Resource,
		EventType: node.EventType,
		Detail:    watch.JsonString(detail),
	}, nil
}

// watchWithCursor get events with the start cursor which is offered by user.
// it will hold the request for timeout seconds if no matched event is hit.
// if event has been hit in a round, then events will be returned immediately.
// if no events hit, then will loop the event every 200ms until timeout and return
// with a special cursor named "NoEventCursor", then we will help the user watch
// event from the head cursor.
func (c *Client) WatchWithCursor(kit *rest.Kit, key event.Key, opts *watch.WatchEventOptions) (
	[]*watch.WatchEventDetail, error) {

	rid := kit.Rid
	start := time.Now().Unix()

	exists, nodes, err := c.searchFollowingEventChainNodes(kit, opts.Cursor, eventStep, key)
	if err != nil {
		blog.Errorf("search nodes after cursor %s failed, err: %v, rid: %s", opts.Cursor, err, kit.Rid)
		return nil, err
	}

	if !exists {
		return nil, kit.CCError.CCError(common.CCErrEventChainNodeNotExist)
	}

	for {
		if len(nodes) == 0 {
			if time.Now().Unix()-start > timeoutWatchLoopSeconds {
				// has already looped for timeout seconds, and we still got no event.
				// return with NoEventCursor and empty detail
				return []*watch.WatchEventDetail{{
					Cursor:    watch.NoEventCursor,
					Resource:  opts.Resource,
					EventType: "",
					Detail:    nil,
				}}, nil
			}

			// we got not event one event, sleep a little, and then try to continue the loop watch
			time.Sleep(loopInternal)
			blog.V(5).Infof("watch key: %s with resource: %s, got nothing, try next round. rid: %s",
				key.Namespace(), opts.Resource, rid)
			continue
		}

		hitNodes := c.GetHitNodeWithEventType(nodes, opts.EventTypes)
		if len(hitNodes) != 0 {
			// matched event has been found, get them all.
			blog.V(5).Infof("watch key: %s with resource: %s, hit events, return immediately. rid: %s", key.Namespace(), opts.Resource, rid)
			return c.getEventDetailsWithNodes(kit, opts, hitNodes, key)
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

		nodes, err = c.searchFollowingEventChainNodesByID(kit, nodes[len(nodes)-1].ID, eventStep, key)
		if err != nil {
			blog.Errorf("watch event from cursor: %s failed, err: %v, rid: %s", opts.Cursor, err, rid)
			return nil, err
		}
		continue
	}
}

func (c *Client) GetHitNodeWithEventType(nodes []*watch.ChainNode, typs []watch.EventType) []*watch.ChainNode {
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
