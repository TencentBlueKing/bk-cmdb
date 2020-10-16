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

package y3_9_202010151650

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func addCloudIDIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	index := types.Index{
		Keys:       map[string]int32{"bk_cloud_id": 1},
		Name:       "bk_cloud_id_1",
		Unique:     false,
		Background: true,
	}

	tableNames := []string{
		common.BKTableNameBaseHost,
		common.BKTableNameBasePlat,
	}

	for _, tableName := range tableNames {
		err := db.Table(tableName).CreateIndex(ctx, index)
		if err != nil {
			blog.ErrorJSON("add index %s for table %s failed, err:%s", index, tableName, err)
			return err
		}
	}

	return nil
}
