/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/index"
	"configcenter/src/common/metadata"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"
)

/*
 如何展示错误给用户
*/

// DBSync TODO
func DBSync(e *backbone.Engine, db dal.RDB, options options.Config) {
	f := func() {
		defaultDBTable = db
		fmt.Println(defaultDBTable)
	}
	once.Do(f)

	go func() {
		RunSyncDBTableIndex(context.Background(), e, db, options)
	}()

}

var (
	once sync.Once

	defaultDBTable dal.RDB
)

type dbTable struct {
	db                         dal.RDB
	preCleanRedundancyTableMap map[string]struct{}
	rid                        string
	options                    options.Config
}

// RunSyncDBTableIndex TODO
func RunSyncDBTableIndex(ctx context.Context, e *backbone.Engine, db dal.RDB,
	options options.Config) {

	rid := util.GenerateRID()
	for dbReady := false; !dbReady; {
		// 等待数据库初始化
		if !dbReady {
			var err error
			dbReady, err = upgrader.DBReady(ctx, db)
			if err != nil {
				blog.Errorf("Check whether the db initialization is complete error. err: %s rid: %s", err.Error(), rid)
			}
			if !dbReady {
				blog.Errorf("db not initialization is complete, rid: %s", rid)
			}

			time.Sleep(20 * time.Second)
		}
	}

	syncWorker := func(dt *dbTable, isTable bool) {
		for {
			rid := util.GenerateRID()
			dt.rid = rid
			blog.Infof("start sync table or index worker rid: %s", rid)

			if !e.ServiceManageInterface.IsMaster() {
				blog.Infof("skip sync table or index worker. reason: not master. rid: %s", rid)
				time.Sleep(5 * time.Second)
				continue
			}

			if isTable {
				blog.Infof("start object sharding table rid: %s", rid)
				// 先处理模型实例和关联关系表
				if err := dt.syncModelShardingTable(ctx); err != nil {
					blog.Errorf("model table sync error. err: %s, rid: %s", err.Error(), dt.rid)
				}
				blog.Infof("end sync table rid: %s", rid)
				time.Sleep(time.Second * time.Duration(options.ShardingTable.TableInterval))

			} else {
				blog.Infof("start table common index rid: %s", rid)
				if err := dt.syncIndexes(ctx); err != nil {
					blog.Errorf("model table sync error. err: %s, rid: %s", err.Error(), dt.rid)
				}
				blog.Infof("end sync table index rid: %s", rid)
				time.Sleep(time.Minute * time.Duration(options.ShardingTable.IndexesInterval))
			}

		}
	}

	dtTable := &dbTable{db: db, rid: rid, options: options}
	go syncWorker(dtTable, true)
	dtIndex := &dbTable{db: db, rid: rid, options: options}
	go syncWorker(dtIndex, false)

}

// RunSyncDBIndex TODO
func RunSyncDBIndex(ctx context.Context, e *backbone.Engine) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	ccErr := e.CCErr.CreateDefaultCCErrorIf("en")

	if defaultDBTable == nil {
		blog.Errorf("db client not initialization is complete, rid: %s", rid)
		return ccErr.CCError(common.CCErrCommDBSelectFailed)
	}
	// defaultDBTable DBSync 负责在启动时候初始化
	dbReady, err := upgrader.DBReady(ctx, defaultDBTable)
	if err != nil {
		blog.Errorf("Check whether the db initialization is complete error. err: %s rid: %s", err.Error(), rid)
		return ccErr.CCError(common.CCErrCommDBSelectFailed)
	}
	if !dbReady {
		blog.Errorf("db not initialization is complete, rid: %s", rid)
		return ccErr.CCError(common.CCErrCommDBSelectFailed)

	}

	dt := &dbTable{db: defaultDBTable, rid: rid}
	blog.Infof("start table common index rid: %s", rid)
	if err := dt.syncIndexes(ctx); err != nil {
		blog.Errorf("model table sync error. err: %s, rid: %s", err.Error(), dt.rid)
	}
	blog.Infof("end sync table index rid: %s", rid)

	return nil
}

// syncIndexes 同步表中定义的索引
func (dt *dbTable) syncIndexes(ctx context.Context) error {
	if err := dt.syncDBTableIndexes(ctx); err != nil {
		blog.Warnf("sync table index to db error. err: %s, rid: %s", err.Error(), dt.rid)
		// 不影响后需任务
	}

	return nil

}

func (dt *dbTable) syncDBTableIndexes(ctx context.Context) error {
	deprecatedIndexNames := index.DeprecatedIndexName()
	tableIndexes := index.TableIndexes()

	dtIndexesMap, err := dt.findSyncIndexesLogicUnique(ctx)
	if err != nil {
		blog.ErrorJSON("find db logic unique error. err: %s, rid: %s", err, dt.rid)
		return err
	}

	for tableName, indexes := range tableIndexes {
		blog.Infof("start sync table(%s) index, rid: %s", tableName, dt.rid)
		deprecatedTableIndexNames := deprecatedIndexNames[tableName]

		indexes = append(indexes, dtIndexesMap[tableName]...)
		delete(dtIndexesMap, tableName)
		if err := dt.syncIndexesToDB(ctx, tableName, indexes, deprecatedTableIndexNames); err != nil {
			blog.Warnf("sync table (%s) index failed. err: %v, rid: %s", tableName, err, dt.rid)
			continue
		}
	}

	for tableName, indexes := range dtIndexesMap {
		blog.Infof("start sync table(%s) index, rid: %s", tableName, dt.rid)
		deprecatedTableIndexNames := deprecatedIndexNames[tableName]

		if err := dt.syncIndexesToDB(ctx, tableName, indexes, deprecatedTableIndexNames); err != nil {
			blog.Warnf("sync table (%s) index failed. err: %v, rid: %s", tableName, err, dt.rid)
			continue
		}
	}

	return nil
}

// findObjAttrs 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
func (dt *dbTable) findObjAttrs(ctx context.Context, objID string) ([]metadata.Attribute, error) {
	// 获取字段类型,只需要共有字段
	attrFilter := map[string]interface{}{
		common.BKObjIDField: objID,
		common.BKAppIDField: 0,
	}
	attrs := make([]metadata.Attribute, 0)
	fields := []string{common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID}
	if err := dt.db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Fields(fields...).All(ctx, &attrs); err != nil {
		newErr := fmt.Errorf("get obj(%s) property error. err: %s", objID, err.Error())
		blog.Errorf("%s, rid: %s", newErr.Error(), dt.rid)
		return nil, newErr
	}

	return attrs, nil
}

func (dt *dbTable) syncIndexesToDB(ctx context.Context, tableName string,
	logicIndexes []types.Index, deprecatedTableIndexNames []string) (err error) {

	dbIndexList, err := dt.db.Table(tableName).Indexes(ctx)
	if err != nil {
		blog.Errorf("find db table(%s) index list error. err: %s, rid: %s", tableName, err.Error(), dt.rid)
		return err
	}

	dbIdxNameMap := make(map[string]types.Index, len(dbIndexList))
	// 等待删除处理索引名字,所有规范后索引名字
	waitDelIdxNameMap := make(map[string]struct{}, 0)
	for _, index := range dbIndexList {

		dbIdxNameMap[index.Name] = index
		if strings.HasPrefix(index.Name, common.CCLogicIndexNamePrefix) ||
			strings.HasPrefix(index.Name, common.CCLogicUniqueIdxNamePrefix) {
			waitDelIdxNameMap[index.Name] = struct{}{}
		}

	}
	//  加入所有不规范索引的索引名字
	for _, indexName := range deprecatedTableIndexNames {
		waitDelIdxNameMap[indexName] = struct{}{}
	}

	for _, logicIndex := range logicIndexes {
		delete(waitDelIdxNameMap, logicIndex.Name)
	}

	for indexName := range waitDelIdxNameMap {
		if err := dt.db.Table(tableName).DropIndex(ctx, indexName); err != nil &&
			!ErrDropIndexNameNotFound(err) {
			blog.Errorf("remove redundancy table(%s) index(%s) error. err: %s, rid: %s",
				tableName, indexName, err.Error(), dt.rid)
			return err
		}
	}

	for _, logicIndex := range logicIndexes {
		// 是否存在同名
		dbIndex, indexNameExist := dbIdxNameMap[logicIndex.Name]
		if indexNameExist {
			if err := dt.tryUpdateTableIndex(ctx, tableName, dbIndex, logicIndex); err != nil {
				// 不影响后需执行，
				blog.Error("try update table index err: %s, rid: %s", err.Error(), dt.rid)
				continue
			}
		} else {
			dt.createIndexes(ctx, tableName, []types.Index{logicIndex})
		}
	}

	return nil

}

func (dt *dbTable) findSyncIndexesLogicUnique(ctx context.Context) (map[string][]types.Index, error) {
	objs := make([]metadata.Object, 0)
	if err := dt.db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		common.BKIsPre, common.BKOwnerIDField).All(ctx, &objs); err != nil {
		blog.Errorf("get all common object id  error. err: %s, rid: %s", err.Error(), dt.rid)
		return nil, err
	}

	tbIndexes := make(map[string][]types.Index)
	for _, obj := range objs {
		blog.Infof("start object(%s) sharding table rid: %s", obj.ObjectID, dt.rid)

		instTable := common.GetObjectInstTableName(obj.ObjectID, obj.OwnerID)
		instAsstTable := common.GetObjectInstAsstTableName(obj.ObjectID, obj.OwnerID)

		uniques, err := dt.findObjUniques(ctx, obj.ObjectID)
		if err != nil {
			blog.Errorf("object(%s) logic unique to db index error. err: %s, rid: %s",
				obj.ObjectID, err.Error(), dt.rid)
			return nil, err
		}
		// 内置模型不需要简表
		if !obj.IsPre {
			tbIndexes[instTable] = append(index.InstanceIndexes(), uniques...)
		} else {
			tb := ""
			switch obj.ObjectID {
			case common.BKInnerObjIDHost:
				tb = common.BKTableNameBaseHost
			case common.BKInnerObjIDBizSet:
				tb = common.BKTableNameBaseBizSet
			case common.BKInnerObjIDApp:
				tb = common.BKTableNameBaseApp
			case common.BKInnerObjIDModule:
				tb = common.BKTableNameBaseModule
			case common.BKInnerObjIDSet:
				tb = common.BKTableNameBaseSet
			case common.BKInnerObjIDPlat:
				tb = common.BKTableNameBasePlat
			case common.BKInnerObjIDProc:
				tb = common.BKTableNameBaseProcess
			}
			if tb != "" {
				tbIndexes[tb] = uniques
			}

		}
		tbIndexes[instAsstTable] = index.InstanceAssociationIndexes()

	}

	return tbIndexes, nil
}

func (dt *dbTable) tryUpdateTableIndex(ctx context.Context, tableName string,
	dbIndex, logicIndex types.Index) error {
	if index.IndexEqual(dbIndex, logicIndex) {
		// db collection 中的索引和定义所以的索引一致，无需处理
		return nil
	} else {
		// 说明索引不等， 删除原有的索引，
		if err := dt.db.Table(tableName).DropIndex(ctx, logicIndex.Name); err != nil {
			blog.Errorf("remove table(%s) index(%s) error. err: %s, rid: %s",
				tableName, logicIndex.Name, err.Error(), dt.rid)
			return err
		}
		if err := dt.db.Table(tableName).CreateIndex(ctx, logicIndex); err != nil {
			blog.Errorf("create table(%s) index(%s) error. err: %s, rid: %s",
				tableName, logicIndex.Name, err.Error(), dt.rid)
			monitor.Collect(&meta.Alarm{
				RequestID: dt.rid,
				Type:      meta.MongoDDLFatalError,
				Detail:    fmt.Sprintf("collection(%s)  create index failed", tableName),
				Module:    types2.CC_MODULE_MIGRATE,
				Dimension: map[string]string{"hit_create_index": "yes"},
			})
			return err
		}
	}
	return nil
}

func (dt *dbTable) syncModelShardingTable(ctx context.Context) error {

	allDBTables, err := dt.db.ListTables(ctx)
	if err != nil {
		blog.Errorf("show tables error. err: %s, rid: %s", err.Error(), dt.rid)
		return err
	}

	modelDBTableNameMap := make(map[string]struct{}, 0)
	for _, name := range allDBTables {
		if common.IsObjectShardingTable(name) {
			modelDBTableNameMap[name] = struct{}{}
		}
	}

	objs := make([]metadata.Object, 0)
	if err := dt.db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		common.BKIsPre, common.BKOwnerIDField).All(ctx, &objs); err != nil {
		blog.Errorf("get all common object id  error. err: %s, rid: %s", err.Error(), dt.rid)
		return err
	}

	for _, obj := range objs {
		blog.Infof("start object(%s) sharding table rid: %s", obj.ObjectID, dt.rid)

		instTable := common.GetObjectInstTableName(obj.ObjectID, obj.OwnerID)
		instAsstTable := common.GetObjectInstAsstTableName(obj.ObjectID, obj.OwnerID)

		uniques, err := dt.findObjUniques(ctx, obj.ObjectID)
		if err != nil {
			blog.Errorf("object(%s) logic unique to db index error. err: %s, rid: %s",
				obj.ObjectID, err.Error(), dt.rid)
			monitor.Collect(&meta.Alarm{
				RequestID: dt.rid,
				Type:      meta.MongoDDLFatalError,
				Detail:    fmt.Sprintf("query %s collection logic unique detail failed", instTable),
				Module:    types2.CC_MODULE_MIGRATE,
				Dimension: map[string]string{"hit_create_collection": "yes"},
			})
			return err
		}

		objIndexes := append(index.InstanceIndexes(), uniques...)
		// 内置模型不需要简表
		if !obj.IsPre {
			// 判断模型实例表是否存在, 不存在新建
			if _, exist := modelDBTableNameMap[instTable]; !exist {
				if err := dt.db.CreateTable(ctx, instTable); err != nil {
					// TODO: 需要报警，但是不影响后需逻辑继续执行下一个
					blog.Errorf("create table(%s) error. err: %s, rid: %s", instTable, err.Error(), dt.rid)
					// NOTICE: 索引有专门的任务处理
					monitor.Collect(&meta.Alarm{
						RequestID: dt.rid,
						Type:      meta.MongoDDLFatalError,
						Detail:    fmt.Sprintf("create %s collection failed", instTable),
						Module:    types2.CC_MODULE_MIGRATE,
						Dimension: map[string]string{"hit_create_collection": "yes"},
					})
				}
				dt.createIndexes(ctx, instTable, objIndexes)
			} else {
				if err := dt.syncIndexesToDB(ctx, instTable, objIndexes, nil); err != nil {
					blog.Errorf("sync table(%s) definition index to table error. err: %s, rid: %s",
						instTable, err.Error(), dt.rid)
				}
			}
		}

		// 判断模型实例关联关系表是否存在， 不存在新建
		if _, exist := modelDBTableNameMap[instAsstTable]; !exist {
			if err := dt.db.CreateTable(ctx, instAsstTable); err != nil {
				// TODO: 需要报警，但是不影响后需逻辑继续执行下一个
				blog.Errorf("create table(%s) error. err: %s, rid: %s", instAsstTable, err.Error(), dt.rid)
				// NOTICE: 索引有专门的任务处理
			}
			monitor.Collect(&meta.Alarm{
				RequestID: dt.rid,
				Type:      meta.MongoDDLFatalError,
				Detail:    fmt.Sprintf("create %s collection failed", instAsstTable),
				Module:    types2.CC_MODULE_MIGRATE,
				Dimension: map[string]string{"hit_create_collection": "yes"},
			})
			dt.createIndexes(ctx, instAsstTable, index.InstanceAssociationIndexes())
		} else {
			if err := dt.syncIndexesToDB(ctx, instAsstTable, index.InstanceAssociationIndexes(), nil); err != nil {
				blog.Errorf("sync table(%s) definition index to table error. err: %s, rid: %s",
					instTable, err.Error(), dt.rid)
			}
		}

		delete(modelDBTableNameMap, instTable)
		delete(modelDBTableNameMap, instAsstTable)
	}

	if err := dt.cleanRedundancyTable(ctx, modelDBTableNameMap); err != nil {
		blog.Errorf("clean redundancy table Name map:(%#v) error. err: %s, rid: %s",
			modelDBTableNameMap, err.Error(), dt.rid)

	}
	return nil
}

func (dt *dbTable) cleanRedundancyTable(ctx context.Context, modelDBTableNameMap map[string]struct{}) error {
	objs := make([]metadata.Object, 0)
	if err := dt.db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		common.BKIsPre, common.BKOwnerIDField).All(ctx, &objs); err != nil {
		blog.Errorf("get all common object id  error. err: %s, rid: %s", err.Error(), dt.rid)
		// NOTICE: 错误直接忽略不行后需功能
		return err
	}

	// 再次确认数据，保证存在模型的的表不被删除
	for _, obj := range objs {
		instTable := common.GetObjectInstTableName(obj.ObjectID, obj.OwnerID)
		instAsstTable := common.GetObjectInstAsstTableName(obj.ObjectID, obj.OwnerID)
		delete(modelDBTableNameMap, instTable)
		delete(modelDBTableNameMap, instAsstTable)

	}

	// 清理表的时候，需要有延时，表至少要生存两个定时删除的周期
	// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
	// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
	// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
	// Collection minimum is Timestamp(1616747878, 5)
	if len(dt.preCleanRedundancyTableMap) == 0 {
		dt.preCleanRedundancyTableMap = modelDBTableNameMap
		return nil
	}

	preCleanRedundancyTableMap := make(map[string]struct{}, 0)
	for name := range modelDBTableNameMap {
		// 上个周期不存在，不删除表
		if _, exists := dt.preCleanRedundancyTableMap[name]; !exists {
			// 下个周可以删除表的
			blog.Errorf("skip redundant table(%s), reason: first appearance, rid: %s", name, dt.rid)
			preCleanRedundancyTableMap[name] = struct{}{}
			continue
		}
		// 检查是否有数据
		cnt, err := dt.db.Table(name).Find(nil).Count(ctx)
		if err != nil {
			blog.Errorf("count table(%s) failed, skip, err: %v, rid: %s", name, err, dt.rid)
			continue
		}

		blog.Infof("find redundant table(%s), try to delete now, rid: %s", name, dt.rid)
		if cnt == 0 {
			blog.Infof("delete sharding table(%s) rid: %s", name, dt.rid)
			// 检查表最后操作事件， 如果最后操作事件小于60s不删除， 避免count 查询不到事务里面的数据

			// 没有数据删除
			if err := dt.db.DropTable(ctx, name); err != nil {
				blog.Errorf("delete table(%s) error. err: %s, rid: %s", name, err.Error(), dt.rid)
				monitor.Collect(&meta.Alarm{
					RequestID: dt.rid,
					Type:      meta.MongoDDLFatalError,
					Detail:    fmt.Sprintf("drop collection(%s) failed", name),
					Module:    types2.CC_MODULE_MIGRATE,
					Dimension: map[string]string{"hit_clean_redundancy_table": "yes"},
				})
				continue

			} else {
				monitor.Collect(&meta.Alarm{
					RequestID: dt.rid,
					Type:      meta.MongoDDLFatalError,
					Detail:    fmt.Sprintf("drop collection(%s) failed, reason: find table has error", name),
					Module:    types2.CC_MODULE_MIGRATE,
					Dimension: map[string]string{"hit_clean_redundancy_table": "yes"},
				})
				blog.Errorf("drop collection(%s) failed, reason: find table has error, rid: %s", name, dt.rid)

			}

		} else {
			monitor.Collect(&meta.Alarm{
				RequestID: dt.rid,
				Type:      meta.MongoDDLFatalError,
				Detail:    fmt.Sprintf("drop collection(%s) failed, reason: non-empty sharding table", name),
				Module:    types2.CC_MODULE_MIGRATE,
				Dimension: map[string]string{"hit_clean_redundancy_table": "yes"},
			})
			blog.Errorf("can't drop the non-empty sharding table, table name: %s, rid: %s", name, dt.rid)
		}

	}

	dt.preCleanRedundancyTableMap = preCleanRedundancyTableMap
	return nil
}

func (dt *dbTable) createIndexes(ctx context.Context, tableName string, indexes []types.Index) {

	for _, index := range indexes {
		if err := dt.db.Table(tableName).CreateIndex(ctx, index); err != nil {
			// 不影响后需执行，
			// TODO: 报警
			blog.WarnJSON("create table(%s) error. index: %s, err: %s, rid: %s", tableName, index, err, dt.rid)
			monitor.Collect(&meta.Alarm{
				RequestID: dt.rid,
				Type:      meta.MongoDDLFatalError,
				Detail:    fmt.Sprintf("collection(%s) create index(%s) failed", tableName, index.Name),
				Module:    types2.CC_MODULE_MIGRATE,
				Dimension: map[string]string{"hit_create_index": "yes"},
			})
		}
	}

}

func (dt *dbTable) findObjUniques(ctx context.Context, objID string) ([]types.Index, error) {

	filter := map[string]interface{}{
		common.BKObjIDField: objID,
	}
	uniqueIdxs := make([]metadata.ObjectUnique, 0)
	if err := dt.db.Table(common.BKTableNameObjUnique).Find(filter).All(ctx, &uniqueIdxs); err != nil {
		newErr := fmt.Errorf("get obj(%s) logic unique index error. err: %s", objID, err.Error())
		blog.ErrorJSON("%s, rid: %s", newErr.Error(), dt.rid)
		return nil, newErr
	}

	// 返回的数据只有common.BKPropertyIDField, common.BKPropertyTypeField, common.BKFieldID 三个字段
	attrs, err := dt.findObjAttrs(ctx, objID)
	if err != nil {
		return nil, err
	}

	var indexes []types.Index
	for _, idx := range uniqueIdxs {
		newDBIndex, err := index.ToDBUniqueIndex(objID, idx.ID, idx.Keys, attrs)
		if err != nil {
			newErr := fmt.Errorf("obj(%s). %s", objID, err.Error())
			blog.ErrorJSON("%s, rid: %s", newErr.Error(), dt.rid)
			return nil, newErr
		}
		indexes = append(indexes, newDBIndex)

	}

	return indexes, nil
}

// ErrDropIndexNameNotFound TODO
func ErrDropIndexNameNotFound(err error) bool {
	if strings.HasPrefix(err.Error(), "(IndexNotFound) index not found with name") {
		return true
	}
	return false
}
