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

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	idx "configcenter/src/common/index"
	"configcenter/src/common/metadata"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/dal/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"
)

/*
 如何展示错误给用户
*/

// DBSync do sync db table index background task
func DBSync(e *backbone.Engine, db dal.Dal, options options.Config) {
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

	defaultDBTable dal.Dal
)

type shardingDBTable struct {
	db                         dal.Dal
	preCleanRedundancyTableMap map[string]map[string]struct{}
	rid                        string
}

type dbTable struct {
	db                         dal.RDB
	preCleanRedundancyTableMap map[string]struct{}
	rid                        string
	tenantID                   string
}

// RunSyncDBTableIndex run sync db table index task
func RunSyncDBTableIndex(ctx context.Context, e *backbone.Engine, db dal.Dal, options options.Config) {
	rid := util.GenerateRID()
	// do not sync db table index if this admin-server is for ci,
	// because suite test contains clear database operations that collides with this task
	if version.CCRunMode == version.CCRunModeForCI {
		blog.Infof("run mode is for ci, skip sync db table index task, rid: %s", rid)
		return
	}

	for dbReady := false; !dbReady; {
		// 等待数据库初始化
		var err error
		dbReady, err = upgrader.DBReady(ctx, db)
		if err != nil {
			blog.Errorf("check whether the db initialization is complete failed, err: %v, rid: %s", err, rid)
		}
		if !dbReady {
			blog.Errorf("db initialization is not complete, rid: %s", rid)
		}
		time.Sleep(20 * time.Second)
	}
	syncWorker := func(st *shardingDBTable, isTable bool) {
		for {
			rid := util.GenerateRID()
			st.rid = rid
			blog.Infof("start sync table or index worker rid: %s", rid)

			if !e.ServiceManageInterface.IsMaster() {
				blog.Infof("skip sync table or index worker. reason: not master. rid: %s", rid)
				time.Sleep(5 * time.Second)
				continue
			}

			if isTable {
				preCleanRedundancyTableMap := make(map[string]map[string]struct{})
				err := tenant.ExecForAllTenants(func(tenantID string) error {
					tenantRid := fmt.Sprintf("%s-%s", tenantID, rid)
					blog.Infof("start object sharding table rid: %s", tenantRid)

					// 先处理模型实例和关联关系表
					dt := &dbTable{
						db:                         st.db.Shard(sharding.NewShardOpts().WithTenant(tenantID)),
						preCleanRedundancyTableMap: st.preCleanRedundancyTableMap[tenantID],
						rid:                        tenantRid,
						tenantID:                   tenantID,
					}
					if err := dt.syncModelShardingTable(ctx); err != nil {
						blog.Errorf("model table sync failed, err: %v, rid: %s", err, tenantRid)
						return err
					}
					preCleanRedundancyTableMap[tenantID] = dt.preCleanRedundancyTableMap

					blog.Infof("end sync table rid: %s", tenantRid)
					return nil
				})
				if err != nil {
					blog.Errorf("model table sync failed, err: %v, rid: %s", err, st.rid)
					time.Sleep(time.Second * time.Duration(options.ShardingTable.TableInterval))
					continue
				}
				st.preCleanRedundancyTableMap = preCleanRedundancyTableMap
				time.Sleep(time.Second * time.Duration(options.ShardingTable.TableInterval))
				continue
			}

			blog.Infof("start table common index rid: %s", rid)
			if err := st.syncIndexes(ctx); err != nil {
				blog.Errorf("model table sync failed, err: %v, rid: %s", err, st.rid)
				time.Sleep(time.Minute * time.Duration(options.ShardingTable.IndexesInterval))
				continue
			}
			blog.Infof("end sync table index rid: %s", rid)
			time.Sleep(time.Minute * time.Duration(options.ShardingTable.IndexesInterval))
		}
	}

	dtTable := &shardingDBTable{db: db, rid: rid, preCleanRedundancyTableMap: make(map[string]map[string]struct{})}
	go syncWorker(dtTable, true)
	dtIndex := &shardingDBTable{db: db, rid: rid}
	go syncWorker(dtIndex, false)
}

// RunSyncDBIndex run sync db index task
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
		blog.Errorf("check whether the db initialization is complete failed, err: %v, rid: %s", err, rid)
		return ccErr.CCError(common.CCErrCommDBSelectFailed)
	}
	if !dbReady {
		blog.Errorf("db not initialization is complete, rid: %s", rid)
		return ccErr.CCError(common.CCErrCommDBSelectFailed)

	}

	dt := &shardingDBTable{db: defaultDBTable, rid: rid}
	blog.Infof("start table common index rid: %s", rid)
	if err := dt.syncIndexes(ctx); err != nil {
		blog.Errorf("model table sync failed, err: %v, rid: %s", err, dt.rid)
	}
	blog.Infof("end sync table index rid: %s", rid)

	return nil
}

// syncIndexes 同步表中定义的索引
func (st *shardingDBTable) syncIndexes(ctx context.Context) error {
	if err := st.syncDBTableIndexes(ctx); err != nil {
		blog.Warnf("sync table index to db failed, err: %v, rid: %s", err, st.rid)
		// 不影响后需任务
	}

	return nil
}

func (st *shardingDBTable) syncDBTableIndexes(ctx context.Context) error {
	deprecatedIndexNames := index.DeprecatedIndexName()
	tableIndexes := index.TableIndexes()

	platDBTable := &dbTable{db: st.db.Shard(sharding.NewShardOpts().WithIgnoreTenant()), rid: "platform-" + st.rid}
	for _, tableName := range common.PlatformTables() {
		indexes, exists := tableIndexes[tableName]
		if !exists {
			continue
		}

		blog.Infof("start sync table(%s) index, rid: %s", tableName, platDBTable.rid)

		if common.IsPlatformTableWithTenant(tableName) {
			err := st.db.ExecForAllDB(func(db local.DB) error {
				dbTableInst := &dbTable{db: db, rid: "platform-with-tenant-" + st.rid}
				return dbTableInst.syncIndexesToDB(ctx, tableName, indexes, deprecatedIndexNames[tableName])
			})
			if err != nil {
				blog.Warnf("sync table (%s) index failed. err: %v, rid: %s", tableName, err, platDBTable.rid)
				continue
			}
		}

		if err := platDBTable.syncIndexesToDB(ctx, tableName, indexes, deprecatedIndexNames[tableName]); err != nil {
			blog.Warnf("sync table (%s) index failed. err: %v, rid: %s", tableName, err, platDBTable.rid)
			continue
		}
	}

	err := tenant.ExecForAllTenants(func(tenantID string) error {
		dt := &dbTable{
			db:       st.db.Shard(sharding.NewShardOpts().WithTenant(tenantID)),
			rid:      tenantID + "-" + st.rid,
			tenantID: tenantID,
		}

		dtIndexesMap, err := dt.findSyncIndexesLogicUnique(ctx)
		if err != nil {
			blog.Errorf("find db logic unique failed, err: %v, rid: %s", err, dt.rid)
			return err
		}

		for tableName, indexes := range tableIndexes {
			if common.IsPlatformTable(tableName) {
				continue
			}

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
	})
	if err != nil {
		blog.Errorf("sync db table indexes failed, err: %v, rid: %s", err, st.rid)
		return err
	}

	return nil
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
			sameIndex := false
			for _, dbIndex := range dbIndexList {
				if index.IndexEqual(dbIndex, logicIndex) {
					sameIndex = true
					break
				}
			}
			if sameIndex {
				continue
			}

			if err := dt.db.Table(tableName).CreateIndex(ctx, logicIndex); err != nil {
				blog.Errorf("create table(%s) index(%s) error. err: %s, rid: %s",
					tableName, logicIndex.Name, err.Error(), dt.rid)
				return err
			}
		}
	}

	return nil

}

// countQuotedWithDestModel count quoted object id by dest_model.
func (dt *dbTable) countQuotedWithDestModel(ctx context.Context, destModel string) (uint64, error) {

	cond := map[string]interface{}{
		common.BKDestModelField: destModel,
	}
	count, err := dt.db.Table(common.BKTableNameModelQuoteRelation).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("count quoted object failed, cond: %+v, err: %v, rid: %s", cond, err, dt.rid)
		return 0, err
	}

	return count, nil
}

func (dt *dbTable) findSyncIndexesLogicUnique(ctx context.Context) (map[string][]types.Index, error) {
	objs := make([]metadata.Object, 0)
	if err := dt.db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField, common.BKIsPre).All(ctx,
		&objs); err != nil {
		blog.Errorf("get all common object id  error. err: %s, rid: %s", err.Error(), dt.rid)
		return nil, err
	}

	tbIndexes := make(map[string][]types.Index)
	for _, obj := range objs {
		blog.Infof("start object(%s) sharding table rid: %s", obj.ObjectID, dt.rid)

		instTable := common.GetObjectInstTableName(obj.ObjectID, dt.tenantID)
		instAsstTable := common.GetObjectInstAsstTableName(obj.ObjectID, dt.tenantID)

		uniques, err := dt.findObjUniques(ctx, obj.ObjectID)
		if err != nil {
			blog.Errorf("object(%s) logic unique to db index error. err: %s, rid: %s",
				obj.ObjectID, err.Error(), dt.rid)
			return nil, err
		}
		// 内置模型不需要简表
		if !obj.IsPre {
			tbIndexes[instTable] = append(index.InstanceIndexes(), uniques...)
			// if there is a table instance table, an index needs to be created.
			count, err := dt.countQuotedWithDestModel(ctx, obj.ObjectID)
			if err != nil {
				return nil, err
			}
			if count > 0 {
				tbIndexes[instTable] = append(index.TableInstanceIndexes(), uniques...)
			}

		} else {
			tb := ""
			switch obj.ObjectID {
			case common.BKInnerObjIDHost:
				tb = common.BKTableNameBaseHost
			case common.BKInnerObjIDBizSet:
				tb = common.BKTableNameBaseBizSet
			case common.BKInnerObjIDApp:
				tb = common.BKTableNameBaseApp
			case common.BKInnerObjIDProject:
				tb = common.BKTableNameBaseProject
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
		if common.IsObjectInstAsstShardingTable(name) {
			modelDBTableNameMap[name] = struct{}{}
		}
	}

	objs := make([]metadata.Object, 0)
	if err := dt.db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField, common.BKIsPre).All(ctx,
		&objs); err != nil {
		blog.Errorf("get all common object id  error. err: %s, rid: %s", err.Error(), dt.rid)
		return err
	}

	for _, obj := range objs {
		blog.Infof("start object(%s) sharding table rid: %s", obj.ObjectID, dt.rid)

		instTable := common.GetObjectInstTableName(obj.ObjectID, dt.tenantID)
		instAsstTable := common.GetObjectInstAsstTableName(obj.ObjectID, dt.tenantID)

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
		// if there is a table instance table, an index needs to be created.
		count, err := dt.countQuotedWithDestModel(ctx, obj.ObjectID)
		if err != nil {
			return err
		}
		if count > 0 {
			objIndexes = append(index.TableInstanceIndexes(), uniques...)
		}

		dt.createTable(ctx, obj, modelDBTableNameMap, instTable, objIndexes)

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

func (dt *dbTable) createTable(ctx context.Context, obj metadata.Object, modelDBTableNameMap map[string]struct{},
	instTable string, objIndexes []types.Index) {

	// 内置模型不需要建表
	if obj.IsPre {
		return
	}

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
	return
}

func (dt *dbTable) cleanRedundancyTable(ctx context.Context, modelDBTableNameMap map[string]struct{}) error {
	objs := make([]metadata.Object, 0)
	if err := dt.db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField, common.BKIsPre).All(ctx,
		&objs); err != nil {
		blog.Errorf("get all common object id  error. err: %s, rid: %s", err.Error(), dt.rid)
		// NOTICE: 错误直接忽略不行后需功能
		return err
	}

	// 再次确认数据，保证存在模型的的表不被删除
	for _, obj := range objs {
		instTable := common.GetObjectInstTableName(obj.ObjectID, dt.tenantID)
		instAsstTable := common.GetObjectInstAsstTableName(obj.ObjectID, dt.tenantID)
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
			blog.WarnJSON("create table(%s) error. index: %s, err: %v, rid: %s", tableName, index, err, dt.rid)
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
		newErr := fmt.Errorf("get obj(%s) logic unique index error. err: %v", objID, err)
		blog.Errorf("%v, rid: %s", newErr, dt.rid)
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
			blog.Errorf("%s, rid: %s", newErr.Error(), dt.rid)
			return nil, newErr
		}
		indexes = append(indexes, newDBIndex)

	}

	return indexes, nil
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
		newErr := fmt.Errorf("get obj(%s) property failed, err: %v", objID, err)
		blog.Errorf("%v, rid: %s", newErr, dt.rid)
		return nil, newErr
	}

	return attrs, nil
}

// ErrDropIndexNameNotFound TODO
func ErrDropIndexNameNotFound(err error) bool {
	if strings.HasPrefix(err.Error(), "(IndexNotFound) index not found with name") {
		return true
	}
	return false
}

// CreateIndexes create table index
func CreateIndexes(kit *rest.Kit, db local.DB, table string, indexes []types.Index) error {
	existIndexes, err := db.Table(table).Indexes(kit.Ctx)
	if err != nil {
		blog.Errorf("get %s table exist index failed, err: %v", table, err)
		return err
	}

	existIndexMap := make(map[string]struct{})
	for _, index := range existIndexes {
		existIndexMap[index.Name] = struct{}{}
	}

	insertIndexes := make([]types.Index, 0)
	for _, index := range indexes {
		if _, exist := existIndexMap[index.Name]; exist {
			continue
		}

		isExist := false
		for _, existIndex := range existIndexes {
			if idx.IndexEqual(existIndex, index) {
				isExist = true
				break
			}
		}
		if isExist {
			continue
		}
		insertIndexes = append(insertIndexes, index)
	}

	if len(insertIndexes) == 0 {
		blog.Infof("table %s index is up to date", table)
		return nil
	}
	if err = db.Table(table).BatchCreateIndexes(kit.Ctx, insertIndexes); err != nil {
		blog.Errorf("create %s table index %+v failed, err: %v", table, insertIndexes, err)
		return err
	}

	return nil
}
