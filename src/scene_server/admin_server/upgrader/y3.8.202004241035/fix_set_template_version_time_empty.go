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

package y3_8_202004241035

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// fixSetTemplateVersionTimeEmpty fix time type field value empty for set template version attribute
func fixSetTemplateVersionTimeEmpty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKPropertyIDField: common.BKSetTemplateVersionField,
		common.BKOwnerIDField:    conf.OwnerID,
		common.BKObjIDField:      common.BKInnerObjIDSet,
	}
	now := metadata.Now()
	doc := map[string]interface{}{
		common.CreateTimeField: &now,
		common.LastTimeField:   &now,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.ErrorJSON("update set template version attribute failed, filter: %s, doc: %s, err: %s", filter, doc, err)
		return err
	}
	return nil
}
