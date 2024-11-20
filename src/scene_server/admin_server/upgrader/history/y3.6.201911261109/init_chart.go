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
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func initInnerChart(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idArr := make([]uint64, 0)
	idArr, err := db.NextSequences(ctx, BKTableNameChartConfig, len(InnerChartsArr))
	if err != nil {
		return fmt.Errorf("get next sequences failed, tableName: %s, err: %+v", BKTableNameChartConfig, err)
	}

	for index, chart := range InnerChartsArr {
		innerChart := InnerChartsMap[chart]
		innerChart.ConfigID = idArr[index]
		innerChart.CreateTime.Time = time.Now()
		innerChart.OwnerID = conf.OwnerID
		if err := db.Table(BKTableNameChartConfig).Insert(ctx, innerChart); err != nil {
			return fmt.Errorf("insert chart config failed, tableName: %s, chart: %+v, err: %+v",
				BKTableNameChartConfig, innerChart, err)
		}
	}

	position := ChartPosition{
		BizID: 0,
		Position: PositionInfo{
			Host: idArr[2:6],
			Inst: idArr[6:],
		},
		OwnerID: "0",
	}

	if err := db.Table(BKTableNameChartPosition).Insert(ctx, position); err != nil {
		return fmt.Errorf("insert cahrt position data failed, table: %s, position: %+v, err: %s",
			BKTableNameChartPosition, position, err)
	}

	return nil
}

// Time Time
type Time struct {
	time.Time `bson:",inline" json:",inline"`
}

// ChartConfig chart config
type ChartConfig struct {
	ConfigID   uint64 `json:"config_id" bson:"config_id"`
	ReportType string `json:"report_type" bson:"report_type"`
	Name       string `json:"name" bson:"name"`
	CreateTime Time   `json:"create_time" bson:"create_time"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjID      string `json:"bk_obj_id" bson:"bk_obj_id"`
	Width      string `json:"width" bson:"width"`
	ChartType  string `json:"chart_type" bson:"chart_type"`
	Field      string `json:"field" bson:"field"`
	XAxisCount int64  `json:"x_axis_count" bson:"x_axis_count"`
}

// ChartPosition chart position
type ChartPosition struct {
	BizID    int64        `json:"bk_biz_id" bson:"bk_biz_id"`
	Position PositionInfo `json:"position" bson:"position"`
	OwnerID  string       `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// PositionInfo position info
type PositionInfo struct {
	Host []uint64 `json:"host" bson:"host"`
	Inst []uint64 `json:"inst" bson:"inst"`
}

// CloudMapping cloud mapping
type CloudMapping struct {
	CreateTime Time   `json:"create_time" bson:"create_time"`
	LastTime   Time   `json:"last_time" bson:"lsat_time"`
	CloudName  string `json:"bk_cloud_name" bson:"bk_cloud_name"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	CloudID    int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
}

var (
	// BizModuleHostChart biz module host chart
	BizModuleHostChart = ChartConfig{
		ReportType: "biz_module_host_chart",
	}

	// HostOsChart HostOs Chart
	HostOsChart = ChartConfig{
		ReportType: "host_os_chart",
		Name:       "按操作系统类型统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "pie",
		Field:      "bk_os_type",
		XAxisCount: 10,
	}

	// HostBizChart Host Biz Chart
	HostBizChart = ChartConfig{
		ReportType: "host_biz_chart",
		Name:       "按业务统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// HostCloudChart Host Cloud Chart
	HostCloudChart = ChartConfig{
		ReportType: "host_cloud_chart",
		Name:       "按管控区域统计",
		Width:      "100",
		ObjID:      "host",
		ChartType:  "bar",
		Field:      common.BKCloudIDField,
		XAxisCount: 20,
	}

	// HostChangeBizChart Host Change Biz Chart
	HostChangeBizChart = ChartConfig{
		ReportType: "host_change_biz_chart",
		Name:       "主机数量变化趋势",
		Width:      "100",
		XAxisCount: 20,
	}

	// ModelAndInstCountChart Model And Inst Count Chart
	ModelAndInstCountChart = ChartConfig{
		ReportType: "model_and_inst_count",
	}

	// ModelInstChart Model Inst Chart
	ModelInstChart = ChartConfig{
		ReportType: "model_inst_chart",
		Name:       "实例数量统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// ModelInstChangeChart Model Inst Change Chart
	ModelInstChangeChart = ChartConfig{
		ReportType: "model_inst_change_chart",
		Name:       "实例变更统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// InnerChartsMap Inner Charts Map
	InnerChartsMap = map[string]ChartConfig{
		"biz_module_host_chart":   BizModuleHostChart,
		"model_and_inst_count":    ModelAndInstCountChart,
		"host_os_chart":           HostOsChart,
		"host_biz_chart":          HostBizChart,
		"host_cloud_chart":        HostCloudChart,
		"host_change_biz_chart":   HostChangeBizChart,
		"model_inst_chart":        ModelInstChart,
		"model_inst_change_chart": ModelInstChangeChart,
	}

	// InnerChartsArr Inner Charts Arr
	InnerChartsArr = []string{
		"biz_module_host_chart",
		"model_and_inst_count",
		"host_os_chart",
		"host_biz_chart",
		"host_cloud_chart",
		"host_change_biz_chart",
		"model_inst_chart",
		"model_inst_change_chart",
	}
)
