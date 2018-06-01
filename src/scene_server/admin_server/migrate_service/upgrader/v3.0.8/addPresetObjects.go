package v3v0v8

import (
	"configcenter/src/scene_server/admin_server/migrate_service/upgrader"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/storage"
)

func addPresetObjects(db storage.DI, conf *upgrader.Config) (err error) {
	for tablename, indexs := range classificationRows {
		if err = db.CreateTable(tablename); err != nil {
			return err
		}
		for index := range indexs {
			if err = db.Index(tablename, &indexs[index]); err != nil {
				return err
			}
		}
	}
	return nil
}

func addClassifications(db storage.DI, conf *upgrader.Config) (err error) {

}

var classificationRows = []*metadata.ObjClassification{
	&metadata.ObjClassification{ClassificationID: "bk_host_manage", ClassificationName: "主机管理", ClassificationType: "inner", ClassificationIcon: "icon-cc-business"},
	&metadata.ObjClassification{ClassificationID: "bk_biz_topo", ClassificationName: "业务拓扑", ClassificationType: "inner", ClassificationIcon: "icon-cc-square"},
	&metadata.ObjClassification{ClassificationID: "bk_organization", ClassificationName: "组织架构", ClassificationType: "inner", ClassificationIcon: "icon-cc-free-pool"},
	&metadata.ObjClassification{ClassificationID: "bk_network", ClassificationName: "网络", ClassificationType: "inner", ClassificationIcon: "icon-cc-networks"},
	&metadata.ObjClassification{ClassificationID: "bk_middleware", ClassificationName: "中间件", ClassificationType: "inner", ClassificationIcon: "icon-cc-record"},
}
