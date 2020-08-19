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

package y3_9_202008101530

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
)

// changeMetadataToBizID transfer field metadata.label.bk_biz_id to bk_biz_id and delete unnecessary metadata
func changeMetadataToBizID(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	mongo, ok := db.(*local.Mongo)
	if !ok {
		return fmt.Errorf("db is not *local.Mongo type")
	}
	dbc := mongo.GetDBClient()

	type ret struct {
		ID       interface{}   `bson:"_id" json:"_id"`
		Metadata mapstr.MapStr `bson:"metadata" json:"metadata"`
	}
	fields := []string{"_id", "metadata"}
	existsMetadataFilter := map[string]interface{}{
		"metadata": map[string]interface{}{
			common.BKDBExists: true,
		},
	}

	for _, tableName := range common.AllTables {
		count, err := db.Table(tableName).Find(existsMetadataFilter).Count(ctx)
		if err != nil {
			blog.Errorf("count table %s failed, err: %s", tableName, err.Error())
			return err
		}

		for i := uint64(0); i < count; i += common.BKMaxPageSize {
			result := make([]ret, 0)
			err = db.Table(tableName).Find(existsMetadataFilter).Sort("_id").Fields(fields...).Start(uint64(i)).Limit(uint64(common.BKMaxPageSize)).All(ctx, &result)
			if err != nil {
				blog.Errorf("changeMetadataToBizID starting from %d failed, err: %s", i, err)
				return err
			}

			delMeta := make([]interface{}, 0)
			delLabel := make([]interface{}, 0)
			updateMap := make(map[interface{}]int64)
			for _, r := range result {
				if label, ok := r.Metadata.Get("label"); ok {
					if labelMap, ok := label.(map[string]interface{}); ok {
						if bizIDStr, ok := labelMap["bk_biz_id"].(string); ok {
							bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
							if err != nil {
								blog.Errorf("the bk_biz_id is not int, it's value is %#v", labelMap["bk_biz_id"])
								continue
							}
							updateMap[r.ID] = bizID

						}
					}
					if len(r.Metadata) >= 2 {
						delLabel = append(delLabel, r.ID)
					} else {
						delMeta = append(delMeta, r.ID)
					}
				}
				if len(r.Metadata) == 0 {
					delMeta = append(delMeta, r.ID)
				}
			}

			// transfer field metadata.label.bk_biz_id to bk_biz_id
			for ID, bizID := range updateMap {

				filter := map[string]interface{}{
					"_id": ID,
				}
				doc := map[string]int64{"bk_biz_id": bizID}

				if err := db.Table(tableName).Update(ctx, filter, doc); err != nil {
					blog.ErrorJSON("update bizID failed, filter: %s, doc: %s, err: %s", filter, doc, err)
					return err
				}
			}

			// delete metadata
			delMetaFilter := map[string]interface{}{
				"_id": map[string]interface{}{
					"$in": delMeta,
				},
			}
			delMetaDoc := map[string]interface{}{
				"$unset": map[string]interface{}{
					"metadata": "",
				},
			}

			if _, err := dbc.Database(mongo.GetDBName()).Collection(tableName).UpdateMany(ctx, delMetaFilter, delMetaDoc); err != nil {
				blog.ErrorJSON("del metadata failed, filter: %s, doc: %s, err: %s", delMetaFilter, delMetaDoc, err)
				return err
			}

			// delete metadata.label
			delLabelFilter := map[string]interface{}{
				"_id": map[string]interface{}{
					"$in": delLabel,
				},
			}
			delLabelDoc := map[string]interface{}{
				"$unset": map[string]interface{}{
					"metadata.label": "",
				},
			}

			if _, err := dbc.Database(mongo.GetDBName()).Collection(tableName).UpdateMany(ctx, delLabelFilter, delLabelDoc); err != nil {
				blog.ErrorJSON("del metadata.label failed, filter: %s, doc: %s, err: %s", delLabelFilter, delLabelDoc, err)
				return err
			}
		}
	}

	return nil
}
