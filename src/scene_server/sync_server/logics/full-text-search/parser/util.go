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
	"fmt"
	"regexp"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	"configcenter/src/storage/driver/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tidwall/gjson"
)

// convMetaIDToStr convert meta id(objID/hostID/setID/moduleID/instanceID/bizID...) to string.
func convMetaIDToStr(data mapstr.MapStr, idField string) (string, error) {
	id, exists := data[idField]
	if !exists || id == nil || id == "" {
		return "", fmt.Errorf("document id %+v is invalid", data[idField])
	}
	return fmt.Sprintf("%v", id), nil
}

// cleanCommonKeywordData cleans common fields that do not need to be saved in es
func cleanCommonKeywordData(document mapstr.MapStr, index string) mapstr.MapStr {
	if len(document) == 0 {
		return make(mapstr.MapStr)
	}

	for _, field := range baseCleanFields {
		delete(document, field)
	}

	for _, field := range indexKeywordCleanFieldsMap[index] {
		delete(document, field)
	}

	return document
}

// parseKeywords parse es keywords by index
func parseKeywords(data mapstr.MapStr) ([]string, error) {
	jsonDoc, err := json.MarshalToString(data)
	if err != nil {
		return nil, err
	}

	keywords := analysisJSONKeywords(gjson.Parse(jsonDoc))
	return compressKeywords(keywords), nil
}

// analysisJSONKeywords analysis the given json style document,
// and extract all the keywords as elastic document content.
func analysisJSONKeywords(result gjson.Result) []string {
	keywords := make([]string, 0)
	if !result.IsObject() && !result.IsArray() {
		keyword := result.String()
		if len(keyword) != 0 {
			keywords = append(keywords, keyword)
		}
		return keywords
	}

	result.ForEach(func(key, value gjson.Result) bool {
		keywords = append(keywords, analysisJSONKeywords(value)...)
		return true
	})

	return keywords
}

// compressKeywords compress the keywords, unique the keywords array.
func compressKeywords(keywords []string) []string {
	compressedKeywords := make([]string, 0)
	// keywordsMap control repeated or screened keywords.
	keywordsMap := make(map[string]struct{})
	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}
		if _, exist := keywordsMap[keyword]; exist {
			continue
		}
		compressedKeywords = append(compressedKeywords, keyword)
		keywordsMap[keyword] = struct{}{}
	}

	return compressedKeywords
}

// GetObjIDByData get object id by collection & instance data
func GetObjIDByData(coll string, data mapstr.MapStr) string {
	switch coll {
	case common.BKTableNameBaseBizSet:
		return common.BKInnerObjIDBizSet
	case common.BKTableNameBaseApp:
		return common.BKInnerObjIDApp
	case common.BKTableNameBaseSet:
		return common.BKInnerObjIDSet
	case common.BKTableNameBaseModule:
		return common.BKInnerObjIDModule
	case common.BKTableNameBaseHost:
		return common.BKInnerObjIDHost
	default:
		if !common.IsObjectInstShardingTable(coll) {
			return ""
		}

		if data == nil {
			return ""
		}

		objID := util.GetStrByInterface(data[common.BKObjIDField])
		if objID != "" {
			return objID
		}

		// parse obj id from table name, NOTE: this is only a compatible logics
		regex := regexp.MustCompile(`cc_ObjectBase_(.*)_pub_(.*)`)
		if regex.MatchString(coll) {
			matches := regex.FindStringSubmatch(coll)
			return matches[2]
		}

		return ""
	}
}

var objEsIndexMap = map[string]string{
	common.BKInnerObjIDBizSet: metadata.IndexNameBizSet,
	common.BKInnerObjIDApp:    metadata.IndexNameBiz,
	common.BKInnerObjIDSet:    metadata.IndexNameSet,
	common.BKInnerObjIDModule: metadata.IndexNameModule,
	common.BKInnerObjIDHost:   metadata.IndexNameHost,
}

// getEsIndexByObjID get the es index by object id.
func getEsIndexByObjID(objID string) string {
	index, exists := objEsIndexMap[objID]
	if exists {
		return types.GetIndexName(index)
	}

	return types.GetIndexName(metadata.IndexNameObjectInstance)
}

type delArchive struct {
	Oid    string        `bson:"oid"`
	Coll   string        `bson:"coll"`
	Detail mapstr.MapStr `bson:"detail"`
}

// getDelArchive get deleted data by collOidMap, returns es id to deleted mongo data map
func getDelArchive(collOidMap map[string][]string, rid string) []delArchive {
	orCond := make([]mapstr.MapStr, 0)
	for coll, oids := range collOidMap {
		orCond = append(orCond, mapstr.MapStr{
			"coll": coll,
			"oid":  oids,
		})
	}

	filter := mapstr.MapStr{common.BKDBOR: orCond}

	docs := make([]delArchive, 0)

	ferrors.FatalErrHandler(200, 100, func() error {
		err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).All(context.Background(), &docs)
		if err != nil {
			blog.Errorf("get del archive failed, filter: %+v, err: %v, rid: %s", filter, err, rid)
			return err
		}
		return nil
	})

	return docs
}

func parseOid(oid interface{}) (string, error) {
	switch t := oid.(type) {
	case primitive.ObjectID:
		return t.Hex(), nil
	case string:
		return t, nil
	default:
		return "", fmt.Errorf("oid %+v is invalid", t)
	}
}
