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

package x19_05_16_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addDefaultCategory(ctx context.Context, db dal.RDB, conf *upgrader.Config) (int64, error) {

	firstCategory := metadata.ServiceCategory{}
	// insert first category
	cond := metadata.BizLabelNotExist.Clone()
	cond.Set(common.BKFieldName, common.DefaultServiceCategoryName)
	cond.Set(common.BKParentIDField, 0)
	err := db.Table(common.BKTableNameServiceCategory).Find(cond).One(ctx, &firstCategory)
	if db.IsNotFoundError(err) {
		firstID, err := db.NextSequence(ctx, common.BKTableNameServiceCategory)
		if err != nil {
			return 0, err
		}

		firstCategory = metadata.ServiceCategory{
			ID:              int64(firstID),
			Name:            common.DefaultServiceCategoryName,
			RootID:          int64(firstID),
			ParentID:        0,
			SupplierAccount: "0",
			IsBuiltIn:       true,
		}
		err = db.Table(common.BKTableNameServiceCategory).Insert(ctx, firstCategory)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	// insert second category
	secondCategory := metadata.ServiceCategory{}
	cond = metadata.BizLabelNotExist.Clone()
	cond.Set(common.BKFieldName, common.DefaultServiceCategoryName)
	cond.Set(common.BKParentIDField, firstCategory.ID)
	err = db.Table(common.BKTableNameServiceCategory).Find(cond).One(ctx, &secondCategory)
	if db.IsNotFoundError(err) {
		secondID, err := db.NextSequence(ctx, common.BKTableNameServiceCategory)
		if err != nil {
			return 0, err
		}

		secondCategory = metadata.ServiceCategory{
			ID:              int64(secondID),
			Name:            common.DefaultServiceCategoryName,
			RootID:          int64(firstCategory.ID),
			ParentID:        int64(firstCategory.ID),
			SupplierAccount: "0",
			IsBuiltIn:       true,
		}
		err = db.Table(common.BKTableNameServiceCategory).Insert(ctx, secondCategory)
		if err != nil {
			return 0, err
		}

	} else if err != nil {
		return 0, err
	}

	return secondCategory.ID, nil
}
