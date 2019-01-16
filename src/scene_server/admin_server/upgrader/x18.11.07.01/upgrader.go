package x18_11_07_01

import (
	"context"
	"gopkg.in/mgo.v2"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addCloudTaskTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameCloudTask
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_task_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"bk_task_name": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil
}

func addCloudResourceConfirmTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameCloudResourceConfirm
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_resource_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"bk_obj_id": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil

}

func addCloudSyncHistoryTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameCloudSyncHistory
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_task_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"bk_history_id": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil

}

func addCloudConfirmHistoryTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameResourceConfirmHistory
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_resource_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"confirm_history_id": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil
}
