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

type ServiceCategory struct {
	Metadata metadata.Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID   int64  `field:"id" json:"id,omitempty" bson:"id"`
	Name string `field:"name" json:"name,omitempty" bson:"name"`

	RootID          int64  `field:"bk_root_id" json:"bk_root_id,omitempty" bson:"bk_root_id"`
	ParentID        int64  `field:"bk_parent_id" json:"bk_parent_id,omitempty" bson:"bk_parent_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`

	// IsBuiltIn indicates internal system service category, which shouldn't be modified.
	IsBuiltIn bool `field:"is_built_in" json:"is_built_in" bson:"is_built_in"`
}

func addDefaultCategory(ctx context.Context, db dal.RDB, conf *upgrader.Config) (int64, error) {

	firstCategory := ServiceCategory{}
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

		firstCategory = ServiceCategory{
			ID:              int64(firstID),
			Name:            common.DefaultServiceCategoryName,
			RootID:          int64(firstID),
			ParentID:        0,
			SupplierAccount: "0",
			IsBuiltIn:       true,
			Metadata:        metadata.NewMetadata(0),
		}
		err = db.Table(common.BKTableNameServiceCategory).Insert(ctx, firstCategory)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	// insert second category
	secondCategory := ServiceCategory{}
	cond = metadata.BizLabelNotExist.Clone()
	cond.Set(common.BKFieldName, common.DefaultServiceCategoryName)
	cond.Set(common.BKParentIDField, firstCategory.ID)
	err = db.Table(common.BKTableNameServiceCategory).Find(cond).One(ctx, &secondCategory)
	if db.IsNotFoundError(err) {
		secondID, err := db.NextSequence(ctx, common.BKTableNameServiceCategory)
		if err != nil {
			return 0, err
		}

		secondCategory = ServiceCategory{
			ID:              int64(secondID),
			Name:            common.DefaultServiceCategoryName,
			RootID:          int64(firstCategory.RootID),
			ParentID:        int64(firstCategory.ID),
			SupplierAccount: "0",
			IsBuiltIn:       true,
			Metadata:        metadata.NewMetadata(0),
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
