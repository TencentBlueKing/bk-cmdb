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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/parser"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"

	"github.com/olivere/elastic/v7"
)

// dataSyncer is the es data syncer
type dataSyncer struct {
	index    string
	parser   parser.Parser
	bulk     *elastic.BulkService
	requests []elastic.BulkableRequest
}

// newDataSyncer new dataSyncer
func newDataSyncer(esCli *elastic.Client, index string) (*dataSyncer, error) {
	_, exists := types.IndexMap[index]
	if !exists {
		return nil, fmt.Errorf("index %s is invalid", index)
	}

	return &dataSyncer{
		index:  index,
		parser: parser.IndexParserMap[index],
		bulk:   esCli.Bulk(),
	}, nil
}

// addUpsertReq add upsert request to es bulk request, returns if the data is valid and needs to be upserted
func (s *dataSyncer) addUpsertReq(kit *rest.Kit, coll, oid string, data []mapstr.MapStr) bool {
	if len(data) == 0 {
		blog.Errorf("upsert data is empty, coll: %s, oid: %s, rid: %s", coll, oid, kit.Rid)
		return false
	}

	skip, doc, err := s.parser.ParseData(kit, data, coll)
	if err != nil {
		blog.Errorf("parse %s data %+v failed, err: %v, rid: %s", coll, data, err, kit.Rid)
		return false
	}

	if skip {
		return false
	}

	id := s.parser.GenEsID(kit.TenantID, oid)
	if coll == common.BKTableNameObjAttDes {
		id = s.parser.GenEsID(kit.TenantID, util.GetStrByInterface(doc[common.MongoMetaID]))
		delete(doc, common.MongoMetaID)
	}

	req := elastic.NewBulkUpdateRequest().
		Index(types.GetIndexName(s.index)).
		RetryOnConflict(10).
		Id(id)

	_, exists := doc[metadata.TablePropertyName]
	if exists {
		// upsert document with nested table fields by script, this will upsert the nested data to the exact value
		req.Script(elastic.NewScriptInline(`ctx._source=params`).Params(doc)).Upsert(doc)
	} else {
		req.DocAsUpsert(true).Doc(doc)
	}

	if _, err = req.Source(); err != nil {
		blog.Errorf("upsert data is invalid, err: %v, id: %s, data: %+v, rid: %s", err, id, data, kit.Rid)
		return false
	}

	s.requests = append(s.requests, req)

	return true
}

// addWatchDeleteReq add watch data delete request to es bulk request, returns if the data needs to be deleted
func (s *dataSyncer) addWatchDeleteReq(kit *rest.Kit, coll, oid string, data mapstr.MapStr) bool {
	if len(data) == 0 {
		blog.Errorf("delete data is empty, coll: %s, oid: %s, rid: %s", coll, oid, kit.Rid)
		return false
	}

	needDel, extraRequest, skip := s.parser.ParseWatchDeleteData(kit, data, coll, oid)
	if skip {
		return false
	}

	if needDel {
		req := elastic.NewBulkDeleteRequest().Index(types.GetIndexName(s.index)).Id(s.parser.GenEsID(kit.TenantID, oid))
		s.requests = append(s.requests, req)
		return true
	}

	s.requests = append(s.requests, extraRequest)
	return true
}

// addDeleteReq add es data delete request to es bulk request, returns if the data needs to be deleted
func (s *dataSyncer) addEsDeleteReq(delEsIDs []string, rid string) bool {
	if len(delEsIDs) == 0 {
		blog.Errorf("es delete ids is empty, rid: %s", rid)
		return false
	}

	for _, id := range delEsIDs {
		req := elastic.NewBulkDeleteRequest().Index(types.GetIndexName(s.index)).Id(id)
		s.requests = append(s.requests, req)
	}

	return true
}

// doBulk do es bulk request
func (s *dataSyncer) doBulk(ctx context.Context, rid string) error {
	defer func() {
		s.requests = make([]elastic.BulkableRequest, 0)
	}()

	return ferrors.EsRespErrHandler(func() (bool, error) {
		if len(s.requests) == 0 {
			return false, nil
		}

		s.bulk.Reset()
		for _, req := range s.requests {
			s.bulk.Add(req)
		}

		resp, err := s.bulk.Do(ctx)
		if err != nil {
			blog.Errorf("do bulk request failed, err: %v, requests: %+v, rid: %s", err, s.requests, rid)
			return false, err
		}

		if resp == nil || !resp.Errors {
			return false, nil
		}

		if len(resp.Items) != len(s.requests) {
			blog.Errorf("bulk response length %d != request length %d, rid: %s", len(resp.Items), len(s.requests), rid)
			return false, errors.New("bulk response length != request length")
		}

		var retry, fatal bool
		retryRequests := make([]elastic.BulkableRequest, 0)
		for i, item := range resp.Items {
			for _, result := range item {
				retry, fatal = ferrors.EsStatusHandler(result.Status)
				if !retry {
					break
				}

				blog.Errorf("do request %+v failed, resp: %+v, rid: %s", s.requests[i], result, rid)
				retryRequests = append(retryRequests, s.requests[i])
				break
			}
		}

		if len(retryRequests) > 0 {
			return fatal, errors.New("do bulk request failed")
		}

		return false, nil
	})
}
