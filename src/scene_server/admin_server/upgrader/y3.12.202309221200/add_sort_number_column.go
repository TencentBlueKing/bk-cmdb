/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_12_202309221200

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

func addSortNumberColumnToObjDes(ctx context.Context, db dal.RDB) error {
	// 查询 obj_sort_number 字段是否已存在，存在则返回
	objSortNumberList := make([]map[string]int64, 0)
	objSortNumberFilter := map[string]interface{}{
		common.ObjSortNumberField: map[string]interface{}{
			common.BKDBExists: true,
		},
	}
	err := db.Table(common.BKTableNameObjDes).Find(objSortNumberFilter).Fields(common.ObjSortNumberField).
		All(ctx, &objSortNumberList)
	if err != nil {
		blog.Errorf("list object sort number failed, db find failed, err: %v", err)
		return err
	}
	if len(objSortNumberList) > 0 {
		return nil
	}

	// 获取所有模型信息
	objectList := make([]metadata.Object, 0)
	err = db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKFieldID, common.BKClassificationIDField).
		All(ctx, &objectList)
	if err != nil {
		blog.Errorf("list object ids failed, db find failed, err: %v", err)
		return err
	}

	clsMap := make(map[string][]int64)
	for _, object := range objectList {
		clsMap[object.ObjCls] = append(clsMap[object.ObjCls], object.ID)
	}

	for _, ids := range clsMap {
		if len(ids) == 0 {
			continue
		}
		for index, id := range ids {
			filter := map[string]int64{
				common.BKFieldID: id,
			}
			doc := map[string]int64{
				common.ObjSortNumberField: int64(index),
			}
			err = db.Table(common.BKTableNameObjDes).Update(ctx, filter, doc)
			if err != nil {
				blog.Errorf("add column to table failed, err: %v, column: %s, table: %s", err,
					common.ObjSortNumberField, common.BKTableNameObjDes)
				return err
			}
		}
	}
	return nil
}
