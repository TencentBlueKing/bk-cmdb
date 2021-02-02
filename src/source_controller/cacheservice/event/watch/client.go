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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	daltypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
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

// GetLatestEvent get latest event chain node for resource
func (c *Client) GetLatestEvent(kit *rest.Kit, opts *metadata.GetLatestEventOption) (
	*metadata.EventNode, error) {

	key, err := event.GetResourceKeyWithCursorType(opts.Resource)
	if err != nil {
		blog.Errorf("get resource key with cursor type %s failed, err: %v, rid: %s", opts.Resource, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_resource")
	}

	node, exists, err := c.getLatestEvent(kit, key)
	if err != nil {
		blog.Errorf("get latest event for resource %s failed, err: %v", opts.Resource, err)
		return nil, err
	}

	return &metadata.EventNode{Node: node, ExistsNode: exists}, nil
}

// getLatestEvent search latest event chain node in not expired nodes
func (c *Client) getLatestEvent(kit *rest.Kit, key event.Key) (*watch.ChainNode, bool, error) {
	filter := map[string]interface{}{
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGTE: metadata.Time{Time: time.Now().Add(-time.Duration(key.TTLSeconds()) * time.Second).UTC()},
		},
	}

	node := new(watch.ChainNode)
	err := c.watchDB.Table(key.ChainCollection()).Find(filter).Sort(common.BKFieldID+":-1").One(kit.Ctx, node)
	if err != nil {
		blog.ErrorJSON("get chain node from mongo failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
		if !c.watchDB.IsNotFoundError(err) {
			return nil, false, fmt.Errorf("get last chain node from mongo failed, err: %v", err)
		}
		return nil, false, nil
	}
	return node, true, nil
}

// getEarliestEvent search earliest event chain node in not expired nodes
func (c *Client) getEarliestEvent(kit *rest.Kit, key event.Key) (*watch.ChainNode, bool, error) {
	filter := map[string]interface{}{
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGTE: metadata.Time{Time: time.Now().Add(-time.Duration(key.TTLSeconds()) * time.Second).UTC()},
		},
	}

	node := new(watch.ChainNode)
	err := c.watchDB.Table(key.ChainCollection()).Find(filter).Sort(common.BKFieldID).One(kit.Ctx, node)
	if err != nil {
		blog.ErrorJSON("get chain node from mongo failed, err: %s, collection: %s, filter: %s, rid: %s", err,
			key.ChainCollection(), filter, kit.Rid)
		if !c.watchDB.IsNotFoundError(err) {
			return nil, false, fmt.Errorf("get first chain node from mongo failed, err: %v", err)
		}
		return nil, false, nil
	}
	return node, true, nil
}

// getEventDetail get event detail with the needed fields by chain node, first get from redis, if failed, get from mongo
func (c *Client) getEventDetail(kit *rest.Kit, node *watch.ChainNode, fields []string, key event.Key) (*string,
	bool, error) {

	detail, err := c.getEventDetailFromRedis(kit, node.Cursor, fields, key)
	if err == nil {
		return detail, true, nil
	}

	return c.getEventDetailFromMongo(kit, node, fields, key)
}

// getEventDetailFromRedis get event detail with the needed fields by cursor from redis
func (c *Client) getEventDetailFromRedis(kit *rest.Kit, cursor string, fields []string, key event.Key) (
	*string, error) {

	detail, err := c.cache.Get(kit.Ctx, key.DetailKey(cursor)).Result()
	if err != nil {
		blog.Errorf("get watch detail from redis failed, err: %v, cursor: %s", err, cursor)
		return nil, fmt.Errorf("get event detail from redis failed, err: %v", err)
	}

	jsonStr := types.GetEventDetail(&detail)
	return json.CutJsonDataWithFields(jsonStr, fields), nil
}

// getEventDetailFromMongo get event detail with the needed fields by chain node from mongo
func (c *Client) getEventDetailFromMongo(kit *rest.Kit, node *watch.ChainNode, fields []string, key event.Key) (
	*string, bool, error) {

	// get delete events' detail with oid from cmdb
	if node.EventType == watch.Delete {
		filter := map[string]interface{}{
			"oid":  node.Oid,
			"coll": key.Collection(),
		}

		detailFields := make([]string, len(fields))
		for index, field := range fields {
			detailFields[index] = "detail." + field
		}

		if key.Collection() == common.BKTableNameBaseHost {
			doc := new(event.HostArchive)
			err := c.db.Table(common.BKTableNameDelArchive).Find(filter).Fields(detailFields...).One(kit.Ctx, doc)
			if err != nil {
				blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oid: %, err: %v",
					key.Collection(), node.Oid, err)
				if c.db.IsNotFoundError(err) {
					return nil, false, nil
				}
				return nil, false, fmt.Errorf("get archive deleted doc from mongo failed, err: %v", err)
			}

			byt, err := json.Marshal(doc.Detail)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					key.Collection(), doc.Oid, err)
				return nil, false, fmt.Errorf("marshal detail failed, err: %v", err)
			}
			detail := string(byt)
			return &detail, true, nil
		} else {
			doc := new(bsonx.Doc)
			err := c.db.Table(common.BKTableNameDelArchive).Find(filter).Fields(detailFields...).One(kit.Ctx, doc)
			if err != nil {
				blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oid: %, err: %v",
					key.Collection(), node.Oid, err)
				if c.db.IsNotFoundError(err) {
					return nil, false, nil
				}
				return nil, false, fmt.Errorf("get archive deleted doc from mongo failed, err: %v", err)
			}

			byt, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					key.Collection(), doc.Lookup("oid").String(), err)
				return nil, false, fmt.Errorf("marshal detail failed, err: %v", err)
			}
			detail := string(byt)
			return &detail, true, nil
		}
	}

	// get current detail from mongodb with oid
	objectId, err := primitive.ObjectIDFromHex(node.Oid)
	if err != nil {
		blog.ErrorJSON("get mongodb _id from oid(%s) failed, err: %s, rid: %s", node.Oid, err, kit.Rid)
		return nil, false, fmt.Errorf("get mongodb _id from oid(%s) failed, err: %s", node.Oid, err)
	}
	filter := map[string]interface{}{
		"_id": objectId,
	}

	var detailMap interface{}
	if key.Collection() == common.BKTableNameBaseHost {
		detailMap = new(metadata.HostMapStr)
	} else {
		detailMap = new(map[string]interface{})
	}

	if err := c.db.Table(key.Collection()).Find(filter).Fields(fields...).One(kit.Ctx, detailMap); err != nil {
		blog.ErrorJSON("get detail from db failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
		if c.db.IsNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get event detail from mongo failed, err: %v", err)
	}

	detailJson, _ := json.Marshal(detailMap)
	detail := string(detailJson)
	return &detail, true, nil
}

// SearchFollowingEventNodes search nodes after the node(excluding itself) by cursor and resource
func (c *Client) SearchFollowingEventChainNodes(kit *rest.Kit, opts *metadata.SearchEventNodesOption) (
	*metadata.EventNodes, error) {

	if opts.Limit > common.BKMaxPageSize {
		return nil, kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "limit", common.BKMaxPageSize)
	}

	key, err := event.GetResourceKeyWithCursorType(opts.Resource)
	if err != nil {
		blog.Errorf("get resource key with cursor type %s failed, err: %v, rid: %s", opts.Resource, err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_resource")
	}

	exists, nodes, _, err := c.searchFollowingEventChainNodes(kit, opts.StartCursor, uint64(opts.Limit), nil, key)
	if err != nil {
		blog.Errorf("search nodes after cursor %s failed, err: %v, rid: %s", opts.StartCursor, err, kit.Rid)
		return nil, err
	}

	return &metadata.EventNodes{Nodes: nodes, ExistsStartNode: exists}, nil
}

// searchFollowingEventNodes search nodes after the node(excluding itself) by cursor
func (c *Client) searchFollowingEventChainNodes(kit *rest.Kit, startCursor string, limit uint64, types []watch.EventType,
	key event.Key) (bool, []*watch.ChainNode, uint64, error) {

	// if start cursor is no event cursor, start from the beginning
	if startCursor == watch.NoEventCursor {
		node, exists, err := c.getEarliestEvent(kit, key)
		if err != nil {
			blog.Errorf("get earliest event for kwy %s failed, err: %v", key.Namespace(), err)
			return false, nil, 0, err
		}

		if !exists {
			return false, make([]*watch.ChainNode, 0), 0, nil
		}

		nodes, err := c.searchFollowingEventChainNodesByID(kit, node.ID, limit, types, key)
		if err != nil {
			return false, nil, 0, err
		}

		if c.isNodeWithEventType(node, types) {
			return true, append([]*watch.ChainNode{node}, nodes...), node.ID, nil
		}
		return true, nodes, node.ID, nil
	}

	filter := map[string]interface{}{
		common.BKCursorField: startCursor,
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGTE: metadata.Time{Time: time.Now().Add(-time.Duration(key.TTLSeconds()) * time.Second).UTC()},
		},
	}

	node := new(watch.ChainNode)
	err := c.watchDB.Table(key.ChainCollection()).Find(filter).Fields(common.BKFieldID).One(kit.Ctx, node)
	if err != nil {
		blog.ErrorJSON("get chain node from mongo failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
		if !c.watchDB.IsNotFoundError(err) {
			return false, nil, 0, err
		}

		filter := map[string]interface{}{
			"_id":                key.Collection(),
			common.BKCursorField: startCursor,
		}

		data := new(watch.LastChainNodeData)
		err := c.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(common.BKFieldID).One(kit.Ctx, data)
		if err != nil {
			blog.ErrorJSON("get last watch id failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
			if !c.watchDB.IsNotFoundError(err) {
				return false, nil, 0, err
			}
			return false, nil, 0, nil
		}

		nodes, err := c.searchFollowingEventChainNodesByID(kit, data.ID, limit, types, key)
		if err != nil {
			return false, nil, 0, err
		}
		return true, nodes, data.ID, nil
	}

	nodes, err := c.searchFollowingEventChainNodesByID(kit, node.ID, limit, types, key)
	if err != nil {
		return false, nil, 0, err
	}
	return true, nodes, node.ID, nil
}

// searchFollowingEventChainNodes search nodes after the node(excluding itself) by id
func (c *Client) searchFollowingEventChainNodesByID(kit *rest.Kit, id uint64, limit uint64, types []watch.EventType,
	key event.Key) ([]*watch.ChainNode, error) {

	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{common.BKDBGT: id},
	}

	if len(types) > 0 {
		filter[common.BKEventTypeField] = map[string]interface{}{common.BKDBIN: types}
	}

	nodes := make([]*watch.ChainNode, 0)
	if err := c.watchDB.Table(key.ChainCollection()).Find(filter).Sort(common.BKFieldID).Limit(limit).
		All(kit.Ctx, &nodes); err != nil {
		blog.Errorf("get chain nodes from mongo failed, err: %v, start id: %d, rid: %s", err, id, kit.Rid)
		return nil, fmt.Errorf("get chain nodes from mongo failed, err: %v, start id: %d", err, id)
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

	details, errCursors, errCursorIndexMap, err := c.searchEventDetailsFromRedis(kit, opts.Cursors, key)
	if err != nil {
		return nil, err
	}
	if len(errCursors) == 0 {
		return details, nil
	}

	// get event chain nodes from db for cursors that failed when reading redis
	chainFilter := map[string]interface{}{
		common.BKCursorField: map[string]interface{}{common.BKDBIN: errCursors},
	}
	nodes := make([]*watch.ChainNode, 0)
	if err := c.watchDB.Table(key.ChainCollection()).Find(chainFilter).Fields(common.BKCursorField,
		common.BKOIDField).All(kit.Ctx, &nodes); err != nil {
		blog.Errorf("get chain nodes failed, err: %v, cursor: %+v, rid: %s", err, errCursors, kit.Rid)
		return nil, fmt.Errorf("get chain nodes from mongo failed, err: %v, cursor: %+v", err, errCursors)
	}

	indexDetailMap, err := c.searchEventDetailsFromMongo(kit, nodes, make([]string, 0), errCursorIndexMap, key)
	if err != nil {
		blog.Errorf("get details from mongo failed, err: %v, cursors: %+v, rid: %s", err, errCursors, kit.Rid)
		return nil, err
	}

	for index, detail := range indexDetailMap {
		details[index] = `{"detail":` + detail + "}"
	}

	return details, nil
}

/** searchEventDetailsFromRedis get event details by cursors from redis, record the failed ones.
  returns the details that are successfully acquired from redis, the cursors with no detail in redis, and their map to
  the index in detail array so that we can get detail from mongo and fill them into the proper location in detail array
*/
func (c *Client) searchEventDetailsFromRedis(kit *rest.Kit, cursors []string, key event.Key) (
	[]string, []string, map[string]int, error) {

	if len(cursors) == 0 {
		return make([]string, 0), make([]string, 0), make(map[string]int), nil
	}

	detailKeys := make([]string, len(cursors))
	for idx, cursor := range cursors {
		detailKeys[idx] = key.DetailKey(cursor)
	}

	results, err := c.cache.MGet(kit.Ctx, detailKeys...).Result()
	if err != nil {
		blog.Errorf("search event details by cursors(%+v) failed, err: %v, rid: %s", cursors, err, kit.Rid)
		return nil, nil, nil, fmt.Errorf("search event details by cursors(%+v) failed, err: %v", cursors, err)
	}

	details := make([]string, len(results))
	errCursorIndexMap := make(map[string]int)
	errCursors := make([]string, 0)
	for index, result := range results {
		if result == nil {
			cursor := cursors[index]
			blog.Errorf("event detail for cursor(%s) do not exist in redis, rid: %s", cursor, kit.Rid)
			errCursorIndexMap[cursor] = index
			errCursors = append(errCursors, cursor)
			continue
		}

		resultStr, ok := result.(string)
		if !ok {
			blog.ErrorJSON("search event details from redis, but result(%s) is not string", result)
			return nil, nil, nil, fmt.Errorf("search event details from redis, but result is not string")
		}

		details[index] = resultStr
	}
	return details, errCursors, errCursorIndexMap, nil
}

// searchEventDetailsFromMongo get map of array index and detail with the needed fields by nodes from mongo
func (c *Client) searchEventDetailsFromMongo(kit *rest.Kit, nodes []*watch.ChainNode, fields []string,
	errCursorIndexMap map[string]int, key event.Key) (map[int]string, error) {

	if len(nodes) == 0 {
		return make(map[int]string, 0), nil
	}

	// get oids and its mapping with the detail array indexes
	oids := make([]primitive.ObjectID, 0)
	deletedOids := make([]string, 0)
	oidIndexMap := make(map[string][]int)
	for _, node := range nodes {
		if node.EventType == watch.Delete {
			deletedOids = append(deletedOids, node.Oid)
		} else {
			objectId, err := primitive.ObjectIDFromHex(node.Oid)
			if err != nil {
				blog.Errorf("get mongodb _id from oid(%s) failed, err: %v, rid: %s", node.Oid, err, kit.Rid)
				return nil, fmt.Errorf("get mongodb _id from oid(%s) failed, err: %v", node.Oid, err)
			}
			oids = append(oids, objectId)
		}

		oidIndexMap[node.Oid] = append(oidIndexMap[node.Oid], errCursorIndexMap[node.Cursor])
	}
	coll := key.Collection()
	oidDetailMap := make(map[int]string)

	// get details in collection by oids, need to get _id to do mapping
	if len(oids) > 0 {
		filter := map[string]interface{}{
			"_id": map[string]interface{}{common.BKDBIN: oids},
		}

		findOpts := daltypes.FindOpts{
			WithObjectID: true,
		}

		fields := fields
		if len(fields) > 0 {
			fields = append(fields, "_id")
		}

		if coll == common.BKTableNameBaseHost {
			detailArr := make([]metadata.HostMapStr, 0)
			if err := c.db.Table(coll).Find(filter, findOpts).Fields(fields...).All(kit.Ctx, &detailArr); err != nil {
				blog.Errorf("get details from db failed, err: %s, oids: %+v, rid: %s", err, oids, kit.Rid)
				return nil, fmt.Errorf("get details from mongo failed, err: %v, oids: %+v", err, oids)
			}

			for _, detailMap := range detailArr {
				objectID, ok := detailMap["_id"].(primitive.ObjectID)
				if !ok {
					return nil, fmt.Errorf("parse detail oid failed, oid: %+v", detailMap["_id"])
				}
				oid := objectID.Hex()
				delete(detailMap, "_id")
				detailJson, _ := json.Marshal(detailMap)
				for _, index := range oidIndexMap[oid] {
					oidDetailMap[index] = string(detailJson)
				}
			}
		} else {
			detailArr := make([]mapStrWithOid, 0)
			if err := c.db.Table(coll).Find(filter, findOpts).Fields(fields...).All(kit.Ctx, &detailArr); err != nil {
				blog.Errorf("get details from db failed, err: %s, oids: %+v, rid: %s", err, oids, kit.Rid)
				return nil, fmt.Errorf("get details from mongo failed, err: %v, oids: %+v", err, oids)
			}

			for _, detailMap := range detailArr {
				oid := detailMap.Oid.Hex()
				detailJson, _ := json.Marshal(detailMap.MapStr)
				for _, index := range oidIndexMap[oid] {
					oidDetailMap[index] = string(detailJson)
				}
			}
		}
	}

	if len(deletedOids) == 0 {
		return oidDetailMap, nil
	}

	oidDetailMap, err := c.searchDeletedEventDetailsFromMongo(kit, coll, deletedOids, fields, oidIndexMap, oidDetailMap)
	if err != nil {
		blog.Errorf("get delete details from db failed, err: %s, oids: %+v, rid: %s", err, deletedOids, kit.Rid)
		return nil, err
	}

	return oidDetailMap, nil
}

// searchDeletedEventDetailsFromMongo search delete events' details from the cc_DelArchive table by oids
func (c *Client) searchDeletedEventDetailsFromMongo(kit *rest.Kit, coll string, deletedOids []string, fields []string,
	oidIndexMap map[string][]int, oidDetailMap map[int]string) (map[int]string, error) {

	detailFields := make([]string, len(fields))
	for index, field := range fields {
		detailFields[index] = "detail." + field
	}
	detailFields = append(detailFields, "oid")

	deleteFilter := map[string]interface{}{
		"oid":  map[string]interface{}{common.BKDBIN: deletedOids},
		"coll": coll,
	}

	if coll == common.BKTableNameBaseHost {
		docs := make([]event.HostArchive, 0)
		err := c.db.Table(common.BKTableNameDelArchive).Find(deleteFilter).Fields(detailFields...).All(kit.Ctx, &docs)
		if err != nil {
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v, "+
				"rid: %s", coll, deletedOids, err, kit.Rid)
			return nil, fmt.Errorf("get archive deleted docs from mongo failed, err: %v, oids: %+v", err, deletedOids)
		}

		for _, doc := range docs {
			detailJson, _ := json.Marshal(doc.Detail)
			for _, index := range oidIndexMap[doc.Oid] {
				oidDetailMap[index] = string(detailJson)
			}
		}
	} else {
		docs := make([]bsonx.Doc, 0)
		err := c.db.Table(common.BKTableNameDelArchive).Find(deleteFilter).Fields(detailFields...).All(kit.Ctx, &docs)
		if err != nil {
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v, "+
				"rid: %s", coll, deletedOids, err, kit.Rid)
			return nil, fmt.Errorf("get archive deleted docs from mongo failed, err: %v, oids: %+v", err, deletedOids)
		}

		for _, doc := range docs {
			oid := doc.Lookup("oid").String()
			detailJson, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
			if err != nil {
				blog.Errorf("marshal detail failed, oid: %s, err: %v, rid: %s", oid, err, kit.Rid)
				return nil, fmt.Errorf("marshal detail failed, oid: %s, err: %v", oid, err)
			}
			for _, index := range oidIndexMap[oid] {
				oidDetailMap[index] = string(detailJson)
			}
		}
	}

	return oidDetailMap, nil
}

type mapStrWithOid struct {
	Oid    primitive.ObjectID     `bson:"_id"`
	MapStr map[string]interface{} `bson:",inline"`
}
