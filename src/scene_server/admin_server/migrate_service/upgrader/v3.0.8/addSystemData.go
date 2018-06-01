package v3v0v8

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/upgrader"
	"configcenter/src/storage"
)

func addSystemData(db storage.DI, conf *upgrader.Config) error {
	tablename := "cc_System"
	blog.Errorf("add data for  %s table ", tablename)
	data := map[string]interface{}{
		common.HostCrossBizField: common.HostCrossBizValue}
	isExist, err := db.GetCntByCondition(tablename, data)
	if nil != err {
		blog.Errorf("add data for  %s table error  %s", tablename, err)
		return err
	}
	if isExist > 0 {
		return nil
	}
	_, err = db.Insert(tablename, data)
	if nil != err {
		blog.Errorf("add data for  %s table error  %s", tablename, err)
		return err
	}

	return nil
}
