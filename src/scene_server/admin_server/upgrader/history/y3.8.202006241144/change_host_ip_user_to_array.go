/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_8_202006241144

import (
	"context"
	"fmt"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// changeHostIPAndUserToArray TODO
// change host inner ip and outer ip and operator and bak operator value from string split by comma to array
func changeHostIPAndUserToArray(ctx context.Context, db dal.RDB, conf *history.Config) error {
	count, err := db.Table(common.BKTableNameBaseHost).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("count hosts failed, err: %s", err.Error())
		return err
	}
	mongo, ok := db.(*local.OldMongo)
	if !ok {
		return fmt.Errorf("db is not *local.OldMongo type")
	}
	dbc := mongo.GetDBClient()
	needChangeFields := []string{common.BKHostInnerIPField, common.BKHostOuterIPField, common.BKOperatorField,
		common.BKBakOperatorField}
	for i := uint64(0); i < count; i += common.BKMaxPageSize {
		hosts := make([]map[string]interface{}, 0)
		findOpts := &options.FindOptions{}
		findOpts.SetSkip(int64(i))
		findOpts.SetLimit(common.BKMaxPageSize)
		cursor, err := dbc.Database(mongo.GetDBName()).Collection(common.BKTableNameBaseHost).Find(ctx, bson.M{},
			findOpts)
		if err != nil {
			blog.Errorf("find hosts starting from %d failed, err: %s", i, err.Error())
			return err
		}
		if err := cursor.All(ctx, &hosts); err != nil {
			blog.Errorf("decode hosts failed, err: %s", i, err.Error())
			return err
		}
		for _, host := range hosts {
			doc := make(map[string]interface{})
			for _, field := range needChangeFields {
				if host[field] == nil {
					doc[field] = make([]string, 0)
					continue
				}
				if value, ok := host[field].(string); ok {
					if len(value) == 0 {
						doc[field] = make([]string, 0)
					} else {
						doc[field] = strings.Split(value, ",")
					}
				}
			}
			if len(doc) == 0 {
				continue
			}
			filter := map[string]interface{}{
				common.BKHostIDField: host[common.BKHostIDField],
			}
			if err := db.Table(common.BKTableNameBaseHost).Update(ctx, filter, doc); err != nil {
				blog.ErrorJSON("update host ip to array failed, filter: %s, doc: %s, err: %s", filter, doc, err)
				return err
			}
		}
	}
	return nil
}
