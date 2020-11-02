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

package y3_9_202010281615

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// updateInstTimeVal update the value of the instance time type
func updateInstTimeVal(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	objIDArray := make([]map[string]string, 0)
	if err := db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField).All(ctx, &objIDArray); err != nil {
		blog.ErrorJSON("find model %s field failed, err: %s", common.BKObjIDField, err)
		return err
	}

	for _, objIDmap := range objIDArray {
		objID := objIDmap[common.BKObjIDField]
		// find attributes of model time type
		propertyIDArray := make([]map[string]string, 0)
		filter := mapstr.MapStr{
			common.BKObjIDField: objID,
			common.BKPropertyTypeField: common.FieldTypeTime,
		}
		if err := db.Table(common.BKTableNameObjAttDes).Find(filter).Fields(common.BKPropertyIDField).All(ctx, &propertyIDArray); err != nil {
			blog.ErrorJSON("find object attribute field failed, filter: %s, err: %s", filter, err)
			return err
		}

		var timeTypeAttr []string
		var isTimeTypeAttrExist []interface{}
		for _, propertyID := range propertyIDArray {
			timeTypeAttr = append(timeTypeAttr, propertyID[common.BKPropertyIDField])
			isTimeTypeAttrExist = append(isTimeTypeAttrExist, mapstr.MapStr{propertyID[common.BKPropertyIDField]: mapstr.MapStr{"$exists": true}})
		}

		// start to find model instances
		instTable := common.GetInstTableName(objID)
		filter = mapstr.MapStr{}

		if isTimeTypeAttrExist == nil {
			continue
		}

		filter.Set(common.BKDBOR, isTimeTypeAttrExist)

		if instTable == common.BKTableNameBaseInst {
			filter.Set(common.BKObjIDField, objID)
		}

		count, err := db.Table(instTable).Find(filter).Count(ctx)
		if err != nil {
			blog.Errorf("count table %s failed, err: %s", instTable, err.Error())
			return err
		}

		instIDField := common.GetInstIDField(objID)
		instFields := make([]string, 0)
		instFields = append(instFields, instIDField)
		instFields = append(instFields, timeTypeAttr...)

		cursor := db.Table(instTable).Find(filter).Fields(instFields...)
		for i := uint64(0); i < count; i += 2000 {
			insts := make([]map[string]interface{}, 0)
			if err := cursor.Start(i).Limit(2000).All(ctx, &insts); err != nil {
				blog.Errorf("search inst failed, err: %s", i, err.Error())
				return err
			}
			for _, inst := range insts {
				doc := make(map[string]interface{})
				// if it is a string value, convert it to time
				for field, val := range inst {
					if val == nil || field == instIDField {
						continue
					}
					valStr, ok := val.(string)
					if ok == false {
						continue
					}
					if util.IsTime(valStr) {
						doc[field] = util.Str2Time(valStr)
						continue
					}
					blog.ErrorJSON("It is not a time type string, table: %s, filed: %s, val: %s", instTable, field, val)
				}

				if len(doc) == 0 {
					continue
				}

				filter := map[string]interface{}{
					instIDField : inst[instIDField],
				}

				if err := db.Table(instTable).Update(ctx, filter, doc); err != nil {
					blog.ErrorJSON("update the value of the instance time type failed, " +
						"table: %s, filter: %s, doc: %s, err: %s", instTable, filter, doc, err)
					return err
				}
			}
		}
	}

	return nil
}
