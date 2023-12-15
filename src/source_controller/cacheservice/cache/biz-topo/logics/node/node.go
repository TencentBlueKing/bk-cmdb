/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package node defines business topology node related common logics
package node

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"

	"github.com/tidwall/gjson"
)

const step = 100

// PagedGetNodes loop getting paged nodes from db
func PagedGetNodes(ctx context.Context, table string, nodeCond mapstr.MapStr, fields []string, parser NodeParser,
	rid string) ([]types.Node, error) {

	nodes := make([]types.Node, 0)

	cond := nodeCond.Clone()
	findOpt := dbtypes.NewFindOpts().SetWithObjectID(true)
	for {
		data := make([]mapstr.MapStr, 0)
		err := mongodb.Client().Table(table).Find(cond, findOpt).Fields(fields...).Sort("_id").Limit(step).
			All(ctx, &data)
		if err != nil {
			blog.Errorf("get node data failed, table: %s, cond: %+v, err: %v, rid: %s", table, cond, err, rid)
			return nil, err
		}

		if len(data) == 0 {
			break
		}

		parsed, err := parser(ctx, data, rid)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, parsed...)

		if len(data) < step {
			break
		}

		cond["_id"] = mapstr.MapStr{common.BKDBGT: data[len(data)-1]["_id"]}
	}

	return nodes, nil
}

// NodeParser parse db node data to topo nodes
type NodeParser func(ctx context.Context, data []mapstr.MapStr, rid string) ([]types.Node, error)

// CombineChildNodes combine parent nodes with child nodes info
func CombineChildNodes(nodes, childNodes []types.Node) []types.Node {
	parentChildNodeMap := make(map[int64][]types.Node)
	hasCnt := false
	parentCntMap := make(map[int64]int64)
	for _, node := range childNodes {
		parentChildNodeMap[node.ParentID] = append(parentChildNodeMap[node.ParentID], node)

		if node.Count != nil {
			hasCnt = true
			parentCntMap[node.ParentID] += *node.Count
		}
	}

	for i, node := range nodes {
		nodes[i].SubNodes = parentChildNodeMap[node.ID]
		if hasCnt {
			cnt := parentCntMap[node.ID]
			nodes[i].Count = &cnt
		}
	}

	return nodes
}

// GenBizNodeListKey generate biz to topo node ids list cache key
func GenBizNodeListKey(topoKey key.Key, bizID int64, kind string) string {
	return fmt.Sprintf("%s:%s:list", topoKey.BizTopoKey(bizID), kind)
}

// GenNodeInfoKey generate biz topo node info separate cache key
func GenNodeInfoKey(topoKey key.Key, bizID int64, kind string, id int64) string {
	return fmt.Sprintf("%s:%s:%d", topoKey.BizTopoKey(bizID), kind, id)
}

// AddNodeInfoCache add biz topo nodes info cache by kind
func AddNodeInfoCache(topoKey key.Key, bizID int64, kind string, nodes []types.Node, rid string) error {
	pip := redis.Client().Pipeline()
	defer pip.Close()

	listKey := GenBizNodeListKey(topoKey, bizID, kind)

	ids := make([]interface{}, len(nodes))
	for i, node := range nodes {
		nodeKey := GenNodeInfoKey(topoKey, bizID, kind, node.ID)
		pip.Set(nodeKey, fmt.Sprintf(`{"id":%d,"nm":"%s","par":%d}`, node.ID, node.Name, node.ParentID), topoKey.TTL())
		ids[i] = node.ID
	}
	pip.SAdd(listKey, ids...)
	pip.Expire(listKey, topoKey.TTL())

	_, err := pip.Exec()
	if err != nil {
		blog.Errorf("cache biz %d topo nodes info failed, err: %v, nodes: %+v, rid: %s", bizID, err, nodes, rid)
		return err
	}
	return nil
}

// DeleteNodeInfoCache delete biz topo nodes info cache by kind
func DeleteNodeInfoCache(topoKey key.Key, bizID int64, kind string, ids []int64, rid string) error {
	pip := redis.Client().Pipeline()
	defer pip.Close()

	listKey := GenBizNodeListKey(topoKey, bizID, kind)
	pip.Expire(listKey, topoKey.TTL())

	idList := make([]interface{}, len(ids))
	for i, id := range ids {
		nodeKey := GenNodeInfoKey(topoKey, bizID, kind, id)
		pip.Del(nodeKey)
		idList[i] = id
	}
	pip.SRem(listKey, idList...)

	_, err := pip.Exec()
	if err != nil {
		blog.Errorf("delete biz %d topo nodes info cache failed, err: %v, ids: %+v, rid: %s", bizID, err, ids, rid)
		return err
	}
	return nil
}

// GetNodeInfoCache get biz topo nodes info cache by kind
func GetNodeInfoCache(topoKey key.Key, bizID int64, kind string, rid string) ([]types.Node, error) {
	ctx := context.Background()

	listKey := GenBizNodeListKey(topoKey, bizID, kind)

	cursor := uint64(0)

	nodes := make([]types.Node, 0)
	for {
		ids, nextCursor, err := redis.Client().SScan(listKey, cursor, "", step).Result()
		if err != nil {
			blog.Errorf("scan topo node cache list %s %d failed, err: %v, rid: %s", listKey, cursor, err, rid)
			return nil, err
		}
		cursor = nextCursor

		if len(ids) == 0 {
			if nextCursor == uint64(0) {
				return nodes, nil
			}
			continue
		}

		detailKeys := make([]string, len(ids))
		for i, idStr := range ids {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				blog.Errorf("parse node id %s failed, err: %v, rid: %s", idStr, err, rid)
				continue
			}
			detailKeys[i] = GenNodeInfoKey(topoKey, bizID, kind, id)
		}

		details, err := redis.Client().MGet(ctx, detailKeys...).Result()
		if err != nil {
			blog.Errorf("get topo node cache details by keys: %+v failed, err: %v, rid: %s", detailKeys, err, rid)
			return nil, err
		}

		for _, detail := range details {
			if detail == nil {
				continue
			}

			strVal, ok := detail.(string)
			if !ok {
				blog.Errorf("node info cache detail type %T is invalid, detail: %v, rid: %s", detail, detail, rid)
				continue
			}

			nodes = append(nodes, types.Node{
				Kind:     kind,
				ID:       gjson.Get(strVal, "id").Int(),
				Name:     gjson.Get(strVal, "nm").String(),
				ParentID: gjson.Get(strVal, "par").Int(),
			})
		}

		if nextCursor == uint64(0) {
			return nodes, nil
		}
	}
}

// CrossCompareNodeInfoCache cross compare biz topo nodes info cache by kind
func CrossCompareNodeInfoCache(topoKey key.Key, bizID int64, kind string, nodes []types.Node, rid string) error {
	nodeMap := make(map[int64]struct{}, len(nodes))

	// paged add biz topo node info cache
	pagedNodes := make([]types.Node, 0)
	for _, node := range nodes {
		nodeMap[node.ID] = struct{}{}

		pagedNodes = append(pagedNodes, node)
		if len(pagedNodes) == step {
			if err := AddNodeInfoCache(topoKey, bizID, kind, pagedNodes, rid); err != nil {
				return err
			}
			pagedNodes = make([]types.Node, 0)
		}
	}

	if len(pagedNodes) > 0 {
		if err := AddNodeInfoCache(topoKey, bizID, kind, pagedNodes, rid); err != nil {
			return err
		}
	}

	listKey := GenBizNodeListKey(topoKey, bizID, kind)
	cursor := uint64(0)

	// paged delete redundant biz topo node info cache
	for {
		ids, nextCursor, err := redis.Client().SScan(listKey, cursor, "", step).Result()
		if err != nil {
			blog.Errorf("scan topo node cache list %s %d failed, err: %v, rid: %s", listKey, cursor, err, rid)
			return err
		}
		cursor = nextCursor

		if len(ids) == 0 {
			if nextCursor == uint64(0) {
				return nil
			}
			continue
		}

		delIDs := make([]int64, 0)
		for _, idStr := range ids {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				blog.Errorf("parse node id %s failed, err: %v, rid: %s", idStr, err, rid)
				continue
			}

			_, exists := nodeMap[id]
			if !exists {
				delIDs = append(delIDs, id)
			}
		}

		if len(delIDs) > 0 {
			if err = DeleteNodeInfoCache(topoKey, bizID, kind, delIDs, rid); err != nil {
				return err
			}
		}

		if nextCursor == uint64(0) {
			return nil
		}
	}
}

// GenNodeInfoCntKey generate biz topo node info count cache key
func GenNodeInfoCntKey(topoKey key.Key, bizID int64, kind string, id int64) string {
	return GenNodeInfoKey(topoKey, bizID, kind, id) + ":count"
}

// AddNodeCountCache add biz topo nodes count cache by kind
func AddNodeCountCache(topoKey key.Key, bizID int64, kind string, cntMap map[int64]int64, rid string) error {
	pip := redis.Client().Pipeline()
	defer pip.Close()

	for id, cnt := range cntMap {
		cntKey := GenNodeInfoCntKey(topoKey, bizID, kind, id)
		pip.Set(cntKey, cnt, topoKey.TTL())
	}

	_, err := pip.Exec()
	if err != nil {
		blog.Errorf("cache biz %d topo node count info %+v failed, err: %v, rid: %s", bizID, cntMap, err, rid)
		return err
	}
	return nil
}

// DeleteNodeCountCache delete biz topo node count cache by kind
func DeleteNodeCountCache(topoKey key.Key, bizID int64, kind string, ids []int64, rid string) error {
	pip := redis.Client().Pipeline()
	defer pip.Close()

	for _, id := range ids {
		nodeKey := GenNodeInfoCntKey(topoKey, bizID, kind, id)
		pip.Del(nodeKey)
	}

	_, err := pip.Exec()
	if err != nil {
		blog.Errorf("delete biz %d topo nodes count cache failed, err: %v, ids: %+v, rid: %s", bizID, err, ids, rid)
		return err
	}
	return nil
}
