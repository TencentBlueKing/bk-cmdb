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

package x19_05_22_01

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"context"
	"time"
)

func initInnerChart(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idArr := make([]uint64, 0)
	for _, chart := range InnerChartsArr {
		configID, err := db.NextSequence(ctx, common.BKTableNameCloudTask)
		idArr = append(idArr, configID)
		if err != nil {
			return err
		}
		innerChart := InnerChartsMap[chart]
		innerChart.ConfigID = configID
		innerChart.CreateTime = time.Now()
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

var (
	BizModuleHostChart = metadata.ChartConfig{
		ReportType: common.BizModuleHostChart,
	}

	HostOsChart = metadata.ChartConfig{
		ReportType: common.HostOSChart,
		Name:       "按操作系统类型统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "pie",
		Field:      "bk_os_type",
		XAxisCount: 10,
	}

	HostBizChart = metadata.ChartConfig{
		ReportType: common.HostBizChart,
		Name:       "按业务统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	HostCloudChart = metadata.ChartConfig{
		ReportType: common.HostCloudChart,
		Name:       "按云区域统计",
		Width:      "100",
		ObjID:      "host",
		ChartType:  "bar",
		Field:      common.BKCloudIDField,
		XAxisCount: 20,
	}

	HostChangeBizChart = metadata.ChartConfig{
		ReportType: common.HostChangeBizChart,
		Name:       "主机数量变化趋势",
		Width:      "100",
		XAxisCount: 20,
	}

	ModelAndInstCountChart = metadata.ChartConfig{
		ReportType: common.ModelAndInstCount,
	}

	ModelInstChart = metadata.ChartConfig{
		ReportType: common.ModelInstChart,
		Name:       "实例数量统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	ModelInstChangeChart = metadata.ChartConfig{
		ReportType: common.ModelInstChangeChart,
		Name:       "实例变更统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	InnerChartsMap = map[string]metadata.ChartConfig{
		common.BizModuleHostChart:   BizModuleHostChart,
		common.ModelAndInstCount:    ModelAndInstCountChart,
		common.HostOSChart:          HostOsChart,
		common.HostBizChart:         HostBizChart,
		common.HostCloudChart:       HostCloudChart,
		common.HostChangeBizChart:   HostChangeBizChart,
		common.ModelInstChart:       ModelInstChart,
		common.ModelInstChangeChart: ModelInstChangeChart,
	}

	InnerChartsArr = []string{
		common.BizModuleHostChart,
		common.ModelAndInstCount,
		common.HostOSChart,
		common.HostBizChart,
		common.HostCloudChart,
		common.HostChangeBizChart,
		common.ModelInstChart,
		common.ModelInstChangeChart,
	}
)
