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

package fulltextsearch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	ftypes "configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	"configcenter/src/scene_server/sync_server/logics/token_handler"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
)

// incrSyncer is the full-text search data incremental syncer
type incrSyncer struct {
	index        string
	metrics      *event.EventMetrics
	tokenHandler *tokenhandler.TokenHandler
	dataSyncer   *dataSyncer
}

const (
	// watchKeyPrefix is the prefix of full-text search watch key
	watchKeyPrefix = "full_text_search:"
	// watchBatchSize is the batch size of full-text search watch
	watchBatchSize = 1000
)

// incrementalSync do full-text search incremental sync
func (f fullTextSearch) incrementalSync(watcher stream.LoopInterface) error {
	if !f.enableSync {
		return errors.New("full text search sync is disabled")
	}

	ctx := context.Background()
	rid := util.GenerateRID()

	blog.Infof("start full-text search incremental sync, rid: %s", rid)

	for _, index := range ftypes.AllIndexNames {
		key := genWatchKey(index)

		dataSyncer, err := newDataSyncer(f.esCli.Client, index)
		if err != nil {
			return err
		}

		syncer := incrSyncer{
			index:      index,
			metrics:    event.InitialMetrics(ftypes.IndexCollMap[index], "watch"),
			dataSyncer: dataSyncer,
		}

		// create token handler
		tokenHandler, err := tokenhandler.New(key, mongodb.Client("watch"), syncer.metrics)
		if err != nil {
			return err
		}
		syncer.tokenHandler = tokenHandler

		exists, startAtTime, err := syncer.tokenHandler.GetStartWatchTime(ctx)
		if err != nil {
			blog.Errorf("get start watch time for %s failed, err: %v, rid: %s", key, err, rid)
			return err
		}

		// init watch token data and all es info if it's the first time to sync this index(token not exists)
		if !exists {
			if err = syncer.tokenHandler.InitWatchToken(ctx); err != nil {
				return err
			}

			if err = f.syncDataByIndex(ctx, index, rid); err != nil {
				return err
			}
		}

		watchOpts := &types.WatchOptions{
			Options: types.Options{
				EventStruct:             new(mapstr.MapStr),
				StartAtTime:             &startAtTime,
				WatchFatalErrorCallback: syncer.tokenHandler.ResetWatchToken,
			},
		}

		addCollWatchOpt(index, watchOpts)

		opts := &types.LoopBatchOptions{
			LoopOptions: types.LoopOptions{
				Name:         key,
				WatchOpt:     watchOpts,
				TokenHandler: syncer.tokenHandler,
			},
			EventHandler: &types.BatchHandler{
				DoBatch: syncer.doBatch,
			},
			BatchSize: watchBatchSize,
		}

		if err = watcher.WithBatch(opts); err != nil {
			blog.Errorf("watch batch failed, err: %v, opt: %+v, rid: %s", err, opts, rid)
			return err
		}
	}

	return nil
}

// addCollWatchOpt add collection watch options
func addCollWatchOpt(index string, watchOpts *types.WatchOptions) {
	switch index {
	case metadata.IndexNameModel:
		// watch model and attribute table
		watchOpts.Options.CollectionFilter = mapstr.MapStr{
			common.BKDBIN: []string{common.BKTableNameObjDes, common.BKTableNameObjAttDes},
		}
	case metadata.IndexNameObjectInstance:
		// watch all tables with the prefix of instance table
		watchOpts.Options.CollectionFilter = mapstr.MapStr{
			common.BKDBLIKE: event.ObjInstTablePrefixRegex,
		}
	default:
		watchOpts.Options.Collection = ftypes.IndexCollMap[index]
	}
}

// genWatchKey generate full-text search watch key
func genWatchKey(index string) string {
	return watchKeyPrefix + index
}

// doBatch handle one batch of full-text search watch events
func (s incrSyncer) doBatch(es []*types.Event) (retry bool) {
	eventLen := len(es)
	if eventLen == 0 {
		return false
	}

	rid := es[0].ID()
	hasError := true

	// collect event related metrics
	start := time.Now()
	defer func() {
		if retry {
			s.metrics.CollectRetryError()
		}
		if hasError {
			return
		}
		s.metrics.CollectCycleDuration(time.Since(start))
	}()

	// aggregate events
	es, keys := s.aggregateEvents(es)

	// parse events and do es operation
	hasChange := false
	collOidMap := make(map[string][]string)

	for _, e := range es {
		switch e.OperationType {
		case types.Insert, types.Update, types.Replace:
			docPtr, ok := e.Document.(*mapstr.MapStr)
			if !ok || docPtr == nil || *docPtr == nil {
				blog.Errorf("[%s] document is empty, event: %+v, rid: %s", s.index, e, rid)
				continue
			}

			if s.dataSyncer.addUpsertReq(e.Collection, e.Oid, []mapstr.MapStr{*docPtr}, rid) {
				hasChange = true
			}
		case types.Delete:
			collOidMap[e.Collection] = append(collOidMap[e.Collection], e.Oid)
		// since following event cannot be parsed, skip them and do not retry
		case types.Invalidate:
			blog.Errorf("[%s] received invalid event operation type, doc: %s, rid: %s", s.index, e.DocBytes, rid)
		case types.Drop:
			blog.Errorf("[%s] received drop table event operation type, doc: %s, rid: %s", s.index, e.DocBytes, rid)
		default:
			blog.Errorf("[%s] received unsupported event operation type: %s, doc: %s, rid: %s", s.index,
				e.OperationType, e.DocBytes, rid)
		}
	}

	if s.dataSyncer.addWatchDeleteReq(collOidMap, rid) {
		hasChange = true
	}

	if !hasChange {
		blog.Infof("[%s] all events are invalid or do not need to sync, keys: %+v, rid: %s", s.index, keys, rid)
		return false
	}

	if err := s.dataSyncer.doBulk(context.Background(), rid); err != nil {
		blog.Infof("[%s] do es bulk request failed, err: %v, rid: %s", s.index, err, rid)
		return false
	}

	blog.Infof("[%s] full-text search incremental sync success, keys: %+v, rid: %s", s.index, keys, rid)
	hasError = false
	return false
}

// aggregateEvents deduplicate events by collection and oid
func (s incrSyncer) aggregateEvents(es []*types.Event) ([]*types.Event, []string) {
	eventMap := make(map[string]*types.Event)
	keys := make([]string, 0)

	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]

		// collect event's basic metrics
		s.metrics.CollectBasic(e)

		key := genEventUniqueKey(e.Collection, e.Oid)
		keys = append(keys, key)

		lastEvent, exists := eventMap[key]
		if !exists {
			eventMap[key] = e
			continue
		}

		switch e.OperationType {
		case types.Insert:
			if lastEvent.OperationType == types.Delete {
				delete(eventMap, key)
				continue
			}
			eventMap[key].OperationType = types.Insert
		}
	}

	events := make([]*types.Event, 0)
	for _, e := range eventMap {
		events = append(events, e)
	}

	return events, keys
}

// genEventUniqueKey generate event unique key
func genEventUniqueKey(coll, oid string) string {
	return fmt.Sprintf("%s.%s", coll, oid)
}
