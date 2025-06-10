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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/cache"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// modelParser is the general data parser for model/attribute
type modelParser struct {
	*commonParser
}

func newModelParser(index string, cs *ClientSet, cache *cache.Cache) *modelParser {
	return &modelParser{newCommonParser(index, cs, cache)}
}

// ParseData parse mongo data to es data
func (p *modelParser) ParseData(kit *rest.Kit, info []mapstr.MapStr, coll string) (bool, mapstr.MapStr, error) {
	if len(info) == 0 {
		return false, nil, errors.New("data is empty")
	}
	data := info[0]

	objID, ok := data[common.BKObjIDField].(string)
	if !ok {
		blog.Errorf("[%s] object id is invalid, data: %+v, rid: %s", p.index, data, kit.Rid)
		return false, nil, errors.New("object id is invalid")
	}

	// skip table object
	isQuoted, _, _ := metadata.ParseModelQuoteDestObjID(objID)
	if isQuoted {
		return true, nil, nil
	}

	// parse model and attributes info
	model, attributes, skip := p.parseModelAndAttributes(kit, info, objID)
	if skip {
		return true, nil, nil
	}

	keywords := []string{objID, util.GetStrByInterface(model[common.BKObjNameField])}

	// all attributes with model metadata is ONE elastic document.
	tableAttrs := make([]mapstr.MapStr, 0)
	enumAttrs := make([]mapstr.MapStr, 0)
	for _, attribute := range attributes {
		propertyType, err := convMetaIDToStr(attribute, common.BKPropertyTypeField)
		if err != nil {
			blog.Errorf("[%s] property type is invalid, data: %+v, rid: %s", p.index, data, kit.Rid)
			continue
		}

		switch propertyType {
		case common.FieldTypeInnerTable:
			tableAttrs = append(tableAttrs, attribute)
		case common.FieldTypeEnum, common.FieldTypeEnumMulti:
			enumAttrs = append(enumAttrs, attribute)
		}

		keywords = append(keywords, util.GetStrByInterface(attribute[common.BKPropertyIDField]),
			util.GetStrByInterface(attribute[common.BKPropertyNameField]))
	}

	p.cache.SetObjEnumInfo(kit.TenantID, objID, enumAttrs)

	// build elastic document.
	document := mapstr.MapStr{
		// we use meta_bk_obj_id to search model, set this id to special null value
		metadata.IndexPropertyID:       nullMetaID,
		metadata.IndexPropertyDataKind: metadata.DataKindModel,
		metadata.IndexPropertyBKObjID:  objID,
		metadata.IndexPropertyTenantID: kit.TenantID,
		metadata.IndexPropertyBKBizID:  model[common.BKAppIDField],
		metadata.IndexPropertyKeywords: compressKeywords(keywords),
	}

	if err := p.updateModelTableProperties(document, tableAttrs, kit.Rid); err != nil {
		blog.Errorf("parse model table attributes failed, table attr: %+v, rid: %s", tableAttrs, kit.Rid)
		return false, nil, err
	}

	if coll == common.BKTableNameObjAttDes {
		oid, ok := model[common.MongoMetaID].(primitive.ObjectID)
		if !ok {
			blog.Errorf("parse model(%+v) oid failed, obj id: %s, rid: %s", model, objID, kit.Rid)
			return true, nil, nil
		}
		document[common.MongoMetaID] = oid.Hex()
	}

	return false, document, nil
}

func (p *modelParser) parseModelAndAttributes(kit *rest.Kit, info []mapstr.MapStr, objID string) (mapstr.MapStr,
	[]mapstr.MapStr, bool) {

	var model mapstr.MapStr
	attributes := make([]mapstr.MapStr, 0)

	if len(info) > 1 {
		// info contains: model, attributes
		model = info[0]
		for i := 1; i < len(info); i++ {
			attributes = append(attributes, info[i])
		}
	} else {
		// get model info from db by object id
		model = make(mapstr.MapStr)
		exists := false

		cond := mapstr.MapStr{common.BKObjIDField: objID}
		ferrors.FatalErrHandler(200, 100, func() error {
			opts := dbtypes.NewFindOpts().SetWithObjectID(true)
			err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(cond, opts).One(kit.Ctx, &model)
			if err != nil {
				if mongodb.IsNotFoundError(err) {
					return nil
				}
				blog.Errorf("get model data failed, cond: %+v, err: %v", cond, err)
				return err
			}

			exists = true
			return nil
		})

		if !exists {
			// skip not exists model
			return nil, nil, true
		}

		// get object attributes from db
		attrCond := mapstr.MapStr{common.BKObjIDField: objID}
		ferrors.FatalErrHandler(200, 100, func() error {
			if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Find(attrCond).All(kit.Ctx,
				&attributes); err != nil {
				blog.Errorf("get model attribute data failed, cond: %+v, err: %v", attrCond, err)
				return err
			}
			return nil
		})
	}
	return model, attributes, false
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
		tables[propertyID] = mapstr.MapStr{nullMetaID: compressKeywords(keywords)}
	}

	document[metadata.TablePropertyName] = tables
	return nil
}

// ParseWatchDeleteData parse delete model attribute data from mongodb watch
func (p *modelParser) ParseWatchDeleteData(kit *rest.Kit, info mapstr.MapStr, coll, oid string) (bool,
	elastic.BulkableRequest, bool) {

	switch coll {
	case common.BKTableNameObjDes:
		return true, nil, false
	case common.BKTableNameObjAttDes:
		skip, data, err := p.ParseData(kit, []mapstr.MapStr{info}, coll)
		if err != nil || skip {
			return false, nil, true
		}

		id := p.GenEsID(kit.TenantID, util.GetStrByInterface(data[common.MongoMetaID]))
		delete(data, common.MongoMetaID)

		req := elastic.NewBulkUpdateRequest().Index(types.GetIndexName(metadata.IndexNameModel)).
			DocAsUpsert(true).RetryOnConflict(10).Id(id).Doc(data)

		if _, err = req.Source(); err != nil {
			blog.Errorf("upsert data is invalid, err: %v, id: %s, data: %+v, rid: %s", err, id, data, kit.Rid)
			return false, nil, true
		}

		return false, req, false
	default:
		blog.Errorf("unsupported collection: %s for model parser, data: %+v, rid: %s", coll, info, kit.Rid)
		return false, nil, true
	}
}
