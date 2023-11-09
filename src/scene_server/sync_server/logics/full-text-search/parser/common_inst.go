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

package parser

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/cache"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	"configcenter/src/storage/driver/mongodb"

	"github.com/olivere/elastic/v7"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// commonObjInstParser is the data parser for common object instance, including table instance
type commonObjInstParser struct {
	*objInstParser
}

func newCommonObjInstParser(index string, cs *ClientSet) *commonObjInstParser {
	return &commonObjInstParser{newObjInstParser(index, cs)}
}

// GenEsID generate es id from mongo oid
func (p *commonObjInstParser) GenEsID(coll, oid string) string {
	return fmt.Sprintf("%s:%s", oid, indexIdentifierMap[p.index])
}

// ParseData parse mongo data to es data
func (p *commonObjInstParser) ParseData(info []mapstr.MapStr, coll string, rid string) (bool, mapstr.MapStr, error) {
	if len(info) == 0 {
		return false, nil, errors.New("data is empty")
	}
	data := info[0]

	objID := GetObjIDByData(coll, data)

	// parse table instance separately
	supplierAccount := util.GetStrByInterface(data[common.BkSupplierAccount])
	isQuoted, propertyID, srcObjID := cache.GetQuotedInfoByObjID(p.cs.CacheCli, objID, supplierAccount)
	if isQuoted {
		return p.parseQuotedInst(data, propertyID, srcObjID, rid)
	}

	return p.objInstParser.ParseData(info, coll, rid)
}

// parseQuotedInst parse quoted instance mongo data to es data
func (p *commonObjInstParser) parseQuotedInst(data mapstr.MapStr, propertyID, objID string, rid string) (bool,
	mapstr.MapStr, error) {

	ctx := context.Background()

	instID, err := util.GetInt64ByInterface(data[common.BKInstIDField])
	if err != nil {
		blog.Errorf("[%s] parse quote inst id failed, err: %v, data: %+v, rid: %s", err, data, rid)
		return false, nil, errors.New("quote inst id is invalid")
	}

	// Note: instID == 0 表明表格实例没有与模型实例表进行关联，无需处理
	if instID == 0 {
		return true, nil, nil
	}

	oid, err := parseOid(data[common.MongoMetaID])
	if err != nil {
		return false, nil, err
	}

	account, err := convMetaIDToStr(data, common.BKOwnerIDField)
	if err != nil {
		blog.Errorf("[%s] parse supplier account failed, err: %v, data: %+v, rid: %s", err, data, rid)
		return false, nil, errors.New("supplier account is invalid")
	}

	index := getEsIndexByObjID(objID)

	document, keywords, err := p.analysisTableDocument(propertyID, oid, data)
	if err != nil {
		blog.Errorf("analysis table document failed, err: %v", err)
		return false, nil, err
	}

	// 直接更新 es文档
	succeed, err := p.updateTablePropertyEsDoc(index, strconv.FormatInt(instID, 10), propertyID, oid, keywords)
	if err != nil {
		blog.Errorf("update table property es doc failed, err: %v", err)
		return false, nil, err
	}

	if succeed {
		return true, nil, nil
	}

	// 更新败降级处理，查询实例数据，如果es文档不存在，直接创建es文档
	id, err := p.getEsIDByMongoID(objID, account, instID, rid)
	if err != nil {
		return false, nil, err
	}

	err = ferrors.EsRespErrHandler(func() (bool, error) {
		resp, err := p.cs.EsCli.Update().Index(index).DocAsUpsert(true).RetryOnConflict(10).
			Doc(document).Id(id).Do(ctx)
		if err != nil {
			blog.Errorf("upsert parent inst failed, err: %v, id: %s, doc: %+v, rid: %s", err, id, document, rid)
			return false, err
		}

		retry, fatal := ferrors.EsStatusHandler(resp.Status)
		if !retry {
			return false, nil
		}

		return fatal, errors.New("upsert parent inst failed")
	})
	return true, nil, nil
}

// updateTablePropertyEsDoc update table property es doc.
func (p *commonObjInstParser) updateTablePropertyEsDoc(index, instIDStr, propID, oid string, keywords []string) (bool,
	error) {

	keywordStr, err := json.MarshalToString(keywords)
	if err != nil {
		return false, err
	}

	var succeed bool
	err = ferrors.EsRespErrHandler(func() (bool, error) {
		resp, err := p.cs.EsCli.UpdateByQuery(index).
			ProceedOnVersionConflict().
			Query(elastic.NewMatchQuery(metadata.IndexPropertyID, instIDStr)).
			Script(elastic.NewScriptInline(fmt.Sprintf(updateTableScript, propID, propID, propID, oid,
				keywordStr))).
			Do(context.Background())
		if err != nil {
			blog.Errorf("update table property failed, err: %v, inst id: %s, property id: %s", err, instIDStr, propID)
			return false, err
		}

		for _, failure := range resp.Failures {
			retry, fatal := ferrors.EsStatusHandler(failure.Status)
			if !retry {
				break
			}

			return fatal, errors.New("update table property failed")
		}

		succeed = resp.Total == 1
		return false, nil
	})

	return succeed, err
}

// getEsIDByMongoID get the es id by mongo document id.
// 如果mongo的实例数据不存在，说明是脏数据，直接返回错误。
func (p *commonObjInstParser) getEsIDByMongoID(objID, supplierAccount string, id int64, rid string) (string, error) {
	coll := common.GetInstTableName(objID, supplierAccount)
	filter := mapstr.MapStr{common.GetInstIDField(objID): id}

	doc := make(mapstr.MapStr)
	ferrors.FatalErrHandler(200, 100, func() error {
		err := mongodb.Client().Table(coll).Find(filter).Fields(common.MongoMetaID).One(context.Background(), &doc)
		if err != nil {
			blog.Errorf("get mongo _id failed, obj: %s, id: %d, err: %v, rid: %s", objID, id, err, rid)
			return err
		}
		return nil
	})

	documentID, ok := doc[common.MongoMetaID].(primitive.ObjectID)
	if !ok {
		return "", errors.New("missing document metadata id")
	}

	return p.GenEsID(coll, documentID.Hex()), nil
}

// analysisTableDocument analysis the table property document.
func (p *commonObjInstParser) analysisTableDocument(propertyID, oid string, originDoc mapstr.MapStr) (
	mapstr.MapStr, []string, error) {

	originDoc = cleanCommonKeywordData(originDoc, p.index)

	delete(originDoc, common.BKFieldID)
	delete(originDoc, common.BKInstIDField)

	jsonDoc, err := json.MarshalToString(originDoc)
	if err != nil {
		return nil, nil, err
	}

	keywords := analysisJSONKeywords(gjson.Parse(jsonDoc))
	document := mapstr.MapStr{
		metadata.TablePropertyName: mapstr.MapStr{
			propertyID: mapstr.MapStr{
				oid: keywords,
			},
		},
	}
	return document, keywords, nil
}

// ParseWatchDeleteData parse delete data from mongodb watch
func (p *commonObjInstParser) ParseWatchDeleteData(collOidMap map[string][]string, rid string) ([]string,
	[]elastic.BulkableRequest, bool) {

	delArchives := getDelArchive(collOidMap, rid)

	needDelIDs := make([]string, 0)

	for _, archive := range delArchives {
		objID := util.GetStrByInterface(archive.Detail[common.BKObjIDField])
		esID := p.GenEsID(archive.Coll, archive.Oid)

		supplierAccount := util.GetStrByInterface(archive.Detail[common.BkSupplierAccount])
		isQuoted, propID, objID := cache.GetQuotedInfoByObjID(p.cs.CacheCli, objID, supplierAccount)
		if !isQuoted {
			needDelIDs = append(needDelIDs, esID)
			continue
		}

		err := p.deleteTablePropertyEsDoc(getEsIndexByObjID(objID), propID, archive.Oid)
		if err != nil {
			blog.Errorf("delete table property es document failed, err: %v, rid: %s", err, rid)
			continue
		}
	}

	return needDelIDs, nil, true
}

// deleteTablePropertyEsDoc delete table property instance from es.
func (p *commonObjInstParser) deleteTablePropertyEsDoc(index, propertyID, oid string) error {
	return ferrors.EsRespErrHandler(func() (bool, error) {
		resp, err := p.cs.EsCli.UpdateByQuery(index).
			ProceedOnVersionConflict().
			Query(elastic.NewExistsQuery(fmt.Sprintf(deleteTableQueryScript, propertyID, oid))).
			Script(elastic.NewScriptInline(fmt.Sprintf(deleteTableScript, propertyID, oid, propertyID, propertyID))).
			Do(context.Background())
		if err != nil {
			blog.Errorf("delete table inst failed, err: %v", err)
			return false, err
		}

		for _, failure := range resp.Failures {
			retry, fatal := ferrors.EsStatusHandler(failure.Status)
			if !retry {
				break
			}

			return fatal, errors.New("delete table inst failed")
		}
		return false, nil
	})
}
