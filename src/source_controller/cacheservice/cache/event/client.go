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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	daltypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"

	rawRedis "github.com/go-redis/redis/v7"
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
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_resource")
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
	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

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
			return nil, false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		return nil, false, nil
	}
	return node, true, nil
}

// getEarliestEvent search earliest event chain node in not expired nodes
func (c *Client) getEarliestEvent(kit *rest.Kit, key event.Key) (*watch.ChainNode, bool, error) {
	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	filter := map[string]interface{}{
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGTE: metadata.Time{Time: time.Now().Add(-time.Duration(key.TTLSeconds()) * time.Second).UTC()},
		},
	}

	node := new(watch.ChainNode)
	err := c.watchDB.Table(key.ChainCollection()).Find(filter).Sort(common.BKFieldID).One(kit.Ctx, node)
	if err != nil {
		blog.ErrorJSON("get chain node from mongo failed, err: %s, filter: %s, rid: %s", err, filter, kit.Rid)
		if !c.watchDB.IsNotFoundError(err) {
			return nil, false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		return nil, false, nil
	}
	return node, true, nil
}

// getEventDetail get event detail with the needed fields by chain node, first get from redis, if failed, get from mongo
func (c *Client) getEventDetail(kit *rest.Kit, node *watch.ChainNode, fields []string, key event.Key) (string,
	bool, error) {

	detail, err := c.getEventDetailFromRedis(kit, node.Cursor, fields, key)
	if err == nil {
		return detail, true, nil
	}

	return c.getEventDetailFromMongo(kit, node, fields, key)
}

// getEventDetailFromRedis get event detail with the needed fields by cursor from redis
func (c *Client) getEventDetailFromRedis(kit *rest.Kit, cursor string, fields []string, key event.Key) (
	string, error) {

	detail, err := c.cache.Get(kit.Ctx, key.DetailKey(cursor)).Result()
	if err != nil {
		blog.Errorf("get watch detail from redis failed, err: %v, cursor: %s", err, cursor)
		return "", kit.CCError.CCError(common.CCErrCommRedisOPErr)
	}

	jsonStr := types.GetEventDetail(detail)
	return *json.CutJsonDataWithFields(&jsonStr, fields), nil
}

// getEventDetailFromMongo get event detail with the needed fields by chain node from mongo
func (c *Client) getEventDetailFromMongo(kit *rest.Kit, node *watch.ChainNode, fields []string, key event.Key) (
	string, bool, error) {

	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	// get delete events' detail with oid from cmdb
	if node.EventType == watch.Delete {
		filter := map[string]interface{}{
			"oid":  node.Oid,
			"coll": key.Collection(),
		}

		if key.Collection() == common.BKTableNameBaseHost {
			doc := new(event.HostArchive)
			err := c.db.Table(common.BKTableNameDelArchive).Find(filter).One(kit.Ctx, doc)
			if err != nil {
				blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oid: %, err: %v",
					key.Collection(), node.Oid, err)
				if c.db.IsNotFoundError(err) {
					return "", false, nil
				}
				return "", false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}

			byt, err := json.Marshal(doc.Detail)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					key.Collection(), doc.Oid, err)
				return "", false, kit.CCError.CCError(common.CCErrCommJSONMarshalFailed)
			}
			detail := string(byt)
			return *json.CutJsonDataWithFields(&detail, fields), true, nil
		} else {
			doc := new(bsonx.Doc)
			err := c.db.Table(common.BKTableNameDelArchive).Find(filter).One(kit.Ctx, doc)
			if err != nil {
				blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oid: %, err: %v",
					key.Collection(), node.Oid, err)
				if c.db.IsNotFoundError(err) {
					return "", false, nil
				}
				return "", false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}

			byt, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					key.Collection(), doc.Lookup("oid").String(), err)
				return "", false, kit.CCError.CCError(common.CCErrCommJSONMarshalFailed)
			}
			detail := string(byt)
			return *json.CutJsonDataWithFields(&detail, fields), true, nil
		}
	}

	// get current detail from mongodb with oid
	objectId, err := primitive.ObjectIDFromHex(node.Oid)
	if err != nil {
		blog.ErrorJSON("get mongodb _id from oid(%s) failed, err: %s, rid: %s", node.Oid, err, kit.Rid)
		return "", false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "oid")
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
			return "", false, nil
		}
		return "", false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	detailJson, _ := json.Marshal(detailMap)
	return string(detailJson), true, nil
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

	exists, nodes, err := c.searchFollowingEventChainNodes(kit, opts.StartCursor, uint64(opts.Limit), key)
	if err != nil {
		blog.Errorf("search nodes after cursor %s failed, err: %v, rid: %s", opts.StartCursor, err, kit.Rid)
		return nil, err
	}

	return &metadata.EventNodes{Nodes: nodes, ExistsStartNode: exists}, nil
}

// searchFollowingEventNodes search nodes after the node(excluding itself) by cursor
func (c *Client) searchFollowingEventChainNodes(kit *rest.Kit, startCursor string, limit uint64, key event.Key) (bool,
	[]*watch.ChainNode, error) {

	// if start cursor is no event cursor, start from the beginning
	if startCursor == watch.NoEventCursor {
		node, exists, err := c.getEarliestEvent(kit, key)
		if err != nil {
			blog.Errorf("get earliest event for kwy %s failed, err: %v", key.Namespace(), err)
			return false, nil, err
		}

		if !exists {
			return false, make([]*watch.ChainNode, 0), nil
		}

		nodes, err := c.searchFollowingEventChainNodesByID(kit, node.ID, limit, key)
		if err != nil {
			return false, nil, err
		}
		return true, append([]*watch.ChainNode{node}, nodes...), nil
	}

	cursor := new(watch.Cursor)
	if err := cursor.Decode(startCursor); err != nil {
		blog.Errorf("decode cursor %s failed, err: %v", startCursor, err)
		return false, nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "start_cursor")
	}

	exists, err := c.isEventChainNodeExist(kit, cursor.EventID, key)
	if err != nil {
		return false, nil, err
	}

	if !exists {
		return false, nil, nil
	}

	nodes, err := c.searchFollowingEventChainNodesByID(kit, cursor.EventID, limit, key)
	if err != nil {
		return false, nil, err
	}
	return true, nodes, nil
}

func (c *Client) isEventChainNodeExist(kit *rest.Kit, id uint64, key event.Key) (bool, error) {
	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	filter := map[string]interface{}{
		common.BKFieldID: id,
		common.BKClusterTimeField: map[string]interface{}{
			common.BKDBGTE: metadata.Time{Time: time.Now().Add(-time.Duration(key.TTLSeconds()) * time.Second).UTC()},
		},
	}

	cnt, err := c.watchDB.Table(key.ChainCollection()).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count chain node from mongo failed, err: %v, id: %d, rid: %s", err, id, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

// searchFollowingEventChainNodes search nodes after the node(excluding itself) by id
func (c *Client) searchFollowingEventChainNodesByID(kit *rest.Kit, id uint64, limit uint64, key event.Key) (
	[]*watch.ChainNode, error) {

	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{common.BKDBGT: id},
	}

	nodes := make([]*watch.ChainNode, 0)
	if err := c.watchDB.Table(key.ChainCollection()).Find(filter).Sort(common.BKFieldID).Limit(limit).
		All(kit.Ctx, &nodes); err != nil {
		blog.Errorf("get chain nodes from mongo failed, err: %v, start id: %d, rid: %s", err, id, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
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

	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	details, errCursors, errCursorIndexMap := c.searchEventDetailsFromRedis(kit, opts.Cursors, make([]string, 0), key)
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
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	indexDetailMap, err := c.searchEventDetailsFromMongo(kit, nodes, make([]string, 0), errCursorIndexMap, key)
	if err != nil {
		blog.Errorf("get details from mongo failed, err: %v, cursors: %+v, rid: %s", err, errCursors, kit.Rid)
		return nil, err
	}

	for index, detail := range indexDetailMap {
		details[index] = detail
	}

	return details, nil
}

// searchEventDetailsFromRedis get event details with the needed fields by cursors from redis, record the failed ones
func (c *Client) searchEventDetailsFromRedis(kit *rest.Kit, cursors []string, fields []string, key event.Key) (
	[]string, []string, map[string]int) {

	if len(cursors) == 0 {
		return make([]string, 0), make([]string, 0), make(map[string]int)
	}

	results := make([]*rawRedis.StringCmd, len(cursors))
	pipe := c.cache.Pipeline()
	for idx, cursor := range cursors {
		results[idx] = pipe.Get(key.DetailKey(cursor))
	}

	_, err := pipe.Exec()
	if err != nil {
		blog.Errorf("search event details by cursors(%+v) failed, err: %v, rid: %s", cursors, err, kit.Rid)
	}

	details := make([]string, len(results))
	errCursorIndexMap := make(map[string]int)
	errCursors := make([]string, 0)

	for idx, result := range results {
		if err := result.Err(); err != nil {
			cursor := cursors[idx]
			blog.Errorf("search event detail by cursor(%s) failed, err: %v, rid: %s", cursor, err, kit.Rid)
			errCursorIndexMap[cursor] = idx
			errCursors = append(errCursors, cursor)
			continue
		}
		jsonStr := types.GetEventDetail(result.Val())
		details[idx] = *json.CutJsonDataWithFields(&jsonStr, fields)
	}
	return details, errCursors, errCursorIndexMap
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
				blog.ErrorJSON("get mongodb _id from oid(%s) failed, err: %s, rid: %s", node.Oid, err, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "oid")
			}
			oids = append(oids, objectId)
		}

		oidIndexMap[node.Oid] = append(oidIndexMap[node.Oid], errCursorIndexMap[node.Cursor])
	}

	util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)
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
				return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}

			for _, detailMap := range detailArr {
				objectID, ok := detailMap["_id"].(primitive.ObjectID)
				if !ok {
					return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
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
				return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
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

	// get details in delete archive collection by oids
	if len(deletedOids) > 0 {
		deleteFilter := map[string]interface{}{
			"oid":  map[string]interface{}{common.BKDBIN: deletedOids},
			"coll": coll,
		}

		if coll == common.BKTableNameBaseHost {
			docs := make([]event.HostArchive, 0)
			err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(deleteFilter).All(kit.Ctx, &docs)
			if err != nil {
				blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
					coll, deletedOids, err)
				return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}

			for _, doc := range docs {
				detailJson, _ := json.Marshal(doc.Detail)
				detail := string(detailJson)
				detailWithFields := *json.CutJsonDataWithFields(&detail, fields)
				for _, index := range oidIndexMap[doc.Oid] {
					oidDetailMap[index] = detailWithFields
				}
			}
		} else {
			docs := make([]bsonx.Doc, 0)
			err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(deleteFilter).All(kit.Ctx, &docs)
			if err != nil {
				blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
					coll, deletedOids, err)
				return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}

			for _, doc := range docs {
				detailJson, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
				if err != nil {
					blog.Errorf("marshal detail failed, oid: %s, err: %v", doc.Lookup("oid").String(), err)
					return nil, kit.CCError.CCError(common.CCErrCommJSONMarshalFailed)
				}
				detail := string(detailJson)
				oid := doc.Lookup("oid").String()
				detailWithFields := *json.CutJsonDataWithFields(&detail, fields)
				for _, index := range oidIndexMap[oid] {
					oidDetailMap[index] = detailWithFields
				}
			}
		}
	}

	return oidDetailMap, nil
}

type mapStrWithOid struct {
	Oid    primitive.ObjectID     `bson:"_id"`
	MapStr map[string]interface{} `bson:",inline"`
}
