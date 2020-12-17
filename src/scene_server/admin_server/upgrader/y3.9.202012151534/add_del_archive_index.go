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

package y3_9_202012151534

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

var oidCollIndex = types.Index{
	Keys:       map[string]int32{"oid": 1, "coll": 1},
	Unique:     true,
	Background: true,
	Name:       "idx_oid_coll",
}

// addDelArchiveIndex add unique index for coll and oid
func addDelArchiveIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	existIndexes, err := db.Table(common.BKTableNameDelArchive).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("find indexes for del archive table failed. err: %v", err)
		return err
	}

	for _, index := range existIndexes {
		if index.Name == oidCollIndex.Name {
			return nil
		}
	}

	err = db.Table(common.BKTableNameDelArchive).CreateIndex(ctx, oidCollIndex)
	if err != nil {
		blog.ErrorJSON("add index %s for del archive table failed, err: %s", oidCollIndex, err)
		return err
	}

	return nil
}
