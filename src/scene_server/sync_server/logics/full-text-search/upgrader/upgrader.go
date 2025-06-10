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

// Package upgrader defines the upgrade logics of full-text-search sync
package upgrader

import (
	"context"
	"sort"
	"sync"

	"configcenter/pkg/tenant"
	ftypes "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"

	"github.com/olivere/elastic/v7"
)

// upgraderInst is an instance of upgrader
var upgraderInst = &upgrader{
	upgraderPool: make(map[int]UpgraderFunc),
	registerLock: sync.Mutex{},
}

// InitUpgrader initialize global upgrader
func InitUpgrader(esCli *elastic.Client, syncer types.SyncDataI, indexSetting metadata.ESIndexMetaSettings) {
	upgraderInst.esCli = esCli
	upgraderInst.syncer = syncer
	upgraderInst.indexSetting = indexSetting
	upgraderInst.initCurrentEsIndex()
}

// upgrader is the full-text-search sync upgrade structure
type upgrader struct {
	// esCli is the elasticsearch client
	esCli *elastic.Client
	// syncer is the full-text-search data syncer
	syncer types.SyncDataI
	// indexSetting is the es index meta setting
	indexSetting metadata.ESIndexMetaSettings
	// upgraderPool is the mapping of all upgrader version -> upgrader function
	upgraderPool map[int]UpgraderFunc
	// registerLock is the lock for registering upgrader function to avoid conflict
	registerLock sync.Mutex
	// currentIndexNameMap stores the current index names, is used to check whether the index is current
	currentIndexNameMap map[string]struct{}
}

// UpgraderFunc is upgrader function definition
// NOTE: do not need to add new index, only update/remove old index and migrate data
type UpgraderFunc func(ctx context.Context, rid string) (*UpgraderFuncResult, error)

// UpgraderFuncResult is upgrader function return result
type UpgraderFuncResult struct {
	// Indexes is all indexes in this version of upgrader
	Indexes []string
	// ReindexInfo is the reindex info of the pre version index to new version index
	ReindexInfo map[string]string
	// NeedSyncAll defines whether we need to sync all data
	NeedSyncAll bool
}

// RegisterUpgrader register upgrader
func RegisterUpgrader(version int, handler UpgraderFunc) {
	upgraderInst.registerLock.Lock()
	defer upgraderInst.registerLock.Unlock()

	upgraderInst.upgraderPool[version] = handler
}

// Upgrade es index to the newest version
func Upgrade(ctx context.Context, rid string) (*ftypes.MigrateResult, error) {
	// compare version to get the needed upgraders
	dbVersion, versions, result, err := compareVersions(ctx, rid)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return result, nil
	}

	// add current version indexes first
	newIndexMap, err := upgraderInst.createCurrentEsIndex(ctx, rid)
	if err != nil {
		return nil, err
	}

	delIndexMap := make(map[string]struct{})
	reIndexInfo := make(map[string]string)
	// sync all data if it is the first time to upgrade
	needSyncAll := dbVersion.CurrentVersion == 0

	// do all the upgrader
	for _, version := range versions {
		upgraderFunc := upgraderInst.upgraderPool[version]
		res, err := upgraderFunc(ctx, rid)
		if err != nil {
			blog.Errorf("upgrade full-text search sync failed, version: %d, err: %v, rid: %s", version, err, rid)
			return nil, err
		}

		for _, index := range res.Indexes {
			_, exists := upgraderInst.currentIndexNameMap[index]
			if !exists {
				delIndexMap[index] = struct{}{}
			}
		}

		if res.NeedSyncAll {
			needSyncAll = true
		}

		for oldIdx, newIdx := range res.ReindexInfo {
			reIndexInfo[oldIdx] = newIdx
			delete(newIndexMap, newIdx)
		}

		dbVersion.CurrentVersion = version
	}

	if err = migrateIndexData(ctx, needSyncAll, reIndexInfo, newIndexMap, rid); err != nil {
		return nil, err
	}

	// delete all old indexes
	if err = upgraderInst.deleteIndex(ctx, delIndexMap, rid); err != nil {
		return nil, err
	}

	if err = updateVersion(ctx, dbVersion, rid); err != nil {
		return nil, err
	}

	return result, nil
}

func migrateIndexData(ctx context.Context, needSyncAll bool, reIndexInfo map[string]string,
	newIndexMap map[string]struct{}, rid string) error {

	syncIndexes := make([]string, 0)
	if needSyncAll {
		// sync all index data
		syncIndexes = types.AllIndexNames
	} else {
		if len(reIndexInfo) > 0 {
			// TODO complete the reindex logics in next version that needs reindex
		}

		// need sync new indexes that have no reindex data
		for index := range newIndexMap {
			syncIndexes = append(syncIndexes, index)
		}
	}

	// sync all data in the newly created index
	for _, index := range syncIndexes {
		syncOpts := &ftypes.SyncDataOption{
			Index: index,
		}

		err := tenant.ExecForAllTenants(func(tenantID string) error {
			kit := rest.NewKit().WithCtx(ctx).WithTenant(tenantID).WithRid(rid)
			if err := upgraderInst.syncer.SyncData(kit, syncOpts); err != nil {
				blog.Errorf("sync data by index %s after migration failed, err: %v, rid: %s", index, err, rid)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func compareVersions(ctx context.Context, rid string) (*Version, []int, *ftypes.MigrateResult, error) {
	dbVersion, err := getVersion(ctx, rid)
	if err != nil {
		return nil, nil, nil, err
	}

	result := &ftypes.MigrateResult{
		PreVersion:       dbVersion.CurrentVersion,
		CurrentVersion:   dbVersion.CurrentVersion,
		FinishedVersions: make([]int, 0),
	}

	var versions []int
	for version := range upgraderInst.upgraderPool {
		if version > dbVersion.CurrentVersion {
			versions = append(versions, version)
		}
	}

	if len(versions) == 0 {
		return dbVersion, versions, result, nil
	}

	dbVersion.InitVersion = dbVersion.CurrentVersion
	sort.Ints(versions)
	return dbVersion, versions, result, nil
}

// fullTextVersion is the full-text search sync version type
const fullTextVersion = "full_text_search_version"

// Version is the full-text search sync version info
type Version struct {
	Type           string `bson:"type"`
	CurrentVersion int    `bson:"current_version"`
	InitVersion    int    `bson:"init_version"`
}

// getVersion get full-text search sync version info from db
func getVersion(ctx context.Context, rid string) (*Version, error) {
	condition := map[string]interface{}{
		"type": fullTextVersion,
	}

	data := new(Version)
	db := mongodb.Shard(sharding.NewShardOpts().WithIgnoreTenant())
	err := db.Table(common.BKTableNameSystem).Find(condition).One(ctx, &data)
	if err != nil {
		if !mongodb.IsNotFoundError(err) {
			blog.Errorf("get full-text search sync version failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		data.Type = fullTextVersion

		err = db.Table(common.BKTableNameSystem).Insert(ctx, data)
		if err != nil {
			blog.Errorf("insert full-text search sync version failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
		return data, nil
	}

	return data, nil
}

// updateVersion update full-text search sync version info to db
func updateVersion(ctx context.Context, version *Version, rid string) error {
	condition := map[string]interface{}{
		"type": fullTextVersion,
	}

	err := mongodb.Shard(sharding.NewShardOpts().WithIgnoreTenant()).Table(common.BKTableNameSystem).
		Update(ctx, condition, version)
	if err != nil {
		blog.Errorf("update full-text search sync version %+v failed, err: %v, rid: %s", version, err, rid)
		return err
	}

	return nil
}
