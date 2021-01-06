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

package flow

import (
	"context"
	"strings"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type flowOptions struct {
	key      event.Key
	watch    stream.LoopInterface
	isMaster discovery.ServiceManageInterface
	watchDB  *local.Mongo
	ccDB     dal.DB
}

func newFlow(ctx context.Context, opts flowOptions) error {
	flow := Flow{
		flowOptions: opts,
		metrics:     initialMetrics(opts.key.Collection()),
	}

	flow.cleanDelArchiveData(ctx)

	return flow.RunFlow(ctx)
}

type Flow struct {
	flowOptions
	metrics      *eventMetrics
	tokenHandler *flowTokenHandler
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

	f.tokenHandler = NewFlowTokenHandler(f.key, f.watchDB, f.metrics)

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         f.key.Namespace(),
			WatchOpt:     watchOpts,
			TokenHandler: f.tokenHandler,
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
	if eventLen == 0 {
		return false
	}

	chainNodes := make([]*watch.ChainNode, 0)
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

	oidDetailMap, retry, err := f.getDeleteEventDetails(es)
	if err != nil {
		blog.Errorf("get deleted event details failed, err: %v", err)
		return retry
	}

	ids, err := f.watchDB.NextSequences(context.Background(), f.key.ChainCollection(), eventLen)
	if err != nil {
		blog.Errorf("get event ids for collection: %s failed, err: %v", f.key.ChainCollection(), err)
		return true
	}

	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	lastTokenData := make(map[string]interface{})
	for index, e := range es {
		// collect event's basic metrics
		f.metrics.collectBasic(e)
		lastTokenData[common.BKTokenField] = e.Token.Data

		switch e.OperationType {
		case types.Insert, types.Update, types.Replace, types.Delete:
		case types.Invalidate:
			blog.Errorf("loop flow, received invalid event operation type, doc: %s", e.DocBytes)
			continue
		default:
			blog.Errorf("loop flow, received unsupported event operation type, doc: %s", e.DocBytes)
			continue
		}

		oids[index] = e.Oid
		id := ids[index]
		name := f.key.Name(e.DocBytes)
		currentCursor, err := watch.GetEventCursor(f.key.Collection(), e)
		if err != nil {
			blog.Errorf("get event cursor failed, name: %s, err: %v, oid: %s ", name, err, e.Oid)
			return false
		}

		chainNode := &watch.ChainNode{
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
		chainNodes = append(chainNodes, chainNode)

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

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		if err := f.tokenHandler.setLastWatchToken(context.Background(), lastTokenData); err != nil {
			f.metrics.collectMongoError()
			return false
		}
		return false
	}

	// store details at first, in case those watching cmdb events read chain when details are not inserted yet
	if _, err := pipe.Exec(); err != nil {
		f.metrics.collectRedisError()
		blog.Errorf("run flow, but insert details for %s failed, err: %v, oids: %+v", err, f.key.Collection(), oids)
		return true
	}

	watchDBClient := f.watchDB.GetDBClient()

	session, err := watchDBClient.StartSession()
	if err != nil {
		blog.Errorf("run flow, but start mongo session failed, err: %v, coll: %s", err, f.key.Collection())
		return true
	}
	defer session.EndSession(context.Background())

	if txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			blog.Errorf("run flow, but start mongo transaction failed, err: %v, coll: %s", err, f.key.Collection())
			return err
		}

		if err := f.watchDB.Table(f.key.ChainCollection()).Insert(sc, chainNodes); err != nil {
			blog.Errorf("run flow, but insert chain nodes for %s failed, err: %v, oids: %+v", f.key.Collection(), err, oids)
			if !strings.Contains(err.Error(), "duplicate key error") {
				f.metrics.collectMongoError()
			}
			_ = session.AbortTransaction(context.Background())
			return err
		}

		lastNode := chainNodes[len(chainNodes)-1]
		lastTokenData[common.BKFieldID] = lastNode.ID
		lastTokenData[common.BKCursorField] = lastNode.Cursor
		if err := f.tokenHandler.setLastWatchToken(sc, lastTokenData); err != nil {
			f.metrics.collectMongoError()
			_ = session.AbortTransaction(context.Background())
			return err
		}

		// Use context.Background() to ensure that the commit can complete successfully even if the context passed to
		// mongo.WithSession is changed to have a timeout.
		if err = session.CommitTransaction(context.Background()); err != nil {
			blog.Errorf("run flow, but commit mongo transaction failed, err: %v", err)
			return err
		}
		return nil
	}); txnErr != nil {
		// if an error occurred, roll back and re-watch again
		return true
	}

	blog.Infof("insert watch event for %s success, oids: %+v", f.key.Collection(), oids)
	hasError = false
	return false
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
		docs := make([]event.HostArchive, 0)
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

var _ = types.TokenHandler(&flowTokenHandler{})

type flowTokenHandler struct {
	key     event.Key
	watchDB dal.DB
	metrics *eventMetrics
}

func NewFlowTokenHandler(key event.Key, watchDB dal.DB, metrics *eventMetrics) *flowTokenHandler {
	return &flowTokenHandler{
		key:     key,
		watchDB: watchDB,
		metrics: metrics,
	}
}

/* SetLastWatchToken do not set last watch token in the do batch action(set it after events are successfully inserted)
   when there are several masters watching db event, we use db transaction to avoid inserting duplicate data by setting
   the last token after the insertion of db chain nodes in one transaction, since we have a unique index on the cursor
   field, the later one will encounters an error when inserting nodes and roll back without setting the token and watch
   another round from the last token of the last inserted node, thus ensures the sequence of db chain nodes.
*/
func (f *flowTokenHandler) SetLastWatchToken(ctx context.Context, token string) error {
	return nil
}

// setLastWatchToken set last watch token(used after events are successfully inserted)
func (f *flowTokenHandler) setLastWatchToken(ctx context.Context, data map[string]interface{}) error {
	filter := map[string]interface{}{
		"_id": f.key.Collection(),
	}
	if err := f.watchDB.Table(common.BKTableNameWatchToken).Update(context.Background(), filter, data); err != nil {
		blog.Errorf("set last watch token failed, err: %v, data: %+v", err, data)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db first, if an error occurred, get from chain db
func (f *flowTokenHandler) GetStartWatchToken(ctx context.Context) (token string, err error) {
	filter := map[string]interface{}{
		"_id": f.key.Collection(),
	}

	data := new(watch.LastChainNodeData)
	if err := f.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(common.BKTokenField).One(ctx, data); err != nil {
		if !f.watchDB.IsNotFoundError(err) {
			f.metrics.collectMongoError()
			blog.ErrorJSON("run flow, but get start watch token failed, err: %v, filter: %+v", err, filter)
		}

		tailNode := new(watch.ChainNode)
		if err := f.watchDB.Table(f.key.ChainCollection()).Find(map[string]interface{}{}).Fields(common.BKTokenField).
			Sort(common.BKFieldID+":-1").One(context.Background(), tailNode); err != nil {

			if !f.watchDB.IsNotFoundError(err) {
				f.metrics.collectMongoError()
				blog.Errorf("get last watch token from mongo failed, err: %v", err)
				return "", err
			}
			// the tail node is not exist.
			return "", nil
		}
		return tailNode.Token, nil
	}

	return data.Token, nil
}
