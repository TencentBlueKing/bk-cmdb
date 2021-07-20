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

package y3_9_202002131522

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func upgradeCloudArea(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	opt := mapstr.MapStr{}
	cloudMapping := make([]metadata.CloudMapping, 0)
	if err := db.Table(common.BKTableNameBasePlat).Find(opt).All(ctx, &cloudMapping); err != nil {
		return fmt.Errorf("upgrade y3.9.202002131522, upgradeCloudArea failed because get cloud area data failed, err: %v", err)
	}

	for _, cloud := range cloudMapping {
		cloudArea := mapstr.MapStr{
			"bk_status":        "1",
			"bk_status_detail": "",
			"bk_account_id":    0,
			"last_time":        metadata.Now(),
			"bk_last_editor":   conf.User,
			"bk_cloud_vendor":  "",
			"bk_region":        "",
			"bk_vpc_id":        "",
			"bk_vpc_name":      "",
			"bk_creator":       conf.User,
		}

		cond := mapstr.MapStr{common.BKCloudIDField: cloud.CloudID}
		if err := db.Table(common.BKTableNameBasePlat).Update(ctx, cond, cloudArea); err != nil {
			return err
		}
	}

	return nil
}
