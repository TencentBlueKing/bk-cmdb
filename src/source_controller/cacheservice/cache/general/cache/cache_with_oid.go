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

package cache

import (
	"fmt"
	"time"

	"configcenter/pkg/cache/general"
	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type dataWithOid[T any] struct {
	Oid  primitive.ObjectID `json:"-" bson:"_id"`
	Data T                  `json:",inline" bson:",inline"`
}

// MarshalJSON marshal json
func (data dataWithOid[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(data.Data)
}

// newCacheWithOid new general cache whose data uses id as id key
func newCacheWithOid[T any](key *general.Key, needCacheAll bool, table string,
	parser func(data T) (*basicInfo, error)) *Cache {

	cache := NewCache()
	cache.key = key
	cache.expireSeconds = 30 * 60 * time.Second
	cache.expireRangeSeconds = [2]int{-600, 600}
	cache.needCacheAll = needCacheAll
	cache.parseData = parseDataWithOid[T](parser)
	cache.getDataByID = getDataByOid[T](table)
	cache.listData = listDataWithOid[T](table)
	return cache
}

func parseDataWithOid[T any](parser func(data T) (*basicInfo, error)) dataParser {
	return func(data any) (*basicInfo, error) {
		var info *basicInfo
		switch val := data.(type) {
		case dataWithOid[T]:
			var err error
			info, err = parser(val.Data)
			if err != nil {
				return nil, err
			}
			info.oid = val.Oid.Hex()
		case filter.JsonString:
			info = &basicInfo{
				oid: gjson.Get(string(val), common.MongoMetaID).String(),
			}
		default:
			return nil, fmt.Errorf("data type %T is invalid", data)
		}

		if info.oid == "" {
			return nil, fmt.Errorf("oid is zero")
		}

		return info, nil
	}
}

func getDataByOid[T any](table string) dataGetterByKeys {
	return func(kit *rest.Kit, opt *getDataByKeysOpt) ([]any, error) {
		kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

		if len(opt.Keys) == 0 {
			return make([]any, 0), nil
		}

		oids := make([]primitive.ObjectID, len(opt.Keys))
		for i, key := range opt.Keys {
			oid, err := primitive.ObjectIDFromHex(key)
			if err != nil {
				blog.Errorf("parse oid (index: %d, key: %s) failed, err: %v, rid: %s", i, key, err, kit.Rid)
				return nil, err
			}
			oids[i] = oid
		}

		cond := mapstr.MapStr{
			common.MongoMetaID: mapstr.MapStr{common.BKDBIN: oids},
		}

		dbOpts := dbtypes.NewFindOpts().SetWithObjectID(true)

		dataArr := make([]dataWithOid[T], 0)
		if err := mongodb.Shard(kit.ShardOpts()).Table(table).Find(cond, dbOpts).All(kit.Ctx, &dataArr); err != nil {
			blog.Errorf("get %s data by cond(%+v) failed, err: %v, rid: %s", table, cond, err, kit.Rid)
			return nil, err
		}

		return convertToAnyArr(dataArr), nil
	}
}

func listDataWithOid[T any](table string) dataLister {
	return func(kit *rest.Kit, opt *listDataOpt) (*listDataRes, error) {
		kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

		if rawErr := opt.Validate(false); rawErr.ErrCode != 0 {
			blog.Errorf("list general data option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
			return nil, fmt.Errorf("list data option is invalid")
		}

		if opt.OnlyListID {
			opt.Fields = []string{common.MongoMetaID}
		}

		cond := opt.Cond
		if cond == nil {
			cond = make(mapstr.MapStr)
		}

		if opt.Page.EnableCount {
			cnt, err := mongodb.Shard(kit.ShardOpts()).Table(table).Find(cond).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("count %s data by cond(%+v) failed, err: %v, rid: %s", table, cond, err, kit.Rid)
				return nil, err
			}

			return &listDataRes{Count: cnt}, nil
		}

		if opt.Page.StartOid != "" {
			oid, err := primitive.ObjectIDFromHex(opt.Page.StartOid)
			if err != nil {
				blog.Errorf("parse start oid %s failed, err: %v, rid: %s", opt.Page.StartOid, err, kit.Rid)
				return nil, err
			}

			_, exists := cond[common.MongoMetaID]
			if exists {
				cond = mapstr.MapStr{
					common.BKDBAND: []mapstr.MapStr{{common.MongoMetaID: mapstr.MapStr{common.BKDBGT: oid}}, cond},
				}
			} else {
				cond[common.MongoMetaID] = mapstr.MapStr{common.BKDBGT: oid}
			}
		}

		dbOpts := dbtypes.NewFindOpts().SetWithObjectID(true)

		dataArr := make([]dataWithOid[T], 0)
		err := mongodb.Shard(kit.ShardOpts()).Table(table).Find(cond, dbOpts).Sort(common.MongoMetaID).Start(
			uint64(opt.Page.StartIndex)).Limit(uint64(opt.Page.Limit)).Fields(opt.Fields...).All(kit.Ctx, &dataArr)
		if err != nil {
			blog.Errorf("list %s data by cond(%+v) failed, err: %v, rid: %s", table, cond, err, kit.Rid)
			return nil, err
		}

		return &listDataRes{Data: convertToAnyArr(dataArr)}, nil
	}
}
