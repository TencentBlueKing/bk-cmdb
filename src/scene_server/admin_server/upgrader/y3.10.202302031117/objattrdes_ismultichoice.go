/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package y3_10_202302031117

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// addObjAttrDesIsMultipleColumn add is multiple field to objAttrDes table.
func addObjAttrDesIsMultipleColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	err := db.Table(common.BKTableNameObjAttDes).AddColumn(ctx, common.BKIsMultipleField, false)
	if err != nil {
		blog.Errorf("add multiple column to objAttrDes failed, err: %v", err)
		return err
	}

	return nil
}

// updateUserFieldTypeObjAttrIsMultiple update the default value of user type field is multiple to true
func updateUserFieldTypeObjAttrIsMultiple(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKPropertyTypeField: common.FieldTypeUser,
			},
			{
				common.BKPropertyTypeField: common.FieldTypeOrganization,
			},
		},
	}
	doc := map[string]interface{}{
		common.BKIsMultipleField: true,
	}
	_, err := db.Table(common.BKTableNameObjAttDes).UpdateMany(ctx, filter, doc)
	if err != nil {
		blog.Errorf("update the default value of user type field is multiple to true failed, err: %v", err)
		return err
	}

	return nil
}
