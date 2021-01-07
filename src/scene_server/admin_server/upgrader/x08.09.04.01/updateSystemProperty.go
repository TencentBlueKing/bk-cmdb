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
	"context"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updateSystemProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	objs := []metadata.Object{}
	condition := map[string]interface{}{
		"bk_classification_id": "bk_biz_topo",
	}
	err = db.Table(common.BKTableNameObjDes).Find(condition).All(ctx, &objs)
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

	return db.Table(tablename).Update(ctx, condition, data)
}

func fixesProcess(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: map[string]interface{}{"$in": []string{"priority", "proc_num", "auto_time_gap", "timeout"}},
	}
	data := map[string]interface{}{
		"option": metadata.IntOption{Min: "1", Max: "10000"},
	}
	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, condition, data)
	if nil != err {
		return err
	}
	return nil
}
