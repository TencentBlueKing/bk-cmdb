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

package y3_9_202010151455

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// addBizIDOnCommonAttr add bk_biz_id field with its value 0 to common attributes whose bk_biz_id field is not exist
func addBizIDOnCommonAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	tableName := common.BKTableNameObjAttDes

	filter := map[string]interface{}{
		"bk_biz_id": map[string]interface{}{
			"$exists": 0,
		},
	}

	doc := map[string]int64{"bk_biz_id": 0}

	if err := db.Table(tableName).Update(ctx, filter, doc); err != nil {
		blog.ErrorJSON("addBizIDOnCommonAttr update failed, filter: %s, doc: %s, err: %s", filter, doc, err)
		return err
	}

	return nil
}
