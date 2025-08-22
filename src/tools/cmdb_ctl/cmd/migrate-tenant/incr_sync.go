/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package migratetenant

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// copyIncrSyncData copy incremental sync data from old db to new db
func (s *migrateTenantService) copyIncrSyncData(kit *rest.Kit, startFrom uint32) error {
	// get full sync start time for incremental sync if not specified
	if startFrom == 0 {
		startTimeCond := mapstr.MapStr{"_id": "full_sync_start_time"}
		startTime := make(map[string]uint32)
		err := s.sysDB.Table(common.BKTableNameSystem).Find(startTimeCond).Fields("start_from").One(kit.Ctx, &startTime)
		if err != nil {
			return fmt.Errorf("get full sync start time failed, err: %v", err)
		}
		startFrom = startTime["start_from"]
	}

	ctx, cancel := context.WithCancel(kit.Ctx)
	kit.Ctx = ctx

	// refresh watch info before exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		printInfo("receive signal %v, begin to shutdown", sig)
		if err := s.copyWatchInfo(kit); err != nil {
			fmt.Errorf("copy watch info failed, err: %v", err)
		}
		cancel()
	}()

	// watch old db events for incremental sync
	if err := s.watchOldDB(kit, startFrom); err != nil {
		return err
	}
	return nil
}

// storeLastEventTime store last event cluster time
func (s *migrateTenantService) storeLastEventTime(kit *rest.Kit, eventClusterTime uint32) error {
	timeCond := mapstr.MapStr{
		"_id": "full_sync_start_time",
	}
	timeData := mapstr.MapStr{
		"start_from": eventClusterTime,
	}
	if err := s.sysDB.Table(common.BKTableNameSystem).Upsert(kit.Ctx, timeCond, timeData); err != nil {
		return fmt.Errorf("store last event cluster time %d failed, err: %v", eventClusterTime, err)
	}
	return nil
}

// watchOldDB watch old db for incremental sync
func (s *migrateTenantService) watchOldDB(kit *rest.Kit, startFrom uint32) error {
	stream, err := s.streamWatchOldDB(kit, startFrom)
	if err != nil {
		return err
	}

	eventTime := startFrom
	eventCnt := 0
	for stream.Next(kit.Ctx) {
		// still have events, try cancel steam here.
		select {
		case <-kit.Ctx.Done():
			_ = stream.Close(context.Background())
			return nil
		default:
		}

		coll := stream.Current.Lookup("ns", "coll").StringValue()
		if coll == "cc_idgenerator" {
			id := strings.TrimPrefix(stream.Current.Lookup("fullDocument", "_id").StringValue(), "cc_")
			if strings.HasPrefix(id, "id_rule:incr_id:") {
				id = "id_rule:incr_id:default:" + strings.TrimPrefix(id, "id_rule:incr_id:")
			}
			cond := mapstr.MapStr{"_id": id}
			data := mapstr.MapStr{"SequenceID": stream.Current.Lookup("fullDocument", "SequenceID").AsInt64()}
			if err = s.sysDB.Table(strings.TrimPrefix(coll, "cc_")).Upsert(kit.Ctx, cond, data); err != nil {
				return fmt.Errorf("upsert table %s oid %s data(%+v) failed, err: %v", coll, id, data, err)
			}
			continue
		}

		event := new(types.RawEvent)
		if err = stream.Decode(event); err != nil {
			printInfo("decode event %s failed, err: %v", stream.Current, err)
			continue
		}
		eventTime = event.ClusterTime.T

		if err = s.handleIncrSyncEvent(kit, event); err != nil {
			if lastEventErr := s.storeLastEventTime(kit, eventTime); lastEventErr != nil {
				printInfo(lastEventErr.Error())
			}
			return fmt.Errorf("handle incr sync event failed, err: %v", err)
		}

		// store last event time and reset event counter if handled events reach max page size
		eventCnt++
		if eventCnt > common.BKMaxPageSize {
			eventCnt = 0
			if err = s.storeLastEventTime(kit, eventTime); err != nil {
				printInfo(err.Error())
			}
		}
	}

	return nil
}

func (s *migrateTenantService) streamWatchOldDB(kit *rest.Kit, startFrom uint32) (*mongo.ChangeStream, error) {
	oldMongo, ok := s.oldDB.(*local.OldMongo)
	if !ok {
		return nil, fmt.Errorf("old db type(%T) is invalid", s.oldDB)
	}
	oldMongoCli := oldMongo.GetMongoClient()

	major := options.UpdateLookup
	preImage := options.WhenAvailable
	batchSize := int32(2000)
	defaultMaxAwaitTime := time.Second
	streamOptions := &options.ChangeStreamOptions{
		FullDocument:             &major,
		FullDocumentBeforeChange: &preImage,
		BatchSize:                &batchSize,
		MaxAwaitTime:             &defaultMaxAwaitTime,
		StartAtOperationTime:     &primitive.Timestamp{T: startFrom},
	}

	pipeline := []bson.D{{{Key: common.BKDBMatch, Value: bson.D{{Key: "ns.coll", Value: bson.M{
		common.BKDBNIN: []string{"cc_ChartConfig", "cc_ChartData", "cc_ChartPosition", "cc_CloudAccount",
			"cc_CloudSyncHistory", "cc_CloudSyncTak", "cc_KubeDelArchive", "cc_DelArchive", "cc_History",
			"cc_NetcollectDevice", "cc_NetcollectProperty", "cc_Subscription", "cc_System"}}}}}}}
	stream, err := oldMongoCli.Client().
		Database(oldMongoCli.DBName()).
		Watch(kit.Ctx, pipeline, streamOptions)
	if err != nil {
		return nil, fmt.Errorf("watch old db failed, err: %v", err)
	}
	return stream, nil
}

// handleIncrSyncEvent handle incremental sync event, copy data from old db to new db
func (s *migrateTenantService) handleIncrSyncEvent(kit *rest.Kit, event *types.RawEvent) error {
	coll := strings.TrimPrefix(event.Namespace.Collection, "cc_")
	db := s.newDB
	if common.IsPlatformTable(coll) {
		db = s.sysDB
	}

	// handle object instance & instance association data
	for _, prefix := range []string{common.BKObjectInstShardingTablePrefix,
		common.BKObjectInstAsstShardingTablePrefix} {

		if strings.HasPrefix(coll, prefix) {
			uuid, ok := s.objUUIDMap[strings.TrimPrefix(coll, prefix+"0_pub_")]
			if !ok {
				return nil
			}
			coll = common.GetObjInstTableName(uuid)
			s.tableHandlers[coll] = s.removeSupplierAccount
			break
		}
	}

	switch event.OperationType {
	case types.Insert, types.Replace, types.Update:
		handler, exists := s.tableHandlers[coll]
		if !exists {
			return nil
		}

		// convert event detail to new version and upsert to new db
		doc := make(mapstr.MapStr)
		if err := bson.Unmarshal(event.FullDoc, &doc); err != nil {
			printInfo("decode event data(%s) failed, err: %v", event.FullDoc, err)
			return nil
		}

		data, err := handler(doc)
		if err != nil {
			printInfo("handle table %s data(%+v) failed, err: %v", coll, data, err)
			return nil
		}

		cond := mapstr.MapStr{"_id": event.DocumentKey.ID}
		if err = db.Table(coll).Upsert(kit.Ctx, cond, data); err != nil {
			return fmt.Errorf("upsert table %s oid %s data(%+v) failed, err: %v", coll, event.DocumentKey.ID, data, err)
		}
	case types.Delete:
		cond := mapstr.MapStr{"_id": event.DocumentKey.ID}
		err := db.Table(coll).Delete(kit.Ctx, cond)
		if err != nil {
			return fmt.Errorf("delete table %s data %s failed, err: %v", coll, event.DocumentKey.ID, err)
		}
	}

	printInfo("handle event, coll: %s, operation: %s, oid: %s, timestamp: %d\n", event.Namespace.Collection,
		event.OperationType, event.DocumentKey.ID, event.ClusterTime.T)
	return nil
}
