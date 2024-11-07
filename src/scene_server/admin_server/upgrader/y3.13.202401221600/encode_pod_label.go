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

package y3_13_202401221600

import (
	"context"
	"strings"

	"configcenter/pkg/conv"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal"
)

func encodePodLabel(ctx context.Context, db dal.RDB) error {
	start := uint64(0)
	for {
		pods := make([]types.Pod, 0)
		if err := db.Table(types.BKTableNameBasePod).Find(nil).Fields(common.BKFieldID, types.LabelsField).
			Start(start).Limit(common.BKMaxPageSize).Sort(common.BKFieldID).All(ctx, &pods); err != nil {
			blog.Errorf("list pod failed, start: %d, err: %v", start, err)
			return err
		}

		if len(pods) == 0 {
			break
		}

		for _, pod := range pods {
			if pod.Labels == nil {
				continue
			}

			hasDot := false
			newLabels := make(map[string]string)
			for key, val := range *pod.Labels {
				if strings.Contains(key, ".") {
					key = conv.EncodeDot(key)
					hasDot = true
				}

				newLabels[key] = val
			}

			if !hasDot {
				continue
			}

			cond := map[string]interface{}{common.BKFieldID: pod.ID}
			updateData := map[string]interface{}{types.LabelsField: newLabels}
			if err := db.Table(types.BKTableNameBasePod).Update(ctx, cond, updateData); err != nil {
				blog.Errorf("update pod failed, cond: %+v, data: %+v, err: %v", cond, updateData, err)
				return err
			}
		}

		start += common.BKMaxPageSize
	}

	return nil
}
