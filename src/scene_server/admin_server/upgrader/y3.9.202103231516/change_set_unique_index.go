package y3_9_202103231516

import (
	"configcenter/src/storage/dal/types"
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)
var (
	sortFlag      = int32(1)
	idUniqueIndex = types.Index{
		Keys:       map[string]int32{common.BKFieldID: sortFlag},
		Unique:     true,
		Background: true,
		Name:       "idx_unique_id",
	}
)

func changeSetUniqueIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	idxUniqueParentIDSetName := types.Index{
		Keys:       map[string]int32{common.BKParentIDField: sortFlag, common.BKSetNameField: sortFlag},
		Unique:     true,
		Background: true,
		Name:       "idx_unique_parentID_setName",
	}
	tableName := common.BKTableNameBaseSet
	dbIndexes, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("find table(%s) index error. err: %s", tableName, err.Error())
		return err
	}
	idxUniqueSetName := "idx_unique_bizID_setName"
	for _, index := range dbIndexes {

		if index.Keys == nil {
			continue
		}
		if index.Name == idxUniqueSetName {
			if err := db.Table(tableName).DropIndex(ctx, idxUniqueSetName); err != nil {
				blog.ErrorJSON("drop table(%s) index error. idx name: %s, err: %s",
					tableName, idxUniqueSetName, err.Error())
				return err
			}
			if err := db.Table(tableName).CreateIndex(ctx, idxUniqueParentIDSetName); err != nil {
				blog.ErrorJSON("create table(%s) index error. idx name: %s, index: %s, err: %s",
					tableName, idxUniqueSetName, idxUniqueParentIDSetName, err.Error())
				return err
			}
		}
	}
	return nil
}