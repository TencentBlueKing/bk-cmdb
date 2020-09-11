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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func initInnerChart(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idArr := make([]uint64, 0)
	idArr, err := db.NextSequences(ctx, common.BKTableNameChartConfig, len(metadata.InnerChartsArr))
	if err != nil {
		return fmt.Errorf("get next sequences failed, tableName: %s, err: %+v", common.BKTableNameChartConfig, err)
	}

	for index, chart := range metadata.InnerChartsArr {
		innerChart := metadata.InnerChartsMap[chart]
		innerChart.ConfigID = idArr[index]
		innerChart.CreateTime.Time = time.Now()
		innerChart.OwnerID = conf.OwnerID
		if err := db.Table(common.BKTableNameChartConfig).Insert(ctx, innerChart); err != nil {
			return fmt.Errorf("insert chart config failed, tableName: %s, chart: %+v, err: %+v", common.BKTableNameChartConfig, innerChart, err)
		}
	}

	position := metadata.ChartPosition{
		BizID: 0,
		Position: metadata.PositionInfo{
			Host: idArr[2:6],
			Inst: idArr[6:],
		},
		OwnerID: "0",
	}

	if err := db.Table(common.BKTableNameChartPosition).Insert(ctx, position); err != nil {
		return fmt.Errorf("insert cahrt position data failed, table: %s, position: %+v, err: %s", common.BKTableNameChartPosition, position, err)
	}

	return nil
}
