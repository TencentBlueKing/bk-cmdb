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

package y3_10_202305151505

import (
	"configcenter/src/scene_server/admin_server/upgrader/history"

	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
)

func updateBKCloudWord(ctx context.Context, db dal.RDB) error {

	err := db.Table(common.BKTableNameObjAttDes).
		Update(ctx,
			map[string]interface{}{"creator": common.CCSystemOperatorUserName, common.BKPropertyNameField: "云区域"},
			map[string]interface{}{common.BKPropertyNameField: "管控区域"})
	if err != nil {
		blog.Errorf("update bk cloud word failed, err: %v", err)
		return err
	}

	return nil
}

func updateDefaultArea(ctx context.Context, db dal.RDB, conf *history.Config) error {

	cond := map[string]interface{}{
		common.BKCloudNameField: "default area",
	}

	data := map[string]interface{}{
		common.BKCloudNameField: "Default Area",
	}
	if err := db.Table(common.BKTableNameBasePlat).Update(ctx, cond, data); err != nil {
		blog.Errorf("update default area failed, err: %v", err)
		return err
	}

	return nil
}
