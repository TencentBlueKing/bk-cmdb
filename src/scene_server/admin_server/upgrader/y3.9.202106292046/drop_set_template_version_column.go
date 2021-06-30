package y3_9_202106292046

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func dropVersionColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	existsVersionFilter := map[string]interface{}{
		"version": map[string]interface{}{
			common.BKDBExists: true,
		},
	}

	count, err := db.Table(common.BKTableNameSetTemplate).Find(existsVersionFilter).Count(ctx)
	if err != nil {
		blog.Errorf("count table %s failed, err: %s", common.BKTableNameSetTemplate, err.Error())
		return err
	}

	for i := uint64(0); i < count; i += common.BKMaxPageSize {
		if err := db.Table(common.BKTableNameSetTemplate).DropColumn(ctx, "version"); err != nil {
			blog.Errorf("drop column failed, field:%s, err:%v", "version", err)
			return err
		}
	}

	return nil
}

func dropSetTplVersionColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	existsVersionFilter := map[string]interface{}{
		common.BKSetTemplateVersionField: map[string]interface{}{
			common.BKDBExists: true,
		},
	}

	count, err := db.Table(common.BKTableNameBaseSet).Find(existsVersionFilter).Count(ctx)
	if err != nil {
		blog.Errorf("count table %s failed, err: %s", common.BKTableNameBaseSet, err.Error())
		return err
	}

	for i := uint64(0); i < count; i += common.BKMaxPageSize {
		if err := db.Table(common.BKTableNameBaseSet).DropColumn(ctx, common.BKSetTemplateVersionField); err != nil {
			blog.Errorf("drop column failed, field:%s, err:%v", common.BKSetTemplateVersionField, err)
			return err
		}
	}

	return nil
}
