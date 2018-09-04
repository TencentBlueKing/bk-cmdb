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

package v3v0v1alpha2

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
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
