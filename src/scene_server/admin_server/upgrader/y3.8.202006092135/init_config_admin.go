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

package y3_8_202006092135

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func initConfigAdmin(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	cnt, err := db.Table(common.BKTableNameSystem).Find(cond).Count(ctx)
	if err != nil {
		blog.ErrorJSON("insert failed, find err:%s, cond:%s", "", err, cond)
		return err
	}

	if cnt == 0 {
		doc := map[string]interface{}{
			"_id":                  common.ConfigAdminID,
			common.CreateTimeField: time.Now(),
			common.LastTimeField:   time.Now(),
		}
		if err := db.Table(common.BKTableNameSystem).Insert(ctx, doc); err != nil {
			blog.ErrorJSON("insert failed, insert err:%s, doc: %s", err, doc)
			return err
		}
	}

	return upgrader.UpgradeConfigAdmin(ctx, db)
}
