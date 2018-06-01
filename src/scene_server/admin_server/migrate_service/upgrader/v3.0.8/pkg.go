package v3v0v8

import (
	"configcenter/src/scene_server/admin_server/migrate_service/upgrader"
	"configcenter/src/storage"
)

func init() {
	upgrader.RegistUpgrader("v3.0.8", upgrade)
}

func upgrade(db storage.DI, conf *upgrader.Config) (err error) {
	err = createTable(db, conf)
	if err != nil {
		return err
	}
	err = addPresetObjects(db, conf)
	if err != nil {
		return err
	}
	err = addPlatData(db, conf)
	if err != nil {
		return err
	}
	err = addSystemData(db, conf)
	if err != nil {
		return err
	}

	return
}
