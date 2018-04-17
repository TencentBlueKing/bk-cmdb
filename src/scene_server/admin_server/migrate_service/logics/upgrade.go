package logics

import (
	"configcenter/src/common"
	dbStorage "configcenter/src/storage"
)

func Upgrade(instData dbStorage.DI) error {
	condition := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDApp,
		common.BKPropertyIDField: map[string]interface{}{
			"$in": []string{
				"time_zone",
				"language",
			},
		},
	}
	data := map[string]interface{}{
		"isrequired": true,
	}
	instData.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)
	return nil
}
