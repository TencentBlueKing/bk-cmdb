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

package y3_13_202402191140

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func migrateInnerClassification(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := mapstr.MapStr{
		common.BKClassificationIDField: mapstr.MapStr{
			common.BKDBIN: []string{metadata.ClassificationHostManageID, metadata.ClassificationBizTopoID,
				metadata.ClassificationOrganizationID, metadata.ClassificationNetworkID},
		},
	}

	updateData := mapstr.MapStr{
		common.BKClassificationTypeField: "",
		common.BkSupplierAccount:         conf.OwnerID,
	}

	if err := db.Table(common.BKTableNameObjClassification).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("migrate inner classification failed, err: %v, cond: %+v, data: %+v", err, cond, updateData)
		return err
	}

	return nil
}
