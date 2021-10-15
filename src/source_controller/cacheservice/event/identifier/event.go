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

package identifier

import (
	"context"
	"fmt"
	"strings"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"go.mongodb.org/mongo-driver/mongo"
)

type identityOptions struct {
	key         event.Key
	watchFields []string
	watch       stream.LoopInterface
	isMaster    discovery.ServiceManageInterface
	watchDB     *local.Mongo
	ccDB        dal.DB
}

func newIdentity(ctx context.Context, opts identityOptions) error {
	identity := hostIdentity{
		identityOptions: opts,
		metrics:         event.InitialMetrics(opts.key.Collection(), "host_identifier"),
	}

	return identity.Run(ctx)
}

type hostIdentity struct {
	identityOptions
	metrics      *event.EventMetrics
	tokenHandler *identityHandler
}

const batchSize = 500

func (f *hostIdentity) Run(ctx context.Context) error {
	blog.Infof("start host identity events, for key: %s.", f.key.Namespace())

	es := make(map[string]interface{})
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct:     &es,
			Collection:      f.key.Collection(),
			StartAfterToken: nil,
		},
	}
	if f.key.Collection() == common.BKTableNameBaseHost {
		watchOpts.EventStruct = new(metadata.HostMapStr)
	}

	f.tokenHandler = newIdentityTokenHandler(f.key, f.watchDB, f.metrics)

	startAtTime, err := f.tokenHandler.getStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get start watch time for %s failed, err: %v", f.key.Collection(), err)
		return err
	}
	watchOpts.StartAtTime = startAtTime
	watchOpts.WatchFatalErrorCallback = f.tokenHandler.resetWatchToken
	watchOpts.Fields = f.watchFields

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         "host_identity_" + f.key.Collection(),
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
		blog.Errorf("host identity events, but watch batch failed, err: %v", err)
		return err
	}

	return nil
}

func (f *hostIdentity) doBatch(es []*types.Event) (retry bool) {

	if len(es) == 0 {
		return false
	}

	rid := es[0].ID()
	hasError := true

	// collect event related metrics
	start := time.Now()
	defer func() {
		if retry {
			f.metrics.CollectRetryError()
		}
		if hasError {
			return
		}
		f.metrics.CollectCycleDuration(time.Since(start))
	}()

	// convert events to host id events
	identityEvents, err := f.rearrangeEvents(rid, es)
	if err != nil {
		blog.Errorf("host identify event, rearrange events failed, will retry, err: %v, rid: %s", err, rid)
		return true
	}

	// get the lock to get sequences ids.
	// otherwise, we can not guarantee the host,host relation,process's event's id is in the right order/sequences
	// it should be a natural increase order.
	if err = getLock(rid); err != nil {
		blog.Errorf("get host identity lock failed, err: %v, rid: %s", err, rid)
		return true
	}

	// release the lock when the job is done or failed,
	defer releaseLock(rid)

	eventIDs, err := f.watchDB.NextSequences(context.Background(), event.HostIdentityKey.Collection(), len(identityEvents))
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.key.ChainCollection(), err, rid)
		return true
	}

	chainNodes := make([]*watch.ChainNode, 0)
	oids := make([]string, 0)
	cursorMap := make(map[string]struct{})
	for index, e := range identityEvents {
		// collect event's basic metrics
		f.metrics.CollectBasic(e)

		switch e.OperationType {
		case types.Insert, types.Update, types.Replace, types.Delete:
		case types.Invalidate:
			blog.Errorf("host identify event, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
			continue
		default:
			blog.Errorf("host identify event, received unsupported event operation type: %s, doc: %s, rid: %s",
				e.OperationType, e.DocBytes, rid)
			continue
		}

		oids = append(oids, e.ID())
		id := eventIDs[index]
		name := f.key.Name(e.DocBytes)
		cursor, err := genHostIdentifyCursor(f.key.Collection(), e, rid)
		if err != nil {
			blog.Errorf("get %s event cursor failed, name: %s, err: %v, oid: %s, rid: %s", f.key.Collection(), name,
				err, e.ID(), rid)
			return false
		}

		// validate if the cursor is already exists, this is happens when the concurrent operation is very high.
		// which will generate the same operation event with same cluster time, and generate with the same cursor
		// in the end. if this happens, drop this event directly, because we only care this host's identifier is
		// changed or not.
		if _, exists := cursorMap[cursor]; exists {
			// skip this event.
			continue
		}
		cursorMap[cursor] = struct{}{}

		chainNode := &watch.ChainNode{
			ID:          id,
			ClusterTime: e.ClusterTime,
			Oid:         e.Oid,
			// redirect all the event type to update.
			EventType: watch.ConvertOperateType(types.Update),
			Token:     e.Token.Data,
			Cursor:    cursor,
		}

		if instanceID := event.HostIdentityKey.InstanceID(e.DocBytes); instanceID > 0 {
			chainNode.InstanceID = instanceID
		}
		chainNodes = append(chainNodes, chainNode)

	}

	lastEvents := es[len(es)-1]
	lastTokenData := mapstr.MapStr{
		common.BKTokenField:       lastEvents.Token.Data,
		common.BKStartAtTimeField: lastEvents.ClusterTime,
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		if err := f.tokenHandler.setLastWatchToken(context.Background(), lastTokenData); err != nil {
			f.metrics.CollectMongoError()
			return false
		}
		return false
	}

	retry, err = f.doInsertEvents(chainNodes, lastTokenData, rid)
	if err != nil {
		return retry
	}

	blog.Infof("insert host identity event for %s success, oid: %v, rid: %s", f.key.Collection(), oids, rid)
	hasError = false
	return false
}

func (f *hostIdentity) rearrangeEvents(rid string, es []*types.Event) ([]*types.Event, error) {
	switch f.key.Collection() {
	case event.HostKey.Collection():
		return f.rearrangeHostEvents(es, rid), nil
	case event.ModuleHostRelationKey.Collection():
		return f.rearrangeHostRelationEvents(es, rid)
	case event.ProcessKey.Collection():
		return f.rearrangeProcessEvents(es, rid)
	default:
		blog.ErrorJSON("received unsupported host identity event, skip, es: %s, rid :%s", es, rid)
		return es[:0], nil
	}
}

func (f *hostIdentity) doInsertEvents(chainNodes []*watch.ChainNode, lastTokenData map[string]interface{}, rid string) (
	bool, error) {

	count := len(chainNodes)

	if count == 0 {
		return false, nil
	}

	watchDBClient := f.watchDB.GetDBClient()

	session, err := watchDBClient.StartSession()
	if err != nil {
		blog.Errorf("host identity events, but start session failed, coll: %s, err: %v, rid: %s", f.key.Collection(), err, rid)
		return true, err
	}
	defer session.EndSession(context.Background())

	// conflictError record the conflict cursor error
	var conflictError error

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			blog.Errorf("host identity events, but start transaction failed, coll: %s, err: %v, rid: %s",
				f.key.Collection(), err, rid)
			return err
		}

		if err := f.watchDB.Table(event.HostIdentityKey.ChainCollection()).Insert(sc, chainNodes); err != nil {
			blog.ErrorJSON("host identity events, but insert chain nodes for %s failed, nodes: %s, err: %v, rid: %s",
				f.key.Collection(), chainNodes, err, rid)
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())

			if event.IsConflictError(err) {
				conflictError = err
			}
			return err
		}

		lastNode := chainNodes[len(chainNodes)-1]
		lastTokenData[common.BKFieldID] = lastNode.ID
		lastTokenData[common.BKCursorField] = lastNode.Cursor
		lastTokenData[common.BKStartAtTimeField] = lastNode.ClusterTime
		if err := f.tokenHandler.setLastWatchToken(sc, lastTokenData); err != nil {
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())
			return err
		}

		// Use context.Background() to ensure that the commit can complete successfully even if the context passed to
		// mongo.WithSession is changed to have a timeout.
		if err = session.CommitTransaction(context.Background()); err != nil {
			blog.Errorf("host identity events, but commit mongo transaction failed, err: %v", err)
			f.metrics.CollectMongoError()
			return err
		}
		return nil
	})

	if txnErr != nil {
		blog.Errorf("do insert host identity events failed, err: %v, rid: %s", txnErr, rid)

		rid = rid + ":" + chainNodes[0].Oid
		if conflictError != nil && len(chainNodes) >= 1 {
			monitor.Collect(&meta.Alarm{
				RequestID: rid,
				Type:      meta.EventFatalError,
				Detail: fmt.Sprintf("host identifier, but got conflict %s cursor with chain nodes",
					f.key.Collection()),
				Module:    types2.CCModuleCacheService,
				Dimension: map[string]string{"retry_conflict_nodes": "yes"},
			})

			var conflictNode *watch.ChainNode
			// get the conflict cursor
			for idx := range chainNodes {
				if strings.Contains(conflictError.Error(), chainNodes[idx].Cursor) {
					// record conflict node
					conflictNode = chainNodes[idx]
					// remove the conflict cursor
					chainNodes = append(chainNodes[0:idx], chainNodes[idx+1:]...)
					break
				}
			}

			if conflictNode == nil {
				// this should not happen
				// reduce event's one by one, then retry again.
				blog.ErrorJSON("host identity events, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
					f.key.Collection(), chainNodes[0], chainNodes[1:], rid)

				// retry insert events
				return f.doInsertEvents(chainNodes[1:], lastTokenData, rid)
			}

			blog.ErrorJSON("host identity events, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
				f.key.Collection(), conflictNode, chainNodes, rid)

			// retry insert events
			return f.doInsertEvents(chainNodes, lastTokenData, rid)
		}

		// if an error occurred, roll back and re-watch again
		blog.Warnf("do insert host identity events, do retry insert with rid: %s", rid)
		return true, err
	}

	return false, nil
}

const (
	hostIdentityLockKey = common.BKCacheKeyV3Prefix + "host_identity:event_lock"
	hostIdentityLockTTL = 1 * time.Minute
)

func getLock(rid string) error {
	timeout := time.After(hostIdentityLockTTL)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("get host identity: %s lock timeout", hostIdentityLockKey)
		default:
		}

		success, err := redis.Client().SetNX(context.Background(), hostIdentityLockKey, 1, hostIdentityLockTTL).Result()
		if err != nil {
			blog.Errorf("get host identity: %s lock, err: %v, rid: %s", hostIdentityLockKey, err, rid)
			return err
		}

		if !success {
			blog.V(3).Infof("get host identity: %s lock failed, will retry later, rid: %s", hostIdentityLockKey, rid)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// get lock success.
		return nil
	}

}

func releaseLock(rid string) {
	_, err := redis.Client().Del(context.Background(), hostIdentityLockKey).Result()
	if err != nil {
		blog.Errorf("delete host identity redis lock key: %s failed, err: %v, rid: %s", hostIdentityLockKey, err, rid)
		return
	}
	return
}

func genHostIdentifyCursor(coll string, e *types.Event, rid string) (string, error) {
	curType := watch.UnknownType
	switch coll {
	case common.BKTableNameBaseHost:
		curType = watch.Host
	case common.BKTableNameModuleHostConfig:
		curType = watch.ModuleHostRelation
	case common.BKTableNameBaseProcess:
		curType = watch.Process
	default:
		blog.ErrorJSON("unsupported host identity cursor type collection: %s, event: %s, oid: %s", coll, e, rid)
		return "", fmt.Errorf("unsupported host identity cursor type collection: %s", coll)
	}

	hCursor := &watch.Cursor{
		Type:        curType,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		Oper:        e.OperationType,
		// UniqKey:     coll,
	}

	hCursorEncode, err := hCursor.Encode()
	if err != nil {
		blog.ErrorJSON("encode head node cursor failed, cursor: %s, err: %s, rid: %s", hCursor, err, rid)
		return "", err
	}

	return hCursorEncode, nil
}
