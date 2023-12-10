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

// Package cache defines full-text search caching logics
package cache

import (
	"context"
	"errors"
	"sync"

	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	cachekey "configcenter/src/source_controller/cacheservice/cache/common/key"
	"configcenter/src/storage/driver/mongodb"
)

const synUser = "cc_full_text_search_sync"

func getCacheInfo(cli cacheservice.Cache, data, desInfo interface{}, typ cachekey.KeyType,
	kind cachekey.KeyKind) error {

	ctx := context.Background()
	header := util.BuildHeader(synUser, common.BKDefaultOwnerID)

	generator, err := cachekey.GetKeyGenerator(typ)
	if err != nil {
		blog.Errorf("get %s key generator failed, err: %v", typ, err)
		return err
	}

	redisKey, err := generator.GenerateRedisKey(kind, data)
	if err != nil {
		blog.Errorf("get %s kind: %s redis key from data: %+v failed, err: %v", typ, kind, data, err)
		return err
	}

	opt := &metadata.ListCommonCacheWithKeyOpt{
		Kind: string(kind),
		Keys: []string{redisKey},
	}

	infoJs, err := cli.CommonRes().ListWithKey(ctx, header, string(typ), opt)
	if err != nil {
		blog.Errorf("list %s data from cache failed, err: %v, opt: %+v", typ, err, opt)
		return err
	}

	err = json.Unmarshal([]byte(infoJs), desInfo)
	if err != nil {
		blog.Errorf("unmarshal %s cache info %s failed, err: %v", typ, infoJs, err)
		return err
	}

	return nil
}

// GetQuotedInfoByObjID get the quoted object id related property id and src obj id
func GetQuotedInfoByObjID(cli cacheservice.Cache, objID, supplierAccount string) (bool, string, string) {
	// get quoted info from cache
	quoteInfo := make([]metadata.ModelQuoteRelation, 0)
	err := getCacheInfo(cli, objID, &quoteInfo, cachekey.ModelQuoteRelType, cachekey.DestModelKind)
	if err == nil {
		for _, relation := range quoteInfo {
			if relation.SupplierAccount == supplierAccount {
				return true, quoteInfo[0].PropertyID, quoteInfo[0].SrcModel
			}
		}

		return false, "", ""
	}

	// get quoted info from db by dest model id
	cond := mapstr.MapStr{
		common.BKDestModelField:  objID,
		common.BkSupplierAccount: supplierAccount,
	}
	exists := false
	rel := new(metadata.ModelQuoteRelation)

	ferrors.FatalErrHandler(200, 100, func() error {
		err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(cond).One(context.Background(), &rel)
		if err != nil {
			if mongodb.Client().IsNotFoundError(err) {
				return nil
			}
			blog.Errorf("get model quote relation failed, cond: %+v, err: %v", cond, err)
			return err
		}

		exists = true
		return nil
	})

	if !exists {
		return false, "", ""
	}

	return true, rel.PropertyID, rel.SrcModel
}

// GetModelInfoByObjID get model info by object id
func GetModelInfoByObjID(cli cacheservice.Cache, objID string) (mapstr.MapStr, bool) {
	// get model info from cache
	objects := make([]mapstr.MapStr, 0)
	err := getCacheInfo(cli, objID, &objects, cachekey.ModelType, cachekey.ObjIDKind)
	if err == nil {
		if len(objects) == 0 {
			return make(mapstr.MapStr), false
		}

		return objects[0], true
	}

	// get model info from db by object id
	cond := mapstr.MapStr{common.BKObjIDField: objID}
	exists := false
	data := make(mapstr.MapStr)

	ferrors.FatalErrHandler(200, 100, func() error {
		err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond).One(context.Background(), &data)
		if err != nil {
			if mongodb.Client().IsNotFoundError(err) {
				return nil
			}
			blog.Errorf("get model data failed, cond: %+v, err: %v", cond, err)
			return err
		}

		exists = true
		return nil
	})

	return data, exists
}

// GetPropertyInfoByObjID get property id to info map by object id
func GetPropertyInfoByObjID(cli cacheservice.Cache, objID string) (map[string]mapstr.MapStr, bool) {
	properties := getPropertiesByObjID(cli, objID)

	if len(properties) == 0 {
		return make(map[string]mapstr.MapStr), false
	}

	propertyInfo := make(map[string]mapstr.MapStr)
	for _, property := range properties {
		propID := util.GetStrByInterface(property[common.BKPropertyIDField])
		propertyInfo[propID] = property
	}

	return propertyInfo, true
}

func getPropertiesByObjID(cli cacheservice.Cache, objID string) []mapstr.MapStr {
	// get model info from cache
	properties := make([]mapstr.MapStr, 0)
	err := getCacheInfo(cli, objID, &properties, cachekey.AttributeType, cachekey.ObjIDKind)
	if err != nil {
		// get model info from db to compensate
		cond := mapstr.MapStr{common.BKObjIDField: objID}

		ferrors.FatalErrHandler(200, 100, func() error {
			err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(context.Background(), &properties)
			if err != nil {
				blog.Errorf("get model attribute data failed, cond: %+v, err: %v", cond, err)
				return err
			}

			return nil
		})
	}

	return properties
}

// EnumIDToName change instance data enum id to enum name.
func EnumIDToName(cli cacheservice.Cache, document mapstr.MapStr, objID string) mapstr.MapStr {
	properties := getPropertiesByObjID(cli, objID)

	if len(properties) == 0 {
		return document
	}

	for _, property := range properties {
		propType := util.GetStrByInterface(property[common.BKPropertyTypeField])

		if propType != common.FieldTypeEnum {
			continue
		}

		propID := util.GetStrByInterface(property[common.BKPropertyIDField])
		if _, ok := document[propID]; !ok {
			continue
		}

		docVal, ok := document[propID].(string)
		if !ok {
			continue
		}

		option, err := metadata.ParseEnumOption(property[common.BKOptionField])
		if err != nil {
			blog.Errorf("parse %v enum option failed, err: %v", property, err)
			continue
		}

		for _, opt := range option {
			if opt.ID == docVal {
				document[propID] = opt.Name
				break
			}
		}
	}

	return document
}

// ResPoolBizIDMap is used to judge if biz/set... is in resource pool
var ResPoolBizIDMap = sync.Map{}

// InitResourcePoolBiz initialize resource pool biz info
// NOTE: right now resource pool cannot be operated, so we don't need to change it.
func InitResourcePoolBiz() error {
	resPoolCond := mapstr.MapStr{common.BKDefaultField: common.DefaultAppFlag}
	bizs := make([]metadata.BizInst, 0)
	err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(resPoolCond).Fields(common.BKAppIDField).
		All(context.Background(), &bizs)
	if err != nil {
		return err
	}

	if len(bizs) == 0 {
		return errors.New("there's no resource pool biz")
	}

	for _, biz := range bizs {
		ResPoolBizIDMap.Store(biz.BizID, struct{}{})
	}

	return nil
}
