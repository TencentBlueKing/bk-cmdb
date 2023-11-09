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

	"github.com/olivere/elastic/v7"
)

// modelParser is the general data parser for model/attribute
type modelParser struct {
	*commonParser
}

func newModelParser(index string, cs *ClientSet) *modelParser {
	return &modelParser{newCommonParser(index, cs)}
}

// ParseData parse mongo data to es data
func (p *modelParser) ParseData(info []mapstr.MapStr, coll string, rid string) (bool, mapstr.MapStr, error) {
	if len(info) == 0 {
		return false, nil, errors.New("data is empty")
	}
	data := info[0]

	objID, ok := data[common.BKObjIDField].(string)
	if !ok {
		blog.Errorf("[%s] object id is invalid, data: %+v, rid: %s", p.index, data, rid)
		return false, nil, errors.New("object id is invalid")
	}

	// skip table object
	supplierAccount := util.GetStrByInterface(data[common.BkSupplierAccount])
	isQuoted, _, _ := cache.GetQuotedInfoByObjID(p.cs.CacheCli, objID, supplierAccount)
	if isQuoted {
		return true, nil, nil
	}

	var model mapstr.MapStr
	var propertyData map[string]mapstr.MapStr

	if len(info) > 1 {
		// info contains: model, properties
		model = info[0]
		for i := 1; i < len(info); i++ {
			propertyData[util.GetStrByInterface(info[i][common.BKPropertyIDField])] = info[i]
		}
	} else {
		// get model and property info from cache
		var exists bool
		model, exists = cache.GetModelInfoByObjID(p.cs.CacheCli, objID)
		if !exists {
			// skip not exists model
			return true, nil, nil
		}

		propertyData, exists = cache.GetPropertyInfoByObjID(p.cs.CacheCli, objID)
		if !exists {
			propertyData = make(map[string]mapstr.MapStr)
		}
	}

	keywords := []string{objID, util.GetStrByInterface(model[common.BKObjNameField])}

	// all attributes with model metadata is ONE elastic document.
	tableAttrs := make([]mapstr.MapStr, 0)
	for _, attribute := range propertyData {
		propertyType, err := convMetaIDToStr(attribute, common.BKPropertyTypeField)
		if err != nil {
			blog.Errorf("[%s] property type is invalid, data: %+v, rid: %s", p.index, data, rid)
			continue
		}

		if propertyType == common.FieldTypeInnerTable {
			tableAttrs = append(tableAttrs, attribute)
		}

		keywords = append(keywords, util.GetStrByInterface(attribute[common.BKPropertyIDField]),
			util.GetStrByInterface(attribute[common.BKPropertyNameField]))
	}

	// build elastic document.
	document := mapstr.MapStr{
		// we use meta_bk_obj_id to search model, set this id to special null value
		metadata.IndexPropertyID:                nullMetaID,
		metadata.IndexPropertyDataKind:          metadata.DataKindModel,
		metadata.IndexPropertyBKObjID:           objID,
		metadata.IndexPropertyBKSupplierAccount: model[common.BKOwnerIDField],
		metadata.IndexPropertyBKBizID:           model[common.BKAppIDField],
		metadata.IndexPropertyKeywords:          compressKeywords(keywords),
	}

	if err := p.updateModelTableProperties(document, tableAttrs, rid); err != nil {
		blog.Errorf("parse model table attributes failed, table attr: %+v, rid: %s", tableAttrs, rid)
		return false, nil, err
	}

	return false, document, nil
}

// updateModelTableProperties update model table property.
func (p *modelParser) updateModelTableProperties(document mapstr.MapStr, attrs []mapstr.MapStr, rid string) error {
	if len(attrs) == 0 {
		return nil
	}

	tables := make(mapstr.MapStr)
	for _, attribute := range attrs {
		propertyID, err := convMetaIDToStr(attribute, common.BKPropertyIDField)
		if err != nil {
			blog.Errorf("parse property id failed, err: %v, attr: %+v, rid: %s", err, attribute, rid)
			continue
		}

		option, err := metadata.ParseTableAttrOption(attribute[common.BKOptionField])
		if err != nil {
			blog.Errorf("parse table attr option failed, err: %v, attr: %+v, rid: %s", err, attribute, rid)
			continue
		}

		if len(option.Header) == 0 {
			continue
		}

		keywords := make([]string, 0)
		for _, header := range option.Header {
			keywords = append(keywords, header.PropertyID, header.PropertyName)
		}

		// 0 为占位符，保持搜索时模型和实例的统一
		// todo 临时方案，后续优化
		tables[propertyID] = mapstr.MapStr{nullMetaID: compressKeywords(keywords)}
	}

	document[metadata.TablePropertyName] = tables
	return nil
}

// ParseWatchDeleteData parse delete model attribute data from mongodb watch
func (p *modelParser) ParseWatchDeleteData(collOidMap map[string][]string, rid string) ([]string,
	[]elastic.BulkableRequest, bool) {

	needDelIDs := make([]string, 0)
	requests := make([]elastic.BulkableRequest, 0)

	for coll, oids := range collOidMap {
		switch coll {
		case common.BKTableNameObjDes:
			for _, oid := range oids {
				needDelIDs = append(needDelIDs, p.GenEsID(coll, oid))
			}
		case common.BKTableNameObjAttDes:
			delArchives := getDelArchive(collOidMap, rid)

			for _, archive := range delArchives {
				skip, data, err := p.ParseData([]mapstr.MapStr{archive.Detail}, coll, rid)
				if err != nil || skip {
					continue
				}

				id := p.GenEsID(coll, archive.Oid)

				req := elastic.NewBulkUpdateRequest().DocAsUpsert(true).RetryOnConflict(10).Id(id).Doc(data)

				if _, err = req.Source(); err != nil {
					blog.Errorf("upsert data is invalid, err: %v, id: %s, data: %+v, rid: %s", err, id, data, rid)
					continue
				}

				requests = append(requests, req)
			}

		}
	}

	return needDelIDs, requests, false
}
