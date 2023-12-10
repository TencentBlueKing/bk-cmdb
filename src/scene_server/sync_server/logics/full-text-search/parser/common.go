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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"

	"github.com/olivere/elastic/v7"
)

// commonParser is the data parser for common object instance, including table instance
type commonParser struct {
	index string
	cs    *ClientSet
}

func newCommonParser(index string, cs *ClientSet) *commonParser {
	return &commonParser{index: index, cs: cs}
}

// GenEsID generate es id from mongo oid
func (p *commonParser) GenEsID(coll, oid string) string {
	return fmt.Sprintf("%s:%s", oid, indexIdentifierMap[p.index])
}

// ParseData parse mongo data to es data
func (p *commonParser) ParseData(info []mapstr.MapStr, coll string, rid string) (bool, mapstr.MapStr, error) {
	if len(info) == 0 {
		return false, nil, errors.New("data is empty")
	}
	data := info[0]

	// generate es doc
	esDoc := mapstr.MapStr{
		metadata.IndexPropertyBKObjID:           data[common.BKObjIDField],
		metadata.IndexPropertyBKSupplierAccount: data[common.BKOwnerIDField],
		metadata.IndexPropertyBKBizID:           data[common.BKAppIDField],
	}

	for _, field := range types.IndexExtraFieldsMap[p.index] {
		esDoc[field] = data[extraEsFieldMap[field]]
	}

	for _, field := range types.IndexExcludeFieldsMap[p.index] {
		delete(esDoc, field)
	}

	// parse es keywords
	data = cleanCommonKeywordData(data, p.index)
	keywords, err := parseKeywords(data)
	if err != nil {
		blog.Errorf("parse keywords failed, err: %v, data: %+v, index: %s, rid: %s", err, data, p.index, rid)
		return false, nil, err
	}

	esDoc[metadata.IndexPropertyKeywords] = keywords

	return false, esDoc, nil
}

// ParseWatchDeleteData parse delete data from mongodb watch
func (p *commonParser) ParseWatchDeleteData(collOidMap map[string][]string, rid string) ([]string,
	[]elastic.BulkableRequest, bool) {
	return nil, nil, true
}
