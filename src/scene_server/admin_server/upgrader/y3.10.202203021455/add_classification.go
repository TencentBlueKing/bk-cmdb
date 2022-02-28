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

package y3_10_202203021455

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const (
	innerType   = "inner"
	defaultIcon = "icon-cc-default"
)

// changeModelClassification put all hidden models into the unclassified group.
func changeModelClassification(ctx context.Context, db dal.RDB) error {

	cond := mapstr.MapStr{
		metadata.ModelFieldIsHidden: true,
	}
	doc := mapstr.MapStr{
		common.BKClassificationIDField: metadata.ClassificationUncategorizedID,
	}

	if err := db.Table(common.BKTableNameObjDes).Update(ctx, cond, doc); err != nil {
		blog.Errorf("Failed to update model grouping to uncategorized, err: %v", err)
		return err
	}
	return nil
}

// changeClassificationType For this scenario, just update the type directly, without judging the original type value.
func changeClassificationType(ctx context.Context, db dal.RDB) error {

	cond := mapstr.MapStr{
		common.BKClassificationIDField: mapstr.MapStr{
			common.BKDBIN: []string{metadata.ClassificationHostManageID, metadata.ClassificationBizTopoID,
				metadata.ClassificationOrganizationID, metadata.ClassificationNetworkID},
		},
	}

	data := mapstr.MapStr{
		"bk_classification_type": "",
	}

	err := db.Table(common.BKTableNameObjClassification).Update(ctx, cond, data)
	if err != nil {
		blog.Errorf("update classification type failed, err: %v", err)
		return err
	}

	return nil
}

func addUncategorizedClassification(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	defaultClassificationRow := metadata.Classification{
		ClassificationID:   metadata.ClassificationUncategorizedID,
		ClassificationName: metadata.ClassificationUncategorized,
		ClassificationType: innerType,
		ClassificationIcon: defaultIcon,
		OwnerID:            conf.OwnerID,
	}

	// To check whether the id and name are repeated, search by "bk_classification_id" or "bk_classification_name". It
	// is necessary to judge whether it was created by the system. If the other fields "bk_classification_type" and
	// "bk_classification_icon" are consistent with the built-in definitions, it is considered that the system has been
	// initialized, and the return is successful. If it does not return an error, it needs to be handled manually.
	result := new(metadata.Classification)
	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				metadata.ClassFieldClassificationID: metadata.ClassificationUncategorizedID,
			},
			{
				metadata.ClassFieldClassificationName: metadata.ClassificationUncategorized,
			},
		},
	}

	err := db.Table(common.BKTableNameObjClassification).Find(filter).One(ctx, result)
	if err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("find obj classification failed, err: %v", err)
		return err
	}
	// there is a category item with the same "bk_classification_name" or the same "bk_classification_id"
	if result.ID != 0 {
		// if the icon and type are the same, return success directly,Otherwise, it means that it is a category item
		// created by the user, which needs to be manually processed and upgraded by the user.
		if result.ClassificationIcon == defaultIcon && result.ClassificationType == innerType &&
			result.ClassificationID == metadata.ClassificationUncategorizedID &&
			result.ClassificationName == metadata.ClassificationUncategorized {
			return nil
		}
		blog.Errorf("classification name %s or id %s already exists, user needs to deal with duplicates first to "+
			"successfully upgrade", metadata.ClassificationUncategorized, metadata.ClassificationUncategorizedID)

		return fmt.Errorf("failed to add classification %s", metadata.ClassificationUncategorizedID)
	}

	table := common.BKTableNameObjClassification
	if _, _, err := upgrader.Upsert(ctx, db, table, &defaultClassificationRow, "id",
		[]string{common.BKClassificationIDField}, []string{"id"}); err != nil {
		blog.Errorf("add default classification failed for %s table, err: %v", table, err)
		return err
	}
	return nil
}
