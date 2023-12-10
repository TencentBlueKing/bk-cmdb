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
	"io"
	"time"

	ftypes "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/cache"
	ferrors "configcenter/src/scene_server/sync_server/logics/full-text-search/errors"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/parser"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SyncData sync full-text search data by index
func (f fullTextSearch) syncDataByIndex(ctx context.Context, index string, rid string) error {
	collections := make([]string, 0)

	switch index {
	case metadata.IndexNameObjectInstance:
		// get all object instance collections by objects
		objs := make([]metadata.Object, 0)

		ferrors.FatalErrHandler(200, 100, func() error {
			err := mongodb.Client().Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
				common.BkSupplierAccount).All(ctx, &objs)
			if err != nil {
				blog.Errorf("get objects failed, err: %v, rid: %s", err, rid)
				return err
			}
			return nil
		})

		for _, obj := range objs {
			isQuoted, _, _ := cache.GetQuotedInfoByObjID(f.cacheCli, obj.ObjectID, obj.OwnerID)
			if isQuoted {
				continue
			}
			collections = append(collections, common.GetObjectInstTableName(obj.ObjectID, obj.OwnerID))
		}
	default:
		collections = append(collections, types.IndexCollMap[index])
	}

	existEsIDMap := make(map[string]struct{})
	for _, collection := range collections {
		existEsIDs, err := f.syncCollection(ctx, index, collection, nil, rid)
		if err != nil {
			return err
		}
		for _, id := range existEsIDs {
			existEsIDMap[id] = struct{}{}
		}
	}

	return f.cleanEsData(ctx, index, existEsIDMap, rid)

}

func (f fullTextSearch) cleanEsData(ctx context.Context, index string, existEsIDMap map[string]struct{},
	rid string) error {

	syncer, err := newDataSyncer(f.esCli.Client, index)
	if err != nil {
		return err
	}

	var scrollID string

	for {
		var scrollRes *elastic.SearchResult
		err = ferrors.EsRespErrHandler(func() (bool, error) {
			scrollRes, err = f.esCli.Client.Scroll(index).ScrollId(scrollID).Do(ctx)
			if err != nil && err != io.EOF {
				blog.Errorf("scroll get %s es data failed, err: %v, rid: %s", index, err, rid)
				return false, err
			}
			return false, nil
		})

		if err == io.EOF {
			return nil
		}

		if err != nil || scrollRes.Hits == nil || scrollRes.Hits.TotalHits == nil {
			blog.Errorf("scroll get %s es data failed, err: %v, res: %v, rid: %s", index, err, scrollRes, rid)
			return err
		}

		if len(scrollRes.Hits.Hits) == 0 {
			return nil
		}

		scrollID = scrollRes.ScrollId

		// delete not exist data in this range
		delEsIDs := make([]string, 0)
		for _, hit := range scrollRes.Hits.Hits {
			_, exists := existEsIDMap[hit.Id]
			if !exists {
				delEsIDs = append(delEsIDs, hit.Id)
			}
		}

		if len(delEsIDs) > 0 {
			syncer.addEsDeleteReq(delEsIDs, rid)
			if err = syncer.doBulk(ctx, rid); err != nil {
				blog.Infof("do %s es bulk delete failed, err: %v, del ids: %+v, rid: %s", index, err, delEsIDs, rid)
				continue
			}
		}
	}
}

// syncCollection upsert full-text search data to es by collection
func (f fullTextSearch) syncCollection(ctx context.Context, index, coll string, oids []string, rid string) (
	[]string, error) {

	syncer, err := newDataSyncer(f.esCli.Client, index)
	if err != nil {
		return nil, err
	}

	// sync data by oids
	if len(oids) > 0 {
		mongoOids := make([]primitive.ObjectID, len(oids))
		for i, oid := range oids {
			mongoOid, err := primitive.ObjectIDFromHex(oid)
			if err != nil {
				blog.Errorf("parse mongodb oid from %s failed, err: %v, rid: %s", oid, err, rid)
				return nil, err
			}
			mongoOids[i] = mongoOid
		}
		cond := mapstr.MapStr{common.MongoMetaID: mapstr.MapStr{common.BKDBIN: mongoOids}}
		f.upsertDataByCond(ctx, syncer, coll, cond, rid)
		return nil, err
	}

	// lock full-text search compensate sync operation
	locker := lock.NewLocker(redis.Client())
	lockKey := genSyncLockKey(coll)
	locked, err := locker.Lock(lockKey, 10*time.Minute)
	if err != nil {
		blog.Errorf("lock full-text search sync failed, key: %s, err: %v, rid: %s", lockKey, err, rid)
		return nil, err
	}

	if !locked {
		return nil, errors.New("there is another sync task running, please wait until it's done")
	}
	defer func() {
		if err = locker.Unlock(); err != nil {
			blog.Errorf("unlock full-text search sync failed, key: %s, err: %v, rid: %s", lockKey, err, rid)
		}
	}()

	// paged get data by _id and sync to es
	existEsIDs := make([]string, 0)
	cond := mapstr.MapStr{}
	for {
		oids := f.upsertDataByCond(ctx, syncer, coll, cond, rid)
		if len(oids) == 0 {
			return existEsIDs, nil
		}

		for _, oid := range oids {
			existEsIDs = append(existEsIDs, syncer.parser.GenEsID(coll, oid.Hex()))
		}

		cond = mapstr.MapStr{common.MongoMetaID: mapstr.MapStr{common.BKDBGT: oids[len(oids)-1]}}
	}
}

// genSyncLockKey generate full-text search sync lock key by collection
func genSyncLockKey(collection string) lock.StrFormat {
	return lock.GetLockKey("full:text:search:sync:%s", collection)
}

// collIndexMap is the map of cmdb collection -> es index name
var collIndexMap = map[string]string{
	common.BKTableNameBaseBizSet: metadata.IndexNameBizSet,
	common.BKTableNameBaseApp:    metadata.IndexNameBiz,
	common.BKTableNameBaseSet:    metadata.IndexNameSet,
	common.BKTableNameBaseModule: metadata.IndexNameModule,
	common.BKTableNameBaseHost:   metadata.IndexNameHost,
	common.BKTableNameObjDes:     metadata.IndexNameModel,
}

// getIndexByColl get es index name by cmdb collection name
func getIndexByColl(collection string) (string, error) {
	index, exists := collIndexMap[collection]
	if exists {
		return index, nil
	}

	if common.IsObjectInstShardingTable(collection) {
		return metadata.IndexNameObjectInstance, nil
	}

	return "", errors.New("collection is invalid")
}

type mapStrWithOid struct {
	Oid    primitive.ObjectID     `bson:"_id"`
	MapStr map[string]interface{} `bson:",inline"`
}

// upsertDataByCond upsert data to es by mongo condition, returns all oids in mongo data
func (f fullTextSearch) upsertDataByCond(ctx context.Context, syncer *dataSyncer, coll string, cond mapstr.MapStr,
	rid string) []primitive.ObjectID {

	findOpt := dbtypes.NewFindOpts().SetWithObjectID(true)
	allData := make([]mapStrWithOid, 0)
	ferrors.FatalErrHandler(200, 100, func() error {
		err := mongodb.Client().Table(coll).Find(cond, findOpt).Sort(common.MongoMetaID).Limit(ftypes.SyncDataPageSize).
			All(ctx, &allData)
		if err != nil {
			blog.Errorf("get data failed, cond: %+v, err: %v, rid: %s", cond, err, rid)
			return err
		}
		return nil
	})

	if len(allData) == 0 {
		return make([]primitive.ObjectID, 0)
	}

	dataMap := dataGetterMap[syncer.index](ctx, coll, allData, rid)

	oids := make([]primitive.ObjectID, len(allData))
	for i, data := range allData {
		oids[i] = data.Oid
		syncer.addUpsertReq(coll, data.Oid.Hex(), dataMap[data.Oid], rid)
	}

	if err := syncer.doBulk(context.Background(), rid); err != nil {
		blog.Errorf("do es bulk request failed, err: %v, coll: %s, cond: %+v, rid: %s", err, coll, cond, rid)
		return oids
	}

	return oids
}

// dataGetter is the data getter to get all oid related sync data
type dataGetter func(context.Context, string, []mapStrWithOid, string) map[primitive.ObjectID][]mapstr.MapStr

// collIndexMap is the map of index -> sync data getter
var dataGetterMap = map[string]dataGetter{
	metadata.IndexNameBizSet:         objInstDataGetter,
	metadata.IndexNameBiz:            objInstDataGetter,
	metadata.IndexNameSet:            objInstDataGetter,
	metadata.IndexNameModule:         objInstDataGetter,
	metadata.IndexNameHost:           objInstDataGetter,
	metadata.IndexNameModel:          modelDataGetter,
	metadata.IndexNameObjectInstance: objInstDataGetter,
}

func modelDataGetter(ctx context.Context, coll string, allData []mapStrWithOid,
	rid string) map[primitive.ObjectID][]mapstr.MapStr {

	objIDs := make([]string, 0)
	for _, data := range allData {
		objIDs = append(objIDs, util.GetStrByInterface(data.MapStr[common.BKObjIDField]))
	}

	// get model related attributes
	attrMap := make(map[string][]mapstr.MapStr)
	cond := mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}}
	fields := []string{common.MongoMetaID, common.BKPropertyTypeField, common.BKPropertyIDField,
		common.BKPropertyNameField}

	attributes := pagedGetMongoData(common.BKTableNameObjAttDes, cond, fields)
	for _, attribute := range attributes {
		objID := util.GetStrByInterface(attribute[common.BKObjIDField])
		attrMap[objID] = append(attrMap[objID], attribute)
	}

	// model info contains the model data and its attributes
	dataMap := make(map[primitive.ObjectID][]mapstr.MapStr)
	for _, data := range allData {
		attr := attrMap[util.GetStrByInterface(data.MapStr[common.BKObjIDField])]
		dataMap[data.Oid] = append([]mapstr.MapStr{data.MapStr}, attr...)
	}
	return dataMap
}

func objInstDataGetter(ctx context.Context, coll string, allData []mapStrWithOid,
	rid string) map[primitive.ObjectID][]mapstr.MapStr {

	objID := parser.GetObjIDByData(coll, allData[0].MapStr)

	instIDs := make([]int64, len(allData))
	for i, data := range allData {
		id, err := util.GetInt64ByInterface(data.MapStr[common.GetInstIDField(objID)])
		if err != nil {
			blog.Errorf("get instance id failed, err: %v, data: %+v, rid: %s", err, data.MapStr, rid)
			continue
		}
		instIDs[i] = id
	}

	// get all model quote relations by src obj id
	supplierAccount := util.GetStrByInterface(allData[0].MapStr[common.BkSupplierAccount])
	relCond := mapstr.MapStr{
		common.BKSrcModelField:   objID,
		common.BkSupplierAccount: supplierAccount,
	}

	relations := make([]metadata.ModelQuoteRelation, 0)
	ferrors.FatalErrHandler(200, 100, func() error {
		err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(relCond).All(ctx, &relations)
		if err != nil {
			blog.Errorf("get model quote relation failed, cond: %+v, err: %v", relCond, err)
			return err
		}
		return nil
	})

	// get all quote instances
	instMap := make(map[int64][]mapstr.MapStr)

	for _, relation := range relations {
		table := common.GetObjectInstTableName(relation.DestModel, supplierAccount)
		cond := mapstr.MapStr{common.BKInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}}

		instances := pagedGetMongoData(table, cond, make([]string, 0))
		for _, instance := range instances {
			instID, err := util.GetInt64ByInterface(instance[common.BKInstIDField])
			if err != nil {
				blog.Errorf("get quote instance id failed, err: %v, instance: %+v, rid: %s", err, instance, rid)
				continue
			}
			instance[common.BKPropertyIDField] = relation.PropertyID
			instMap[instID] = append(instMap[instID], instance)
		}
	}

	// instance info contains the instance data and its quote instances
	dataMap := make(map[primitive.ObjectID][]mapstr.MapStr)
	for _, data := range allData {
		id, _ := util.GetInt64ByInterface(data.MapStr[common.GetInstIDField(objID)])
		quote := instMap[id]
		dataMap[data.Oid] = append([]mapstr.MapStr{data.MapStr}, quote...)
	}

	return dataMap
}

func pagedGetMongoData(table string, cond mapstr.MapStr, fields []string) []mapstr.MapStr {
	allData := make([]mapstr.MapStr, 0)

	findOpt := dbtypes.NewFindOpts().SetWithObjectID(true)
	for {
		data := make([]mapstr.MapStr, 0)
		ferrors.FatalErrHandler(200, 100, func() error {
			err := mongodb.Client().Table(table).Find(cond, findOpt).Fields(fields...).Sort(common.MongoMetaID).
				Limit(ftypes.SyncDataPageSize).All(context.Background(), &data)
			if err != nil {
				blog.Errorf("get quote instance failed, table: %s, cond: %+v, err: %v", table, cond, err)
				return err
			}
			return nil
		})

		if len(data) == 0 {
			return allData
		}
		allData = append(allData, data...)

		cond[common.MongoMetaID] = mapstr.MapStr{common.BKDBGT: data[len(data)-1][common.MongoMetaID]}
	}
}
