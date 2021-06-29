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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/stream/types"
)

func newHostFlow(ctx context.Context, opts flowOptions) error {
	flow := HostFlow{
		Flow{
			flowOptions: opts,
			metrics:     event.InitialMetrics(opts.key.Collection(), "watch"),
		},
	}

	return flow.RunFlow(ctx)
}

type HostFlow struct {
	Flow
}

func (f *HostFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	opts, err := f.generateLoopBatchOptions(ctx)
	if err != nil {
		return err
	}

	opts.WatchOpt.EventStruct = &metadata.HostMapStr{}
	opts.EventHandler.DoBatch = f.doBatchWrapper(f.doBatch)

	if err := f.watch.WithBatch(opts); err != nil {
		blog.Errorf("run flow, but watch batch failed, err: %v", err)
		return err
	}
	return nil
}

func (f *HostFlow) doBatch(es []*types.Event, rid string) (bool, error) {
	oidDetailMap, retry, err := f.getDeleteHostEventDetails(es)
	if err != nil {
		blog.Errorf("get deleted event details failed, err: %v, rid: %s", err, rid)
		return retry, err
	}

	return f.batchHandleEvents(es, oidDetailMap, rid)
}

// getDeleteHostEventDetails get delete host events' oid and related detail map from cmdb
func (f *HostFlow) getDeleteHostEventDetails(es []*types.Event) (map[string][]byte, bool, error) {
	oidDetailMap := make(map[string][]byte)

	deletedEventOidMap := make(map[string][]string, 0)
	for _, e := range es {
		if e.OperationType == types.Delete {
			deletedEventOidMap[e.Collection] = append(deletedEventOidMap[e.Collection], e.Oid)
		}
	}

	if len(deletedEventOidMap) == 0 {
		return oidDetailMap, false, nil
	}

	for collection, deletedEventOids := range deletedEventOidMap {
		filter := map[string]interface{}{
			"oid":  map[string]interface{}{common.BKDBIN: deletedEventOids},
			"coll": collection,
		}

		docs := make([]event.HostArchive, 0)
		err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
		if err != nil {
			f.metrics.CollectMongoError()
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
				collection, deletedEventOids, err)
			return nil, true, err
		}

		for _, doc := range docs {
			byt, err := json.Marshal(doc.Detail)
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					collection, doc.Oid, err)
				return nil, false, err
			}
			oidDetailMap[doc.Oid] = byt
		}
	}

	return oidDetailMap, false, nil
}
