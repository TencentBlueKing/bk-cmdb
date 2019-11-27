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

package y3_6_201911261109

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func initInnerChart(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idArr := make([]uint64, 0)
	for _, chart := range metadata.InnerChartsArr {
		configID, err := db.NextSequence(ctx, common.BKTableNameCloudTask)
		idArr = append(idArr, configID)
		if err != nil {
			return err
		}
		innerChart := metadata.InnerChartsMap[chart]
		innerChart.ConfigID = configID
		innerChart.CreateTime.Time = time.Now()
		innerChart.OwnerID = conf.OwnerID
		if err := db.Table(common.BKTableNameChartConfig).Insert(ctx, innerChart); err != nil {
			return err
		}
	}

	position := metadata.ChartPosition{}
	position.Position.Host = idArr[2:6]
	position.Position.Inst = idArr[6:]
	position.OwnerID = "0"

	if err := db.Table(common.BKTableNameChartPosition).Insert(ctx, position); err != nil {
		return err
	}

	return nil
}
