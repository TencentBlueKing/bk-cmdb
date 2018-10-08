/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package x08_09_04_01

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/scene_server/validator"
	"configcenter/src/storage"
)

func updateSystemProperty(db storage.DI, conf *upgrader.Config) (err error) {
	objs := []metadata.ObjectDes{}
	condition := map[string]interface{}{
		"bk_classification_id": "bk_biz_topo",
	}
	err = db.GetMutilByCondition(common.BKTableNameObjDes, nil, condition, &objs, "", 0, 0)
	if err != nil {
		return err
	}

	objIDs := []string{}
	for _, obj := range objs {
		objIDs = append(objIDs, obj.ObjectID)
	}

	tablename := common.BKTableNameObjAttDes
	condition = map[string]interface{}{
		"bk_property_id": map[string]interface{}{"$in": []string{common.BKChildStr, common.BKInstParentStr}},
		"bk_obj_id":      map[string]interface{}{"$in": objIDs},
	}
	data := map[string]interface{}{
		"bk_issystem": true,
	}

	return db.UpdateByCondition(tablename, data, condition)
}

func updateIcon(db storage.DI, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		"bk_obj_id": common.BKInnerObjIDTomcat,
	}
	data := map[string]interface{}{
		"bk_obj_icon": "icon-cc-tomcat",
	}
	err = db.UpdateByCondition(common.BKTableNameObjDes, data, condition)
	if err != nil {
		return err
	}
	condition = map[string]interface{}{
		"bk_obj_id": common.BKInnerObjIDApache,
	}
	data = map[string]interface{}{
		"bk_obj_icon": "icon-cc-apache",
	}

	return db.UpdateByCondition(common.BKTableNameObjDes, data, condition)
}

func fixesProcess(db storage.DI, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: map[string]interface{}{"$in": []string{"priority", "proc_num", "auto_time_gap", "timeout"}},
	}
	data := map[string]interface{}{
		"option": validator.IntOption{Min: "1", Max: "10000"},
	}
	err = db.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)
	if nil != err {
		return err
	}
	return nil
}
