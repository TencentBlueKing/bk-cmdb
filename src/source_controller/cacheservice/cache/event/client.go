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
	"github.com/tidwall/gjson"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/stream/types"

	rawRedis "github.com/go-redis/redis/v7"
)

type Client struct {
	// cache is cc redis client.
	cache redis.Client

	// watchDB is cc event watch database.
	watchDB dal.DB

	// db is cc main database.
	db dal.DB
}

func NewClient(watchDB dal.DB, db dal.DB, cache redis.Client) *Client {
	return &Client{watchDB: watchDB, db: db, cache: cache}
}

// SearchEventChainNode search event chain node by condition, then search for node detail if needed
func (c *Client) SearchEventChainNode(kit *rest.Kit, opts *metadata.SearchEventNodeOption) (
	*metadata.EventNodeWithDetail, error) {

	key, err := event.GetResourceKeyWithCursorType(opts.Resource)
	if err != nil {
		blog.Errorf("get resource key with cursor type %s failed, err: %v, rid: %s", opts.Resource, err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_resource")
	}

	filter := setClusterTimeFilter(opts.Filter, key)
	query := c.watchDB.Table(key.ChainCollection()).Find(filter)
	if len(opts.Sort) > 0 {
		query.Sort(opts.Sort)
	}

	node := new(watch.ChainNode)
	if err := query.One(kit.Ctx, node); err != nil {
		blog.ErrorJSON("get start chain node from mongo failed, err: %s, opts: %s", err, opts)
		if !c.watchDB.IsNotFoundError(err) {
			return nil, err
		}
		return nil, kit.CCError.CCError(common.CCErrEventChainNodeNotExist)
	}

	if !opts.WithDetail {
		return &metadata.EventNodeWithDetail{Node: node}, nil
	}

	detail, err := c.cache.Get(kit.Ctx, key.DetailKey(node.Cursor)).Result()
	if err != nil {
		blog.Errorf("get watch detail from redis failed, err: %v, cursor: %s", err, node.Cursor)
		if !redis.IsNilErr(err) {
			return nil, err
		}

		detailFilter := map[string]interface{}{
			"_id": node.Oid,
		}

		detailStr := types.JsonString("")
		if err := c.watchDB.Table(key.Collection()).Find(detailFilter).One(kit.Ctx, &detailStr); err != nil {
			if c.watchDB.IsNotFoundError(err) {
				return nil, kit.CCError.CCError(common.CCErrEventDetailNotExist)
			}
			blog.ErrorJSON("get detail from db failed, err: %s, filter: %s, rid: %s", err, detailFilter, kit.Rid)
			return nil, err
		}

		detail = string(detailStr)
	}

	return &metadata.EventNodeWithDetail{Node: node, Detail: detail}, nil
}

// SearchFollowingEventNodes search nodes after the node(including itself) specified by filter
func (c *Client) SearchFollowingEventChainNodes(kit *rest.Kit, opts *metadata.SearchFollowingEventNodesOption) (
	[]*watch.ChainNode, error) {

	if opts.Limit > common.BKMaxPageSize {
		return nil, kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "limit", common.BKMaxPageSize)
	}

	key, err := event.GetResourceKeyWithCursorType(opts.Resource)
	if err != nil {
		blog.Errorf("get resource key with cursor type %s failed, err: %v, rid: %s", opts.Resource, err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_resource")
	}

	filter := setClusterTimeFilter(opts.Filter, key)
	collection := key.ChainCollection()
	node := new(watch.ChainNode)
	if err := c.watchDB.Table(collection).Find(filter).Fields(common.BKFieldID).One(kit.Ctx, node); err != nil {
		if !c.watchDB.IsNotFoundError(err) {
			blog.ErrorJSON("get start node failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
			return nil, err
		}
		return nil, kit.CCError.Error(common.CCErrEventChainNodeNotExist)
	}

	nodeFilter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{common.BKDBGTE: node.ID},
	}

	query := c.watchDB.Table(collection).Find(nodeFilter).Limit(uint64(opts.Limit))
	if len(opts.Sort) > 0 {
		query.Sort(opts.Sort)
	}

	nodes := make([]*watch.ChainNode, 0)
	if err := query.All(kit.Ctx, &nodes); err != nil {
		blog.Errorf("get chain nodes from mongo failed, err: %v, start id: %d, rid: %s", err, node.ID, kit.Rid)
		return nil, err
	}

	return nodes, nil
}

// SearchEventDetails search event details by cursors
func (c *Client) SearchEventDetails(kit *rest.Kit, opts *metadata.SearchEventDetailsOption) ([]string, error) {
	if len(opts.Cursors) == 0 {
		return make([]string, 0), nil
	}

	key, err := event.GetResourceKeyWithCursorType(opts.Resource)
	if err != nil {
		blog.Errorf("get resource key with cursor type %s failed, err: %v, rid: %s", opts.Resource, err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_resource")
	}

	results := make([]*rawRedis.StringCmd, len(opts.Cursors))
	pipe := c.cache.Pipeline()
	for idx, cursor := range opts.Cursors {
		results[idx] = pipe.Get(key.DetailKey(cursor))
	}

	_, err = pipe.Exec()
	if err != nil {
		blog.Errorf("search event details by cursors(%+v) failed, err: %v, rid: %s", opts.Cursors, err, kit.Rid)

		// get details from db for cursors
		errCursorIndexMap := make(map[string]int)
		errCursors := make([]string, 0)
		details := make([]string, len(results))
		for idx, result := range results {
			if err := result.Err(); err != nil {
				cursor := opts.Cursors[idx]
				blog.Errorf("search event detail by cursor(%s) failed, err: %v, rid: %s", cursor, err, kit.Rid)
				errCursorIndexMap[cursor] = idx
				errCursors = append(errCursors, cursor)
				continue
			}
			details[idx] = result.Val()
		}

		chainFilter := map[string]interface{}{
			common.BKCursorField: map[string]interface{}{common.BKDBIN: errCursors},
		}
		nodes := make([]*watch.ChainNode, 0)
		if err := c.watchDB.Table(key.ChainCollection()).Find(chainFilter).Fields(common.BKCursorField,
			common.BKOIDField).All(kit.Ctx, &nodes); err != nil {
			blog.Errorf("get chain nodes failed, err: %v, cursor: %+v, rid: %s", err, errCursors, kit.Rid)
			return nil, err
		}

		oids := make([]string, len(nodes))
		oidCursorMap := make(map[string]string)
		for index, node := range nodes {
			oids[index] = node.Oid
			oidCursorMap[node.Oid] = node.Cursor
		}

		detailFilter := map[string]interface{}{
			"_id": map[string]interface{}{common.BKDBIN: oids},
		}

		detailArr := make([]types.JsonString, 0)
		if err := c.watchDB.Table(key.Collection()).Find(detailFilter).All(kit.Ctx, &detailArr); err != nil {
			blog.Errorf("get details from db failed, err: %s, oids: %+v, rid: %s", err, oids, kit.Rid)
			return nil, err
		}

		for _, detailStr := range detailArr {
			oid := gjson.Get(string(detailStr), "_id").String()
			cursor := oidCursorMap[oid]
			index := errCursorIndexMap[cursor]
			details[index] = string(detailStr)
		}
		return details, nil
	}

	details := make([]string, len(results))
	for idx, result := range results {
		details[idx] = result.Val()
	}
	return details, nil
}

// set cluster time filter condition because db nodes have a longer ttl than the valid event's life time
func setClusterTimeFilter(filter map[string]interface{}, key event.Key) map[string]interface{} {
	return map[string]interface{}{
		common.BKDBAND: []map[string]interface{}{
			filter,
			{
				common.BKClusterTimeField: map[string]interface{}{
					common.BKDBGTE: time.Now().Add(-time.Duration(key.TTLSeconds()) * time.Second),
				},
			},
		},
	}
}
