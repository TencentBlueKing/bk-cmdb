// Package upgrader defines the upgrade logics of full-text-search sync
package upgrader

import (
	"context"
	"sort"
	"sync"

	ftypes "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
	"configcenter/src/storage/driver/mongodb"

	"github.com/olivere/elastic/v7"
)

// upgraderInst is an instance of upgrader
var upgraderInst = &upgrader{
	upgraderPool: make(map[int]UpgraderFunc),
	registerLock: sync.Mutex{},
}

// InitUpgrader initialize global upgrader
func InitUpgrader(esCli *elastic.Client, indexSetting metadata.ESIndexMetaSettings) {
	upgraderInst.esCli = esCli
	upgraderInst.indexSetting = indexSetting
	upgraderInst.initCurrentEsIndex()
}

// upgrader is the full-text-search sync upgrade structure
type upgrader struct {
	// esCli is the elasticsearch client
	esCli *elastic.Client
	// indexSetting is the es index meta setting
	indexSetting metadata.ESIndexMetaSettings
	// upgraderPool is the mapping of all upgrader version -> upgrader function
	upgraderPool map[int]UpgraderFunc
	// registerLock is the lock for registering upgrader function to avoid conflict
	registerLock sync.Mutex
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
}

// RegisterUpgrader register upgrader
func RegisterUpgrader(version int, handler UpgraderFunc) {
	upgraderInst.registerLock.Lock()
	defer upgraderInst.registerLock.Unlock()

	upgraderInst.upgraderPool[version] = handler
}

// Upgrade es index to the newest version
func Upgrade(ctx context.Context, rid string) (*ftypes.MigrateResult, []string, error) {
	// compare version to get the needed upgraders
	dbVersion, versions, result, err := compareVersions(ctx, rid)
	if err != nil {
		return nil, nil, err
	}

	if len(versions) == 0 {
		return result, nil, nil
	}

	// add current version indexes first
	newIndexMap, err := upgraderInst.createCurrentEsIndex(ctx, rid)
	if err != nil {
		return nil, nil, err
	}

	currentIndexMap := make(map[string]struct{})
	for _, indexes := range types.IndexMap {
		for _, index := range indexes {
			currentIndexMap[index.Name()] = struct{}{}
		}
	}

	delIndexMap := make(map[string]struct{})
	reIndexInfo := make(map[string]string)

	// do all the upgrader
	for _, version := range versions {
		upgraderFunc := upgraderInst.upgraderPool[version]
		res, err := upgraderFunc(ctx, rid)
		if err != nil {
			blog.Errorf("upgrade full-text search sync failed, version: %d, err: %v, rid: %s", version, err, rid)
			return nil, nil, err
		}

		for _, index := range res.Indexes {
			_, exists := currentIndexMap[index]
			if !exists {
				delIndexMap[index] = struct{}{}
			}
		}

		for oldIdx, newIdx := range res.ReindexInfo {
			reIndexInfo[oldIdx] = newIdx
			delete(newIndexMap, newIdx)
		}

		dbVersion.CurrentVersion = version
		if err = updateVersion(ctx, dbVersion, rid); err != nil {
			return nil, nil, err
		}
	}

	// TODO complete these logics in next version, right now there's only one version
	// TODO delete all old indexes
	// TODO reindex all data

	// returns all new indexes that requires data sync
	syncIndexes := make([]string, 0)
	for index := range newIndexMap {
		syncIndexes = append(syncIndexes, index)
	}

	return result, syncIndexes, nil
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
		return nil, versions, result, nil
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
	err := mongodb.Client().Table(common.BKTableNameSystem).Find(condition).One(ctx, &data)
	if err != nil {
		if !mongodb.Client().IsNotFoundError(err) {
			blog.Errorf("get full-text search sync version failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		data.Type = fullTextVersion

		err = mongodb.Client().Table(common.BKTableNameSystem).Insert(ctx, data)
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

	err := mongodb.Client().Table(common.BKTableNameSystem).Update(ctx, condition, version)
	if err != nil {
		blog.Errorf("update full-text search sync version %+v failed, err: %v, rid: %s", version, err, rid)
		return err
	}

	return nil
}
