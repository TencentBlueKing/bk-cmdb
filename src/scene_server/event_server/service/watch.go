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

package service

import (
	ejson "encoding/json"
	"errors"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/coreservice/event"
	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"
)

func (s *Service) WatchEvent(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	resource := req.PathParameter("resource")
	options := new(watch.WatchEventOptions)
	if err := ejson.NewDecoder(req.Request.Body).Decode(options); err != nil {
		blog.Errorf("watch event, but decode request body failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		resp.WriteError(http.StatusOK, result)
		return
	}
	options.Resource = watch.CursorType(resource)

	if err := options.Validate(); err != nil {
		blog.Errorf("watch event, but got invalid request options, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	key, err := event.GetResourceKeyWithCursorType(options.Resource)
	if err != nil {
		blog.Errorf("watch event, but get resource key with cursor type failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	// watch with cursor
	if len(options.Cursor) != 0 {
		events, err := s.watchWithCursor(key, options, rid)
		if err != nil {
			blog.Errorf("watch event with cursor failed, cursor: %s, err: %v, rid: %s", options.Cursor, err, rid)
			resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
			return
		}

		// if not events is hit, then we return user's cursor, so that they can watch with this cursor again.
		resp.WriteEntity(s.generateResp(options.Cursor, options.Resource, events))
		return
	}

	// watch with start from
	if options.StartFrom != 0 {
		events, err := s.watchWithStartFrom(key, options, rid)
		if err != nil {
			blog.Errorf("watch event with start from: %s, err: %v, rid: %s", time.Unix(options.StartFrom, 0).Format(time.RFC3339), err, rid)
			resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
			return
		}

		resp.WriteEntity(s.generateResp("", options.Resource, events))
		return
	}

	// watch from now
	events, err := s.watchFromNow(key, options, rid)
	if err != nil {
		blog.Errorf("watch event from now, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	resp.WriteEntity(s.generateResp("", options.Resource, []*watch.WatchEventDetail{events}))
}

func (s *Service) generateResp(startCursor string, rsc watch.CursorType, events []*watch.WatchEventDetail) *metadata.
	Response {
	result := new(watch.WatchResp)
	if len(events) == 0 {
		result.Watched = false
		if len(startCursor) == 0 {
			result.Events = []*watch.WatchEventDetail{
				{
					Cursor:   watch.NoEventCursor,
					Resource: rsc,
				},
			}
		} else {
			// if user's watch with a start cursor, but we do not find event after this cursor,
			// then we return this start cursor directly, so that they can watch with this cursor for next round.
			result.Events = []*watch.WatchEventDetail{
				{
					Cursor:   startCursor,
					Resource: rsc,
				},
			}
		}

	} else {
		if events[0].Cursor == watch.NoEventCursor {
			result.Watched = false

			if len(startCursor) == 0 {
				// user watch with start form time, or watch from now, then return with NoEventCursor cursor.
				result.Events = []*watch.WatchEventDetail{
					{
						Cursor:   watch.NoEventCursor,
						Resource: rsc,
					},
				}
			} else {
				// if user's watch with a start cursor, but hit a NoEventCursor cursor,
				// then we return this start cursor directly, so that they can watch with this cursor for next round.
				result.Events = []*watch.WatchEventDetail{
					{
						Cursor:   startCursor,
						Resource: rsc,
					},
				}
			}

		} else {
			result.Watched = true
			result.Events = events
		}
	}

	return metadata.NewSuccessResp(result)
}

func (s *Service) watchWithStartFrom(key event.Key, opts *watch.WatchEventOptions, rid string) ([]*watch.WatchEventDetail, error) {

	// validate start from value is in the range or not
	headTarget, tailTarget, err := s.getHeadTailNodeTargetNode(key)
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

		latestEvent, err := s.watchFromNow(key, opts, rid)
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
	timeout := time.After(25 * time.Second)
	for {
		select {
		case <-timeout:
			// scan the event's too long time, need to exist immediately.
			blog.Errorf("watch with start from: %d, scan the cursor chain, but scan too long time, rid: %s", opts.StartFrom, rid)
			return nil, errors.New("scan the event cost too long time")
		default:

		}

		// scan event node from head, returned nodes does not contain tail node.
		nodes, err := s.getNodesFromCursor(eventStep, nextCursor, key)
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

		hitNodes := getHitNodeWithEventType(nodes, opts.EventTypes)
		matchedNodes := make([]*watch.ChainNode, 0)
		for _, node := range hitNodes {
			// find node that cluster time is larger than the start from seconds.
			if int64(node.ClusterTime.Sec) >= opts.StartFrom {
				matchedNodes = append(matchedNodes, node)
			}
		}

		if len(matchedNodes) != 0 {
			// matched event has been found, get them all.
			return s.getEventsWithCursorNodes(opts, matchedNodes, key, rid)
		}

		// not even one is hit.
		// check if nodes has already scan to the end
		lastNode := nodes[len(nodes)-1]
		if lastNode.NextCursor == key.TailKey() {
			// has already scan to the end, no need to scan anymore
			// get event detail.
			detail, err := s.cache.Get(key.DetailKey(lastNode.Cursor)).Result()
			if err != nil {
				blog.Errorf("get cursor: %s detail failed, err: %v, rid: %s", lastNode.Cursor, err, rid)
				return nil, err
			}

			resp := &watch.WatchEventDetail{
				Cursor:    lastNode.Cursor,
				Resource:  opts.Resource,
				EventType: lastNode.EventType,
				Detail:    watch.JsonString(detail),
			}
			return []*watch.WatchEventDetail{resp}, nil
		}

		// update nextCursor and do next scan round.
		nextCursor = lastNode.Cursor
	}
}

func (s *Service) getEventsWithCursorNodes(opts *watch.WatchEventOptions, hitNodes []*watch.ChainNode, key event.Key, rid string) ([]*watch.WatchEventDetail, error) {
	results := make([]*redis.StringCmd, 0)
	pipe := s.cache.Pipeline()
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
		blog.ErrorJSON("watch with start from: %d, resource: %s, hit events, but get event detail failed, hit nodes: %s, err: %v, rid: %s",
			opts.StartFrom, opts.Resource, hitNodes, err, rid)
		return nil, err
	}
	resp := make([]*watch.WatchEventDetail, 0)
	for idx, result := range results {
		jsonStr := result.Val()
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

func (s *Service) watchFromNow(key event.Key, opts *watch.WatchEventOptions, rid string) (*watch.WatchEventDetail, error) {
	node, tailTarget, err := s.getLatestEventDetail(key)
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

	hit := getHitNodeWithEventType([]*watch.ChainNode{node}, opts.EventTypes)
	if len(hit) == 0 {
		// not matched, set to no event cursor with empty detail
		return &watch.WatchEventDetail{
			Cursor:    watch.NoEventCursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}, nil
	}
	cut := json.CutJsonDataWithFields(&tailTarget, opts.Fields)
	// matched the event type.
	return &watch.WatchEventDetail{
		Cursor:    node.Cursor,
		Resource:  opts.Resource,
		EventType: node.EventType,
		Detail:    watch.JsonString(*cut),
	}, nil
}

const (
	// 25s timeout
	timeoutWatchLoopSeconds = 25
	// watch loop internal duration
	loopInternal = 250 * time.Millisecond
)

// watchWithCursor get events with the start cursor which is offered by user.
// it will hold the request for timeout seconds if no matched event is hit.
// if event has been hit in a round, then events will be returned immediately.
// if no events hit, then will loop the event every 200ms until timeout and return
// with a special cursor named "NoEventCursor", then we will help the user watch
// event from the head cursor.
func (s *Service) watchWithCursor(key event.Key, opts *watch.WatchEventOptions, rid string) ([]*watch.WatchEventDetail, error) {
	startCursor := opts.Cursor
	if startCursor == watch.NoEventCursor {
		// user got no events because of no event occurs in the system in the previous watch around,
		// we should watch from the head cursor in this round, so that user can not miss any events.
		startCursor = key.HeadKey()
	}

	start := time.Now().Unix()
	for {
		nodes, err := s.getNodesFromCursor(eventStep, startCursor, key)
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

		hitNodes := getHitNodeWithEventType(nodes, opts.EventTypes)
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
			return s.getEventsWithCursorNodes(opts, hitNodes, key, rid)
		}

		if time.Now().Unix()-start > timeoutWatchLoopSeconds {
			// no event is hit, but timeout, we return the last event cursor with nil detail
			// because it's not what the use want, return the last cursor to help user can
			// watch from here later for next watch round.
			lastNode := nodes[len(nodes)-1]
			resp := &watch.WatchEventDetail{
				Cursor:    lastNode.Cursor,
				Resource:  opts.Resource,
				EventType: lastNode.EventType,
				Detail:    nil,
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

func getHitNodeWithEventType(nodes []*watch.ChainNode, typs []watch.EventType) []*watch.ChainNode {
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
