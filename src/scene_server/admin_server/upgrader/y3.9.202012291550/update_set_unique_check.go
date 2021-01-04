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

package y3_9_202012291550

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// updateSetUniqueCheck update the document about set-unique-check in 'cc_ObjectUnique'
func updateSetUniqueCheck(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	oldAttributes := []metadata.Attribute{}
	cond := map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"bk_obj_id":      "set",
				"bk_property_id": "bk_biz_id",
			},
			{
				"bk_obj_id":      "set",
				"bk_property_id": "bk_set_name",
			},
		},
	}
	err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &oldAttributes)
	if err != nil {
		return err
	}
	if len(oldAttributes) != 2 {
		blog.Errorf("find object in cc_ObjAttDes, but got unexpected result number which should be 2, condition:%#v, result:%v", cond, oldAttributes)
		return fmt.Errorf("find object in cc_ObjAttDes, but got unexpected result number which should be 2, condition:%#v, result:%v", cond, oldAttributes)
	}

	filter := map[string]interface{}{
		"bk_obj_id":  "set",
		"ispre":      true,
		"must_check": true,
	}

	uniques := make([]metadata.ObjectUnique, 0)
	// locate the set-unique-check document in 'cc_ObjectUnique'
	err = db.Table(common.BKTableNameObjUnique).Find(filter).All(ctx, &uniques)
	if err != nil {
		blog.Errorf("find object in cc_ObjectUnique failed, filter:%#v, err:%v", filter, err)
		return err
	}
	if len(uniques) != 1 {
		blog.Errorf("find object in cc_ObjectUnique ,but got unexpected result number which should be single, filter:%#v, result:%v", filter, uniques)
		return fmt.Errorf("find object in cc_ObjectUnique ,but got unexpected result number which should be single, filter:%#v, result:%v", filter, uniques)
	}

	cond = map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"key_kind": "property",
				"key_id":   oldAttributes[0].ID,
			},
			{
				"key_kind": "property",
				"key_id":   oldAttributes[1].ID,
			},
		},
	}
	// update the set-unique-check document in 'cc_ObjectUnique'
	err = db.Table(common.BKTableNameObjUnique).Update(ctx, filter, cond)
	if err != nil {
		blog.Errorf("update cc_ObjectUnique document failed, filter:%#v, cond:%#v, err:%v", filter, cond, err)
		return err
	}
	return nil
}
