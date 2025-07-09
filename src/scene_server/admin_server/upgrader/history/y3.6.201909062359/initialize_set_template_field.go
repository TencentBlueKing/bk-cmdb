/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_6_201909062359

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func initializeSetTemplateField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKSetTemplateIDField: map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	doc := map[string]interface{}{
		common.BKSetTemplateIDField: 0,
	}
	if err := db.Table(common.BKTableNameBaseSet).Update(ctx, filter, doc); err != nil {
		blog.Errorf("[upgrade v19.08.24.01] initializeSetTemplateField set_template_id on set failed, err: %+v", err)
		return err
	}
	if err := db.Table(common.BKTableNameBaseModule).Update(ctx, filter, doc); err != nil {
		blog.Errorf("[upgrade v19.08.24.01] initializeSetTemplateField set_template_id on module failed, err: %+v", err)
		return err
	}

	return nil
}
