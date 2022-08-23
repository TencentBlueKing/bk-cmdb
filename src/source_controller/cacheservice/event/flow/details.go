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

package flow

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

// getDeleteEventDetailsFunc function type for getting delete events' oid+collection to related detail map
type getDeleteEventDetailsFunc func(es []*types.Event, db dal.DB, metrics *event.EventMetrics) (map[oidCollKey][]byte,
	bool, error)

// getDeleteEventDetails get delete events' oid+collection to related detail map from cmdb
func getDeleteEventDetails(es []*types.Event, db dal.DB, metrics *event.EventMetrics) (map[oidCollKey][]byte, bool,
	error) {

	oidDetailMap := make(map[oidCollKey][]byte)
	if len(es) == 0 {
		return oidDetailMap, false, nil
	}

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

		docs := make([]map[string]interface{}, 0)
		err := db.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
		if err != nil {
			metrics.CollectMongoError()
			blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
				collection, deletedEventOids, err)
			return nil, true, err
		}

		for _, doc := range docs {
			oid := util.GetStrByInterface(doc["oid"])
			byt, err := json.Marshal(doc["detail"])
			if err != nil {
				blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
					collection, oid, err)
				return nil, false, err
			}
			oidDetailMap[oidCollKey{oid: oid, coll: collection}] = byt
		}
	}

	return oidDetailMap, false, nil
}

// getDeleteEventDetails get delete events' oid+collection to related detail map from cmdb
func getHostDeleteEventDetails(es []*types.Event, db dal.DB, metrics *event.EventMetrics) (map[oidCollKey][]byte, bool,
	error) {

	oidDetailMap := make(map[oidCollKey][]byte)
	if len(es) == 0 {
		return oidDetailMap, false, nil
	}

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
		"coll": common.BKTableNameBaseHost,
	}

	docs := make([]event.HostArchive, 0)
	err := db.Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
	if err != nil {
		metrics.CollectMongoError()
		blog.Errorf("get archive deleted doc for collection %s from mongodb failed, oids: %+v, err: %v",
			common.BKTableNameBaseHost, deletedEventOids, err)
		return nil, true, err
	}

	for _, doc := range docs {
		byt, err := json.Marshal(doc.Detail)
		if err != nil {
			blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v",
				common.BKTableNameBaseHost, doc.Oid, err)
			return nil, false, err
		}
		oidDetailMap[oidCollKey{oid: doc.Oid, coll: common.BKTableNameBaseHost}] = byt
	}

	return oidDetailMap, false, nil
}
