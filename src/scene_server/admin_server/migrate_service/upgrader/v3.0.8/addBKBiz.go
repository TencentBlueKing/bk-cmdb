package v3v0v8

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/upgrader"
	"configcenter/src/storage"
	"time"
)

func addBKBiz(db storage.DI, conf *upgrader.Config) error {
	return nil
}

func createBiz(db storage.DI, data map[string]interface{}) error {
	tablename := "cc_ApplicationBase"

	return storage.Upsert(db, tablename, data, []string{common.BKOwnerIDField, common.BKDefaultField}, []string{})
}
