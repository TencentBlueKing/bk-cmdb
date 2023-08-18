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

package y3_11_202308161500

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
)

func addSortNumberColumnToObjDes(ctx context.Context, db dal.RDB) error {
	//查询 obj_sort_number 字段是否已存在，存在则返回
	objSortNumberList := make([]map[string]int64, 0)
	objSortNumberFilter := map[string]interface{}{
		common.ObjSortNumberField: map[string]interface{}{
			common.BKDBExists: true,
		},
	}
	err := db.Table(common.BKTableNameObjDes).Find(objSortNumberFilter).
		Fields(common.ObjSortNumberField).All(ctx, &objSortNumberList)
	if err != nil {
		blog.Errorf("list object sort number failed, db find failed, err: %s", err.Error())
		return err
	}
	if len(objSortNumberList) > 0 {
		return nil
	}

	//不存在则新建字段 obj_sort_number
	if err := db.Table(common.BKTableNameObjDes).AddColumn(ctx, common.ObjSortNumberField, 0); err != nil {
		blog.Errorf("add %s column to table %s failed, err: %v",
			common.ObjSortNumberField, common.BKTableNameObjDes, err)
		return err
	}

	//获取所有模型分组信息
	classificationIds := make([]map[string]string, 0)
	err = db.Table(common.BKTableNameObjClassification).Find(nil).
		Fields(common.BKClassificationIDField).All(ctx, &classificationIds)
	if err != nil {
		blog.Errorf("list classification ids failed, db find failed, err: %s", err.Error())
		return err
	}

	for _, classification := range classificationIds {
		//获取分组下模型信息
		objectList := make([]map[string]int64, 0)
		objectFilter := map[string]interface{}{
			common.BKClassificationIDField: map[string]interface{}{
				common.BKDBEQ: classification[common.BKClassificationIDField],
			},
		}
		err := db.Table(common.BKTableNameObjDes).Find(objectFilter).
			Fields(common.BKFieldID).All(ctx, &objectList)
		if err != nil {
			blog.Errorf("list object ids failed, db find failed, err: %s", err.Error())
			return err
		}
		if len(objectList) == 0 {
			continue
		}

		//更新字段值
		for index, objectId := range objectList {
			updateField := map[string]interface{}{
				common.BKFieldID: objectId[common.BKFieldID],
			}
			data := map[string]interface{}{
				common.ObjSortNumberField: index,
			}

			err := db.Table(common.BKTableNameObjDes).Update(ctx, updateField, data)
			if err != nil {
				blog.Errorf("update sort number failed, objectId: d%, err: %s", objectId[common.BKFieldID], err.Error())
				return err
			}
		}
	}
	return nil
}
