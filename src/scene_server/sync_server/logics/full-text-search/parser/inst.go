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
	"errors"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/cache"
)

// objInstParser is the general data parser for object instance
type objInstParser struct {
	*commonParser
}

func newObjInstParser(index string, cs *ClientSet) *objInstParser {
	return &objInstParser{commonParser: newCommonParser(index, cs)}
}

// ParseData parse mongo data to es data
func (p *objInstParser) ParseData(info []mapstr.MapStr, coll string, rid string) (bool, mapstr.MapStr, error) {
	if len(info) == 0 {
		return false, nil, errors.New("data is empty")
	}
	data := info[0]

	objID := GetObjIDByData(coll, data)

	data = cache.EnumIDToName(p.cs.CacheCli, data, objID)

	// get es id
	id, err := convMetaIDToStr(data, metadata.GetInstIDFieldByObjID(objID))
	if err != nil {
		blog.Errorf("get meta id failed, err: %v, data: %+v, obj: %s, rid: %s", err, data, objID, rid)
		return false, nil, err
	}

	_, esDoc, err := p.commonParser.ParseData(info, coll, rid)
	if err != nil {
		return false, nil, err
	}

	esDoc[metadata.IndexPropertyID] = id
	esDoc[metadata.IndexPropertyDataKind] = metadata.DataKindInstance
	esDoc[metadata.IndexPropertyBKObjID] = objID

	if len(info) > 1 {
		// parse quoted instance data
		quotedData := make(map[string]mapstr.MapStr)

		for i := 1; i < len(info); i++ {
			quotedInst := info[i]
			quotedOid, err := parseOid(quotedInst[common.MongoMetaID])
			if err != nil {
				return false, nil, err
			}
			propertyID := util.GetStrByInterface(quotedInst[common.BKPropertyIDField])

			quotedInst = cleanCommonKeywordData(quotedInst, p.index)
			delete(quotedInst, common.BKFieldID)
			delete(quotedInst, common.BKInstIDField)
			delete(quotedInst, common.BKPropertyIDField)

			quotedKeywords, err := parseKeywords(quotedInst)
			if err != nil {
				blog.Errorf("parse quoted inst %+v keywords failed, err: %v, rid: %s", quotedInst, err, rid)
				return false, nil, err
			}

			_, exists := quotedData[propertyID]
			if !exists {
				quotedData[propertyID] = make(mapstr.MapStr)
			}
			quotedData[propertyID][quotedOid] = quotedKeywords
		}

		esDoc[metadata.TablePropertyName] = quotedData
	}

	return false, esDoc, nil
}
