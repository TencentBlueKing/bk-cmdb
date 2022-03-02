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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func instanceObjectIDMapping(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	exists, err := db.HasTable(ctx, objectBaseMapping)
	if err != nil {
		blog.ErrorJSON("check table(%s) exist error, err: %s", objectBaseMapping, err)
		return err

	}
	if !exists {
		if err := db.CreateTable(ctx, objectBaseMapping); err != nil {
			blog.ErrorJSON("create table(%s) error, err: %s", objectBaseMapping, err)
			return err
		}

	}

	index := types.Index{
		Name: common.CCLogicIndexNamePrefix + "InstID",
		Keys: map[string]int32{
			common.BKInstIDField: 1,
		},
		Background: true,
		Unique:     true,
	}
	if err := db.Table(objectBaseMapping).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
		return err
	}

	return nil
}
