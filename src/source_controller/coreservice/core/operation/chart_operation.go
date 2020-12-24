/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operation

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

func (m *operationManager) SearchOperationChart(kit *rest.Kit, inputParam interface{}) (*metadata.ChartClassification, error) {
	opt := map[string]interface{}{}
	chartConfig := make([]metadata.ChartConfig, 0)

	if err := mongodb.Client().Table(common.BKTableNameChartConfig).Find(inputParam).All(kit.Ctx, &chartConfig); err != nil {
		blog.Errorf("SearchOperationChart fail, err: %v, rid: %v", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrOperationSearchChartFail)
	}

	chartPosition := make([]metadata.ChartPosition, 0)
	if err := mongodb.Client().Table(common.BKTableNameChartPosition).Find(opt).All(kit.Ctx, &chartPosition); err != nil {
		blog.Errorf("SearchOperationChart fail, err: %v, rid: %v", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrOperationSearchChartFail)
	}

	if len(chartConfig) == 0 || len(chartPosition) == 0 {
		return nil, nil
	}

	// 两个for循环是为了确定图表位置
	chartsInfo := &metadata.ChartClassification{}
	for _, id := range chartPosition[0].Position.Host {
		for _, chart := range chartConfig {
			if chart.ConfigID == id {
				chartsInfo.Host = append(chartsInfo.Host, chart)
				break
			}
		}
	}
	for _, id := range chartPosition[0].Position.Inst {
		for _, chart := range chartConfig {
			if chart.ConfigID == id {
				chartsInfo.Inst = append(chartsInfo.Inst, chart)
				break
			}
		}
	}

	for _, chart := range chartConfig {
		if chart.ReportType == common.BizModuleHostChart {
			chartsInfo.Nav = append(chartsInfo.Nav, chart)
		}
		if chart.ReportType == common.ModelAndInstCount {
			chartsInfo.Nav = append(chartsInfo.Nav, chart)
		}
	}

	return chartsInfo, nil
}

func (m *operationManager) CreateOperationChart(kit *rest.Kit, inputParam metadata.ChartConfig) (uint64, error) {
	configID, err := mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameChartConfig)
	if err != nil {
		return 0, err
	}
	inputParam.ConfigID = configID

	if err := mongodb.Client().Table(common.BKTableNameChartConfig).Insert(kit.Ctx, inputParam); err != nil {
		blog.Errorf("CreateOperationChart fail, err: %v, rid: %v", err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrOperationNewAddStatisticFail)
	}

	return configID, nil
}

func (m *operationManager) UpdateChartPosition(kit *rest.Kit, inputParam interface{}) (interface{}, error) {
	opt := map[string]interface{}{}

	if err := mongodb.Client().Table(common.BKTableNameChartPosition).Delete(kit.Ctx, opt); err != nil {
		blog.Errorf("UpdateChartPosition, delete chart position info fail, err: %v, rid: %v", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrOperationUpdateChartPositionFail)
	}

	if err := mongodb.Client().Table(common.BKTableNameChartPosition).Insert(kit.Ctx, inputParam); err != nil {
		blog.Errorf("UpdateChartPosition, update chart position fail, err: %v, rid: %v", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrOperationUpdateChartPositionFail)
	}

	return nil, nil
}

func (m *operationManager) DeleteOperationChart(kit *rest.Kit, id int64) (interface{}, error) {
	opt := map[string]interface{}{}
	opt[common.OperationConfigID] = id
	if err := mongodb.Client().Table(common.BKTableNameChartConfig).Delete(kit.Ctx, opt); err != nil {
		blog.Errorf("DeleteOperationChart fail, err: %v, rid: %v", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrOperationDeleteChartFail)
	}

	return nil, nil
}

func (m *operationManager) UpdateOperationChart(kit *rest.Kit, inputParam map[string]interface{}) (interface{}, error) {
	opt := map[string]interface{}{}
	opt[common.OperationConfigID] = inputParam[common.OperationConfigID]
	if err := mongodb.Client().Table(common.BKTableNameChartConfig).Update(kit.Ctx, opt, inputParam); err != nil {
		blog.Errorf("UpdateOperationChart fail,chartName: %v, id: %v err: %v, rid: %v", opt["name"], inputParam[common.OperationConfigID], err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrOperationUpdateChartFail)
	}

	return nil, nil
}
