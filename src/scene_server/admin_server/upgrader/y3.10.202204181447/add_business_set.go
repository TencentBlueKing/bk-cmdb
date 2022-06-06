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

package y3_10_202204181447

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const (
	bizSetID         = 9991001
	bizSetName       = "BlueKing"
	bizSetMaintainer = "admin"
)

func addDefaultBusinessSet(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	now := metadata.Time{
		Time: time.Now(),
	}
	defaultBusinessSet := &metadata.BizSetInst{
		BizSetID:         bizSetID,
		BizSetName:       bizSetName,
		BizSetMaintainer: bizSetMaintainer,
		SupplierAccount:  conf.OwnerID,
		CreateTime:       now,
		LastTime:         now,
		Default:          common.DefaultResBusinessSetFlag,
		Scope: metadata.BizSetScope{
			MatchAll: true,
		},
	}

	result := make([]metadata.BizSetInst, 0)
	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKBizSetIDField: bizSetID,
			},
			{
				common.BKBizSetNameField: bizSetName,
			},
		},
	}

	err := db.Table(common.BKTableNameBaseBizSet).Find(filter).All(ctx, &result)
	if err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("find business set failed, err: %v", err)
		return err
	}

	if len(result) >= 2 {
		blog.Errorf("business set id %s and name %s already exists, the number of business set that meet the "+
			"condition is greater than 2, val: %v", bizSetID, bizSetName, result)
		return fmt.Errorf("failed to add default business set, business set id: %d or name: %s already exists, the "+
			"number of business set that meet the condition is greater than 2", bizSetID, bizSetName)
	}

	if len(result) == 1 {
		if result[0].BizSetID == bizSetID && result[0].BizSetName == bizSetName && result[0].Scope.MatchAll &&
			result[0].SupplierAccount == conf.OwnerID {

			if result[0].Default != common.DefaultResBusinessSetFlag {
				data := map[string]interface{}{
					common.BKDefaultField: common.DefaultResBusinessSetFlag,
				}

				if err := db.Table(common.BKTableNameBaseBizSet).Update(ctx, filter, data); err != nil {
					blog.Errorf("update business set default val failed, err: %v", err)
					return err
				}
			}

			return nil
		}

		blog.Errorf("business set id %s or name %s already exists, user needs to deal with duplicates first to "+
			"successfully upgrade", bizSetID, bizSetName)
		return fmt.Errorf("failed to add default business set, business set id: %d or name: %s already exists, user "+
			"needs to deal with duplicates first to successfully upgrade", bizSetID, bizSetName)
	}

	if err := db.Table(common.BKTableNameBaseBizSet).Insert(ctx, defaultBusinessSet); err != nil {
		blog.Errorf("insert default business set failed, val: %v, err: %v", defaultBusinessSet, err)
		return err
	}

	return nil
}
