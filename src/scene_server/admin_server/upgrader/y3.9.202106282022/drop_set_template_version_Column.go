package y3_9_202106282022

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func dropSetTplVersionColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := db.Table(common.BKTableNameSetTemplate).DropColumn(ctx, "version"); err != nil {
		blog.Errorf("update failed, field:%s, err:%v", "version", err)
		return err
	}

	if err := db.Table(common.BKTableNameBaseSet).DropColumn(ctx, common.BKSetTemplateVersionField); err != nil {
		blog.Errorf("update failed, field:%s, err:%v", common.BKSetTemplateVersionField, err)
		return err
	}

	return nil
}
