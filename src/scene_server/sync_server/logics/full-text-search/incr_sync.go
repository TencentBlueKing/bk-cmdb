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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	ftypes "configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	"configcenter/src/scene_server/sync_server/logics/token_handler"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/task"
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

// addIncrementalSyncTask add full-text search incremental sync task
func (f *fullTextSearch) addIncrementalSyncTask() error {
	if !f.enableSync {
		return errors.New("full text search sync is disabled")
	}

	ctx := context.Background()
	rid := util.GenerateRID()

	if err := tokenhandler.InitWatchTokenTable(ctx, mongodb.Dal("watch")); err != nil {
		blog.Errorf("init full-text search watch token table failed, err: %v, rid: %s", err, rid)
		return err
	}

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
		tokenHandler, err := tokenhandler.New(key, syncer.metrics)
		if err != nil {
			return err
		}
		syncer.tokenHandler = tokenHandler

		opts := &types.LoopBatchTaskOptions{
			WatchTaskOptions: &types.WatchTaskOptions{
				Name: key,
				CollOpts: &types.WatchCollOptions{
					CollectionOptions: types.CollectionOptions{
						CollectionFilter: getWatchCollFilter(index),
						EventStruct:      new(mapstr.MapStr),
					},
				},
				TokenHandler: tokenHandler,
			},
			EventHandler: &types.TaskBatchHandler{
				DoBatch: syncer.doBatch,
			},
			BatchSize: watchBatchSize,
		}

		flowTask, err := task.NewLoopBatchTask(opts)
		if err != nil {
			blog.Errorf("generate full-text search %s watch task failed, err: %v", key, err)
			return err
		}

		f.tasks = append(f.tasks, flowTask)
	}

	return nil
}

// getWatchCollFilter get watch collection filter by index
func getWatchCollFilter(index string) *types.CollectionFilter {
	switch index {
	case metadata.IndexNameModel:
		// watch model and attribute table
		return &types.CollectionFilter{
			Regex: fmt.Sprintf("_%s$|_%s$", common.BKTableNameObjDes, common.BKTableNameObjAttDes),
		}
	case metadata.IndexNameObjectInstance:
		return &types.CollectionFilter{
			Regex: fmt.Sprintf("_%s", common.BKObjectInstShardingTablePrefix),
		}
	default:
		return &types.CollectionFilter{
			Regex: fmt.Sprintf("_%s$", ftypes.IndexCollMap[index]),
		}
	}
}

// genWatchKey generate full-text search watch key
func genWatchKey(index string) string {
	return watchKeyPrefix + index
}

// doBatch handle one batch of full-text search watch events
func (s incrSyncer) doBatch(dbInfo *types.DBInfo, es []*types.Event) (retry bool) {
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

	for _, e := range es {
		switch e.OperationType {
		case types.Insert, types.Update, types.Replace, types.Delete:
			docPtr, ok := e.Document.(*mapstr.MapStr)
			if !ok || docPtr == nil || *docPtr == nil {
				blog.Errorf("[%s] document is empty, event: %+v, rid: %s", s.index, e, rid)
				continue
			}

			kit := rest.NewKit().WithRid(rid).WithTenant(e.TenantID)

			if e.OperationType == types.Delete {
				if s.dataSyncer.addWatchDeleteReq(kit, e.ParsedColl, e.Oid, *docPtr) {
					hasChange = true
				}
				continue
			}

			if s.dataSyncer.addUpsertReq(kit, e.ParsedColl, e.Oid, []mapstr.MapStr{*docPtr}) {
				hasChange = true
			}
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
