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
	"context"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/dal"
	daltypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// dbChainTTLTime the ttl time seconds of the db event chain, used to set the ttl index of mongodb
const dbChainTTLTime = 48 * 60 * 60

func newFlow(ctx context.Context, opts FlowOptions) error {
	flow := Flow{
		FlowOptions: opts,
		metrics:     initialMetrics(opts.key.Collection()),
	}

	if err := flow.createDBChainCollection(ctx); err != nil {
		return err
	}

	flow.cleanDelArchiveData(ctx)

	return flow.RunFlow(ctx)
}

func (f *Flow) createDBChainCollection(ctx context.Context) error {
	exists, err := f.watchDB.HasTable(ctx, f.key.ChainCollection())
	if err != nil {
		blog.Errorf("check if table %s exists failed, err: %v", f.key.ChainCollection(), err)
		return err
	}

	if !exists {
		if err = f.watchDB.CreateTable(ctx, f.key.ChainCollection()); err != nil && !f.watchDB.IsDuplicatedError(err) {
			blog.Errorf("create table %s failed, err: %v", f.key.ChainCollection(), err)
			return err
		}
	}

	indexes := []daltypes.Index{
		{Name: "index_id", Keys: map[string]int32{common.BKFieldID: -1}, Background: true},
		{Name: "index_cursor", Keys: map[string]int32{common.BKCursorField: -1}, Background: true, Unique: true},
		{Name: "index_cluster_time", Keys: map[string]int32{common.BKClusterTimeField: -1}, Background: true,
			ExpireAfterSeconds: dbChainTTLTime},
	}

	existIndexArr, err := f.watchDB.Table(f.key.ChainCollection()).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist indexes for event chain table %s failed, err: %v", f.key.ChainCollection(), err)
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		err = f.watchDB.Table(f.key.ChainCollection()).CreateIndex(ctx, index)
		if err != nil && !f.watchDB.IsDuplicatedError(err) {
			blog.Errorf("create indexes for event chain table %s failed, err: %v", f.key.ChainCollection(), err)
			return err
		}
	}

	return nil
}

type Flow struct {
	FlowOptions
	metrics *eventMetrics
}

const batchSize = 200

func (f *Flow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	event := make(map[string]interface{})
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct:     &event,
			Collection:      f.key.Collection(),
			StartAfterToken: nil,
		},
	}
	if f.key.Collection() == common.BKTableNameBaseHost {
		watchOpts.EventStruct = &metadata.HostMapStr{}
	}

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         f.key.Namespace(),
			WatchOpt:     watchOpts,
			TokenHandler: NewFlowTokenHandler(f.key, f.watchDB, f.metrics),
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: batchSize,
	}

	if err := f.watch.WithBatch(opts); err != nil {
		blog.Errorf("run flow, but watch batch failed, err: %v", err)
		return err
	}

	return nil
}

func (f *Flow) doBatch(es []*types.Event) (retry bool) {
	eventLen := len(es)
	chainNodes := make([]watch.ChainNode, eventLen)
	oids := make([]string, eventLen)
	hasError := true

	// collect event related metrics
	start := time.Now()
	defer func() {
		if retry {
			f.metrics.collectRetryError()
		}
		if hasError {
			return
		}
		f.metrics.collectCycleDuration(time.Since(start) / time.Duration(eventLen))
	}()

	ids, err := f.watchDB.NextSequences(context.Background(), f.key.ChainCollection(), eventLen)
	if err != nil {
		blog.Errorf("get event ids for collection: %s failed, err: %v", f.key.ChainCollection(), err)
		return true
	}

	oidDetailMap, retry, err := f.getDeleteEventDetails(es)
	if err != nil {
		blog.Errorf("get deleted event details failed, err: %v", err)
		return retry
	}

	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	for index, e := range es {
		// collect event's basic metrics
		f.metrics.collectBasic(e)

		oids[index] = e.Oid
		id := ids[index]
		name := f.key.Name(e.DocBytes)
		currentCursor, err := watch.GetEventCursor(f.key.Collection(), e, id)
		if err != nil {
			blog.Errorf("get event cursor failed, name: %s, err: %v, oid: %s ", name, err, e.Oid)
			return false
		}

		chainNode := watch.ChainNode{
			ID:          id,
			ClusterTime: e.ClusterTime,
			Oid:         e.Oid,
			EventType:   watch.ConvertOperateType(e.OperationType),
			Token:       e.Token.Data,
			Cursor:      currentCursor,
		}

		if instanceID := f.key.InstanceID(e.DocBytes); instanceID > 0 {
			chainNode.InstanceID = instanceID
		}
		chainNodes[index] = chainNode

		docBytes := e.DocBytes
		if e.OperationType == types.Delete {
			docBytes = oidDetailMap[e.Oid]
		}

		detail := types.EventDetail{
			Detail:        types.JsonString(docBytes),
			UpdatedFields: e.ChangeDesc.UpdatedFields,
			RemovedFields: e.ChangeDesc.RemovedFields,
		}
		detailBytes, err := json.Marshal(detail)
		if err != nil {
			blog.Errorf("run flow, marshal detail failed, name: %s, detail: %+v, err: %v, oid: %s", name,
				detail, err, e.Oid)
			return false
		}
		pipe.Set(f.key.DetailKey(currentCursor), string(detailBytes), time.Duration(f.key.TTLSeconds())*time.Second)
	}

	getLock := f.getLockWithRetry(f.key.LockKey(), oids)
	if !getLock {
		blog.Errorf("run flow, insert nodes, do not get lock and return, oids: %+v", oids)
		return true
	}

	// already get the lock. prepare to release the lock.
	defer func() {
		if err := redis.Client().Del(context.Background(), f.key.LockKey()).Err(); err != nil {
			f.metrics.collectLockError()
			blog.ErrorJSON("run flow, insert nodes, release lock failed, err: %s, oids: %+v", err, oids)
		}
	}()

	// store details at first, in case those watching cmdb events read chain when details are not inserted yet
	if _, err := pipe.Exec(); err != nil {
		f.metrics.collectRedisError()
		blog.Errorf("run flow, but insert details failed, err: %v, oids: %+v", err, oids)
		return true
	}

	// store chain node to db, if has duplicate events, update their ids, then insert the new events
	if err := f.watchDB.Table(f.key.ChainCollection()).Insert(context.Background(), chainNodes); err != nil {
		if !strings.Contains(err.Error(), "duplicate key error") {
			f.metrics.collectMongoError()
			blog.Errorf("run flow, but insert chain nodes failed, err: %v, oids: %+v", err, oids)
			return true
		}

		if retry, err := f.handleInsertChainNodeDuplicateError(chainNodes); err != nil {
			blog.Errorf("handle insert chain node duplicated error failed, err: %v", err)
			return retry
		}
	}

	blog.Infof("insert watch event for %s success, oids: %+v", f.key.Collection(), oids)
	hasError = false
	return false
}

func (f *Flow) getLockWithRetry(name string, oids []string) bool {
	getLock := false
	for retry := 0; retry < 10; retry++ {
		// get operate lock to avoid concurrent revise the chain
		success, err := redis.Client().SetNX(context.Background(), f.key.LockKey(), "lock", 5*time.Second).Result()
		if err != nil {
			f.metrics.collectLockError()
			blog.Errorf("get lock failed, name: %s, err: %v, oids: %+v", name, err, oids)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if !success {
			blog.Warnf("do not get lock, name: %s, oids: %+v", name, oids)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		getLock = true
		break
	}
	return getLock
}

// getDeleteEventDetails get delete events' oid and related detail map from cmdb
func (f *Flow) getDeleteEventDetails(es []*types.Event) (map[string][]byte, bool, error) {
	oidDetailMap := make(map[string][]byte)

	deletedEventOids := make([]string, 0)
	for _, e := range es {
		if e.OperationType == types.Delete {
			deletedEventOids = append(deletedEventOids, e.Oid)
		}
	}

	if len(deletedEventOids) == 0 {
		return oidDetailMap, false, nil
	}

	filter := map[string]interface{}{
		"oid":  map[string]interface{}{common.BKDBIN: deletedEventOids},
		"coll": f.key.Collection(),
	}

	if f.key.Collection() == common.BKTableNameBaseHost {
		docs := make([]HostArchive, 0)
		err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
		if err != nil {
			f.metrics.collectMongoError()
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
				f.key.Collection(), deletedEventOids, err)
			return nil, true, err
		}

		for _, doc := range docs {
			byt, err := json.Marshal(doc.Detail)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					f.key.Collection(), doc.Oid, err)
				return nil, false, err
			}
			oidDetailMap[doc.Oid] = byt
		}
	} else {
		docs := make([]bsonx.Doc, 0)
		err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
		if err != nil {
			f.metrics.collectMongoError()
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
				f.key.Collection(), deletedEventOids, err)
			return nil, true, err
		}

		for _, doc := range docs {
			byt, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					f.key.Collection(), doc.Lookup("oid").String(), err)
				return nil, false, err
			}
			oidDetailMap[doc.Lookup("oid").String()] = byt
		}
	}

	return oidDetailMap, false, nil
}

/** handleInsertChainNodeDuplicateError if an duplicate error occurred when inserting chain nodes, we find all exist
  nodes at first, update their ids to the bigger one(if several masters exists in cases like switching master & slave,
  the fist one that gets the ids might not get the lock to insert nodes, when it gets the lock, others have already
  inserted former nodes with a larger id, we need to update the exist nodes' ids to ensure the sequence of them), and
  then insert the newer nodes one by one in case another duplicate error occurred.
*/
func (f *Flow) handleInsertChainNodeDuplicateError(chainNodes []watch.ChainNode) (bool, error) {
	cursors := make([]string, len(chainNodes))
	chainNodeMap := make(map[string]watch.ChainNode)
	for index, chainNode := range chainNodes {
		cursors[index] = chainNode.Cursor
		chainNodeMap[chainNode.Cursor] = chainNode
	}
	filter := map[string]interface{}{
		common.BKCursorField: map[string]interface{}{common.BKDBIN: cursors},
	}

	existNodes := make([]watch.ChainNode, 0)
	if err := f.watchDB.Table(f.key.ChainCollection()).Find(filter).All(context.Background(), &existNodes); err != nil {
		f.metrics.collectMongoError()
		blog.Errorf("run flow, but get chain nodes failed, err: %v, cursors: %+v", err, cursors)
		return true, err
	}

	for _, existNode := range existNodes {
		updateFilter := map[string]interface{}{
			common.BKCursorField: existNode.Cursor,
		}
		updateData := map[string]interface{}{
			common.BKFieldID: chainNodeMap[existNode.Cursor].ID,
		}

		err := f.watchDB.Table(f.key.ChainCollection()).Update(context.Background(), updateFilter, updateData)
		if err != nil {
			f.metrics.collectMongoError()
			blog.Errorf("run flow, but update exist chain node(%s) failed, err: %s", existNode, err)
			return true, err
		}

		delete(chainNodeMap, existNode.Cursor)
	}

	for _, chainNode := range chainNodeMap {
		if err := f.watchDB.Table(f.key.ChainCollection()).Insert(context.Background(), chainNode); err != nil {
			if strings.Contains(err.Error(), "duplicate key error") {
				continue
			}
			f.metrics.collectMongoError()
			blog.ErrorJSON("run flow, but insert chain nodes failed, err: %v, chain node: %s", err, chainNode)
			return true, err
		}
	}

	return false, nil
}

type flowTokenHandler struct {
	key     Key
	watchDB dal.DB
	metrics *eventMetrics
}

func NewFlowTokenHandler(key Key, watchDB dal.DB, metrics *eventMetrics) types.TokenHandler {
	return &flowTokenHandler{
		key:     key,
		watchDB: watchDB,
		metrics: metrics,
	}
}

// SetLastWatchToken set start watch token to redis
func (f *flowTokenHandler) SetLastWatchToken(ctx context.Context, token string) error {
	lastTokenKey := f.key.LastTokenKey()

	err := redis.Client().Set(ctx, lastTokenKey, token, time.Duration(f.key.TTLSeconds())*time.Second).Err()
	if err != nil {
		f.metrics.collectRedisError()
		blog.Errorf("set last watch token failed, key: %s, token: %s, err: %v", lastTokenKey, token, err)
		return err
	}

	return nil
}

// GetStartWatchToken get start watch token from redis first, if an error occurred, get from db and set redis
func (f *flowTokenHandler) GetStartWatchToken(ctx context.Context) (token string, err error) {
	lastTokenKey := f.key.LastTokenKey()

	token, err = redis.Client().Get(context.Background(), lastTokenKey).Result()
	if err != nil {
		blog.Errorf("get last watch token from redis failed, key: %s, err: %v", lastTokenKey, err)
		if !redis.IsNilErr(err) {
			f.metrics.collectRedisError()
		}

		tailNode := new(watch.ChainNode)
		if err := f.watchDB.Table(f.key.ChainCollection()).Find(map[string]interface{}{}).Fields(common.BKTokenField).
			Sort(common.BKFieldID+":-1").One(context.Background(), tailNode); err != nil {

			if !f.watchDB.IsNotFoundError(err) {
				f.metrics.collectMongoError()
				blog.Errorf("get last watch token from mongo failed,err: %v", err)
				return "", err
			}
			// the tail node is not exist.
			return "", nil
		}

		err := redis.Client().Set(ctx, lastTokenKey, token, time.Duration(f.key.TTLSeconds())*time.Second).Err()
		if err != nil {
			f.metrics.collectRedisError()
			blog.Errorf("set last watch token failed, key: %s, token: %s, err: %v", lastTokenKey, token, err)
		}

		return tailNode.Token, nil
	}

	return token, nil
}
