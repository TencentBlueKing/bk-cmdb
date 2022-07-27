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

package y3_10_202104221702

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func dropPlatVpcIDIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		"bk_obj_id":      "plat",
		"bk_property_id": "bk_vpc_id",
	}

	attrIDs, err := db.Table("cc_ObjAttDes").Distinct(ctx, "id", filter)
	if err != nil {
		blog.ErrorJSON("find cloud bk_vpc_id attribute id error. err: %s, filter: %s", err, filter)
		return err
	}

	if len(attrIDs) > 0 {
		deleteObjectUniqueFilter := map[string]interface{}{
			"keys": map[string]interface{}{
				"key_kind": "property",
				"key_id":   attrIDs[0],
			},
		}
		if err := db.Table("cc_ObjectUnique").Delete(ctx, deleteObjectUniqueFilter); err != nil {
			return err
		}

	}

	indexes, err := db.Table("cc_PlatBase").Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("find cloud indexes error. err: %s ", err)
		return err
	}
	for _, index := range indexes {
		if len(index.Keys) == 1 && index.Unique {
			if _, exist := index.Keys.Map()["bk_vpc_id"]; exist {
				if err := db.Table("cc_PlatBase").DropIndex(ctx, index.Name); err != nil {
					blog.ErrorJSON("delete cloud  bk_vpc_id index error. err: %s ", err)
					return err
				}
			}
		}
	}

	return nil
}
