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
	"configcenter/src/common/metadata"
	"context"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// removeResourceFaultRecycle 资源池目录不需要"故障机"和"待回收"
func removeResourceFaultRecycle(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	bizInfo := &metadata.BizInst{}
	cond := mapstr.MapStr{common.BKDefaultField: 1}
	if err := db.Table(common.BKTableNameBaseApp).Find(cond).One(ctx, bizInfo); err != nil {
		return err
	}

	shouldRemoveModules := []int64{2, 3}
	filter := mapstr.MapStr{
		common.BKAppIDField:   bizInfo.BizID,
		common.BKDefaultField: mapstr.MapStr{common.BKDBIN: shouldRemoveModules},
	}
	if err := db.Table(common.BKTableNameBaseModule).Delete(ctx, filter); err != nil {
		return err
	}

	return nil
}
