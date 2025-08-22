/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// Package migratetenant defines the migrate tenant command
package migratetenant

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

// NewMigrateTenantCommand new tool command for migrating data from old version to multi-tenant version
func NewMigrateTenantCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-tenant",
		Short: "migrate data for multi-tenant version that disables multi-tenant mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	watchConf := new(config.MongoConfig)
	skipRemoveSupplierAccount := false
	inPlaceUpgradeCmd := &cobra.Command{
		Use:   "in-place-upgrade",
		Short: "in-place upgrade old version data to new version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return inPlaceUpgrade(watchConf, skipRemoveSupplierAccount)
		},
	}
	inPlaceUpgradeCmd.Flags().StringVar(&watchConf.MongoURI, "watch-mongo-uri", "",
		"watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb")
	inPlaceUpgradeCmd.Flags().StringVar(&watchConf.MongoRsName, "watch-mongo-rs-name", "rs0", "watch replica set name")
	inPlaceUpgradeCmd.Flags().BoolVar(&skipRemoveSupplierAccount, "skip-remove-supplier-account", false,
		"skip removing all bk_supplier_account fields, upgrade will be faster, but data will have redundant fields")
	cmd.AddCommand(inPlaceUpgradeCmd)

	conf := &copyToNewDBConf{
		oldMongo:      new(config.MongoConfig),
		oldWatchMongo: new(config.MongoConfig),
		newWatchMongo: new(config.MongoConfig),
	}
	copyToNewDBCmd := &cobra.Command{
		Use:   "copy-to-new-db",
		Short: "copy old version data to new db",
		RunE: func(cmd *cobra.Command, args []string) error {
			return copyToNewDB(conf)
		},
	}

	copyToNewDBCmd.Flags().StringVar(&conf.oldMongo.MongoURI, "old-mongo-uri", "",
		"old mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb")
	copyToNewDBCmd.Flags().StringVar(&conf.oldMongo.MongoRsName, "old-mongo-rs-name", "rs0", "old db replica set name")
	copyToNewDBCmd.Flags().StringVar(&conf.oldWatchMongo.MongoURI, "old-watch-mongo-uri", "",
		"old watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb")
	copyToNewDBCmd.Flags().StringVar(&conf.oldWatchMongo.MongoRsName, "old-watch-mongo-rs-name", "rs0",
		"old watch db replica set name")
	copyToNewDBCmd.Flags().StringVar(&conf.newWatchMongo.MongoURI, "watch-mongo-uri", "",
		"new watch mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb")
	copyToNewDBCmd.Flags().StringVar(&conf.newWatchMongo.MongoRsName, "watch-mongo-rs-name", "rs0",
		"new watch db replica set name")
	copyToNewDBCmd.Flags().BoolVar(&conf.isFullSync, "full-sync", false, "is full sync")
	copyToNewDBCmd.Flags().Uint32Var(&conf.startFrom, "start-from", 0, "unix timestamp to start incremental sync from")
	cmd.AddCommand(copyToNewDBCmd)

	return cmd
}

type migrateTenantService struct {
	*migrateTenantDBInfo
	watchDB *migrateTenantDBInfo

	objUUIDMap    map[string]string
	tableHandlers map[string]func(data mapstr.MapStr) (mapstr.MapStr, error)
}

type migrateTenantDBInfo struct {
	oldDB  local.DB
	newDB  local.DB
	sysDB  local.DB
	dbUUID string
}

func newMigrateTenantService(kit *rest.Kit, oldConf, newConf, oldWatchConf, newWatchConf *config.MongoConfig) (
	*migrateTenantService, error) {

	dbInfo, err := newMigrateTenantDBInfo(kit, oldConf, newConf, false)
	if err != nil {
		return nil, err
	}

	// check if current version is the last version of 3.14
	versionCond := mapstr.MapStr{
		"type": "version",
	}
	versionInfo := make(mapstr.MapStr)
	err = dbInfo.oldDB.Table(common.BKTableNameSystem).Find(versionCond).One(context.Background(), &versionInfo)
	if err != nil {
		return nil, fmt.Errorf("get system version info failed, err: %v", err)
	}
	if util.GetStrByInterface(versionInfo["current_version"]) != "y3.14.202502101200" {
		return nil, fmt.Errorf("current version %v is not the last version of 3.14", versionInfo["current_version"])
	}

	watchDBInfo, err := newMigrateTenantDBInfo(kit, oldWatchConf, newWatchConf, true)
	if err != nil {
		return nil, err
	}

	srv := &migrateTenantService{
		migrateTenantDBInfo: dbInfo,
		watchDB:             watchDBInfo,
		objUUIDMap:          make(map[string]string),
		tableHandlers:       make(map[string]func(data mapstr.MapStr) (mapstr.MapStr, error)),
	}

	// init table data handlers for copy to new db operation
	for table := range index.TableIndexes() {
		if common.IsPlatformTableWithTenant(table) {
			srv.tableHandlers[table] = srv.migrateTenantIDField
			continue
		}

		switch table {
		case common.BKTableNameBasePlat, common.BKTableNameBaseHost:
			srv.tableHandlers[table] = srv.migrateCloudIDField
		case common.BKTableNameObjDes:
			srv.tableHandlers[table] = srv.migrateObject
		default:
			srv.tableHandlers[table] = srv.removeSupplierAccount
		}
	}
	srv.tableHandlers[common.BKTableNameHostFavorite] = srv.removeSupplierAccount
	srv.tableHandlers[common.BKTableNameUserAPI] = srv.removeSupplierAccount
	srv.tableHandlers[common.BKTableNameUserCustom] = srv.removeSupplierAccount

	return srv, nil
}

func newMigrateTenantDBInfo(kit *rest.Kit, oldConf, newConf *config.MongoConfig, isWatchDB bool) (*migrateTenantDBInfo,
	error) {

	oldDB, err := local.NewOldMgo(local.MongoConf{
		MaxOpenConns: mongo.MinimumMaxIdleOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          oldConf.MongoURI,
		RsName:       oldConf.MongoRsName,
	}, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client for previous version failed, err: %v", err)
	}

	newMongoConf := local.MongoConf{
		MaxOpenConns: mongo.MinimumMaxIdleOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          newConf.MongoURI,
		RsName:       newConf.MongoRsName,
	}
	var newShardingDB sharding.ShardingDB
	if isWatchDB {
		newShardingDB, err = sharding.NewWatchMongo(newMongoConf, time.Minute, nil)
	} else {
		newShardingDB, err = sharding.NewShardingMongo(newMongoConf, time.Minute, nil)
	}
	if err != nil {
		return nil, err
	}

	newDB, dbUUID, err := logics.GetNewTenantCli(kit, newShardingDB)
	if err != nil {
		return nil, fmt.Errorf("get new tenant db failed, err: %v", err)
	}

	return &migrateTenantDBInfo{
		oldDB:  oldDB,
		newDB:  newDB,
		sysDB:  newShardingDB.Shard(kit.SysShardOpts()),
		dbUUID: dbUUID,
	}, nil
}

func inPlaceUpgrade(watchMongo *config.MongoConfig, skipRemoveSupplierAccount bool) error {
	kit := &rest.Kit{
		Ctx:      util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode),
		TenantID: common.BKSingleTenantID,
	}
	srv, err := newMigrateTenantService(kit, config.Conf.MongoConf, config.Conf.MongoConf, watchMongo, watchMongo)
	if err != nil {
		return err
	}

	if err = srv.inPlaceUpgrade(kit, skipRemoveSupplierAccount); err != nil {
		return err
	}
	return srv.insertTenant(kit)
}

type copyToNewDBConf struct {
	oldMongo      *config.MongoConfig
	oldWatchMongo *config.MongoConfig
	newWatchMongo *config.MongoConfig
	isFullSync    bool
	startFrom     uint32
}

func copyToNewDB(conf *copyToNewDBConf) error {
	kit := &rest.Kit{
		Ctx:      util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode),
		TenantID: common.BKSingleTenantID,
	}
	srv, err := newMigrateTenantService(kit, conf.oldMongo, config.Conf.MongoConf, conf.oldWatchMongo,
		conf.newWatchMongo)
	if err != nil {
		return err
	}

	// init object id to uuid map
	objects := make([]metadata.Object, 0)
	if err = srv.newDB.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField,
		metadata.ModelFieldObjUUID).All(kit.Ctx, &objects); err != nil {
		return fmt.Errorf("get object id and uuid info failed, err: %v", err)
	}

	for _, object := range objects {
		srv.objUUIDMap[object.ObjectID] = object.UUID
	}

	// copy data to new db
	if conf.isFullSync {
		if err = srv.copyFullSyncData(kit); err != nil {
			return err
		}
		return srv.insertTenant(kit)
	}

	return srv.copyIncrSyncData(kit, conf.startFrom)
}
