package y3_9_202103231516

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.9.202103231516", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute y3.9.202103231516")

	err = changeSetUniqueIndex(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.9.202103231516] change unique index failed, err: %v", err)
		return err
	}
	return nil
}