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

// Package watch TODO
package watch

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/source_controller/coreservice/core/host/identifier"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"
)

/* eventserver watcher defines, just created base on old service/watch.go */

const (
	// timeoutWatchLoopSeconds 20s timeout
	timeoutWatchLoopSeconds = 20

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
		if !c.isNodeHitEventType(tailNode, opts.EventTypes) || !c.isNodeHitSubResource(tailNode, opts.Filter.SubResource) {
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

		event := &watch.WatchEventDetail{
			Cursor:    tailNode.Cursor,
			Resource:  opts.Resource,
			EventType: tailNode.EventType,
		}

		if detail == nil {
			// convert to a no event cursor
			event.Detail = nil
		} else {
			if len(*detail) == 0 {
				// convert to a no event cursor
				event.Detail = nil
			} else {
				event.Detail = watch.JsonString(*detail)
			}

		}

		// matched the event type.
		return []*watch.WatchEventDetail{event}, nil
	}

	// find the first node with a larger cluster time than the start from parameter
	filter := map[string]interface{}{
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGT: metadata.Time{Time: time.Unix(opts.StartFrom, 0).Local()},
		},
	}

	// filters out the previous version where sub resource is string type // TODO remove this
	if key.Collection() == common.BKTableNameBaseInst ||
		key.Collection() == common.BKTableNameMainlineInstance {
		filter[common.BKSubResourceField] = map[string]interface{}{common.BKDBType: "array"}
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

	searchOpt := &searchFollowingChainNodesOption{
		id:          node.ID,
		limit:       eventStep,
		types:       opts.EventTypes,
		key:         key,
		subResource: opts.Filter.SubResource,
	}
	nodes, err := c.searchFollowingEventChainNodesByID(kit, searchOpt)
	if err != nil {
		blog.ErrorJSON("get event failed, err: %s, rid: %s, filter: %s", err, rid, filter)
		return nil, err
	}

	// since the first node is after the start time, we need to include it in the nodes after the start time
	if c.isNodeHitEventType(node, opts.EventTypes) && c.isNodeHitSubResource(node, opts.Filter.SubResource) {
		nodes = append([]*watch.ChainNode{node}, nodes...)
	}

	// nodes has already scan to the end
	if len(nodes) == 0 {
		resp := &watch.WatchEventDetail{
			Cursor:    tailNode.Cursor,
			Resource:  opts.Resource,
			EventType: "",
			Detail:    nil,
		}
		return []*watch.WatchEventDetail{resp}, nil
	}

	// matched event has been found, get them all.
	return c.getEventDetailsWithNodes(kit, opts, nodes, key)
}

// getEventDetailsWithNodes get event details with nodes, first get from redis, then get failed ones from mongo
func (c *Client) getEventDetailsWithNodes(kit *rest.Kit, opts *watch.WatchEventOptions, hitNodes []*watch.ChainNode, key event.Key) (
	[]*watch.WatchEventDetail, error) {

	if len(hitNodes) == 0 {
		return make([]*watch.WatchEventDetail, 0), nil
	}

	if opts.Resource == watch.HostIdentifier {
		// get from db directly.
		return c.getHostIdentityEventDetailWithNodes(kit, hitNodes)
	}

	if opts.Resource == watch.BizSetRelation {
		// get from db directly.
		return c.getBizSetRelationEventDetailWithNodes(kit, hitNodes)
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
			jsonStr := types.GetEventDetail(&detail)
			detail = *json.CutJsonDataWithFields(jsonStr, opts.Fields)
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

// getHostIdentityEventDetailWithNodes TODO
// get host identity from db directly
func (c *Client) getHostIdentityEventDetailWithNodes(kit *rest.Kit, hitNodes []*watch.ChainNode) (
	[]*watch.WatchEventDetail, error) {

	if len(hitNodes) == 0 {
		return nil, errors.New("no hit host identity event nodes")
	}

	hostIDs := make([]int64, 0)
	for idx := range hitNodes {
		if hitNodes[idx].InstanceID <= 0 {
			monitor.Collect(&meta.Alarm{
				RequestID: kit.Rid,
				Type:      meta.EventFatalError,
				Detail: fmt.Sprintf("host identity, instance id: %d is invalid, cursor: %s",
					hitNodes[idx].InstanceID, hitNodes[idx].Cursor),
				Module:    types2.CC_MODULE_CACHESERVICE,
				Dimension: map[string]string{"host_identifier": "yes"},
			})

			blog.ErrorJSON("get host identity with chain nodes, but got invalid host id, skip, detail: %s, rid: %s",
				hitNodes[idx], kit.Rid)
			continue
		}

		hostIDs = append(hostIDs, hitNodes[idx].InstanceID)
	}

	hostIDs = util.IntArrayUnique(hostIDs)
	// read from secondary, but this may get host identity may not same with master.
	// kit.Ctx, kit.Header = util.SetReadPreference(kit.Ctx, kit.Header, common.SecondaryPreferredMode)
	list, err := identifier.NewIdentifier().Identifier(kit, hostIDs)
	if err != nil {
		blog.Errorf("get host identity from db failed, host id: %v, err: %v, rid: %s", hostIDs, err, kit.Rid)
		return nil, err
	}

	identityMap := make(map[int64]*metadata.HostIdentifier)
	for idx := range list {
		identityMap[list[idx].HostID] = &list[idx]
	}

	details := make([]*watch.WatchEventDetail, 0)
	for _, one := range hitNodes {

		if one.InstanceID <= 0 {
			// skip
			continue
		}

		identity, exists := identityMap[one.InstanceID]
		if !exists {
			// host already be deleted, skip this event.
			continue
		}

		js, err := json.Marshal(identity)
		if err != nil {
			blog.Errorf("marshal host identity failed, skip, detail: %+v, err :%v, rid: %s", *identity, err, kit.Rid)
			continue
		}

		details = append(details, &watch.WatchEventDetail{
			Cursor:    one.Cursor,
			Resource:  watch.HostIdentifier,
			EventType: one.EventType,
			Detail:    watch.JsonString(js),
		})
	}

	if len(details) == 0 {
		// return the last node's cursor without details, so that use can watch from
		// this nodes continually.
		lastOne := hitNodes[len(hitNodes)-1]
		one := &watch.WatchEventDetail{
			Cursor:    lastOne.Cursor,
			Resource:  watch.HostIdentifier,
			EventType: lastOne.EventType,
			Detail:    nil,
		}
		return []*watch.WatchEventDetail{one}, nil
	}

	return details, nil
}

// getBizSetRelationEventDetailWithNodes get biz set relation event detail by chain nodes
func (c *Client) getBizSetRelationEventDetailWithNodes(kit *rest.Kit, hitNodes []*watch.ChainNode) (
	[]*watch.WatchEventDetail, error) {

	if len(hitNodes) == 0 {
		return make([]*watch.WatchEventDetail, 0), nil
	}

	cursors := make([]string, len(hitNodes))
	for index, node := range hitNodes {
		cursors[index] = node.Cursor
	}

	// get biz set relation event detail from redis, ignore the fields in watch option, return the whole detail
	details, errCursors, errCursorIndexMap, err := c.searchEventDetailsFromRedis(kit, cursors, event.BizSetRelationKey)
	if err != nil {
		return nil, err
	}

	if len(errCursors) == 0 {
		resp := make([]*watch.WatchEventDetail, len(details))
		for idx, detail := range details {
			resp[idx] = &watch.WatchEventDetail{
				Cursor:    hitNodes[idx].Cursor,
				Resource:  watch.BizSetRelation,
				EventType: hitNodes[idx].EventType,
				Detail:    watch.JsonString(*types.GetEventDetail(&detail)),
			}
		}
		return resp, nil
	}

	// get event chain nodes related details from db for cursors that failed when reading redis
	errCursorsExistMap := make(map[string]struct{})
	for _, errCursor := range errCursors {
		errCursorsExistMap[errCursor] = struct{}{}
	}

	bizSetIDs := make([]int64, 0)
	indexBizSetIDMap := make(map[int]int64)
	for _, node := range hitNodes {
		if _, exists := errCursorsExistMap[node.Cursor]; exists {
			bizSetIDs = append(bizSetIDs, node.InstanceID)
			indexBizSetIDMap[errCursorIndexMap[node.Cursor]] = node.InstanceID
		}
	}

	bizSetDetailMap, err := c.getBizSetRelationEventDetailFromMongo(kit, bizSetIDs)
	if err != nil {
		return nil, err
	}

	// generate event details, if mongo detail not exists(biz set is deleted), return biz set id with empty relations
	resp := make([]*watch.WatchEventDetail, len(details))
	for idx, detail := range details {
		bizSetID, exists := indexBizSetIDMap[idx]
		if exists {
			detail, exists = bizSetDetailMap[bizSetID]
			if !exists {
				detail = event.GenBizSetRelationDetail(bizSetID, "")
			}
		} else {
			detail = *types.GetEventDetail(&detail)
		}

		resp[idx] = &watch.WatchEventDetail{
			Cursor:    hitNodes[idx].Cursor,
			Resource:  watch.BizSetRelation,
			EventType: hitNodes[idx].EventType,
			Detail:    watch.JsonString(detail),
		}
	}
	return resp, nil
}

// getBizSetRelationEventDetailFromMongo get biz set relation event details by biz set ids from mongo
func (c *Client) getBizSetRelationEventDetailFromMongo(kit *rest.Kit, bizSetIDs []int64) (map[int64]string, error) {
	if len(bizSetIDs) == 0 {
		return make(map[int64]string), nil
	}

	// get biz sets by chain nodes instance ids from db
	bizSetCond := map[string]interface{}{
		common.BKBizSetIDField: map[string]interface{}{common.BKDBIN: bizSetIDs},
	}

	bizSets := make([]metadata.BizSetInst, 0)
	err := mongodb.Client().Table(common.BKTableNameBaseBizSet).Find(bizSetCond).
		Fields(common.BKBizSetIDField, common.BKBizSetScopeField).All(kit.Ctx, &bizSets)
	if err != nil {
		blog.Errorf("get biz sets by cond(%+v) failed, err: %v, rid: %s", bizSetCond, err, kit.Rid)
		return nil, err
	}

	// get biz set id to detail map by searching for biz ids by biz set scope
	var allBizIDStr string
	bizSetDetailMap := make(map[int64]string)
	for _, bizSet := range bizSets {
		// save a cache of all biz ids string form for biz set scope that matches all, use it to gen relation detail
		if bizSet.Scope.MatchAll {
			if len(allBizIDStr) == 0 {
				// do not include resource pool and disabled biz in biz set
				allBizIDCond := map[string]interface{}{
					common.BKDefaultField:    mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag},
					common.BKDataStatusField: map[string]interface{}{common.BKDBNE: common.DataStatusDisabled},
				}

				allBizIDStr, err = c.getBizIDArrStrByCond(kit, allBizIDCond)
				if err != nil {
					return nil, err
				}
			}
			bizSetDetailMap[bizSet.BizSetID] = event.GenBizSetRelationDetail(bizSet.BizSetID, allBizIDStr)
			continue
		}

		// parse biz condition from biz set scope filter, get biz ids using it to gen relation detail
		if bizSet.Scope.Filter == nil {
			continue
		}

		bizSetBizCond, errKey, rawErr := bizSet.Scope.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse biz set scope(%#v) failed, err: %v, rid: %s", bizSet.Scope, rawErr, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey)
		}

		// do not include resource pool and disabled biz in biz set
		bizSetBizCond[common.BKDefaultField] = mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag}
		bizSetBizCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}

		bizIDStr, err := c.getBizIDArrStrByCond(kit, bizSetBizCond)
		if err != nil {
			return nil, err
		}
		bizSetDetailMap[bizSet.BizSetID] = event.GenBizSetRelationDetail(bizSet.BizSetID, bizIDStr)
	}

	return bizSetDetailMap, nil
}

func (c *Client) getBizIDArrStrByCond(kit *rest.Kit, cond map[string]interface{}) (string, error) {
	const step = 500

	bizIDJson := bytes.Buffer{}

	for start := uint64(0); ; start += step {
		oneStep := make([]metadata.BizInst, 0)

		err := c.db.Table(common.BKTableNameBaseApp).Find(cond).Fields(common.BKAppIDField).Start(start).
			Limit(step).Sort(common.BKAppIDField).All(kit.Ctx, &oneStep)
		if err != nil {
			blog.Errorf("get biz by cond(%+v) failed, err: %v, rid: %s", cond, err, kit.Rid)
			return "", err
		}

		for _, biz := range oneStep {
			bizIDJson.WriteString(strconv.FormatInt(biz.BizID, 10))
			bizIDJson.WriteByte(',')
		}

		if len(oneStep) < step {
			break
		}
	}

	// returns biz ids string form joined by comma, trim the extra trilling comma
	if bizIDJson.Len() == 0 {
		return "", nil
	}
	return bizIDJson.String()[:bizIDJson.Len()-1], nil
}

// WatchFromNow watches target resource events from now.
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

	if !c.isNodeHitEventType(node, opts.EventTypes) || !c.isNodeHitSubResource(node, opts.Filter.SubResource) {
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

	e := &watch.WatchEventDetail{
		Cursor:    node.Cursor,
		Resource:  opts.Resource,
		EventType: node.EventType,
	}

	if detail == nil {
		// convert to a no event cursor
		e.Detail = nil
	} else {
		if len(*detail) == 0 {
			// convert to a no event cursor
			e.Detail = nil
		} else {
			e.Detail = watch.JsonString(*detail)
		}

	}

	// matched the event type.
	return e, nil
}

// WatchWithCursor get events with the start cursor which is offered by user.
// it will hold the request for timeout seconds if no matched event is hit.
// if event has been hit in a round, then events will be returned immediately.
// if no events hit, then will loop the event every 200ms until timeout and return
// with a special cursor named "NoEventCursor", then we will help the user watch
// event from the head cursor.
func (c *Client) WatchWithCursor(kit *rest.Kit, key event.Key, opts *watch.WatchEventOptions) (
	[]*watch.WatchEventDetail, error) {

	rid := kit.Rid
	start := time.Now().Unix()

	searchOpt := &searchFollowingChainNodesOption{
		startCursor: opts.Cursor,
		limit:       eventStep,
		types:       opts.EventTypes,
		key:         key,
		subResource: opts.Filter.SubResource,
	}
	exists, nodes, nodeID, err := c.searchFollowingEventChainNodes(kit, searchOpt)
	if err != nil {
		blog.Errorf("search nodes after cursor %s failed, err: %v, rid: %s", opts.Cursor, err, kit.Rid)
		return nil, err
	}

	if !exists && opts.Cursor != watch.NoEventCursor {
		return nil, kit.CCError.CCError(common.CCErrEventChainNodeNotExist)
	}

	for {
		if len(nodes) != 0 {
			return c.getEventDetailsWithNodes(kit, opts, nodes, key)
		}

		// we got not even one event, sleep a little, and then try to continue the loop watch
		time.Sleep(loopInternal)
		blog.V(5).Infof("watch key: %s with resource: %s, got nothing, try next round. rid: %s",
			key.Namespace(), opts.Resource, rid)

		if time.Now().Unix()-start > timeoutWatchLoopSeconds {
			lastNode, exists, err := c.getLatestEvent(kit, key)
			if err != nil {
				blog.Errorf("watch from now, but get latest event failed, err: %v, rid: %s", err, rid)
				return nil, err
			}

			if !exists {
				// has already looped for timeout seconds, and we still got no event.
				// return with NoEventCursor and empty detail
				opts.Cursor = watch.NoEventCursor
				return []*watch.WatchEventDetail{{
					Cursor:    watch.NoEventCursor,
					Resource:  opts.Resource,
					EventType: "",
					Detail:    nil,
				}}, nil
			} else {
				// 如果最后一个事件存在，则重新拉取匹配watch条件(type和sub resource)的事件，防止最后一个事件正好在超时之后但是
				// 拉取之前产生的情况下丢失从超时起到最后一个事件之间的事件。如果从起始cursor到最后一个事件之间没有匹配事件的话，
				// 返回最后一个事件，以免下次拉取时需要从起始cursor再重新拉取一遍不匹配的事件
				searchOpt.id = nodeID
				nodes, err = c.searchFollowingEventChainNodesByID(kit, searchOpt)
				if err != nil {
					blog.Errorf("watch event from cursor: %s failed, err: %v, rid: %s", opts.Cursor, err, rid)
					return nil, err
				}
				if len(nodes) != 0 {
					return c.getEventDetailsWithNodes(kit, opts, nodes, key)
				}

				resp := &watch.WatchEventDetail{
					Cursor:   lastNode.Cursor,
					Resource: opts.Resource,
					Detail:   nil,
				}

				// at least the tail node should can be scan, so something goes wrong.
				return []*watch.WatchEventDetail{resp}, nil
			}
		}

		searchOpt.id = nodeID
		nodes, err = c.searchFollowingEventChainNodesByID(kit, searchOpt)
		if err != nil {
			blog.Errorf("watch event from cursor: %s failed, err: %v, rid: %s", opts.Cursor, err, rid)
			return nil, err
		}
	}
}

func (c *Client) isNodeHitEventType(node *watch.ChainNode, types []watch.EventType) bool {
	if len(types) == 0 {
		return true
	}

	for _, eventType := range types {
		if node.EventType == eventType {
			return true
		}
	}
	return false
}

// isNodeHitSubResource check if node hit the sub resource, not specifying sub resource means matching all and is hit
// only used for common and mainline instances that contains sub resource
func (c *Client) isNodeHitSubResource(node *watch.ChainNode, subResource string) bool {
	if len(subResource) == 0 {
		return true
	}

	for _, subRes := range node.SubResource {
		if subRes == subResource {
			return true
		}
	}

	return false
}
