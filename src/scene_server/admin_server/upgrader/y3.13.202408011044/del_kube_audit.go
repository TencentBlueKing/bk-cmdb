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

package y3_13_202408011044

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func deleteKubeAudit(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := mapstr.MapStr{
		common.BKAuditTypeField: metadata.KubeType,
	}

	for {
		logs := make([]metadata.AuditLog, 0)
		err := db.Table(common.BKTableNameAuditLog).Find(filter).Fields(common.BKFieldID).Limit(common.BKMaxPageSize).
			All(ctx, &logs)
		if err != nil {
			blog.Errorf("get kube audit logs failed, err: %v", err)
			return err
		}

		if len(logs) == 0 {
			return nil
		}

		delIDs := make([]int64, len(logs))
		for i, log := range logs {
			delIDs[i] = log.ID
		}

		delCond := map[string]interface{}{
			common.BKFieldID: map[string]interface{}{common.BKDBIN: delIDs},
		}
		if err := db.Table(common.BKTableNameAuditLog).Delete(ctx, delCond); err != nil {
			blog.Errorf("delete kube audit logs failed, err: %v, cond: %+v", err, delCond)
			return err
		}

		if len(logs) < common.BKMaxPageSize {
			break
		}

		time.Sleep(time.Millisecond * 5)
	}

	return nil
}
